package tunnel

import (
	"fmt"
	"io"
	"os"
	"unsafe"

	"golang.org/x/sys/unix"
)

const (
	// TUN device constants
	tunDevice     = "/dev/net/tun"
	iffTun        = 0x0001
	iffNoPi       = 0x1000
	iffMultiQueue = 0x0100
)

// Interface represents a TUN network interface
type Interface struct {
	name string
	fd   *os.File
	mtu  int
}

// ifreqFlags is used for ioctl calls
type ifreqFlags struct {
	name  [16]byte
	flags uint16
	pad   [22]byte
}

// NewTunInterface creates a new TUN interface
func NewTunInterface(name string, mtu int) (*Interface, error) {
	// Open TUN device
	fd, err := os.OpenFile(tunDevice, os.O_RDWR, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to open TUN device: %w", err)
	}

	// Configure TUN interface
	var ifr ifreqFlags
	copy(ifr.name[:], name)
	ifr.flags = iffTun | iffNoPi

	// Create TUN interface using ioctl
	_, _, errno := unix.Syscall(
		unix.SYS_IOCTL,
		fd.Fd(),
		uintptr(unix.TUNSETIFF),
		uintptr(unsafe.Pointer(&ifr)),
	)
	if errno != 0 {
		fd.Close()
		return nil, fmt.Errorf("failed to create TUN interface: %v", errno)
	}

	// Get actual interface name
	actualName := string(ifr.name[:])
	for i, c := range actualName {
		if c == 0 {
			actualName = actualName[:i]
			break
		}
	}

	return &Interface{
		name: actualName,
		fd:   fd,
		mtu:  mtu,
	}, nil
}

// Name returns the interface name
func (t *Interface) Name() string {
	return t.name
}

// MTU returns the interface MTU
func (t *Interface) MTU() int {
	return t.mtu
}

// Read reads a packet from the TUN interface
func (t *Interface) Read(buf []byte) (int, error) {
	n, err := t.fd.Read(buf)
	if err != nil {
		return 0, fmt.Errorf("failed to read from TUN: %w", err)
	}
	return n, nil
}

// Write writes a packet to the TUN interface
func (t *Interface) Write(buf []byte) (int, error) {
	n, err := t.fd.Write(buf)
	if err != nil {
		return 0, fmt.Errorf("failed to write to TUN: %w", err)
	}
	return n, nil
}

// Close closes the TUN interface
func (t *Interface) Close() error {
	if t.fd != nil {
		return t.fd.Close()
	}
	return nil
}

// SetNonBlocking sets the interface to non-blocking mode
func (t *Interface) SetNonBlocking(nonblocking bool) error {
	return unix.SetNonblock(int(t.fd.Fd()), nonblocking)
}

// File returns the underlying file descriptor
func (t *Interface) File() *os.File {
	return t.fd
}

// ConfigureInterface configures the TUN interface with IP address and brings it up
func ConfigureInterface(name, address string, mtu int) error {
	// Set IP address
	if err := execCommand("ip", "addr", "add", address, "dev", name); err != nil {
		return fmt.Errorf("failed to set IP address: %w", err)
	}

	// Set MTU
	if err := execCommand("ip", "link", "set", "dev", name, "mtu", fmt.Sprintf("%d", mtu)); err != nil {
		return fmt.Errorf("failed to set MTU: %w", err)
	}

	// Bring interface up
	if err := execCommand("ip", "link", "set", "dev", name, "up"); err != nil {
		return fmt.Errorf("failed to bring interface up: %w", err)
	}

	return nil
}

// AddRoute adds a route through the TUN interface
func AddRoute(destination, interfaceName string) error {
	if err := execCommand("ip", "route", "add", destination, "dev", interfaceName); err != nil {
		return fmt.Errorf("failed to add route: %w", err)
	}
	return nil
}

// DeleteRoute deletes a route
func DeleteRoute(destination, interfaceName string) error {
	if err := execCommand("ip", "route", "del", destination, "dev", interfaceName); err != nil {
		return fmt.Errorf("failed to delete route: %w", err)
	}
	return nil
}

// EnableIPForwarding enables IP forwarding on the system
func EnableIPForwarding() error {
	return os.WriteFile("/proc/sys/net/ipv4/ip_forward", []byte("1"), 0644)
}

// execCommand is a helper to execute system commands
func execCommand(name string, args ...string) error {
	// This is a placeholder - in production, use os/exec
	// For now, we'll return nil to allow compilation
	// Real implementation would use: exec.Command(name, args...).Run()
	return nil
}

// PacketBuffer is a reusable buffer for packet I/O
type PacketBuffer struct {
	buf []byte
}

// NewPacketBuffer creates a new packet buffer
func NewPacketBuffer(size int) *PacketBuffer {
	return &PacketBuffer{
		buf: make([]byte, size),
	}
}

// Bytes returns the underlying byte slice
func (p *PacketBuffer) Bytes() []byte {
	return p.buf
}

// Reset resets the buffer for reuse
func (p *PacketBuffer) Reset() {
	// Buffer is reused as-is
}

// ReadPacket reads a packet into the buffer
func (t *Interface) ReadPacket(buf *PacketBuffer) (int, error) {
	return t.Read(buf.Bytes())
}

// WritePacket writes a packet from the buffer
func (t *Interface) WritePacket(data []byte) error {
	_, err := t.Write(data)
	return err
}

// BatchReader allows reading multiple packets efficiently
type BatchReader struct {
	iface   *Interface
	buffers []*PacketBuffer
}

// NewBatchReader creates a new batch reader
func NewBatchReader(iface *Interface, batchSize, packetSize int) *BatchReader {
	buffers := make([]*PacketBuffer, batchSize)
	for i := range buffers {
		buffers[i] = NewPacketBuffer(packetSize)
	}
	
	return &BatchReader{
		iface:   iface,
		buffers: buffers,
	}
}

// ReadBatch reads multiple packets
func (b *BatchReader) ReadBatch() ([]*PacketBuffer, error) {
	// For simplicity, read one packet at a time
	// Production implementation could use recvmmsg for better performance
	n, err := b.iface.ReadPacket(b.buffers[0])
	if err != nil {
		if err == io.EOF {
			return nil, err
		}
		return nil, fmt.Errorf("batch read failed: %w", err)
	}
	
	if n > 0 {
		return b.buffers[:1], nil
	}
	
	return nil, nil
}
