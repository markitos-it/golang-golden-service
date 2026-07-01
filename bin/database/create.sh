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

DB_NAME=$(echo $DATABASE_DSN | awk -F'[ =]' '{for(i=1;i<=NF;i++){if($i=="dbname"){print $(i+1)}}}')
DB_USER=$(echo $DATABASE_DSN | awk -F'[ =]' '{for(i=1;i<=NF;i++){if($i=="user"){print $(i+1)}}}')
DB_PASS=$(echo $DATABASE_DSN | awk -F'[ =]' '{for(i=1;i<=NF;i++){if($i=="password"){print $(i+1)}}}')

function database_exists() {
    docker exec $POSTGRES_CONTAINER_NAME psql -U markitos_it_svc_golden -d postgres -tAc "SELECT 1 FROM pg_database WHERE datname='$DB_NAME'" | grep -q 1
}

function user_exists() {
    docker exec $POSTGRES_CONTAINER_NAME psql -U markitos_it_svc_golden -d postgres -tAc "SELECT 1 FROM pg_roles WHERE rolname='$DB_USER'" | grep -q 1
}

log_info "Creating database $DB_NAME"

if database_exists "$DB_NAME"; then
    log_info "Database $DB_NAME already exists"
else
    docker exec $POSTGRES_CONTAINER_NAME createdb --username=markitos_it_svc_golden --owner=markitos_it_svc_golden "$DB_NAME"
    log_info "Database $DB_NAME created in container $POSTGRES_CONTAINER_NAME"
fi

if user_exists "$DB_USER"; then
    log_info "User $DB_USER already exists"
else
    docker exec $POSTGRES_CONTAINER_NAME psql -U markitos_it_svc_golden -d postgres -c "CREATE USER \"$DB_USER\" WITH PASSWORD '$DB_PASS';"
    log_info "User $DB_USER created"
fi

docker exec $POSTGRES_CONTAINER_NAME psql -U markitos_it_svc_golden -d postgres -c "GRANT ALL PRIVILEGES ON DATABASE \"$DB_NAME\" TO \"$DB_USER\";"
docker exec $POSTGRES_CONTAINER_NAME psql -U markitos_it_svc_golden -d postgres -c "GRANT ALL PRIVILEGES ON SCHEMA public TO \"$DB_USER\";"
log_info "Granted privileges to user $DB_USER on database $DB_NAME"
#:[.'.]:>-------------------------------------