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

// Client represents the VPN client
type Client struct {
	config       *config.ClientConfig
	log          *logger.Logger
	metrics      *metrics.Metrics
	
	tunInterface *tunnel.Interface
	udpConn      *net.UDPConn
	serverAddr   *net.UDPAddr
	
	privateKey   [crypto.KeySize]byte
	serverPubKey [crypto.KeySize]byte
	
	connection   *protocol.Connection
	connMu       sync.RWMutex
	
	router       *network.Router
	packetPool   *network.PacketPool
	
	ctx          context.Context
	cancel       context.CancelFunc
	wg           sync.WaitGroup
}

// NewClient creates a new VPN client
func NewClient(cfg *config.ClientConfig, log *logger.Logger) (*Client, error) {
	// Decode private key
	privateKey, err := crypto.DecodeKey(cfg.Crypto.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("invalid private key: %w", err)
	}

	// Decode server public key
	serverPubKey, err := crypto.DecodeKey(cfg.Crypto.ServerPublicKey)
	if err != nil {
		return nil, fmt.Errorf("invalid server public key: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &Client{
		config:       cfg,
		log:          log,
		metrics:      metrics.NewMetrics(),
		privateKey:   privateKey,
		serverPubKey: serverPubKey,
		router:       network.NewRouter(),
		packetPool:   network.NewPacketPool(256, cfg.Network.MTU+100),
		ctx:          ctx,
		cancel:       cancel,
	}, nil
}

// Connect connects to the VPN server
func (c *Client) Connect() error {
	c.log.Info("connecting to VPN server",
		zap.String("server", c.config.Client.ServerAddress),
		zap.String("interface", c.config.Network.Interface),
	)

	// Create TUN interface
	tunIface, err := tunnel.NewTunInterface(c.config.Network.Interface, c.config.Network.MTU)
	if err != nil {
		return fmt.Errorf("failed to create TUN interface: %w", err)
	}
	c.tunInterface = tunIface
	c.log.Info("TUN interface created", zap.String("name", tunIface.Name()))

	// Configure TUN interface
	if err := tunnel.ConfigureInterface(tunIface.Name(), c.config.Network.Address, c.config.Network.MTU); err != nil {
		return fmt.Errorf("failed to configure interface: %w", err)
	}
	c.log.Info("TUN interface configured", zap.String("address", c.config.Network.Address))

	// Add routes
	for _, route := range c.config.Routing.Routes {
		if err := tunnel.AddRoute(route, tunIface.Name()); err != nil {
			c.log.Warn("failed to add route", zap.String("route", route), zap.Error(err))
		} else {
			c.log.Info("route added", zap.String("route", route))
		}
	}

	// Resolve server address
	serverAddr, err := net.ResolveUDPAddr("udp", c.config.Client.ServerAddress)
	if err != nil {
		return fmt.Errorf("invalid server address: %w", err)
	}
	c.serverAddr = serverAddr

	// Create UDP socket
	conn, err := net.DialUDP("udp", nil, serverAddr)
	if err != nil {
		return fmt.Errorf("failed to connect to server: %w", err)
	}
	c.udpConn = conn
	c.log.Info("UDP socket connected", zap.String("server", serverAddr.String()))

	// Perform handshake
	if err := c.performHandshake(); err != nil {
		return fmt.Errorf("handshake failed: %w", err)
	}

	// Start worker goroutines
	c.wg.Add(3)
	go c.tunReader()
	go c.udpReader()
	go c.keepaliveLoop()

	c.log.Info("VPN client connected successfully")

	// Wait for shutdown signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	c.log.Info("disconnecting from VPN server")
	return c.Disconnect()
}

// Disconnect disconnects from the VPN server
func (c *Client) Disconnect() error {
	c.cancel()
	
	if c.udpConn != nil {
		c.udpConn.Close()
	}
	
	if c.tunInterface != nil {
		c.tunInterface.Close()
	}
	
	c.wg.Wait()
	
	c.log.Info("VPN client disconnected")
	return nil
}

// performHandshake performs the initial handshake with the server
func (c *Client) performHandshake() error {
	c.log.Info("performing handshake")

	// TODO: Implement full Noise handshake
	// For now, create a simple connection
	
	connID := uint32(time.Now().UnixNano())
	
	// Create cipher (in production, derive from handshake)
	key := c.privateKey
	cipher, err := crypto.NewCipher(&key)
	if err != nil {
		return fmt.Errorf("failed to create cipher: %w", err)
	}
	
	conn := protocol.NewConnection(connID, c.serverAddr, cipher, cipher)
	
	c.connMu.Lock()
	c.connection = conn
	c.connMu.Unlock()
	
	// Send handshake packet
	pkt := &network.Packet{
		Header: network.PacketHeader{
			Type:   network.PacketTypeHandshakeInit,
			ConnID: connID,
		},
		Payload: []byte("HELLO"),
	}
	
	data, err := pkt.Marshal()
	if err != nil {
		return fmt.Errorf("failed to marshal handshake: %w", err)
	}
	
	if _, err := c.udpConn.Write(data); err != nil {
		return fmt.Errorf("failed to send handshake: %w", err)
	}
	
	c.metrics.ConnectionOpened()
	c.log.Info("handshake completed", zap.Uint32("conn_id", connID))
	
	return nil
}

// tunReader reads packets from TUN interface and sends them to server
func (c *Client) tunReader() {
	defer c.wg.Done()
	
	c.log.Info("TUN reader started")
	buf := c.packetPool.Get()
	
	for {
		select {
		case <-c.ctx.Done():
			return
		default:
		}
		
		n, err := c.tunInterface.Read(buf)
		if err != nil {
			if c.ctx.Err() != nil {
				return
			}
			c.log.Error("failed to read from TUN", zap.Error(err))
			continue
		}
		
		// Get connection
		c.connMu.RLock()
		conn := c.connection
		c.connMu.RUnlock()
		
		if conn == nil {
			c.log.Debug("no active connection")
			continue
		}
		
		// Encrypt and send packet
		encrypted, err := conn.EncryptPacket(network.PacketTypeData, buf[:n])
		if err != nil {
			c.log.Error("failed to encrypt packet", zap.Error(err))
			c.metrics.RecordEncryptionError()
			continue
		}
		
		if _, err := c.udpConn.Write(encrypted); err != nil {
			c.log.Error("failed to send packet", zap.Error(err))
		} else {
			c.metrics.RecordPacketSent()
		}
	}
}

// udpReader reads packets from UDP socket and processes them
func (c *Client) udpReader() {
	defer c.wg.Done()
	
	c.log.Info("UDP reader started")
	buf := make([]byte, 65535)
	
	for {
		select {
		case <-c.ctx.Done():
			return
		default:
		}
		
		n, err := c.udpConn.Read(buf)
		if err != nil {
			if c.ctx.Err() != nil {
				return
			}
			c.log.Error("failed to read from UDP", zap.Error(err))
			continue
		}
		
		// Parse packet
		var pkt network.Packet
		if err := pkt.Unmarshal(buf[:n]); err != nil {
			c.log.Debug("failed to parse packet", zap.Error(err))
			continue
		}
		
		// Handle packet based on type
		switch pkt.Header.Type {
		case network.PacketTypeHandshakeResp:
			c.handleHandshakeResponse(&pkt)
		case network.PacketTypeData:
			c.handleDataPacket(&pkt)
		case network.PacketTypeKeepalive:
			c.handleKeepalive(&pkt)
		}
	}
}

// handleHandshakeResponse handles handshake response packets
func (c *Client) handleHandshakeResponse(pkt *network.Packet) {
	c.log.Info("handshake response received")
}

// handleDataPacket handles data packets
func (c *Client) handleDataPacket(pkt *network.Packet) {
	c.connMu.RLock()
	conn := c.connection
	c.connMu.RUnlock()
	
	if conn == nil {
		c.log.Debug("no active connection")
		return
	}
	
	// Decrypt packet
	decrypted, err := conn.DecryptPacket(pkt)
	if err != nil {
		c.log.Error("failed to decrypt packet", zap.Error(err))
		c.metrics.RecordDecryptionError()
		return
	}
	
	// Write to TUN interface
	if err := c.tunInterface.WritePacket(decrypted); err != nil {
		c.log.Error("failed to write to TUN", zap.Error(err))
	} else {
		c.metrics.RecordPacketReceived()
	}
}

// handleKeepalive handles keepalive packets
func (c *Client) handleKeepalive(pkt *network.Packet) {
	c.log.Debug("keepalive received")
}

// keepaliveLoop sends periodic keepalive packets
func (c *Client) keepaliveLoop() {
	defer c.wg.Done()
	
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	
	for {
		select {
		case <-c.ctx.Done():
			return
		case <-ticker.C:
			c.sendKeepalive()
		}
	}
}

// sendKeepalive sends a keepalive packet
func (c *Client) sendKeepalive() {
	c.connMu.RLock()
	conn := c.connection
	c.connMu.RUnlock()
	
	if conn == nil {
		return
	}
	
	pkt := &network.Packet{
		Header: network.PacketHeader{
			Type:   network.PacketTypeKeepalive,
			ConnID: conn.ID(),
		},
		Payload: []byte{},
	}
	
	data, err := pkt.Marshal()
	if err != nil {
		c.log.Error("failed to marshal keepalive", zap.Error(err))
		return
	}
	
	if _, err := c.udpConn.Write(data); err != nil {
		c.log.Error("failed to send keepalive", zap.Error(err))
	}
}
