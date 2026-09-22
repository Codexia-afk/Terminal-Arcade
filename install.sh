#!/usr/bin/env bash
# Terminal Arcade - Universal One-Line Installer
# Usage: curl -fsSL https://raw.githubusercontent.com/Codexia-afk/Terminal-Arcade/main/install.sh | bash

set -euo pipefail

REPO="Codexia-afk/Terminal-Arcade"
BINARY_NAME="arcade"

# Styling & Colors
BOLD="$(tput bold 2>/dev/null || printf '')"
GREEN="$(tput setaf 2 2>/dev/null || printf '')"
CYAN="$(tput setaf 6 2>/dev/null || printf '')"
YELLOW="$(tput setaf 3 2>/dev/null || printf '')"
RED="$(tput setaf 1 2>/dev/null || printf '')"
RESET="$(tput sgr0 2>/dev/null || printf '')"

info() {
    printf "${CYAN}==>${RESET} ${BOLD}%s${RESET}\n" "$1"
}

success() {
    printf "${GREEN}✔${RESET} ${BOLD}%s${RESET}\n" "$1"
}

warn() {
    printf "${YELLOW}⚠${RESET} %s\n" "$1"
}

error() {
    printf "${RED}✘ Error:${RESET} %s\n" "$1" >&2
    exit 1
}

# ASCII Logo
print_banner() {
    cat << "EOF"
  _______                  _             _                         _      
 |__   __|                (_)           | |      /\               | |     
    | | ___ _ __ _ __ ___  _ _ __   __ _| |     /  \   _ __ ___ __| | ___ 
    | |/ _ \ '__| '_ ` _ \| | '_ \ / _` | |    / /\ \ | '__/ __/ _` |/ _ \
    | |  __/ |  | | | | | | | | | | (_| | |   / ____ \| | | (_| (_| |  __/
    |_|\___|_|  |_| |_| |_|_|_| |_|\__,_|_|  /_/    \_\_|  \___\__,_|\___|
EOF
    printf "   %sRetro Terminal Gaming Suite (Snake, Pacman, Breakout)%s\n\n" "${CYAN}" "${RESET}"
}

print_banner

# 1. Detect Operating System
OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
case "$OS" in
    darwin)
        OS="darwin"
        ;;
    linux)
        OS="linux"
        ;;
    mingw*|msys*|cygwin*)
        OS="windows"
        ;;
    *)
        error "Unsupported operating system: $OS. Supported: macOS, Linux, Windows."
        ;;
esac

# 2. Detect Hardware Architecture
ARCH="$(uname -m)"
case "$ARCH" in
    x86_64|amd64)
        ARCH="amd64"
        ;;
    arm64|aarch64)
        ARCH="arm64"
        ;;
    *)
        error "Unsupported architecture: $ARCH. Supported: x86_64, arm64."
        ;;
esac

info "Detected platform: ${OS}/${ARCH}"

# 3. Determine Installation Destination
INSTALL_DIR=""
if [ -d "/usr/local/bin" ] && [ -w "/usr/local/bin" ]; then
    INSTALL_DIR="/usr/local/bin"
elif mkdir -p "$HOME/.local/bin" 2>/dev/null && [ -d "$HOME/.local/bin" ] && [ -w "$HOME/.local/bin" ]; then
    INSTALL_DIR="$HOME/.local/bin"
elif command -v go >/dev/null 2>&1 && [ -d "$(go env GOPATH)/bin" ]; then
    INSTALL_DIR="$(go env GOPATH)/bin"
elif [ -d "$HOME/bin" ] && [ -w "$HOME/bin" ]; then
    INSTALL_DIR="$HOME/bin"
else
    INSTALL_DIR="/usr/local/bin"
fi

TMP_DIR="$(mktemp -d 2>/dev/null || mktemp -d -t 'arcade-install')"
trap 'rm -rf "$TMP_DIR"' EXIT

TARGET_BINARY="${TMP_DIR}/${BINARY_NAME}"
DOWNLOAD_FILE="arcade_${OS}_${ARCH}"
if [ "$OS" = "windows" ]; then
    DOWNLOAD_FILE="${DOWNLOAD_FILE}.exe"
    TARGET_BINARY="${TARGET_BINARY}.exe"
