package protocol

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/lahcenassmira/open-source-vpn/internal/crypto"
	"github.com/lahcenassmira/open-source-vpn/internal/network"
)

// Connection represents a VPN connection
type Connection struct {
	id          uint32
	remoteAddr  *net.UDPAddr
	sendCipher  *crypto.Cipher
	recvCipher  *crypto.Cipher
	
	sendNonce   uint64
	recvNonce   uint64
	replayWindow *ReplayWindow
	
	lastSeen    time.Time
	lastSeenMu  sync.RWMutex
	
	sendQueue   chan []byte
	recvQueue   chan []byte
	
	closed      atomic.Bool
	closeChan   chan struct{}
	
	stats       *ConnectionStats
}

// ConnectionStats tracks connection statistics
type ConnectionStats struct {
	BytesSent     atomic.Uint64
	BytesReceived atomic.Uint64
	PacketsSent   atomic.Uint64
	PacketsReceived atomic.Uint64
	LastHandshake time.Time
}

// NewConnection creates a new VPN connection
func NewConnection(id uint32, remoteAddr *net.UDPAddr, sendCipher, recvCipher *crypto.Cipher) *Connection {
	return &Connection{
		id:           id,
		remoteAddr:   remoteAddr,
		sendCipher:   sendCipher,
		recvCipher:   recvCipher,
		replayWindow: NewReplayWindow(1024),
		lastSeen:     time.Now(),
		sendQueue:    make(chan []byte, 256),
		recvQueue:    make(chan []byte, 256),
		closeChan:    make(chan struct{}),
		stats:        &ConnectionStats{},
	}
}

// ID returns the connection ID
func (c *Connection) ID() uint32 {
	return c.id
}

// RemoteAddr returns the remote address
func (c *Connection) RemoteAddr() *net.UDPAddr {
	return c.remoteAddr
}

// UpdateLastSeen updates the last seen timestamp
func (c *Connection) UpdateLastSeen() {
	c.lastSeenMu.Lock()
	c.lastSeen = time.Now()
	c.lastSeenMu.Unlock()
}

// LastSeen returns the last seen timestamp
func (c *Connection) LastSeen() time.Time {
	c.lastSeenMu.RLock()
	defer c.lastSeenMu.RUnlock()
	return c.lastSeen
}

// IsExpired checks if the connection has expired
func (c *Connection) IsExpired(timeout time.Duration) bool {
	return time.Since(c.LastSeen()) > timeout
}

// EncryptPacket encrypts a packet for transmission
func (c *Connection) EncryptPacket(packetType network.PacketType, payload []byte) ([]byte, error) {
	if c.closed.Load() {
		return nil, fmt.Errorf("connection closed")
	}

	// Generate nonce
	nonce := make([]byte, crypto.NonceSize)
	nonceValue := atomic.AddUint64(&c.sendNonce, 1)
	binary.BigEndian.PutUint64(nonce[16:], nonceValue)
	if _, err := rand.Read(nonce[:16]); err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %w", err)
	}

	// Encrypt payload
	encrypted, err := c.sendCipher.Encrypt(nonce, payload, nil)
	if err != nil {
		return nil, fmt.Errorf("encryption failed: %w", err)
	}

	// Build packet
	packet := &network.Packet{
		Header: network.PacketHeader{
			Type:   packetType,
			ConnID: c.id,
		},
		Payload: encrypted,
	}
	copy(packet.Header.Nonce[:], nonce)

	data, err := packet.Marshal()
	if err != nil {
		return nil, err
	}

	// Update stats
	c.stats.BytesSent.Add(uint64(len(data)))
	c.stats.PacketsSent.Add(1)

	return data, nil
}

// DecryptPacket decrypts a received packet
func (c *Connection) DecryptPacket(packet *network.Packet) ([]byte, error) {
	if c.closed.Load() {
		return nil, fmt.Errorf("connection closed")
	}

	// Extract nonce
	nonce := packet.Header.Nonce[:]

	// Check for replay attacks
	nonceValue := binary.BigEndian.Uint64(nonce[16:])
	if !c.replayWindow.Check(nonceValue) {
		return nil, fmt.Errorf("replay attack detected")
	}

	// Decrypt payload
	decrypted, err := c.recvCipher.Decrypt(nonce, packet.Payload, nil)
	if err != nil {
		return nil, fmt.Errorf("decryption failed: %w", err)
	}

	// Mark nonce as seen
	c.replayWindow.Update(nonceValue)

	// Update stats
	c.stats.BytesReceived.Add(uint64(len(packet.Payload)))
	c.stats.PacketsReceived.Add(1)
	c.UpdateLastSeen()

	return decrypted, nil
}

// SendPacket queues a packet for sending
func (c *Connection) SendPacket(data []byte) error {
	if c.closed.Load() {
		return fmt.Errorf("connection closed")
	}

	select {
	case c.sendQueue <- data:
		return nil
	case <-c.closeChan:
		return fmt.Errorf("connection closed")
	default:
		return fmt.Errorf("send queue full")
	}
}

// ReceivePacket retrieves a received packet
func (c *Connection) ReceivePacket() ([]byte, error) {
	if c.closed.Load() {
		return nil, fmt.Errorf("connection closed")
	}

	select {
	case data := <-c.recvQueue:
		return data, nil
	case <-c.closeChan:
		return nil, fmt.Errorf("connection closed")
	}
}

// Close closes the connection
func (c *Connection) Close() error {
	if c.closed.Swap(true) {
		return nil // Already closed
	}

	close(c.closeChan)
	close(c.sendQueue)
	close(c.recvQueue)

	return nil
}

// IsClosed returns true if the connection is closed
func (c *Connection) IsClosed() bool {
	return c.closed.Load()
}

// Stats returns connection statistics
func (c *Connection) Stats() *ConnectionStats {
	return c.stats
}

// ReplayWindow implements an anti-replay window
type ReplayWindow struct {
	window   []uint64
	size     uint64
	lastSeen uint64
	mu       sync.Mutex
}

// NewReplayWindow creates a new replay window
func NewReplayWindow(size uint64) *ReplayWindow {
	return &ReplayWindow{
		window: make([]uint64, (size+63)/64),
		size:   size,
	}
}

// Check checks if a nonce is valid (not replayed)
func (r *ReplayWindow) Check(nonce uint64) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Nonce too old
	if nonce+r.size < r.lastSeen {
		return false
	}

	// Nonce in the future
	if nonce > r.lastSeen {
		return true
	}

	// Check if already seen
	diff := r.lastSeen - nonce
	index := diff / 64
	bit := diff % 64

	if index >= uint64(len(r.window)) {
		return false
	}

	return (r.window[index] & (1 << bit)) == 0
}

// Update marks a nonce as seen
func (r *ReplayWindow) Update(nonce uint64) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Nonce in the future, shift window
	if nonce > r.lastSeen {
		diff := nonce - r.lastSeen

		// Shift window
		if diff < r.size {
			shift := diff / 64
			for i := len(r.window) - 1; i >= int(shift); i-- {
				r.window[i] = r.window[i-int(shift)]
			}
			for i := 0; i < int(shift); i++ {
				r.window[i] = 0
			}
		} else {
			// Clear entire window
			for i := range r.window {
				r.window[i] = 0
			}
		}

		r.lastSeen = nonce
	}

	// Mark as seen
	diff := r.lastSeen - nonce
	index := diff / 64
	bit := diff % 64

	if index < uint64(len(r.window)) {
		r.window[index] |= 1 << bit
	}
}
