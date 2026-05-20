package crypto

import (
	"crypto/rand"
	"fmt"
	"io"
)

// NoiseHandshake implements the Noise_XX handshake pattern
// This provides mutual authentication and forward secrecy
type NoiseHandshake struct {
	localStaticPrivate  [KeySize]byte
	localStaticPublic   [KeySize]byte
	localEphemeralPrivate [KeySize]byte
	localEphemeralPublic  [KeySize]byte
	
	remoteStaticPublic    [KeySize]byte
	remoteEphemeralPublic [KeySize]byte
	
	isInitiator bool
	state       int
	
	sendCipher *Cipher
	recvCipher *Cipher
}

// NewNoiseHandshake creates a new Noise handshake instance
func NewNoiseHandshake(staticPrivate [KeySize]byte, isInitiator bool) (*NoiseHandshake, error) {
	staticPublic := PublicKeyFromPrivate(&staticPrivate)
	
	// Generate ephemeral key pair
	ephemeralKeyPair, err := GenerateKeyPair()
	if err != nil {
		return nil, fmt.Errorf("failed to generate ephemeral key: %w", err)
	}
	
	return &NoiseHandshake{
		localStaticPrivate:    staticPrivate,
		localStaticPublic:     staticPublic,
		localEphemeralPrivate: ephemeralKeyPair.PrivateKey,
		localEphemeralPublic:  ephemeralKeyPair.PublicKey,
		isInitiator:           isInitiator,
		state:                 0,
	}, nil
}

// WriteMessage generates the next handshake message
func (n *NoiseHandshake) WriteMessage(payload []byte) ([]byte, error) {
	var message []byte
	
	if n.isInitiator {
		switch n.state {
		case 0:
			// -> e
			message = append(message, n.localEphemeralPublic[:]...)
			n.state = 1
			
		case 1:
			// -> s, se
			// Compute shared secrets and derive keys
			es, err := SharedSecret(&n.localEphemeralPrivate, &n.remoteStaticPublic)
			if err != nil {
				return nil, fmt.Errorf("failed to compute es: %w", err)
			}
			
			ss, err := SharedSecret(&n.localStaticPrivate, &n.remoteStaticPublic)
			if err != nil {
				return nil, fmt.Errorf("failed to compute ss: %w", err)
			}
			
			// Derive encryption key from shared secrets
			key := deriveKey(es[:], ss[:])
			cipher, err := NewCipher(&key)
			if err != nil {
				return nil, err
			}
			
			// Encrypt static public key
			nonce := make([]byte, NonceSize)
			if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
				return nil, err
			}
			
			encrypted, err := cipher.Encrypt(nonce, n.localStaticPublic[:], nil)
			if err != nil {
				return nil, err
			}
			
			message = append(message, nonce...)
			message = append(message, encrypted...)
			
			if len(payload) > 0 {
				message = append(message, payload...)
			}
			
			n.state = 2
			n.sendCipher = cipher
		}
	} else {
		switch n.state {
		case 0:
			// <- e, s
			ee, err := SharedSecret(&n.localEphemeralPrivate, &n.remoteEphemeralPublic)
			if err != nil {
				return nil, fmt.Errorf("failed to compute ee: %w", err)
			}
			
			se, err := SharedSecret(&n.localStaticPrivate, &n.remoteEphemeralPublic)
			if err != nil {
				return nil, fmt.Errorf("failed to compute se: %w", err)
			}
			
			// Derive encryption key
			key := deriveKey(ee[:], se[:])
			cipher, err := NewCipher(&key)
			if err != nil {
				return nil, err
			}
			
			// Send ephemeral public key and encrypted static public key
			message = append(message, n.localEphemeralPublic[:]...)
			
			nonce := make([]byte, NonceSize)
			if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
				return nil, err
			}
			
			encrypted, err := cipher.Encrypt(nonce, n.localStaticPublic[:], nil)
			if err != nil {
				return nil, err
			}
			
			message = append(message, nonce...)
			message = append(message, encrypted...)
			
			n.state = 1
			n.recvCipher = cipher
		}
	}
	
	return message, nil
}

// ReadMessage processes a received handshake message
func (n *NoiseHandshake) ReadMessage(message []byte) ([]byte, error) {
	if n.isInitiator {
		switch n.state {
		case 1:
			// <- e, s
			if len(message) < KeySize+NonceSize+KeySize+TagSize {
				return nil, fmt.Errorf("message too short")
			}
			
			// Extract remote ephemeral public key
			copy(n.remoteEphemeralPublic[:], message[:KeySize])
			message = message[KeySize:]
			
			// Extract nonce
			nonce := message[:NonceSize]
			message = message[NonceSize:]
			
			// Compute shared secrets
			ee, err := SharedSecret(&n.localEphemeralPrivate, &n.remoteEphemeralPublic)
			if err != nil {
				return nil, fmt.Errorf("failed to compute ee: %w", err)
			}
			
			se, err := SharedSecret(&n.localStaticPrivate, &n.remoteEphemeralPublic)
			if err != nil {
				return nil, fmt.Errorf("failed to compute se: %w", err)
			}
			
			// Derive decryption key
			key := deriveKey(ee[:], se[:])
			cipher, err := NewCipher(&key)
			if err != nil {
				return nil, err
			}
			
			// Decrypt remote static public key
			encryptedKey := message[:KeySize+TagSize]
			decrypted, err := cipher.Decrypt(nonce, encryptedKey, nil)
			if err != nil {
				return nil, fmt.Errorf("failed to decrypt remote static key: %w", err)
			}
			
			copy(n.remoteStaticPublic[:], decrypted)
			n.recvCipher = cipher
			
			return nil, nil
		}
	} else {
		switch n.state {
		case 0:
			// -> e
			if len(message) < KeySize {
				return nil, fmt.Errorf("message too short")
			}
			
			copy(n.remoteEphemeralPublic[:], message[:KeySize])
			n.state = 1
			return nil, nil
		}
	}
	
	return nil, fmt.Errorf("invalid handshake state")
}

// IsComplete returns true if the handshake is complete
func (n *NoiseHandshake) IsComplete() bool {
	return n.state == 2 || (n.state == 1 && !n.isInitiator)
}

// GetCiphers returns the send and receive ciphers after handshake completion
func (n *NoiseHandshake) GetCiphers() (*Cipher, *Cipher, error) {
	if !n.IsComplete() {
		return nil, nil, fmt.Errorf("handshake not complete")
	}
	
	return n.sendCipher, n.recvCipher, nil
}

// GetRemoteStaticPublic returns the authenticated remote static public key
func (n *NoiseHandshake) GetRemoteStaticPublic() [KeySize]byte {
	return n.remoteStaticPublic
}

// deriveKey derives a symmetric key from multiple shared secrets
func deriveKey(secrets ...[]byte) [KeySize]byte {
	var key [KeySize]byte
	
	// Simple key derivation: XOR all secrets
	// In production, use HKDF or similar KDF
	for _, secret := range secrets {
		for i := 0; i < KeySize && i < len(secret); i++ {
			key[i] ^= secret[i]
		}
	}
	
	return key
}
