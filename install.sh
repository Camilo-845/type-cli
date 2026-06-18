#!/usr/bin/env bash
set -euo pipefail

REPO="Camilo-845/typingame"
BINARY="tpg"
DEFAULT_DIR="$HOME/.local/bin"

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

usage() {
    cat <<EOF
Usage: curl -fsSL https://raw.githubusercontent.com/$REPO/main/install.sh | bash [-s -- OPTIONS]

Options:
  --dir PATH     Install directory (default: $DEFAULT_DIR)
  --version TAG  Install a specific version (e.g., --version v0.1.0)
  --help         Show this message
EOF
    exit 0
}

die() {
    echo -e "${RED}Error: $*${NC}" >&2
    exit 1
}

INSTALL_DIR="$DEFAULT_DIR"
VERSION=""

while [[ $# -gt 0 ]]; do
    case "$1" in
        --dir)
            INSTALL_DIR="$2"
            shift 2
            ;;
        --version)
            VERSION="$2"
            shift 2
            ;;
        --help)
            usage
            ;;
        *)
            die "Unknown option: $1 (use --help)"
            ;;
    esac
done

OS=$(uname -s | tr '[:upper:]' '[:lower:]')
case "$OS" in
    linux|darwin) ;;
    *) die "Unsupported OS: $OS (only linux and darwin are supported)" ;;
esac

ARCH=$(uname -m)
case "$ARCH" in
    x86_64)  ARCH="amd64" ;;
    aarch64) ARCH="arm64" ;;
    arm64)   ARCH="arm64" ;;
    *) die "Unsupported architecture: $ARCH" ;;
esac

if [[ -n "$VERSION" ]]; then
    DOWNLOAD_URL="https://github.com/$REPO/releases/download/$VERSION/$BINARY-$OS-$ARCH"
else
    DOWNLOAD_URL="https://github.com/$REPO/releases/latest/download/$BINARY-$OS-$ARCH"
fi

mkdir -p "$INSTALL_DIR"

echo "Downloading $BINARY $OS/$ARCH ..."
if command -v curl &> /dev/null; then
    curl -fsSL "$DOWNLOAD_URL" -o "$INSTALL_DIR/$BINARY"
elif command -v wget &> /dev/null; then
    wget -q "$DOWNLOAD_URL" -O "$INSTALL_DIR/$BINARY"
else
    die "Neither curl nor wget found. Install one of them and try again."
fi

chmod +x "$INSTALL_DIR/$BINARY"
echo -e "${GREEN}Installed $BINARY to $INSTALL_DIR/$BINARY${NC}"

if [[ ":$PATH:" != *":$INSTALL_DIR:"* ]]; then
    echo
    echo -e "${YELLOW}Warning: $INSTALL_DIR is not in your PATH.${NC}"
    SHELL_RC=""
    case "${SHELL##*/}" in
        zsh)  SHELL_RC="$HOME/.zshrc" ;;
        bash) SHELL_RC="$HOME/.bashrc" ;;
    esac
    if [[ -n "$SHELL_RC" ]]; then
        echo "  Add this to $SHELL_RC:"
        echo -e "  ${GREEN}export PATH=\"\$PATH:$INSTALL_DIR\"${NC}"
        echo "  Then run: source $SHELL_RC"
    else
        echo "  Ensure $INSTALL_DIR is in your PATH."
    fi
fi
