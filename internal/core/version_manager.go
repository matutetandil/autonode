package core

// VersionManager interface defines version manager operations
// Interface Segregation Principle: Focused interface for version management
// Open/Closed Principle: New managers can be added without modifying existing code
// Liskov Substitution Principle: All implementations (nvm, nvs, volta) are interchangeable
type VersionManager interface {
	GetName() string
	IsInstalled() bool
	IsVersionInstalled(version string) (bool, error)
	InstallVersion(version string) error
	UseVersion(version string) error
}
