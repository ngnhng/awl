package main

import (
	"fmt"
	"net"
	"strings"
	"sync"
)

// Simplified routing table for mesh VPN
type RoutingTable struct {
	routes map[string]string // IP -> PeerID
	peers  map[string]PeerInfo // PeerID -> PeerInfo
	mutex  sync.RWMutex
}

type PeerInfo struct {
	ID       string
	Name     string
	Address  string
	Status   string
}

func NewRoutingTable() *RoutingTable {
	return &RoutingTable{
		routes: make(map[string]string),
		peers:  make(map[string]PeerInfo),
	}
}

func (rt *RoutingTable) AddPeer(info PeerInfo) {
	rt.mutex.Lock()
	defer rt.mutex.Unlock()
	rt.peers[info.ID] = info
	fmt.Printf("Added peer: %s (%s)\n", info.Name, info.ID)
}

func (rt *RoutingTable) AddRoute(ip, peerID string) error {
	rt.mutex.Lock()
	defer rt.mutex.Unlock()
	
	// Check if peer exists
	if _, exists := rt.peers[peerID]; !exists {
		return fmt.Errorf("peer %s not found", peerID)
	}
	
	rt.routes[ip] = peerID
	peer := rt.peers[peerID]
	fmt.Printf("Added route: %s -> %s (%s)\n", ip, peer.Name, peerID)
	return nil
}

func (rt *RoutingTable) GetPeer(ip string) (PeerInfo, bool) {
	rt.mutex.RLock()
	defer rt.mutex.RUnlock()
	
	peerID, exists := rt.routes[ip]
	if !exists {
		return PeerInfo{}, false
	}
	
	peer, exists := rt.peers[peerID]
	return peer, exists
}

func (rt *RoutingTable) ListRoutes() {
	rt.mutex.RLock()
	defer rt.mutex.RUnlock()
	
	fmt.Println("\n=== Routing Table ===")
	fmt.Printf("%-15s | %-20s | %-15s | %s\n", "IP Address", "Peer Name", "Peer ID", "Status")
	fmt.Println(strings.Repeat("-", 70))
	
	for ip, peerID := range rt.routes {
		if peer, exists := rt.peers[peerID]; exists {
			fmt.Printf("%-15s | %-20s | %-15s | %s\n", 
				ip, peer.Name, peerID, peer.Status)
		}
	}
	fmt.Println()
}

func (rt *RoutingTable) ListPeers() {
	rt.mutex.RLock()
	defer rt.mutex.RUnlock()
	
	fmt.Println("\n=== Known Peers ===")
	fmt.Printf("%-15s | %-20s | %-25s | %s\n", "Peer ID", "Name", "Address", "Status")
	fmt.Println(strings.Repeat("-", 80))
	
	for _, peer := range rt.peers {
		fmt.Printf("%-15s | %-20s | %-25s | %s\n", 
			peer.ID, peer.Name, peer.Address, peer.Status)
	}
	fmt.Println()
}

// Packet represents a simplified IP packet
type Packet struct {
	SrcIP    net.IP
	DstIP    net.IP
	Protocol byte
	Data     []byte
}

func parsePacket(data []byte) (*Packet, error) {
	if len(data) < 20 {
		return nil, fmt.Errorf("packet too short")
	}
	
	return &Packet{
		SrcIP:    net.IP(data[12:16]),
		DstIP:    net.IP(data[16:20]),
		Protocol: data[9],
		Data:     data,
	}, nil
}

// Simplified packet processing pipeline
func processVPNPacket(packetData []byte, routingTable *RoutingTable) {
	packet, err := parsePacket(packetData)
	if err != nil {
		fmt.Printf("Error parsing packet: %v\n", err)
		return
	}
	
	dstIP := packet.DstIP.String()
	
	fmt.Printf("\n--- Processing Packet ---\n")
	fmt.Printf("Source IP: %s\n", packet.SrcIP)
	fmt.Printf("Destination IP: %s\n", dstIP)
	fmt.Printf("Protocol: %d\n", packet.Protocol)
	fmt.Printf("Size: %d bytes\n", len(packet.Data))
	
	// Look up destination peer
	peer, exists := routingTable.GetPeer(dstIP)
	if !exists {
		fmt.Printf("❌ No route to %s - dropping packet\n", dstIP)
		showRoutingSuggestions(dstIP, routingTable)
		return
	}
	
	if peer.Status != "connected" {
		fmt.Printf("⚠️  Peer %s is %s - queueing packet\n", peer.Name, peer.Status)
		return
	}
	
	fmt.Printf("✅ Routing packet to %s via peer %s (%s)\n", 
		dstIP, peer.Name, peer.ID)
	
	// In real implementation:
	// 1. Find P2P connection to peer
	// 2. Send packet over the connection  
	// 3. Handle connection errors
	// 4. Implement retry logic
	fmt.Printf("   -> Sending to peer at %s\n", peer.Address)
	fmt.Printf("   -> Connection status: %s\n", peer.Status)
}

