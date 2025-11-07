# AutoNode

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)](https://go.dev/)
[![Release](https://img.shields.io/github/v/release/matutetandil/autonode)](https://github.com/matutetandil/autonode/releases)

Automatically detect and switch to the correct Node.js version for your project.

## Description

AutoNode is a professional CLI tool written in **Go** that automatically detects the Node.js version your project requires and switches to it using your installed version manager (nvm, nvs, or Volta).

No more manual version switching or compatibility issues - just run `autonode` and get to work!

**Why Go?** Unlike other Node.js version switchers, AutoNode is a native binary with zero dependencies. It doesn't require Node.js to run, making it the perfect tool for managing Node.js versions without circular dependencies.

## Features

- 🔍 **Smart Detection**: Automatically detects Node.js version from multiple sources:
  1. `.nvmrc` file (highest priority)
  2. `.node-version` file
  3. `engines.node` field in `package.json`
  4. `FROM node:X` in `Dockerfile` (supports numeric versions, LTS codenames like `iron`/`jod`, and tags like `lts`/`latest`)

- 🔄 **Multi-Manager Support**: Works with popular Node.js version managers:
  - **nvm** (Node Version Manager)
  - **nvs** (Node Version Switcher)
  - **Volta**

- 🎨 **User-Friendly**: Clear, colorful output with helpful messages

- 🏗️ **SOLID Architecture**: Built with clean architecture principles for easy extensibility

- 🌐 **Cross-Platform**: Native binaries for Linux, macOS (Intel & Apple Silicon), and Windows

- ⚡ **Fast & Lightweight**: Single binary, no runtime dependencies, instant startup

- 🚀 **Dynamic LTS Resolution**: Intelligently resolves Node.js LTS codenames (iron, jod, krypton) from live API
  - Automatic caching in `~/.autonode/` for offline use
  - 24-hour cache validity with automatic refresh
  - No recompilation needed when new Node versions are released

## Installation

### Automatic Installation (Recommended)

Install AutoNode with shell integration in one command:

```bash
curl -fsSL https://raw.githubusercontent.com/matutetandil/autonode/main/install.sh | bash
```

This will:
- Download the correct binary for your platform
- Install it to `/usr/local/bin/` (or `~/.local/bin/` if no sudo)
- Set up automatic version switching in your shell (bash/zsh/fish)

After installation, **restart your terminal** or run:
```bash
source ~/.bashrc   # for bash
source ~/.zshrc    # for zsh
```

Now **AutoNode will automatically switch Node.js versions** when you `cd` into a project directory!

### Manual Installation

#### Download Pre-built Binaries

Download the appropriate binary for your platform from the [releases page](https://github.com/matutetandil/autonode/releases):

**macOS (Apple Silicon)**
```bash
curl -L https://github.com/matutetandil/autonode/releases/latest/download/autonode-darwin-arm64 -o autonode
chmod +x autonode
sudo mv autonode /usr/local/bin/
```

**macOS (Intel)**
```bash
curl -L https://github.com/matutetandil/autonode/releases/latest/download/autonode-darwin-amd64 -o autonode
chmod +x autonode
sudo mv autonode /usr/local/bin/
```

**Linux (amd64)**
```bash
curl -L https://github.com/matutetandil/autonode/releases/latest/download/autonode-linux-amd64 -o autonode
chmod +x autonode
sudo mv autonode /usr/local/bin/
```

**Windows (amd64)**
Download `autonode-windows-amd64.exe` and add it to your PATH.

#### Shell Integration (Manual Setup)

To enable automatic version switching, add this to your shell config:

**Bash/Zsh** (`~/.bashrc` or `~/.zshrc`):
```bash
autonode_hook() {
  eval "$(autonode shell 2>/dev/null)"
}
autonode_cd() {
  builtin cd "$@" && autonode_hook
}
alias cd='autonode_cd'
autonode_hook  # Run on shell startup
```

**Fish** (`~/.config/fish/config.fish`):
```fish
function autonode_hook
    autonode shell 2>/dev/null | source
end
function cd
    builtin cd $argv; and autonode_hook
end
autonode_hook  # Run on shell startup
```

### Build from Source

**Prerequisites:**
- Go 1.21 or later

```bash
git clone https://github.com/matutetandil/autonode.git
cd autonode
make build
sudo make install
```

## Updating

### Automatic Update

Update to the latest version with a single command:

```bash
autonode update
```

This will:
- Check for the latest release on GitHub
- Download the correct binary for your platform
- Replace the current binary automatically

### Manual Update

Re-run the installation script:

```bash
curl -fsSL https://raw.githubusercontent.com/matutetandil/autonode/main/install.sh | bash
```

## Uninstalling

### Automatic Uninstall

Remove AutoNode and its shell integration:

```bash
curl -fsSL https://raw.githubusercontent.com/matutetandil/autonode/main/uninstall.sh | bash
```

This will:
- Remove the AutoNode binary
- Remove shell integration from your config files (`.bashrc`, `.zshrc`, etc.)
- Create backups of modified files

### Manual Uninstall

1. Remove the binary:
```bash
sudo rm /usr/local/bin/autonode
# or
rm ~/.local/bin/autonode
```

2. Remove shell integration from your shell config file (`.bashrc`, `.zshrc`, or `~/.config/fish/config.fish`):
   - Delete the section marked with `# AutoNode - automatic Node.js version switching`

3. Restart your terminal or source your shell config

## Usage

### Automatic Version Switching (With Shell Integration)

If you installed with the automatic installer or manually set up shell integration, **AutoNode works automatically**!

Just `cd` into a project directory:

```bash
cd my-project  # AutoNode automatically switches to the version in .nvmrc!
node --version  # Shows the correct version
```

No manual commands needed! AutoNode detects `.nvmrc`, `.node-version`, `package.json` (engines.node), or `Dockerfile` and switches versions automatically.

### Manual Usage

You can also run AutoNode manually in any project directory:

```bash
autonode
```

This will:
1. Detect the required Node.js version
2. Install it if not already present
3. Switch to that version

**Note:** Manual usage requires you to run `nvm use` afterwards or source your shell, as a child process cannot modify the parent shell's environment.

### Command Options

#### Check Version Only (`--check` or `-c`)

Display the detected version without switching:

```bash
autonode --check
autonode -c
```

#### Force Reinstall (`--force` or `-f`)

Reinstall the version even if it's already installed:

```bash
autonode --force
autonode -f
```

#### Version Information

Display AutoNode version:

```bash
autonode --version
autonode -v
```

#### Shell Integration Command

Output shell commands for eval (used internally by the shell hook):

```bash
autonode shell
```

This is used by the automatic integration. You can also use it manually:

```bash
eval "$(autonode shell)"
```

#### Update AutoNode

Update to the latest version:

```bash
autonode update
```

This downloads and installs the latest release from GitHub automatically.

### Examples

```bash
# In a project with .nvmrc
$ autonode
Scanning project at: /Users/dev/myproject
✓ Detected Node.js version 18.17.0 from .nvmrc
Using version manager: nvm
Node.js 18.17.0 is already installed
Switching to Node.js 18.17.0...
✓ Successfully switched to Node.js 18.17.0

# Check what version would be used
$ autonode --check
Scanning project at: /Users/dev/myproject
✓ Detected Node.js version 16.20.0 from package.json (engines.node)

# No version file found
$ autonode
Scanning project at: /Users/dev/myproject
✗ No Node.js version specification found in project
```

## Version Detection Priority

AutoNode looks for version information in this order:

1. **`.nvmrc`** - Most common, highest priority
2. **`.node-version`** - Alternative version file
3. **`package.json`** - Reads `engines.node` field
4. **`Dockerfile`** - Parses `FROM node:X` instruction
   - Supports numeric versions: `FROM node:18.17.0`, `FROM node:20`
   - Supports LTS codenames: `FROM node:iron-alpine`, `FROM node:jod`
   - Supports special tags: `FROM node:lts`, `FROM node:latest`
   - Works with all image variants: `-alpine`, `-slim`, `-bullseye`, etc.

The first source that provides a version will be used.

### Dockerfile LTS Codenames

AutoNode dynamically resolves Node.js LTS codenames by fetching release data from the official Node.js API. Supported codenames include:

| Codename | Node Version | LTS Start |
|----------|--------------|-----------|
| krypton  | 24           | 2025-10   |
| jod      | 22           | 2024-10   |
| iron     | 20           | 2023-10   |
| hydrogen | 18           | 2022-10   |
| gallium  | 16           | 2021-10   |
| fermium  | 14           | 2020-10   |

The mapping is cached locally in `~/.autonode/node-releases.json` and refreshed every 24 hours.

## Supported Version Managers

### nvm (Node Version Manager)

Install from: https://github.com/nvm-sh/nvm

```bash
autonode  # Will use: nvm install <version> && nvm use <version>
```

### nvs (Node Version Switcher)

Install from: https://github.com/jasongin/nvs

```bash
autonode  # Will use: nvs add <version> && nvs use <version>
```

### Volta

Install from: https://volta.sh

```bash
autonode  # Will use: volta install node@<version> && volta pin node@<version>
```

## Architecture

AutoNode is built with **Go** following SOLID principles and clean architecture:

### Project Structure

```
autonode/
├── cmd/autonode/          # CLI entry point
│   └── main.go            # Main CLI with dependency injection
├── internal/
│   ├── core/              # Core abstractions and implementations
│   │   ├── config.go              # Config type
│   │   ├── detection_result.go   # DetectionResult type
│   │   ├── operation_result.go   # OperationResult type
│   │   ├── logger.go              # Logger interface
│   │   ├── console_logger.go     # Logger implementation
│   │   ├── shell_executor.go     # ShellExecutor interface
│   │   ├── exec_shell.go         # ShellExecutor implementation
│   │   ├── version_detector.go   # VersionDetector interface
│   │   ├── version_manager.go    # VersionManager interface
│   │   └── service.go            # AutoNodeService orchestrator
│   ├── detectors/         # Version detection implementations
│   │   ├── nvmrc.go              # .nvmrc detector
│   │   ├── node_version.go       # .node-version detector
│   │   ├── package_json.go       # package.json detector
│   │   └── dockerfile.go         # Dockerfile detector
│   └── managers/          # Version manager implementations
│       ├── nvm.go                # nvm manager
│       ├── nvs.go                # nvs manager
│       └── volta.go              # Volta manager
├── go.mod                 # Go module definition
├── go.sum                 # Go dependencies
└── Makefile              # Build automation
```

### SOLID Principles in Go

- **Single Responsibility**: Each file/type has one clear purpose
- **Open/Closed**: Easy to add new detectors or managers without modifying existing code
- **Liskov Substitution**: All managers and detectors are interchangeable through interfaces
- **Interface Segregation**: Small, focused interfaces
- **Dependency Inversion**: Depends on interfaces, not concrete implementations

## Extending AutoNode

### Adding a New Version Detector

1. Create a new file in `internal/detectors/`:

```go
package detectors

import (
    "github.com/matutetandil/autonode/internal/core"
)

type MyCustomDetector struct{}

func NewMyCustomDetector() *MyCustomDetector {
    return &MyCustomDetector{}
}

func (d *MyCustomDetector) Detect(projectPath string) (core.DetectionResult, error) {
    // Your detection logic here
    return core.DetectionResult{
        Found:   true,
        Version: "18.0.0",
        Source:  "my-custom-file",
    }, nil
}

func (d *MyCustomDetector) GetPriority() int {
    return 5 // Set priority
}

func (d *MyCustomDetector) GetSourceName() string {
    return "my-custom-file"
}
```

2. Add it in `cmd/autonode/main.go`:

```go
detectorsList := []core.VersionDetector{
    detectors.NewNvmrcDetector(),
    detectors.NewNodeVersionDetector(),
    detectors.NewPackageJsonDetector(),
    detectors.NewDockerfileDetector(),
    detectors.NewMyCustomDetector(), // Add your detector
}
```

### Adding a New Version Manager

1. Create a new file in `internal/managers/`:

```go
package managers

import (
    "github.com/matutetandil/autonode/internal/core"
)

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

// Implement other methods: IsVersionInstalled, InstallVersion, UseVersion
```

2. Add it in `cmd/autonode/main.go`:

```go
managersList := []core.VersionManager{
    managers.NewNvmManager(shell),
    managers.NewNvsManager(shell),
    managers.NewVoltaManager(shell),
    managers.NewMyManager(shell), // Add your manager
}
```

## Development

### Prerequisites

- Go 1.21 or later
- Make (optional, for using Makefile commands)

### Build Commands

```bash
# Build for current platform
make build

# Build for all platforms
make build-all

# Build for specific platform
make linux    # Linux (amd64 and arm64)
make darwin   # macOS (Intel and Apple Silicon)
make windows  # Windows (amd64)

# Install locally
make install

# Clean build artifacts
make clean

# Run tests
make test

# Format code
make fmt

# Tidy dependencies
make tidy
```

### Manual Build

```bash
# Build for current platform
go build -o bin/autonode ./cmd/autonode

# Build for specific platform
GOOS=linux GOARCH=amd64 go build -o bin/autonode-linux-amd64 ./cmd/autonode
GOOS=darwin GOARCH=arm64 go build -o bin/autonode-darwin-arm64 ./cmd/autonode
GOOS=windows GOARCH=amd64 go build -o bin/autonode-windows-amd64.exe ./cmd/autonode
```

### Project Dependencies

- [cobra](https://github.com/spf13/cobra) - CLI framework
- [color](https://github.com/fatih/color) - Terminal colors

All dependencies are managed via `go.mod`.

## License

MIT

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Roadmap

- [x] Add unit tests
- [x] Add CI/CD pipeline
- [x] Shell integration for auto-switching on directory change
- [ ] Add support for more version managers (fnm, asdf)
- [ ] Package managers (Homebrew, apt, etc.)
