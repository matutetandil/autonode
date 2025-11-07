package core

// Logger interface defines logging operations
// Interface Segregation Principle: Small, focused interface with only logging methods
// Dependency Inversion Principle: High-level modules depend on this abstraction
type Logger interface {
	Info(message string)
	Success(message string)
	Error(message string)
	Warning(message string)
}