fi

# Ensure DEVELOPER_DIR is set on macOS if CommandLineTools is available
if [ "$OS" = "darwin" ] && [ -d "/Library/Developer/CommandLineTools" ] && [ -z "${DEVELOPER_DIR:-}" ]; then
    export DEVELOPER_DIR="/Library/Developer/CommandLineTools"
fi

INSTALLED_VIA=""

# 4. Attempt Installation:
# Check 1: If installer is run from inside the local repository, build directly
if [ -f "./cmd/arcade/main.go" ] && command -v go >/dev/null 2>&1; then
    info "Repository checkout detected. Compiling binary..."
    go build -buildvcs=false -o "$TARGET_BINARY" ./cmd/arcade
    INSTALLED_VIA="local source build"
fi

# Check 2: If local pre-built ./arcade exists in current directory
if [ ! -f "$TARGET_BINARY" ] && [ -f "./arcade" ] && [ -x "./arcade" ]; then
    info "Found pre-built arcade binary in current directory..."
    cp "./arcade" "$TARGET_BINARY"
    INSTALLED_VIA="local binary"
fi

# Check 3: Go toolchain via canonical package
if [ ! -f "$TARGET_BINARY" ] && command -v go >/dev/null 2>&1; then
    info "Go toolchain detected. Installing canonical package..."
    if go install -buildvcs=false "github.com/${REPO}/cmd/arcade@latest" 2>/dev/null; then
        GOPATH_BIN="$(go env GOPATH)/bin"
        if [ -f "${GOPATH_BIN}/${BINARY_NAME}" ]; then
            cp "${GOPATH_BIN}/${BINARY_NAME}" "$TARGET_BINARY"
            INSTALLED_VIA="go install"
        fi
    fi
fi

