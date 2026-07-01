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
echo "║               MARKITOS APPSEC HOOKS INSTALLER                  ║"
echo "║                                                                ║"
echo "║  Installing: pre-commit                                        ║"
echo "║  Target: .git/hooks/pre-commit                                 ║"
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

mkdir -p .git/hooks
PRECOMMIT_SRC="etc/pre-commit"
PRECOMMIT_DEST=".git/hooks/pre-commit"
if [[ -d ".git" && -d ".git/hooks" ]]; then
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

echo
echo -e "${GREEN}${BOLD}"
echo "╔════════════════════════════════════════════════════════════════╗"
echo "║                    🎉 INSTALLATION COMPLETE 🎉                 ║"
echo "║                                                                ║"
echo "║  Hooks installed in: .git/hooks                                ║"
echo "║                                                                ║"
echo "║  Next steps:                                                   ║"
echo "║  1. Install pre-commit hook                                    ║"
echo "║                                                                ║"
echo "║  'Security is a feature, not a bug fix.' - Markitos            ║"
echo "╚════════════════════════════════════════════════════════════════╝"
echo -e "${NC}"
echo
echo