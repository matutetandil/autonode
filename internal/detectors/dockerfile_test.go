package detectors

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDockerfileDetector_Detect(t *testing.T) {
	detector := NewDockerfileDetector()

	tests := []struct {
		name        string
		fileContent string
		wantFound   bool
		wantVersion string
	}{
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
			name: "with alpine tag",
			fileContent: `FROM node:18.17.0-alpine
RUN npm install`,
			wantFound:   true,
			wantVersion: "18.17.0",
		},
		{
			name: "with AS alias",
			fileContent: `FROM node:18.17.0 AS builder
RUN npm install`,
			wantFound:   true,
			wantVersion: "18.17.0",
		},
		{
			name: "comment lines ignored",
			fileContent: `# FROM node:99.99.99
FROM node:18.17.0
# Another comment`,
			wantFound:   true,
			wantVersion: "18.17.0",
		},
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
			name:        "node without version",
			fileContent: "FROM node:latest",
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
		})
	}
}

func TestDockerfileDetector_NoFile(t *testing.T) {
	detector := NewDockerfileDetector()
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
	detector := NewDockerfileDetector()

	priority := detector.GetPriority()
	if priority != 4 {
		t.Errorf("GetPriority() = %d, want 4 (lowest priority)", priority)
	}
}
