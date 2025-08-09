package pkg

import (
	"fmt"
)

func Usage() {
	fmt.Println(`
  ╔══════════════════════════════════════════════╗
  ║                  GoScan                      ║
  ╠══════════════════════════════════════════════╣
  ║ Usage:                                       ║
  ║   goscan -host 192.168.1.1 -w 50             ║
  ║                                              ║
  ║ Options:                                     ║
  ║   -host string  Target host to scan          ║
  ║   -w int        Worker threads (default 25)  ║
  ║   -timeout int  Timeout in sec (default 5)   ║
  ╚══════════════════════════════════════════════╝
`)
}
