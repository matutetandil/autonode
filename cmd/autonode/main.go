package main

import (
	"fmt"
	"os"
	"time"

	// Import commands package to trigger init() functions that register commands
	_ "github.com/matutetandil/autonode/cmd/autonode/commands"
	"github.com/matutetandil/autonode/cmd/autonode/commands"
	"github.com/matutetandil/autonode/internal/core"
	"github.com/spf13/cobra"
)

var version = "0.7.0"

// noUpdateCheck disables the automatic update check (set via --no-update-check flag)
var noUpdateCheck bool

// GetVersion returns the current version of autonode
func GetVersion() string {
	return version
}

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

	// Add global flag to disable update check (useful for CI/CD)
	rootCmd.PersistentFlags().BoolVar(&noUpdateCheck, "no-update-check", false, "Disable automatic update check")

	// Add all other commands as subcommands
	for _, cmd := range allCommands {
		cobraCmd := cmd.GetCobraCommand()
		// Skip the root command itself
		if cobraCmd.Use != "autonode" {
			rootCmd.AddCommand(cobraCmd)
		}
	}

	// Start async update check before executing command
	var updateChecker *core.UpdateChecker
	cache, cacheErr := core.NewCacheManager()
	if !noUpdateCheck && cacheErr == nil {
		// Load global config to check if update check is disabled
		globalConfig, _ := core.LoadGlobalConfig(cache)
		if !globalConfig.DisableUpdateCheck {
			updateChecker = core.NewUpdateChecker(cache, version)
			// Apply custom interval if configured
			if globalConfig.UpdateCheckIntervalDays > 0 {
				updateChecker.SetCheckInterval(
					time.Duration(globalConfig.UpdateCheckIntervalDays) * 24 * time.Hour,
				)
			}
			updateChecker.StartAsyncCheck()
		}
	}

	// Execute
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// Show update notification if available (after command completes)
	if updateChecker != nil {
		result := updateChecker.GetResult()
		if result != nil && result.UpdateAvailable {
			notifier := core.NewUpdateNotifier(core.NewConsoleLogger())
			notifier.ShowUpdateBanner(result)
		}
	}
}
