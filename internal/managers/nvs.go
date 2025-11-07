package managers

import (
	"fmt"
	"os"
	"path/filepath"
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
// nvs is a shell function, not a binary, so we check for the nvs.sh script
func (m *NvsManager) IsInstalled() bool {
	nvsHome := m.getNvsHome()
	nvsScript := filepath.Join(nvsHome, "nvs.sh")

	// Check if nvs.sh exists
	if _, err := os.Stat(nvsScript); err == nil {
		return true
	}

	return false
}

// getNvsHome returns the nvs installation directory
// Checks NVS_HOME environment variable, otherwise defaults to ~/.nvs
func (m *NvsManager) getNvsHome() string {
	if nvsHome := os.Getenv("NVS_HOME"); nvsHome != "" {
		return nvsHome
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	return filepath.Join(homeDir, ".nvs")
}

// sourceNvs returns the shell command to source nvs.sh before executing nvs commands
func (m *NvsManager) sourceNvs() string {
	nvsHome := m.getNvsHome()
	nvsScript := filepath.Join(nvsHome, "nvs.sh")
	return fmt.Sprintf(". %s && ", nvsScript)
}

// IsVersionInstalled checks if a specific Node.js version is installed via nvs
func (m *NvsManager) IsVersionInstalled(version string) (bool, error) {
	// List installed versions (need to source nvs.sh first)
	command := m.sourceNvs() + "nvs list"
	output, err := m.shell.ExecuteInShell(command)
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
	command := m.sourceNvs() + fmt.Sprintf("nvs add %s", normalizedVersion)
	_, err := m.shell.ExecuteInShell(command)
	if err != nil {
		return fmt.Errorf("failed to install version %s: %w", normalizedVersion, err)
	}
	return nil
}

// UseVersion switches to a specific Node.js version using nvs
func (m *NvsManager) UseVersion(version string) error {
	normalizedVersion := normalizeNvsVersion(version)
	command := m.sourceNvs() + fmt.Sprintf("nvs use %s", normalizedVersion)
	_, err := m.shell.ExecuteInShell(command)
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
