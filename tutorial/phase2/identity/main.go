package main

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

type PeerIdentity struct {
	PublicKey  ed25519.PublicKey
	PrivateKey ed25519.PrivateKey
	ID         string
}

func GenerateIdentity() (*PeerIdentity, error) {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}

	// Use base64 encoded public key as peer ID
	id := base64.StdEncoding.EncodeToString(pub)

	return &PeerIdentity{
		PublicKey:  pub,
		PrivateKey: priv,
		ID:         id,
	}, nil
}

func (pi *PeerIdentity) Sign(message []byte) []byte {
	return ed25519.Sign(pi.PrivateKey, message)
}

func VerifySignature(publicKey ed25519.PublicKey, message, signature []byte) bool {
	return ed25519.Verify(publicKey, message, signature)
}

func main() {
	fmt.Println("=== AWL Tutorial: Cryptographic Identity ===\n")

	// Generate identity
	identity, err := GenerateIdentity()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Generated peer identity:\n")
	fmt.Printf("ID: %s\n", identity.ID[:20]+"...") // Show first 20 chars
	fmt.Printf("Public Key: %x\n", identity.PublicKey[:8])  // Show first 8 bytes
	fmt.Printf("Private Key: %x (never share this!)\n", identity.PrivateKey[:8])

	// Test signing
	message := []byte("Hello, mesh VPN world!")
	signature := identity.Sign(message)

	fmt.Printf("\nMessage: %s\n", message)
	fmt.Printf("Signature: %x...\n", signature[:8]) // Show first 8 bytes

	// Verify signature
	valid := VerifySignature(identity.PublicKey, message, signature)
	fmt.Printf("Signature valid: %v\n", valid)

	// Test with wrong message (should fail)
	wrongMessage := []byte("Wrong message")
	wrongValid := VerifySignature(identity.PublicKey, wrongMessage, signature)
	fmt.Printf("Wrong message valid: %v\n", wrongValid)

	// Test with another identity (should fail)
	fmt.Printf("\n--- Testing with different identity ---\n")
	otherIdentity, _ := GenerateIdentity()
	otherValid := VerifySignature(otherIdentity.PublicKey, message, signature)
	fmt.Printf("Other identity verification: %v\n", otherValid)

	fmt.Printf("\n--- Key Security Properties ---\n")
	fmt.Printf("✓ Ed25519 provides strong cryptographic security\n")
	fmt.Printf("✓ Private key never leaves this device\n")
	fmt.Printf("✓ Public key can be safely shared\n")
	fmt.Printf("✓ Signatures prove message authenticity\n")
	fmt.Printf("✓ Cannot forge signatures without private key\n")
}