func showRoutingSuggestions(ip string, rt *RoutingTable) {
	fmt.Printf("💡 Suggestions:\n")
	fmt.Printf("   - Add peer that owns network containing %s\n", ip)
	fmt.Printf("   - Check if %s should be routed through existing peer\n", ip)
	
	// Show available routes for reference
	rt.mutex.RLock()
	if len(rt.routes) > 0 {
		fmt.Printf("   - Available routes: ")
		for routeIP := range rt.routes {
			fmt.Printf("%s ", routeIP)
		}
		fmt.Println()
	}
	rt.mutex.RUnlock()
}

func main() {
	fmt.Println("=== AWL Tutorial: Packet Routing ===\n")
	
	rt := NewRoutingTable()

	// Add some example peers
	peers := []PeerInfo{
		{ID: "peer-alice", Name: "Alice's Computer", Address: "192.168.1.100:9001", Status: "connected"},
		{ID: "peer-bob", Name: "Bob's Laptop", Address: "192.168.1.101:9001", Status: "connected"},
		{ID: "peer-charlie", Name: "Charlie's Phone", Address: "10.0.0.50:9001", Status: "connecting"},
		{ID: "peer-david", Name: "David's Server", Address: "203.0.113.1:9001", Status: "disconnected"},
	}
	
	for _, peer := range peers {
		rt.AddPeer(peer)
	}

	// Add some routes (IP assignments for each peer)
	routes := []struct {
		ip     string
		peerID string
	}{
		{"10.66.0.2", "peer-alice"},    // Alice gets 10.66.0.2
		{"10.66.0.3", "peer-bob"},      // Bob gets 10.66.0.3
		{"10.66.0.4", "peer-charlie"},  // Charlie gets 10.66.0.4
		{"10.66.0.5", "peer-david"},    // David gets 10.66.0.5
	}
	
	fmt.Println()
	for _, route := range routes {
		rt.AddRoute(route.ip, route.peerID)
	}

	// Show current state
	rt.ListPeers()
	rt.ListRoutes()

	// Simulate some packets
	fmt.Println("=== Simulating Packet Processing ===")
	
	testPackets := [][]byte{
		createFakePacket("10.66.0.1", "10.66.0.2"), // To Alice (connected)
		createFakePacket("10.66.0.1", "10.66.0.3"), // To Bob (connected)
		createFakePacket("10.66.0.1", "10.66.0.4"), // To Charlie (connecting)
		createFakePacket("10.66.0.1", "10.66.0.5"), // To David (disconnected)
		createFakePacket("10.66.0.1", "10.66.0.6"), // No route
		createFakePacket("10.66.0.1", "8.8.8.8"),   // Internet (no route)
	}

	for i, packet := range testPackets {
		fmt.Printf("\n--- Test Packet %d ---", i+1)
		processVPNPacket(packet, rt)
	}
	
	fmt.Println("\n=== Routing Concepts Demonstrated ===")
	fmt.Printf("✓ Peer registration and management\n")
	fmt.Printf("✓ IP to peer mapping (routing table)\n")
	fmt.Printf("✓ Packet destination lookup\n") 
	fmt.Printf("✓ Connection status checking\n")
	fmt.Printf("✓ Error handling for unknown destinations\n")
	fmt.Printf("✓ Graceful handling of disconnected peers\n")
}

// Create a minimal fake IP packet with src and dst
func createFakePacket(src, dst string) []byte {
	packet := make([]byte, 20) // Minimal IP header
	
	// Set version (4) and header length (5 * 4 = 20 bytes)
	packet[0] = 0x45
	
	// Set protocol (1 = ICMP for simplicity)
	packet[9] = 1
	
	// Copy source IP
	srcIP := net.ParseIP(src).To4()
	copy(packet[12:16], srcIP)
	
	// Copy destination IP  
	dstIP := net.ParseIP(dst).To4()
	copy(packet[16:20], dstIP)
	
	return packet
}