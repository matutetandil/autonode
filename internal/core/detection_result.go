package core

// DetectionResult represents the result of version detection
// Single Responsibility Principle: Only responsible for holding detection result data
type DetectionResult struct {
	Found   bool
	Version string
	Source  string
}
