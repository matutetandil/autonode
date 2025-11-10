package detectors

import (
	"os"
	"path/filepath"
	"testing"
)

func TestAutonodeYmlProfileDetector_Detect(t *testing.T) {
	detector := NewAutonodeYmlProfileDetector()

	tests := []struct {
		name            string
		fileContent     string
		wantFound       bool
		wantProfileName string
		wantSource      string
	}{
		{
			name:            "valid npm profile",
			fileContent:     "npmProfile: work",
			wantFound:       true,
			wantProfileName: "work",
			wantSource:      ".autonode.yml",
		},
		{
			name:            "profile with quotes",
			fileContent:     "npmProfile: \"personal\"",
			wantFound:       true,
			wantProfileName: "personal",
			wantSource:      ".autonode.yml",
		},
		{
			name:            "profile with extra fields",
			fileContent:     "npmProfile: work\notherField: value",
			wantFound:       true,
			wantProfileName: "work",
			wantSource:      ".autonode.yml",
		},
		{
			name:        "empty profile",
			fileContent: "npmProfile:",
			wantFound:   false,
		},
		{
			name:        "missing npmProfile field",
			fileContent: "otherField: value",
			wantFound:   false,
		},
		{
			name:        "empty file",
			fileContent: "",
			wantFound:   false,
		},
		{
			name:        "invalid yaml",
			fileContent: "npmProfile: work\n  invalid: indentation",
			wantFound:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temp directory
			tmpDir := t.TempDir()
			configPath := filepath.Join(tmpDir, ".autonode.yml")

			// Write test file
			err := os.WriteFile(configPath, []byte(tt.fileContent), 0644)
			if err != nil {
				t.Fatalf("failed to write test file: %v", err)
			}

			// Test detection
			result, err := detector.Detect(tmpDir)

			// For invalid YAML, we expect an error but not Found
			if !tt.wantFound && err != nil {
				// This is expected for invalid YAML
				return
			}

			if err != nil {
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

func TestAutonodeYmlProfileDetector_NoFile(t *testing.T) {
	detector := NewAutonodeYmlProfileDetector()
	tmpDir := t.TempDir()

	result, err := detector.Detect(tmpDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Found {
		t.Errorf("Found = true, want false when file doesn't exist")
	}
}

func TestAutonodeYmlProfileDetector_GetPriority(t *testing.T) {
	detector := NewAutonodeYmlProfileDetector()

	priority := detector.GetPriority()
	if priority != 1 {
		t.Errorf("GetPriority() = %d, want 1 (highest priority)", priority)
	}
}

func TestAutonodeYmlProfileDetector_GetSourceName(t *testing.T) {
	detector := NewAutonodeYmlProfileDetector()

	sourceName := detector.GetSourceName()
	if sourceName != ".autonode.yml" {
		t.Errorf("GetSourceName() = %q, want %q", sourceName, ".autonode.yml")
	}
}
