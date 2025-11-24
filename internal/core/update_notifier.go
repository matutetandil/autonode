package core

import (
	"fmt"
	"strings"
)

// UpdateNotifier displays update notifications to the user
// Single Responsibility Principle: Only responsible for displaying update notifications
type UpdateNotifier struct {
	logger Logger
}

// NewUpdateNotifier creates a new UpdateNotifier instance
// Dependency Inversion Principle: Depends on Logger interface
func NewUpdateNotifier(logger Logger) *UpdateNotifier {
	return &UpdateNotifier{
		logger: logger,
	}
}

// ShowUpdateBanner displays an update notification banner if an update is available
func (n *UpdateNotifier) ShowUpdateBanner(result *UpdateCheckResult) {
	if result == nil || !result.UpdateAvailable {
		return
	}

	// Format versions for display
	current := result.CurrentVersion
	latest := result.LatestVersion

	// Build the banner
	banner := n.buildBanner(current, latest)

	// Print the banner
	fmt.Println()
	fmt.Println(banner)
}

// buildBanner creates a nice-looking update notification banner
func (n *UpdateNotifier) buildBanner(current, latest string) string {
	// Box characters
	topLeft := "╭"
	topRight := "╮"
	bottomLeft := "╰"
	bottomRight := "╯"
	horizontal := "─"
	vertical := "│"

	// Content lines
	line1 := "A new version of autonode is available!"
	line2 := fmt.Sprintf("Current: %s → Latest: %s", current, latest)
	line3 := "Run 'autonode update' to upgrade"

	// Calculate the width (max line length + padding)
	maxLen := max(len(line1), max(len(line2), len(line3)))
	width := maxLen + 4 // 2 spaces padding on each side

	// Build the banner
	var sb strings.Builder

	// Top border
	sb.WriteString(topLeft)
	sb.WriteString(strings.Repeat(horizontal, width))
	sb.WriteString(topRight)
	sb.WriteString("\n")

	// Content lines with padding
	for _, line := range []string{line1, line2, line3} {
		padding := width - len(line) - 2 // -2 for the spaces
		leftPad := 1
		rightPad := padding + 1

		sb.WriteString(vertical)
		sb.WriteString(strings.Repeat(" ", leftPad))
		sb.WriteString(line)
		sb.WriteString(strings.Repeat(" ", rightPad))
		sb.WriteString(vertical)
		sb.WriteString("\n")
	}

	// Bottom border
	sb.WriteString(bottomLeft)
	sb.WriteString(strings.Repeat(horizontal, width))
	sb.WriteString(bottomRight)

	return sb.String()
}
