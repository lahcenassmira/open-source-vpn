package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// ServerConfig represents server configuration
type ServerConfig struct {
	Server  ServerSettings  `yaml:"server"`
	Network NetworkSettings `yaml:"network"`
	Crypto  CryptoSettings  `yaml:"crypto"`
	Clients []ClientPeer    `yaml:"clients"`
	Routing RoutingSettings `yaml:"routing"`
	DNS     DNSSettings     `yaml:"dns"`
	Logging LoggingSettings `yaml:"logging"`
}

// ClientConfig represents client configuration
type ClientConfig struct {
	Client  ClientSettings  `yaml:"client"`
	Network NetworkSettings `yaml:"network"`
	Crypto  ClientCryptoSettings `yaml:"crypto"`
	Routing ClientRoutingSettings `yaml:"routing"`
	DNS     DNSSettings     `yaml:"dns"`
	Logging LoggingSettings `yaml:"logging"`
}

// ServerSettings contains server-specific settings
type ServerSettings struct {
	ListenAddress string `yaml:"listen_address"`
	Protocol      string `yaml:"protocol"`
	Workers       int    `yaml:"workers"`
}

// ClientSettings contains client-specific settings
type ClientSettings struct {
	ServerAddress string `yaml:"server_address"`
	Protocol      string `yaml:"protocol"`
	ReconnectDelay int   `yaml:"reconnect_delay"`
}

// NetworkSettings contains network configuration
type NetworkSettings struct {
	Interface string `yaml:"interface"`
	Address   string `yaml:"address"`
	MTU       int    `yaml:"mtu"`
}

// CryptoSettings contains server crypto configuration
type CryptoSettings struct {
	PrivateKey string `yaml:"private_key"`
}

// ClientCryptoSettings contains client crypto configuration
type ClientCryptoSettings struct {
	PrivateKey      string `yaml:"private_key"`
	ServerPublicKey string `yaml:"server_public_key"`
}

// ClientPeer represents an allowed client
type ClientPeer struct {
	PublicKey  string   `yaml:"public_key"`
	AllowedIPs []string `yaml:"allowed_ips"`
	Name       string   `yaml:"name,omitempty"`
}

// RoutingSettings contains server routing configuration
type RoutingSettings struct {
	ForwardAllTraffic bool     `yaml:"forward_all_traffic"`
	AllowedNetworks   []string `yaml:"allowed_networks"`
}

// ClientRoutingSettings contains client routing configuration
type ClientRoutingSettings struct {
	DefaultRoute bool     `yaml:"default_route"`
	Routes       []string `yaml:"routes"`
}

// DNSSettings contains DNS configuration
type DNSSettings struct {
	Enabled bool     `yaml:"enabled"`
	Servers []string `yaml:"servers"`
}

// LoggingSettings contains logging configuration
type LoggingSettings struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
	Output string `yaml:"output"`
}

// LoadServerConfig loads server configuration from a file
func LoadServerConfig(path string) (*ServerConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config ServerConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	// Set defaults
	if config.Server.Protocol == "" {
		config.Server.Protocol = "udp"
	}
	if config.Server.Workers == 0 {
		config.Server.Workers = 4
	}
	if config.Network.MTU == 0 {
		config.Network.MTU = 1420
	}
	if config.Logging.Level == "" {
		config.Logging.Level = "info"
	}
	if config.Logging.Format == "" {
		config.Logging.Format = "json"
	}

	return &config, nil
}

// LoadClientConfig loads client configuration from a file
func LoadClientConfig(path string) (*ClientConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config ClientConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	// Set defaults
	if config.Client.Protocol == "" {
		config.Client.Protocol = "udp"
	}
	if config.Client.ReconnectDelay == 0 {
		config.Client.ReconnectDelay = 5
	}
	if config.Network.MTU == 0 {
		config.Network.MTU = 1420
	}
	if config.Logging.Level == "" {
		config.Logging.Level = "info"
	}
	if config.Logging.Format == "" {
		config.Logging.Format = "json"
	}

	return &config, nil
}

// SaveServerConfig saves server configuration to a file
func SaveServerConfig(path string, config *ServerConfig) error {
	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// SaveClientConfig saves client configuration to a file
func SaveClientConfig(path string, config *ClientConfig) error {
	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// Validate validates server configuration
func (c *ServerConfig) Validate() error {
	if c.Server.ListenAddress == "" {
		return fmt.Errorf("server.listen_address is required")
	}
	if c.Network.Interface == "" {
		return fmt.Errorf("network.interface is required")
	}
	if c.Network.Address == "" {
		return fmt.Errorf("network.address is required")
	}
	if c.Crypto.PrivateKey == "" {
		return fmt.Errorf("crypto.private_key is required")
	}
	return nil
}

// Validate validates client configuration
func (c *ClientConfig) Validate() error {
	if c.Client.ServerAddress == "" {
		return fmt.Errorf("client.server_address is required")
	}
	if c.Network.Interface == "" {
		return fmt.Errorf("network.interface is required")
	}
	if c.Network.Address == "" {
		return fmt.Errorf("network.address is required")
	}
	if c.Crypto.PrivateKey == "" {
		return fmt.Errorf("crypto.private_key is required")
	}
	if c.Crypto.ServerPublicKey == "" {
		return fmt.Errorf("crypto.server_public_key is required")
	}
	return nil
}
