package switchers

import (
	"fmt"
	"strings"
	"testing"
)

func TestTsNpmrcSwitcher_GetName(t *testing.T) {
	shell := &MockShell{}
	switcher := NewTsNpmrcSwitcher(shell)

	name := switcher.GetName()
	if name != "ts-npmrc" {
		t.Errorf("GetName() = %q, want %q", name, "ts-npmrc")
	}
}

func TestTsNpmrcSwitcher_IsInstalled(t *testing.T) {
	t.Run("ts-npmrc is installed", func(t *testing.T) {
		shell := &MockShell{
			CommandExistsFunc: func(command string) bool {
				return command == "ts-npmrc"
			},
		}

		switcher := NewTsNpmrcSwitcher(shell)
		got := switcher.IsInstalled()

		if !got {
			t.Errorf("IsInstalled() = false, want true when ts-npmrc is in PATH")
		}
	})

	// Note: We don't test the "not installed" case because findExecutable()
	// searches the real filesystem (~/.nvm/versions/node/*/bin/ts-npmrc).
	// If ts-npmrc is installed on the test system, the test would incorrectly pass/fail.
	// In production, the behavior is correct: it finds ts-npmrc even if not in current PATH.
}

func TestTsNpmrcSwitcher_ProfileExists(t *testing.T) {
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
			name:        "profile exists in formatted output",
			profileName: "work",
			listOutput:  "Available profiles:\n  - default\n  - work\n  - personal",
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
					return command == "ts-npmrc"
				},
				ExecuteFunc: func(command string, args ...string) (string, error) {
					if command == "ts-npmrc" && len(args) == 1 && args[0] == "list" {
						return tt.listOutput, tt.listError
					}
					return "", nil
				},
			}

			switcher := NewTsNpmrcSwitcher(shell)
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

func TestTsNpmrcSwitcher_SwitchProfile(t *testing.T) {
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
					return command == "ts-npmrc"
				},
				ExecuteFunc: func(command string, args ...string) (string, error) {
					if command == "ts-npmrc" && len(args) == 3 &&
					   args[0] == "link" && args[1] == "-p" && args[2] == tt.profileName {
						return "", tt.executeError
					}
					return "", fmt.Errorf("unexpected command: %s %v", command, args)
				},
			}

			switcher := NewTsNpmrcSwitcher(shell)
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
