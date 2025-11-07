package detectors

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNvmrcDetector_Detect(t *testing.T) {
	detector := NewNvmrcDetector()

	tests := []struct {
		name        string
		fileContent string
		wantFound   bool
		wantVersion string
		wantSource  string
	}{
		{
			name:        "valid version",
			fileContent: "18.17.0",
			wantFound:   true,
			wantVersion: "18.17.0",
			wantSource:  ".nvmrc",
		},
		{
			name:        "version with whitespace",
			fileContent: "  16.20.0  \n",
			wantFound:   true,
			wantVersion: "16.20.0",
			wantSource:  ".nvmrc",
		},
		{
			name:        "version with v prefix",
			fileContent: "v20.0.0",
			wantFound:   true,
			wantVersion: "v20.0.0",
			wantSource:  ".nvmrc",
		},
		{
			name:        "lts alias",
			fileContent: "lts/*",
			wantFound:   true,
			wantVersion: "lts/*",
			wantSource:  ".nvmrc",
		},
		{
			name:        "empty file",
			fileContent: "",
			wantFound:   false,
		},
		{
			name:        "whitespace only",
			fileContent: "   \n  \t  ",
			wantFound:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temp directory
			tmpDir := t.TempDir()
			nvmrcPath := filepath.Join(tmpDir, ".nvmrc")

			// Write test file
			err := os.WriteFile(nvmrcPath, []byte(tt.fileContent), 0644)
			if err != nil {
				t.Fatalf("failed to write test file: %v", err)
			}

			// Test detection
			result, err := detector.Detect(tmpDir)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if result.Found != tt.wantFound {
				t.Errorf("Found = %v, want %v", result.Found, tt.wantFound)
			}

			if tt.wantFound {
				if result.Version != tt.wantVersion {
					t.Errorf("Version = %v, want %v", result.Version, tt.wantVersion)
				}
				if result.Source != tt.wantSource {
					t.Errorf("Source = %v, want %v", result.Source, tt.wantSource)
				}
			}
		})
	}
}

func TestNvmrcDetector_NoFile(t *testing.T) {
	detector := NewNvmrcDetector()
	tmpDir := t.TempDir()

	result, err := detector.Detect(tmpDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Found {
		t.Errorf("Found = true, want false when file doesn't exist")
	}
}

func TestNvmrcDetector_GetPriority(t *testing.T) {
	detector := NewNvmrcDetector()

	priority := detector.GetPriority()
	if priority != 1 {
		t.Errorf("GetPriority() = %d, want 1 (highest priority)", priority)
	}
}

func TestNvmrcDetector_GetSourceName(t *testing.T) {
	detector := NewNvmrcDetector()

	sourceName := detector.GetSourceName()
	if sourceName != ".nvmrc" {
		t.Errorf("GetSourceName() = %q, want %q", sourceName, ".nvmrc")
	}
}
