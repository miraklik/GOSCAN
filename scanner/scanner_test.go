package scanner

import (
	"testing"
	"time"
)

func TestNewTCPScanner(t *testing.T) {
	s := NewTCPScanner(5 * time.Second)
	if s == nil {
		t.Fatal("NewTCPScanner returned nil")
	}
	if s.Timeout != 5*time.Second {
		t.Errorf("expected timeout 5s, got %v", s.Timeout)
	}
}

func TestNewUDPScanner(t *testing.T) {
	s := NewUDPScanner(3 * time.Second)
	if s == nil {
		t.Fatal("NewUDPScanner returned nil")
	}
	if s.Timeout != 3*time.Second {
		t.Errorf("expected timeout 3s, got %v", s.Timeout)
	}
}

func TestTCPScanner_InvalidAddr(t *testing.T) {
	s := NewTCPScanner(100 * time.Millisecond)
	_, err := s.ScanPort("invalid:address")
	if err == nil {
		t.Error("expected error for invalid address, got nil")
	}
}

func TestUDPScanner_InvalidAddr(t *testing.T) {
	s := NewUDPScanner(100 * time.Millisecond)
	_, err := s.ScanPort("invalid:address")
	if err == nil {
		t.Error("expected error for invalid address, got nil")
	}
}

func TestResult_Fields(t *testing.T) {
	r := Result{
		Port:    22,
		Service: "SSH",
		Banner:  "SSH-2.0-OpenSSH",
		Status:  "open",
	}

	if r.Port != 22 {
		t.Errorf("expected port 22, got %d", r.Port)
	}
	if r.Service != "SSH" {
		t.Errorf("expected service SSH, got %s", r.Service)
	}
	if r.Status != "open" {
		t.Errorf("expected status open, got %s", r.Status)
	}
}