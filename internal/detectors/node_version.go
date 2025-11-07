package detectors

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/matutetandil/autonode/internal/core"
)

// NodeVersionDetector detects Node.js version from .node-version file
// Single Responsibility Principle: Only responsible for detecting version from .node-version
// Open/Closed Principle: Implements VersionDetector interface
// Liskov Substitution Principle: Can be used anywhere a VersionDetector is expected
type NodeVersionDetector struct{}

// NewNodeVersionDetector creates a new NodeVersionDetector instance
func NewNodeVersionDetector() *NodeVersionDetector {
	return &NodeVersionDetector{}
}

// Detect reads the .node-version file and returns the version
func (d *NodeVersionDetector) Detect(projectPath string) (core.DetectionResult, error) {
	nodeVersionPath := filepath.Join(projectPath, ".node-version")

	// Check if file exists
	if _, err := os.Stat(nodeVersionPath); os.IsNotExist(err) {
		return core.DetectionResult{Found: false}, nil
	}

	// Read file content
	content, err := os.ReadFile(nodeVersionPath)
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
		Source:  ".node-version",
	}, nil
}

// GetPriority returns the priority of this detector (2 = second priority)
func (d *NodeVersionDetector) GetPriority() int {
	return 2
}

// GetSourceName returns the name of the version source
func (d *NodeVersionDetector) GetSourceName() string {
	return ".node-version"
}
