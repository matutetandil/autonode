package detectors

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPackageJsonDetector_Detect(t *testing.T) {
	detector := NewPackageJsonDetector()

	tests := []struct {
		name        string
		fileContent string
		wantFound   bool
		wantVersion string
	}{
		{
			name: "exact version",
			fileContent: `{
				"engines": {
					"node": "18.17.0"
				}
			}`,
			wantFound:   true,
			wantVersion: "18.17.0",
		},
		{
			name: "version with >= operator",
			fileContent: `{
				"engines": {
					"node": ">=16.0.0"
				}
			}`,
			wantFound:   true,
			wantVersion: "16.0.0",
		},
		{
			name: "version with caret",
			fileContent: `{
				"engines": {
					"node": "^18.0.0"
				}
			}`,
			wantFound:   true,
			wantVersion: "18.0.0",
		},
		{
			name: "version with tilde",
			fileContent: `{
				"engines": {
					"node": "~16.20.0"
				}
			}`,
			wantFound:   true,
			wantVersion: "16.20.0",
		},
		{
			name: "version range",
			fileContent: `{
				"engines": {
					"node": "16.0.0 - 18.0.0"
				}
			}`,
			wantFound:   true,
			wantVersion: "16.0.0",
		},
		{
			name: "version with OR",
			fileContent: `{
				"engines": {
					"node": "16.0.0 || 18.0.0"
				}
			}`,
			wantFound:   true,
			wantVersion: "16.0.0",
		},
		{
			name: "no engines field",
			fileContent: `{
				"name": "test-package"
			}`,
			wantFound: false,
		},
		{
			name: "engines without node",
			fileContent: `{
				"engines": {
					"npm": "8.0.0"
				}
			}`,
			wantFound: false,
		},
		{
			name:        "invalid JSON",
			fileContent: `{invalid json`,
			wantFound:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			filePath := filepath.Join(tmpDir, "package.json")

			err := os.WriteFile(filePath, []byte(tt.fileContent), 0644)
			if err != nil {
				t.Fatalf("failed to write test file: %v", err)
			}

			result, err := detector.Detect(tmpDir)

			// For invalid JSON, we expect an error to be handled gracefully
			if tt.name == "invalid JSON" {
				if result.Found {
					t.Errorf("Found = true for invalid JSON, want false")
				}
				return
			}

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

func TestPackageJsonDetector_GetPriority(t *testing.T) {
	detector := NewPackageJsonDetector()

	priority := detector.GetPriority()
	if priority != 3 {
		t.Errorf("GetPriority() = %d, want 3", priority)
	}
}

func TestCleanVersionSpecifier(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"exact version", "18.17.0", "18.17.0"},
		{"caret", "^18.17.0", "18.17.0"},
		{"tilde", "~18.17.0", "18.17.0"},
		{"gte", ">=18.17.0", "18.17.0"},
		{"lte", "<=18.17.0", "18.17.0"},
		{"gt", ">18.17.0", "18.17.0"},
		{"lt", "<18.17.0", "18.17.0"},
		{"equals", "=18.17.0", "18.17.0"},
		{"range", "16.0.0 - 18.0.0", "16.0.0"},
		{"or", "16.0.0 || 18.0.0", "16.0.0"},
		{"with spaces", "  >= 18.17.0  ", "18.17.0"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := cleanVersionSpecifier(tt.input)
			if got != tt.want {
				t.Errorf("cleanVersionSpecifier(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
