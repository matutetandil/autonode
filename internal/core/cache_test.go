package core

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNewCacheManager(t *testing.T) {
	// Use a temporary directory for testing
	tempHome := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempHome)
	defer os.Setenv("HOME", originalHome)

	cache, err := NewCacheManager()
	if err != nil {
		t.Fatalf("NewCacheManager() error = %v", err)
	}

	expectedDir := filepath.Join(tempHome, ".autonode")
	if cache.cacheDir != expectedDir {
		t.Errorf("cacheDir = %v, want %v", cache.cacheDir, expectedDir)
	}

	// Verify directory was created
	if _, err := os.Stat(expectedDir); os.IsNotExist(err) {
		t.Errorf("cache directory not created at %s", expectedDir)
	}
}

func TestCacheManager_GetCacheFilePath(t *testing.T) {
	tempHome := t.TempDir()
	cache := &CacheManager{cacheDir: filepath.Join(tempHome, ".autonode")}

	tests := []struct {
		name     string
		filename string
		want     string
	}{
		{
			name:     "simple filename",
			filename: "test.json",
			want:     filepath.Join(tempHome, ".autonode", "test.json"),
		},
		{
			name:     "filename with extension",
			filename: "node-releases.json",
			want:     filepath.Join(tempHome, ".autonode", "node-releases.json"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := cache.GetCacheFilePath(tt.filename)
			if got != tt.want {
				t.Errorf("GetCacheFilePath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCacheManager_WriteAndReadCache(t *testing.T) {
	tempHome := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempHome)
	defer os.Setenv("HOME", originalHome)

	cache, err := NewCacheManager()
	if err != nil {
		t.Fatalf("NewCacheManager() error = %v", err)
	}

	// Test data
	type testData struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}

	writeData := testData{
		Name:  "test",
		Value: 42,
	}

	filename := "test-cache.json"

	// Write cache
	err = cache.WriteCache(filename, writeData)
	if err != nil {
		t.Fatalf("WriteCache() error = %v", err)
	}

	// Read cache
	var readData testData
	err = cache.ReadCache(filename, &readData)
	if err != nil {
		t.Fatalf("ReadCache() error = %v", err)
	}

	// Verify data matches
	if readData.Name != writeData.Name {
		t.Errorf("Name = %v, want %v", readData.Name, writeData.Name)
	}
	if readData.Value != writeData.Value {
		t.Errorf("Value = %v, want %v", readData.Value, writeData.Value)
	}
}

func TestCacheManager_ReadCache_FileNotExists(t *testing.T) {
	tempHome := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempHome)
	defer os.Setenv("HOME", originalHome)

	cache, err := NewCacheManager()
	if err != nil {
		t.Fatalf("NewCacheManager() error = %v", err)
	}

	var data map[string]string
	err = cache.ReadCache("nonexistent.json", &data)
	if err == nil {
		t.Error("ReadCache() expected error for nonexistent file, got nil")
	}
}

func TestCacheManager_IsCacheValid(t *testing.T) {
	tempHome := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempHome)
	defer os.Setenv("HOME", originalHome)

	cache, err := NewCacheManager()
	if err != nil {
		t.Fatalf("NewCacheManager() error = %v", err)
	}

	filename := "test-validity.json"

	// Test 1: Non-existent file
	if cache.IsCacheValid(filename, 1*time.Hour) {
		t.Error("IsCacheValid() = true for non-existent file, want false")
	}

	// Test 2: Fresh file
	testData := map[string]string{"key": "value"}
	err = cache.WriteCache(filename, testData)
	if err != nil {
		t.Fatalf("WriteCache() error = %v", err)
	}

	if !cache.IsCacheValid(filename, 1*time.Hour) {
		t.Error("IsCacheValid() = false for fresh file, want true")
	}

	// Test 3: Expired file (simulate by setting maxAge to 0)
	if cache.IsCacheValid(filename, 0*time.Second) {
		t.Error("IsCacheValid() = true for expired file, want false")
	}
}

func TestCacheManager_ClearCache(t *testing.T) {
	tempHome := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempHome)
	defer os.Setenv("HOME", originalHome)

	cache, err := NewCacheManager()
	if err != nil {
		t.Fatalf("NewCacheManager() error = %v", err)
	}

	filename := "test-clear.json"

	// Create a cache file
	testData := map[string]string{"key": "value"}
	err = cache.WriteCache(filename, testData)
	if err != nil {
		t.Fatalf("WriteCache() error = %v", err)
	}

	// Verify file exists
	filePath := cache.GetCacheFilePath(filename)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Fatal("Cache file was not created")
	}

	// Clear cache
	err = cache.ClearCache(filename)
	if err != nil {
		t.Fatalf("ClearCache() error = %v", err)
	}

	// Verify file no longer exists
	if _, err := os.Stat(filePath); !os.IsNotExist(err) {
		t.Error("Cache file still exists after ClearCache()")
	}
}

func TestCacheManager_ClearCache_NonExistent(t *testing.T) {
	tempHome := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempHome)
	defer os.Setenv("HOME", originalHome)

	cache, err := NewCacheManager()
	if err != nil {
		t.Fatalf("NewCacheManager() error = %v", err)
	}

	// Clearing non-existent file should not error
	err = cache.ClearCache("nonexistent.json")
	if err != nil {
		t.Errorf("ClearCache() error = %v, want nil for non-existent file", err)
	}
}
