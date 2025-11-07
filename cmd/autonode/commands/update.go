package commands

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"

	"github.com/matutetandil/autonode/internal/core"
	"github.com/spf13/cobra"
)

// UpdateCommand implements the self-update command
// Single Responsibility Principle: Only responsible for updating the binary
type UpdateCommand struct{}

// init registers this command automatically when the package is imported
func init() {
	Register(&UpdateCommand{})
}

// GetCobraCommand returns the cobra command for this command
func (c *UpdateCommand) GetCobraCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "update",
		Short: "Update autonode to the latest version",
		Long:  `Downloads and installs the latest version of autonode from GitHub releases.`,
		RunE:  c.run,
	}
}

// run downloads and installs the latest version of autonode
func (c *UpdateCommand) run(cmd *cobra.Command, args []string) error {
	logger := core.NewConsoleLogger()

	logger.Info("Checking for updates...")

	// Detect current platform
	binaryName := fmt.Sprintf("autonode-%s-%s", runtime.GOOS, runtime.GOARCH)
	if runtime.GOOS == "windows" {
		binaryName += ".exe"
	}

	// GitHub release URL (latest)
	downloadURL := fmt.Sprintf("https://github.com/matutetandil/autonode/releases/latest/download/%s", binaryName)

	logger.Info(fmt.Sprintf("Downloading latest version for %s/%s...", runtime.GOOS, runtime.GOARCH))

	// Download the binary
	resp, err := http.Get(downloadURL)
	if err != nil {
		return fmt.Errorf("failed to download update: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download update: HTTP %d", resp.StatusCode)
	}

	// Get current executable path
	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}

	// Resolve symlinks
	exePath, err = filepath.EvalSymlinks(exePath)
	if err != nil {
		return fmt.Errorf("failed to resolve executable path: %w", err)
	}

	// Create temporary file
	tmpFile, err := os.CreateTemp("", "autonode-update-*")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	tmpPath := tmpFile.Name()
	defer os.Remove(tmpPath)

	// Copy downloaded content to temp file
	_, err = io.Copy(tmpFile, resp.Body)
	if err != nil {
		tmpFile.Close()
		return fmt.Errorf("failed to write update: %w", err)
	}
	tmpFile.Close()

	// Make executable
	if err := os.Chmod(tmpPath, 0755); err != nil {
		return fmt.Errorf("failed to make executable: %w", err)
	}

	// Replace current binary
	logger.Info("Installing update...")

	// On Unix, we can rename directly. On Windows, we might need different approach
	if err := os.Rename(tmpPath, exePath); err != nil {
		// If rename fails (e.g., cross-device), try copy + remove
		if err := c.copyFile(tmpPath, exePath); err != nil {
			return fmt.Errorf("failed to install update: %w", err)
		}
		os.Remove(tmpPath)
	}

	logger.Success("AutoNode has been updated successfully!")
	logger.Info("Run 'autonode --version' to verify")

	return nil
}

// copyFile copies a file from src to dst
func (c *UpdateCommand) copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	if _, err := io.Copy(destFile, sourceFile); err != nil {
		return err
	}

	// Copy permissions
	sourceInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	return os.Chmod(dst, sourceInfo.Mode())
}
