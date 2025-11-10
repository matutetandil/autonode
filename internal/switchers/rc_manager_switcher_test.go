package switchers

import (
	"fmt"
	"strings"
	"testing"
)

func TestRcManagerSwitcher_GetName(t *testing.T) {
	shell := &MockShell{}
	switcher := NewRcManagerSwitcher(shell)

	name := switcher.GetName()
	if name != "rc-manager" {
		t.Errorf("GetName() = %q, want %q", name, "rc-manager")
	}
}

func TestRcManagerSwitcher_IsInstalled(t *testing.T) {
	t.Run("rc-manager is installed", func(t *testing.T) {
		shell := &MockShell{
			CommandExistsFunc: func(command string) bool {
				return command == "rc-manager"
			},
		}

		switcher := NewRcManagerSwitcher(shell)
		got := switcher.IsInstalled()

		if !got {
			t.Errorf("IsInstalled() = false, want true when rc-manager is in PATH")
		}
	})

	// Note: We don't test the "not installed" case because findExecutable()
	// searches the real filesystem (~/.nvm/versions/node/*/bin/rc-manager).
	// If rc-manager is installed on the test system, the test would incorrectly pass/fail.
	// In production, the behavior is correct: it finds rc-manager even if not in current PATH.
}

func TestRcManagerSwitcher_ProfileExists(t *testing.T) {
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
			name:        "profile exists with formatting",
			profileName: "work",
			listOutput:  "Available profiles:\ndefault\nwork\npersonal",
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
					return command == "rc-manager"
				},
				ExecuteFunc: func(command string, args ...string) (string, error) {
					if command == "rc-manager" && len(args) == 1 && args[0] == "list" {
						return tt.listOutput, tt.listError
					}
					return "", nil
				},
			}

			switcher := NewRcManagerSwitcher(shell)
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

func TestRcManagerSwitcher_SwitchProfile(t *testing.T) {
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
					return command == "rc-manager"
				},
				ExecuteFunc: func(command string, args ...string) (string, error) {
					if command == "rc-manager" && len(args) == 2 &&
					   args[0] == "load" && args[1] == tt.profileName {
						return "", tt.executeError
					}
					return "", fmt.Errorf("unexpected command: %s %v", command, args)
				},
			}

			switcher := NewRcManagerSwitcher(shell)
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
