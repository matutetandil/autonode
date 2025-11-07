package main

import (
	"fmt"
	"os"

	// Import commands package to trigger init() functions that register commands
	_ "github.com/matutetandil/autonode/cmd/autonode/commands"
	"github.com/matutetandil/autonode/cmd/autonode/commands"
	"github.com/spf13/cobra"
)

var version = "0.3.1"

func main() {
	// Create root command
	// This will be replaced by the RunCommand, but we need it for version flag
	var rootCmd *cobra.Command

	// Get all registered commands
	allCommands := commands.GetAll()

	// Find the run command to use as root
	// The run command should not be a subcommand, it should be the root command itself
	for _, cmd := range allCommands {
		cobraCmd := cmd.GetCobraCommand()
		if cobraCmd.Use == "autonode" {
			// This is the main run command - use it as root
			rootCmd = cobraCmd
			rootCmd.Version = version
			break
		}
	}

	if rootCmd == nil {
		fmt.Fprintln(os.Stderr, "Error: root command not found")
		os.Exit(1)
	}

	// Add all other commands as subcommands
	for _, cmd := range allCommands {
		cobraCmd := cmd.GetCobraCommand()
		// Skip the root command itself
		if cobraCmd.Use != "autonode" {
			rootCmd.AddCommand(cobraCmd)
		}
	}

	// Execute
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
