# open-source-vpn Quick Start Guide

Get up and running with open-source-vpn in 5 minutes!

## Prerequisites

- Linux or macOS
- Go 1.21+
- Root/sudo access

## Installation

```bash
# Clone the repository
git clone https://github.com/lahcenassmira/open-source-vpn.git
cd open-source-vpn

# Download dependencies
go mod download

# Build binaries
make build
```

## Server Setup (5 steps)

### 1. Generate Server Keys

```bash
sudo ./bin/vpn-server keygen --output server-keys.json
```

Copy the private key from the output.

### 2. Create Server Configuration

```bash
cp configs/server.example.yaml server.yaml
```

Edit `server.yaml` and paste your private key:

```yaml
crypto:
  private_key: "YOUR_PRIVATE_KEY_HERE"
```

### 3. Enable IP Forwarding

```bash
sudo sysctl -w net.ipv4.ip_forward=1
```

### 4. Configure Firewall

```bash
# Allow VPN port
sudo ufw allow 51820/udp

# Setup NAT (replace eth0 with your interface)
sudo iptables -t nat -A POSTROUTING -s 10.8.0.0/24 -o eth0 -j MASQUERADE
```

### 5. Start Server

```bash
sudo ./bin/vpn-server start --config server.yaml
```

## Client Setup (4 steps)

### 1. Generate Client Keys

```bash
./bin/vpn-client keygen --output client-keys.json
```

Copy both the private key and public key from the output.

### 2. Add Client to Server

On the server, edit `server.yaml` and add:

```yaml
clients:
  - public_key: "CLIENT_PUBLIC_KEY_HERE"
    allowed_ips:
      - "10.8.0.2/32"
    name: "my-laptop"
```

Restart the server.

### 3. Create Client Configuration

```bash
cp configs/client.example.yaml client.yaml
```

Edit `client.yaml`:

```yaml
client:
  server_address: "YOUR_SERVER_IP:51820"

crypto:
  private_key: "YOUR_CLIENT_PRIVATE_KEY"
  server_public_key: "SERVER_PUBLIC_KEY"
```

### 4. Connect

```bash
sudo ./bin/vpn-client connect --config client.yaml
```

## Verify Connection

```bash
# Check interface
ip addr show tun0

# Ping server
ping 10.8.0.1

# Check routing
ip route
```

## Common Issues

### "Permission denied"
Run with `sudo` - VPN needs root for TUN device.

### "Failed to create TUN interface"
```bash
sudo modprobe tun
```

### "Connection timeout"
Check firewall allows UDP port 51820:
```bash
sudo ufw status
```

### "No route to host"
Verify server IP is correct and reachable:
```bash
ping YOUR_SERVER_IP
```

## Next Steps

- [Full Setup Guide](docs/SETUP.md) - Detailed installation
- [Architecture](docs/ARCHITECTURE.md) - How it works
- [Configuration](configs/) - Advanced options

## Quick Commands Reference

```bash
# Server
sudo ./bin/vpn-server keygen --output keys.json
sudo ./bin/vpn-server start --config server.yaml
sudo ./bin/vpn-server status

# Client
./bin/vpn-client keygen --output keys.json
sudo ./bin/vpn-client connect --config client.yaml
sudo ./bin/vpn-client status

# Build
make build          # Build binaries
make test           # Run tests
make clean          # Clean build artifacts

# Docker
make docker-build   # Build Docker image
docker-compose up   # Run server in Docker
```

## Security Reminder

⚠️ **Never share your private keys!**
- Keep private keys secure (0600 permissions)
- Use separate keys for each client
- Don't commit keys to version control
- Rotate keys periodically

## Getting Help

- 📖 [Documentation](docs/)
- 🐛 [Report Issues](https://github.com/lahcenassmira/open-source-vpn/issues)
- 💬 [Discussions](https://github.com/lahcenassmira/open-source-vpn/discussions)

---

**Ready to go!** 🚀 Your VPN is now running securely.
