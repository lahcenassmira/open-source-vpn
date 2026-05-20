package network

import (
	"encoding/binary"
	"fmt"
	"net"
)

// PacketType represents the type of VPN packet
type PacketType uint8

const (
	// PacketTypeHandshakeInit is the initial handshake packet
	PacketTypeHandshakeInit PacketType = 1
	// PacketTypeHandshakeResp is the handshake response packet
	PacketTypeHandshakeResp PacketType = 2
	// PacketTypeData is a data packet
	PacketTypeData PacketType = 3
	// PacketTypeKeepalive is a keepalive packet
	PacketTypeKeepalive PacketType = 4
)

const (
	// PacketHeaderSize is the size of the packet header
	PacketHeaderSize = 1 + 4 + 24 // Type(1) + ConnID(4) + Nonce(24)
	// MaxPacketSize is the maximum packet size
	MaxPacketSize = 65535
)

// PacketHeader represents a VPN packet header
type PacketHeader struct {
	Type   PacketType
	ConnID uint32
	Nonce  [24]byte
}

// Packet represents a VPN packet
type Packet struct {
	Header  PacketHeader
	Payload []byte
}

// Marshal serializes a packet to bytes
func (p *Packet) Marshal() ([]byte, error) {
	size := PacketHeaderSize + len(p.Payload)
	if size > MaxPacketSize {
		return nil, fmt.Errorf("packet too large: %d bytes", size)
	}

	buf := make([]byte, size)
	
	// Write header
	buf[0] = byte(p.Header.Type)
	binary.BigEndian.PutUint32(buf[1:5], p.Header.ConnID)
	copy(buf[5:29], p.Header.Nonce[:])
	
	// Write payload
	copy(buf[29:], p.Payload)
	
	return buf, nil
}

// Unmarshal deserializes a packet from bytes
func (p *Packet) Unmarshal(data []byte) error {
	if len(data) < PacketHeaderSize {
		return fmt.Errorf("packet too short: %d bytes", len(data))
	}

	// Read header
	p.Header.Type = PacketType(data[0])
	p.Header.ConnID = binary.BigEndian.Uint32(data[1:5])
	copy(p.Header.Nonce[:], data[5:29])
	
	// Read payload
	p.Payload = make([]byte, len(data)-PacketHeaderSize)
	copy(p.Payload, data[29:])
	
	return nil
}

// IPPacket represents an IP packet
type IPPacket struct {
	Version  uint8
	Protocol uint8
	SrcIP    net.IP
	DstIP    net.IP
	Payload  []byte
}

// ParseIPPacket parses an IP packet from raw bytes
func ParseIPPacket(data []byte) (*IPPacket, error) {
	if len(data) < 20 {
		return nil, fmt.Errorf("packet too short for IP header")
	}

	pkt := &IPPacket{}
	
	// Parse IP version
	pkt.Version = data[0] >> 4
	
	if pkt.Version == 4 {
		return parseIPv4Packet(data)
	} else if pkt.Version == 6 {
		return parseIPv6Packet(data)
	}
	
	return nil, fmt.Errorf("unsupported IP version: %d", pkt.Version)
}

// parseIPv4Packet parses an IPv4 packet
func parseIPv4Packet(data []byte) (*IPPacket, error) {
	if len(data) < 20 {
		return nil, fmt.Errorf("packet too short for IPv4 header")
	}

	pkt := &IPPacket{
		Version: 4,
	}
	
	// Header length in 32-bit words
	headerLen := int(data[0]&0x0F) * 4
	if len(data) < headerLen {
		return nil, fmt.Errorf("packet too short for IPv4 header length")
	}
	
	// Protocol
	pkt.Protocol = data[9]
	
	// Source and destination IP
	pkt.SrcIP = net.IP(data[12:16])
	pkt.DstIP = net.IP(data[16:20])
	
	// Payload
	pkt.Payload = data[headerLen:]
	
	return pkt, nil
}

// parseIPv6Packet parses an IPv6 packet
func parseIPv6Packet(data []byte) (*IPPacket, error) {
	if len(data) < 40 {
		return nil, fmt.Errorf("packet too short for IPv6 header")
	}

	pkt := &IPPacket{
		Version: 6,
	}
	
	// Next header (protocol)
	pkt.Protocol = data[6]
	
	// Source and destination IP
	pkt.SrcIP = net.IP(data[8:24])
	pkt.DstIP = net.IP(data[24:40])
	
	// Payload
	pkt.Payload = data[40:]
	
	return pkt, nil
}

// IsIPv4 returns true if this is an IPv4 packet
func (p *IPPacket) IsIPv4() bool {
	return p.Version == 4
}

// IsIPv6 returns true if this is an IPv6 packet
func (p *IPPacket) IsIPv6() bool {
	return p.Version == 6
}

// String returns a string representation of the packet
func (p *IPPacket) String() string {
	return fmt.Sprintf("IPv%d %s -> %s (proto: %d, len: %d)",
		p.Version, p.SrcIP, p.DstIP, p.Protocol, len(p.Payload))
}

// PacketPool manages a pool of reusable packet buffers
type PacketPool struct {
	buffers chan []byte
	size    int
}

// NewPacketPool creates a new packet pool
func NewPacketPool(poolSize, bufferSize int) *PacketPool {
	pool := &PacketPool{
		buffers: make(chan []byte, poolSize),
		size:    bufferSize,
	}
	
	// Pre-allocate buffers
	for i := 0; i < poolSize; i++ {
		pool.buffers <- make([]byte, bufferSize)
	}
	
	return pool
}

// Get retrieves a buffer from the pool
func (p *PacketPool) Get() []byte {
	select {
	case buf := <-p.buffers:
		return buf
	default:
		// Pool exhausted, allocate new buffer
		return make([]byte, p.size)
	}
}

// Put returns a buffer to the pool
func (p *PacketPool) Put(buf []byte) {
	if len(buf) != p.size {
		return // Don't return wrong-sized buffers
	}
	
	select {
	case p.buffers <- buf:
	default:
		// Pool full, let GC handle it
	}
}
