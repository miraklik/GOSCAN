package main

import (
	"flag"
	"fmt"
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
	var flags utils.Flags
	flags.InitializeFlags()

	if flags.Host == "" {
		fmt.Println("Error: host is required")
		flag.Usage()
		os.Exit(1)
	}

	var sc scanner.Scanner
	switch flags.Protocol {
	case "tcp":
		sc = scanner.NewTCPScanner(flags.Timeout)
	case "udp":
		sc = scanner.NewUDPScanner(flags.Timeout)
	default:
		fmt.Printf("Error: unsupported protocol '%s'\n", flags.Protocol)
		os.Exit(1)
	}

	ports := utils.ParsePorts(flags.PortRange)
	if len(ports) == 0 {
		fmt.Println("Error: invalid port range")
		os.Exit(1)
	}

	rateLimiter := rate.NewLimiter(rate.Limit(flags.NumWorkers), flags.NumWorkers)

	fmt.Printf("[*] Starting %s scan on %s\n", flags.Protocol, flags.Host)
	fmt.Printf("[*] Port range: %s (%d ports)\n", flags.PortRange, len(ports))
	fmt.Printf("[*] Workers: %d | Timeout: %s\n", flags.NumWorkers, flags.Timeout)
	fmt.Println(utils.Separator())

	start := time.Now()

	portsChan := make(chan int, 1000)
	resultsChan := make(chan scanner.Result, 100)
	var wg sync.WaitGroup

	var scannedPorts int32
	var openPorts int32

	if flags.ShowProgress {
		go utils.ProgressMonitor(&scannedPorts, &openPorts, len(ports))
	}

	for i := 0; i < flags.NumWorkers; i++ {
		wg.Add(1)
		go utils.Worker(i+1, flags.Host, sc, rateLimiter, portsChan, resultsChan, &wg, &scannedPorts, flags.Verbose)
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
			if !flags.ShowProgress {
				utils.PrintResult(res)
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

	if flags.ShowProgress {
		fmt.Print("\r" + utils.Separator() + "\n")
	}

	if len(allResults) == 0 {
		fmt.Println("[*] No open ports found")
	} else {
		fmt.Printf("[+] Found %d open port(s)\n", len(allResults))

		if flags.ShowProgress {
			fmt.Println()
			for _, res := range allResults {
				utils.PrintResult(res)
			}
		}

		filename := saveResults(allResults, flags.OutputFormat, flags.OutputFile)
		if filename != "" {
			fmt.Printf("[+] Results saved to: %s\n", filename)
		}
	}

	utils.PrintStats(utils.Stats{
		TotalPorts:  len(ports),
		OpenPorts:   len(allResults),
		ClosedPorts: len(ports) - len(allResults),
		Duration:    duration,
	})
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
