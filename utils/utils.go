package utils

import (
	"strconv"
	"strings"
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
