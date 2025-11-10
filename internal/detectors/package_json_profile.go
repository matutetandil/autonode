package detectors

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/matutetandil/autonode/internal/core"
)

// PackageJsonProfileDetector detects npm profile configuration from package.json file.
//
// This detector adheres to:
// - Single Responsibility Principle (SRP): Only handles package.json profile detection
// - Open/Closed Principle (OCP): Part of an extensible detection system
// - Liskov Substitution Principle (LSP): Implements ProfileDetector interface
type PackageJsonProfileDetector struct{}

// packageJsonWithProfile represents the relevant fields in package.json
type packageJsonWithProfile struct {
	Autonode *autonodeConfig `json:"autonode"`
}

// autonodeConfig represents the autonode configuration in package.json
type autonodeConfig struct {
	NpmProfile string `json:"npmProfile"`
}

// NewPackageJsonProfileDetector creates a new PackageJsonProfileDetector instance.
func NewPackageJsonProfileDetector() *PackageJsonProfileDetector {
	return &PackageJsonProfileDetector{}
}

// Detect searches for package.json file in the project directory
// and extracts the npm profile from the autonode.npmProfile field.
func (d *PackageJsonProfileDetector) Detect(projectPath string) (core.ProfileDetectionResult, error) {
	filePath := filepath.Join(projectPath, "package.json")

	// Check if file exists
	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return core.ProfileDetectionResult{Found: false}, nil
		}
		return core.ProfileDetectionResult{Found: false}, err
	}

	// Parse JSON
	var pkg packageJsonWithProfile
	if err := json.Unmarshal(data, &pkg); err != nil {
		return core.ProfileDetectionResult{Found: false}, err
	}

	// Check if autonode.npmProfile is specified
	if pkg.Autonode == nil || pkg.Autonode.NpmProfile == "" {
		return core.ProfileDetectionResult{Found: false}, nil
	}

	return core.ProfileDetectionResult{
		Found:       true,
		ProfileName: pkg.Autonode.NpmProfile,
		Source:      "package.json",
	}, nil
}

// GetPriority returns the priority of this detector.
// Priority 2 means lower priority than .autonode.yml (checked after .autonode.yml).
func (d *PackageJsonProfileDetector) GetPriority() int {
	return 2
}

// GetSourceName returns a human-readable name of the source.
func (d *PackageJsonProfileDetector) GetSourceName() string {
	return "package.json"
}
