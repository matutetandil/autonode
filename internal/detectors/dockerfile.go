package detectors

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/matutetandil/autonode/internal/core"
)

// releasesClient interface for dependency injection
// This allows us to mock the client in tests
type releasesClient interface {
	GetVersionForCodename(codename string) (string, error)
}

// DockerfileDetector detects Node.js version from Dockerfile FROM node:X instruction
// Single Responsibility Principle: Only responsible for detecting version from Dockerfile
// Open/Closed Principle: Implements VersionDetector interface
// Liskov Substitution Principle: Can be used anywhere a VersionDetector is expected
// Dependency Inversion Principle: Depends on releasesClient interface abstraction
type DockerfileDetector struct {
	releasesClient releasesClient
}

// NewDockerfileDetector creates a new DockerfileDetector instance
func NewDockerfileDetector(releasesClient *core.NodeReleasesClient) *DockerfileDetector {
	return &DockerfileDetector{
		releasesClient: releasesClient,
	}
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

	// Regular expression to match:
	// - FROM node:18.17.0
	// - FROM node:18
	// - FROM node:iron
	// - FROM node:iron-alpine
	// - FROM node:18-alpine
	// - FROM node:lts
	// Captures tag after "node:" (before any variant like -alpine, -slim, etc.)
	re := regexp.MustCompile(`(?i)FROM\s+node:([a-z0-9][a-z0-9.]*)(?:-[a-z0-9]+)?`)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip comments
		if strings.HasPrefix(line, "#") {
			continue
		}

		matches := re.FindStringSubmatch(line)
		if len(matches) > 1 {
			tag := matches[1]
			version := d.resolveTag(tag)

			if version != "" {
				return core.DetectionResult{
					Found:   true,
					Version: version,
					Source:  "Dockerfile",
				}, nil
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return core.DetectionResult{Found: false}, err
	}

	return core.DetectionResult{Found: false}, nil
}

// resolveTag converts Docker image tags to Node.js versions
func (d *DockerfileDetector) resolveTag(tag string) string {
	tag = strings.ToLower(strings.TrimSpace(tag))

	// If it's already a numeric version, return it
	if regexp.MustCompile(`^[0-9]+(?:\.[0-9]+)?(?:\.[0-9]+)?$`).MatchString(tag) {
		return tag
	}

	// Try to resolve as LTS codename from cache/API
	if version, err := d.releasesClient.GetVersionForCodename(tag); err == nil {
		return version
	}

	// Handle special tags
	switch tag {
	case "lts":
		// Return latest LTS version (try to get from cache, fallback to hardcoded)
		if version, err := d.releasesClient.GetVersionForCodename("jod"); err == nil {
			return version
		}
		return "22" // Fallback to current LTS
	case "latest", "current":
		return "23" // Latest stable as of January 2025
	}

	// Unknown tag - return empty to signal failure
	return ""
}

// GetPriority returns the priority of this detector (4 = fifth/lowest priority)
func (d *DockerfileDetector) GetPriority() int {
	return 4
}

// GetSourceName returns the name of the version source
func (d *DockerfileDetector) GetSourceName() string {
	return "Dockerfile"
}
