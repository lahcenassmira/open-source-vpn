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
		Use:   "vpn-server",
		Short: "open-source-vpn Server",
		Long:  "A secure, high-performance VPN server written in Go",
	}

	startCmd := &cobra.Command{
		Use:   "start",
		Short: "Start the VPN server",
		RunE:  runStart,
	}
	startCmd.Flags().StringVarP(&configPath, "config", "c", "server.yaml", "Path to configuration file")

	keygenCmd := &cobra.Command{
		Use:   "keygen",
		Short: "Generate a new key pair",
		RunE:  runKeygen,
	}
	keygenCmd.Flags().StringVarP(&outputPath, "output", "o", "server-keys.json", "Output file for keys")

	statusCmd := &cobra.Command{
		Use:   "status",
		Short: "Show server status",
		RunE:  runStatus,
	}

	rootCmd.AddCommand(startCmd, keygenCmd, statusCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func runStart(cmd *cobra.Command, args []string) error {
	// Load configuration
	cfg, err := config.LoadServerConfig(configPath)
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

	log := logger.Default().WithComponent("server")
	log.Info("Starting VPN server")

	// Create and start server
	server, err := NewServer(cfg, log)
	if err != nil {
		return fmt.Errorf("failed to create server: %w", err)
	}

	if err := server.Start(); err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}

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

	return nil
}

func runStatus(cmd *cobra.Command, args []string) error {
	fmt.Println("Server status: Not implemented yet")
	fmt.Println("Use 'systemctl status vpn-server' or check logs for server status")
	return nil
}
