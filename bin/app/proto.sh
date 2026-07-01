#!/usr/bin/env bash

set -euo pipefail
IFS=$'\n\t'
echo "🔧 Compiling proto files from $(pwd)"
ENVIRONMENT_FILE="bin/shared/environment.sh"
source $ENVIRONMENT_FILE

function log_info() {
    echo "[INFO] $*"
}
function log_error() {
    echo "[ERROR] $*" >&2
}

setup_environment
show_config "full"

#:[.'.]:>-------------------------------------
show_banner

export PATH="$PATH:$(go env GOPATH)/bin"
protoc \
    --proto_path=internal/infrastructure/proto \
    --go_out=internal/infrastructure/gapi \
    --go-grpc_out=internal/infrastructure/gapi \
    --go_opt=paths=source_relative \
    --go-grpc_opt=paths=source_relative \
    internal/infrastructure/proto/*.proto
#:[.'.]:>-------------------------------------

log_info "Proto files compiled successfully."