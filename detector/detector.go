package detector

import "strings"

var serviceSignatures = map[string]string{
	"SSH":        "SSH-",
	"HTTP":       "HTTP/",
	"FTP":        "220",
	"SMTP":       "220",
	"MySQL":      "mysql",
	"PostgreSQL": "PostgreSQL",
	"Redis":      "REDIS",
	"MongoDB":    "MongoDB",
	"Telnet":     "Telnet",
	"POP3":       "+OK",
	"IMAP":       "* OK",
}

var knownPorts = map[int]string{
	20:    "FTP-Data",
	21:    "FTP",
	22:    "SSH",
	23:    "Telnet",
	25:    "SMTP",
	53:    "DNS",
	80:    "HTTP",
	110:   "POP3",
	143:   "IMAP",
	443:   "HTTPS",
	445:   "SMB",
	465:   "SMTPS",
	587:   "SMTP",
	993:   "IMAPS",
	995:   "POP3S",
	1433:  "MSSQL",
	1521:  "Oracle",
	3306:  "MySQL",
	3389:  "RDP",
	5432:  "PostgreSQL",
	5900:  "VNC",
	6379:  "Redis",
	8080:  "HTTP-Proxy",
	8443:  "HTTPS-Alt",
	9200:  "Elasticsearch",
	27017: "MongoDB",
	27018: "MongoDB",
	50000: "DB2",
}

func DetectService(banner string, port int) string {
	for service, signature := range serviceSignatures {
		if strings.Contains(strings.ToUpper(banner), strings.ToUpper(signature)) {
			return service
		}
	}

	if service, ok := knownPorts[port]; ok {
		return service
	}

	return "Unknown"
}

func GetServiceDescription(service string) string {
	descriptions := map[string]string{
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
	}

	if desc, ok := descriptions[service]; ok {
		return desc
	}
	return "Unknown Service"
}
