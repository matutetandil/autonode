package switchers

import "github.com/matutetandil/autonode/internal/core"

// MockShell is a mock implementation of core.ShellExecutor for testing
type MockShell struct {
	CommandExistsFunc func(command string) bool
	ExecuteFunc       func(command string, args ...string) (string, error)
}

// Ensure MockShell implements core.ShellExecutor
var _ core.ShellExecutor = (*MockShell)(nil)

// CommandExists calls the mock function
func (m *MockShell) CommandExists(command string) bool {
	if m.CommandExistsFunc != nil {
		return m.CommandExistsFunc(command)
	}
	return false
}

// Execute calls the mock function
func (m *MockShell) Execute(command string, args ...string) (string, error) {
	if m.ExecuteFunc != nil {
		return m.ExecuteFunc(command, args...)
	}
	return "", nil
}

// ExecuteInShell calls Execute (same behavior for mock)
func (m *MockShell) ExecuteInShell(command string) (string, error) {
	return m.Execute(command)
}
