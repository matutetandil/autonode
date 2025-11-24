package core

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestUpdateChecker_IsNewerVersion(t *testing.T) {
	cache, _ := NewCacheManager()
	checker := NewUpdateChecker(cache, "0.5.0")

	tests := []struct {
		name     string
		latest   string
		current  string
		expected bool
	}{
		{
			name:     "newer version available",
			latest:   "0.6.0",
			current:  "0.5.0",
			expected: true,
		},
		{
			name:     "same version",
			latest:   "0.5.0",
			current:  "0.5.0",
			expected: false,
		},
		{
			name:     "older version (no update)",
			latest:   "0.4.0",
			current:  "0.5.0",
			expected: false,
		},
		{
			name:     "version with v prefix in latest",
			latest:   "v0.6.0",
			current:  "0.5.0",
			expected: true,
		},
		{
			name:     "version with v prefix in current",
			latest:   "0.6.0",
			current:  "v0.5.0",
			expected: true,
		},
		{
			name:     "both with v prefix",
			latest:   "v0.6.0",
			current:  "v0.5.0",
			expected: true,
		},
		{
			name:     "major version update",
			latest:   "1.0.0",
			current:  "0.9.9",
			expected: true,
		},
		{
			name:     "patch version update",
			latest:   "0.5.1",
			current:  "0.5.0",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := checker.isNewerVersion(tt.latest, tt.current)
			if result != tt.expected {
				t.Errorf("isNewerVersion(%q, %q) = %v, want %v",
					tt.latest, tt.current, result, tt.expected)
			}
		})
	}
}

func TestUpdateChecker_SetDisabled(t *testing.T) {
	cache, _ := NewCacheManager()
	checker := NewUpdateChecker(cache, "0.5.0")

	// Initially not disabled
	if checker.disabled {
		t.Error("UpdateChecker should not be disabled by default")
	}

	// Disable
	checker.SetDisabled(true)
	if !checker.disabled {
		t.Error("UpdateChecker should be disabled after SetDisabled(true)")
	}

	// Enable again
	checker.SetDisabled(false)
	if checker.disabled {
		t.Error("UpdateChecker should not be disabled after SetDisabled(false)")
	}
}

func TestUpdateChecker_SetCheckInterval(t *testing.T) {
	cache, _ := NewCacheManager()
	checker := NewUpdateChecker(cache, "0.5.0")

	// Default interval
	if checker.checkInterval != UpdateCheckInterval {
		t.Errorf("Default interval = %v, want %v", checker.checkInterval, UpdateCheckInterval)
	}

	// Custom interval
	customInterval := 24 * time.Hour
	checker.SetCheckInterval(customInterval)
	if checker.checkInterval != customInterval {
		t.Errorf("Custom interval = %v, want %v", checker.checkInterval, customInterval)
	}
}

func TestUpdateChecker_CacheReading(t *testing.T) {
	// Create temp directory for cache
	tmpDir, err := os.MkdirTemp("", "autonode-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a mock cache manager that uses the temp dir
	cache := &CacheManager{cacheDir: tmpDir}

	// Create a cached result
	cachedResult := &UpdateCheckResult{
		LastCheck:       time.Now(),
		LatestVersion:   "v0.8.0",
		CurrentVersion:  "0.5.0",
		UpdateAvailable: true,
	}

	// Write to cache
	data, _ := json.MarshalIndent(cachedResult, "", "  ")
	cachePath := filepath.Join(tmpDir, UpdateCheckCacheFile)
	if err := os.WriteFile(cachePath, data, 0644); err != nil {
		t.Fatalf("Failed to write cache: %v", err)
	}

	// Create checker and check from cache
	checker := NewUpdateChecker(cache, "0.5.0")
	checker.checkForUpdates()

	result := checker.result
	if result == nil {
		t.Fatal("Expected result to be set from cache")
	}

	if result.LatestVersion != "v0.8.0" {
		t.Errorf("LatestVersion = %q, want %q", result.LatestVersion, "v0.8.0")
	}

	if !result.UpdateAvailable {
		t.Error("Expected UpdateAvailable to be true")
	}
}

func TestUpdateChecker_StartAsyncCheckWhenDisabled(t *testing.T) {
	cache, _ := NewCacheManager()
	checker := NewUpdateChecker(cache, "0.5.0")
	checker.SetDisabled(true)

	// Start should immediately close the done channel
	checker.StartAsyncCheck()

	// Should not block
	select {
	case <-checker.done:
		// Success - channel was closed
	case <-time.After(100 * time.Millisecond):
		t.Error("StartAsyncCheck did not close done channel when disabled")
	}

	// Result should be nil when disabled
	result := checker.GetResult()
	if result != nil {
		t.Error("Expected nil result when checker is disabled")
	}
}

func TestUpdateCheckResult_Serialization(t *testing.T) {
	result := &UpdateCheckResult{
		LastCheck:       time.Date(2025, 11, 23, 10, 0, 0, 0, time.UTC),
		LatestVersion:   "v0.7.0",
		CurrentVersion:  "0.6.0",
		UpdateAvailable: true,
	}

	// Serialize
	data, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	// Deserialize
	var decoded UpdateCheckResult
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if decoded.LatestVersion != result.LatestVersion {
		t.Errorf("LatestVersion = %q, want %q", decoded.LatestVersion, result.LatestVersion)
	}

	if decoded.CurrentVersion != result.CurrentVersion {
		t.Errorf("CurrentVersion = %q, want %q", decoded.CurrentVersion, result.CurrentVersion)
	}

	if decoded.UpdateAvailable != result.UpdateAvailable {
		t.Errorf("UpdateAvailable = %v, want %v", decoded.UpdateAvailable, result.UpdateAvailable)
	}
}
