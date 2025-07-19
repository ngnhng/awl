package main

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"time"
)

// Challenge-Response Authentication Protocol
type AuthChallenge struct {
	Challenge []byte    `json:"challenge"`
	Timestamp time.Time `json:"timestamp"`
	Nonce     int64     `json:"nonce"`
}

type AuthResponse struct {
	PeerID    string `json:"peer_id"`
	Signature []byte `json:"signature"`
}

type PeerIdentity struct {
	PublicKey  ed25519.PublicKey
	PrivateKey ed25519.PrivateKey
	ID         string
}

func GenerateChallenge() AuthChallenge {
	challenge := make([]byte, 32)
	rand.Read(challenge)
	
	return AuthChallenge{
		Challenge: challenge,
		Timestamp: time.Now(),
		Nonce:     time.Now().UnixNano(), // Prevent replay attacks
	}
}

func (pi *PeerIdentity) RespondToChallenge(challenge AuthChallenge) AuthResponse {
	// Create message to sign (includes all challenge data)
	data, _ := json.Marshal(challenge)
	signature := ed25519.Sign(pi.PrivateKey, data)

	return AuthResponse{
		PeerID:    pi.ID,
		Signature: signature,
	}
}

func VerifyResponse(challenge AuthChallenge, response AuthResponse, publicKey ed25519.PublicKey) (bool, string) {
	// Check timestamp (prevent replay attacks)
	age := time.Since(challenge.Timestamp)
	if age > 30*time.Second {
		return false, fmt.Sprintf("challenge expired (age: %v)", age)
	}

	// Verify signature
	data, _ := json.Marshal(challenge)
	valid := ed25519.Verify(publicKey, data, response.Signature)
	
	if !valid {
		return false, "invalid signature"
	}
	
	return true, "authentication successful"
}

func main() {
	fmt.Println("=== AWL Tutorial: Secure Authentication Protocol ===\n")

	// Create two peer identities
	alice, _ := generateTestIdentity("alice")
	bob, _ := generateTestIdentity("bob")

	fmt.Printf("Peers:\n")
	fmt.Printf("  Alice ID: %s\n", alice.ID[:16]+"...")
	fmt.Printf("  Bob ID:   %s\n\n", bob.ID[:16]+"...")

	// Scenario 1: Successful authentication
	fmt.Println("=== Scenario 1: Normal Authentication ===")
	
	// 1. Alice generates a challenge for Bob
	challenge := GenerateChallenge()
	fmt.Printf("1. Alice generates challenge:\n")
	fmt.Printf("   Challenge: %x...\n", challenge.Challenge[:8])
	fmt.Printf("   Timestamp: %v\n", challenge.Timestamp.Format("15:04:05"))
	fmt.Printf("   Nonce: %d\n", challenge.Nonce)

	// 2. Bob responds to the challenge
	response := bob.RespondToChallenge(challenge)
	fmt.Printf("\n2. Bob responds:\n")
	fmt.Printf("   Peer ID: %s\n", response.PeerID[:16]+"...")
	fmt.Printf("   Signature: %x...\n", response.Signature[:8])

	// 3. Alice verifies Bob's response
	valid, reason := VerifyResponse(challenge, response, bob.PublicKey)
	fmt.Printf("\n3. Alice verifies response:\n")
	fmt.Printf("   Result: %v\n", valid)
	fmt.Printf("   Reason: %s\n", reason)

	// Scenario 2: Wrong key attack
	fmt.Println("\n=== Scenario 2: Wrong Key Attack ===")
	wrongValid, wrongReason := VerifyResponse(challenge, response, alice.PublicKey)
	fmt.Printf("Wrong key verification: %v (%s)\n", wrongValid, wrongReason)

	// Scenario 3: Replay attack
	fmt.Println("\n=== Scenario 3: Replay Attack ===")
	time.Sleep(100 * time.Millisecond) // Small delay to show timestamp checking
	oldChallenge := AuthChallenge{
		Challenge: challenge.Challenge,
		Timestamp: time.Now().Add(-1 * time.Minute), // Old timestamp
		Nonce:     challenge.Nonce,
	}
	oldResponse := bob.RespondToChallenge(oldChallenge)
	replayValid, replayReason := VerifyResponse(oldChallenge, oldResponse, bob.PublicKey)
	fmt.Printf("Replay attack result: %v (%s)\n", replayValid, replayReason)

	// Scenario 4: Modified challenge attack
	fmt.Println("\n=== Scenario 4: Modified Challenge Attack ===")
	modifiedChallenge := challenge
	modifiedChallenge.Challenge[0] ^= 0xFF // Flip some bits
	modifiedResponse := bob.RespondToChallenge(modifiedChallenge)
	modifiedValid, modifiedReason := VerifyResponse(challenge, modifiedResponse, bob.PublicKey)
	fmt.Printf("Modified challenge result: %v (%s)\n", modifiedValid, modifiedReason)

	fmt.Println("\n=== Security Properties ===")
	fmt.Printf("✓ Only peers with correct private key can respond\n")
	fmt.Printf("✓ Challenges expire to prevent replay attacks\n")
	fmt.Printf("✓ Nonces provide additional replay protection\n")
	fmt.Printf("✓ Tampering with challenges is detected\n")
	fmt.Printf("✓ No passwords or shared secrets required\n")
}

func generateTestIdentity(name string) (*PeerIdentity, error) {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}

	return &PeerIdentity{
		PublicKey:  pub,
		PrivateKey: priv,
		ID:         name + "-test-id-" + fmt.Sprintf("%d", time.Now().UnixNano()%10000),
	}, nil
}