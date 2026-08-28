package utils

import (
	"fmt"
	"os"
	"time"
)

func LogToFile(filename string, message any) error {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file %s: %w", filename, err)
	}
	defer f.Close()

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	if _, err := fmt.Fprintf(f, "[%s] %s\n", timestamp, message); err != nil {
		return fmt.Errorf("failed to write to log file %s: %w", filename, err)
	}

	return nil
}
