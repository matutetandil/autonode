package core

// Config holds the configuration for AutoNode execution
// Single Responsibility Principle: Only responsible for holding configuration data
type Config struct {
	ProjectPath string
	CheckOnly   bool
	Force       bool
	ShellMode   bool // When true, outputs shell commands instead of executing them
}
