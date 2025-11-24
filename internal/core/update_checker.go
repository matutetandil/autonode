package core

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

const (
	// UpdateCheckCacheFile is the filename for update check cache
	UpdateCheckCacheFile = "update-check.json"
	// UpdateCheckInterval is the default interval between update checks (7 days)
	UpdateCheckInterval = 7 * 24 * time.Hour
	// UpdateCheckTimeout is the timeout for the GitHub API request
	UpdateCheckTimeout = 3 * time.Second
	// GitHubAPIURL is the GitHub releases API endpoint
	GitHubAPIURL = "https://api.github.com/repos/matutetandil/autonode/releases/latest"
)

// UpdateCheckResult contains the result of an update check
// Single Responsibility Principle: Only holds update check data
type UpdateCheckResult struct {
	LastCheck      time.Time `json:"lastCheck"`
	LatestVersion  string    `json:"latestVersion"`
	CurrentVersion string    `json:"currentVersion"`
	UpdateAvailable bool     `json:"updateAvailable"`
}

// UpdateChecker handles automatic update checking
// Single Responsibility Principle: Only responsible for checking updates
type UpdateChecker struct {
	cache          *CacheManager
	currentVersion string
	checkInterval  time.Duration
	disabled       bool
	result         *UpdateCheckResult
	mu             sync.Mutex
	done           chan struct{}
}

// NewUpdateChecker creates a new UpdateChecker instance
// Dependency Inversion Principle: Depends on CacheManager abstraction
func NewUpdateChecker(cache *CacheManager, currentVersion string) *UpdateChecker {
	return &UpdateChecker{
		cache:          cache,
		currentVersion: currentVersion,
		checkInterval:  UpdateCheckInterval,
		disabled:       false,
		done:           make(chan struct{}),
	}
}

// SetDisabled enables or disables update checking
func (u *UpdateChecker) SetDisabled(disabled bool) {
	u.disabled = disabled
}

// SetCheckInterval sets a custom check interval
func (u *UpdateChecker) SetCheckInterval(interval time.Duration) {
	u.checkInterval = interval
}

// StartAsyncCheck starts an asynchronous update check
// Returns immediately, check runs in background
func (u *UpdateChecker) StartAsyncCheck() {
	if u.disabled {
		close(u.done)
		return
	}

	go func() {
		defer close(u.done)
		u.checkForUpdates()
	}()
}

// GetResult returns the update check result (blocks until check completes or timeout)
func (u *UpdateChecker) GetResult() *UpdateCheckResult {
	// Wait for async check to complete with a short timeout
	select {
	case <-u.done:
	case <-time.After(UpdateCheckTimeout + 500*time.Millisecond):
		// Timeout waiting for result
	}

	u.mu.Lock()
	defer u.mu.Unlock()
	return u.result
}

// checkForUpdates performs the actual update check
func (u *UpdateChecker) checkForUpdates() {
	// First, try to load from cache
	var cached UpdateCheckResult
	err := u.cache.ReadCache(UpdateCheckCacheFile, &cached)

	if err == nil && u.cache.IsCacheValid(UpdateCheckCacheFile, u.checkInterval) {
		// Cache is valid, use cached result
		cached.CurrentVersion = u.currentVersion
		cached.UpdateAvailable = u.isNewerVersion(cached.LatestVersion, u.currentVersion)
		u.setResult(&cached)
		return
	}

	// Cache is invalid or missing, fetch from GitHub
	latestVersion, err := u.fetchLatestVersion()
	if err != nil {
		// Failed to fetch, try to use stale cache if available
		if cached.LatestVersion != "" {
			cached.CurrentVersion = u.currentVersion
			cached.UpdateAvailable = u.isNewerVersion(cached.LatestVersion, u.currentVersion)
			u.setResult(&cached)
		}
		return
	}

	// Create new result
	result := &UpdateCheckResult{
		LastCheck:       time.Now(),
		LatestVersion:   latestVersion,
		CurrentVersion:  u.currentVersion,
		UpdateAvailable: u.isNewerVersion(latestVersion, u.currentVersion),
	}

	// Save to cache (ignore errors, it's non-critical)
	_ = u.cache.WriteCache(UpdateCheckCacheFile, result)

	u.setResult(result)
}

// fetchLatestVersion fetches the latest version from GitHub API
func (u *UpdateChecker) fetchLatestVersion() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), UpdateCheckTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, GitHubAPIURL, nil)
	if err != nil {
		return "", err
	}

	// Set user agent (GitHub API requires it)
	req.Header.Set("User-Agent", "autonode/"+u.currentVersion)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	var release struct {
		TagName string `json:"tag_name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "", err
	}

	return release.TagName, nil
}

// isNewerVersion checks if latest is newer than current
func (u *UpdateChecker) isNewerVersion(latest, current string) bool {
	// Normalize versions (remove 'v' prefix)
	if len(latest) > 0 && latest[0] == 'v' {
		latest = latest[1:]
	}
	if len(current) > 0 && current[0] == 'v' {
		current = current[1:]
	}

	// Simple string comparison works for semver
	return latest != current && latest > current
}

// setResult safely sets the result
func (u *UpdateChecker) setResult(result *UpdateCheckResult) {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.result = result
}
