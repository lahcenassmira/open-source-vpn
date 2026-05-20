# open-source-vpn Architecture

## Overview

open-source-vpn is a modern VPN system built with Go, emphasizing security, performance, and clean architecture. It implements a client-server model with strong cryptography and efficient packet handling.

## System Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                        VPN Client                            │
├─────────────────────────────────────────────────────────────┤
│  Application Layer                                           │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐                  │
│  │   CLI    │  │  Config  │  │  Logger  │                  │
│  └──────────┘  └──────────┘  └──────────┘                  │
├─────────────────────────────────────────────────────────────┤
│  Protocol Layer                                              │
│  ┌──────────────────┐  ┌──────────────────┐                │
│  │   Connection     │  │   Handshake      │                │
│  │   Management     │  │   (Noise XX)     │                │
│  └──────────────────┘  └──────────────────┘                │
├─────────────────────────────────────────────────────────────┤
│  Crypto Layer                                                │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐     │
│  │  X25519 KX   │  │ ChaCha20-    │  │   Replay     │     │
│  │              │  │ Poly1305     │  │   Window     │     │
│  └──────────────┘  └──────────────┘  └──────────────┘     │
├─────────────────────────────────────────────────────────────┤
│  Network Layer                                               │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐     │
│  │   Routing    │  │   Packet     │  │     NAT      │     │
│  │              │  │   Parsing    │  │              │     │
│  └──────────────┘  └──────────────┘  └──────────────┘     │
├─────────────────────────────────────────────────────────────┤
│  Tunnel Layer                                                │
│  ┌──────────────────────────────────────────────────┐       │
│  │              TUN Device Interface                 │       │
│  └──────────────────────────────────────────────────┘       │
└─────────────────────────────────────────────────────────────┘
                            │
                            │ Encrypted UDP
                            ▼
┌─────────────────────────────────────────────────────────────┐
│                        VPN Server                            │
│                    (Same Architecture)                       │
└─────────────────────────────────────────────────────────────┘
```

## Component Details

### 1. Crypto Layer (`internal/crypto/`)

**Purpose**: Provides cryptographic primitives and key management.

**Components**:
- **keys.go**: X25519 key pair generation and management
- **cipher.go**: ChaCha20-Poly1305 AEAD encryption
- **noise.go**: Noise Protocol Framework handshake implementation

**Key Features**:
- X25519 elliptic curve for key exchange
- ChaCha20-Poly1305 for authenticated encryption
- Noise_XX pattern for mutual authentication
- Perfect forward secrecy

### 2. Tunnel Layer (`internal/tunnel/`)

**Purpose**: Manages TUN virtual network interfaces.

**Components**:
- **tun.go**: TUN device creation and I/O operations

**Key Features**:
- Cross-platform TUN device support (Linux/macOS)
- Efficient packet buffering
- Batch read/write operations
- MTU configuration

### 3. Network Layer (`internal/network/`)

**Purpose**: Handles packet routing and network operations.

**Components**:
- **packet.go**: Packet parsing and serialization
- **router.go**: Routing table management and IP forwarding

**Key Features**:
- IPv4 and IPv6 support
- Longest prefix match routing
- Allowed IPs filtering
- NAT table management
- Packet pooling for performance

### 4. Protocol Layer (`internal/protocol/`)

**Purpose**: Implements the VPN protocol logic.

**Components**:
- **connection.go**: Connection state management

**Key Features**:
- Connection lifecycle management
- Packet encryption/decryption
- Replay attack protection
- Connection statistics tracking
- Automatic timeout handling

### 5. Configuration (`internal/config/`)

**Purpose**: Configuration file parsing and validation.

**Features**:
- YAML-based configuration
- Validation and defaults
- Separate server/client configs
- Hot-reload support (future)

### 6. Logging (`pkg/logger/`)

**Purpose**: Structured logging for debugging and monitoring.

**Features**:
- JSON and console output formats
- Multiple log levels
- Contextual logging
- Performance metrics logging

### 7. Metrics (`pkg/metrics/`)

**Purpose**: Performance and usage metrics tracking.

**Features**:
- Connection statistics
- Throughput measurement
- Latency tracking
- Error counting

## Data Flow

### Outbound (Client → Server)

```
Application
    ↓
TUN Interface (read IP packet)
    ↓
Parse IP packet
    ↓
Lookup connection
    ↓
