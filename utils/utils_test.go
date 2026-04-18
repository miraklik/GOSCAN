package utils

import (
	"testing"
	"time"
)

func TestParsePorts_SinglePort(t *testing.T) {
	tests := []struct {
		name      string
		portRange string
		wantPorts []int
	}{
		{"single port 80", "80", []int{80}},
		{"single port 443", "443", []int{443}},
		{"single port 8080", "8080", []int{8080}},
		{"single port 1", "1", []int{1}},
		{"single port 65535", "65535", []int{65535}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParsePorts(tt.portRange)
			if len(got) != len(tt.wantPorts) || (len(got) > 0 && got[0] != tt.wantPorts[0]) {
				t.Errorf("ParsePorts(%q) = %v, want %v", tt.portRange, got, tt.wantPorts)
			}
		})
	}
}

func TestParsePorts_PortRange(t *testing.T) {
	tests := []struct {
		name      string
		portRange string
		wantLen   int
		wantFirst int
		wantLast  int
	}{
		{"range 1-10", "1-10", 10, 1, 10},
		{"range 80-80", "80-80", 1, 80, 80},
		{"range 1-100", "1-100", 100, 1, 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParsePorts(tt.portRange)
			if len(got) != tt.wantLen {
				t.Errorf("ParsePorts(%q) length = %d, want %d", tt.portRange, len(got), tt.wantLen)
			}
			if len(got) > 0 && got[0] != tt.wantFirst {
				t.Errorf("ParsePorts(%q) first = %d, want %d", tt.portRange, got[0], tt.wantFirst)
			}
			if len(got) > 0 && got[len(got)-1] != tt.wantLast {
				t.Errorf("ParsePorts(%q) last = %d, want %d", tt.portRange, got[len(got)-1], tt.wantLast)
			}
		})
	}
}

func TestParsePorts_PortList(t *testing.T) {
	tests := []struct {
		name      string
		portRange string
		wantPorts []int
	}{
		{"list 80,443", "80,443", []int{80, 443}},
		{"list 22,80,443", "22,80,443", []int{22, 80, 443}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParsePorts(tt.portRange)
			if len(got) != len(tt.wantPorts) {
				t.Errorf("ParsePorts(%q) = %v, want %v", tt.portRange, got, tt.wantPorts)
			}
			for i, want := range tt.wantPorts {
				if i >= len(got) || got[i] != want {
					t.Errorf("ParsePorts(%q)[%d] = %d, want %d", tt.portRange, i, got[i], want)
				}
			}
		})
	}
}

func TestParsePorts_InvalidInput(t *testing.T) {
	tests := []struct {
		name      string
		portRange string
	}{
		{"port too high", "70000"},
		{"port 0", "0"},
		{"invalid string", "abc"},
		{"empty range", "-"},
		{"start > end", "100-50"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParsePorts(tt.portRange)
			if got != nil {
				t.Errorf("ParsePorts(%q) = %v, want nil", tt.portRange, got)
			}
		})
	}
}

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		name     string
		d        time.Duration
		want     string
	}{
		{"milliseconds", 500 * time.Millisecond, "500ms"},
		{"seconds", 5 * time.Second, "5.00s"},
		{"minutes", 2 * time.Minute, "2m 0s"},
		{"minutes seconds", 2*time.Minute + 30*time.Second, "2m 30s"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FormatDuration(tt.d)
			if got != tt.want {
				t.Errorf("FormatDuration(%v) = %q, want %q", tt.d, got, tt.want)
			}
		})
	}
}

func TestSeparator(t *testing.T) {
	got := Separator()
	if len(got) != 60 {
		t.Errorf("Separator() len = %d, want 60", len(got))
	}
}

func TestSeparatorLength(t *testing.T) {
	got := Separator()
	if len(got) != 60 {
		t.Errorf("Separator() len = %d, want 60", len(got))
	}
}