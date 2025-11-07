package core

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

// CacheManager handles cache directory and file operations
// Single Responsibility Principle: Only responsible for cache management
type CacheManager struct {
	cacheDir string
}

// NewCacheManager creates a new CacheManager instance
func NewCacheManager() (*CacheManager, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	cacheDir := filepath.Join(homeDir, ".autonode")

	// Create cache directory if it doesn't exist
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return nil, err
	}

	return &CacheManager{
		cacheDir: cacheDir,
	}, nil
}

// GetCacheFilePath returns the full path to a cache file
func (c *CacheManager) GetCacheFilePath(filename string) string {
	return filepath.Join(c.cacheDir, filename)
}

// ReadCache reads and unmarshals JSON data from a cache file
func (c *CacheManager) ReadCache(filename string, v interface{}) error {
	filePath := c.GetCacheFilePath(filename)

	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, v)
}

// WriteCache marshals and writes JSON data to a cache file
func (c *CacheManager) WriteCache(filename string, v interface{}) error {
	filePath := c.GetCacheFilePath(filename)

	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, data, 0644)
}

// IsCacheValid checks if a cache file exists and is not older than maxAge
func (c *CacheManager) IsCacheValid(filename string, maxAge time.Duration) bool {
	filePath := c.GetCacheFilePath(filename)

	info, err := os.Stat(filePath)
	if err != nil {
		return false
	}

	age := time.Since(info.ModTime())
	return age < maxAge
}

// ClearCache removes a cache file
func (c *CacheManager) ClearCache(filename string) error {
	filePath := c.GetCacheFilePath(filename)
	err := os.Remove(filePath)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}
