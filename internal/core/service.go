package core

import (
	"fmt"
	"sort"
)

// AutoNodeService orchestrates version detection and switching, as well as npm profile switching
// Single Responsibility Principle: Only responsible for orchestrating the workflow
// Dependency Inversion Principle: Depends on abstractions (interfaces), not concrete implementations
// Open/Closed Principle: New detectors and managers can be added without modifying this service
type AutoNodeService struct {
	logger           Logger
	detectors        []VersionDetector
	managers         []VersionManager
	profileDetectors []ProfileDetector
	profileSwitchers []ProfileSwitcher
}

// NewAutoNodeService creates a new AutoNodeService with injected dependencies
// Dependency Injection: All dependencies are injected through the constructor
func NewAutoNodeService(
	logger Logger,
	detectors []VersionDetector,
	managers []VersionManager,
	profileDetectors []ProfileDetector,
	profileSwitchers []ProfileSwitcher,
) *AutoNodeService {
	// Sort detectors by priority (lower number = higher priority)
	sortedDetectors := make([]VersionDetector, len(detectors))
	copy(sortedDetectors, detectors)
	sort.Slice(sortedDetectors, func(i, j int) bool {
		return sortedDetectors[i].GetPriority() < sortedDetectors[j].GetPriority()
	})

	// Sort profile detectors by priority (lower number = higher priority)
	sortedProfileDetectors := make([]ProfileDetector, len(profileDetectors))
	copy(sortedProfileDetectors, profileDetectors)
	sort.Slice(sortedProfileDetectors, func(i, j int) bool {
		return sortedProfileDetectors[i].GetPriority() < sortedProfileDetectors[j].GetPriority()
	})

	return &AutoNodeService{
		logger:           logger,
		detectors:        sortedDetectors,
		managers:         managers,
		profileDetectors: sortedProfileDetectors,
		profileSwitchers: profileSwitchers,
	}
}

// Run executes the main workflow: detect version, find manager, and switch version
func (s *AutoNodeService) Run(config Config) error {
	s.logger.Info(fmt.Sprintf("Scanning project at: %s", config.ProjectPath))

	// Step 1: Detect Node.js version
	result, err := s.detectVersion(config.ProjectPath)
	if err != nil {
		return err
	}

	if !result.Found {
		s.logger.Error("No Node.js version specification found in project")
		return fmt.Errorf("no version found")
	}

	s.logger.Success(fmt.Sprintf("Detected Node.js version %s from %s", result.Version, result.Source))

	// Detect npm profile configuration (for dry-run display in check mode)
	profileResult, _ := s.detectProfile(config.ProjectPath)
	if profileResult.Found {
		s.logger.Success(fmt.Sprintf("Detected npm profile '%s' from %s", profileResult.ProfileName, profileResult.Source))
	}

	// If check-only mode, stop here (dry-run completed)
	if config.CheckOnly {
		return nil
	}

	// Step 2: Find an installed version manager
	manager, err := s.findVersionManager()
	if err != nil {
		return err
	}

	s.logger.Info(fmt.Sprintf("Using version manager: %s", manager.GetName()))

	// Step 3: Check if version is already installed
	installed, err := manager.IsVersionInstalled(result.Version)
	if err != nil {
		s.logger.Warning(fmt.Sprintf("Could not check if version is installed: %v", err))
	}

	// Step 4: Install version if needed
	if !installed || config.Force {
		if config.Force {
			s.logger.Info(fmt.Sprintf("Force installing Node.js %s...", result.Version))
		} else {
			s.logger.Info(fmt.Sprintf("Installing Node.js %s...", result.Version))
		}

		err = manager.InstallVersion(result.Version)
		if err != nil {
			s.logger.Error(fmt.Sprintf("Failed to install version: %v", err))
			return err
		}

		s.logger.Success(fmt.Sprintf("Node.js %s installed successfully", result.Version))
	} else {
		s.logger.Info(fmt.Sprintf("Node.js %s is already installed", result.Version))
	}

	// Step 5: Switch to the version
	s.logger.Info(fmt.Sprintf("Switching to Node.js %s...", result.Version))
	err = manager.UseVersion(result.Version)
	if err != nil {
		s.logger.Error(fmt.Sprintf("Failed to switch version: %v", err))
		return err
	}

	s.logger.Success(fmt.Sprintf("Successfully switched to Node.js %s", result.Version))

	// Step 6: Switch npm profile if configured
	s.switchProfileIfConfigured(config.ProjectPath)

	return nil
}

