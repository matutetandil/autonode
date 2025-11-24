package core

import (
	"encoding/json"
	"os"
)

const (
	// GlobalConfigFile is the filename for global configuration
	GlobalConfigFile = "config.json"
)

// GlobalConfig represents the global autonode configuration stored in ~/.autonode/config.json
// Single Responsibility Principle: Only holds global configuration data
type GlobalConfig struct {
	// DisableUpdateCheck disables automatic update checking
	DisableUpdateCheck bool `json:"disableUpdateCheck,omitempty"`
	// UpdateCheckInterval is the interval between update checks in days (default: 7)
	UpdateCheckIntervalDays int `json:"updateCheckIntervalDays,omitempty"`
}

// LoadGlobalConfig loads the global configuration from ~/.autonode/config.json
func LoadGlobalConfig(cache *CacheManager) (*GlobalConfig, error) {
	config := &GlobalConfig{
		// Defaults
		DisableUpdateCheck:      false,
		UpdateCheckIntervalDays: 7,
	}

	err := cache.ReadCache(GlobalConfigFile, config)
	if err != nil {
		if os.IsNotExist(err) {
			return config, nil // Return defaults if file doesn't exist
		}
		return config, nil // Return defaults on any error
	}

	return config, nil
}

// SaveGlobalConfig saves the global configuration to ~/.autonode/config.json
func SaveGlobalConfig(cache *CacheManager, config *GlobalConfig) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(cache.GetCacheFilePath(GlobalConfigFile), data, 0644)
}
