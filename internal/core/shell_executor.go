package core

// ShellExecutor interface defines shell command execution operations
// Interface Segregation Principle: Focused interface for command execution only
// Dependency Inversion Principle: Managers depend on this abstraction, not concrete implementations
type ShellExecutor interface {
	Execute(command string, args ...string) (string, error)
	ExecuteInShell(command string) (string, error)
	CommandExists(command string) bool
}
