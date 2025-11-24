package core

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestLoadGlobalConfig_DefaultValues(t *testing.T) {
	// Create temp directory
	tmpDir, err := os.MkdirTemp("", "autonode-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	cache := &CacheManager{cacheDir: tmpDir}

	// Load config (file doesn't exist)
	config, err := LoadGlobalConfig(cache)
	if err != nil {
		t.Fatalf("LoadGlobalConfig failed: %v", err)
	}

	// Check defaults
	if config.DisableUpdateCheck {
		t.Error("DisableUpdateCheck should default to false")
	}
	if config.UpdateCheckIntervalDays != 7 {
		t.Errorf("UpdateCheckIntervalDays = %d, want 7", config.UpdateCheckIntervalDays)
	}
}

func TestLoadGlobalConfig_ExistingConfig(t *testing.T) {
	// Create temp directory
	tmpDir, err := os.MkdirTemp("", "autonode-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Write config file
	configData := GlobalConfig{
		DisableUpdateCheck:      true,
		UpdateCheckIntervalDays: 14,
	}
	data, _ := json.MarshalIndent(configData, "", "  ")
	configPath := filepath.Join(tmpDir, GlobalConfigFile)
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		t.Fatalf("Failed to write config: %v", err)
	}

	cache := &CacheManager{cacheDir: tmpDir}

	// Load config
	config, err := LoadGlobalConfig(cache)
	if err != nil {
		t.Fatalf("LoadGlobalConfig failed: %v", err)
	}

	if !config.DisableUpdateCheck {
		t.Error("DisableUpdateCheck should be true")
	}
	if config.UpdateCheckIntervalDays != 14 {
		t.Errorf("UpdateCheckIntervalDays = %d, want 14", config.UpdateCheckIntervalDays)
	}
}

func TestLoadGlobalConfig_InvalidJSON(t *testing.T) {
	// Create temp directory
	tmpDir, err := os.MkdirTemp("", "autonode-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Write invalid JSON
	configPath := filepath.Join(tmpDir, GlobalConfigFile)
	if err := os.WriteFile(configPath, []byte("not valid json"), 0644); err != nil {
		t.Fatalf("Failed to write config: %v", err)
	}

	cache := &CacheManager{cacheDir: tmpDir}

	// Should return defaults on invalid JSON
	config, _ := LoadGlobalConfig(cache)

	// Check defaults are returned
	if config.DisableUpdateCheck {
		t.Error("Should return default DisableUpdateCheck=false on invalid JSON")
	}
	if config.UpdateCheckIntervalDays != 7 {
		t.Errorf("Should return default UpdateCheckIntervalDays=7, got %d", config.UpdateCheckIntervalDays)
	}
}

func TestSaveGlobalConfig(t *testing.T) {
	// Create temp directory
	tmpDir, err := os.MkdirTemp("", "autonode-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	cache := &CacheManager{cacheDir: tmpDir}

	config := &GlobalConfig{
		DisableUpdateCheck:      true,
		UpdateCheckIntervalDays: 30,
	}

	// Save config
	if err := SaveGlobalConfig(cache, config); err != nil {
		t.Fatalf("SaveGlobalConfig failed: %v", err)
	}

	// Verify file was created
	configPath := filepath.Join(tmpDir, GlobalConfigFile)
	data, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("Failed to read config file: %v", err)
	}

	// Verify content
	var loaded GlobalConfig
	if err := json.Unmarshal(data, &loaded); err != nil {
		t.Fatalf("Failed to parse saved config: %v", err)
	}

	if !loaded.DisableUpdateCheck {
		t.Error("Saved DisableUpdateCheck should be true")
	}
	if loaded.UpdateCheckIntervalDays != 30 {
		t.Errorf("Saved UpdateCheckIntervalDays = %d, want 30", loaded.UpdateCheckIntervalDays)
	}
}

func TestGlobalConfig_OmitEmpty(t *testing.T) {
	config := &GlobalConfig{}
	data, err := json.Marshal(config)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	// With omitempty, empty values should not appear
	jsonStr := string(data)
	if jsonStr != "{}" {
		t.Errorf("Empty config should serialize to {}, got: %s", jsonStr)
	}
}
