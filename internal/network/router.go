package network

import (
	"fmt"
	"net"
	"sync"
)

// Router handles packet routing decisions
type Router struct {
	routes map[string]*Route
	mu     sync.RWMutex
}

// Route represents a routing entry
type Route struct {
	Destination *net.IPNet
	Gateway     net.IP
	Interface   string
	Metric      int
}

// NewRouter creates a new router
func NewRouter() *Router {
	return &Router{
		routes: make(map[string]*Route),
	}
}

// AddRoute adds a route to the routing table
func (r *Router) AddRoute(destination *net.IPNet, gateway net.IP, iface string, metric int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	key := destination.String()
	r.routes[key] = &Route{
		Destination: destination,
		Gateway:     gateway,
		Interface:   iface,
		Metric:      metric,
	}

	return nil
}

// DeleteRoute removes a route from the routing table
func (r *Router) DeleteRoute(destination *net.IPNet) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	key := destination.String()
	delete(r.routes, key)

	return nil
}

// Lookup finds the best matching route for a destination IP
func (r *Router) Lookup(dstIP net.IP) (*Route, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var bestRoute *Route
	var bestPrefixLen int

	for _, route := range r.routes {
		if route.Destination.Contains(dstIP) {
			prefixLen, _ := route.Destination.Mask.Size()
			if bestRoute == nil || prefixLen > bestPrefixLen {
				bestRoute = route
				bestPrefixLen = prefixLen
			}
		}
	}

	if bestRoute == nil {
		return nil, fmt.Errorf("no route to host: %s", dstIP)
	}

	return bestRoute, nil
}

// GetRoutes returns all routes
func (r *Router) GetRoutes() []*Route {
	r.mu.RLock()
	defer r.mu.RUnlock()

	routes := make([]*Route, 0, len(r.routes))
	for _, route := range r.routes {
		routes = append(routes, route)
	}

	return routes
}

// Clear removes all routes
func (r *Router) Clear() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.routes = make(map[string]*Route)
}

// AllowedIPsChecker checks if an IP is in the allowed list
type AllowedIPsChecker struct {
	networks []*net.IPNet
	mu       sync.RWMutex
}

// NewAllowedIPsChecker creates a new allowed IPs checker
func NewAllowedIPsChecker() *AllowedIPsChecker {
	return &AllowedIPsChecker{
		networks: make([]*net.IPNet, 0),
	}
}

// Add adds an allowed IP network
func (a *AllowedIPsChecker) Add(cidr string) error {
	_, network, err := net.ParseCIDR(cidr)
	if err != nil {
		return fmt.Errorf("invalid CIDR: %w", err)
	}

	a.mu.Lock()
	defer a.mu.Unlock()

	a.networks = append(a.networks, network)
	return nil
}

// Contains checks if an IP is in the allowed list
func (a *AllowedIPsChecker) Contains(ip net.IP) bool {
	a.mu.RLock()
	defer a.mu.RUnlock()

	for _, network := range a.networks {
		if network.Contains(ip) {
			return true
		}
	}

	return false
}

// Clear removes all allowed IPs
func (a *AllowedIPsChecker) Clear() {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.networks = make([]*net.IPNet, 0)
}

// NATTable manages NAT translations
type NATTable struct {
	translations map[string]*NATEntry
	mu           sync.RWMutex
}

// NATEntry represents a NAT translation entry
type NATEntry struct {
	OriginalSrc  net.IP
	OriginalDst  net.IP
	TranslatedSrc net.IP
	TranslatedDst net.IP
	Protocol     uint8
	OriginalPort uint16
	TranslatedPort uint16
}

// NewNATTable creates a new NAT table
func NewNATTable() *NATTable {
	return &NATTable{
		translations: make(map[string]*NATEntry),
	}
}

// Add adds a NAT translation
func (n *NATTable) Add(key string, entry *NATEntry) {
	n.mu.Lock()
	defer n.mu.Unlock()

	n.translations[key] = entry
}

// Lookup finds a NAT translation
func (n *NATTable) Lookup(key string) (*NATEntry, bool) {
	n.mu.RLock()
	defer n.mu.RUnlock()

	entry, ok := n.translations[key]
	return entry, ok
}

// Delete removes a NAT translation
func (n *NATTable) Delete(key string) {
	n.mu.Lock()
	defer n.mu.Unlock()

	delete(n.translations, key)
}

// Clear removes all NAT translations
func (n *NATTable) Clear() {
	n.mu.Lock()
	defer n.mu.Unlock()

	n.translations = make(map[string]*NATEntry)
}
