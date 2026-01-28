package main

import (
	"context"
	"flag"
	"fmt"
	"goscan/detector"
	"goscan/output"
	"goscan/scanner"
	"goscan/utils"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"golang.org/x/time/rate"
)

func main() {
	host := flag.String("host", "127.0.0.1", "Target host to scan")
	numWorkers := flag.Int("w", 25, "Number of concurrent workers")
	portRange := flag.String("p", "1-1024", "Port range (e.g., 80,443 or 1-1024 or 1-65535)")
	protocol := flag.String("proto", "tcp", "Protocol: tcp or udp")
	timeout := flag.Duration("t", 5*time.Second, "Connection timeout")
	outputFormat := flag.String("o", "json", "Output format: json, csv, txt")
	outputFile := flag.String("output", "scan_results", "Output file name (without extension)")
	showProgress := flag.Bool("progress", true, "Show progress indicator")
	verbose := flag.Bool("v", false, "Verbose mode (show closed ports)")
	flag.Parse()

	if *host == "" {
		fmt.Println("Error: host is required")
		flag.Usage()
		os.Exit(1)
	}

	var sc scanner.Scanner
	switch *protocol {
	case "tcp":
		sc = scanner.NewTCPScanner(*timeout)
	case "udp":
		sc = scanner.NewUDPScanner(*timeout)
	default:
		fmt.Printf("Error: unsupported protocol '%s'\n", *protocol)
		os.Exit(1)
	}

	ports := utils.ParsePorts(*portRange)
	if len(ports) == 0 {
		fmt.Println("Error: invalid port range")
		os.Exit(1)
	}

	rateLimiter := rate.NewLimiter(rate.Limit(*numWorkers), *numWorkers)

	fmt.Printf("[*] Starting %s scan on %s\n", *protocol, *host)
	fmt.Printf("[*] Port range: %s (%d ports)\n", *portRange, len(ports))
	fmt.Printf("[*] Workers: %d | Timeout: %s\n", *numWorkers, *timeout)
	fmt.Println(utils.Separator())

	start := time.Now()

	portsChan := make(chan int, 1000)
	resultsChan := make(chan scanner.Result, 100)
	var wg sync.WaitGroup

	var scannedPorts int32
	var openPorts int32

	if *showProgress {
		go progressMonitor(&scannedPorts, &openPorts, len(ports))
	}

	for i := 0; i < *numWorkers; i++ {
		wg.Add(1)
		go worker(i+1, *host, sc, rateLimiter, portsChan, resultsChan, &wg, &scannedPorts, *verbose)
	}

	var allResults []scanner.Result
	var mu sync.Mutex
	done := make(chan struct{})

	go func() {
		for res := range resultsChan {
			atomic.AddInt32(&openPorts, 1)
			mu.Lock()
			allResults = append(allResults, res)
			mu.Unlock()
			if !*showProgress {
				printResult(res)
			}
		}
		close(done)
	}()

	go func() {
		for _, p := range ports {
			portsChan <- p
		}
		close(portsChan)
	}()

	wg.Wait()
	close(resultsChan)
	<-done

	duration := time.Since(start)

	if *showProgress {
		fmt.Print("\r" + utils.Separator() + "\n")
	}

	if len(allResults) == 0 {
		fmt.Println("[*] No open ports found")
	} else {
		fmt.Printf("[+] Found %d open port(s)\n", len(allResults))

		if *showProgress {
			fmt.Println()
			for _, res := range allResults {
				printResult(res)
			}
		}

		filename := saveResults(allResults, *outputFormat, *outputFile)
		if filename != "" {
			fmt.Printf("[+] Results saved to: %s\n", filename)
		}
	}

	printStats(utils.Stats{
		TotalPorts:  len(ports),
		OpenPorts:   len(allResults),
		ClosedPorts: len(ports) - len(allResults),
		Duration:    duration,
	})
}

func worker(id int, host string, sc scanner.Scanner, limiter *rate.Limiter, ports <-chan int, results chan<- scanner.Result, wg *sync.WaitGroup, scanned *int32, verbose bool) {
	defer wg.Done()

	for port := range ports {
		if limiter != nil {
			limiter.Wait(context.Background())
		}

		addr := fmt.Sprintf("%s:%d", host, port)
		banner, err := sc.ScanPort(addr)

		atomic.AddInt32(scanned, 1)

		if err == nil {
			service := detector.DetectService(banner, port)
			results <- scanner.Result{
				Port:    port,
				Service: service,
				Banner:  banner,
				Status:  "open",
			}
		} else if verbose {
			fmt.Printf("[-] Port %5d closed\n", port)
		}
	}
}

func progressMonitor(scanned, open *int32, total int) {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for range ticker.C {
		current := atomic.LoadInt32(scanned)
		openCount := atomic.LoadInt32(open)

		if current >= int32(total) {
			break
		}

		percentage := float64(current) / float64(total) * 100
		bar := progressBar(int(percentage))

		fmt.Printf("\r[%s] %d/%d ports (%.1f%%) | %d open   ",
			bar, current, total, percentage, openCount)
	}
}

func progressBar(percentage int) string {
	width := 30
	filled := percentage * width / 100

	bar := ""
	for i := 0; i < width; i++ {
		if i < filled {
			bar += "█"
		} else {
			bar += "░"
		}
	}
	return bar
}

func printResult(r scanner.Result) {
	fmt.Printf("%s[+]%s Port %s%5d%s | %-12s | %s\n",
		utils.ColorGreen, utils.ColorReset,
		utils.ColorYellow, r.Port, utils.ColorReset,
		r.Service, r.Banner)
}

func saveResults(results []scanner.Result, format, filename string) string {
	var err error
	var fullPath string

	switch format {
	case "json":
		fullPath = filename + ".json"
		err = output.SaveAsJSON(results, fullPath)
	case "csv":
		fullPath = filename + ".csv"
		err = output.SaveAsCSV(results, fullPath)
	case "txt":
		fullPath = filename + ".txt"
		err = output.SaveAsTXT(results, fullPath)
	default:
		fmt.Printf("Error: unsupported output format '%s'\n", format)
		return ""
	}

	if err != nil {
		fmt.Printf("Error saving results: %v\n", err)
		return ""
	}

	return fullPath
}

func printStats(stats utils.Stats) {
	fmt.Println(utils.Separator())
	fmt.Println("Scan Statistics:")
	fmt.Printf("  Total Ports:    %d\n", stats.TotalPorts)
	if stats.TotalPorts > 0 {
		percentage := float64(stats.OpenPorts) / float64(stats.TotalPorts) * 100
		fmt.Printf("  Open Ports:     %d (%.2f%%)\n", stats.OpenPorts, percentage)
	} else {
		fmt.Printf("  Open Ports:     %d\n", stats.OpenPorts)
	}
	fmt.Printf("  Closed Ports:   %d\n", stats.ClosedPorts)
	fmt.Printf("  Duration:       %s\n", stats.Duration)
	if stats.Duration.Seconds() > 0 {
		fmt.Printf("  Ports/second:   %.2f\n", float64(stats.TotalPorts)/stats.Duration.Seconds())
	}
	fmt.Println(utils.Separator())
}
