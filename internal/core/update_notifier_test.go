package core

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func TestUpdateNotifier_ShowUpdateBanner_NilResult(t *testing.T) {
	logger := NewConsoleLogger()
	notifier := NewUpdateNotifier(logger)

	// Capture stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	notifier.ShowUpdateBanner(nil)

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	// Should not print anything for nil result
	if output != "" {
		t.Errorf("Expected no output for nil result, got: %q", output)
	}
}

func TestUpdateNotifier_ShowUpdateBanner_NoUpdateAvailable(t *testing.T) {
	logger := NewConsoleLogger()
	notifier := NewUpdateNotifier(logger)

	result := &UpdateCheckResult{
		UpdateAvailable: false,
		CurrentVersion:  "0.7.0",
		LatestVersion:   "0.7.0",
	}

	// Capture stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	notifier.ShowUpdateBanner(result)

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	// Should not print anything when no update available
	if output != "" {
		t.Errorf("Expected no output when no update available, got: %q", output)
	}
}

func TestUpdateNotifier_ShowUpdateBanner_UpdateAvailable(t *testing.T) {
	logger := NewConsoleLogger()
	notifier := NewUpdateNotifier(logger)

	result := &UpdateCheckResult{
		UpdateAvailable: true,
		CurrentVersion:  "0.5.0",
		LatestVersion:   "v0.7.0",
	}

	// Capture stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	notifier.ShowUpdateBanner(result)

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	// Should contain key information
	if !strings.Contains(output, "new version") {
		t.Error("Banner should contain 'new version'")
	}
	if !strings.Contains(output, "0.5.0") {
		t.Error("Banner should contain current version '0.5.0'")
	}
	if !strings.Contains(output, "v0.7.0") {
		t.Error("Banner should contain latest version 'v0.7.0'")
	}
	if !strings.Contains(output, "autonode update") {
		t.Error("Banner should contain 'autonode update' command")
	}

	// Should have box characters
	if !strings.Contains(output, "╭") || !strings.Contains(output, "╯") {
		t.Error("Banner should have box drawing characters")
	}
}

func TestUpdateNotifier_BuildBanner(t *testing.T) {
	logger := NewConsoleLogger()
	notifier := NewUpdateNotifier(logger)

	banner := notifier.buildBanner("0.5.0", "v0.8.0")

	// Check structure
	lines := strings.Split(banner, "\n")
	if len(lines) != 5 {
		t.Errorf("Expected 5 lines in banner, got %d", len(lines))
	}

	// First line should be top border
	if !strings.HasPrefix(lines[0], "╭") || !strings.HasSuffix(lines[0], "╮") {
		t.Errorf("First line should be top border, got: %q", lines[0])
	}

	// Last line should be bottom border
	if !strings.HasPrefix(lines[4], "╰") || !strings.HasSuffix(lines[4], "╯") {
		t.Errorf("Last line should be bottom border, got: %q", lines[4])
	}

	// Middle lines should be content with vertical bars
	for i := 1; i <= 3; i++ {
		if !strings.HasPrefix(lines[i], "│") || !strings.HasSuffix(lines[i], "│") {
			t.Errorf("Line %d should have vertical bars, got: %q", i, lines[i])
		}
	}
}
