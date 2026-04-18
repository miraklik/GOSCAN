package utils

import (
	"context"
	"fmt"
	"goscan/detector"
	"goscan/scanner"
	"sync"
	"sync/atomic"

	"golang.org/x/time/rate"
)

func Worker(id int, host string, sc scanner.Scanner, limiter *rate.Limiter, ports <-chan int, results chan<- scanner.Result, wg *sync.WaitGroup, scanned *int32, verbose bool) {
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
