package core

// ProfileDetector defines the interface for detecting npm profile configuration
// in a project directory.
//
// This interface adheres to:
// - Interface Segregation Principle (ISP): Small, focused interface
// - Open/Closed Principle (OCP): New detectors can be added without modifying existing code
// - Liskov Substitution Principle (LSP): All implementations are interchangeable
type ProfileDetector interface {
	// Detect searches for npm profile configuration in the given project path
	// and returns a ProfileDetectionResult indicating whether a profile was found
	// and which profile should be used.
	Detect(projectPath string) (ProfileDetectionResult, error)

	// GetPriority returns the priority of this detector.
	// Lower numbers indicate higher priority.
	// Example: .autonode.yml (priority 1) takes precedence over package.json (priority 2)
	GetPriority() int

	// GetSourceName returns a human-readable name of the source this detector checks
	// (e.g., ".autonode.yml", "package.json")
	GetSourceName() string
}
