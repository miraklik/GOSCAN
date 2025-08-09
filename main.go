package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

func scanPort(host string, port int) (string, error) {
	addr := fmt.Sprintf("%s:%v", host, port)
	conn, err := net.DialTimeout("tcp", addr, time.Second*5)
	if err != nil {
		return "", fmt.Errorf("failed to connect to host port %d: %v", port, err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			log.Printf("Failed to close connection to host port %d: %s", port, err)
		}
	}()

	return conn.RemoteAddr().String(), nil
}

func worker(host string, ports <-chan int, result chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()
	for p := range ports {
		addr, err := scanPort(host, p)
		if err == nil {
			result <- fmt.Sprintf("[OPEN]%s", addr)
		}
	}
}

func banner() string {
	return `
 ██████╗  ██████╗ ███████╗ ██████╗ █████╗ ███╗   ██╗
██╔════╝ ██╔═══██╗██╔════╝██╔════╝██╔══██╗████╗  ██║
██║  ███╗██║   ██║███████╗██║     ███████║██╔██╗ ██║
██║   ██║██║   ██║╚════██║██║     ██╔══██║██║╚██╗██║
╚██████╔╝╚██████╔╝███████║╚██████╗██║  ██║██║ ╚████║
 ╚═════╝  ╚═════╝ ╚══════╝ ╚═════╝╚═╝  ╚═╝╚═╝  ╚═══╝
`
}

func main() {
	fmt.Println(banner())

	host := flag.String("host", "127.0.0.1", "Define a host for the service")
	NumWorkers := flag.Int("w", 25, "Number of workers")
	flag.Parse()

	start := time.Now()
	var wg sync.WaitGroup
	ports := make(chan int, 1000)
	results := make(chan string)

	for i := 0; i < *NumWorkers; i++ {
		wg.Add(1)
		go worker(*host, ports, results, &wg)
	}

	go func() {
		for result := range results {
			fmt.Println(result)
		}
	}()

	go func() {
		for p := 1; p <= 65535; p++ {
			ports <- p
		}

		close(ports)
		wg.Wait()
		close(results)
	}()

	fmt.Printf("[*]Scanning %s...\n", *host)
	<-time.After(1 * time.Second)
	fmt.Println("[*]Scan Completed in ", time.Since(start))
}
