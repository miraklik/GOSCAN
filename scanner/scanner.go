package scanner

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

type Scanner interface {
	ScanPort(addr string) (string, error)
}

type Result struct {
	Port    int    `json:"port"`
	Service string `json:"service"`
	Banner  string `json:"banner"`
	Status  string `json:"status"`
}

type TCPScanner struct {
	Timeout time.Duration
}

type UDPScanner struct {
	Timeout time.Duration
}

func NewTCPScanner(timeout time.Duration) *TCPScanner {
	return &TCPScanner{Timeout: timeout}
}

func NewUDPScanner(timeout time.Duration) *UDPScanner {
	return &UDPScanner{Timeout: timeout}
}

func (t *TCPScanner) ScanPort(addr string) (string, error) {
	conn, err := net.DialTimeout("tcp", addr, t.Timeout)
	if err != nil {
		return "", fmt.Errorf("failed to connect: %v", err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			log.Printf("Failed to close connection to %s: %v", addr, err)
		}
	}()

	conn.SetReadDeadline(time.Now().Add(t.Timeout))

	fmt.Fprintf(conn, "Hello\n\n")

	reader := bufio.NewReader(conn)
	banner, err := reader.ReadString('\n')
	if err != nil {
		return "No banner", nil
	}

	return strings.TrimSpace(banner), nil
}

func (u *UDPScanner) ScanPort(addr string) (string, error) {
	conn, err := net.DialTimeout("udp", addr, u.Timeout)
	if err != nil {
		return "", fmt.Errorf("failed to connect: %v", err)
	}
	defer conn.Close()

	conn.SetDeadline(time.Now().Add(u.Timeout))

	_, err = conn.Write([]byte("probe"))
	if err != nil {
		return "", fmt.Errorf("failed to send probe: %v", err)
	}

	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		return "open|filtered", nil
	}

	return fmt.Sprintf("open - %s", string(buf[:n])), nil
}
