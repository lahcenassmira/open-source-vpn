# Multi-stage build for open-source-vpn Server
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git make

# Set working directory
WORKDIR /build

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build server binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o vpn-server ./cmd/server

# Final stage
FROM alpine:latest

# Install runtime dependencies
RUN apk add --no-cache \
    iptables \
    ip6tables \
    iproute2 \
    ca-certificates

# Create non-root user (but VPN needs root for TUN device)
RUN addgroup -g 1000 vpn && \
    adduser -D -u 1000 -G vpn vpn

# Copy binary from builder
COPY --from=builder /build/vpn-server /usr/local/bin/vpn-server

# Create directories
RUN mkdir -p /etc/vpn /var/log

# Set permissions
RUN chmod +x /usr/local/bin/vpn-server

# Expose VPN port
EXPOSE 51820/udp

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD pgrep vpn-server || exit 1

# Volume for configuration
VOLUME ["/etc/vpn"]

# Run as root (required for TUN device)
USER root

# Entry point
ENTRYPOINT ["/usr/local/bin/vpn-server"]
CMD ["start", "--config", "/etc/vpn/server.yaml"]
