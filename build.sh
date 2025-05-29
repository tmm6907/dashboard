#!/bin/bash

# --- Functions ---
log_info() {
  echo "INFO: $1"
}

log_error() {
  echo "ERROR: $1" >&2
}

# --- Main Script ---

# Pre-flight Check: Ensure Docker is installed
if ! command -v docker &> /dev/null; then
    log_error "Docker is not installed. Please install Docker before running this script."
    exit 1
fi

log_info "Starting Docker image build process for local use..."

# Build server image
SERVER_IMAGE="mashboard-server:local"
log_info "Building $SERVER_IMAGE..."
docker build -t "$SERVER_IMAGE" -f src/api/Dockerfile src/api || { log_error "Failed to build $SERVER_IMAGE."; exit 1; }
log_info "$SERVER_IMAGE built successfully."

# Build client image
CLIENT_IMAGE="mashboard-client:local"
log_info "Building $CLIENT_IMAGE..."
docker build -t "$CLIENT_IMAGE" -f src/client/Dockerfile src/client || { log_error "Failed to build $CLIENT_IMAGE."; exit 1; }
log_info "$CLIENT_IMAGE built successfully."

# Build mail image
MAIL_IMAGE="mashboard-mail:local"
log_info "Building $MAIL_IMAGE..."
docker build -t "$MAIL_IMAGE" -f src/mailserver/Dockerfile src/mailserver || { log_error "Failed to build $MAIL_IMAGE."; exit 1; }
log_info "$MAIL_IMAGE built successfully."

# Build caddy image (context is the current directory, Dockerfile is at root)
CADDY_IMAGE="mashboard-caddy:local"
log_info "Building $CADDY_IMAGE..."
docker build -t "$CADDY_IMAGE" -f Dockerfile . || { log_error "Failed to build $CADDY_IMAGE."; exit 1; }
log_info "$CADDY_IMAGE built successfully."

log_info "All Docker images built and tagged for local use!"

---