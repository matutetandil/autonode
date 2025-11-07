package core

import (
	"bytes"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

// ExecShell implements the ShellExecutor interface using os/exec
// Single Responsibility Principle: Only responsible for executing shell commands
// Dependency Inversion Principle: Depends on ShellExecutor interface
type ExecShell struct{}

// NewExecShell creates a new ExecShell instance
func NewExecShell() *ExecShell {
	return &ExecShell{}
}

// Execute runs a command with arguments and returns the output
func (e *ExecShell) Execute(command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("command failed: %w, stderr: %s", err, stderr.String())
	}

	return strings.TrimSpace(stdout.String()), nil
}

// ExecuteInShell runs a command in a shell environment
func (e *ExecShell) ExecuteInShell(command string) (string, error) {
	var cmd *exec.Cmd

	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", command)
	} else {
		cmd = exec.Command("sh", "-c", command)
	}

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("command failed: %w, stderr: %s", err, stderr.String())
	}

	return strings.TrimSpace(stdout.String()), nil
}

// CommandExists checks if a command is available in the system
func (e *ExecShell) CommandExists(command string) bool {
	var cmd *exec.Cmd

	if runtime.GOOS == "windows" {
		cmd = exec.Command("where", command)
	} else {
		cmd = exec.Command("which", command)
	}

	err := cmd.Run()
	return err == nil
}
