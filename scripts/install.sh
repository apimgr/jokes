#!/usr/bin/env bash
#===============================================================================
# Jokes API - Installation Script (OS-agnostic)
# Supports: Linux, BSD, macOS
#===============================================================================

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Variables
PROJECT_NAME="jokes"
BINARY_NAME="jokes-api"
VERSION="1.0.0"

# Detect OS
detect_os() {
    case "$(uname -s)" in
        Linux*)     OS="linux";;
        Darwin*)    OS="macos";;
        FreeBSD*)   OS="freebsd";;
        OpenBSD*)   OS="openbsd";;
        *)          OS="unknown";;
    esac

    case "$(uname -m)" in
        x86_64)     ARCH="amd64";;
        aarch64)    ARCH="arm64";;
        arm64)      ARCH="arm64";;
        *)          ARCH="unknown";;
    esac

    echo -e "${GREEN}✓${NC} Detected: $OS/$ARCH"
}

# Check if running as root/sudo
check_root() {
    if [ "$EUID" -ne 0 ] && [ "$OS" != "macos" ]; then
        echo -e "${RED}✗${NC} Please run as root or with sudo"
        exit 1
    fi
}

# Delegate to OS-specific script
delegate_install() {
    SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

    case "$OS" in
        linux|freebsd|openbsd)
            bash "$SCRIPT_DIR/linux.sh"
            ;;
        macos)
            bash "$SCRIPT_DIR/macos.sh"
            ;;
        *)
            echo -e "${RED}✗${NC} Unsupported OS: $OS"
            exit 1
            ;;
    esac
}

# Main installation
main() {
    echo "🎭 Jokes API Installer v$VERSION"
    echo "=================================="
    echo ""

    detect_os
    check_root
    delegate_install

    echo ""
    echo -e "${GREEN}✅ Installation complete!${NC}"
    echo ""
    echo "Next steps:"
    echo "1. Check service status"
    echo "2. Visit http://localhost:64xxx/ (check config for actual port)"
    echo "3. View API docs at http://localhost:64xxx/docs"
}

main "$@"
