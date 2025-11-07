package managers

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/matutetandil/autonode/internal/core"
)

// NvmManager manages Node.js versions using nvm (Node Version Manager)
// Single Responsibility Principle: Only responsible for nvm operations
// Dependency Inversion Principle: Depends on ShellExecutor abstraction, not concrete implementation
// Open/Closed Principle: Implements VersionManager interface
type NvmManager struct {
	shell core.ShellExecutor
}

// NewNvmManager creates a new NvmManager with injected ShellExecutor
func NewNvmManager(shell core.ShellExecutor) *NvmManager {
	return &NvmManager{
		shell: shell,
	}
}

// GetName returns the name of this version manager
func (m *NvmManager) GetName() string {
	return "nvm"
}

// IsInstalled checks if nvm is installed on the system
// nvm is a shell function, not a binary, so we check for the nvm directory
func (m *NvmManager) IsInstalled() bool {
	nvmDir := m.getNvmDir()
	nvmScript := filepath.Join(nvmDir, "nvm.sh")

	// Check if nvm.sh exists
	if _, err := os.Stat(nvmScript); err == nil {
		return true
	}

	return false
}

// getNvmDir returns the nvm installation directory
// Checks NVM_DIR environment variable, otherwise defaults to ~/.nvm
func (m *NvmManager) getNvmDir() string {
	if nvmDir := os.Getenv("NVM_DIR"); nvmDir != "" {
		return nvmDir
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	return filepath.Join(homeDir, ".nvm")
}

// sourceNvm returns the shell command to source nvm.sh before executing nvm commands
func (m *NvmManager) sourceNvm() string {
	nvmDir := m.getNvmDir()
	nvmScript := filepath.Join(nvmDir, "nvm.sh")
	return fmt.Sprintf(". %s && ", nvmScript)
}

// IsVersionInstalled checks if a specific Node.js version is installed via nvm
func (m *NvmManager) IsVersionInstalled(version string) (bool, error) {
	// List installed versions (need to source nvm.sh first)
	command := m.sourceNvm() + "nvm list"
	output, err := m.shell.ExecuteInShell(command)
	if err != nil {
		return false, fmt.Errorf("failed to list nvm versions: %w", err)
	}

	// Check if the version appears in the list
	// nvm list output includes version numbers like "v18.17.0" or "18.17.0"
	normalizedVersion := normalizeVersion(version)
	return strings.Contains(output, normalizedVersion), nil
}

// InstallVersion installs a specific Node.js version using nvm
func (m *NvmManager) InstallVersion(version string) error {
	normalizedVersion := normalizeVersion(version)
	command := m.sourceNvm() + fmt.Sprintf("nvm install %s", normalizedVersion)
	_, err := m.shell.ExecuteInShell(command)
	if err != nil {
		return fmt.Errorf("failed to install version %s: %w", normalizedVersion, err)
	}
	return nil
}

// UseVersion switches to a specific Node.js version using nvm
func (m *NvmManager) UseVersion(version string) error {
	normalizedVersion := normalizeVersion(version)
	command := m.sourceNvm() + fmt.Sprintf("nvm use %s", normalizedVersion)
	_, err := m.shell.ExecuteInShell(command)
	if err != nil {
		return fmt.Errorf("failed to use version %s: %w", normalizedVersion, err)
	}
	return nil
}

// normalizeVersion ensures version has consistent format
// Examples: "18" -> "18", "v18.17.0" -> "18.17.0", "18.17.0" -> "18.17.0"
func normalizeVersion(version string) string {
	return strings.TrimPrefix(version, "v")
}
