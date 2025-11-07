package detectors

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNodeVersionDetector_Detect(t *testing.T) {
	detector := NewNodeVersionDetector()

	tests := []struct {
		name        string
		fileContent string
		wantFound   bool
		wantVersion string
	}{
		{
			name:        "valid version",
			fileContent: "18.17.0",
			wantFound:   true,
			wantVersion: "18.17.0",
		},
		{
			name:        "version with whitespace",
			fileContent: "  16.20.0  \n",
			wantFound:   true,
			wantVersion: "16.20.0",
		},
		{
			name:        "version with v prefix",
			fileContent: "v20.0.0",
			wantFound:   true,
			wantVersion: "v20.0.0",
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
			filePath := filepath.Join(tmpDir, ".node-version")

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

func TestNodeVersionDetector_GetPriority(t *testing.T) {
	detector := NewNodeVersionDetector()

	priority := detector.GetPriority()
	if priority != 2 {
		t.Errorf("GetPriority() = %d, want 2", priority)
	}
}
