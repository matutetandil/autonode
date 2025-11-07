package detectors

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/matutetandil/autonode/internal/core"
)

// PackageJsonDetector detects Node.js version from package.json engines.node field
// Single Responsibility Principle: Only responsible for detecting version from package.json
// Open/Closed Principle: Implements VersionDetector interface
// Liskov Substitution Principle: Can be used anywhere a VersionDetector is expected
type PackageJsonDetector struct{}

// packageJSON represents the structure of package.json we care about
type packageJSON struct {
	Engines struct {
		Node string `json:"node"`
	} `json:"engines"`
}

// NewPackageJsonDetector creates a new PackageJsonDetector instance
func NewPackageJsonDetector() *PackageJsonDetector {
	return &PackageJsonDetector{}
}

// Detect reads the package.json file and extracts the Node.js version from engines.node
func (d *PackageJsonDetector) Detect(projectPath string) (core.DetectionResult, error) {
	packageJsonPath := filepath.Join(projectPath, "package.json")

	// Check if file exists
	if _, err := os.Stat(packageJsonPath); os.IsNotExist(err) {
		return core.DetectionResult{Found: false}, nil
	}

	// Read file content
	content, err := os.ReadFile(packageJsonPath)
	if err != nil {
		return core.DetectionResult{Found: false}, err
	}

	// Parse JSON
	var pkg packageJSON
	if err := json.Unmarshal(content, &pkg); err != nil {
		return core.DetectionResult{Found: false}, err
	}

	version := strings.TrimSpace(pkg.Engines.Node)
	if version == "" {
		return core.DetectionResult{Found: false}, nil
	}

	// Clean up version specifiers (e.g., ">=16.0.0" -> "16.0.0")
	version = cleanVersionSpecifier(version)

	return core.DetectionResult{
		Found:   true,
		Version: version,
		Source:  "package.json (engines.node)",
	}, nil
}

// GetPriority returns the priority of this detector (3 = third priority)
func (d *PackageJsonDetector) GetPriority() int {
	return 3
}

// GetSourceName returns the name of the version source
func (d *PackageJsonDetector) GetSourceName() string {
	return "package.json"
}

// cleanVersionSpecifier removes version range operators and extracts a specific version
func cleanVersionSpecifier(version string) string {
	// Trim spaces first
	version = strings.TrimSpace(version)

	// Remove common prefixes
	version = strings.TrimPrefix(version, ">=")
	version = strings.TrimPrefix(version, "<=")
	version = strings.TrimPrefix(version, ">")
	version = strings.TrimPrefix(version, "<")
	version = strings.TrimPrefix(version, "^")
	version = strings.TrimPrefix(version, "~")
	version = strings.TrimPrefix(version, "=")
	version = strings.TrimSpace(version)

	// If it's a range (e.g., "16.0.0 - 18.0.0"), take the first version
	if strings.Contains(version, " - ") {
		parts := strings.Split(version, " - ")
		version = strings.TrimSpace(parts[0])
	}

	// If it contains ||, take the first alternative
	if strings.Contains(version, "||") {
		parts := strings.Split(version, "||")
		version = strings.TrimSpace(parts[0])
		version = cleanVersionSpecifier(version) // Recursively clean
	}

	return version
}
