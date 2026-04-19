package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
)

var (
	Host         string
	NumWorkers   int
	PortRange    string
	Protocol     string
	Timeout      time.Duration
	OutputFormat string
	OutputFile   string
	ShowProgress bool
	Verbose      bool
	TUI          bool
	ConfigFile   string
)

var commandCMD = &cobra.Command{
	Use:   "goscan",
	Short: "goscan - fast network scanner",
	Long: `goscan - This is a powerful port scanning tool.
    
TARGET SPECIFICATION:
Valid targets are hostnames, IP addresses, networks, etc.
Ex: scanme.nmap.org, 192.168.0.1, 192.168.0.0/24

OUTPUT:
Specify a file to write the output.
Ex: --output-format json --output-file results`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Scanning host: %s | Workers: %d\n", Host, NumWorkers)
	},
}

func Execute() {
	if err := commandCMD.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	commandCMD.Flags().StringVarP(&Host, "host", "H", "", "Target host or IP address to scan (required)")
	commandCMD.Flags().IntVarP(&NumWorkers, "workers", "w", 16, "Number of concurrent workers for scanning")
	commandCMD.Flags().StringVarP(&PortRange, "ports", "p", "1-65535", "Port range to scan (e.g., 1-1000)")
	commandCMD.Flags().StringVarP(&Protocol, "protocol", "", "tcp", "Protocol to scan (tcp or udp)")
	commandCMD.Flags().DurationVarP(&Timeout, "timeout", "", 5*time.Second, "Timeout for each scan attempt")
	commandCMD.Flags().StringVarP(&OutputFormat, "output-format", "", "text", "Output format (text, json, xml)")
	commandCMD.Flags().StringVarP(&OutputFile, "output-file", "", "", "File to save the scan results")
	commandCMD.Flags().BoolVarP(&ShowProgress, "progress", "", false, "Show progress during scanning")
	commandCMD.Flags().BoolVarP(&Verbose, "verbose", "v", false, "Enable verbose output")
	commandCMD.Flags().BoolVarP(&TUI, "tui", "", false, "Enable Text User Interface (TUI) mode")
	commandCMD.Flags().StringVarP(&ConfigFile, "config", "c", "", "Path to configuration file")
	commandCMD.MarkFlagRequired("host")
}
