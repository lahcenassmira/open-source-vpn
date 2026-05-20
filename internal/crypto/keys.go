package crypto

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"golang.org/x/crypto/curve25519"
)

const (
	// KeySize is the size of X25519 keys in bytes
	KeySize = 32
)

// KeyPair represents a public/private key pair
type KeyPair struct {
	PrivateKey [KeySize]byte
	PublicKey  [KeySize]byte
}

// GenerateKeyPair generates a new X25519 key pair
func GenerateKeyPair() (*KeyPair, error) {
	var privateKey [KeySize]byte
	
	// Generate random private key
	if _, err := rand.Read(privateKey[:]); err != nil {
		return nil, fmt.Errorf("failed to generate private key: %w", err)
	}

	// Clamp the private key (X25519 requirement)
	privateKey[0] &= 248
	privateKey[31] &= 127
	privateKey[31] |= 64

	// Derive public key
	var publicKey [KeySize]byte
	curve25519.ScalarBaseMult(&publicKey, &privateKey)

	return &KeyPair{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
	}, nil
}

// SharedSecret computes the shared secret using X25519
func SharedSecret(privateKey, publicKey *[KeySize]byte) (*[KeySize]byte, error) {
	var secret [KeySize]byte
	
	curve25519.ScalarMult(&secret, privateKey, publicKey)
	
	// Check for all-zero output (invalid public key)
	var zero [KeySize]byte
	if secret == zero {
		return nil, fmt.Errorf("invalid public key")
	}
	
	return &secret, nil
}

// EncodeKey encodes a key to base64 string
func EncodeKey(key []byte) string {
	return base64.StdEncoding.EncodeToString(key)
}

// DecodeKey decodes a base64 string to key bytes
func DecodeKey(encoded string) ([KeySize]byte, error) {
	var key [KeySize]byte
	
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return key, fmt.Errorf("failed to decode key: %w", err)
	}
	
	if len(decoded) != KeySize {
		return key, fmt.Errorf("invalid key size: expected %d, got %d", KeySize, len(decoded))
	}
	
	copy(key[:], decoded)
	return key, nil
}

// PublicKeyFromPrivate derives the public key from a private key
func PublicKeyFromPrivate(privateKey *[KeySize]byte) [KeySize]byte {
	var publicKey [KeySize]byte
	curve25519.ScalarBaseMult(&publicKey, privateKey)
	return publicKey
}
