package main

import (
	"flag"
	"fmt"
	"goscan/cmd"
	"goscan/output"
	"goscan/pkg/logger"
	"goscan/scanner"
	"goscan/tui"
	"goscan/utils"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"golang.org/x/time/rate"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			logger.Error("Panic recovered: %v", r)
			os.Exit(1)
		}
	}()

	cmd.Execute()

	logger.Init(cmd.Verbose)
	logger.Info("Starting goscan v1.0.0")
	logger.Debug("Verbose mode enabled")

	if cmd.TUI {
		launchTUI()
		return
	}

	if cmd.Host == "" {
		fmt.Println("Error: host is required")
		flag.Usage()
		os.Exit(1)
	}

	var sc scanner.Scanner
	switch cmd.Protocol {
	case "tcp":
		sc = scanner.NewTCPScanner(cmd.Timeout)
	case "udp":
		sc = scanner.NewUDPScanner(cmd.Timeout)
	default:
		fmt.Printf("Error: unsupported protocol '%s'\n", cmd.Protocol)
		os.Exit(1)
	}

	ports := utils.ParsePorts(cmd.PortRange)
	if len(ports) == 0 {
		fmt.Println("Error: invalid port range")
		os.Exit(1)
	}

	rateLimiter := rate.NewLimiter(rate.Limit(cmd.NumWorkers), cmd.NumWorkers)

	fmt.Printf("[*] Starting %s scan on %s\n", cmd.Protocol, cmd.Host)
	fmt.Printf("[*] Port range: %s (%d ports)\n", cmd.PortRange, len(ports))
	fmt.Printf("[*] Workers: %d | Timeout: %s\n", cmd.NumWorkers, cmd.Timeout)
	fmt.Println(utils.Separator())

	start := time.Now()

	portsChan := make(chan int, 1000)
	resultsChan := make(chan scanner.Result, 100)
	var wg sync.WaitGroup

	var scannedPorts int32
	var openPorts int32

	if cmd.ShowProgress {
		go utils.ProgressMonitor(&scannedPorts, &openPorts, len(ports))
	}

	for i := 0; i < cmd.NumWorkers; i++ {
		wg.Add(1)
		go utils.Worker(i+1, cmd.Host, sc, rateLimiter, portsChan, resultsChan, &wg, &scannedPorts, cmd.Verbose)
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
			if !cmd.ShowProgress {
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

	if cmd.ShowProgress {
		fmt.Print("\r" + utils.Separator() + "\n")
	}

	if len(allResults) == 0 {
		fmt.Println("[*] No open ports found")
	} else {
		fmt.Printf("[+] Found %d open port(s)\n", len(allResults))

		if cmd.ShowProgress {
			fmt.Println()
			for _, res := range allResults {
				utils.PrintResult(res)
			}
		}

		filename := saveResults(allResults, cmd.OutputFormat, cmd.OutputFile)
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

func launchTUI() {
	tui.Run()
}
