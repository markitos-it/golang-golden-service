#!/usr/bin/env bash

##############################################################################
#                       🥷 MARKITOS APPSEC HOOKSTOOLS UNINSTALLER 🧹         #
#                                                                            #
#           Elimina hook files instaladas en .git/hooks/pre-commit           #
#                                                                            #
#                     Markitos DevSecOps Kulture - 2025                      #
##############################################################################
set -euo pipefail
IFS=$'\n\t'

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR/../../" || true

readonly RED='\033[0;31m'
readonly GREEN='\033[0;32m'
readonly YELLOW='\033[0;33m'
readonly BLUE='\033[0;34m'
readonly CYAN='\033[0;36m'
readonly BOLD='\033[1m'
readonly NC='\033[0m'

log_info() { echo -e "${BLUE}[INFO]${NC} $*"; }
log_success() { echo -e "${GREEN}[OK]${NC} $*"; }
log_warning() { echo -e "${YELLOW}[WARN]${NC} $*"; }
log_error() { echo -e "${RED}[ERR]${NC} $*" >&2; }

echo -e "${CYAN}${BOLD}Markitos AppSec Uninstaller (hooks)${NC}"

PRECOMMIT_PATH=".git/hooks/pre-commit"
if [[ -f "$PRECOMMIT_PATH" ]]; then
    rm -f "$PRECOMMIT_PATH" && echo "Removed $PRECOMMIT_PATH" || echo "Failed to remove $PRECOMMIT_PATH"
else
    echo "$PRECOMMIT_PATH not present, skipping"
fi

echo "Uninstall (minimal) complete."