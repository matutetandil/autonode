package detectors

import (
	"os"
	"path/filepath"

	"github.com/matutetandil/autonode/internal/core"
	"gopkg.in/yaml.v3"
)

// AutonodeYmlProfileDetector detects npm profile configuration from .autonode.yml file.
//
// This detector adheres to:
// - Single Responsibility Principle (SRP): Only handles .autonode.yml detection
// - Open/Closed Principle (OCP): Part of an extensible detection system
// - Liskov Substitution Principle (LSP): Implements ProfileDetector interface
type AutonodeYmlProfileDetector struct{}

// autonodeYmlConfig represents the structure of .autonode.yml file
type autonodeYmlConfig struct {
	NpmProfile string `yaml:"npmProfile"`
}

// NewAutonodeYmlProfileDetector creates a new AutonodeYmlProfileDetector instance.
func NewAutonodeYmlProfileDetector() *AutonodeYmlProfileDetector {
	return &AutonodeYmlProfileDetector{}
}

// Detect searches for .autonode.yml file in the project directory
// and extracts the npm profile configuration.
func (d *AutonodeYmlProfileDetector) Detect(projectPath string) (core.ProfileDetectionResult, error) {
	filePath := filepath.Join(projectPath, ".autonode.yml")

	// Check if file exists
	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return core.ProfileDetectionResult{Found: false}, nil
		}
		return core.ProfileDetectionResult{Found: false}, err
	}

	// Parse YAML
	var config autonodeYmlConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return core.ProfileDetectionResult{Found: false}, err
	}

	// Check if npmProfile is specified
	if config.NpmProfile == "" {
		return core.ProfileDetectionResult{Found: false}, nil
	}

	return core.ProfileDetectionResult{
		Found:       true,
		ProfileName: config.NpmProfile,
		Source:      ".autonode.yml",
	}, nil
}

// GetPriority returns the priority of this detector.
// Priority 1 means highest priority (checked first).
func (d *AutonodeYmlProfileDetector) GetPriority() int {
	return 1
}

// GetSourceName returns a human-readable name of the source.
func (d *AutonodeYmlProfileDetector) GetSourceName() string {
	return ".autonode.yml"
}
