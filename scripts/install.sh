#!/bin/sh
# install.sh - Install dotfiles binary from GitHub releases
# Usage: curl -fsSL https://raw.githubusercontent.com/yveskaufmann/dotfiles-cli/main/scripts/install.sh | sh
# Usage with flags: curl -fsSL https://raw.githubusercontent.com/yveskaufmann/dotfiles-cli/main/scripts/install.sh | sh -s -- --home
# Usage with token: curl -fsSL https://raw.githubusercontent.com/yveskaufmann/dotfiles-cli/main/scripts/install.sh | sh -s -- --github-token YOUR_TOKEN

set -e

# Configuration
REPO_OWNER="yveskaufmann"
REPO_NAME="dotfiles-cli"
BINARY_NAME="dotfiles"
INSTALL_DIR=""
FORCE_HOME=false
FORCE_SYSTEM=false
GITHUB_TOKEN=""

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Helper functions
print_info() {
    printf "${CYAN}[INFO]${NC} %s\n" "$1"
}

print_success() {
    printf "${GREEN}[SUCCESS]${NC} %s\n" "$1"
}

print_warning() {
    printf "${YELLOW}[WARNING]${NC} %s\n" "$1"
}

print_error() {
    printf "${RED}[ERROR]${NC} %s\n" "$1" >&2
}

cleanup() {
    if [ -n "$TEMP_DIR" ] && [ -d "$TEMP_DIR" ]; then
        rm -rf "$TEMP_DIR"
    fi
}

trap cleanup EXIT

# Parse command line arguments
while [ $# -gt 0 ]; do
    case "$1" in
        --home)
            FORCE_HOME=true
            shift
            ;;
        --system)
            FORCE_SYSTEM=true
            shift
            ;;
        --dir)
            INSTALL_DIR="$2"
            shift 2
            ;;
        --github-token)
            GITHUB_TOKEN="$2"
            shift 2
            ;;
        *)
            print_error "Unknown option: $1"
            echo "Usage: $0 [--home] [--system] [--dir DIR] [--github-token TOKEN]"
            exit 1
            ;;
    esac
done

# Check for GITHUB_TOKEN environment variable if no token parameter provided
if [ -z "$GITHUB_TOKEN" ] && [ -n "${GITHUB_TOKEN:-}" ]; then
    GITHUB_TOKEN="${GITHUB_TOKEN:-}"
fi

# Detect OS
detect_os() {
    OS=$(uname -s)
    case "$OS" in
        Linux*)
            OS="linux"
            ;;
        Darwin*)
            OS="darwin"
            ;;
        MINGW*|MSYS*|CYGWIN*)
            OS="windows"
            ;;
        *)
            print_error "Unsupported operating system: $OS"
            exit 1
            ;;
    esac
    return 0
}

# Detect architecture
detect_arch() {
    ARCH=$(uname -m)
    case "$ARCH" in
        x86_64|amd64)
            ARCH="amd64"
            ;;
        aarch64|arm64)
            ARCH="arm64"
            ;;
        *)
            print_error "Unsupported architecture: $ARCH"
            exit 1
            ;;
    esac
    return 0
}

# Determine install directory
determine_install_dir() {
    if [ -n "$INSTALL_DIR" ]; then
        echo "$INSTALL_DIR"
        return
    fi

    if [ "$FORCE_SYSTEM" = true ]; then
        if [ -w "/usr/local/bin" ]; then
            echo "/usr/local/bin"
        else
            print_error "/usr/local/bin is not writable. Try with sudo or use --home flag."
            exit 1
        fi
    elif [ "$FORCE_HOME" = true ]; then
        echo "$HOME/.local/bin"
    else
        # Default: try system directory first, fallback to user home if not writable
        if [ -w "/usr/local/bin" ]; then
            echo "/usr/local/bin"
        else
            echo "$HOME/.local/bin"
        fi
    fi
}

# Check if directory is in PATH
is_in_path() {
    case ":$PATH:" in
        *:"$1":*)
            return 0
            ;;
        *)
            return 1
            ;;
    esac
}

