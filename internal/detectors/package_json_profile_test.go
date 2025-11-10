package detectors

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPackageJsonProfileDetector_Detect(t *testing.T) {
	detector := NewPackageJsonProfileDetector()

	tests := []struct {
		name            string
		fileContent     string
		wantFound       bool
		wantProfileName string
		wantSource      string
	}{
		{
			name: "valid npm profile in autonode field",
			fileContent: `{
				"name": "test-project",
				"autonode": {
					"npmProfile": "work"
				}
			}`,
			wantFound:       true,
			wantProfileName: "work",
			wantSource:      "package.json",
		},
		{
			name: "profile with other autonode fields",
			fileContent: `{
				"name": "test-project",
				"autonode": {
					"npmProfile": "personal",
					"otherField": "value"
				}
			}`,
			wantFound:       true,
			wantProfileName: "personal",
			wantSource:      "package.json",
		},
		{
			name: "missing autonode field",
			fileContent: `{
				"name": "test-project",
				"version": "1.0.0"
			}`,
			wantFound: false,
		},
		{
			name: "empty autonode field",
			fileContent: `{
				"name": "test-project",
				"autonode": {}
			}`,
			wantFound: false,
		},
		{
			name: "autonode with empty npmProfile",
			fileContent: `{
				"name": "test-project",
				"autonode": {
					"npmProfile": ""
				}
			}`,
			wantFound: false,
		},
		{
			name: "autonode null",
			fileContent: `{
				"name": "test-project",
				"autonode": null
			}`,
			wantFound: false,
		},
		{
			name:        "empty package.json",
			fileContent: "{}",
			wantFound:   false,
		},
		{
			name:        "invalid json",
			fileContent: "{invalid json}",
			wantFound:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temp directory
			tmpDir := t.TempDir()
			packageJsonPath := filepath.Join(tmpDir, "package.json")

			// Write test file
			err := os.WriteFile(packageJsonPath, []byte(tt.fileContent), 0644)
			if err != nil {
				t.Fatalf("failed to write test file: %v", err)
			}

			// Test detection
			result, err := detector.Detect(tmpDir)

			// For invalid JSON, we don't expect an error to be returned,
			// just Found = false
			if err != nil && tt.wantFound {
				t.Fatalf("unexpected error: %v", err)
			}

			if result.Found != tt.wantFound {
				t.Errorf("Found = %v, want %v", result.Found, tt.wantFound)
			}

			if tt.wantFound {
				if result.ProfileName != tt.wantProfileName {
					t.Errorf("ProfileName = %v, want %v", result.ProfileName, tt.wantProfileName)
				}
				if result.Source != tt.wantSource {
					t.Errorf("Source = %v, want %v", result.Source, tt.wantSource)
				}
			}
		})
	}
}

func TestPackageJsonProfileDetector_NoFile(t *testing.T) {
	detector := NewPackageJsonProfileDetector()
	tmpDir := t.TempDir()

	result, err := detector.Detect(tmpDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Found {
		t.Errorf("Found = true, want false when file doesn't exist")
	}
}

func TestPackageJsonProfileDetector_GetPriority(t *testing.T) {
	detector := NewPackageJsonProfileDetector()

	priority := detector.GetPriority()
	if priority != 2 {
		t.Errorf("GetPriority() = %d, want 2 (lower priority than .autonode.yml)", priority)
	}
}

func TestPackageJsonProfileDetector_GetSourceName(t *testing.T) {
	detector := NewPackageJsonProfileDetector()

	sourceName := detector.GetSourceName()
	if sourceName != "package.json" {
		t.Errorf("GetSourceName() = %q, want %q", sourceName, "package.json")
	}
}
