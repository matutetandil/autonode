package detectors

import "fmt"

// mockReleasesClient is a mock implementation for testing
// It implements the same methods as NodeReleasesClient
type mockReleasesClient struct {
	codenameMap map[string]string
}

// newMockReleasesClient creates a new mock with predefined codenames
func newMockReleasesClient() *mockReleasesClient {
	return &mockReleasesClient{
		codenameMap: map[string]string{
			"krypton":  "24",
			"jod":      "22",
			"iron":     "20",
			"hydrogen": "18",
			"gallium":  "16",
			"fermium":  "14",
			"erbium":   "12",
		},
	}
}

// GetVersionForCodename returns the version for a given codename (mock implementation)
func (m *mockReleasesClient) GetVersionForCodename(codename string) (string, error) {
	if version, found := m.codenameMap[codename]; found {
		return version, nil
	}
	return "", fmt.Errorf("codename '%s' not found", codename)
}
