# open-source-vpn Setup Guide

## Prerequisites

### System Requirements

- **Operating System**: Linux (Ubuntu 20.04+, Debian 11+, CentOS 8+) or macOS
- **Go**: Version 1.21 or higher
- **Root Access**: Required for TUN device management
- **Network**: UDP port 51820 (or custom port) open in firewall

### Dependencies

```bash
# Ubuntu/Debian
sudo apt-get update
sudo apt-get install -y build-essential git iproute2 iptables

# CentOS/RHEL
sudo yum install -y gcc git iproute iptables

# macOS
brew install go
```

## Installation

### Option 1: Build from Source

```bash
# Clone the repository
git clone https://github.com/lahcenassmira/open-source-vpn.git
cd open-source-vpn

# Download dependencies
go mod download

# Build binaries
make build

# Binaries will be in bin/
ls -la bin/
```

### Option 2: Install to System

```bash
# Build and install to $GOPATH/bin
make install

# Verify installation
which vpn-server
which vpn-client
```

### Option 3: Docker

```bash
# Build Docker image
make docker-build

# Or use docker-compose
cd docker
docker-compose up -d
```

## Server Setup

### Step 1: Generate Server Keys

```bash
sudo ./bin/vpn-server keygen --output server-keys.json
```

Output:
```
Key pair generated successfully!
Private key: <base64-encoded-private-key>
Public key:  <base64-encoded-public-key>

Keys saved to: server-keys.json

⚠️  Keep your private key secure and never share it!
```

### Step 2: Configure Server

```bash
# Copy example configuration
cp configs/server.example.yaml configs/server.yaml

# Edit configuration
nano configs/server.yaml
```

Update the following fields:

```yaml
server:
  listen_address: "0.0.0.0:51820"  # Your server IP and port

network:
  interface: "tun0"
  address: "10.8.0.1/24"  # VPN network range
  mtu: 1420

crypto:
  private_key: "<paste-server-private-key-here>"

clients: []  # Will add clients later
```

### Step 3: Configure Firewall

```bash
# Allow VPN port
sudo ufw allow 51820/udp

# Enable IP forwarding
sudo sysctl -w net.ipv4.ip_forward=1
sudo sysctl -w net.ipv6.conf.all.forwarding=1

# Make permanent
echo "net.ipv4.ip_forward=1" | sudo tee -a /etc/sysctl.conf
echo "net.ipv6.conf.all.forwarding=1" | sudo tee -a /etc/sysctl.conf

# Configure NAT (if routing all traffic)
sudo iptables -t nat -A POSTROUTING -s 10.8.0.0/24 -o eth0 -j MASQUERADE
sudo iptables -A FORWARD -i tun0 -j ACCEPT
sudo iptables -A FORWARD -o tun0 -j ACCEPT

# Save iptables rules
sudo iptables-save | sudo tee /etc/iptables/rules.v4
```

### Step 4: Start Server

```bash
# Start server (foreground)
sudo ./bin/vpn-server start --config configs/server.yaml

# Or run as systemd service (see below)
```

### Step 5: Create Systemd Service (Optional)

```bash
# Create service file
sudo nano /etc/systemd/system/vpn-server.service
```

Content:
```ini
[Unit]
Description=open-source-vpn Server
After=network.target

[Service]
Type=simple
User=root
ExecStart=/usr/local/bin/vpn-server start --config /etc/vpn/server.yaml
Restart=on-failure
RestartSec=5s

[Install]
WantedBy=multi-user.target
```

```bash
# Enable and start service
sudo systemctl daemon-reload
sudo systemctl enable vpn-server
sudo systemctl start vpn-server

# Check status
sudo systemctl status vpn-server

# View logs
sudo journalctl -u vpn-server -f
```

## Client Setup

### Step 1: Generate Client Keys

```bash
./bin/vpn-client keygen --output client-keys.json
```

Output:
```
Key pair generated successfully!
Private key: <base64-encoded-private-key>
Public key:  <base64-encoded-public-key>

Keys saved to: client-keys.json

⚠️  Keep your private key secure and never share it!

📋 Add this public key to the server configuration:
   public_key: "<client-public-key>"
```

### Step 2: Add Client to Server

Edit `configs/server.yaml` on the server:

```yaml
clients:
  - public_key: "<client-public-key>"
    allowed_ips:
      - "10.8.0.2/32"
    name: "client1"
```

