package metrics

import (
	"sync"
	"sync/atomic"
	"time"
)

// Metrics tracks VPN performance metrics
type Metrics struct {
	// Connection metrics
	ActiveConnections atomic.Int64
	TotalConnections  atomic.Uint64
	
	// Traffic metrics
	BytesSent     atomic.Uint64
	BytesReceived atomic.Uint64
	PacketsSent   atomic.Uint64
	PacketsReceived atomic.Uint64
	
	// Error metrics
	EncryptionErrors atomic.Uint64
	DecryptionErrors atomic.Uint64
	HandshakeErrors  atomic.Uint64
	
	// Latency tracking
	latencies []time.Duration
	latencyMu sync.RWMutex
	
	startTime time.Time
}

// NewMetrics creates a new metrics instance
func NewMetrics() *Metrics {
	return &Metrics{
		latencies: make([]time.Duration, 0, 1000),
		startTime: time.Now(),
	}
}

// ConnectionOpened increments active connections
func (m *Metrics) ConnectionOpened() {
	m.ActiveConnections.Add(1)
	m.TotalConnections.Add(1)
}

// ConnectionClosed decrements active connections
func (m *Metrics) ConnectionClosed() {
	m.ActiveConnections.Add(-1)
}

// RecordBytesSent records bytes sent
func (m *Metrics) RecordBytesSent(bytes uint64) {
	m.BytesSent.Add(bytes)
}

// RecordBytesReceived records bytes received
func (m *Metrics) RecordBytesReceived(bytes uint64) {
	m.BytesReceived.Add(bytes)
}

// RecordPacketSent records a packet sent
func (m *Metrics) RecordPacketSent() {
	m.PacketsSent.Add(1)
}

// RecordPacketReceived records a packet received
func (m *Metrics) RecordPacketReceived() {
	m.PacketsReceived.Add(1)
}

// RecordEncryptionError records an encryption error
func (m *Metrics) RecordEncryptionError() {
	m.EncryptionErrors.Add(1)
}

// RecordDecryptionError records a decryption error
func (m *Metrics) RecordDecryptionError() {
	m.DecryptionErrors.Add(1)
}

// RecordHandshakeError records a handshake error
func (m *Metrics) RecordHandshakeError() {
	m.HandshakeErrors.Add(1)
}

// RecordLatency records a latency measurement
func (m *Metrics) RecordLatency(latency time.Duration) {
	m.latencyMu.Lock()
	defer m.latencyMu.Unlock()
	
	// Keep only last 1000 measurements
	if len(m.latencies) >= 1000 {
		m.latencies = m.latencies[1:]
	}
	m.latencies = append(m.latencies, latency)
}

// GetActiveConnections returns the number of active connections
func (m *Metrics) GetActiveConnections() int64 {
	return m.ActiveConnections.Load()
}

// GetTotalConnections returns the total number of connections
func (m *Metrics) GetTotalConnections() uint64 {
	return m.TotalConnections.Load()
}

// GetBytesSent returns total bytes sent
func (m *Metrics) GetBytesSent() uint64 {
	return m.BytesSent.Load()
}

// GetBytesReceived returns total bytes received
func (m *Metrics) GetBytesReceived() uint64 {
	return m.BytesReceived.Load()
}

// GetPacketsSent returns total packets sent
func (m *Metrics) GetPacketsSent() uint64 {
	return m.PacketsSent.Load()
}

// GetPacketsReceived returns total packets received
func (m *Metrics) GetPacketsReceived() uint64 {
	return m.PacketsReceived.Load()
}

// GetAverageLatency returns the average latency
func (m *Metrics) GetAverageLatency() time.Duration {
	m.latencyMu.RLock()
	defer m.latencyMu.RUnlock()
	
	if len(m.latencies) == 0 {
		return 0
	}
	
	var total time.Duration
	for _, lat := range m.latencies {
		total += lat
	}
	
	return total / time.Duration(len(m.latencies))
}

// GetUptime returns the uptime duration
func (m *Metrics) GetUptime() time.Duration {
	return time.Since(m.startTime)
}

// Snapshot returns a snapshot of current metrics
func (m *Metrics) Snapshot() *MetricsSnapshot {
	return &MetricsSnapshot{
		ActiveConnections: m.GetActiveConnections(),
		TotalConnections:  m.GetTotalConnections(),
		BytesSent:         m.GetBytesSent(),
		BytesReceived:     m.GetBytesReceived(),
		PacketsSent:       m.GetPacketsSent(),
		PacketsReceived:   m.GetPacketsReceived(),
		EncryptionErrors:  m.EncryptionErrors.Load(),
		DecryptionErrors:  m.DecryptionErrors.Load(),
		HandshakeErrors:   m.HandshakeErrors.Load(),
		AverageLatency:    m.GetAverageLatency(),
		Uptime:            m.GetUptime(),
	}
}

// MetricsSnapshot represents a point-in-time snapshot of metrics
type MetricsSnapshot struct {
	ActiveConnections int64
	TotalConnections  uint64
	BytesSent         uint64
	BytesReceived     uint64
	PacketsSent       uint64
	PacketsReceived   uint64
	EncryptionErrors  uint64
	DecryptionErrors  uint64
	HandshakeErrors   uint64
	AverageLatency    time.Duration
	Uptime            time.Duration
}

// Throughput calculates throughput in bytes per second
func (s *MetricsSnapshot) Throughput() (sent, received float64) {
	seconds := s.Uptime.Seconds()
	if seconds == 0 {
		return 0, 0
	}
	
	sent = float64(s.BytesSent) / seconds
	received = float64(s.BytesReceived) / seconds
	return
}

// PacketRate calculates packet rate per second
func (s *MetricsSnapshot) PacketRate() (sent, received float64) {
	seconds := s.Uptime.Seconds()
	if seconds == 0 {
		return 0, 0
	}
	
	sent = float64(s.PacketsSent) / seconds
	received = float64(s.PacketsReceived) / seconds
	return
}
