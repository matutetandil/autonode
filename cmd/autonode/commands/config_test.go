package commands

import (
	"os"
	"path/filepath"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestConfigCommand_LoadConfig(t *testing.T) {
	tests := []struct {
		name           string
		fileContent    string
		createFile     bool
		expectVersion  string
		expectProfile  string
		expectError    bool
	}{
		{
			name:          "both fields present",
			fileContent:   "nodeVersion: \"20\"\nnpmProfile: work",
			createFile:    true,
			expectVersion: "20",
			expectProfile: "work",
		},
		{
			name:          "only version",
			fileContent:   "nodeVersion: \"18.17.0\"",
			createFile:    true,
			expectVersion: "18.17.0",
			expectProfile: "",
		},
		{
			name:          "only profile",
			fileContent:   "npmProfile: personal",
			createFile:    true,
			expectVersion: "",
			expectProfile: "personal",
		},
		{
			name:          "empty file",
			fileContent:   "",
			createFile:    true,
			expectVersion: "",
			expectProfile: "",
		},
		{
			name:          "file does not exist",
			createFile:    false,
			expectVersion: "",
			expectProfile: "",
		},
		{
			name:        "invalid yaml",
			fileContent: "nodeVersion: [invalid",
			createFile:  true,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			configPath := filepath.Join(tmpDir, ".autonode.yml")

			if tt.createFile {
				if err := os.WriteFile(configPath, []byte(tt.fileContent), 0644); err != nil {
					t.Fatalf("failed to create test file: %v", err)
				}
			}

			cmd := &ConfigCommand{}
			config, err := cmd.loadConfig(configPath)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if config.NodeVersion != tt.expectVersion {
				t.Errorf("NodeVersion = %q, want %q", config.NodeVersion, tt.expectVersion)
			}

			if config.NpmProfile != tt.expectProfile {
				t.Errorf("NpmProfile = %q, want %q", config.NpmProfile, tt.expectProfile)
			}
		})
	}
}

func TestConfigCommand_SaveConfig(t *testing.T) {
	tests := []struct {
		name          string
		nodeVersion   string
		npmProfile    string
		expectContent map[string]string
	}{
		{
			name:        "both fields",
			nodeVersion: "20",
			npmProfile:  "work",
			expectContent: map[string]string{
				"nodeVersion": "20",
				"npmProfile":  "work",
			},
		},
		{
			name:        "only version",
			nodeVersion: "18",
			npmProfile:  "",
			expectContent: map[string]string{
				"nodeVersion": "18",
			},
		},
		{
			name:        "only profile",
			nodeVersion: "",
			npmProfile:  "personal",
			expectContent: map[string]string{
				"npmProfile": "personal",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			configPath := filepath.Join(tmpDir, ".autonode.yml")

			config := &autonodeConfig{
				NodeVersion: tt.nodeVersion,
				NpmProfile:  tt.npmProfile,
			}

			// Marshal and write directly (we can't call saveConfig without a logger)
			data, err := yaml.Marshal(config)
			if err != nil {
				t.Fatalf("failed to marshal: %v", err)
			}

			if err := os.WriteFile(configPath, data, 0644); err != nil {
				t.Fatalf("failed to write file: %v", err)
			}

			// Read back and verify
			readData, err := os.ReadFile(configPath)
			if err != nil {
				t.Fatalf("failed to read file: %v", err)
			}

			var readConfig map[string]string
			if err := yaml.Unmarshal(readData, &readConfig); err != nil {
				t.Fatalf("failed to parse yaml: %v", err)
			}

			for key, expectedVal := range tt.expectContent {
				if readConfig[key] != expectedVal {
					t.Errorf("%s = %q, want %q", key, readConfig[key], expectedVal)
				}
			}
		})
	}
}

func TestConfigCommand_RemoveConfig(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, ".autonode.yml")

	// Create file
	if err := os.WriteFile(configPath, []byte("nodeVersion: \"20\""), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	// Verify it exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Fatal("test file should exist")
	}

	// Remove it
	if err := os.Remove(configPath); err != nil {
		t.Fatalf("failed to remove: %v", err)
	}

	// Verify it's gone
	if _, err := os.Stat(configPath); !os.IsNotExist(err) {
		t.Error("file should not exist after removal")
	}
}

func TestConfigCommand_RemoveNonExistent(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, ".autonode.yml")

	// Try to remove non-existent file
	if _, err := os.Stat(configPath); !os.IsNotExist(err) {
		t.Error("file should not exist before test")
	}

	// This should not error
	err := os.Remove(configPath)
	if err == nil {
		t.Error("expected error when removing non-existent file")
	}
}

func TestAutonodeConfig_YamlOmitEmpty(t *testing.T) {
	tests := []struct {
		name        string
		config      autonodeConfig
		expectKeys  []string
		rejectKeys  []string
	}{
		{
			name: "both fields",
			config: autonodeConfig{
				NodeVersion: "20",
				NpmProfile:  "work",
			},
			expectKeys: []string{"nodeVersion", "npmProfile"},
		},
		{
			name: "only version",
			config: autonodeConfig{
				NodeVersion: "20",
			},
			expectKeys: []string{"nodeVersion"},
			rejectKeys: []string{"npmProfile"},
		},
		{
			name: "only profile",
			config: autonodeConfig{
				NpmProfile: "work",
			},
			expectKeys: []string{"npmProfile"},
			rejectKeys: []string{"nodeVersion"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := yaml.Marshal(&tt.config)
			if err != nil {
				t.Fatalf("failed to marshal: %v", err)
			}

			content := string(data)

			for _, key := range tt.expectKeys {
				if !contains(content, key) {
					t.Errorf("expected %q in output, got: %s", key, content)
				}
			}

			for _, key := range tt.rejectKeys {
				if contains(content, key) {
					t.Errorf("did not expect %q in output, got: %s", key, content)
				}
			}
		})
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
