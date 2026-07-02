#!/bin/bash

##############################################################################
#                       MARKITOS APPSEC TOOLS INSTALLER                       #
#                                                                            #
#               Installs security tools into ~/.local/bin                    #
#                                                                            #
#                        Markitos DevSecOps Kulture                          #
##############################################################################

# Bash strict mode
set -euo pipefail
IFS=$'\n\t'

# Change to script directory and then to project root
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR/../../" || exit 1

# Colors
readonly RED='\033[0;31m'
readonly GREEN='\033[0;32m'
readonly YELLOW='\033[0;33m'
readonly BLUE='\033[0;34m'
readonly CYAN='\033[0;36m'
readonly BOLD='\033[1m'
readonly NC='\033[0m'

echo -e "${CYAN}${BOLD}"
echo "╔════════════════════════════════════════════════════════════════╗"
echo "║               MARKITOS APPSEC INSTALLER                        ║"
echo "║                                                                ║"
echo "║  Installing: Snyk • Gitleaks                                   ║"
echo "║  Target: ~/.local/bin                                          ║"
echo "╚════════════════════════════════════════════════════════════════╝"
echo -e "${NC}"

# Logging functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $*"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $*"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $*" >&2
}

log_warning() {
    echo -e "${YELLOW}[WARN]${NC} $*"
}

TMPDIR="/tmp"
cleanup() {
    rm -f "$TMPDIR/snyk" "$TMPDIR/gitleaks" "$TMPDIR/gitleaks_*.tar.gz" "$TMPDIR/README.md" "$TMPDIR/LICENSE" || true
}
trap cleanup EXIT

mkdir -p ~/.local/bin

log_info "Installing Snyk CLI..."
SNYK_TMP="$TMPDIR/snyk"
if curl -sL https://static.snyk.io/cli/latest/snyk-linux -o "$SNYK_TMP"; then
    chmod +x "$SNYK_TMP"
    mv "$SNYK_TMP" ~/.local/bin/snyk
    log_success "Snyk CLI installed successfully"
else
    log_error "Failed to download Snyk CLI"
    exit 1
fi

log_info "Installing Gitleaks..."
GITLEAKS_VERSION="8.24.0"
GITLEAKS_URL="https://github.com/gitleaks/gitleaks/releases/download/v${GITLEAKS_VERSION}/gitleaks_${GITLEAKS_VERSION}_linux_x64.tar.gz"
GITLEAKS_TMP="$TMPDIR/gitleaks.tar.gz"

if wget -q "$GITLEAKS_URL" -O "$GITLEAKS_TMP"; then
    tar xfz "$GITLEAKS_TMP" -C "$TMPDIR"
    chmod +x "$TMPDIR/gitleaks" || true
    mv "$TMPDIR/gitleaks" ~/.local/bin/gitleaks
    log_success "Gitleaks installed successfully"
else
    log_error "Failed to download Gitleaks"
    exit 1
fi

log_info "Updating PATH in ~/.bashrc..."
LOCAL_BIN_PATH='export PATH=$PATH:$HOME/.local/bin'

if grep -Fq '.local/bin' ~/.bashrc 2>/dev/null; then
    log_info "PATH already includes ~/.local/bin"
else
    echo >> ~/.bashrc
    echo "# Added by Markitos AppSec Installer" >> ~/.bashrc
    echo "$LOCAL_BIN_PATH" >> ~/.bashrc
    log_success "PATH updated in ~/.bashrc"
fi

export PATH="$PATH:$HOME/.local/bin"

PRECOMMIT_SRC="etc/pre-commit"
if [[ -d ".git" && -d ".git/hooks" ]]; then
    PRECOMMIT_DEST=".git/hooks/pre-commit"
    if [[ -f "$PRECOMMIT_SRC" ]]; then
        if [[ -f "$PRECOMMIT_DEST" ]]; then
            cp "$PRECOMMIT_DEST" "$PRECOMMIT_DEST.bak" || true
            log_info "Backed up existing pre-commit hook to $PRECOMMIT_DEST.bak"
        fi
        cp "$PRECOMMIT_SRC" "$PRECOMMIT_DEST"
        chmod +x "$PRECOMMIT_DEST"
        log_success "Installed pre-commit hook to $PRECOMMIT_DEST"
    else
        log_warning "$PRECOMMIT_SRC not found; skipping pre-commit hook installation"
    fi
elif [[ -d ".git" ]]; then
    log_warning ".git/hooks directory not present; skipping pre-commit hook installation (hooks might be managed externally)"
else
    log_warning ".git directory not found in working directory; skipping pre-commit hook installation"
fi

log_info "Verifying installations..."
echo

if ~/.local/bin/snyk --version &>/dev/null; then
    SNYK_VERSION=$(~/.local/bin/snyk --version)
    log_success "Snyk CLI: $SNYK_VERSION"
else
    log_error "Snyk CLI verification failed"
fi

if ~/.local/bin/gitleaks version &>/dev/null; then
    GITLEAKS_VERSION=$(~/.local/bin/gitleaks version)
    log_success "Gitleaks: $GITLEAKS_VERSION"
else
    log_error "Gitleaks verification failed"
fi

echo
echo -e "${GREEN}${BOLD}"
echo "╔════════════════════════════════════════════════════════════════╗"
echo "║                    🎉 INSTALLATION COMPLETE 🎉                 ║"
echo "║                                                                ║"
echo "║  Tools installed in: ~/.local/bin                              ║"
echo "║                                                                ║"
echo "║  Next steps:                                                   ║"
echo "║  1. Restart terminal or run: source ~/.bashrc                  ║"
echo "║  2. Authenticate Snyk: snyk auth                               ║"
echo "║  3. Install pre-commit hook                                    ║"
echo "║                                                                ║"
echo "║  'Security is a feature, not a bug fix.' - Markitos            ║"
echo "╚════════════════════════════════════════════════════════════════╝"
echo -e "${NC}"

# Show authentication reminder
echo
log_warning "Don't forget to authenticate Snyk:"
echo -e "${YELLOW}  snyk auth${NC}"
echo