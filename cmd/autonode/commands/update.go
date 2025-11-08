package commands

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"encoding/json"
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

// GitHubRelease represents the latest release info from GitHub API
type GitHubRelease struct {
	TagName string `json:"tag_name"`
}

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

	// Get current version (import from main package would create circular dependency,
	// so we get it from the cobra command which has it set)
	currentVersion := cmd.Root().Version

	logger.Info("Checking for updates...")

	// Fetch latest release info from GitHub API
	latestVersion, err := c.getLatestVersion()
	if err != nil {
		return fmt.Errorf("failed to check latest version: %w", err)
	}

	// Normalize versions for comparison (remove 'v' prefix if present)
	normalizedCurrent := currentVersion
	normalizedLatest := latestVersion
	if len(normalizedLatest) > 0 && normalizedLatest[0] == 'v' {
		normalizedLatest = normalizedLatest[1:]
	}
	if len(normalizedCurrent) > 0 && normalizedCurrent[0] == 'v' {
		normalizedCurrent = normalizedCurrent[1:]
	}

	// Check if already on latest version
	if normalizedCurrent == normalizedLatest {
		logger.Success(fmt.Sprintf("You're already on the latest version (%s)", currentVersion))
		return nil
	}

	// Show update info
	logger.Info(fmt.Sprintf("Updating from %s → %s", currentVersion, latestVersion))

	// Detect current platform
	archiveName := fmt.Sprintf("autonode-%s-%s.tar.gz", runtime.GOOS, runtime.GOARCH)
	if runtime.GOOS == "windows" {
		archiveName = fmt.Sprintf("autonode-%s-%s.zip", runtime.GOOS, runtime.GOARCH)
	}

	binaryName := "autonode"
	if runtime.GOOS == "windows" {
		binaryName += ".exe"
	}

	// GitHub release URL (latest)
	downloadURL := fmt.Sprintf("https://github.com/matutetandil/autonode/releases/latest/download/%s", archiveName)

	logger.Info(fmt.Sprintf("Downloading %s...", archiveName))

	// Download the archive
	resp, err := http.Get(downloadURL)
	if err != nil {
		return fmt.Errorf("failed to download update: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download update: HTTP %d", resp.StatusCode)
	}

	// Create temporary directory
	tmpDir, err := os.MkdirTemp("", "autonode-update-*")
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	// Save archive to temp file
	tmpArchive := filepath.Join(tmpDir, archiveName)
	tmpFile, err := os.Create(tmpArchive)
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}

	_, err = io.Copy(tmpFile, resp.Body)
	tmpFile.Close()
	if err != nil {
		return fmt.Errorf("failed to write archive: %w", err)
	}

	logger.Info("Extracting update...")

	// Extract archive based on platform
	var extractedBinary string
	if runtime.GOOS == "windows" {
		// Extract zip for Windows
		extractedBinary, err = c.extractZip(tmpArchive, tmpDir, binaryName)
		if err != nil {
			return fmt.Errorf("failed to extract archive: %w", err)
		}
	} else {
		// Extract tar.gz for Unix
		extractedBinary, err = c.extractTarGz(tmpArchive, tmpDir, binaryName)
		if err != nil {
			return fmt.Errorf("failed to extract archive: %w", err)
		}
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

	// Make extracted binary executable
	if err := os.Chmod(extractedBinary, 0755); err != nil {
		return fmt.Errorf("failed to make executable: %w", err)
	}

	// Replace current binary
	logger.Info("Installing update...")

	// On Unix, we can rename directly. On Windows, we might need different approach
	if err := os.Rename(extractedBinary, exePath); err != nil {
		// If rename fails (e.g., cross-device), try copy + remove
		if err := c.copyFile(extractedBinary, exePath); err != nil {
			return fmt.Errorf("failed to install update: %w", err)
		}
	}

	logger.Success(fmt.Sprintf("AutoNode updated successfully to %s!", latestVersion))
	logger.Info("Run 'autonode --version' to verify")

	return nil
}

// getLatestVersion fetches the latest version from GitHub API
func (c *UpdateCommand) getLatestVersion() (string, error) {
	resp, err := http.Get("https://api.github.com/repos/matutetandil/autonode/releases/latest")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	var release GitHubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "", err
	}

	return release.TagName, nil
}

// extractTarGz extracts a tar.gz archive and returns the path to the binary
func (c *UpdateCommand) extractTarGz(archivePath, destDir, binaryName string) (string, error) {
	file, err := os.Open(archivePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	gzr, err := gzip.NewReader(file)
	if err != nil {
		return "", err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}

		// We're only interested in the binary file
		if header.Name != binaryName {
			continue
		}

		// Extract the binary
		targetPath := filepath.Join(destDir, header.Name)
		outFile, err := os.Create(targetPath)
		if err != nil {
			return "", err
		}

		if _, err := io.Copy(outFile, tr); err != nil {
			outFile.Close()
			return "", err
		}
		outFile.Close()

		return targetPath, nil
	}

	return "", fmt.Errorf("binary '%s' not found in archive", binaryName)
}

// extractZip extracts a zip archive and returns the path to the binary
func (c *UpdateCommand) extractZip(archivePath, destDir, binaryName string) (string, error) {
	r, err := zip.OpenReader(archivePath)
	if err != nil {
		return "", err
	}
	defer r.Close()

	for _, f := range r.File {
		// We're only interested in the binary file
		if f.Name != binaryName {
			continue
		}

		// Extract the binary
		targetPath := filepath.Join(destDir, f.Name)
		outFile, err := os.Create(targetPath)
		if err != nil {
			return "", err
		}

		rc, err := f.Open()
		if err != nil {
			outFile.Close()
			return "", err
		}

		_, err = io.Copy(outFile, rc)
		rc.Close()
		outFile.Close()

		if err != nil {
			return "", err
		}

		return targetPath, nil
	}

	return "", fmt.Errorf("binary '%s' not found in archive", binaryName)
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
