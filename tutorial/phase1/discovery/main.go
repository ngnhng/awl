package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// Peer represents a network peer
type Peer struct {
	ID      string `json:"id"`
	Address string `json:"address"`
	Name    string `json:"name"`
}

// Registry holds known peers
type Registry struct {
	peers map[string]Peer
	mutex sync.RWMutex
}

func NewRegistry() *Registry {
	return &Registry{
		peers: make(map[string]Peer),
	}
}

func (r *Registry) RegisterPeer(peer Peer) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.peers[peer.ID] = peer
	fmt.Printf("Registered peer: %s (%s)\n", peer.Name, peer.Address)
}

func (r *Registry) GetPeers() []Peer {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	peers := make([]Peer, 0, len(r.peers))
	for _, peer := range r.peers {
		peers = append(peers, peer)
	}
	return peers
}

func main() {
	registry := NewRegistry()

	// API handlers
	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", 405)
			return
		}

		var peer Peer
		if err := json.NewDecoder(r.Body).Decode(&peer); err != nil {
			http.Error(w, "Invalid JSON", 400)
			return
		}

		registry.RegisterPeer(peer)
		w.WriteHeader(200)
	})

	http.HandleFunc("/peers", func(w http.ResponseWriter, r *http.Request) {
		peers := registry.GetPeers()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(peers)
	})

	fmt.Println("Bootstrap server starting on :8080")
	fmt.Println("Endpoints:")
	fmt.Println("  POST /register - Register a new peer")
	fmt.Println("  GET  /peers    - List all known peers")
	http.ListenAndServe(":8080", nil)
}