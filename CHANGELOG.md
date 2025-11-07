# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

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
