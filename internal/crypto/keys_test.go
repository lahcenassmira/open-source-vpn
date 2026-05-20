package crypto

import (
	"testing"
)

func TestGenerateKeyPair(t *testing.T) {
	keyPair, err := GenerateKeyPair()
	if err != nil {
		t.Fatalf("GenerateKeyPair() error = %v", err)
	}

	if keyPair == nil {
		t.Fatal("GenerateKeyPair() returned nil")
	}

	// Check key sizes
	if len(keyPair.PrivateKey) != KeySize {
		t.Errorf("Private key size = %d, want %d", len(keyPair.PrivateKey), KeySize)
	}

	if len(keyPair.PublicKey) != KeySize {
		t.Errorf("Public key size = %d, want %d", len(keyPair.PublicKey), KeySize)
	}

	// Check that keys are not all zeros
	var zero [KeySize]byte
	if keyPair.PrivateKey == zero {
		t.Error("Private key is all zeros")
	}
	if keyPair.PublicKey == zero {
		t.Error("Public key is all zeros")
	}
}

func TestSharedSecret(t *testing.T) {
	// Generate two key pairs
	alice, err := GenerateKeyPair()
	if err != nil {
		t.Fatalf("Failed to generate Alice's keys: %v", err)
	}

	bob, err := GenerateKeyPair()
	if err != nil {
		t.Fatalf("Failed to generate Bob's keys: %v", err)
	}

	// Compute shared secrets
	aliceSecret, err := SharedSecret(&alice.PrivateKey, &bob.PublicKey)
	if err != nil {
		t.Fatalf("Alice's SharedSecret() error = %v", err)
	}

	bobSecret, err := SharedSecret(&bob.PrivateKey, &alice.PublicKey)
	if err != nil {
		t.Fatalf("Bob's SharedSecret() error = %v", err)
	}

	// Shared secrets should match
	if *aliceSecret != *bobSecret {
		t.Error("Shared secrets do not match")
	}

	// Shared secret should not be all zeros
	var zero [KeySize]byte
	if *aliceSecret == zero {
		t.Error("Shared secret is all zeros")
	}
}

func TestEncodeDecodeKey(t *testing.T) {
	keyPair, err := GenerateKeyPair()
	if err != nil {
		t.Fatalf("GenerateKeyPair() error = %v", err)
	}

	// Encode
	encoded := EncodeKey(keyPair.PublicKey[:])
	if encoded == "" {
		t.Error("EncodeKey() returned empty string")
	}

	// Decode
	decoded, err := DecodeKey(encoded)
	if err != nil {
		t.Fatalf("DecodeKey() error = %v", err)
	}

	// Should match original
	if decoded != keyPair.PublicKey {
		t.Error("Decoded key does not match original")
	}
}

func TestDecodeKeyInvalid(t *testing.T) {
	tests := []struct {
		name    string
		encoded string
		wantErr bool
	}{
		{
			name:    "invalid base64",
			encoded: "not-valid-base64!@#$",
			wantErr: true,
		},
		{
			name:    "wrong size",
			encoded: EncodeKey([]byte("short")),
			wantErr: true,
		},
		{
			name:    "empty string",
			encoded: "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := DecodeKey(tt.encoded)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecodeKey() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPublicKeyFromPrivate(t *testing.T) {
	keyPair, err := GenerateKeyPair()
	if err != nil {
		t.Fatalf("GenerateKeyPair() error = %v", err)
	}

	// Derive public key from private key
	derivedPublic := PublicKeyFromPrivate(&keyPair.PrivateKey)

	// Should match original public key
	if derivedPublic != keyPair.PublicKey {
		t.Error("Derived public key does not match original")
	}
}

func BenchmarkGenerateKeyPair(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := GenerateKeyPair()
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkSharedSecret(b *testing.B) {
	alice, _ := GenerateKeyPair()
	bob, _ := GenerateKeyPair()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := SharedSecret(&alice.PrivateKey, &bob.PublicKey)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkEncodeKey(b *testing.B) {
	keyPair, _ := GenerateKeyPair()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = EncodeKey(keyPair.PublicKey[:])
	}
}

func BenchmarkDecodeKey(b *testing.B) {
	keyPair, _ := GenerateKeyPair()
	encoded := EncodeKey(keyPair.PublicKey[:])

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := DecodeKey(encoded)
		if err != nil {
			b.Fatal(err)
		}
	}
}