# Get latest release version
get_latest_version() {
    # Try stable latest first
    if command -v curl >/dev/null 2>&1; then
        if [ -n "$GITHUB_TOKEN" ]; then
            VERSION=$( { curl -fsSL -H "Authorization: Bearer $GITHUB_TOKEN" \
                      "https://api.github.com/repos/${REPO_OWNER}/${REPO_NAME}/releases/latest" 2>/dev/null || true; } | \
                      grep '"tag_name":' | \
                      head -n 1 | \
                      sed -E 's/.*"tag_name": "([^"]+)".*/\1/')
        else
            VERSION=$( { curl -fsSL "https://api.github.com/repos/${REPO_OWNER}/${REPO_NAME}/releases/latest" 2>/dev/null || true; } | \
                      grep '"tag_name":' | \
                      head -n 1 | \
                      sed -E 's/.*"tag_name": "([^"]+)".*/\1/')
        fi
    elif command -v wget >/dev/null 2>&1; then
        if [ -n "$GITHUB_TOKEN" ]; then
            VERSION=$(wget --header="Authorization: Bearer $GITHUB_TOKEN" \
                      -qO- "https://api.github.com/repos/${REPO_OWNER}/${REPO_NAME}/releases/latest" 2>/dev/null | \
                      grep '"tag_name":' | \
                      head -n 1 | \
                      sed -E 's/.*"tag_name": "([^"]+)".*/\1/')
        else
            VERSION=$(wget -qO- "https://api.github.com/repos/${REPO_OWNER}/${REPO_NAME}/releases/latest" 2>/dev/null | \
                      grep '"tag_name":' | \
                      head -n 1 | \
                      sed -E 's/.*"tag_name": "([^"]+)".*/\1/')
        fi
    else
        print_error "Neither curl nor wget found. Please install one of them."
        exit 1
    fi

    # Fallback for repos with no stable release (e.g. only prereleases)
    if [ -z "$VERSION" ]; then
        print_warning "No stable latest release found; falling back to newest release tag"

        if command -v curl >/dev/null 2>&1; then
            if [ -n "$GITHUB_TOKEN" ]; then
                VERSION=$( { curl -fsSL -H "Authorization: Bearer $GITHUB_TOKEN" \
                          "https://api.github.com/repos/${REPO_OWNER}/${REPO_NAME}/releases?per_page=1" 2>/dev/null || true; } | \
                          grep '"tag_name":' | \
                          head -n 1 | \
                          sed -E 's/.*"tag_name": "([^"]+)".*/\1/')
            else
                VERSION=$( { curl -fsSL "https://api.github.com/repos/${REPO_OWNER}/${REPO_NAME}/releases?per_page=1" 2>/dev/null || true; } | \
                          grep '"tag_name":' | \
                          head -n 1 | \
                          sed -E 's/.*"tag_name": "([^"]+)".*/\1/')
            fi
        elif command -v wget >/dev/null 2>&1; then
            if [ -n "$GITHUB_TOKEN" ]; then
                VERSION=$(wget --header="Authorization: Bearer $GITHUB_TOKEN" \
                          -qO- "https://api.github.com/repos/${REPO_OWNER}/${REPO_NAME}/releases?per_page=1" 2>/dev/null | \
                          grep '"tag_name":' | \
                          head -n 1 | \
                          sed -E 's/.*"tag_name": "([^"]+)".*/\1/')
            else
                VERSION=$(wget -qO- "https://api.github.com/repos/${REPO_OWNER}/${REPO_NAME}/releases?per_page=1" 2>/dev/null | \
                          grep '"tag_name":' | \
                          head -n 1 | \
                          sed -E 's/.*"tag_name": "([^"]+)".*/\1/')
            fi
        fi
    fi

    if [ -z "$VERSION" ]; then
        print_error "Failed to fetch latest version"
        exit 1
    fi

    VERSION="${VERSION#v}"  

    return 0
}

