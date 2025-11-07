package commands

import "github.com/spf13/cobra"

// Command interface defines a command that can be registered with the CLI
// Open/Closed Principle: New commands can be added without modifying existing code
// Strategy Pattern: Each command is a different strategy for handling user requests
type Command interface {
	// GetCobraCommand returns the cobra.Command for this command
	GetCobraCommand() *cobra.Command
}

// registry holds all registered commands
var registry []Command

// Register adds a command to the registry
// Called automatically by each command's init() function
func Register(cmd Command) {
	registry = append(registry, cmd)
}

// GetAll returns all registered commands
func GetAll() []Command {
	return registry
}
