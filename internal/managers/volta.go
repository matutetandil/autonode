package managers

import (
	"fmt"
	"strings"

	"github.com/matutetandil/autonode/internal/core"
)

// VoltaManager manages Node.js versions using Volta
// Single Responsibility Principle: Only responsible for Volta operations
// Dependency Inversion Principle: Depends on ShellExecutor abstraction
// Open/Closed Principle: Implements VersionManager interface
type VoltaManager struct {
	shell core.ShellExecutor
}

// NewVoltaManager creates a new VoltaManager with injected ShellExecutor
func NewVoltaManager(shell core.ShellExecutor) *VoltaManager {
	return &VoltaManager{
		shell: shell,
	}
}

// GetName returns the name of this version manager
func (m *VoltaManager) GetName() string {
	return "volta"
}

// IsInstalled checks if Volta is installed on the system
func (m *VoltaManager) IsInstalled() bool {
	return m.shell.CommandExists("volta")
}

// IsVersionInstalled checks if a specific Node.js version is installed via Volta
func (m *VoltaManager) IsVersionInstalled(version string) (bool, error) {
	// List installed versions
	output, err := m.shell.Execute("volta", "list", "node")
	if err != nil {
		return false, fmt.Errorf("failed to list volta versions: %w", err)
	}

	// Check if the version appears in the list
	normalizedVersion := normalizeVoltaVersion(version)
	return strings.Contains(output, normalizedVersion), nil
}

// InstallVersion installs a specific Node.js version using Volta
// Note: Volta automatically installs when you use 'volta install node@version'
func (m *VoltaManager) InstallVersion(version string) error {
	normalizedVersion := normalizeVoltaVersion(version)
	_, err := m.shell.Execute("volta", "install", fmt.Sprintf("node@%s", normalizedVersion))
	if err != nil {
		return fmt.Errorf("failed to install version %s: %w", normalizedVersion, err)
	}
	return nil
}

// UseVersion switches to a specific Node.js version using Volta
// Volta pins the version to the project
func (m *VoltaManager) UseVersion(version string) error {
	normalizedVersion := normalizeVoltaVersion(version)
	_, err := m.shell.Execute("volta", "pin", fmt.Sprintf("node@%s", normalizedVersion))
	if err != nil {
		return fmt.Errorf("failed to use version %s: %w", normalizedVersion, err)
	}
	return nil
}

// normalizeVoltaVersion ensures version has consistent format for Volta
// Volta expects versions without 'v' prefix
func normalizeVoltaVersion(version string) string {
	return strings.TrimPrefix(version, "v")
}
