package core

import (
	"fmt"

	"github.com/fatih/color"
)

// ConsoleLogger implements the Logger interface using colored console output
// Single Responsibility Principle: Only responsible for logging to console
// Dependency Inversion Principle: Depends on Logger interface, can be swapped with other implementations
type ConsoleLogger struct{}

// NewConsoleLogger creates a new ConsoleLogger instance
func NewConsoleLogger() *ConsoleLogger {
	return &ConsoleLogger{}
}

// Info logs an informational message in cyan
func (l *ConsoleLogger) Info(message string) {
	cyan := color.New(color.FgCyan)
	cyan.Println(message)
}

// Success logs a success message in green
func (l *ConsoleLogger) Success(message string) {
	green := color.New(color.FgGreen)
	green.Println("✓", message)
}

// Error logs an error message in red
func (l *ConsoleLogger) Error(message string) {
	red := color.New(color.FgRed)
	red.Println("✗", message)
}

// Warning logs a warning message in yellow
func (l *ConsoleLogger) Warning(message string) {
	yellow := color.New(color.FgYellow)
	fmt.Print("⚠ ")
	yellow.Println(message)
}
