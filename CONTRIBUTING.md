# Contributing to AutoNode

Thanks for your interest in contributing!

## Development Setup

### Prerequisites

- Go 1.21 or later
- Make (optional)

### Getting Started

```bash
# Clone the repository
git clone https://github.com/matutetandil/autonode.git
cd autonode

# Install dependencies
go mod tidy

# Build
make build
# or: go build -o bin/autonode ./cmd/autonode

# Run tests
make test
# or: go test ./...

# Install locally
make install
# or: sudo cp bin/autonode /usr/local/bin/
```

## Build Commands

| Command | Description |
|---------|-------------|
| `make build` | Build for current platform |
| `make build-all` | Build for all platforms |
| `make linux` | Build for Linux (amd64, arm64) |
| `make darwin` | Build for macOS (Intel, Apple Silicon) |
| `make windows` | Build for Windows |
| `make test` | Run all tests |
| `make fmt` | Format code |
| `make tidy` | Tidy dependencies |
| `make clean` | Remove build artifacts |

## Project Structure

See [Architecture](docs/architecture.md) for detailed documentation.

Key directories:
- `cmd/autonode/` - CLI entry point and commands
- `internal/core/` - Core interfaces and implementations
- `internal/detectors/` - Version detection strategies
- `internal/managers/` - Version manager integrations
- `internal/switchers/` - npm profile switchers

## Code Style

- Follow Go conventions (`go fmt`)
- One type/struct per file for SRP compliance
- All code and comments in English
- Use descriptive variable names
- Add comments explaining SOLID principles being followed

## Adding Features

### New Version Detector

1. Create `internal/detectors/my_detector.go`
2. Implement `core.VersionDetector` interface
3. Add tests in `internal/detectors/my_detector_test.go`
4. Register in `cmd/autonode/commands/run.go`

### New Version Manager

1. Create `internal/managers/my_manager.go`
2. Implement `core.VersionManager` interface
3. Add tests
4. Register in `cmd/autonode/commands/run.go`

See [Architecture](docs/architecture.md) for code examples.

## Testing

```bash
# Run all tests
go test ./...

# Run with verbose output
go test -v ./...

# Run specific package
go test -v ./internal/detectors/...

# Run with coverage
go test -cover ./...
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Pull Request Process

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/my-feature`)
3. Make your changes
4. Add/update tests
5. Run `make test` and `make fmt`
6. Commit with descriptive message
7. Push and create a Pull Request

## Commit Messages

Use clear, descriptive commit messages:

```
Add support for fnm version manager

- Implement FnmManager with install/use commands
- Add detection of fnm installation
- Add unit tests for fnm manager
```

## Releasing

Releases are automated via GitHub Actions when a tag is pushed:

```bash
# Update version in main.go and Makefile
# Update CHANGELOG.md
git add .
git commit -m "Release v0.8.0"
git tag v0.8.0
git push && git push --tags
```

## Questions?

Open an issue on GitHub for questions or suggestions.
