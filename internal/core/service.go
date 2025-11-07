package core

import (
	"fmt"
	"sort"
)

// AutoNodeService orchestrates version detection and switching
// Single Responsibility Principle: Only responsible for orchestrating the workflow
// Dependency Inversion Principle: Depends on abstractions (interfaces), not concrete implementations
// Open/Closed Principle: New detectors and managers can be added without modifying this service
type AutoNodeService struct {
	logger    Logger
	detectors []VersionDetector
	managers  []VersionManager
}

// NewAutoNodeService creates a new AutoNodeService with injected dependencies
// Dependency Injection: All dependencies are injected through the constructor
func NewAutoNodeService(
	logger Logger,
	detectors []VersionDetector,
	managers []VersionManager,
) *AutoNodeService {
	// Sort detectors by priority (lower number = higher priority)
	sortedDetectors := make([]VersionDetector, len(detectors))
	copy(sortedDetectors, detectors)
	sort.Slice(sortedDetectors, func(i, j int) bool {
		return sortedDetectors[i].GetPriority() < sortedDetectors[j].GetPriority()
	})

	return &AutoNodeService{
		logger:    logger,
		detectors: sortedDetectors,
		managers:  managers,
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

	// If check-only mode, stop here
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
