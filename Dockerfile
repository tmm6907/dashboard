ARG GO_VERSION=1.24
FROM --platform=$BUILDPLATFORM golang:${GO_VERSION} AS build

# Install xcaddy
RUN go install github.com/caddyserver/xcaddy/cmd/xcaddy@latest

# Set working directory to /build
WORKDIR /build

# Build Caddy with the L4 plugin
RUN xcaddy build --with github.com/mholt/caddy-l4

# --- Final Stage ---
FROM alpine:latest AS final

RUN apk add nss-tools

# Copy the Caddy binary from the build stage
COPY --from=build /build/caddy /usr/bin/caddy


# Set a volume for persistent data storage
VOLUME ["/data"]

# Ensure permissions for the /data directory
RUN mkdir -p /data && chmod -R 755 /data

# Expose Caddy ports
EXPOSE 80 443 25

# Default entry point for Caddy
ENTRYPOINT ["/usr/bin/caddy"]

# Default command (can be overridden)
CMD ["run", "--config", "/etc/caddy/Caddyfile", "--adapter", "caddyfile"]
