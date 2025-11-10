package commands

import (
	"fmt"
	"os"

	"github.com/matutetandil/autonode/internal/core"
	"github.com/matutetandil/autonode/internal/detectors"
	"github.com/matutetandil/autonode/internal/managers"
	"github.com/matutetandil/autonode/internal/switchers"
	"github.com/spf13/cobra"
)

// RunCommand implements the main autonode command (detect and switch versions)
// Single Responsibility Principle: Only responsible for version detection and switching
type RunCommand struct {
	checkOnly bool
	force     bool
}

// init registers this command automatically when the package is imported
func init() {
	Register(&RunCommand{})
}

// GetCobraCommand returns the cobra command for this command
func (c *RunCommand) GetCobraCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "autonode",
		Short: "Automatically detect and switch Node.js versions",
		Long: `AutoNode detects the required Node.js version from your project
and automatically switches to it using your installed version manager (nvm, nvs, or volta).`,
		RunE: c.run,
	}

	cmd.Flags().BoolVarP(&c.checkOnly, "check", "c", false, "Only check and display the detected version without switching")
	cmd.Flags().BoolVarP(&c.force, "force", "f", false, "Force reinstall the version even if already installed")

	return cmd
}

// run is the main entry point - this is where dependency injection happens (composition root)
// Dependency Inversion Principle: We create all dependencies here and inject them
func (c *RunCommand) run(cmd *cobra.Command, args []string) error {
	// Get current working directory
	projectPath, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	// Create configuration
	config := core.Config{
		ProjectPath: projectPath,
		CheckOnly:   c.checkOnly,
		Force:       c.force,
	}

	// Dependency Injection: Create all concrete implementations
	logger := core.NewConsoleLogger()
	shell := core.NewExecShell()

	// Create cache manager for Node.js releases
	cache, err := core.NewCacheManager()
	if err != nil {
		return fmt.Errorf("failed to create cache manager: %w", err)
	}

	// Create Node.js releases client (for Dockerfile codename resolution)
	releasesClient := core.NewNodeReleasesClient(cache, logger)

	// Create all version detectors
	// Open/Closed Principle: Adding new detectors doesn't require modifying existing code
	detectorsList := []core.VersionDetector{
		detectors.NewNvmrcDetector(),
		detectors.NewNodeVersionDetector(),
		detectors.NewPackageJsonDetector(),
		detectors.NewDockerfileDetector(releasesClient),
	}

	// Create all version managers
	// Open/Closed Principle: Adding new managers doesn't require modifying existing code
	managersList := []core.VersionManager{
		managers.NewNvmManager(shell),
		managers.NewNvsManager(shell),
		managers.NewVoltaManager(shell),
	}

	// Create all profile detectors
	// Open/Closed Principle: Adding new detectors doesn't require modifying existing code
	profileDetectorsList := []core.ProfileDetector{
		detectors.NewAutonodeYmlProfileDetector(),
		detectors.NewPackageJsonProfileDetector(),
	}

	// Create all profile switchers
	// Open/Closed Principle: Adding new switchers doesn't require modifying existing code
	profileSwitchersList := []core.ProfileSwitcher{
		switchers.NewNpmrcSwitcher(shell),
		switchers.NewTsNpmrcSwitcher(shell),
		switchers.NewRcManagerSwitcher(shell),
	}

	// Create the main service with all dependencies injected
	// Dependency Inversion Principle: Service depends on abstractions (interfaces)
	service := core.NewAutoNodeService(logger, detectorsList, managersList, profileDetectorsList, profileSwitchersList)

	// Run the service
	return service.Run(config)
}
