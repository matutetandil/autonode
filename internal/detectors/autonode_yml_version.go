package detectors

import (
	"os"
	"path/filepath"

	"github.com/matutetandil/autonode/internal/core"
	"gopkg.in/yaml.v3"
)

// AutonodeYmlVersionDetector detects Node.js version from .autonode.yml file.
//
// This detector adheres to:
// - Single Responsibility Principle (SRP): Only handles .autonode.yml version detection
// - Open/Closed Principle (OCP): Part of an extensible detection system
// - Liskov Substitution Principle (LSP): Implements VersionDetector interface
type AutonodeYmlVersionDetector struct{}

// autonodeYmlVersionConfig represents the structure of .autonode.yml file for version detection
type autonodeYmlVersionConfig struct {
	NodeVersion string `yaml:"nodeVersion"`
}

// NewAutonodeYmlVersionDetector creates a new AutonodeYmlVersionDetector instance.
func NewAutonodeYmlVersionDetector() *AutonodeYmlVersionDetector {
	return &AutonodeYmlVersionDetector{}
}

// Detect searches for .autonode.yml file in the project directory
// and extracts the Node.js version configuration.
func (d *AutonodeYmlVersionDetector) Detect(projectPath string) (core.DetectionResult, error) {
	filePath := filepath.Join(projectPath, ".autonode.yml")

	// Check if file exists
	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return core.DetectionResult{Found: false}, nil
		}
		return core.DetectionResult{Found: false}, err
	}

	// Parse YAML
	var config autonodeYmlVersionConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return core.DetectionResult{Found: false}, err
	}

	// Check if nodeVersion is specified
	if config.NodeVersion == "" {
		return core.DetectionResult{Found: false}, nil
	}

	return core.DetectionResult{
		Found:   true,
		Version: config.NodeVersion,
		Source:  ".autonode.yml",
	}, nil
}

// GetPriority returns the priority of this detector.
// Priority 0 means highest priority (checked first, before .nvmrc).
func (d *AutonodeYmlVersionDetector) GetPriority() int {
	return 0
}

// GetSourceName returns a human-readable name of the source.
func (d *AutonodeYmlVersionDetector) GetSourceName() string {
	return ".autonode.yml"
}
