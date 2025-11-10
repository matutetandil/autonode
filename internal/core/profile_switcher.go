package core

// ProfileSwitcher defines the interface for switching npm profiles using
// various npm profile management tools (npmrc, ts-npmrc, rc-manager, etc.).
//
// This interface adheres to:
// - Interface Segregation Principle (ISP): Small, focused interface
// - Open/Closed Principle (OCP): New profile switchers can be added without modifying existing code
// - Liskov Substitution Principle (LSP): All implementations are interchangeable
// - Dependency Inversion Principle (DIP): High-level code depends on this abstraction
type ProfileSwitcher interface {
	// GetName returns the name of the profile management tool
	// (e.g., "npmrc", "ts-npmrc", "rc-manager")
	GetName() string

	// IsInstalled checks if the profile management tool is installed and available
	IsInstalled() bool

	// ProfileExists checks if a specific profile exists in the tool's configuration
	ProfileExists(profileName string) (bool, error)

	// SwitchProfile switches to the specified npm profile
	SwitchProfile(profileName string) error
}
