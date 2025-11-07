#!/bin/bash
# AutoNode Installation Script
# Usage: curl -fsSL https://raw.githubusercontent.com/matutetandil/autonode/main/install.sh | bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# GitHub repository
REPO="matutetandil/autonode"
VERSION="${AUTONODE_VERSION:-latest}"

# Detect OS and architecture
detect_platform() {
    local os=$(uname -s | tr '[:upper:]' '[:lower:]')
    local arch=$(uname -m)

    case "$os" in
        darwin)
            OS="darwin"
            ;;
        linux)
            OS="linux"
            ;;
        *)
            echo -e "${RED}Error: Unsupported operating system: $os${NC}"
            exit 1
            ;;
    esac

    case "$arch" in
        x86_64|amd64)
            ARCH="amd64"
            ;;
        aarch64|arm64)
            ARCH="arm64"
            ;;
        *)
            echo -e "${RED}Error: Unsupported architecture: $arch${NC}"
            exit 1
            ;;
    esac

    # Archive name with extension
    if [ "$OS" = "windows" ]; then
        ARCHIVE_NAME="autonode-${OS}-${ARCH}.zip"
    else
        ARCHIVE_NAME="autonode-${OS}-${ARCH}.tar.gz"
    fi

    BINARY_NAME="autonode"
    if [ "$OS" = "windows" ]; then
        BINARY_NAME="${BINARY_NAME}.exe"
    fi
}

# Download and install binary
install_binary() {
    local download_url="https://github.com/${REPO}/releases/download/${VERSION}/${ARCHIVE_NAME}"

    echo -e "${BLUE}→ Downloading AutoNode ${VERSION} for ${OS}/${ARCH}...${NC}"

    # Determine install directory
    if [ -w "/usr/local/bin" ]; then
        INSTALL_DIR="/usr/local/bin"
    else
        INSTALL_DIR="$HOME/.local/bin"
        mkdir -p "$INSTALL_DIR"

        # Add to PATH if not already there
        if [[ ":$PATH:" != *":$INSTALL_DIR:"* ]]; then
            echo -e "${YELLOW}⚠ Note: $INSTALL_DIR is not in your PATH${NC}"
            echo -e "${YELLOW}  You may need to add it to your shell configuration${NC}"
        fi
    fi

    local tmp_dir="/tmp/autonode-install-${RANDOM}"
    mkdir -p "$tmp_dir"
    local tmp_archive="${tmp_dir}/${ARCHIVE_NAME}"

    # Download with curl or wget
    if command -v curl >/dev/null 2>&1; then
        curl -fsSL "$download_url" -o "$tmp_archive"
    elif command -v wget >/dev/null 2>&1; then
        wget -q "$download_url" -O "$tmp_archive"
    else
        echo -e "${RED}Error: Neither curl nor wget found. Please install one of them.${NC}"
        exit 1
    fi

    echo -e "${BLUE}→ Extracting archive...${NC}"

    # Extract archive based on format
    if [ "$OS" = "windows" ]; then
        # Extract zip
        unzip -q "$tmp_archive" -d "$tmp_dir"
    else
        # Extract tar.gz
        tar -xzf "$tmp_archive" -C "$tmp_dir"
    fi

    # Make executable and move to install directory
    local extracted_binary="${tmp_dir}/${BINARY_NAME}"
    chmod +x "$extracted_binary"

    if [ "$INSTALL_DIR" = "/usr/local/bin" ]; then
        # May need sudo
        if mv "$extracted_binary" "$INSTALL_DIR/autonode" 2>/dev/null; then
            :
        else
            echo -e "${YELLOW}→ Installing to /usr/local/bin (requires sudo)...${NC}"
            sudo mv "$extracted_binary" "$INSTALL_DIR/autonode"
        fi
    else
        mv "$extracted_binary" "$INSTALL_DIR/autonode"
    fi

    # Clean up
    rm -rf "$tmp_dir"

    echo -e "${GREEN}✓ AutoNode installed to ${INSTALL_DIR}/autonode${NC}"
}

# Setup shell integration
setup_shell_integration() {
    echo -e "${BLUE}→ Setting up shell integration...${NC}"

    # Detect user's shell
    local user_shell=$(basename "$SHELL")

    case "$user_shell" in
        bash)
            setup_bash
            ;;
        zsh)
            setup_zsh
            ;;
        fish)
            setup_fish
            ;;
        *)
            echo -e "${YELLOW}⚠ Unknown shell: $user_shell${NC}"
            echo -e "${YELLOW}  Please manually add the hook to your shell configuration${NC}"
            show_manual_instructions
            return
            ;;
    esac
}