// detectVersion tries all detectors in priority order
// Chain of Responsibility Pattern: Try detectors until one succeeds
func (s *AutoNodeService) detectVersion(projectPath string) (DetectionResult, error) {
	for _, detector := range s.detectors {
		result, err := detector.Detect(projectPath)
		if err != nil {
			s.logger.Warning(fmt.Sprintf("Detector %s failed: %v", detector.GetSourceName(), err))
			continue
		}

		if result.Found {
			return result, nil
		}
	}

	return DetectionResult{Found: false}, nil
}

// findVersionManager returns the first installed version manager
// Strategy Pattern: Select the first available strategy
func (s *AutoNodeService) findVersionManager() (VersionManager, error) {
	for _, manager := range s.managers {
		if manager.IsInstalled() {
			return manager, nil
		}
	}

	return nil, fmt.Errorf("no version manager found (nvm, nvs, or volta)")
}

// detectProfile tries all profile detectors in priority order
// Chain of Responsibility Pattern: Try detectors until one succeeds
func (s *AutoNodeService) detectProfile(projectPath string) (ProfileDetectionResult, error) {
	for _, detector := range s.profileDetectors {
		result, err := detector.Detect(projectPath)
		if err != nil {
			// Silent failure - just try next detector
			continue
		}

		if result.Found {
			return result, nil
		}
	}

	return ProfileDetectionResult{Found: false}, nil
}

// findProfileSwitcher returns the first installed profile switcher
// Strategy Pattern: Select the first available strategy
func (s *AutoNodeService) findProfileSwitcher() ProfileSwitcher {
	for _, switcher := range s.profileSwitchers {
		if switcher.IsInstalled() {
			return switcher
		}
	}

	return nil
}

// switchProfileIfConfigured attempts to switch npm profile if one is configured
// This is called after Node.js version switching and operates silently:
// - If no profile is configured: does nothing (silent)
// - If no profile switcher is installed: does nothing (silent)
// - If profile doesn't exist: logs warning
// - If switch succeeds: logs success
func (s *AutoNodeService) switchProfileIfConfigured(projectPath string) {
	// Try to detect profile configuration
	profileResult, err := s.detectProfile(projectPath)
	if err != nil || !profileResult.Found {
		// No profile configured - silent, this is normal
		return
	}

	// Try to find an installed profile switcher
	switcher := s.findProfileSwitcher()
	if switcher == nil {
		// No profile switcher installed - silent, user may not use profile tools
		return
	}

	// Check if the profile exists
	exists, err := switcher.ProfileExists(profileResult.ProfileName)
	if err != nil {
		s.logger.Warning(fmt.Sprintf("Could not verify if npm profile '%s' exists: %v", profileResult.ProfileName, err))
		return
	}

	if !exists {
		s.logger.Warning(fmt.Sprintf("npm profile '%s' (from %s) not found in %s",
			profileResult.ProfileName, profileResult.Source, switcher.GetName()))
		return
	}

	// Switch to the profile
	s.logger.Info(fmt.Sprintf("Switching to npm profile '%s' using %s...",
		profileResult.ProfileName, switcher.GetName()))

	err = switcher.SwitchProfile(profileResult.ProfileName)
	if err != nil {
		s.logger.Warning(fmt.Sprintf("Failed to switch npm profile: %v", err))
		return
	}

	s.logger.Success(fmt.Sprintf("Successfully switched to npm profile '%s'", profileResult.ProfileName))
}
