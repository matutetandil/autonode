package detectors

import (
	"os"
	"path/filepath"
	"testing"
)

func TestAutonodeYmlVersionDetector_Detect(t *testing.T) {
	tests := []struct {
		name         string
		fileContent  string
		createFile   bool
		expectFound  bool
		expectVer    string
		expectSource string
	}{
		{
			name:         "valid version",
			fileContent:  "nodeVersion: \"20\"",
			createFile:   true,
			expectFound:  true,
			expectVer:    "20",
			expectSource: ".autonode.yml",
		},
		{
			name:         "valid version with full semver",
			fileContent:  "nodeVersion: \"18.17.0\"",
			createFile:   true,
			expectFound:  true,
			expectVer:    "18.17.0",
			expectSource: ".autonode.yml",
		},
		{
			name:         "valid version unquoted",
			fileContent:  "nodeVersion: 20",
			createFile:   true,
			expectFound:  true,
			expectVer:    "20",
			expectSource: ".autonode.yml",
		},
		{
			name:         "version with npm profile",
			fileContent:  "nodeVersion: \"20\"\nnpmProfile: work",
			createFile:   true,
			expectFound:  true,
			expectVer:    "20",
			expectSource: ".autonode.yml",
		},
		{
			name:        "empty version",
			fileContent: "nodeVersion: \"\"",
			createFile:  true,
			expectFound: false,
		},
		{
			name:        "missing version field",
			fileContent: "npmProfile: work",
			createFile:  true,
			expectFound: false,
		},
		{
			name:        "empty file",
			fileContent: "",
			createFile:  true,
			expectFound: false,
		},
		{
			name:        "file does not exist",
			createFile:  false,
			expectFound: false,
		},
		{
			name:        "invalid yaml",
			fileContent: "nodeVersion: [invalid",
			createFile:  true,
			expectFound: false,
		},
		{
			name:         "only whitespace version",
			fileContent:  "nodeVersion: \"   \"",
			createFile:   true,
			expectFound:  true,
			expectVer:    "   ",
			expectSource: ".autonode.yml",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temp directory
			tmpDir := t.TempDir()

			// Create .autonode.yml if needed
			if tt.createFile {
				filePath := filepath.Join(tmpDir, ".autonode.yml")
				if err := os.WriteFile(filePath, []byte(tt.fileContent), 0644); err != nil {
					t.Fatalf("failed to create test file: %v", err)
				}
			}

			// Create detector and run
			detector := NewAutonodeYmlVersionDetector()
			result, err := detector.Detect(tmpDir)

			// Check for unexpected errors (invalid YAML is expected to return error)
			if err != nil && tt.name != "invalid yaml" {
				t.Errorf("unexpected error: %v", err)
			}

			// Check found status
			if result.Found != tt.expectFound {
				t.Errorf("Found = %v, want %v", result.Found, tt.expectFound)
			}

			// Check version if found
			if tt.expectFound {
				if result.Version != tt.expectVer {
					t.Errorf("Version = %q, want %q", result.Version, tt.expectVer)
				}
				if result.Source != tt.expectSource {
					t.Errorf("Source = %q, want %q", result.Source, tt.expectSource)
				}
			}
		})
	}
}

func TestAutonodeYmlVersionDetector_GetPriority(t *testing.T) {
	detector := NewAutonodeYmlVersionDetector()
	priority := detector.GetPriority()

	if priority != 0 {
		t.Errorf("GetPriority() = %d, want 0 (highest priority)", priority)
	}
}

func TestAutonodeYmlVersionDetector_GetSourceName(t *testing.T) {
	detector := NewAutonodeYmlVersionDetector()
	name := detector.GetSourceName()

	if name != ".autonode.yml" {
		t.Errorf("GetSourceName() = %q, want %q", name, ".autonode.yml")
	}
}

func TestAutonodeYmlVersionDetector_PriorityOverNvmrc(t *testing.T) {
	autonodeDetector := NewAutonodeYmlVersionDetector()
	nvmrcDetector := NewNvmrcDetector()

	if autonodeDetector.GetPriority() >= nvmrcDetector.GetPriority() {
		t.Errorf(".autonode.yml priority (%d) should be less than .nvmrc priority (%d)",
			autonodeDetector.GetPriority(), nvmrcDetector.GetPriority())
	}
}
