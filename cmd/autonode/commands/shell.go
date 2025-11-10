package commands

import (
	"os"

	"github.com/matutetandil/autonode/internal/core"
	"github.com/matutetandil/autonode/internal/detectors"
	"github.com/matutetandil/autonode/internal/managers"
	"github.com/matutetandil/autonode/internal/switchers"
	"github.com/spf13/cobra"
)

// ShellCommand implements the shell integration command
// Single Responsibility Principle: Only responsible for outputting shell commands for eval
type ShellCommand struct{}

// init registers this command automatically when the package is imported
func init() {
	Register(&ShellCommand{})
}

// GetCobraCommand returns the cobra command for this command
func (c *ShellCommand) GetCobraCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "shell",
		Short: "Output shell commands for eval (used by shell integration)",
		Long: `Outputs shell commands to switch Node.js version.
Used by the shell integration hook. Usage: eval "$(autonode shell)"`,
		RunE: c.run,
	}
}

// run outputs shell commands for eval integration using AutoNodeService
// This is used by the shell hook for automatic version switching
func (c *ShellCommand) run(cmd *cobra.Command, args []string) error {
	// Get current working directory
	projectPath, err := os.Getwd()
	if err != nil {
		// Silent failure - just exit without output
		return nil
	}

	// Create configuration with ShellMode enabled
	config := core.Config{
		ProjectPath: projectPath,
		ShellMode:   true, // This tells the service to output commands instead of executing them
	}

	// Dependency Injection: Create all concrete implementations
	// Use NullLogger for silent operation (no colorful output)
	logger := core.NewNullLogger()
	shell := core.NewExecShell()

	// Create cache manager for Node.js releases
	cache, err := core.NewCacheManager()
	if err != nil {
		// Silent failure - just exit without output
		return nil
	}

	// Create Node.js releases client (for Dockerfile codename resolution)
	releasesClient := core.NewNodeReleasesClient(cache, logger)

	// Create all version detectors
	detectorsList := []core.VersionDetector{
		detectors.NewNvmrcDetector(),
		detectors.NewNodeVersionDetector(),
		detectors.NewPackageJsonDetector(),
		detectors.NewDockerfileDetector(releasesClient),
	}

	// Create all version managers
	managersList := []core.VersionManager{
		managers.NewNvmManager(shell),
		managers.NewNvsManager(shell),
		managers.NewVoltaManager(shell),
	}

	// Create all profile detectors
	profileDetectorsList := []core.ProfileDetector{
		detectors.NewAutonodeYmlProfileDetector(),
		detectors.NewPackageJsonProfileDetector(),
	}

	// Create all profile switchers
	profileSwitchersList := []core.ProfileSwitcher{
		switchers.NewNpmrcSwitcher(shell),
		switchers.NewTsNpmrcSwitcher(shell),
		switchers.NewRcManagerSwitcher(shell),
	}

	// Create the service with all dependencies
	service := core.NewAutoNodeService(logger, detectorsList, managersList, profileDetectorsList, profileSwitchersList)

	// Run the service in shell mode (outputs commands, doesn't execute them)
	return service.Run(config)
}
