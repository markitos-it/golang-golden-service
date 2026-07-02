#!/usr/bin/env bash

set -euo pipefail
IFS=$'\n\t'
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ENVIRONMENT_FILE="$SCRIPT_DIR/../shared/environment.sh"
source $ENVIRONMENT_FILE

readonly RED='\033[0;31m'
readonly GREEN='\033[0;32m'
readonly YELLOW='\033[0;33m'
readonly BLUE='\033[0;34m'
readonly CYAN='\033[0;36m'
readonly BOLD='\033[1m'
readonly NC='\033[0m'

log_info() { echo -e "${BLUE}[INFO]${NC} $*"; }
log_error() { echo -e "${RED}[ERROR]${NC} $*" >&2; }

declare -a report_lines
overall_status=0

run_check() {
    local name="$1"
    local command="$2"
    local output_file
    output_file=$(mktemp)

    printf "  %-30s" "🛡️  $name"

    if eval $command &> "$output_file"; then
        echo -e "[ ${GREEN}OK${NC} ]"
        report_lines+=("✅ $name: OK")
    else
        local exit_code=$?
        if [[ "$name" == Snyk\ IaC* ]]; then
            if [[ $exit_code -eq 1 ]]; then
                echo -e "[ ${RED}KO${NC} ]"
                report_lines+=("❌ $name: KO (Exit code: $exit_code)")
                overall_status=1
                echo -e "${RED}----------------- DETALLES DEL ERROR -----------------${NC}"
                cat "$output_file"
                echo -e "${RED}----------------------------------------------------${NC}"
            else
                echo -e "[ ${YELLOW}WARN${NC} ] (Ignored exit code $exit_code)"
                report_lines+=("🟡 $name: WARNING (Ignored exit code $exit_code)")
            fi
        else
            echo -e "[ ${RED}KO${NC} ]"
            report_lines+=("❌ $name: KO (Exit code: $exit_code)")
            overall_status=1
            echo -e "${RED}----------------- DETALLES DEL ERROR -----------------${NC}"
            cat "$output_file"
            echo -e "${RED}----------------------------------------------------${NC}"
        fi
    fi
    rm -f "$output_file"
}

setup_environment

echo -e "${CYAN}${BOLD}"
echo "╔════════════════════════════════════════════════════════════════╗"
echo "║        🛡️  INICIANDO ANÁLISIS DE SEGURIDAD (LIGERO) 🛡️           ║"
echo "╚════════════════════════════════════════════════════════════════╝"
echo -e "${NC}"

snyk auth $SNYK_TOKEN
log_info "Versión Snyk: $(snyk --version)"
log_info "Versión Gitleaks: $(gitleaks version)"
echo

run_check "Snyk Code (SAST)" "snyk code test --severity-threshold=medium --include-ignores"
run_check "Gitleaks (Secrets)" "gitleaks detect --source . --no-git --verbose"

echo -e "\n${CYAN}${BOLD}╔════════════════════════════════════════════════════════════════╗"
echo "║                     📊 RESUMEN DEL ANÁLISIS 📊                 ║"
echo -e "╚════════════════════════════════════════════════════════════════╝${NC}"
for line in "${report_lines[@]}"; do
    echo -e "  $line"
done
echo

if [ $overall_status -ne 0 ]; then
    log_error "El análisis de seguridad ha fallado."
    exit 1
else
    log_info "🎉 ¡Análisis de seguridad completado con éxito! 🎉"
fi
