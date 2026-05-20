package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/open-source-vpn/vpn/internal/config"
	"github.com/open-source-vpn/vpn/internal/crypto"
	"github.com/open-source-vpn/vpn/pkg/logger"
	"github.com/spf13/cobra"
)

var (
	configPath string
	outputPath string
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "vpn-client",
		Short: "open-source-vpn Client",
		Long:  "A secure, high-performance VPN client written in Go",
	}

	connectCmd := &cobra.Command{
		Use:   "connect",
		Short: "Connect to VPN server",
		RunE:  runConnect,
	}
	connectCmd.Flags().StringVarP(&configPath, "config", "c", "client.yaml", "Path to configuration file")

	disconnectCmd := &cobra.Command{
		Use:   "disconnect",
		Short: "Disconnect from VPN server",
		RunE:  runDisconnect,
	}

	keygenCmd := &cobra.Command{
		Use:   "keygen",
		Short: "Generate a new key pair",
		RunE:  runKeygen,
	}
	keygenCmd.Flags().StringVarP(&outputPath, "output", "o", "client-keys.json", "Output file for keys")

	statusCmd := &cobra.Command{
		Use:   "status",
		Short: "Show connection status",
		RunE:  runStatus,
	}

	rootCmd.AddCommand(connectCmd, disconnectCmd, keygenCmd, statusCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func runConnect(cmd *cobra.Command, args []string) error {
	// Load configuration
	cfg, err := config.LoadClientConfig(configPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("invalid config: %w", err)
	}

	// Initialize logger
	if err := logger.InitDefault(logger.Config{
		Level:  cfg.Logging.Level,
		Format: cfg.Logging.Format,
		Output: cfg.Logging.Output,
	}); err != nil {
		return fmt.Errorf("failed to initialize logger: %w", err)
	}
	defer logger.Sync()

	log := logger.Default().WithComponent("client")
	log.Info("Starting VPN client")

	// Create and start client
	client, err := NewClient(cfg, log)
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	if err := client.Connect(); err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}

	return nil
}

func runDisconnect(cmd *cobra.Command, args []string) error {
	fmt.Println("Disconnecting from VPN...")
	fmt.Println("Not implemented yet - use Ctrl+C to stop the client")
	return nil
}

func runKeygen(cmd *cobra.Command, args []string) error {
	// Generate key pair
	keyPair, err := crypto.GenerateKeyPair()
	if err != nil {
		return fmt.Errorf("failed to generate key pair: %w", err)
	}

	// Encode keys
	privateKey := crypto.EncodeKey(keyPair.PrivateKey[:])
	publicKey := crypto.EncodeKey(keyPair.PublicKey[:])

	// Create output structure
	output := map[string]string{
		"private_key": privateKey,
		"public_key":  publicKey,
	}

	// Write to file
	data, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal keys: %w", err)
	}

	if err := os.WriteFile(outputPath, data, 0600); err != nil {
		return fmt.Errorf("failed to write keys: %w", err)
	}

	fmt.Printf("Key pair generated successfully!\n")
	fmt.Printf("Private key: %s\n", privateKey)
	fmt.Printf("Public key:  %s\n", publicKey)
	fmt.Printf("\nKeys saved to: %s\n", outputPath)
	fmt.Printf("\n⚠️  Keep your private key secure and never share it!\n")
	fmt.Printf("\n📋 Add this public key to the server configuration:\n")
	fmt.Printf("   public_key: \"%s\"\n", publicKey)

	return nil
}

func runStatus(cmd *cobra.Command, args []string) error {
	fmt.Println("Connection status: Not implemented yet")
	fmt.Println("Check logs for connection status")
	return nil
}
