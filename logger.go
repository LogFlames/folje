package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var (
	logFile   *os.File
	logMutex  sync.Mutex
	loggerInit bool
)

// InitLogger initializes file logging to folje_log.txt next to the executable
func InitLogger() error {
	logMutex.Lock()
	defer logMutex.Unlock()

	if loggerInit {
		return nil
	}

	// Get executable path
	exePath, err := os.Executable()
	if err != nil {
		// Fall back to current directory
		exePath = "."
	} else {
		exePath = filepath.Dir(exePath)
	}

	logPath := filepath.Join(exePath, "folje_log.txt")

	// Open log file for append, create if doesn't exist
	logFile, err = os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		// Try current directory as fallback
		logPath = "folje_log.txt"
		logFile, err = os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return fmt.Errorf("failed to open log file: %w", err)
		}
	}

	loggerInit = true
	return nil
}

// CloseLogger closes the log file
func CloseLogger() {
	logMutex.Lock()
	defer logMutex.Unlock()

	if logFile != nil {
		logFile.Close()
		logFile = nil
	}
	loggerInit = false
}

func logMessage(level string, format string, args ...interface{}) {
	logMutex.Lock()
	defer logMutex.Unlock()

	timestamp := time.Now().Format("2006-01-02 15:04:05.000")
	message := fmt.Sprintf(format, args...)
	logLine := fmt.Sprintf("[%s] [%s] %s\n", timestamp, level, message)

	// Write to stderr
	fmt.Fprint(os.Stderr, logLine)

	// Write to file if available
	if logFile != nil {
		logFile.WriteString(logLine)
		logFile.Sync()
	}
}

// LogInfo logs an info message
func LogInfo(format string, args ...interface{}) {
	logMessage("INFO", format, args...)
}

// LogError logs an error message
func LogError(format string, args ...interface{}) {
	logMessage("ERROR", format, args...)
}

// LogDebug logs a debug message
func LogDebug(format string, args ...interface{}) {
	logMessage("DEBUG", format, args...)
}

// LogFatal logs a fatal message and exits
func LogFatal(format string, args ...interface{}) {
	logMessage("FATAL", format, args...)
	os.Exit(1)
}
