package detector

import "testing"

func TestDetectService_SSHBanner(t *testing.T) {
	tests := []struct {
		name   string
		banner string
		port  int
		want   string
	}{
		{"SSH banner SSH-2.0", "SSH-2.0-OpenSSH_8.2", 22, "SSH"},
		{"SSH banner lowercase", "ssh-2.0-drobak", 22, "SSH"},
		{"HTTP banner", "HTTP/1.1 200 OK", 80, "HTTP"},
		{"HTTP banner lowercase", "http/1.1", 8080, "HTTP"},
		{"FTP banner", "220 FTP server ready", 21, "FTP"},
		{"SMTP banner", "220 mail.example.com ESMTP", 25, "SMTP"},
		{"MySQL banner", "5.7.30-log MySQL", 3306, "MySQL"},
		{"PostgreSQL banner", "PostgreSQL 13.0", 5432, "PostgreSQL"},
		{"Redis banner", "REDIS 6.0.5", 6379, "Redis"},
		{"MongoDB banner", "MongoDB 4.4", 27017, "MongoDB"},
		{"Telnet banner", "Telnet server", 23, "Telnet"},
		{"POP3 banner", "+OK POP3 server", 110, "POP3"},
		{"IMAP banner", "* OK IMAP4rev1", 143, "IMAP"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DetectService(tt.banner, tt.port)
			if got != tt.want {
				t.Errorf("DetectService(%q, %d) = %q, want %q", tt.banner, tt.port, got, tt.want)
			}
		})
	}
}

func TestDetectService_ByKnownPort(t *testing.T) {
	tests := []struct {
		name  string
		port  int
		want  string
	}{
		{"FTP port 21", 21, "FTP"},
		{"SSH port 22", 22, "SSH"},
		{"Telnet port 23", 23, "Telnet"},
		{"SMTP port 25", 25, "SMTP"},
		{"DNS port 53", 53, "DNS"},
		{"HTTP port 80", 80, "HTTP"},
		{"POP3 port 110", 110, "POP3"},
		{"IMAP port 143", 143, "IMAP"},
		{"HTTPS port 443", 443, "HTTPS"},
		{"SMB port 445", 445, "SMB"},
		{"MySQL port 3306", 3306, "MySQL"},
		{"RDP port 3389", 3389, "RDP"},
		{"PostgreSQL port 5432", 5432, "PostgreSQL"},
		{"VNC port 5900", 5900, "VNC"},
		{"Redis port 6379", 6379, "Redis"},
		{"HTTP-Proxy port 8080", 8080, "HTTP-Proxy"},
		{"MongoDB port 27017", 27017, "MongoDB"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DetectService("", tt.port)
			if got != tt.want {
				t.Errorf("DetectService(%q, %d) = %q, want %q", "", tt.port, got, tt.want)
			}
		})
	}
}

func TestDetectService_NoSignature(t *testing.T) {
	got := DetectService("some random string that doesn't match", 12345)
	if got != "Unknown" {
		t.Errorf("DetectService() = %q, want Unknown", got)
	}
}

func TestGetServiceDescription(t *testing.T) {
	tests := map[string]string{
		"SSH":           "Secure Shell",
		"HTTP":          "Hypertext Transfer Protocol",
		"HTTPS":         "HTTP Secure",
		"FTP":           "File Transfer Protocol",
		"SMTP":          "Simple Mail Transfer Protocol",
		"MySQL":         "MySQL Database",
		"PostgreSQL":    "PostgreSQL Database",
		"Redis":         "Redis Key-Value Store",
		"MongoDB":       "MongoDB NoSQL Database",
		"Elasticsearch": "Elasticsearch Search Engine",
		"DNS":           "Domain Name System",
		"RDP":           "Remote Desktop Protocol",
		"VNC":           "Virtual Network Computing",
		"Unknown":       "Unknown Service",
	}

	for service, want := range tests {
		t.Run(service, func(t *testing.T) {
			got := GetServiceDescription(service)
			if got != want {
				t.Errorf("GetServiceDescription(%q) = %q, want %q", service, got, want)
			}
		})
	}
}