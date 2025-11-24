# Architecture

AutoNode is built with **Go** following SOLID principles and clean architecture.

## Why Go?

- **No circular dependency**: Doesn't require Node.js to manage Node.js
- **Single binary**: Easy distribution, no npm install needed
- **Fast**: Instant startup, compiled code
- **Cross-platform**: Easy to compile for all platforms
- **Self-contained**: No runtime dependencies (~6MB binary)

## Project Structure

```
autonode/
├── cmd/autonode/              # CLI entry point
│   ├── main.go                # Main + dependency injection
│   └── commands/              # Cobra commands
│       ├── run.go             # Main autonode command
│       ├── shell.go           # Shell integration
│       ├── update.go          # Self-update
│       └── config.go          # Local configuration
│
├── internal/
│   ├── core/                  # Core abstractions
│   │   ├── logger.go          # Logger interface
│   │   ├── shell_executor.go  # ShellExecutor interface
│   │   ├── version_detector.go # VersionDetector interface
│   │   ├── version_manager.go # VersionManager interface
│   │   ├── profile_detector.go # ProfileDetector interface
│   │   ├── profile_switcher.go # ProfileSwitcher interface
│   │   ├── service.go         # AutoNodeService orchestrator
│   │   ├── cache.go           # CacheManager
│   │   ├── update_checker.go  # Automatic update checks
│   │   └── ...                # Implementations
│   │
│   ├── detectors/             # Version detection
│   │   ├── autonode_yml_version.go  # .autonode.yml (priority 0)
│   │   ├── nvmrc.go                 # .nvmrc (priority 1)
│   │   ├── node_version.go          # .node-version (priority 2)
│   │   ├── package_json.go          # package.json (priority 3)
│   │   └── dockerfile.go            # Dockerfile (priority 4)
│   │
│   ├── managers/              # Version managers
│   │   ├── nvm.go             # nvm support
│   │   ├── nvs.go             # nvs support
│   │   └── volta.go           # Volta support
│   │
│   └── switchers/             # npm profile switchers
│       ├── npmrc_switcher.go
│       ├── ts_npmrc_switcher.go
│       └── rc_manager_switcher.go
│
├── go.mod                     # Go module
└── Makefile                   # Build automation
```

## SOLID Principles

### Single Responsibility (SRP)

Each file/type has one clear purpose:
- `ConsoleLogger` only handles logging
- `ExecShell` only handles shell execution
- Each detector handles one version source
- Each manager handles one version manager

### Open/Closed (OCP)

Adding new detectors or managers doesn't require modifying existing code:

```go
// Just create a new detector...
type MyDetector struct{}

func (d *MyDetector) Detect(path string) (DetectionResult, error) { ... }
func (d *MyDetector) GetPriority() int { return 5 }
func (d *MyDetector) GetSourceName() string { return "my-source" }

// ...and add it to the list
detectors := []VersionDetector{
    NewNvmrcDetector(),
    NewMyDetector(), // Added!
}
```

### Liskov Substitution (LSP)

All implementations are interchangeable through interfaces:
- Any `VersionDetector` can be used for detection
- Any `VersionManager` can be used for switching
- Any `ProfileSwitcher` can be used for npm profiles

### Interface Segregation (ISP)

Small, focused interfaces:

```go
type Logger interface {
    Info(message string)
    Success(message string)
    Warning(message string)
    Error(message string)
}

type VersionDetector interface {
    Detect(projectPath string) (DetectionResult, error)
    GetPriority() int
    GetSourceName() string
}
```

### Dependency Inversion (DIP)

High-level modules depend on abstractions:

```go
// Service depends on interfaces, not concrete types
type AutoNodeService struct {
    logger    Logger           // interface
    detectors []VersionDetector // interface
    managers  []VersionManager  // interface
}

// Dependencies injected via constructor
func NewAutoNodeService(
    logger Logger,
    detectors []VersionDetector,
    managers []VersionManager,
) *AutoNodeService
```

## Adding a New Version Detector

1. Create `internal/detectors/my_detector.go`:

```go
package detectors

import "github.com/matutetandil/autonode/internal/core"

type MyDetector struct{}

func NewMyDetector() *MyDetector {
    return &MyDetector{}
}

func (d *MyDetector) Detect(projectPath string) (core.DetectionResult, error) {
    // Read your version source
    // Return DetectionResult{Found: true, Version: "18.0.0", Source: "my-file"}
}

func (d *MyDetector) GetPriority() int {
    return 5 // Lower = higher priority
}

func (d *MyDetector) GetSourceName() string {
    return "my-file"
}
```

2. Add to `cmd/autonode/commands/run.go`:

```go
detectorsList := []core.VersionDetector{
    detectors.NewAutonodeYmlVersionDetector(),
    detectors.NewNvmrcDetector(),
    // ...
    detectors.NewMyDetector(), // Add here
}
```

3. Add tests in `internal/detectors/my_detector_test.go`

## Adding a New Version Manager

1. Create `internal/managers/my_manager.go`:

```go
package managers

import "github.com/matutetandil/autonode/internal/core"

type MyManager struct {
    shell core.ShellExecutor
}

func NewMyManager(shell core.ShellExecutor) *MyManager {
    return &MyManager{shell: shell}
}

func (m *MyManager) GetName() string {
    return "my-manager"
}

func (m *MyManager) IsInstalled() bool {
    return m.shell.CommandExists("my-manager")
}

func (m *MyManager) IsVersionInstalled(version string) (bool, error) {
    // Check if version is installed
}

func (m *MyManager) InstallVersion(version string) error {
    // Install the version
}

func (m *MyManager) UseVersion(version string) error {
    // Switch to the version
}
```

2. Add to `cmd/autonode/commands/run.go`:

```go
managersList := []core.VersionManager{
    managers.NewNvmManager(shell),
    managers.NewNvsManager(shell),
    managers.NewVoltaManager(shell),
    managers.NewMyManager(shell), // Add here
}
```

## Execution Flow

```
1. CLI Entry (main.go)
   ├── Parse flags
   ├── Start async update check (goroutine)
   └── Execute command

2. Run Command (run.go)
   ├── Create dependencies (DI)
   └── Call service.Run(config)

3. AutoNodeService.Run()
   ├── Detect version (iterate detectors by priority)
   ├── Find installed manager
   ├── Install version if needed
   ├── Switch to version
   ├── Detect npm profile
   └── Switch profile if configured

4. After Command
   └── Show update banner if available
```

## Testing

Tests use table-driven approach:

```go
func TestMyDetector(t *testing.T) {
    tests := []struct {
        name    string
        content string
        want    DetectionResult
        wantErr bool
    }{
        {"valid version", "18.0.0", DetectionResult{Found: true, Version: "18.0.0"}, false},
        {"empty file", "", DetectionResult{Found: false}, false},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Create temp file, run detector, assert results
        })
    }
}
```

Run tests:

```bash
go test ./...           # All tests
go test ./... -v        # Verbose
go test -cover ./...    # With coverage
```
