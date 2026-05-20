package crypto

import (
	"crypto/cipher"
	"fmt"

	"golang.org/x/crypto/chacha20poly1305"
)

const (
	// NonceSize is the size of the nonce for ChaCha20-Poly1305
	NonceSize = chacha20poly1305.NonceSizeX
	// TagSize is the size of the authentication tag
	TagSize = 16
)

// Cipher wraps ChaCha20-Poly1305 AEAD cipher
type Cipher struct {
	aead cipher.AEAD
}

// NewCipher creates a new ChaCha20-Poly1305 cipher with the given key
func NewCipher(key *[KeySize]byte) (*Cipher, error) {
	aead, err := chacha20poly1305.NewX(key[:])
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	return &Cipher{aead: aead}, nil
}

// Encrypt encrypts plaintext with the given nonce
// Returns ciphertext with authentication tag appended
func (c *Cipher) Encrypt(nonce []byte, plaintext []byte, additionalData []byte) ([]byte, error) {
	if len(nonce) != NonceSize {
		return nil, fmt.Errorf("invalid nonce size: expected %d, got %d", NonceSize, len(nonce))
	}

	// Seal appends the ciphertext and tag to dst
	ciphertext := c.aead.Seal(nil, nonce, plaintext, additionalData)
	return ciphertext, nil
}

// Decrypt decrypts ciphertext with the given nonce
// Verifies authentication tag and returns plaintext
func (c *Cipher) Decrypt(nonce []byte, ciphertext []byte, additionalData []byte) ([]byte, error) {
	if len(nonce) != NonceSize {
		return nil, fmt.Errorf("invalid nonce size: expected %d, got %d", NonceSize, len(nonce))
	}

	if len(ciphertext) < TagSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	// Open verifies the tag and decrypts
	plaintext, err := c.aead.Open(nil, nonce, ciphertext, additionalData)
	if err != nil {
		return nil, fmt.Errorf("decryption failed: %w", err)
	}

	return plaintext, nil
}

// Overhead returns the number of bytes added by encryption (tag size)
func (c *Cipher) Overhead() int {
	return c.aead.Overhead()
}
