package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

type Peer struct {
	ID      string `json:"id"`
	Address string `json:"address"`
	Name    string `json:"name"`
}

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: go run main.go <id> <address> <name>")
		fmt.Println("Example: go run main.go peer1 192.168.1.100:9001 \"Alice's Computer\"")
		os.Exit(1)
	}

	peer := Peer{
		ID:      os.Args[1],
		Address: os.Args[2],
		Name:    os.Args[3],
	}

	// Register with bootstrap server
	fmt.Printf("Registering peer %s...\n", peer.Name)
	registerPeer(peer)

	// Periodically discover other peers
	fmt.Printf("Starting peer discovery (every 5 seconds)...\n")
	fmt.Printf("Press Ctrl+C to stop\n\n")
	
	for {
		time.Sleep(5 * time.Second)
		discoverPeers()
	}
}

func registerPeer(peer Peer) {
	data, _ := json.Marshal(peer)
	resp, err := http.Post("http://localhost:8080/register", "application/json", bytes.NewBuffer(data))
	if err != nil {
		fmt.Printf("Error registering: %v\n", err)
		return
	}
	defer resp.Body.Close()
	
	if resp.StatusCode == 200 {
		fmt.Printf("Successfully registered as %s\n", peer.Name)
	} else {
		fmt.Printf("Registration failed with status: %d\n", resp.StatusCode)
	}
}

func discoverPeers() {
	resp, err := http.Get("http://localhost:8080/peers")
	if err != nil {
		fmt.Printf("Error discovering peers: %v\n", err)
		return
	}
	defer resp.Body.Close()

	var peers []Peer
	json.NewDecoder(resp.Body).Decode(&peers)

	fmt.Printf("--- Discovered %d peers ---\n", len(peers))
	for _, peer := range peers {
		fmt.Printf("  - %s (%s) at %s\n", peer.Name, peer.ID, peer.Address)
	}
	fmt.Println()
}