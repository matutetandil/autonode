#!/bin/bash
# AutoNode Uninstallation Script
# Usage: curl -fsSL https://raw.githubusercontent.com/matutetandil/autonode/main/uninstall.sh | bash
#    or: bash uninstall.sh

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Hook marker to identify AutoNode lines
HOOK_MARKER="# AutoNode - automatic Node.js version switching"

# Remove binary
remove_binary() {
    echo -e "${BLUE}→ Removing AutoNode binary...${NC}"

    local removed=false

    # Check common installation locations
    if [ -f "/usr/local/bin/autonode" ]; then
        if rm "/usr/local/bin/autonode" 2>/dev/null; then
            echo -e "${GREEN}✓ Removed /usr/local/bin/autonode${NC}"
            removed=true
        else
            echo -e "${YELLOW}→ Removing from /usr/local/bin (requires sudo)...${NC}"
            sudo rm "/usr/local/bin/autonode"
            echo -e "${GREEN}✓ Removed /usr/local/bin/autonode${NC}"
            removed=true
        fi
    fi

    if [ -f "$HOME/.local/bin/autonode" ]; then
        rm "$HOME/.local/bin/autonode"
        echo -e "${GREEN}✓ Removed $HOME/.local/bin/autonode${NC}"
        removed=true
    fi

    if [ "$removed" = false ]; then
        echo -e "${YELLOW}⚠ AutoNode binary not found in standard locations${NC}"
        echo -e "${YELLOW}  Please manually remove it if installed elsewhere${NC}"
    fi
}

# Remove shell integration from a config file
remove_from_config() {
    local config_file="$1"
    local config_name="$2"

    if [ ! -f "$config_file" ]; then
        return
    fi

    if ! grep -q "$HOOK_MARKER" "$config_file" 2>/dev/null; then
        return
    fi

    echo -e "${BLUE}→ Removing shell integration from $config_name...${NC}"

    # Create a backup
    cp "$config_file" "${config_file}.autonode-backup"

    # Remove AutoNode section (from marker to next empty line or EOF)
    # This uses awk to remove the AutoNode block
    awk '
        BEGIN { in_autonode = 0 }
        /# AutoNode - automatic Node\.js version switching/ { in_autonode = 1; next }
        in_autonode == 1 && /^[[:space:]]*$/ { in_autonode = 0; next }
        in_autonode == 0 { print }
    ' "$config_file" > "${config_file}.tmp"

    mv "${config_file}.tmp" "$config_file"

    echo -e "${GREEN}✓ Removed shell integration from $config_name${NC}"
    echo -e "${YELLOW}  Backup saved: ${config_file}.autonode-backup${NC}"
}

# Remove shell integrations
remove_shell_integration() {
    echo -e "${BLUE}→ Removing shell integrations...${NC}"

    # Bash
    remove_from_config "$HOME/.bashrc" "~/.bashrc"
    remove_from_config "$HOME/.bash_profile" "~/.bash_profile"

    # Zsh
    remove_from_config "$HOME/.zshrc" "~/.zshrc"
    remove_from_config "$HOME/.zprofile" "~/.zprofile"

    # Fish
    remove_from_config "$HOME/.config/fish/config.fish" "~/.config/fish/config.fish"

    echo -e "${GREEN}✓ Shell integration cleanup complete${NC}"
}

# Main uninstallation flow
main() {
    echo -e "${GREEN}"
    echo "╔═══════════════════════════════════════╗"
    echo "║    AutoNode Uninstallation            ║"
    echo "╚═══════════════════════════════════════╝"
    echo -e "${NC}"

    # Confirm uninstallation
    echo -e "${YELLOW}This will remove AutoNode and its shell integration.${NC}"
    read -p "Continue? (y/N) " -n 1 -r
    echo

    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        echo -e "${BLUE}Uninstallation cancelled${NC}"
        exit 0
    fi

    remove_binary
    remove_shell_integration

    echo ""
    echo -e "${GREEN}✓ AutoNode has been uninstalled!${NC}"
    echo ""
    echo -e "${BLUE}Note:${NC}"
    echo "  • Restart your terminal or source your shell config to complete removal"
    echo "  • Config backups were created with .autonode-backup extension"
    echo "  • To restore: mv ~/.bashrc.autonode-backup ~/.bashrc (for example)"
    echo ""
    echo -e "${YELLOW}Thanks for using AutoNode!${NC}"
}

main
