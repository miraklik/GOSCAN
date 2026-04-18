package utils

import (
	"fmt"
	"goscan/scanner"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

type Stats struct {
	TotalPorts    int
	OpenPorts     int
	ClosedPorts   int
	FilteredPorts int
	Duration      time.Duration
}

const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorPurple = "\033[35m"
	ColorCyan   = "\033[36m"
	ColorWhite  = "\033[37m"
)

func ParsePorts(portRange string) []int {
	var ports []int

	if !strings.Contains(portRange, "-") && !strings.Contains(portRange, ",") {
		p, err := strconv.Atoi(strings.TrimSpace(portRange))
		if err == nil && p > 0 && p <= 65535 {
			return []int{p}
		}
		return nil
	}

	if strings.Contains(portRange, ",") {
		parts := strings.Split(portRange, ",")
		for _, part := range parts {
			port, err := strconv.Atoi(strings.TrimSpace(part))
			if err == nil && port > 0 && port <= 65535 {
				ports = append(ports, port)
			}
		}
		return ports
	}

	if strings.Contains(portRange, "-") {
		parts := strings.Split(portRange, "-")
		if len(parts) != 2 {
			return nil
		}

		start, err1 := strconv.Atoi(strings.TrimSpace(parts[0]))
		end, err2 := strconv.Atoi(strings.TrimSpace(parts[1]))

		if err1 != nil || err2 != nil || start < 1 || end > 65535 || start > end {
			return nil
		}

		for p := start; p <= end; p++ {
			ports = append(ports, p)
		}
		return ports
	}

	return nil
}

func Separator() string {
	return strings.Repeat("=", 60)
}

func FormatDuration(d time.Duration) string {
	if d < time.Second {
		return d.String()
	}
	if d < time.Minute {
		return strconv.FormatFloat(d.Seconds(), 'f', 2, 64) + "s"
	}
	minutes := int(d.Minutes())
	seconds := int(d.Seconds()) % 60
	return strconv.Itoa(minutes) + "m " + strconv.Itoa(seconds) + "s"
}

func ProgressMonitor(scanned, open *int32, total int) {
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

func PrintResult(r scanner.Result) {
	fmt.Printf("%s[+]%s Port %s%5d%s | %-12s | %s\n",
		ColorGreen, ColorReset,
		ColorYellow, r.Port, ColorReset,
		r.Service, r.Banner)
}

func PrintStats(stats Stats) {
	fmt.Println(Separator())
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
	fmt.Println(Separator())
}
