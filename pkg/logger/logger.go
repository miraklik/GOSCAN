package logger

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

type Level int

const (
	DEBUG Level = iota
	INFO
	WARN
	ERROR
)

var (
	defaultLogger Logger
	once         sync.Once
)

type Logger struct {
	mu      sync.Mutex
	level   Level
	output  *log.Logger
	verbose bool
}

func Init(verbose bool) {
	once.Do(func() {
		defaultLogger = Logger{
			level:  INFO,
			output: log.New(os.Stdout, "", 0),
		}
	})
	defaultLogger.mu.Lock()
	defaultLogger.verbose = verbose
	if verbose {
		defaultLogger.level = DEBUG
	}
	defaultLogger.mu.Unlock()
}

func SetLevel(level Level) {
	defaultLogger.mu.Lock()
	defer defaultLogger.mu.Unlock()
	defaultLogger.level = level
}

func Debug(format string, v ...interface{}) {
	defaultLogger.mu.Lock()
	defer defaultLogger.mu.Unlock()
	if defaultLogger.level <= DEBUG {
		defaultLogger.output.Printf("[DEBUG] "+format, v...)
	}
}

func Info(format string, v ...interface{}) {
	defaultLogger.mu.Lock()
	defer defaultLogger.mu.Unlock()
	if defaultLogger.level <= INFO {
		timestamp := time.Now().Format("2006/01/02 15:04:05")
		defaultLogger.output.Printf(fmt.Sprintf("[%s] [INFO] ", timestamp)+format, v...)
	}
}

func Warn(format string, v ...interface{}) {
	defaultLogger.mu.Lock()
	defer defaultLogger.mu.Unlock()
	if defaultLogger.level <= WARN {
		defaultLogger.output.Printf("[WARN] "+format, v...)
	}
}

func Error(format string, v ...interface{}) {
	defaultLogger.mu.Lock()
	defer defaultLogger.mu.Unlock()
	if defaultLogger.level <= ERROR {
		defaultLogger.output.Printf("[ERROR] "+format, v...)
	}
}