if [ ! -f "$TARGET_BINARY" ]; then
    # Strategy A: GitHub Releases API for GoReleaser archive
    RELEASE_JSON="$(curl -sSL "https://api.github.com/repos/${REPO}/releases/latest" 2>/dev/null || true)"
    LATEST_TAG="$(echo "$RELEASE_JSON" | grep '"tag_name":' | head -n 1 | sed -E 's/.*"([^"]+)".*/\1/')"

    if [ -n "$LATEST_TAG" ]; then
        VERSION_NUM="${LATEST_TAG#v}"
        ARCHIVE_NAME="arcade_${VERSION_NUM}_${OS}_${ARCH}.tar.gz"
        if [ "$OS" = "windows" ]; then
            ARCHIVE_NAME="arcade_${VERSION_NUM}_${OS}_${ARCH}.zip"
        fi
        ARCHIVE_URL="https://github.com/${REPO}/releases/download/${LATEST_TAG}/${ARCHIVE_NAME}"

        info "Downloading release asset ${ARCHIVE_NAME} (${LATEST_TAG})..."
        if curl -fsSL --connect-timeout 15 -o "${TMP_DIR}/${ARCHIVE_NAME}" "$ARCHIVE_URL" 2>/dev/null; then
            if [ "$OS" = "windows" ]; then
                unzip -q "${TMP_DIR}/${ARCHIVE_NAME}" -d "${TMP_DIR}" 2>/dev/null || true
            else
                tar -xzf "${TMP_DIR}/${ARCHIVE_NAME}" -C "${TMP_DIR}" 2>/dev/null || true
            fi
            FOUND="$(find "${TMP_DIR}" -type f \( -name "arcade" -o -name "arcade.exe" \) | head -n 1)"
            if [ -n "$FOUND" ] && [ -f "$FOUND" ]; then
                cp "$FOUND" "$TARGET_BINARY"
                INSTALLED_VIA="github release archive (${LATEST_TAG})"
            fi
        fi
    fi

    # Strategy B: Fallback to standalone direct binary asset or raw repository tree
    if [ ! -f "$TARGET_BINARY" ]; then
        RELEASE_URL="https://github.com/${REPO}/releases/latest/download/${DOWNLOAD_FILE}"
        RAW_FALLBACK_URL="https://raw.githubusercontent.com/${REPO}/main/dist/${DOWNLOAD_FILE}"

        info "Downloading pre-compiled binary: ${DOWNLOAD_FILE}..."

        if curl -fsSL --connect-timeout 10 -o "$TARGET_BINARY" "$RELEASE_URL" 2>/dev/null; then
            INSTALLED_VIA="github releases"
        elif curl -fsSL --connect-timeout 10 -o "$TARGET_BINARY" "$RAW_FALLBACK_URL" 2>/dev/null; then
            INSTALLED_VIA="github source tree"
        fi
    fi

    # Strategy C: Build from source if git and go are present
    if [ ! -f "$TARGET_BINARY" ]; then
        if command -v git >/dev/null 2>&1 && command -v go >/dev/null 2>&1; then
            info "Fetching repository and building from source..."
            git clone --depth 1 "https://github.com/${REPO}.git" "${TMP_DIR}/src" >/dev/null 2>&1
            (cd "${TMP_DIR}/src" && go build -buildvcs=false -o "$TARGET_BINARY" ./cmd/arcade)
            INSTALLED_VIA="source build"
        else
            error "Could not download pre-compiled binary and Go is not installed to compile from source.
Please download manually from: https://github.com/${REPO}/releases"
        fi
    fi
fi

chmod +x "$TARGET_BINARY"

# 5. Place binary into destination directory
DEST_PATH="${INSTALL_DIR}/${BINARY_NAME}"
info "Installing to ${DEST_PATH}..."

mkdir -p "$INSTALL_DIR" 2>/dev/null || true
if mv "$TARGET_BINARY" "$DEST_PATH" 2>/dev/null; then
    :
elif command -v sudo >/dev/null 2>&1; then
    info "Escalating privileges (sudo) to write to ${INSTALL_DIR}..."
    sudo mv "$TARGET_BINARY" "$DEST_PATH"
else
    error "Cannot write to ${DEST_PATH}. Please run the installer with sudo or check directory permissions."
fi

# 6. Verify PATH accessibility
PATH_CONFIGURED=1
if ! echo "$PATH" | tr ':' '\n' | grep -qx "$INSTALL_DIR"; then
    PATH_CONFIGURED=0
    SHELL_PROFILE=""
    CURRENT_SHELL="$(basename "${SHELL:-bash}")"
    
    if [ "$CURRENT_SHELL" = "zsh" ]; then
        SHELL_PROFILE="$HOME/.zshrc"
    elif [ "$CURRENT_SHELL" = "bash" ]; then
        if [ -f "$HOME/.bash_profile" ]; then
            SHELL_PROFILE="$HOME/.bash_profile"
        else
            SHELL_PROFILE="$HOME/.bashrc"
        fi
    else
        SHELL_PROFILE="$HOME/.profile"
    fi

    warn "${INSTALL_DIR} is not in your current PATH."
    printf "Adding %s to %s...\n" "${INSTALL_DIR}" "${SHELL_PROFILE}"
    printf '\n# Terminal Arcade PATH\nexport PATH="%s:$PATH"\n' "${INSTALL_DIR}" >> "$SHELL_PROFILE"
    export PATH="${INSTALL_DIR}:$PATH"
fi

# 7. Verification Test
if ! "$DEST_PATH" -h >/dev/null 2>&1; then
    warn "Binary installed at ${DEST_PATH}, but execution test encountered a warning."
else
    success "Verification check passed! (${INSTALLED_VIA})"
fi

# 8. Success Output
echo ""
echo "============================================================"
printf "${GREEN}${BOLD}🎉 Terminal Arcade successfully installed!${RESET}\n"
echo "============================================================"
echo ""
printf "👉 Just open any terminal and type:  ${CYAN}${BOLD}arcade${RESET}\n\n"
if [ "$PATH_CONFIGURED" -eq 0 ]; then
    printf "Note: Since your PATH was updated, run:\n"
    printf "      ${YELLOW}source %s${RESET} (or open a new terminal window)\n\n" "$SHELL_PROFILE"
fi
printf "Enjoy Snake 🐍, Pacman 🟡, and Breakout 🧱!\n"
