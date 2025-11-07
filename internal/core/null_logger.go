package core

// NullLogger is a logger that discards all output (used for silent operations)
// Single Responsibility Principle: Only responsible for implementing Logger interface silently
type NullLogger struct{}

// NewNullLogger creates a new NullLogger instance
func NewNullLogger() *NullLogger {
	return &NullLogger{}
}

// Info does nothing (silent)
func (l *NullLogger) Info(message string) {}

// Success does nothing (silent)
func (l *NullLogger) Success(message string) {}

// Error does nothing (silent)
func (l *NullLogger) Error(message string) {}

// Warning does nothing (silent)
func (l *NullLogger) Warning(message string) {}
