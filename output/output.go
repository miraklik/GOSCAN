package output

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"goscan/pkg"
	"goscan/scanner"
	"goscan/utils"
	"os"
	"strconv"
	"time"
)

type ScanReport struct {
	Metadata struct {
		ScanDate   time.Time   `json:"scan_date"`
		Target     string      `json:"target"`
		Protocol   string      `json:"protocol"`
		SystemInfo pkg.SysInfo `json:"system_info"`
	} `json:"metadata"`
	Results []scanner.Result `json:"results"`
	Stats   struct {
		TotalPorts  int           `json:"total_ports"`
		OpenPorts   int           `json:"open_ports"`
		ClosedPorts int           `json:"closed_ports"`
		Duration    time.Duration `json:"duration"`
	} `json:"statistics"`
}

func SaveAsJSON(results []scanner.Result, filename string) error {
	jsonData, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %v", err)
	}

	if err := os.WriteFile(filename, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write file: %v", err)
	}

	return nil
}

func SaveAsCSV(results []scanner.Result, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	if err := writer.Write([]string{"Port", "Service", "Banner", "Status"}); err != nil {
		return fmt.Errorf("failed to write headers: %v", err)
	}

	for _, r := range results {
		record := []string{
			strconv.Itoa(r.Port),
			r.Service,
			r.Banner,
			r.Status,
		}
		if err := writer.Write(record); err != nil {
			return fmt.Errorf("failed to write record: %v", err)
		}
	}

	return nil
}

func SaveAsTXT(results []scanner.Result, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()

	fmt.Fprintln(file, "Port Scan Results")
	fmt.Fprintln(file, "=================")
	fmt.Fprintln(file)

	for _, r := range results {
		fmt.Fprintf(file, "[%s] Port %d (%s)\n", r.Status, r.Port, r.Service)
		if r.Banner != "" && r.Banner != "No banner" {
			fmt.Fprintf(file, "    Banner: %s\n", r.Banner)
		}
		fmt.Fprintln(file)
	}

	return nil
}

func SaveDetailedJSON(results []scanner.Result, filename string, target string, protocol string, stats utils.Stats) error {
	report := ScanReport{}

	report.Metadata.ScanDate = time.Now()
	report.Metadata.Target = target
	report.Metadata.Protocol = protocol
	report.Metadata.SystemInfo = pkg.GetSysInfo()

	report.Results = results

	report.Stats.TotalPorts = stats.TotalPorts
	report.Stats.OpenPorts = stats.OpenPorts
	report.Stats.ClosedPorts = stats.ClosedPorts
	report.Stats.Duration = stats.Duration

	jsonData, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %v", err)
	}

	if err := os.WriteFile(filename, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write file: %v", err)
	}

	return nil
}
