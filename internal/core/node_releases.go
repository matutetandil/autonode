package core

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	nodeReleasesURL   = "https://nodejs.org/dist/index.json"
	cacheFileName     = "node-releases.json"
	cacheMaxAge       = 24 * time.Hour // Refresh cache daily
)

// NodeRelease represents a single Node.js release from the API
type NodeRelease struct {
	Version string      `json:"version"`
	Date    string      `json:"date"`
	LTS     interface{} `json:"lts"` // Can be false (bool) or "Codename" (string)
}

// NodeReleasesCache is the cached data structure
type NodeReleasesCache struct {
	CodenameToVersion map[string]string `json:"codename_to_version"`
	LastUpdated       time.Time         `json:"last_updated"`
}

// NodeReleasesClient fetches and caches Node.js release information
type NodeReleasesClient struct {
	cache  *CacheManager
	logger Logger
}

// NewNodeReleasesClient creates a new NodeReleasesClient instance
func NewNodeReleasesClient(cache *CacheManager, logger Logger) *NodeReleasesClient {
	return &NodeReleasesClient{
		cache:  cache,
		logger: logger,
	}
}

// GetVersionForCodename returns the major version for a given LTS codename
// Returns empty string if not found
func (c *NodeReleasesClient) GetVersionForCodename(codename string) (string, error) {
	codename = strings.ToLower(codename)

	// Try to load from cache first
	cached, err := c.loadFromCache()
	if err == nil && cached != nil {
		if version, found := cached.CodenameToVersion[codename]; found {
			return version, nil
		}
	}

	// Cache miss or invalid - fetch from API
	if err := c.refreshCache(); err != nil {
		// If refresh fails, return error
		return "", fmt.Errorf("failed to fetch Node.js releases: %w", err)
	}

	// Try again from fresh cache
	cached, err = c.loadFromCache()
	if err != nil {
		return "", err
	}

	if version, found := cached.CodenameToVersion[codename]; found {
		return version, nil
	}

	return "", fmt.Errorf("codename '%s' not found", codename)
}

// loadFromCache loads the cache if valid, returns nil if invalid or not found
func (c *NodeReleasesClient) loadFromCache() (*NodeReleasesCache, error) {
	// Check if cache is valid (exists and not too old)
	if !c.cache.IsCacheValid(cacheFileName, cacheMaxAge) {
		return nil, fmt.Errorf("cache invalid or expired")
	}

	var cached NodeReleasesCache
	if err := c.cache.ReadCache(cacheFileName, &cached); err != nil {
		return nil, err
	}

	return &cached, nil
}

// refreshCache fetches fresh data from Node.js API and updates cache
func (c *NodeReleasesClient) refreshCache() error {
	c.logger.Info("Fetching Node.js releases from API...")

	// Fetch from Node.js API
	releases, err := c.fetchReleases()
	if err != nil {
		return err
	}

	// Build codename to version map
	codenameMap := make(map[string]string)

	for _, release := range releases {
		// Check if this is an LTS release with a codename
		if ltsCodename, isLTS := release.LTS.(string); isLTS {
			// Extract major version from "v20.11.0" -> "20"
			version := strings.TrimPrefix(release.Version, "v")
			parts := strings.Split(version, ".")
			if len(parts) > 0 {
				majorVersion := parts[0]
				codenameKey := strings.ToLower(ltsCodename)

				// Keep the latest version for each codename
				// (API is sorted newest first, so first occurrence is latest)
				if _, exists := codenameMap[codenameKey]; !exists {
					codenameMap[codenameKey] = majorVersion
				}
			}
		}
	}

	// Save to cache
	cached := NodeReleasesCache{
		CodenameToVersion: codenameMap,
		LastUpdated:       time.Now(),
	}

	if err := c.cache.WriteCache(cacheFileName, cached); err != nil {
		return err
	}

	c.logger.Success(fmt.Sprintf("Node.js releases cache updated (%d codenames)", len(codenameMap)))
	return nil
}

// fetchReleases fetches the releases list from Node.js API
func (c *NodeReleasesClient) fetchReleases() ([]NodeRelease, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(nodeReleasesURL)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var releases []NodeRelease
	if err := json.Unmarshal(body, &releases); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return releases, nil
}