# Bash integration
setup_bash() {
    local bashrc="$HOME/.bashrc"
    local hook_marker="# AutoNode - automatic Node.js version switching"

    if grep -q "$hook_marker" "$bashrc" 2>/dev/null; then
        echo -e "${GREEN}✓ Shell integration already configured in $bashrc${NC}"
        return
    fi

    echo -e "${BLUE}→ Adding hook to $bashrc...${NC}"

    cat >> "$bashrc" << 'EOF'

# AutoNode - automatic Node.js version switching
autonode_hook() {
  eval "$(autonode shell 2>/dev/null)"
}
autonode_cd() {
  builtin cd "$@" && autonode_hook
}
alias cd='autonode_cd'
autonode_hook  # Run on shell startup
EOF

    echo -e "${GREEN}✓ Shell integration configured${NC}"
    echo -e "${YELLOW}→ Run: source ~/.bashrc  (or restart your terminal)${NC}"
}

# Zsh integration
setup_zsh() {
    local zshrc="$HOME/.zshrc"
    local hook_marker="# AutoNode - automatic Node.js version switching"

    if grep -q "$hook_marker" "$zshrc" 2>/dev/null; then
        echo -e "${GREEN}✓ Shell integration already configured in $zshrc${NC}"
        return
    fi

    echo -e "${BLUE}→ Adding hook to $zshrc...${NC}"

    cat >> "$zshrc" << 'EOF'

# AutoNode - automatic Node.js version switching
autonode_hook() {
  eval "$(autonode shell 2>/dev/null)"
}
autonode_cd() {
  builtin cd "$@" && autonode_hook
}
alias cd='autonode_cd'
autonode_hook  # Run on shell startup
EOF

    echo -e "${GREEN}✓ Shell integration configured${NC}"
    echo -e "${YELLOW}→ Run: source ~/.zshrc  (or restart your terminal)${NC}"
}

# Fish integration
setup_fish() {
    local fish_config="$HOME/.config/fish/config.fish"
    local hook_marker="# AutoNode - automatic Node.js version switching"

    mkdir -p "$(dirname "$fish_config")"

    if grep -q "$hook_marker" "$fish_config" 2>/dev/null; then
        echo -e "${GREEN}✓ Shell integration already configured in $fish_config${NC}"
        return
    fi

    echo -e "${BLUE}→ Adding hook to $fish_config...${NC}"

    cat >> "$fish_config" << 'EOF'

# AutoNode - automatic Node.js version switching
function autonode_hook
    autonode shell 2>/dev/null | source
end
function cd
    builtin cd $argv; and autonode_hook
end
autonode_hook  # Run on shell startup
EOF

    echo -e "${GREEN}✓ Shell integration configured${NC}"
    echo -e "${YELLOW}→ Run: source ~/.config/fish/config.fish  (or restart your terminal)${NC}"
}

# Show manual instructions
show_manual_instructions() {
    cat << 'EOF'

Manual Setup Instructions:
--------------------------
Add the following to your shell configuration file:

For Bash (~/.bashrc):
  autonode_hook() {
    eval "$(autonode shell 2>/dev/null)"
  }
  autonode_cd() {
    builtin cd "$@" && autonode_hook
  }
  alias cd='autonode_cd'
  autonode_hook

For Zsh (~/.zshrc):
  Same as Bash

For Fish (~/.config/fish/config.fish):
  function autonode_hook
    autonode shell 2>/dev/null | source
  end
  function cd
    builtin cd $argv; and autonode_hook
  end
  autonode_hook

EOF
}

# Main installation flow
main() {
    echo -e "${GREEN}"
    echo "╔═══════════════════════════════════════╗"
    echo "║      AutoNode Installation            ║"
    echo "╚═══════════════════════════════════════╝"
    echo -e "${NC}"

    detect_platform
    install_binary
    setup_shell_integration

    echo ""
    echo -e "${GREEN}✓ Installation complete!${NC}"
    echo ""
    echo -e "${BLUE}How it works:${NC}"
    echo "  • AutoNode will automatically detect and switch Node.js versions"
    echo "  • Every time you cd into a directory with .nvmrc (or other version files)"
    echo "  • No manual 'nvm use' needed anymore!"
    echo ""
    echo -e "${BLUE}Usage:${NC}"
    echo "  autonode           - Manually switch to detected version"
    echo "  autonode --check   - Check detected version without switching"
    echo "  autonode --help    - Show all options"
    echo ""
    echo -e "${YELLOW}Note: Restart your terminal or source your shell config to activate${NC}"
}

main