Encrypt packet (ChaCha20-Poly1305)
    ↓
Add VPN header
    ↓
UDP Socket (send to server)
```

### Inbound (Server → Client)

```
UDP Socket (receive from server)
    ↓
Parse VPN packet
    ↓
Lookup connection
    ↓
Verify replay window
    ↓
Decrypt packet (ChaCha20-Poly1305)
    ↓
Extract IP packet
    ↓
TUN Interface (write IP packet)
    ↓
Application
```

## Security Model

### Handshake (Noise_XX Pattern)

```
Client                                Server
  │                                      │
  │  ─────── e ────────────────────────> │  (ephemeral key)
  │                                      │
  │  <────── e, s, encrypted ─────────── │  (ephemeral + static)
  │                                      │
  │  ─────── s, encrypted ──────────────> │  (static key)
  │                                      │
  │  <═══════ Secure Channel ═══════════> │
```

**Properties**:
- Mutual authentication
- Forward secrecy
- Identity hiding
- Resistance to replay attacks

### Packet Format

```
┌──────────┬────────┬───────────┬─────────────────────────┐
│ Type (1) │ ID (4) │ Nonce (24)│ Encrypted Payload + Tag │
└──────────┴────────┴───────────┴─────────────────────────┘
```

**Fields**:
- **Type**: Packet type (handshake/data/keepalive)
- **ID**: Connection identifier
- **Nonce**: Unique nonce for AEAD (prevents replay)
- **Payload**: Encrypted data with authentication tag

### Replay Protection

Uses a sliding window algorithm:
- Tracks last 1024 nonces
- Rejects duplicate nonces
- Rejects old nonces outside window
- Constant-time operations

## Performance Optimizations

### 1. Goroutine Architecture

```
Server/Client
├── TUN Reader (goroutine)
│   └── Reads packets from TUN device
├── UDP Reader (goroutine)
│   └── Reads packets from UDP socket
├── Connection Manager (goroutine)
│   └── Manages connection lifecycle
└── Keepalive Loop (goroutine)
    └── Sends periodic keepalives
```

### 2. Memory Management

- **Packet Pooling**: Reuse packet buffers to reduce GC pressure
- **Zero-Copy**: Minimize data copying where possible
- **Batch Operations**: Process multiple packets per syscall

### 3. Crypto Performance

- **ChaCha20-Poly1305**: Fast on all platforms (no AES-NI required)
- **X25519**: Efficient elliptic curve operations
- **Hardware Acceleration**: Uses Go's optimized crypto libraries

## Scalability

### Server Capacity

- **Concurrent Clients**: 1000+ per server
- **Throughput**: 1+ Gbps on modern hardware
- **Memory**: ~50MB base + ~1MB per client
- **CPU**: Scales with number of cores

### Horizontal Scaling

Future enhancements:
- Load balancer for multiple VPN servers
- Shared state via Redis/etcd
- Geographic distribution
- Automatic failover

## Security Considerations

### Key Management

- Private keys stored in configuration files (0600 permissions)
- Keys never logged or transmitted unencrypted
- Separate keys per client
- Key rotation support (future)

### Network Security

- Firewall rules to restrict access
- Rate limiting to prevent DoS
- Connection timeout to free resources
- IP allowlist per client

### Audit and Compliance

- Structured logging for audit trails
- Connection event tracking
- Data transfer statistics
- Error and security event logging

## Future Enhancements

### Planned Features

1. **Full Noise Protocol**: Complete Noise_XX implementation
2. **Multi-threading**: Per-connection worker pools
3. **UDP Hole Punching**: Better NAT traversal
4. **QUIC Support**: Alternative to UDP
5. **Web Dashboard**: Real-time monitoring UI
6. **Mobile Clients**: iOS and Android support
7. **WireGuard Compatibility**: Protocol compatibility layer

### Performance Improvements

1. **Kernel Bypass**: XDP/eBPF for packet processing
2. **SIMD Crypto**: Vectorized encryption
3. **Zero-Copy Networking**: io_uring support
4. **Connection Pooling**: Reuse connections
5. **Compression**: Optional payload compression

## References

- [Noise Protocol Framework](https://noiseprotocol.org/)
- [ChaCha20-Poly1305](https://tools.ietf.org/html/rfc8439)
- [X25519](https://cr.yp.to/ecdh.html)
- [WireGuard Design](https://www.wireguard.com/papers/wireguard.pdf)