# Download and install
install_binary() {
    detect_os
    detect_arch
    get_latest_version

    INSTALL_DIR=$(determine_install_dir)

    print_info "Detected: OS=$OS, ARCH=$ARCH"
    print_info "Latest version: v$VERSION"
    print_info "Install directory: $INSTALL_DIR"

    # Create temp directory
    TEMP_DIR=$(mktemp -d)
    cd "$TEMP_DIR"

    # Construct archive name
    ARCHIVE_NAME="${REPO_NAME}_${VERSION}_${OS}_${ARCH}"
    if [ "$OS" = "windows" ]; then
        ARCHIVE_EXT="zip"
    else
        ARCHIVE_EXT="tar.gz"
    fi

    # Download archive - use API endpoint when token is provided
    if [ -n "$GITHUB_TOKEN" ]; then
        print_info "Fetching release asset info..."
        
        # Get release info via API to find asset ID
        RELEASE_API_URL="https://api.github.com/repos/${REPO_OWNER}/${REPO_NAME}/releases/tags/v${VERSION}"
        
        if command -v curl >/dev/null 2>&1; then
            RELEASE_JSON=$(curl -fsSL -H "Authorization: Bearer $GITHUB_TOKEN" "$RELEASE_API_URL")
        elif command -v wget >/dev/null 2>&1; then
            RELEASE_JSON=$(wget --header="Authorization: Bearer $GITHUB_TOKEN" -qO- "$RELEASE_API_URL")
        else
            print_error "Neither curl nor wget found. Please install one of them."
            exit 1
        fi
        
        # Extract asset ID for matching archive name using awk
        ASSET_NAME_PATTERN="${ARCHIVE_NAME}.${ARCHIVE_EXT}"
        ASSET_ID=$(echo "$RELEASE_JSON" | awk -v pattern="$ASSET_NAME_PATTERN" '
            /"id":/ { 
                id = $0
                gsub(/[^0-9]/, "", id)
            }
            /"name":/ { 
                if (index($0, pattern) > 0 && id != "") {
                    print id
                    exit
                }
            }
        ')
        
        if [ -z "$ASSET_ID" ]; then
            print_error "Failed to find asset ID for ${ASSET_NAME_PATTERN}"
            exit 1
        fi
        
        print_info "Found asset ID: $ASSET_ID"
        DOWNLOAD_URL="https://api.github.com/repos/${REPO_OWNER}/${REPO_NAME}/releases/assets/${ASSET_ID}"
        print_info "Downloading from API: $DOWNLOAD_URL"
        
        # Download using API endpoint with Accept header for binary content
        if command -v curl >/dev/null 2>&1; then
            curl -fsSL -H "Authorization: Bearer $GITHUB_TOKEN" \
                 -H "Accept: application/octet-stream" \
                 -o "archive.${ARCHIVE_EXT}" "$DOWNLOAD_URL"
        elif command -v wget >/dev/null 2>&1; then
            wget --header="Authorization: Bearer $GITHUB_TOKEN" \
                 --header="Accept: application/octet-stream" \
                 -q -O "archive.${ARCHIVE_EXT}" "$DOWNLOAD_URL"
        fi
    else
        # No token - use direct download URL (works for public repos)
        DOWNLOAD_URL="https://github.com/${REPO_OWNER}/${REPO_NAME}/releases/download/v${VERSION}/${ARCHIVE_NAME}.${ARCHIVE_EXT}"
        print_info "Downloading from: $DOWNLOAD_URL"
        
        if command -v curl >/dev/null 2>&1; then
            curl -fsSL -o "archive.${ARCHIVE_EXT}" "$DOWNLOAD_URL"
        elif command -v wget >/dev/null 2>&1; then
            wget -q -O "archive.${ARCHIVE_EXT}" "$DOWNLOAD_URL"
        else
            print_error "Neither curl nor wget found. Please install one of them."
            exit 1
        fi
    fi

    if [ ! -f "archive.${ARCHIVE_EXT}" ]; then
        print_error "Failed to download binary"
        exit 1
    fi

    print_info "Extracting archive..."

    # Extract archive
    if [ "$ARCHIVE_EXT" = "zip" ]; then
        if command -v unzip >/dev/null 2>&1; then
            unzip -q "archive.${ARCHIVE_EXT}"
        else
            print_error "unzip not found. Please install unzip."
            exit 1
        fi
    else
        tar -xzf "archive.${ARCHIVE_EXT}"
    fi

    # Find binary
    if [ "$OS" = "windows" ]; then
        BINARY_PATH="${BINARY_NAME}.exe"
    else
        BINARY_PATH="${BINARY_NAME}"
    fi

    if [ ! -f "$BINARY_PATH" ]; then
        print_error "Binary not found in archive"
        exit 1
    fi

    # Make binary executable
    chmod +x "$BINARY_PATH"

    # Create install directory if it doesn't exist
    if [ ! -d "$INSTALL_DIR" ]; then
        print_info "Creating directory: $INSTALL_DIR"
        mkdir -p "$INSTALL_DIR"
    fi

    # Install binary
    print_info "Installing binary to $INSTALL_DIR/$BINARY_NAME..."
    cp "$BINARY_PATH" "$INSTALL_DIR/$BINARY_NAME"

    # Verify installation
    if [ -x "$INSTALL_DIR/$BINARY_NAME" ]; then
        print_success "Installation complete!"
        echo ""
        
        # Test binary
        if "$INSTALL_DIR/$BINARY_NAME" version >/dev/null 2>&1; then
            INSTALLED_VERSION=$("$INSTALL_DIR/$BINARY_NAME" version | head -n1)
            print_success "Installed: $INSTALLED_VERSION"
        fi

        # Check PATH
        if ! is_in_path "$INSTALL_DIR"; then
            echo ""
            print_warning "$INSTALL_DIR is not in your PATH"
            echo ""
            echo "Add it to your PATH by adding this line to your shell config file:"
            echo ""
            echo "  ${GREEN}export PATH=\"$INSTALL_DIR:\$PATH\"${NC}"
            echo ""
            echo "For bash: ~/.bashrc or ~/.bash_profile"
            echo "For zsh:  ~/.zshrc"
            echo ""
            echo "Then reload your shell or run: source ~/.bashrc (or ~/.zshrc)"
        else
            echo ""
            print_success "Binary is in your PATH and ready to use!"
        fi

        echo ""
        echo "Next steps:"
        echo "  1. Run: ${CYAN}dotfiles bootstrap${NC}"
        echo "  2. Follow the prompts to set up your dotfiles"
        echo ""
    else
        print_error "Installation failed"
        exit 1
    fi
}

# Main
main() {
    echo ""
    echo "${CYAN}╔═══════════════════════════════════════════╗${NC}"
    echo "${CYAN}║     Dotfiles Installer                    ║${NC}"
    echo "${CYAN}║     github.com/${REPO_OWNER}/${REPO_NAME}     ║${NC}"
    echo "${CYAN}╚═══════════════════════════════════════════╝${NC}"
    echo ""

    install_binary
}

main "$@"
