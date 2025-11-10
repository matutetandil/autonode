package core

// ProfileDetectionResult represents the result of detecting an npm profile configuration
// in a project directory.
//
// This type adheres to the Single Responsibility Principle (SRP) by only
// representing detection results for npm profile configuration.
type ProfileDetectionResult struct {
	// Found indicates whether a profile configuration was detected
	Found bool

	// ProfileName is the name of the npm profile to use (e.g., "work", "personal")
	ProfileName string

	// Source indicates where the profile was detected from (e.g., ".autonode.yml", "package.json")
	Source string
}
