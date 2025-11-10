# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.5.1] - 2025-11-10

### Fixed

- **Shell integration npm profile switching**: Fixed critical bug where `autonode shell` (used by cd hooks) was not switching npm profiles, only Node.js versions. Shell integration now correctly switches both Node.js version and npm profile when configured.

### Changed

- **Code refactoring**: Eliminated code duplication between `run.go` and `shell.go` by implementing `ShellMode` flag in `Config`. Both normal mode and shell mode now use the same `AutoNodeService.Run()` method, reducing code by 87 lines and improving maintainability. This change is internal and does not affect the external API or behavior.

### Technical Details

- Added `ShellMode bool` field to `Config` struct
- Implemented `runShellMode()` method in `AutoNodeService` for shell command output
- Refactored `shell.go` to use `AutoNodeService` with dependency injection
- All 116 tests passing ✅
- Zero breaking changes

## [0.5.0] - 2025-11-10

### Added - npm Profile Management

AutoNode can now automatically switch npm profiles (registry configurations) when changing projects. This is perfect for managing different npm registries for work, personal, or client-specific projects.

#### Profile Switching Features

- **Multi-Tool Support**: Works with popular npm profile management tools:
  - **[npmrc](https://github.com/deoxxa/npmrc)** - Most popular npm profile switcher
  - **[ts-npmrc](https://github.com/darsi-an/ts-npmrc)** - Modern TypeScript implementation
  - **[rc-manager](https://github.com/Lalaluka/rc-manager)** - Supports npm and yarn

- **Flexible Configuration**: Configure profiles in two ways (with priority):
  1. `.autonode.yml` file with `npmProfile: work` (highest priority)
  2. `package.json` with `"autonode": { "npmProfile": "work" }`

- **Silent Operation**: Completely opt-in and non-intrusive
  - Only switches profiles if configured
  - No warnings if profile tool not installed
  - Works seamlessly alongside Node.js version switching

- **Profile Detection Priority**:
  1. `.autonode.yml` (checked first)
  2. `package.json` autonode field

- **Tool Priority**: If multiple tools installed, uses first available:
  1. npmrc (most popular)
  2. ts-npmrc
  3. rc-manager

#### Architecture

- **New Core Interfaces** (SOLID principles):
  - `ProfileDetector` interface for detecting profile configuration
  - `ProfileSwitcher` interface for switching npm profiles
  - `ProfileDetectionResult` type for detection results

- **New Detectors** (`internal/detectors/`):
  - `AutonodeYmlProfileDetector` - Reads `.autonode.yml` (priority 1)
  - `PackageJsonProfileDetector` - Reads `package.json` autonode field (priority 2)

- **New Switchers** (`internal/switchers/`):
  - `NpmrcSwitcher` - Manages profiles using npmrc
  - `TsNpmrcSwitcher` - Manages profiles using ts-npmrc
  - `RcManagerSwitcher` - Manages profiles using rc-manager

- **Service Integration**:
  - Updated `AutoNodeService` to handle profile switching
  - Integrated profile detection and switching into main workflow
  - Silent mode: no errors if profile tools missing

#### Testing

- **Comprehensive Test Coverage**: Added 38 new test cases
  - `AutonodeYmlProfileDetector`: 11 tests (YAML parsing, priorities, edge cases)
  - `PackageJsonProfileDetector`: 11 tests (JSON parsing, priorities, edge cases)
  - `NpmrcSwitcher`: 8 tests (installation, profile detection, switching)
  - `TsNpmrcSwitcher`: 8 tests (installation, profile detection, switching)
  - `RcManagerSwitcher`: 8 tests (installation, profile detection, switching)
  - Created `MockShell` for testing switchers without actual shell commands
  - All 119 test cases passing ✅

### Dependencies

- **Added**: `gopkg.in/yaml.v3` for YAML parsing in `.autonode.yml` files

### Documentation

- **README.md**: Added comprehensive "npm Profile Management" section
  - Installation instructions for profile tools
  - Configuration examples (both `.autonode.yml` and `package.json`)
  - How it works explanation
  - Example output
  - Priority documentation

- **Updated Features**: Added npm profile management to feature list
  - Clear explanation of silent mode behavior
  - Links to supported profile management tools

### Example Usage

After configuration, AutoNode automatically switches both Node.js version and npm profile:

```bash
$ cd my-work-project
✓ Detected Node.js version 18.17.0 from .nvmrc
Using version manager: nvm
Switching to Node.js 18.17.0...
✓ Successfully switched to Node.js 18.17.0
Switching to npm profile 'work' using npmrc...
✓ Successfully switched to npm profile 'work'
```

### Why This Matters

Many developers work on multiple projects with different npm registries (company registry, personal registry, client-specific registry). Previously, this required manual profile switching with tools like npmrc. Now, AutoNode handles it automatically based on project configuration, providing a complete environment setup in one step.

## [0.4.1] - 2025-11-07

### Fixed

- **Update Command Archive Handling**: Fixed critical bug in `autonode update`
  - GoReleaser creates `.tar.gz` archives (or `.zip` for Windows), not raw binaries
  - Updated to download and extract correct archive format
  - Added extraction support for both tar.gz (Unix) and zip (Windows)
  - Fixes 'HTTP 404' error when running `autonode update`

### Added

- **Version Check Before Update**: Added intelligent version comparison
  - Shows "You're already on the latest version" if up-to-date
  - Displays "Updating from X → Y" with current and target versions
  - Skips unnecessary downloads when already on latest version
  - Fetches latest version from GitHub API before downloading

- **Improved Update Messages**: Better user feedback during update process
  - Shows current version → target version before downloading
  - Displays archive name being downloaded
  - Clear success message with new version number

### Changed

- **Version Access**: Added `GetVersion()` getter function in main package
  - Maintains proper encapsulation (version variable stays private)
  - Allows update command to access current version without circular dependency
  - Uses Cobra's root command version instead of direct import

## [0.4.0] - 2025-11-07

### Added - Dynamic Node.js Release Cache

This release adds intelligent detection of Node.js LTS codenames in Dockerfiles through a dynamic caching system.

#### Dynamic Codename Resolution

- **Node.js API Integration**: Fetches release information from `https://nodejs.org/dist/index.json`
  - Maps LTS codenames to their major versions dynamically
  - No more hardcoded version mappings
  - Always up-to-date with latest Node.js releases

- **Smart Caching System**: Created `~/.autonode/` directory for persistent cache
  - Cache stored in `~/.autonode/node-releases.json`
  - 24-hour cache validity (refreshes daily)
  - Includes all historical and current LTS codenames (argon, boron, carbon, dubnium, erbium, fermium, gallium, hydrogen, iron, jod, krypton)
  - Graceful offline fallback with cached data

#### Enhanced Dockerfile Detection

- **LTS Codename Support**: Now detects Docker images with Node.js LTS codenames
  - `FROM node:iron-alpine` → Node 20
  - `FROM node:jod-alpine` → Node 22
  - `FROM node:krypton` → Node 24
  - Works with all image variants: `-alpine`, `-slim`, `-bullseye`, etc.

- **Special Tag Handling**:
  - `FROM node:lts` → Latest LTS version (currently 22)
  - `FROM node:latest` → Latest stable (currently 23)
  - `FROM node:current` → Latest stable (currently 23)

- **Backward Compatibility**: All numeric versions still work
  - `FROM node:18.17.0` → 18.17.0
  - `FROM node:20` → 20

#### Architecture Improvements

- **New Components**:
  - `internal/core/cache.go`: Generic cache manager for future extensibility
  - `internal/core/node_releases.go`: Node.js API client with intelligent caching
  - `internal/core/null_logger.go`: Silent logger for background operations

- **Dependency Injection**: `DockerfileDetector` now receives `releasesClient` interface for dynamic resolution and testability

#### Testing

- **Comprehensive Test Coverage**: Added 33 new test cases
  - `DockerfileDetector`: 41 tests total (added 27 new tests)
    - LTS codename detection (iron, jod, hydrogen, gallium, krypton)
    - LTS codenames with variants (alpine, slim, bullseye)
    - Special tags (lts, latest, current)
    - Special tags with variants
    - Complex multi-variant scenarios
  - `CacheManager`: 7 new tests
    - Cache creation, read, write operations
    - Cache validity checking
    - Cache clearing functionality
  - Created mock `releasesClient` for testing Dockerfile detector without HTTP requests
  - All 81 test cases passing ✅

### Changed

- Dockerfile detector now makes HTTP requests on first run or when cache expires
- Cache directory `~/.autonode/` created automatically on first use

### Why This Matters

Previously, Dockerfile detection only worked with numeric versions (`FROM node:18`). Many projects use LTS codenames (`FROM node:iron-alpine`) which were not recognized. This release makes AutoNode aware of all Node.js releases without requiring recompilation when new versions are released.

## [0.3.1] - 2025-11-07

### Fixed

- **Installation Script Archive Handling**: Fixed critical bug in `install.sh`
  - GoReleaser creates `.tar.gz` archives (or `.zip` for Windows), not raw binaries
  - Updated script to download correct archive format
  - Added extraction step using `tar` or `unzip`
  - Added cleanup of temporary files
  - Changed default version from hardcoded to `latest` (can override with `AUTONODE_VERSION` env var)
  - Fixes 404 error when running installation command

### Added

- **Unit Tests**: Complete test coverage for all detectors
  - `NvmrcDetector`: 9 tests covering versions, whitespace, lts aliases, empty files
  - `NodeVersionDetector`: 5 tests covering version formats and edge cases
  - `PackageJsonDetector`: 20 tests covering operators, ranges, OR conditions, invalid JSON
  - `DockerfileDetector`: 14 tests covering FROM patterns, alpine tags, comments, case sensitivity
  - Total: 48 tests, all passing
  - Tests for priority verification and source name verification

- **CI/CD Test Execution**: Added test execution to release workflow
  - All tests run automatically before building releases
  - Release process halts if tests fail
  - Ensures code quality before distribution

### Changed

- **README Roadmap**: Updated to reflect completed features
  - Marked as completed: Unit tests, CI/CD pipeline, Shell integration
  - Removed completed items from pending list

## [0.3.0] - 2025-11-07

### Added - Shell Integration & Auto-Update

This release adds automatic version switching and self-update capabilities, making AutoNode truly seamless.

#### Shell Integration

- **`autonode shell` command**: Outputs shell commands for `eval` integration
  - Returns version manager-specific commands (sources `nvm.sh` for nvm, etc.)
  - Silent operation - only outputs commands, no logging
  - Detects version and manager automatically
  - Works with bash, zsh, and fish shells

- **Automatic Installation Script** (`install.sh`):
  - One-line installation: `curl -fsSL https://raw.githubusercontent.com/matutetandil/autonode/main/install.sh | bash`
  - Auto-detects platform (macOS Intel/ARM, Linux amd64/arm64, Windows)
  - Downloads correct binary from GitHub releases
  - Installs to `/usr/local/bin/` or `~/.local/bin/`
  - **Automatically configures shell hook** for bash/zsh/fish
  - Adds `cd` hook that triggers version switching on directory change
  - Shell hook runs on startup and every `cd` command

#### Auto-Update

- **`autonode update` command**: Self-update to latest version
  - Detects current platform automatically
  - Downloads latest release from GitHub
  - Replaces binary in-place
  - Cross-platform support (macOS, Linux, Windows)
  - Handles symlinks correctly
  - Preserves file permissions

#### Uninstallation

- **Automatic Uninstall Script** (`uninstall.sh`):
  - One-line uninstall: `curl -fsSL https://raw.githubusercontent.com/matutetandil/autonode/main/uninstall.sh | bash`
  - Removes binary from all standard locations
  - Removes shell integration from config files
  - Creates backups before modifying configs
  - Supports bash, zsh, and fish

### Fixed

- **nvm Detection**: Fixed nvm detection by checking for `~/.nvm/nvm.sh` instead of trying to execute shell function
  - Respects `NVM_DIR` environment variable
  - Sources `nvm.sh` before executing nvm commands
  - Works in non-interactive shells

### Changed

- **Documentation**: Complete rewrite of installation instructions
  - Automatic installation is now the recommended method
  - Added "Updating" and "Uninstalling" sections
  - Added examples for shell integration
  - Updated usage instructions for automatic mode

### User Experience

**Before this release:**
- User had to manually run `autonode` in each directory
- Manual `nvm use` required after running autonode
- No easy way to update or uninstall

**After this release:**
- **Fully automatic**: Just `cd` into a directory and Node.js version switches automatically
- No manual commands needed after installation
- One-command update: `autonode update`
- One-command uninstall with cleanup

## [0.2.0] - 2025-11-07

### Changed - Complete Rewrite in Go

**BREAKING CHANGE**: Complete rewrite from TypeScript/Node.js to Go. This eliminates the circular dependency issue where a Node.js tool was managing Node.js versions.

#### Why the Rewrite?

- **No Circular Dependency**: AutoNode no longer requires Node.js to run, solving the fundamental chicken-and-egg problem
- **Native Binaries**: Single executable file per platform, no installation of dependencies required
- **Zero Runtime Dependencies**: No need for Node.js, npm, or any other runtime
- **Instant Startup**: Compiled binary with no interpreter overhead
- **Easy Distribution**: Download a binary and run it, no `npm install` needed
- **Cross-Platform Compilation**: Trivial to build for Linux, macOS (Intel & Apple Silicon), and Windows

#### Core Architecture (Go)

- Rewrote all interfaces and implementations in Go following SOLID principles
- Core interfaces: `Logger`, `ShellExecutor`, `VersionDetector`, `VersionManager`
- Implemented `AutoNodeService` as main orchestrator
- Implemented `ConsoleLogger` using github.com/fatih/color
- Implemented `ExecShell` using os/exec for cross-platform command execution
- Each type in a separate file for strict Single Responsibility Principle compliance

#### Version Detection (Go)

- Rewrote `NvmrcDetector` for reading `.nvmrc` files (priority 1)
- Rewrote `NodeVersionDetector` for reading `.node-version` files (priority 2)
- Rewrote `PackageJsonDetector` for reading `engines.node` from package.json (priority 3)
- Rewrote `DockerfileDetector` for parsing `FROM node:X` in Dockerfile (priority 4)
- Priority-based detection chain maintained

#### Version Managers (Go)

- Rewrote `NvmManager` for nvm support
- Rewrote `NvsManager` for nvs support
- Rewrote `VoltaManager` for Volta support
- Cross-platform command execution with proper shell handling

#### CLI (Go)

- Implemented CLI using github.com/spf13/cobra (industry standard)
- Commands:
  - `autonode` - Detect and switch to required Node.js version
  - `autonode --check` / `-c` - Display detected version without switching
  - `autonode --force` / `-f` - Reinstall version even if present
  - `autonode --version` / `-v` - Display AutoNode version
- Colored output with emojis for better UX

#### Build System

- Created comprehensive Makefile with targets:
  - `make build` - Build for current platform
  - `make build-all` - Build for all platforms
  - `make linux` - Build for Linux (amd64, arm64)
  - `make darwin` - Build for macOS (Intel, Apple Silicon)
  - `make windows` - Build for Windows (amd64)
  - `make install` - Install to /usr/local/bin
  - `make clean`, `make test`, `make fmt`, `make tidy`
- Native cross-compilation support via GOOS/GOARCH

#### Documentation

- Updated README.md for Go implementation with:
  - Installation instructions for pre-built binaries
  - Build from source instructions
  - Updated architecture section
  - Go-specific extension guide
  - Development commands for Go
- Updated CLAUDE.md for Go with:
  - Go development commands
  - Go architecture details
  - Go-specific SOLID implementation
  - Cross-compilation examples
  - Go conventions and best practices

#### Project Structure

- New structure following Go conventions:
  - `cmd/autonode/` - CLI entry point
  - `internal/core/` - Core interfaces and implementations
  - `internal/detectors/` - Version detectors
  - `internal/managers/` - Version managers
  - `go.mod`, `go.sum` - Go module management
  - `Makefile` - Build automation

### Removed

- All TypeScript/Node.js code and dependencies
- `package.json`, `package-lock.json`, `tsconfig.json`
- `node_modules/`, `dist/`
- npm scripts (replaced with Makefile and Go commands)

### Technical Details

- **Language**: Go 1.21+
- **Module**: github.com/matutetandil/autonode
- **Dependencies**:
  - github.com/spf13/cobra (CLI framework)
  - github.com/fatih/color (terminal colors)
- **Architecture**: SOLID principles, Dependency Injection, Strategy Pattern
- **Cross-Platform**: Native binaries for Linux, macOS (Intel & ARM), and Windows
- **Binary Size**: ~5-10MB per platform (self-contained)

### Migration Guide

**For Users:**
- Uninstall the old npm version: `npm uninstall -g autonode`
- Download the appropriate binary for your platform from GitHub releases
- Make it executable: `chmod +x autonode`
- Move to system path: `sudo mv autonode /usr/local/bin/`
- Usage remains the same: `autonode`, `autonode --check`, etc.

**For Contributors:**
- Install Go 1.21 or later
- Clone repository
- Run `make build` to compile
- See CLAUDE.md for Go development guide

## [0.1.0] - 2025-11-07

### Added

#### Core Architecture
- Implemented complete SOLID-compliant architecture with dependency injection
- Created core interfaces: `ILogger`, `IShellExecutor`, `IVersionDetector`, `IVersionManager`
- Implemented `AutoNodeService` as main orchestrator following Single Responsibility Principle
- Implemented `Logger` with colored console output using chalk
- Implemented `ShellExecutor` for cross-platform command execution using execa

#### Version Detection
- Implemented `NvmrcDetector` for reading `.nvmrc` files (priority 1)
- Implemented `NodeVersionDetector` for reading `.node-version` files (priority 2)
- Implemented `PackageJsonDetector` for reading `engines.node` from package.json (priority 3)
- Implemented `DockerfileDetector` for parsing `FROM node:X` in Dockerfile (priority 4)
- Added automatic priority-based detection chain

#### Version Managers
- Implemented `NvmManager` for nvm (Node Version Manager) support
- Implemented `NvsManager` for nvs (Node Version Switcher) support
- Implemented `VoltaManager` for Volta support
- Added automatic detection of installed version manager
- Cross-platform command execution support

#### CLI Features
- Created CLI using Commander.js with the following commands:
  - `autonode` - Detect and switch to the required Node.js version
  - `autonode --check` - Display detected version without switching
  - `autonode --force` - Reinstall version even if already present
- User-friendly colored output with success, error, warning, and tip messages
- Clear error messages with helpful installation instructions

#### Documentation
- Created comprehensive README.md with:
  - Installation instructions
  - Usage examples
  - Architecture overview
  - Extension guide for adding new detectors/managers
- Created CLAUDE.md for AI assistant guidance with:
  - Development commands
  - Architecture details
  - SOLID principles explanation
  - Common development tasks
- Created CHANGELOG.md following Keep a Changelog format

#### Project Configuration
- Configured TypeScript with ES2022 and strict mode
- Configured ESM modules (type: "module")
- Added build, dev, and clean scripts
- Created .gitignore for Node.js projects
- Configured npm package for CLI distribution

### Technical Details

- **Language**: TypeScript 5.3+
- **Module System**: ES2022 (ESM)
- **Target Platform**: Node.js >= 18.0.0
- **Dependencies**: chalk, commander, execa
- **Architecture**: SOLID principles, Dependency Injection, Strategy Pattern
- **Cross-Platform**: Supports Linux, macOS, and Windows

[0.1.0]: https://github.com/yourusername/autonode/releases/tag/v0.1.0