Restart the server:
```bash
sudo systemctl restart vpn-server
```

### Step 3: Configure Client

```bash
# Copy example configuration
cp configs/client.example.yaml configs/client.yaml

# Edit configuration
nano configs/client.yaml
```

Update the following fields:

```yaml
client:
  server_address: "vpn.example.com:51820"  # Your server address

network:
  interface: "tun0"
  address: "10.8.0.2/24"  # Must match server's allowed_ips
  mtu: 1420

crypto:
  private_key: "<paste-client-private-key-here>"
  server_public_key: "<paste-server-public-key-here>"

routing:
  default_route: false  # Set true to route all traffic through VPN
  routes:
    - "10.8.0.0/24"  # VPN network
    # - "0.0.0.0/0"  # Uncomment to route all traffic
```

### Step 4: Connect Client

```bash
# Connect to VPN (foreground)
sudo ./bin/vpn-client connect --config configs/client.yaml

# Press Ctrl+C to disconnect
```

### Step 5: Verify Connection

```bash
# Check TUN interface
ip addr show tun0

# Ping VPN server
ping 10.8.0.1

# Check routing
ip route show

# Test connectivity
curl ifconfig.me  # Should show your real IP
# If default_route: true, should show VPN server IP
```

## Docker Deployment

### Server Deployment

```bash
# Create configuration directory
mkdir -p /opt/vpn/configs

# Copy configuration
cp configs/server.yaml /opt/vpn/configs/

# Run with docker-compose
cd docker
docker-compose up -d

# View logs
docker-compose logs -f

# Stop server
docker-compose down
```

### Manual Docker Run

```bash
docker run -d \
  --name vpn-server \
  --cap-add=NET_ADMIN \
  --device=/dev/net/tun \
  -p 51820:51820/udp \
  -v /opt/vpn/configs:/etc/vpn \
  open-source-vpn-server:latest
```

## Troubleshooting

### Server Issues

**Problem**: Server fails to start with "permission denied"
```bash
# Solution: Run with sudo
sudo ./bin/vpn-server start --config configs/server.yaml
```

**Problem**: "failed to create TUN interface"
```bash
# Solution: Load TUN module
sudo modprobe tun

# Make permanent
echo "tun" | sudo tee -a /etc/modules
```

**Problem**: Clients can't connect
```bash
# Check firewall
sudo ufw status
sudo iptables -L -n

# Check server is listening
sudo netstat -ulnp | grep 51820

# Check logs
sudo journalctl -u vpn-server -n 50
```

### Client Issues

**Problem**: "no route to host"
```bash
# Check server is reachable
ping vpn.example.com
nc -zvu vpn.example.com 51820

# Check DNS resolution
nslookup vpn.example.com
```

**Problem**: Connected but no internet
```bash
# Check routing
ip route show

# Check DNS
cat /etc/resolv.conf

# Test connectivity
ping 10.8.0.1  # VPN server
ping 8.8.8.8   # External IP
```

**Problem**: "failed to decrypt packet"
```bash
# Solution: Keys mismatch
# Verify client public key is in server config
# Verify server public key is in client config
# Regenerate keys if necessary
```

### Performance Issues

**Problem**: High latency
```bash
# Check MTU settings
ip link show tun0

# Reduce MTU if needed
# In config: mtu: 1280

# Check server load
top
htop
```

**Problem**: Low throughput
```bash
# Check CPU usage
mpstat -P ALL 1

# Check network interface
ethtool eth0

# Increase workers in server config
# server:
#   workers: 8
```

## Security Best Practices

1. **Key Management**
   - Store private keys with 0600 permissions
   - Never commit keys to version control
   - Rotate keys periodically
   - Use separate keys per client

2. **Network Security**
   - Use firewall rules to restrict access
   - Enable fail2ban for brute force protection
   - Monitor logs for suspicious activity
   - Use strong passwords for server access

3. **System Hardening**
   - Keep system updated
   - Disable unnecessary services
   - Use SSH key authentication
   - Enable automatic security updates

4. **Monitoring**
   - Set up log aggregation
   - Monitor connection patterns
   - Track bandwidth usage
   - Alert on anomalies

## Next Steps

- [Architecture Documentation](ARCHITECTURE.md)
- [API Reference](API.md)
- [Contributing Guide](../CONTRIBUTING.md)
- [FAQ](FAQ.md)
