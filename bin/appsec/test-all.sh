#!/usr/bin/env bash

set -euo pipefail
IFS=$'\n\t'
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ENVIRONMENT_FILE="$SCRIPT_DIR/../shared/environment.sh"
source $ENVIRONMENT_FILE

log_info() { echo "[INFO] $*"; }
log_error() { echo "[ERROR] $*" >&2; }

setup_environment
show_config "full"

#:[.'.]:>-------------------------------------
snyk auth $SNYK_TOKEN
SNYK_VER=$(snyk --version 2>/dev/null || true)
if [[ -n "$SNYK_VER" ]]; then
    echo "Snyk CLI version: $SNYK_VER"
else
    log_error "Snyk is not installed. Please check your installation."
fi
GITLEAKS_VER=$(gitleaks version 2>/dev/null || true)
if [[ -n "$GITLEAKS_VER" ]]; then
    echo "Gitleaks version: $GITLEAKS_VER"
else
    log_error "Gitleaks is not installed. Please check your installation."
fi
snyk code test --severity-threshold=medium --include-ignores || log_error "Snyk Code encontró problemas."
snyk test --all-projects --severity-threshold=medium --include-ignores || log_error "Snyk Test found issues."
snyk iac test --severity-threshold=high || log_error "Snyk IAC encontró problemas."
gitleaks detect --source . --no-git --verbose || log_error "Gitleaks found issues."
log_info "Security analysis completed."
#:[.'.]:>-------------------------------------
#:[.'.]:> Tu lógica aquí
#:[.'.]:>-------------------------------------
