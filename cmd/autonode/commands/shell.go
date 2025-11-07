package commands

import (
	"fmt"
	"os"
	"sort"

	"github.com/matutetandil/autonode/internal/core"
	"github.com/matutetandil/autonode/internal/detectors"
	"github.com/matutetandil/autonode/internal/managers"
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

// run outputs shell commands for eval integration
// This is used by the shell hook for automatic version switching
func (c *ShellCommand) run(cmd *cobra.Command, args []string) error {
	// Get current working directory
	projectPath, err := os.Getwd()
	if err != nil {
		// Silent failure - just exit without output
		return nil
	}

	// Create detectors (no logger needed for silent mode)
	detectorsList := []core.VersionDetector{
		detectors.NewNvmrcDetector(),
		detectors.NewNodeVersionDetector(),
		detectors.NewPackageJsonDetector(),
		detectors.NewDockerfileDetector(),
	}

	// Sort detectors by priority
	sort.Slice(detectorsList, func(i, j int) bool {
		return detectorsList[i].GetPriority() < detectorsList[j].GetPriority()
	})

	// Detect version silently
	var detectedVersion string
	for _, detector := range detectorsList {
		result, err := detector.Detect(projectPath)
		if err != nil || !result.Found {
			continue
		}
		detectedVersion = result.Version
		break
	}

	// If no version detected, exit silently
	if detectedVersion == "" {
		return nil
	}

	// Create shell executor and managers
	shell := core.NewExecShell()
	managersList := []core.VersionManager{
		managers.NewNvmManager(shell),
		managers.NewNvsManager(shell),
		managers.NewVoltaManager(shell),
	}

	// Find installed manager
	var manager core.VersionManager
	for _, m := range managersList {
		if m.IsInstalled() {
			manager = m
			break
		}
	}

	// If no manager found, exit silently
	if manager == nil {
		return nil
	}

	// Output shell commands based on manager type
	switch manager.GetName() {
	case "nvm":
		// For nvm, output commands to source nvm.sh and use version
		fmt.Println(`export NVM_DIR="${NVM_DIR:-$HOME/.nvm}"`)
		fmt.Println(`[ -s "$NVM_DIR/nvm.sh" ] && \. "$NVM_DIR/nvm.sh"`)
		fmt.Printf("nvm use %s 2>/dev/null\n", detectedVersion)
	case "nvs":
		fmt.Printf("nvs use %s 2>/dev/null\n", detectedVersion)
	case "volta":
		// Volta doesn't need activation per directory, it reads package.json automatically
		// But we can still output a use command
		fmt.Printf("volta pin node@%s 2>/dev/null\n", detectedVersion)
	}

	return nil
}
