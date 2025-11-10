package switchers

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/matutetandil/autonode/internal/core"
)

// NpmrcSwitcher manages npm profiles using the npmrc tool.
// Repository: https://github.com/deoxxa/npmrc
//
// This implementation adheres to:
// - Single Responsibility Principle (SRP): Only handles npmrc profile switching
// - Dependency Inversion Principle (DIP): Depends on ShellExecutor abstraction
// - Liskov Substitution Principle (LSP): Implements ProfileSwitcher interface
type NpmrcSwitcher struct {
	shell core.ShellExecutor
}

// NewNpmrcSwitcher creates a new NpmrcSwitcher instance.
// Follows Dependency Injection pattern (DIP).
func NewNpmrcSwitcher(shell core.ShellExecutor) *NpmrcSwitcher {
	return &NpmrcSwitcher{
		shell: shell,
	}
}

// GetName returns the name of this profile switcher.
func (s *NpmrcSwitcher) GetName() string {
	return "npmrc"
}

// IsInstalled checks if npmrc is installed and available in the system.
func (s *NpmrcSwitcher) IsInstalled() bool {
	_, err := s.findExecutable()
	return err == nil
}

// findExecutable locates the npmrc executable in multiple locations.
// It searches in:
// 1. Current PATH (using shell.CommandExists)
// 2. All Node.js versions installed via nvm (~/.nvm/versions/node/*/bin/npmrc)
//
// This allows npmrc installed with one Node version to work with projects
// using different Node versions, since npmrc only modifies config files.
func (s *NpmrcSwitcher) findExecutable() (string, error) {
	// Try current PATH first
	if s.shell.CommandExists("npmrc") {
		return "npmrc", nil
	}

	// Search in nvm installations
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not get home directory: %w", err)
	}

	nvmDir := filepath.Join(homeDir, ".nvm", "versions", "node")
	pattern := filepath.Join(nvmDir, "*", "bin", "npmrc")

	matches, err := filepath.Glob(pattern)
	if err != nil {
		return "", fmt.Errorf("failed to search for npmrc: %w", err)
	}

	if len(matches) > 0 {
		// Return the first match (any version will work)
		return matches[0], nil
	}

	return "", fmt.Errorf("npmrc not found in PATH or nvm installations")
}

// ProfileExists checks if the specified profile exists in npmrc.
// It executes 'npmrc' (without arguments) to list all profiles
// and searches for the profile name in the output.
func (s *NpmrcSwitcher) ProfileExists(profileName string) (bool, error) {
	npmrcPath, err := s.findExecutable()
	if err != nil {
		return false, fmt.Errorf("failed to find npmrc: %w", err)
	}

	output, err := s.shell.Execute(npmrcPath)
	if err != nil {
		return false, fmt.Errorf("failed to list npmrc profiles: %w", err)
	}

	// Parse output: each line contains a profile name
	// Format: "   * work" or "   default"
	// Active profile is marked with an asterisk
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		// First trim whitespace, then remove asterisk if present
		profileLine := strings.TrimSpace(line)
		profileLine = strings.TrimSpace(strings.TrimPrefix(profileLine, "*"))
		if profileLine == profileName {
			return true, nil
		}
	}

	return false, nil
}

// SwitchProfile switches to the specified npm profile using npmrc.
// Command: npmrc <profile-name>
func (s *NpmrcSwitcher) SwitchProfile(profileName string) error {
	npmrcPath, err := s.findExecutable()
	if err != nil {
		return fmt.Errorf("failed to find npmrc: %w", err)
	}

	_, err = s.shell.Execute(npmrcPath, profileName)
	if err != nil {
		return fmt.Errorf("failed to switch to profile '%s' using npmrc: %w", profileName, err)
	}
	return nil
}
