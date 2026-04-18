package tui

import (
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	"goscan/scanner"
	"goscan/utils"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	Host     string
	Ports    string
	Workers  string
	Timeout  string
	Protocol string

	ProtocolIndex int

	Scanning bool
	Results  []scanner.Result
	Scanned  int
	Open     int
	Total    int

	Focused int

	Scanner scanner.Scanner
	mu      sync.Mutex
}

type tickMsg time.Time

func New() *Model {
	m := &Model{
		Host:          "127.0.0.1",
		Ports:         "1-1024",
		Workers:       "25",
		Timeout:       "5",
		Protocol:      "TCP",
		ProtocolIndex: 0,
		Focused:       0,
	}
	return m
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKey(msg)
	case tickMsg:
		if m.Scanning {
			return m, tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
				return tickMsg(t)
			})
		}
	}
	return m, nil
}

func (m *Model) handleKey(msg tea.Msg) (tea.Model, tea.Cmd) {
	key := msg.(tea.KeyMsg).String()

	if m.Scanning {
		if key == "ctrl+c" {
			m.Scanning = false
		}
		return m, nil
	}

	switch key {
	case "ctrl+c", "ctrl+q":
		return m, tea.Quit

	case "up", "k":
		if m.Focused > 0 {
			m.Focused--
		}

	case "down", "j":
		if m.Focused < 3 {
			m.Focused++
		} else if m.Focused == 3 {
			m.Focused = 4
		} else if m.Focused == 4 {
			m.Focused = 5
		}

	case "left", "h":
		if m.Focused == 4 {
			m.ProtocolIndex = 0
			m.Protocol = "TCP"
		}

	case "right", "l":
		if m.Focused == 4 {
			m.ProtocolIndex = 1
			m.Protocol = "UDP"
		}

	case "enter":
		if m.Focused == 5 {
			return m, m.startScan
		}

	case "backspace":
		if m.Focused < 4 {
			m.deleteChar()
		}

	default:
		if m.Focused < 4 {
			m.insertChar(key)
		}
	}

	return m, nil
}

func (m *Model) insertChar(ch string) {
	if len(ch) != 1 {
		return
	}

	switch m.Focused {
	case 0:
		m.Host += ch
	case 1:
		m.Ports += ch
	case 2:
		m.Workers += ch
	case 3:
		m.Timeout += ch
	}
}

func (m *Model) deleteChar() {
	switch m.Focused {
	case 0:
		if len(m.Host) > 0 {
			m.Host = m.Host[:len(m.Host)-1]
		}
	case 1:
		if len(m.Ports) > 0 {
			m.Ports = m.Ports[:len(m.Ports)-1]
		}
	case 2:
		if len(m.Workers) > 0 {
			m.Workers = m.Workers[:len(m.Workers)-1]
		}
	case 3:
		if len(m.Timeout) > 0 {
			m.Timeout = m.Timeout[:len(m.Timeout)-1]
		}
	}
}

func (m *Model) startScan() tea.Msg {
	workers, _ := strconv.Atoi(m.Workers)
	if workers < 1 {
		workers = 25
	}

	timeout, _ := strconv.Atoi(m.Timeout)
	if timeout < 1 {
		timeout = 5
	}

	ports := utils.ParsePorts(m.Ports)
	if len(ports) == 0 {
		return nil
	}

	m.Scanning = true
	m.Results = nil
	m.Scanned = 0
	m.Open = 0
	m.Total = len(ports)

	if m.ProtocolIndex == 0 {
		m.Scanner = scanner.NewTCPScanner(time.Duration(timeout) * time.Second)
	} else {
		m.Scanner = scanner.NewUDPScanner(time.Duration(timeout) * time.Second)
	}

	go m.runScan(ports, workers)

	return tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m *Model) runScan(ports []int, workers int) {
	portChan := make(chan int, 100)
	resultChan := make(chan scanner.Result, 100)
	var wg sync.WaitGroup

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for port := range portChan {
				addr := fmt.Sprintf("%s:%d", m.Host, port)
				_, err := m.Scanner.ScanPort(addr)

				m.mu.Lock()
				m.Scanned++
				m.mu.Unlock()

				if err == nil {
					result := scanner.Result{
						Port:    port,
						Service: "Unknown",
						Status:  "open",
					}
					resultChan <- result
				}
			}
		}()
	}

	go func() {
		for _, p := range ports {
			portChan <- p
		}
		close(portChan)
	}()

	go func() {
		for res := range resultChan {
			m.mu.Lock()
			m.Results = append(m.Results, res)
			m.Open++
			m.mu.Unlock()
		}
		wg.Wait()

		m.mu.Lock()
		m.Scanning = false
		m.mu.Unlock()
	}()
}

func (m *Model) View() string {
	green := lipgloss.NewStyle().Foreground(lipgloss.Color("42"))
	yellow := lipgloss.NewStyle().Foreground(lipgloss.Color("226"))
	cyan := lipgloss.NewStyle().Foreground(lipgloss.Color("36"))
	gray := lipgloss.NewStyle().Foreground(lipgloss.Color("8"))

	s := ""

	s += cyan.Bold(true).Render("  GOSCAN - Network Scanner\n")
	s += gray.Render("  ════════════════════════════════════════\n\n")

	labels := []string{"Host", "Ports", "Workers", "Timeout"}
	values := []string{m.Host, m.Ports, m.Workers, m.Timeout}

	for i := range labels {
		prefix := "  "
		if m.Focused == i {
			prefix = cyan.Render("► ")
		}
		s += fmt.Sprintf("%s%-10s %s\n", prefix, labels[i]+":", values[i])
	}

	protoPrefix := "  "
	if m.Focused == 4 {
		protoPrefix = cyan.Render("► ")
	}
	s += fmt.Sprintf("\n%sProtocol:   %s\n", protoPrefix, green.Render(m.Protocol))

	s += "\n" + gray.Render("  ───────────────────────────────────────────") + "\n"

	scanBtn := "[ SCAN ]"
	if m.Focused == 5 {
		scanBtn = yellow.Bold(true).Render("[ SCAN ]")
	}

	s += fmt.Sprintf("  %s\n\n", scanBtn)

	if m.Scanning && m.Total > 0 {
		percent := float64(m.Scanned) / float64(m.Total) * 100
		s += yellow.Render("  Scanning: ") + fmt.Sprintf("%d/%d (%.1f%%) | %d open\n", m.Scanned, m.Total, percent, m.Open)

		bar := ""
		width := 30
		filled := int(float64(width) * float64(m.Scanned) / float64(m.Total))
		for i := 0; i < width; i++ {
			if i < filled {
				bar += green.Render("█")
			} else {
				bar += gray.Render("░")
			}
		}
		s += "  " + bar + "\n"
	} else if len(m.Results) > 0 {
		s += green.Bold(true).Render(fmt.Sprintf("\n  ✓ Found %d open port(s):\n\n", len(m.Results)))

		for _, r := range m.Results {
			s += fmt.Sprintf("    %s %d  (%s)\n", green.Render("●"), r.Port, r.Service)
		}
		s += "\n"
	} else if !m.Scanning && m.Total > 0 {
		s += gray.Render("  ✓ Scan complete. No open ports found.\n")
	} else {
		s += gray.Render("  Ready to scan\n")
	}

	s += "\n" + gray.Render("  [↑/↓] Move  [←/→] Protocol  [Enter] Scan  [Ctrl+C] Quit\n")

	return s
}

func Run() {
	p := tea.NewProgram(New(), tea.WithAltScreen())
	if err := p.Start(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}
