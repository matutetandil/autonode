package switchers

import (
	"fmt"
	"strings"
	"testing"
)

func TestNpmrcSwitcher_GetName(t *testing.T) {
	shell := &MockShell{}
	switcher := NewNpmrcSwitcher(shell)

	name := switcher.GetName()
	if name != "npmrc" {
		t.Errorf("GetName() = %q, want %q", name, "npmrc")
	}
}

func TestNpmrcSwitcher_IsInstalled(t *testing.T) {
	t.Run("npmrc is installed", func(t *testing.T) {
		shell := &MockShell{
			CommandExistsFunc: func(command string) bool {
				return command == "npmrc"
			},
		}

		switcher := NewNpmrcSwitcher(shell)
		got := switcher.IsInstalled()

		if !got {
			t.Errorf("IsInstalled() = false, want true when npmrc is in PATH")
		}
	})

	// Note: We don't test the "not installed" case because findExecutable()
	// searches the real filesystem (~/.nvm/versions/node/*/bin/npmrc).
	// If npmrc is installed on the test system, the test would incorrectly pass/fail.
	// In production, the behavior is correct: it finds npmrc even if not in current PATH.
}

func TestNpmrcSwitcher_ProfileExists(t *testing.T) {
	tests := []struct {
		name        string
		profileName string
		listOutput  string
		listError   error
		want        bool
		wantError   bool
	}{
		{
			name:        "profile exists",
			profileName: "work",
			listOutput:  "default\nwork\npersonal",
			want:        true,
			wantError:   false,
		},
		{
			name:        "profile exists with active marker",
			profileName: "work",
			listOutput:  "default\n* work\npersonal",
			want:        true,
			wantError:   false,
		},
		{
			name:        "profile does not exist",
			profileName: "nonexistent",
			listOutput:  "default\nwork\npersonal",
			want:        false,
			wantError:   false,
		},
		{
			name:        "list command fails",
			profileName: "work",
			listOutput:  "",
			listError:   fmt.Errorf("command failed"),
			want:        false,
			wantError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shell := &MockShell{
				CommandExistsFunc: func(command string) bool {
					return command == "npmrc"
				},
				ExecuteFunc: func(command string, args ...string) (string, error) {
					if command == "npmrc" && len(args) == 0 {
						return tt.listOutput, tt.listError
					}
					return "", nil
				},
			}

			switcher := NewNpmrcSwitcher(shell)
			got, err := switcher.ProfileExists(tt.profileName)

			if (err != nil) != tt.wantError {
				t.Errorf("ProfileExists() error = %v, wantError %v", err, tt.wantError)
				return
			}

			if got != tt.want {
				t.Errorf("ProfileExists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNpmrcSwitcher_SwitchProfile(t *testing.T) {
	tests := []struct {
		name         string
		profileName  string
		executeError error
		wantError    bool
	}{
		{
			name:         "switch succeeds",
			profileName:  "work",
			executeError: nil,
			wantError:    false,
		},
		{
			name:         "switch fails",
			profileName:  "work",
			executeError: fmt.Errorf("switch failed"),
			wantError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shell := &MockShell{
				CommandExistsFunc: func(command string) bool {
					return command == "npmrc"
				},
				ExecuteFunc: func(command string, args ...string) (string, error) {
					if command == "npmrc" && len(args) == 1 && args[0] == tt.profileName {
						return "", tt.executeError
					}
					return "", fmt.Errorf("unexpected command: %s %v", command, args)
				},
			}

			switcher := NewNpmrcSwitcher(shell)
			err := switcher.SwitchProfile(tt.profileName)

			if (err != nil) != tt.wantError {
				t.Errorf("SwitchProfile() error = %v, wantError %v", err, tt.wantError)
			}

			if err != nil && !strings.Contains(err.Error(), tt.profileName) {
				t.Errorf("error message should contain profile name %q, got: %v", tt.profileName, err)
			}
		})
	}
}
