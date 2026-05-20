# 🇲🇦 SecureVPN - Open Source VPN for Morocco

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)](https://golang.org/)
[![Made in Morocco](https://img.shields.io/badge/Made%20in-Morocco%20🇲🇦-red)](https://github.com/lahcenassmira/open-source-vpn)
[![Security](https://img.shields.io/badge/Security-ChaCha20--Poly1305-green)](https://github.com/lahcenassmira/open-source-vpn)
[![Legal Use Only](https://img.shields.io/badge/⚠️-Legal%20Use%20Only-red)](https://github.com/lahcenassmira/open-source-vpn#%EF%B8%8F-security--legal-warning)


A secure, high-performance VPN system written in Go with modern cryptography and clean architecture.

**🎯 Perfect for Moroccan TPE/PME (Small & Medium Businesses)**

---

## 🇲🇦 Pourquoi ce VPN? / Why This VPN?

### Pour les Entreprises Marocaines / For Moroccan Businesses

Ce VPN open-source est **idéal pour**:
- 🏢 **TPE/PME** - Petites et moyennes entreprises
- 🔐 **Sécurité** - Protection des données sensibles
- 💰 **Économique** - Solution gratuite et open-source
- 🌍 **Télétravail** - Accès sécurisé pour employés distants
- 🇲🇦 **Local** - Développé avec les besoins marocains en tête

This open-source VPN is **perfect for**:
- 🏢 **Small & Medium Businesses** - TPE/PME
- 🔐 **Security** - Protect sensitive business data
- 💰 **Cost-Effective** - Free and open-source solution
- 🌍 **Remote Work** - Secure access for remote employees
- 🇲🇦 **Local** - Built with Moroccan needs in mind

### 🎓 Use Cases / Cas d'Usage

 - ✅ **Personal VPN Usage** - Protection personnelle
 - ✅ **Corporate Deployments** - Déploiements d'entreprise
 - ✅ **Educational Purposes** - Objectifs éducatifs
 - ✅ **Learning Go & Cryptography** - Apprendre Go et la cryptographie
 - ✅ **Building Custom VPN Solutions** - Solutions VPN personnalisées
 - ✅ **Open-Source Contributions** - Contributions open-source

---

## ⚠️ SECURITY & LEGAL WARNING

> **🚨 IMPORTANT: READ BEFORE USE 🚨**

### 🛡️ This Project is Designed For:

✅ **Legitimate Use Only:**
- 🔐 Secure business communications
- 🏢 Private corporate networks
- 🎓 Educational and research purposes
- 🔒 Protecting sensitive data on untrusted networks
- 🌐 Secure remote access to private resources
- 💼 Professional VPN deployments

### 🚫 STRICTLY PROHIBITED:

❌ **DO NOT USE FOR:**
- Illegal activities of any kind
- Bypassing legal restrictions or regulations
- Unauthorized access to networks or systems
- Circumventing content protection or copyright laws
- Malicious traffic routing or attacks
- Any activity that violates local, national, or international laws

### ⚖️ Legal Compliance:

**YOU ARE RESPONSIBLE FOR:**
- ✅ Compliance with all applicable laws in Morocco and your jurisdiction
- ✅ Obtaining necessary permissions and authorizations
- ✅ Respecting intellectual property rights
- ✅ Following your organization's security policies
- ✅ Ensuring legitimate use at all times

**🇲🇦 Morocco Specific:**
- VPN usage is legal in Morocco for legitimate purposes
- Users must comply with Moroccan telecommunications regulations
- Businesses should consult with legal counsel for compliance

### 🔒 Security Responsibilities:

**AS A USER, YOU MUST:**
- 🔑 Keep private keys secure and never share them
- 🔐 Use strong passwords and secure key storage
- 📊 Monitor logs for suspicious activity
- 🔄 Keep software updated with security patches
- 🛡️ Follow security best practices
- 📋 Implement proper access controls

### ⚠️ Disclaimer:

**THE AUTHORS AND CONTRIBUTORS:**
- Are NOT responsible for misuse of this software
- Do NOT endorse or support illegal activities
- Provide this software "AS IS" without warranties
- Are NOT liable for any damages or legal consequences
- Strongly condemn any misuse of this technology

**BY USING THIS SOFTWARE, YOU AGREE:**
- To use it only for lawful purposes
- To comply with all applicable laws and regulations
- To take full responsibility for your use
- That the authors are not liable for your actions

### 📞 Report Abuse:

If you discover misuse of this software, please report it to:
- GitHub Issues: [Report Here](https://github.com/lahcenassmira/open-source-vpn/issues)
- Local authorities if illegal activity is suspected

---

**🔴 REMEMBER: With great power comes great responsibility. Use this tool ethically and legally. 🔴**

---

## 🏗️ Architecture

```
┌─────────────┐                    ┌─────────────┐
│   Client    │◄──── Encrypted ───►│   Server    │
│  (TUN Dev)  │      Tunnel        │  (TUN Dev)  │
└─────────────┘                    └─────────────┘
      │                                    │
      ▼                                    ▼
  Local Apps                         Internet Gateway
```

### Components

- **Server**: VPN server handling multiple concurrent clients
- **Client**: VPN client establishing secure tunnels
- **Crypto**: X25519 key exchange + ChaCha20-Poly1305 encryption
- **Tunnel**: TUN device management and packet handling
- **Network**: Routing, NAT traversal, DNS forwarding

## 🔐 Security Features

- **Mutual Authentication**: Public/private key pairs (X25519)
- **Modern Encryption**: ChaCha20-Poly1305 AEAD cipher
- **Key Exchange**: Noise Protocol Framework (Noise_XX pattern)
- **Replay Protection**: Nonce-based anti-replay mechanism
- **Perfect Forward Secrecy**: Ephemeral key exchange per session

## 🚀 Quick Start

### Prerequisites

- Go 1.21 or higher
- Linux/macOS (TUN device support)
- Root/sudo privileges (for network interface management)

### Installation

```bash
# Clone the repository
git clone https://github.com/lahcenassmira/open-source-vpn.git
cd open-source-vpn

# Build server and client
make build

# Or build manually
go build -o bin/vpn-server ./cmd/server
go build -o bin/vpn-client ./cmd/client
```

### Generate Keys

```bash
# Generate server keys
./bin/vpn-server keygen --output server-keys.json

# Generate client keys
./bin/vpn-client keygen --output client-keys.json
```

### Start Server

```bash
# Edit server configuration
cp configs/server.example.yaml configs/server.yaml
# Add client public keys to server.yaml

# Start server (requires root)
sudo ./bin/vpn-server start --config configs/server.yaml
```

### Connect Client

```bash
# Edit client configuration
cp configs/client.example.yaml configs/client.yaml
# Add server public key and endpoint to client.yaml

# Connect (requires root)
sudo ./bin/vpn-client connect --config configs/client.yaml
```

## 📁 Project Structure

```
open-source-vpn/
├── cmd/
│   ├── server/          # Server CLI entry point
│   └── client/          # Client CLI entry point
├── internal/
│   ├── crypto/          # Cryptography and key exchange
│   ├── tunnel/          # TUN device management
│   ├── network/         # Routing and packet handling
│   ├── protocol/        # VPN protocol implementation
│   └── config/          # Configuration management
├── pkg/
│   ├── logger/          # Structured logging
│   └── metrics/         # Performance metrics
├── configs/             # Example configurations
├── docker/              # Docker deployment files
├── scripts/             # Utility scripts
└── docs/                # Additional documentation
```

## ⚙️ Configuration

### Server Configuration (YAML)

```yaml
server:
  listen_address: "0.0.0.0:51820"
  protocol: "udp"
  
network:
  interface: "tun0"
  address: "10.8.0.1/24"
  mtu: 1420
  
crypto:
  private_key: "server_private_key_base64"
  
clients:
  - public_key: "client1_public_key_base64"
    allowed_ips: ["10.8.0.2/32"]
  - public_key: "client2_public_key_base64"
    allowed_ips: ["10.8.0.3/32"]
    
routing:
  forward_all_traffic: false
  allowed_networks: ["10.8.0.0/24"]
  
dns:
  enabled: true
  servers: ["8.8.8.8", "8.8.4.4"]
  
logging:
  level: "info"
  format: "json"
  output: "/var/log/vpn-server.log"
```

### Client Configuration (YAML)

```yaml
client:
  server_address: "vpn.example.com:51820"
  protocol: "udp"
  
network:
  interface: "tun0"
  address: "10.8.0.2/24"
  mtu: 1420
  
crypto:
  private_key: "client_private_key_base64"
  server_public_key: "server_public_key_base64"
  
routing:
  default_route: false
  routes: ["10.8.0.0/24"]
  
dns:
  enabled: true
  servers: ["10.8.0.1"]
  
logging:
  level: "info"
  format: "json"
```

## 🐳 Docker Deployment

```bash
# Build Docker image
docker build -t open-source-vpn-server -f docker/Dockerfile .

# Run server container
docker run -d \
  --name vpn-server \
  --cap-add=NET_ADMIN \
  --device=/dev/net/tun \
  -p 51820:51820/udp \
  -v $(pwd)/configs:/etc/vpn \
  open-source-vpn-server
```

## 📊 Monitoring & Logging

### View Connection Status

```bash
# Server status
./bin/vpn-server status

# Client status
./bin/vpn-client status
```

### Logs

Structured JSON logs include:
- Connection events (connect/disconnect)
- Data transfer statistics
- Latency measurements
- Error tracking

Example log entry:
```json
{
  "timestamp": "2026-05-20T10:30:45Z",
  "level": "info",
  "component": "server",
  "event": "client_connected",
  "client_id": "abc123",
  "client_ip": "10.8.0.2",
  "remote_addr": "203.0.113.45:54321"
}
```

## 🔧 CLI Commands

### Server Commands

```bash
vpn-server start --config <path>     # Start VPN server
vpn-server stop                      # Stop VPN server
vpn-server status                    # Show server status
vpn-server keygen --output <path>    # Generate key pair
vpn-server logs --follow             # Tail logs
```

### Client Commands

```bash
vpn-client connect --config <path>   # Connect to VPN
vpn-client disconnect                # Disconnect from VPN
vpn-client status                    # Show connection status
vpn-client keygen --output <path>    # Generate key pair
vpn-client logs --follow             # Tail logs
```

## 🎯 Performance

- **Latency**: < 5ms overhead on local network
- **Throughput**: 1+ Gbps on modern hardware
- **Concurrent Clients**: 1000+ per server
- **Memory**: ~50MB base + ~1MB per client
- **CPU**: Efficient goroutine-based I/O

## 🧪 Testing

```bash
# Run all tests
make test

# Run with coverage
make test-coverage

# Run benchmarks
make benchmark

# Integration tests
make test-integration
```

## 🛡️ Security Considerations

1. **Key Management**: Store private keys securely, never commit to version control
2. **Firewall**: Configure firewall rules to allow only VPN traffic
3. **Updates**: Keep dependencies updated for security patches
4. **Monitoring**: Enable logging and monitor for suspicious activity
5. **Access Control**: Use allowed_ips to restrict client access

## 📚 Protocol Details

### Handshake (Noise_XX Pattern)

1. Client → Server: Ephemeral public key
2. Server → Client: Ephemeral public key + static public key (encrypted)
3. Client → Server: Static public key (encrypted) + payload

### Packet Format

```
┌──────────┬────────┬───────────┬─────────────┐
│ Type (1) │ ID (4) │ Nonce (12)│ Payload (N) │
└──────────┴────────┴───────────┴─────────────┘
```

- **Type**: Packet type (handshake/data/keepalive)
- **ID**: Connection identifier
- **Nonce**: Unique nonce for AEAD
- **Payload**: Encrypted data + authentication tag

## 🤝 Contributing

Contributions are welcome! Please read [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see [LICENSE](LICENSE) file for details.

## ⚠️ Disclaimer

This VPN system is designed for legitimate privacy and security use cases. Users are responsible for compliance with local laws and regulations. The authors are not responsible for misuse of this software.

## 🙏 Acknowledgments

- Inspired by WireGuard's design principles
- Uses Noise Protocol Framework for key exchange
- Built with Go's excellent networking libraries

## 📞 Support

- Documentation: [docs/](docs/)
- Issues: [GitHub Issues](https://github.com/lahcenassmira/open-source-vpn/issues)
- Discussions: [GitHub Discussions](https://github.com/lahcenassmira/open-source-vpn/discussions)

---

**Built with passion for secure and open internet.**
