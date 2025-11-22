package detectors

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/matutetandil/autonode/internal/core"
)

// NvmrcDetector detects Node.js version from .nvmrc file
// Single Responsibility Principle: Only responsible for detecting version from .nvmrc
// Open/Closed Principle: Implements VersionDetector interface, can be added without modifying existing code
// Liskov Substitution Principle: Can be used anywhere a VersionDetector is expected
type NvmrcDetector struct{}

// NewNvmrcDetector creates a new NvmrcDetector instance
func NewNvmrcDetector() *NvmrcDetector {
	return &NvmrcDetector{}
}

// Detect reads the .nvmrc file and returns the version
func (d *NvmrcDetector) Detect(projectPath string) (core.DetectionResult, error) {
	nvmrcPath := filepath.Join(projectPath, ".nvmrc")

	// Check if file exists
	if _, err := os.Stat(nvmrcPath); os.IsNotExist(err) {
		return core.DetectionResult{Found: false}, nil
	}

	// Read file content
	content, err := os.ReadFile(nvmrcPath)
	if err != nil {
		return core.DetectionResult{Found: false}, err
	}

	version := strings.TrimSpace(string(content))
	if version == "" {
		return core.DetectionResult{Found: false}, nil
	}

	return core.DetectionResult{
		Found:   true,
		Version: version,
		Source:  ".nvmrc",
	}, nil
}

// GetPriority returns the priority of this detector (1 = second priority, after .autonode.yml)
func (d *NvmrcDetector) GetPriority() int {
	return 1
}

// GetSourceName returns the name of the version source
func (d *NvmrcDetector) GetSourceName() string {
	return ".nvmrc"
}
