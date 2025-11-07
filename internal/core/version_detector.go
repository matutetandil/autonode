package core

// VersionDetector interface defines version detection operations
// Interface Segregation Principle: Focused interface for version detection
// Open/Closed Principle: New detectors can be added without modifying existing code
// Liskov Substitution Principle: All implementations are interchangeable
type VersionDetector interface {
	Detect(projectPath string) (DetectionResult, error)
	GetPriority() int
	GetSourceName() string
}
