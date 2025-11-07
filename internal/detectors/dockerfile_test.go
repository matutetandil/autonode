package detectors

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDockerfileDetector_Detect(t *testing.T) {
	mockClient := newMockReleasesClient()
	detector := &DockerfileDetector{releasesClient: mockClient}

	tests := []struct {
		name        string
		fileContent string
		wantFound   bool
		wantVersion string
	}{
		// Numeric versions
		{
			name:        "full version",
			fileContent: "FROM node:18.17.0",
			wantFound:   true,
			wantVersion: "18.17.0",
		},
		{
			name:        "major.minor version",
			fileContent: "FROM node:18.17",
			wantFound:   true,
			wantVersion: "18.17",
		},
		{
			name:        "major version only",
			fileContent: "FROM node:18",
			wantFound:   true,
			wantVersion: "18",
		},

		// LTS Codenames
		{
			name:        "LTS codename - iron",
			fileContent: "FROM node:iron",
			wantFound:   true,
			wantVersion: "20",
		},
		{
			name:        "LTS codename - jod",
			fileContent: "FROM node:jod",
			wantFound:   true,
			wantVersion: "22",
		},
		{
			name:        "LTS codename - hydrogen",
			fileContent: "FROM node:hydrogen",
			wantFound:   true,
			wantVersion: "18",
		},
		{
			name:        "LTS codename - gallium",
			fileContent: "FROM node:gallium",
			wantFound:   true,
			wantVersion: "16",
		},
		{
			name:        "LTS codename - krypton",
			fileContent: "FROM node:krypton",
			wantFound:   true,
			wantVersion: "24",
		},

		// LTS Codenames with variants
		{
			name:        "LTS codename with alpine",
			fileContent: "FROM node:iron-alpine",
			wantFound:   true,
			wantVersion: "20",
		},
		{
			name:        "LTS codename with slim",
			fileContent: "FROM node:jod-slim",
			wantFound:   true,
			wantVersion: "22",
		},
		{
			name:        "LTS codename with bullseye",
			fileContent: "FROM node:hydrogen-bullseye",
			wantFound:   true,
			wantVersion: "18",
		},
		{
			name:        "LTS codename with complex variant",
			fileContent: "FROM node:iron-bullseye-slim",
			wantFound:   true,
			wantVersion: "20",
		},

		// Special tags
		{
			name:        "lts tag",
			fileContent: "FROM node:lts",
			wantFound:   true,
			wantVersion: "22", // Current LTS
		},
		{
			name:        "latest tag",
			fileContent: "FROM node:latest",
			wantFound:   true,
			wantVersion: "23", // Current latest
		},
		{
			name:        "current tag",
			fileContent: "FROM node:current",
			wantFound:   true,
			wantVersion: "23", // Current stable
		},

		// Special tags with variants
		{
			name:        "lts with alpine",
			fileContent: "FROM node:lts-alpine",
			wantFound:   true,
			wantVersion: "22",
		},
		{
			name:        "latest with slim",
			fileContent: "FROM node:latest-slim",
			wantFound:   true,
			wantVersion: "23",
		},

		// Numeric versions with variants
		{
			name:        "version with alpine tag",
			fileContent: "FROM node:18.17.0-alpine",
			wantFound:   true,
			wantVersion: "18.17.0",
		},
		{
			name:        "version with slim tag",
			fileContent: "FROM node:20-slim",
			wantFound:   true,
			wantVersion: "20",
		},

		// Multiline and formatting
		{
			name: "multiline with FROM",
			fileContent: `# Use Node.js
FROM node:16.20.0
RUN npm install`,
			wantFound:   true,
			wantVersion: "16.20.0",
		},
		{
			name: "case insensitive FROM",
			fileContent: `from node:20.0.0
RUN npm install`,
			wantFound:   true,
			wantVersion: "20.0.0",
		},
		{
			name: "multiline with codename",
			fileContent: `# Use Node.js LTS
FROM node:iron-alpine
RUN npm install`,
			wantFound:   true,
			wantVersion: "20",
		},
		{
			name: "with AS alias",
			fileContent: `FROM node:18.17.0 AS builder
RUN npm install`,
			wantFound:   true,
			wantVersion: "18.17.0",
		},
		{
			name: "codename with AS alias",
			fileContent: `FROM node:iron-alpine AS builder
RUN npm install`,
			wantFound:   true,
			wantVersion: "20",
		},
		{
			name: "comment lines ignored",
			fileContent: `# FROM node:99.99.99
FROM node:18.17.0
# Another comment`,
			wantFound:   true,
			wantVersion: "18.17.0",
		},

		// Negative cases
		{
			name:        "different base image",
			fileContent: "FROM ubuntu:20.04",
			wantFound:   false,
		},
		{
			name:        "no FROM instruction",
			fileContent: "RUN apt-get update",
			wantFound:   false,
		},
		{
			name:        "unknown codename",
			fileContent: "FROM node:unknowncodename",
			wantFound:   false,
		},
		{
			name:        "empty file",
			fileContent: "",
			wantFound:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			filePath := filepath.Join(tmpDir, "Dockerfile")

			err := os.WriteFile(filePath, []byte(tt.fileContent), 0644)
			if err != nil {
				t.Fatalf("failed to write test file: %v", err)
			}

			result, err := detector.Detect(tmpDir)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if result.Found != tt.wantFound {
				t.Errorf("Found = %v, want %v", result.Found, tt.wantFound)
			}

			if tt.wantFound && result.Version != tt.wantVersion {
				t.Errorf("Version = %v, want %v", result.Version, tt.wantVersion)
			}

			if tt.wantFound && result.Source != "Dockerfile" {
				t.Errorf("Source = %v, want 'Dockerfile'", result.Source)
			}
		})
	}
}

func TestDockerfileDetector_NoFile(t *testing.T) {
	mockClient := newMockReleasesClient()
	detector := &DockerfileDetector{releasesClient: mockClient}
	tmpDir := t.TempDir()

	result, err := detector.Detect(tmpDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Found {
		t.Errorf("Found = true, want false when file doesn't exist")
	}
}

func TestDockerfileDetector_GetPriority(t *testing.T) {
	mockClient := newMockReleasesClient()
	detector := &DockerfileDetector{releasesClient: mockClient}

	priority := detector.GetPriority()
	if priority != 4 {
		t.Errorf("GetPriority() = %d, want 4 (lowest priority)", priority)
	}
}

func TestDockerfileDetector_GetSourceName(t *testing.T) {
	mockClient := newMockReleasesClient()
	detector := &DockerfileDetector{releasesClient: mockClient}

	sourceName := detector.GetSourceName()
	if sourceName != "Dockerfile" {
		t.Errorf("GetSourceName() = %s, want 'Dockerfile'", sourceName)
	}
}
