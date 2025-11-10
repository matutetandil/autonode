package switchers

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/matutetandil/autonode/internal/core"
)

// TsNpmrcSwitcher manages npm profiles using the ts-npmrc tool.
// Repository: https://github.com/darsi-an/ts-npmrc
//
// This implementation adheres to:
// - Single Responsibility Principle (SRP): Only handles ts-npmrc profile switching
// - Dependency Inversion Principle (DIP): Depends on ShellExecutor abstraction
// - Liskov Substitution Principle (LSP): Implements ProfileSwitcher interface
type TsNpmrcSwitcher struct {
	shell core.ShellExecutor
}

// NewTsNpmrcSwitcher creates a new TsNpmrcSwitcher instance.
// Follows Dependency Injection pattern (DIP).
func NewTsNpmrcSwitcher(shell core.ShellExecutor) *TsNpmrcSwitcher {
	return &TsNpmrcSwitcher{
		shell: shell,
	}
}

// GetName returns the name of this profile switcher.
func (s *TsNpmrcSwitcher) GetName() string {
	return "ts-npmrc"
}

// IsInstalled checks if ts-npmrc is installed and available in the system.
func (s *TsNpmrcSwitcher) IsInstalled() bool {
	_, err := s.findExecutable()
	return err == nil
}

// findExecutable locates the ts-npmrc executable in multiple locations.
// It searches in:
// 1. Current PATH (using shell.CommandExists)
// 2. All Node.js versions installed via nvm (~/.nvm/versions/node/*/bin/ts-npmrc)
//
// This allows ts-npmrc installed with one Node version to work with projects
// using different Node versions, since ts-npmrc only modifies config files.
func (s *TsNpmrcSwitcher) findExecutable() (string, error) {
	// Try current PATH first
	if s.shell.CommandExists("ts-npmrc") {
		return "ts-npmrc", nil
	}

	// Search in nvm installations
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not get home directory: %w", err)
	}

	nvmDir := filepath.Join(homeDir, ".nvm", "versions", "node")
	pattern := filepath.Join(nvmDir, "*", "bin", "ts-npmrc")

	matches, err := filepath.Glob(pattern)
	if err != nil {
		return "", fmt.Errorf("failed to search for ts-npmrc: %w", err)
	}

	if len(matches) > 0 {
		// Return the first match (any version will work)
		return matches[0], nil
	}

	return "", fmt.Errorf("ts-npmrc not found in PATH or nvm installations")
}

// ProfileExists checks if the specified profile exists in ts-npmrc.
// It executes 'ts-npmrc list' to list all profiles
// and searches for the profile name in the output.
func (s *TsNpmrcSwitcher) ProfileExists(profileName string) (bool, error) {
	tsNpmrcPath, err := s.findExecutable()
	if err != nil {
		return false, fmt.Errorf("failed to find ts-npmrc: %w", err)
	}

	output, err := s.shell.Execute(tsNpmrcPath, "list")
	if err != nil {
		return false, fmt.Errorf("failed to list ts-npmrc profiles: %w", err)
	}

	// Parse output: each line contains a profile name
	// Look for the profile in the list
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		profileLine := strings.TrimSpace(line)
		// ts-npmrc list output includes markers, so we check if line contains the profile
		if strings.Contains(profileLine, profileName) {
			return true, nil
		}
	}

	return false, nil
}

// SwitchProfile switches to the specified npm profile using ts-npmrc.
// Command: ts-npmrc link -p <profile-name>
// Note: "link" in ts-npmrc is synonymous with "switch" in npmrc
func (s *TsNpmrcSwitcher) SwitchProfile(profileName string) error {
	tsNpmrcPath, err := s.findExecutable()
	if err != nil {
		return fmt.Errorf("failed to find ts-npmrc: %w", err)
	}

	_, err = s.shell.Execute(tsNpmrcPath, "link", "-p", profileName)
	if err != nil {
		return fmt.Errorf("failed to switch to profile '%s' using ts-npmrc: %w", profileName, err)
	}
	return nil
}
