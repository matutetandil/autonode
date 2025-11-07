package managers

import (
	"fmt"
	"strings"

	"github.com/matutetandil/autonode/internal/core"
)

// NvsManager manages Node.js versions using nvs (Node Version Switcher)
// Single Responsibility Principle: Only responsible for nvs operations
// Dependency Inversion Principle: Depends on ShellExecutor abstraction
// Open/Closed Principle: Implements VersionManager interface
type NvsManager struct {
	shell core.ShellExecutor
}

// NewNvsManager creates a new NvsManager with injected ShellExecutor
func NewNvsManager(shell core.ShellExecutor) *NvsManager {
	return &NvsManager{
		shell: shell,
	}
}

// GetName returns the name of this version manager
func (m *NvsManager) GetName() string {
	return "nvs"
}

// IsInstalled checks if nvs is installed on the system
func (m *NvsManager) IsInstalled() bool {
	return m.shell.CommandExists("nvs")
}

// IsVersionInstalled checks if a specific Node.js version is installed via nvs
func (m *NvsManager) IsVersionInstalled(version string) (bool, error) {
	// List installed versions
	output, err := m.shell.Execute("nvs", "list")
	if err != nil {
		return false, fmt.Errorf("failed to list nvs versions: %w", err)
	}

	// Check if the version appears in the list
	normalizedVersion := normalizeNvsVersion(version)
	return strings.Contains(output, normalizedVersion), nil
}

// InstallVersion installs a specific Node.js version using nvs
func (m *NvsManager) InstallVersion(version string) error {
	normalizedVersion := normalizeNvsVersion(version)
	_, err := m.shell.Execute("nvs", "add", normalizedVersion)
	if err != nil {
		return fmt.Errorf("failed to install version %s: %w", normalizedVersion, err)
	}
	return nil
}

// UseVersion switches to a specific Node.js version using nvs
func (m *NvsManager) UseVersion(version string) error {
	normalizedVersion := normalizeNvsVersion(version)
	_, err := m.shell.Execute("nvs", "use", normalizedVersion)
	if err != nil {
		return fmt.Errorf("failed to use version %s: %w", normalizedVersion, err)
	}
	return nil
}

// normalizeNvsVersion ensures version has consistent format for nvs
// nvs expects versions without 'v' prefix
func normalizeNvsVersion(version string) string {
	return strings.TrimPrefix(version, "v")
}
