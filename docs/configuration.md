# Configuration

AutoNode can be configured at multiple levels: per-command flags, per-project files, and global settings.

## Command Line Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--check` | `-c` | Only display detected version, don't switch |
| `--force` | `-f` | Reinstall version even if already installed |
| `--no-update-check` | | Disable automatic update check (useful for CI/CD) |
| `--version` | `-v` | Display AutoNode version |
| `--help` | `-h` | Display help |

## Version Detection Priority

AutoNode looks for Node.js version in this order (first match wins):

| Priority | Source | Example |
|----------|--------|---------|
| 1 | `.autonode.yml` | `nodeVersion: 20` |
| 2 | `.nvmrc` | `18.17.0` |
| 3 | `.node-version` | `20.10.0` |
| 4 | `package.json` | `"engines": { "node": ">=18" }` |
| 5 | `Dockerfile` | `FROM node:20-alpine` |

## Per-Project Configuration

### `.autonode.yml`

Create this file in your project root for the highest priority configuration:

```yaml
# Node.js version (optional)
nodeVersion: 20

# npm profile to switch to (optional)
npmProfile: work
```

Use the `config` command to manage this file:

```bash
autonode config --node 20           # Set Node version
autonode config --profile work      # Set npm profile
autonode config --show              # Show current config
autonode config --remove            # Remove .autonode.yml
```

### `package.json`

Add an `autonode` field:

```json
{
  "name": "my-project",
  "autonode": {
    "npmProfile": "work"
  },
  "engines": {
    "node": ">=18"
  }
}
```

Note: `engines.node` is used for version detection, `autonode.npmProfile` for profile switching.

## Global Configuration

Global settings are stored in `~/.autonode/config.json`:

```json
{
  "disableUpdateCheck": false,
  "updateCheckIntervalDays": 7
}
```

| Setting | Type | Default | Description |
|---------|------|---------|-------------|
| `disableUpdateCheck` | boolean | `false` | Disable automatic update checks |
| `updateCheckIntervalDays` | number | `7` | Days between update checks |

## Dockerfile Detection

AutoNode supports various Dockerfile formats:

```dockerfile
# Numeric versions
FROM node:18
FROM node:18.17.0
FROM node:20-alpine

# LTS codenames
FROM node:iron          # Node 20
FROM node:jod           # Node 22
FROM node:hydrogen      # Node 18

# Special tags
FROM node:lts           # Current LTS
FROM node:latest        # Latest stable
```

### Supported LTS Codenames

| Codename | Node Version |
|----------|--------------|
| krypton | 24 |
| jod | 22 |
| iron | 20 |
| hydrogen | 18 |
| gallium | 16 |
| fermium | 14 |

Codenames are resolved dynamically from the Node.js API and cached in `~/.autonode/node-releases.json`.

## npm Profile Management

AutoNode can automatically switch npm profiles (registry configurations) per project.

### Supported Tools

- **[npmrc](https://github.com/deoxxa/npmrc)** - Most popular
- **[ts-npmrc](https://github.com/darsi-an/ts-npmrc)** - TypeScript version
- **[rc-manager](https://github.com/Lalaluka/rc-manager)** - Supports npm and yarn

### Setup

1. Install a profile tool:
   ```bash
   npm install -g npmrc
   ```

2. Create profiles:
   ```bash
   npmrc -c work      # Create 'work' profile
   npmrc -c personal  # Create 'personal' profile
   ```

3. Configure your project:
   ```yaml
   # .autonode.yml
   npmProfile: work
   ```

### Behavior

- **Silent mode**: No warnings if tool not installed or profile not configured
- **Auto-discovery**: Finds tools installed in any nvm Node version
- **Tool priority**: npmrc > ts-npmrc > rc-manager

## Cache Files

AutoNode stores cache files in `~/.autonode/`:

| File | Purpose | Validity |
|------|---------|----------|
| `node-releases.json` | LTS codename mappings | 24 hours |
| `update-check.json` | Update check results | 7 days (configurable) |
| `config.json` | Global settings | Permanent |

## Environment Variables

| Variable | Description |
|----------|-------------|
| `NVM_DIR` | Custom nvm installation directory |
