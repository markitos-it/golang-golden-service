#!/bin/bash

# Centralized environment configuration

DEFAULT_DATABASE_HOST="localhost"
DEFAULT_DATABASE_USER="markitos_it_svc_golden"
DEFAULT_DATABASE_PASSWORD="markitos_it_svc_golden"
DEFAULT_DATABASE_NAME="markitos_it_svc_golden"
DEFAULT_DATABASE_SSL_MODE="disable"
DEFAULT_DATABASE_DSN="host=${DEFAULT_DATABASE_HOST} user=${DEFAULT_DATABASE_USER} password=${DEFAULT_DATABASE_PASSWORD} dbname=${DEFAULT_DATABASE_NAME} sslmode=${DEFAULT_DATABASE_SSL_MODE}"
DEFAULT_GRPC_SERVER_ADDRESS=":30000"
DEFAULT_POSTGRES_CONTAINER_NAME="markitos-it-svc-golden-postgres"
DEFAULT_UPLOADS_BASEDIR="./uploads"

function setup_environment() {
    : ${DATABASE_DSN:="${DEFAULT_DATABASE_DSN}"}
    : ${GRPC_SERVER_ADDRESS:="${DEFAULT_GRPC_SERVER_ADDRESS}"}
    : ${POSTGRES_CONTAINER_NAME:="${DEFAULT_POSTGRES_CONTAINER_NAME}"}
    : ${GOLDEN_UPLOADS_BASEDIR:="${DEFAULT_UPLOADS_BASEDIR}"}

    export DATABASE_DSN
    export GRPC_SERVER_ADDRESS
    export POSTGRES_CONTAINER_NAME
    export GOLDEN_UPLOADS_BASEDIR
}

function show_config() {
    echo "Starting configuration:"
    echo "DATABASE_DSN=$DATABASE_DSN"
    
    if [[ "${1:-}" == "full" ]]; then
        echo "GRPC_SERVER_ADDRESS=$GRPC_SERVER_ADDRESS"
        echo "POSTGRES_CONTAINER_NAME=$POSTGRES_CONTAINER_NAME"
        echo "GOLDEN_UPLOADS_BASEDIR=$GOLDEN_UPLOADS_BASEDIR"
    fi
    
    echo "-------------------------------------"
}

function show_banner() {
    echo "============================================="
    echo " __  __  ____  _  __"
    echo "|  \/  |  _ \| |/ /"
    echo "| \  / | | | | ' / "
    echo "| |\/| | | | |  <  "
    echo "| |  | | |_| | . \ "
    echo "|_|  |_|____/|_|\\_\\"
    echo ""
    echo "Creator: Marco Antonio - markitos"
    echo "============================================="
    echo " > (mArKit0sDevSecOpsKit)"
    echo " > Markitos DevSecOps Kulture"
    echo ""
}