package switchers

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/matutetandil/autonode/internal/core"
)

// RcManagerSwitcher manages npm and yarn profiles using the rc-manager tool.
// Repository: https://github.com/Lalaluka/rc-manager
//
// This implementation adheres to:
// - Single Responsibility Principle (SRP): Only handles rc-manager profile switching
// - Dependency Inversion Principle (DIP): Depends on ShellExecutor abstraction
// - Liskov Substitution Principle (LSP): Implements ProfileSwitcher interface
type RcManagerSwitcher struct {
	shell core.ShellExecutor
}

// NewRcManagerSwitcher creates a new RcManagerSwitcher instance.
// Follows Dependency Injection pattern (DIP).
func NewRcManagerSwitcher(shell core.ShellExecutor) *RcManagerSwitcher {
	return &RcManagerSwitcher{
		shell: shell,
	}
}

// GetName returns the name of this profile switcher.
func (s *RcManagerSwitcher) GetName() string {
	return "rc-manager"
}

// IsInstalled checks if rc-manager is installed and available in the system.
func (s *RcManagerSwitcher) IsInstalled() bool {
	_, err := s.findExecutable()
	return err == nil
}

// findExecutable locates the rc-manager executable in multiple locations.
// It searches in:
// 1. Current PATH (using shell.CommandExists)
// 2. All Node.js versions installed via nvm (~/.nvm/versions/node/*/bin/rc-manager)
//
// This allows rc-manager installed with one Node version to work with projects
// using different Node versions, since rc-manager only modifies config files.
func (s *RcManagerSwitcher) findExecutable() (string, error) {
	// Try current PATH first
	if s.shell.CommandExists("rc-manager") {
		return "rc-manager", nil
	}

	// Search in nvm installations
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not get home directory: %w", err)
	}

	nvmDir := filepath.Join(homeDir, ".nvm", "versions", "node")
	pattern := filepath.Join(nvmDir, "*", "bin", "rc-manager")

	matches, err := filepath.Glob(pattern)
	if err != nil {
		return "", fmt.Errorf("failed to search for rc-manager: %w", err)
	}

	if len(matches) > 0 {
		// Return the first match (any version will work)
		return matches[0], nil
	}

	return "", fmt.Errorf("rc-manager not found in PATH or nvm installations")
}

// ProfileExists checks if the specified profile exists in rc-manager.
// It executes 'rc-manager list' to list all profiles
// and searches for the profile name in the output.
func (s *RcManagerSwitcher) ProfileExists(profileName string) (bool, error) {
	rcManagerPath, err := s.findExecutable()
	if err != nil {
		return false, fmt.Errorf("failed to find rc-manager: %w", err)
	}

	output, err := s.shell.Execute(rcManagerPath, "list")
	if err != nil {
		return false, fmt.Errorf("failed to list rc-manager profiles: %w", err)
	}

	// Parse output: each line contains a profile name
	// Look for the profile in the list
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		profileLine := strings.TrimSpace(line)
		if profileLine == profileName {
			return true, nil
		}
	}

	return false, nil
}

// SwitchProfile switches to the specified npm/yarn profile using rc-manager.
// Command: rc-manager load <profile-name>
func (s *RcManagerSwitcher) SwitchProfile(profileName string) error {
	rcManagerPath, err := s.findExecutable()
	if err != nil {
		return fmt.Errorf("failed to find rc-manager: %w", err)
	}

	_, err = s.shell.Execute(rcManagerPath, "load", profileName)
	if err != nil {
		return fmt.Errorf("failed to switch to profile '%s' using rc-manager: %w", profileName, err)
	}
	return nil
}
