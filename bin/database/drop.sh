#!/usr/bin/env bash

set -euo pipefail
IFS=$'\n\t'
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

DB_NAME=$(echo "$DATABASE_DSN" | awk -F'[ =]' '{for(i=1;i<=NF;i++){if($i=="dbname"){print $(i+1)}}}')
DB_USER=$(echo "$DATABASE_DSN" | awk -F'[ =]' '{for(i=1;i<=NF;i++){if($i=="user"){print $(i+1)}}}')

log_info "Removing database and associated user"

docker exec $POSTGRES_CONTAINER_NAME psql -U markitos_it_svc_golden -d postgres -c "SELECT pg_terminate_backend(pid) FROM pg_stat_activity WHERE datname='$DB_NAME';" || true
docker exec $POSTGRES_CONTAINER_NAME psql -U markitos_it_svc_golden -d postgres -c "DROP DATABASE IF EXISTS \"$DB_NAME\";" || log_error "Failed to drop database $DB_NAME"
log_info "Database $DB_NAME dropped (if existed)"

docker exec $POSTGRES_CONTAINER_NAME psql -U markitos_it_svc_golden -d postgres -c "DROP USER IF EXISTS \"$DB_USER\";" || log_error "Failed to drop user $DB_USER"
log_info "User $DB_USER dropped (if existed)"

log_info "Removal process completed"
#:[.'.]:>-------------------------------------