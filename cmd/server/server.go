package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/open-source-vpn/vpn/internal/config"
	"github.com/open-source-vpn/vpn/internal/crypto"
	"github.com/open-source-vpn/vpn/internal/network"
	"github.com/open-source-vpn/vpn/internal/protocol"
	"github.com/open-source-vpn/vpn/internal/tunnel"
	"github.com/open-source-vpn/vpn/pkg/logger"
	"github.com/open-source-vpn/vpn/pkg/metrics"
	"go.uber.org/zap"
)

// Server represents the VPN server
type Server struct {
	config      *config.ServerConfig
	log         *logger.Logger
	metrics     *metrics.Metrics
	
	tunInterface *tunnel.Interface
	udpConn      *net.UDPConn
	
	privateKey   [crypto.KeySize]byte
	clients      map[string]*ClientInfo
	connections  map[uint32]*protocol.Connection
	connMu       sync.RWMutex
	
	router       *network.Router
	packetPool   *network.PacketPool
	
	ctx          context.Context
	cancel       context.CancelFunc
	wg           sync.WaitGroup
}

// ClientInfo stores information about an allowed client
type ClientInfo struct {
	PublicKey  [crypto.KeySize]byte
	AllowedIPs []string
	Name       string
}

// NewServer creates a new VPN server
func NewServer(cfg *config.ServerConfig, log *logger.Logger) (*Server, error) {
	// Decode private key
	privateKey, err := crypto.DecodeKey(cfg.Crypto.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("invalid private key: %w", err)
	}

	// Parse allowed clients
	clients := make(map[string]*ClientInfo)
	for _, client := range cfg.Clients {
		pubKey, err := crypto.DecodeKey(client.PublicKey)
		if err != nil {
			log.Warn("skipping client with invalid public key", zap.String("name", client.Name))
			continue
		}
		
		clients[crypto.EncodeKey(pubKey[:])] = &ClientInfo{
			PublicKey:  pubKey,
			AllowedIPs: client.AllowedIPs,
			Name:       client.Name,
		}
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &Server{
		config:      cfg,
		log:         log,
		metrics:     metrics.NewMetrics(),
		privateKey:  privateKey,
		clients:     clients,
		connections: make(map[uint32]*protocol.Connection),
		router:      network.NewRouter(),
		packetPool:  network.NewPacketPool(256, cfg.Network.MTU+100),
		ctx:         ctx,
		cancel:      cancel,
	}, nil
}

// Start starts the VPN server
func (s *Server) Start() error {
	s.log.Info("initializing VPN server",
		zap.String("listen", s.config.Server.ListenAddress),
		zap.String("interface", s.config.Network.Interface),
	)

	// Create TUN interface
	tunIface, err := tunnel.NewTunInterface(s.config.Network.Interface, s.config.Network.MTU)
	if err != nil {
		return fmt.Errorf("failed to create TUN interface: %w", err)
	}
	s.tunInterface = tunIface
	s.log.Info("TUN interface created", zap.String("name", tunIface.Name()))

	// Configure TUN interface
	if err := tunnel.ConfigureInterface(tunIface.Name(), s.config.Network.Address, s.config.Network.MTU); err != nil {
		return fmt.Errorf("failed to configure interface: %w", err)
	}
	s.log.Info("TUN interface configured", zap.String("address", s.config.Network.Address))

	// Enable IP forwarding
	if err := tunnel.EnableIPForwarding(); err != nil {
		s.log.Warn("failed to enable IP forwarding", zap.Error(err))
	}

	// Create UDP socket
	addr, err := net.ResolveUDPAddr("udp", s.config.Server.ListenAddress)
	if err != nil {
		return fmt.Errorf("invalid listen address: %w", err)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen on UDP: %w", err)
	}
	s.udpConn = conn
	s.log.Info("UDP socket listening", zap.String("address", addr.String()))

	// Start worker goroutines
	s.wg.Add(3)
	go s.tunReader()
	go s.udpReader()
	go s.connectionManager()

	s.log.Info("VPN server started successfully")

	// Wait for shutdown signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	s.log.Info("shutting down VPN server")
	return s.Stop()
}

// Stop stops the VPN server
func (s *Server) Stop() error {
	s.cancel()
	
	if s.udpConn != nil {
		s.udpConn.Close()
	}
	
	if s.tunInterface != nil {
		s.tunInterface.Close()
	}
	
	s.wg.Wait()
	
	s.log.Info("VPN server stopped")
	return nil
}

// tunReader reads packets from TUN interface and sends them to clients
func (s *Server) tunReader() {
	defer s.wg.Done()
	
	s.log.Info("TUN reader started")
	buf := s.packetPool.Get()
	
	for {
		select {
		case <-s.ctx.Done():
			return
		default:
		}
		
		n, err := s.tunInterface.Read(buf)
		if err != nil {
			if s.ctx.Err() != nil {
				return
			}
			s.log.Error("failed to read from TUN", zap.Error(err))
			continue
		}
		
		// Parse IP packet
		ipPkt, err := network.ParseIPPacket(buf[:n])
		if err != nil {
			s.log.Debug("failed to parse IP packet", zap.Error(err))
			continue
		}
		
		// Find connection for destination IP
		conn := s.findConnectionByIP(ipPkt.DstIP)
		if conn == nil {
			s.log.Debug("no connection for destination", zap.String("dst", ipPkt.DstIP.String()))
			continue
		}
		
		// Encrypt and send packet
		encrypted, err := conn.EncryptPacket(network.PacketTypeData, buf[:n])
		if err != nil {
			s.log.Error("failed to encrypt packet", zap.Error(err))
			continue
		}
		
		if _, err := s.udpConn.WriteToUDP(encrypted, conn.RemoteAddr()); err != nil {
			s.log.Error("failed to send packet", zap.Error(err))
		}
	}
}

// udpReader reads packets from UDP socket and processes them
func (s *Server) udpReader() {
	defer s.wg.Done()
	
	s.log.Info("UDP reader started")
	buf := make([]byte, 65535)
	
	for {
		select {
		case <-s.ctx.Done():
			return
		default:
		}
		
		n, addr, err := s.udpConn.ReadFromUDP(buf)
		if err != nil {
			if s.ctx.Err() != nil {
				return
			}
			s.log.Error("failed to read from UDP", zap.Error(err))
			continue
		}
		
		// Parse packet
		var pkt network.Packet
		if err := pkt.Unmarshal(buf[:n]); err != nil {
			s.log.Debug("failed to parse packet", zap.Error(err))
			continue
		}
		
		// Handle packet based on type
		switch pkt.Header.Type {
		case network.PacketTypeHandshakeInit:
			s.handleHandshake(addr, &pkt)
		case network.PacketTypeData:
			s.handleDataPacket(addr, &pkt)
		case network.PacketTypeKeepalive:
			s.handleKeepalive(addr, &pkt)
		}
	}
}

// handleHandshake handles handshake packets
func (s *Server) handleHandshake(addr *net.UDPAddr, pkt *network.Packet) {
	s.log.Info("handshake received", zap.String("from", addr.String()))
	
	// TODO: Implement full Noise handshake
	// For now, create a simple connection
	
	connID := pkt.Header.ConnID
	if connID == 0 {
		connID = uint32(time.Now().UnixNano())
	}
	
	// Create temporary cipher (in production, derive from handshake)
	key := s.privateKey
	cipher, _ := crypto.NewCipher(&key)
	
	conn := protocol.NewConnection(connID, addr, cipher, cipher)
	
	s.connMu.Lock()
	s.connections[connID] = conn
	s.connMu.Unlock()
	
	s.metrics.ConnectionOpened()
	s.log.LogConnectionEvent("client_connected", connID, addr.String())
}

// handleDataPacket handles data packets
func (s *Server) handleDataPacket(addr *net.UDPAddr, pkt *network.Packet) {
	s.connMu.RLock()
	conn := s.connections[pkt.Header.ConnID]
	s.connMu.RUnlock()
	
	if conn == nil {
		s.log.Debug("connection not found", zap.Uint32("conn_id", pkt.Header.ConnID))
		return
	}
	
	// Decrypt packet
	decrypted, err := conn.DecryptPacket(pkt)
	if err != nil {
		s.log.Error("failed to decrypt packet", zap.Error(err))
		s.metrics.RecordDecryptionError()
		return
	}
	
	// Write to TUN interface
	if err := s.tunInterface.WritePacket(decrypted); err != nil {
		s.log.Error("failed to write to TUN", zap.Error(err))
	}
	
	s.metrics.RecordPacketReceived()
}

// handleKeepalive handles keepalive packets
func (s *Server) handleKeepalive(addr *net.UDPAddr, pkt *network.Packet) {
	s.connMu.RLock()
	conn := s.connections[pkt.Header.ConnID]
	s.connMu.RUnlock()
	
	if conn != nil {
		conn.UpdateLastSeen()
	}
}

// connectionManager manages connection lifecycle
func (s *Server) connectionManager() {
	defer s.wg.Done()
	
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	
	for {
		select {
		case <-s.ctx.Done():
			return
		case <-ticker.C:
			s.cleanupExpiredConnections()
		}
	}
}

// cleanupExpiredConnections removes expired connections
func (s *Server) cleanupExpiredConnections() {
	s.connMu.Lock()
	defer s.connMu.Unlock()
	
	timeout := 3 * time.Minute
	for id, conn := range s.connections {
		if conn.IsExpired(timeout) {
			conn.Close()
			delete(s.connections, id)
			s.metrics.ConnectionClosed()
			s.log.Info("connection expired", zap.Uint32("conn_id", id))
		}
	}
}

// findConnectionByIP finds a connection by destination IP
func (s *Server) findConnectionByIP(ip net.IP) *protocol.Connection {
	s.connMu.RLock()
	defer s.connMu.RUnlock()
	
	// TODO: Implement proper IP-to-connection mapping
	// For now, return first connection
	for _, conn := range s.connections {
		return conn
	}
	
	return nil
}
