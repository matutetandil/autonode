# AutoNode

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)](https://go.dev/)
[![Release](https://img.shields.io/github/v/release/matutetandil/autonode)](https://github.com/matutetandil/autonode/releases)

Automatically detect and switch to the correct Node.js version for your project.

## What it does

AutoNode reads your project's Node.js version requirement (from `.nvmrc`, `package.json`, `Dockerfile`, etc.) and automatically switches to it using your installed version manager (nvm, nvs, or Volta).

**Why Go?** AutoNode is a native binary with zero dependencies. It doesn't require Node.js to run - no circular dependency problem.

## Installation

### Quick Install (Recommended)

```bash
curl -fsSL https://raw.githubusercontent.com/matutetandil/autonode/main/install.sh | bash
```

This installs the binary and sets up automatic version switching when you `cd` into directories.

**Restart your terminal** after installation.

### Manual Install

Download from [releases](https://github.com/matutetandil/autonode/releases):

```bash
# macOS (Apple Silicon)
curl -L https://github.com/matutetandil/autonode/releases/latest/download/autonode-darwin-arm64.tar.gz | tar xz
sudo mv autonode /usr/local/bin/

# macOS (Intel)
curl -L https://github.com/matutetandil/autonode/releases/latest/download/autonode-darwin-amd64.tar.gz | tar xz
sudo mv autonode /usr/local/bin/

# Linux
curl -L https://github.com/matutetandil/autonode/releases/latest/download/autonode-linux-amd64.tar.gz | tar xz
sudo mv autonode /usr/local/bin/
```

## Usage

### Automatic (with shell integration)

Just `cd` into a project - AutoNode switches versions automatically:

```bash
cd my-project    # Switches to version in .nvmrc
node --version   # v18.17.0
```

### Manual

```bash
autonode              # Detect and switch version
autonode --check      # Show detected version without switching
autonode --force      # Force reinstall even if installed
autonode update       # Update AutoNode to latest version
```

### Configure a directory

```bash
autonode config --node 20           # Set Node version for current directory
autonode config --profile work      # Set npm profile
autonode config --show              # Show configuration
```

## Version Detection

AutoNode checks these sources in order:

1. `.autonode.yml` - `nodeVersion: 20`
2. `.nvmrc` - `18.17.0`
3. `.node-version` - `20.10.0`
4. `package.json` - `"engines": { "node": ">=18" }`
5. `Dockerfile` - `FROM node:20-alpine`

## Supported Version Managers

- **[nvm](https://github.com/nvm-sh/nvm)** - Node Version Manager
- **[nvs](https://github.com/jasongin/nvs)** - Node Version Switcher
- **[Volta](https://volta.sh)** - JavaScript Tool Manager

## Updating

```bash
autonode update
```

AutoNode checks for updates weekly and shows a notification when available. Disable with `--no-update-check`.

## Uninstalling

```bash
curl -fsSL https://raw.githubusercontent.com/matutetandil/autonode/main/uninstall.sh | bash
```

## Documentation

- **[Configuration](docs/configuration.md)** - All options, version detection, npm profiles
- **[Architecture](docs/architecture.md)** - SOLID principles, how to extend
- **[Contributing](CONTRIBUTING.md)** - Development setup, pull requests
- **[Changelog](CHANGELOG.md)** - Version history

## License

MIT
