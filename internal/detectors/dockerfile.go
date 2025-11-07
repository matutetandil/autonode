package detectors

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/matutetandil/autonode/internal/core"
)

// DockerfileDetector detects Node.js version from Dockerfile FROM node:X instruction
// Single Responsibility Principle: Only responsible for detecting version from Dockerfile
// Open/Closed Principle: Implements VersionDetector interface
// Liskov Substitution Principle: Can be used anywhere a VersionDetector is expected
type DockerfileDetector struct{}

// NewDockerfileDetector creates a new DockerfileDetector instance
func NewDockerfileDetector() *DockerfileDetector {
	return &DockerfileDetector{}
}

// Detect reads the Dockerfile and extracts Node.js version from FROM instruction
func (d *DockerfileDetector) Detect(projectPath string) (core.DetectionResult, error) {
	dockerfilePath := filepath.Join(projectPath, "Dockerfile")

	// Check if file exists
	if _, err := os.Stat(dockerfilePath); os.IsNotExist(err) {
		return core.DetectionResult{Found: false}, nil
	}

	// Open file
	file, err := os.Open(dockerfilePath)
	if err != nil {
		return core.DetectionResult{Found: false}, err
	}
	defer file.Close()

	// Regular expression to match: FROM node:18.17.0 or FROM node:18
	// Captures version after "node:"
	re := regexp.MustCompile(`(?i)FROM\s+node:([0-9]+(?:\.[0-9]+)?(?:\.[0-9]+)?)`)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip comments
		if strings.HasPrefix(line, "#") {
			continue
		}

		matches := re.FindStringSubmatch(line)
		if len(matches) > 1 {
			version := matches[1]
			return core.DetectionResult{
				Found:   true,
				Version: version,
				Source:  "Dockerfile",
			}, nil
		}
	}

	if err := scanner.Err(); err != nil {
		return core.DetectionResult{Found: false}, err
	}

	return core.DetectionResult{Found: false}, nil
}

// GetPriority returns the priority of this detector (4 = lowest priority)
func (d *DockerfileDetector) GetPriority() int {
	return 4
}

// GetSourceName returns the name of the version source
func (d *DockerfileDetector) GetSourceName() string {
	return "Dockerfile"
}
