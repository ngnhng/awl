# AWL Mesh VPN Learning Guide for Novice Programmers

## Table of Contents

1. [Prerequisites and Setup](#prerequisites-and-setup)
2. [Foundational Concepts](#foundational-concepts)
3. [Step-by-Step Implementation Tutorial](#step-by-step-implementation-tutorial)
4. [Advanced Topics](#advanced-topics)
5. [Learning Resources](#learning-resources)
6. [Exercises and Projects](#exercises-and-projects)

---

## Prerequisites and Setup

### Programming Knowledge Required

**Essential (Must Know)**:
- Basic Go programming (variables, functions, structs, interfaces, goroutines, channels)
- Basic networking concepts (IP addresses, ports, TCP/UDP)
- Command line usage (basic terminal/shell commands)
- Git version control basics

**Helpful (Good to Know)**:
- HTTP/REST APIs
- JSON data format
- Basic cryptography concepts
- Operating system concepts (processes, file systems)

**Not Required (Will Learn)**:
- Peer-to-peer networking
- VPN internals
- Advanced cryptography
- Cross-platform development
- libp2p networking stack

### Development Environment Setup

1. **Install Go 1.24+**
   ```bash
   # Download from https://golang.org/dl/
   go version  # Should show 1.24+
   ```

2. **Install Git**
   ```bash
   git --version
   ```

3. **Clone the Repository**
   ```bash
   git clone https://github.com/anywherelan/awl.git
   cd awl
   ```

4. **Install Dependencies**
   ```bash
   go mod tidy
   ```

5. **Development Tools (Optional but Recommended)**
   ```bash
   # Install useful Go tools
   go install -a std
   go install honnef.co/go/tools/cmd/staticcheck@latest
   ```

### Learning Timeline Estimate

- **Beginner**: 4-6 weeks (studying part-time)
- **Intermediate**: 2-3 weeks
- **Advanced Go developer**: 1-2 weeks

---

## Foundational Concepts

### 1. What is a Mesh VPN?

A **mesh VPN** creates a virtual private network where each device can connect directly to every other device, forming a "mesh" topology.

#### Traditional VPN vs Mesh VPN

```
Traditional VPN (Hub-and-Spoke):
Device A ──┐
           ├─── Central Server ─── Internet
Device B ──┘

Mesh VPN (Peer-to-Peer):
Device A ←──→ Device B
    ↕           ↕
Device C ←──→ Device D
```

#### Benefits of Mesh VPN:
- **No single point of failure**: No central server required
- **Lower latency**: Direct peer-to-peer connections
- **Better privacy**: Traffic doesn't go through third-party servers
- **Scalability**: Grows organically with more peers

### 2. Peer-to-Peer (P2P) Networking

#### Core P2P Concepts

**Peer Discovery**: How devices find each other on the internet
```go
// Simplified peer discovery concept
type Peer struct {
    ID        string
    Address   string
    PublicKey []byte
}

func DiscoverPeers(bootstrapNodes []string) ([]Peer, error) {
    // Connect to bootstrap nodes
    // Query for known peers
    // Return list of available peers
}
```

**NAT Traversal**: Getting around Network Address Translation
- **STUN**: Simple Traversal of UDP through NATs
- **TURN**: Traversal Using Relays around NAT
- **Hole Punching**: Technique to establish direct connections

**Distributed Hash Table (DHT)**: Decentralized way to store and find data
```go
// DHT stores key-value pairs across multiple nodes
type DHT interface {
    Put(key, value string) error
    Get(key string) (string, error)
    FindPeer(peerID string) (Peer, error)
}
```

### 3. Virtual Network Interfaces (TUN/TAP)

#### What is TUN/TAP?

**TUN** (tunnel) and **TAP** (network tap) are virtual network interfaces:
- **TUN**: Works at Layer 3 (IP packets)
- **TAP**: Works at Layer 2 (Ethernet frames)

AWL uses TUN for simplicity and cross-platform compatibility.

#### How TUN Works
```
Application ──→ TUN Interface ──→ AWL Program ──→ P2P Network
           packet capture       packet processing    send to peer

Remote Peer ──→ AWL Program ──→ TUN Interface ──→ Application
              receive packet    packet injection    deliver locally
```

#### Basic TUN Interface in Go
```go
package main

import (
    "net"
    "golang.zx2c4.com/wireguard/tun"
)

func createTUN() {
    // Create TUN interface
    tunDev, err := tun.CreateTUN("awl0", 1500)
    if err != nil {
        panic(err)
    }
    defer tunDev.Close()

    // Read packets from TUN
    packet := make([]byte, 1500)
    for {
        n, err := tunDev.Read(packet)
        if err != nil {
            continue
        }
        
        // Process packet
        processPacket(packet[:n])
    }
}

func processPacket(packet []byte) {
    // Parse IP header
    // Determine destination
    // Route to appropriate peer
}
```

### 4. Cryptography and Security

#### Public Key Cryptography

AWL uses **Ed25519** for peer identity and authentication:

```go
import "crypto/ed25519"

// Generate key pair
publicKey, privateKey, err := ed25519.GenerateKey(nil)

// Sign message
message := []byte("hello")
signature := ed25519.Sign(privateKey, message)

// Verify signature
valid := ed25519.Verify(publicKey, message, signature)
```

#### Transport Layer Security (TLS)

All peer-to-peer communication is encrypted with TLS:
- **Encryption**: Protects data in transit
- **Authentication**: Verifies peer identity
- **Integrity**: Ensures data hasn't been tampered with

### 5. Go Concurrency Patterns

AWL heavily uses Go's concurrency features:

#### Goroutines and Channels
```go
// Event bus pattern (simplified)
type EventBus struct {
    subscribers map[string][]chan interface{}
}

func (eb *EventBus) Subscribe(event string) <-chan interface{} {
    ch := make(chan interface{}, 10)
    eb.subscribers[event] = append(eb.subscribers[event], ch)
    return ch
}

func (eb *EventBus) Publish(event string, data interface{}) {
    for _, ch := range eb.subscribers[event] {
        select {
        case ch <- data:
        default: // Don't block if channel is full
        }
    }
}
```

#### Context for Cancellation
```go
func runService(ctx context.Context) {
    ticker := time.NewTicker(time.Second)
    defer ticker.Stop()

    for {
        select {
        case <-ctx.Done():
            return // Service shutdown
        case <-ticker.C:
            // Do periodic work
            doWork()
        }
    }
}
```

---

## Step-by-Step Implementation Tutorial

### Phase 1: Basic P2P Connection (Week 1)

#### Goal: Create two programs that can find and connect to each other

**Step 1: Simple TCP P2P Connection**

Create `tutorial/phase1/simple-p2p/main.go`:

```go
package main

import (
    "bufio"
    "fmt"
    "net"
    "os"
    "strings"
)

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: go run main.go [server|client] [port|address:port]")
        os.Exit(1)
    }

    mode := os.Args[1]
    addr := os.Args[2]

    if mode == "server" {
        runServer(addr)
    } else {
        runClient(addr)
    }
}

func runServer(port string) {
    ln, err := net.Listen("tcp", ":"+port)
    if err != nil {
        panic(err)
    }
    defer ln.Close()

    fmt.Printf("Server listening on port %s\n", port)

    for {
        conn, err := ln.Accept()
        if err != nil {
            continue
        }
        go handleConnection(conn)
    }
}

func runClient(address string) {
    conn, err := net.Dial("tcp", address)
    if err != nil {
        panic(err)
    }
    defer conn.Close()

    fmt.Printf("Connected to %s\n", address)

    // Send messages
    go func() {
        scanner := bufio.NewScanner(os.Stdin)
        for scanner.Scan() {
            text := scanner.Text()
            conn.Write([]byte(text + "\n"))
        }
    }()

    // Receive messages
    scanner := bufio.NewScanner(conn)
    for scanner.Scan() {
        fmt.Printf("Received: %s\n", scanner.Text())
    }
}

func handleConnection(conn net.Conn) {
    defer conn.Close()
    
    // Echo server
    scanner := bufio.NewScanner(conn)
    for scanner.Scan() {
        message := scanner.Text()
        fmt.Printf("Received: %s\n", message)
        
        // Echo back with prefix
        response := "Echo: " + message + "\n"
        conn.Write([]byte(response))
    }
}
```

**Exercise 1**: Run this program:
```bash
# Terminal 1
go run main.go server 8080

# Terminal 2  
go run main.go client localhost:8080
```

Type messages in the client terminal and see them echoed back.

**Step 2: Add Peer Discovery with Bootstrap Node**

Create `tutorial/phase1/discovery/main.go`:

```go
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
    http.ListenAndServe(":8080", nil)
}
```

**Step 3: Create a Peer Client**

Create `tutorial/phase1/peer/main.go`:

```go
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
        os.Exit(1)
    }

    peer := Peer{
        ID:      os.Args[1],
        Address: os.Args[2],
        Name:    os.Args[3],
    }

    // Register with bootstrap server
    registerPeer(peer)

    // Periodically discover other peers
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

    fmt.Printf("Discovered %d peers:\n", len(peers))
    for _, peer := range peers {
        fmt.Printf("  - %s (%s) at %s\n", peer.Name, peer.ID, peer.Address)
    }
}
```

**Exercise 2**: Test peer discovery:
```bash
# Terminal 1: Start bootstrap server
cd tutorial/phase1/discovery
go run main.go

# Terminal 2: Start first peer
cd tutorial/phase1/peer
go run main.go peer1 192.168.1.100:9001 "Alice's Computer"

# Terminal 3: Start second peer  
cd tutorial/phase1/peer
go run main.go peer2 192.168.1.101:9001 "Bob's Laptop"
```

### Phase 2: Add Cryptography (Week 2)

#### Goal: Secure the peer-to-peer connections with cryptography

**Step 1: Peer Identity with Ed25519**

Create `tutorial/phase2/identity/main.go`:

```go
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
    // Generate identity
    identity, err := GenerateIdentity()
    if err != nil {
        panic(err)
    }

    fmt.Printf("Generated peer identity:\n")
    fmt.Printf("ID: %s\n", identity.ID[:20]+"...") // Show first 20 chars
    fmt.Printf("Public Key: %x\n", identity.PublicKey[:8])  // Show first 8 bytes

    // Test signing
    message := []byte("Hello, mesh VPN world!")
    signature := identity.Sign(message)

    fmt.Printf("\nMessage: %s\n", message)
    fmt.Printf("Signature: %x...\n", signature[:8]) // Show first 8 bytes

    // Verify signature
    valid := VerifySignature(identity.PublicKey, message, signature)
    fmt.Printf("Signature valid: %v\n", valid)

    // Test with wrong message
    wrongMessage := []byte("Wrong message")
    wrongValid := VerifySignature(identity.PublicKey, wrongMessage, signature)
    fmt.Printf("Wrong message valid: %v\n", wrongValid)
}
```

**Step 2: Secure Authentication Protocol**

Create `tutorial/phase2/auth/main.go`:

```go
package main

import (
    "crypto/ed25519"
    "crypto/rand"
    "encoding/json"
    "fmt"
    "time"
)

// Challenge-Response Authentication
type AuthChallenge struct {
    Challenge []byte    `json:"challenge"`
    Timestamp time.Time `json:"timestamp"`
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
    }
}

func (pi *PeerIdentity) RespondToChallenge(challenge AuthChallenge) AuthResponse {
    // Sign the challenge
    data, _ := json.Marshal(challenge)
    signature := ed25519.Sign(pi.PrivateKey, data)

    return AuthResponse{
        PeerID:    pi.ID,
        Signature: signature,
    }
}

func VerifyResponse(challenge AuthChallenge, response AuthResponse, publicKey ed25519.PublicKey) bool {
    // Check timestamp (prevent replay attacks)
    if time.Since(challenge.Timestamp) > 30*time.Second {
        return false
    }

    // Verify signature
    data, _ := json.Marshal(challenge)
    return ed25519.Verify(publicKey, data, response.Signature)
}

func main() {
    // Create two peer identities
    alice, _ := generateTestIdentity("alice")
    bob, _ := generateTestIdentity("bob")

    fmt.Println("=== Authentication Protocol Demo ===\n")

    // 1. Alice generates a challenge for Bob
    challenge := GenerateChallenge()
    fmt.Printf("1. Alice generates challenge: %x...\n", challenge.Challenge[:8])

    // 2. Bob responds to the challenge
    response := bob.RespondToChallenge(challenge)
    fmt.Printf("2. Bob responds with signature: %x...\n", response.Signature[:8])

    // 3. Alice verifies Bob's response
    valid := VerifyResponse(challenge, response, bob.PublicKey)
    fmt.Printf("3. Alice verifies Bob's response: %v\n", valid)

    // 4. Test with wrong key
    fmt.Printf("\n--- Testing with wrong key ---\n")
    wrongValid := VerifyResponse(challenge, response, alice.PublicKey)
    fmt.Printf("Wrong key verification: %v\n", wrongValid)

    // 5. Test with old challenge (replay attack)
    fmt.Printf("\n--- Testing replay attack ---\n")
    oldChallenge := AuthChallenge{
        Challenge: challenge.Challenge,
        Timestamp: time.Now().Add(-1 * time.Minute), // Old timestamp
    }
    oldResponse := bob.RespondToChallenge(oldChallenge)
    replayValid := VerifyResponse(oldChallenge, oldResponse, bob.PublicKey)
    fmt.Printf("Replay attack blocked: %v\n", !replayValid)
}

func generateTestIdentity(name string) (*PeerIdentity, error) {
    pub, priv, err := ed25519.GenerateKey(rand.Reader)
    if err != nil {
        return nil, err
    }

    return &PeerIdentity{
        PublicKey:  pub,
        PrivateKey: priv,
        ID:         name + "-test-id",
    }, nil
}
```

**Exercise 3**: Run the authentication demo:
```bash
cd tutorial/phase2/auth
go run main.go
```

### Phase 3: Basic VPN Functionality (Week 3)

#### Goal: Create a simple TUN interface and route packets

**Step 1: Create TUN Interface**

Create `tutorial/phase3/tun/main.go`:

```go
package main

import (
    "fmt"
    "net"
    "os"
    "os/signal"
    "syscall"

    "golang.zx2c4.com/wireguard/tun"
)

func main() {
    // Create TUN interface
    tunDevice, err := tun.CreateTUN("awl-tutorial", 1500)
    if err != nil {
        fmt.Printf("Error creating TUN: %v\n", err)
        fmt.Println("Note: This requires admin/root privileges")
        os.Exit(1)
    }
    defer tunDevice.Close()

    fmt.Println("Created TUN interface: awl-tutorial")
    fmt.Println("Run this command to configure it:")
    fmt.Println("sudo ip addr add 10.66.0.1/24 dev awl-tutorial")
    fmt.Println("sudo ip link set awl-tutorial up")
    fmt.Println("\nPress Ctrl+C to stop...")

    // Handle graceful shutdown
    sigCh := make(chan os.Signal, 1)
    signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

    // Read packets in a goroutine
    go readPackets(tunDevice)

    // Wait for signal
    <-sigCh
    fmt.Println("\nShutting down...")
}

func readPackets(tunDevice tun.Device) {
    packet := make([]byte, 1500)
    
    for {
        n, err := tunDevice.Read(packet[:], 0)
        if err != nil {
            fmt.Printf("Error reading packet: %v\n", err)
            continue
        }

        if n > 0 {
            analyzePacket(packet[:n])
        }
    }
}

func analyzePacket(packet []byte) {
    if len(packet) < 20 {
        return // Too short for IP header
    }

    // Parse basic IP header
    version := packet[0] >> 4
    headerLen := int(packet[0]&0x0F) * 4
    protocol := packet[9]
    srcIP := net.IP(packet[12:16])
    dstIP := net.IP(packet[16:20])

    fmt.Printf("Packet: IPv%d, %s -> %s, Protocol: %d, Size: %d bytes\n",
        version, srcIP, dstIP, protocol, len(packet))

    // In a real mesh VPN, we would:
    // 1. Look up the destination peer for dstIP
    // 2. Forward the packet to that peer
    // 3. If no peer found, drop the packet
}
```

**Exercise 4**: Test TUN interface:
```bash
# Run the program (needs sudo/admin privileges)
sudo go run tutorial/phase3/tun/main.go

# In another terminal, configure the interface
sudo ip addr add 10.66.0.1/24 dev awl-tutorial
sudo ip link set awl-tutorial up

# Test with ping
ping 10.66.0.2
```

**Step 2: Packet Routing Logic**

Create `tutorial/phase3/routing/main.go`:

```go
package main

import (
    "fmt"
    "net"
    "sync"
)

// Simplified routing table
type RoutingTable struct {
    routes map[string]string // IP -> PeerID
    mutex  sync.RWMutex
}

func NewRoutingTable() *RoutingTable {
    return &RoutingTable{
        routes: make(map[string]string),
    }
}

func (rt *RoutingTable) AddRoute(ip, peerID string) {
    rt.mutex.Lock()
    defer rt.mutex.Unlock()
    rt.routes[ip] = peerID
    fmt.Printf("Added route: %s -> %s\n", ip, peerID)
}

func (rt *RoutingTable) GetPeer(ip string) (string, bool) {
    rt.mutex.RLock()
    defer rt.mutex.RUnlock()
    peerID, exists := rt.routes[ip]
    return peerID, exists
}

func (rt *RoutingTable) ListRoutes() {
    rt.mutex.RLock()
    defer rt.mutex.RUnlock()
    
    fmt.Println("Routing Table:")
    for ip, peerID := range rt.routes {
        fmt.Printf("  %s -> %s\n", ip, peerID)
    }
}

// Simplified packet processing
func processVPNPacket(packet []byte, routingTable *RoutingTable) {
    if len(packet) < 20 {
        return
    }

    // Extract destination IP
    dstIP := net.IP(packet[16:20]).String()
    
    // Look up peer
    peerID, exists := routingTable.GetPeer(dstIP)
    if !exists {
        fmt.Printf("No route to %s - dropping packet\n", dstIP)
        return
    }

    fmt.Printf("Routing packet to %s via peer %s\n", dstIP, peerID)
    
    // In real implementation:
    // 1. Find P2P connection to peer
    // 2. Send packet over the connection
    // 3. Handle connection errors
}

func main() {
    rt := NewRoutingTable()

    // Add some example routes
    rt.AddRoute("10.66.0.2", "peer-alice")
    rt.AddRoute("10.66.0.3", "peer-bob") 
    rt.AddRoute("10.66.0.4", "peer-charlie")

    fmt.Println()
    rt.ListRoutes()

    // Simulate some packets
    fmt.Println("\n=== Simulating Packet Processing ===")
    
    // Create fake packets (simplified - just destination IP)
    testPackets := [][]byte{
        createFakePacket("10.66.0.1", "10.66.0.2"),
        createFakePacket("10.66.0.1", "10.66.0.3"),
        createFakePacket("10.66.0.1", "10.66.0.5"), // No route
    }

    for _, packet := range testPackets {
        processVPNPacket(packet, rt)
    }
}

// Create a minimal fake IP packet with src and dst
func createFakePacket(src, dst string) []byte {
    packet := make([]byte, 20) // Minimal IP header
    
    // Set version (4) and header length (5 * 4 = 20 bytes)
    packet[0] = 0x45
    
    // Copy source IP
    srcIP := net.ParseIP(src).To4()
    copy(packet[12:16], srcIP)
    
    // Copy destination IP  
    dstIP := net.ParseIP(dst).To4()
    copy(packet[16:20], dstIP)
    
    return packet
}
```

### Phase 4: Integration with libp2p (Week 4)

#### Goal: Replace our simple networking with libp2p

**Step 1: Basic libp2p Host**

Create `tutorial/phase4/libp2p-basic/main.go`:

```go
package main

import (
    "context"
    "fmt"
    "os"
    "os/signal"
    "syscall"

    "github.com/libp2p/go-libp2p"
    "github.com/libp2p/go-libp2p/core/host"
    "github.com/libp2p/go-libp2p/core/network"
    "github.com/libp2p/go-libp2p/core/peer"
    "github.com/libp2p/go-libp2p/core/protocol"
    "github.com/multiformats/go-multiaddr"
)

const ProtocolID = "/awl-tutorial/1.0.0"

func main() {
    ctx := context.Background()

    // Create libp2p host
    h, err := libp2p.New(
        libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/0"),
    )
    if err != nil {
        panic(err)
    }
    defer h.Close()

    // Set stream handler
    h.SetStreamHandler(protocol.ID(ProtocolID), handleStream)

    // Print host information
    fmt.Printf("Host ID: %s\n", h.ID())
    fmt.Printf("Addresses:\n")
    for _, addr := range h.Addrs() {
        fmt.Printf("  %s/p2p/%s\n", addr, h.ID())
    }

    // If peer address provided, connect to it
    if len(os.Args) > 1 {
        peerAddr := os.Args[1]
        connectToPeer(ctx, h, peerAddr)
    }

    fmt.Println("\nWaiting for connections... Press Ctrl+C to exit")

    // Wait for signal
    sigCh := make(chan os.Signal, 1)
    signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
    <-sigCh
}

func handleStream(s network.Stream) {
    defer s.Close()
    
    remotePeer := s.Conn().RemotePeer()
    fmt.Printf("New stream from peer: %s\n", remotePeer)

    // Simple echo server
    buf := make([]byte, 1024)
    for {
        n, err := s.Read(buf)
        if err != nil {
            break
        }
        
        message := string(buf[:n])
        fmt.Printf("Received: %s", message)
        
        // Echo back
        s.Write([]byte("Echo: " + message))
    }
}

func connectToPeer(ctx context.Context, h host.Host, peerAddr string) {
    // Parse peer address
    maddr, err := multiaddr.NewMultiaddr(peerAddr)
    if err != nil {
        fmt.Printf("Error parsing address: %v\n", err)
        return
    }

    // Extract peer info
    peerInfo, err := peer.AddrInfoFromP2pAddr(maddr)
    if err != nil {
        fmt.Printf("Error extracting peer info: %v\n", err)
        return
    }

    // Connect to peer
    err = h.Connect(ctx, *peerInfo)
    if err != nil {
        fmt.Printf("Error connecting to peer: %v\n", err)
        return
    }

    fmt.Printf("Connected to peer: %s\n", peerInfo.ID)

    // Open stream
    s, err := h.NewStream(ctx, peerInfo.ID, protocol.ID(ProtocolID))
    if err != nil {
        fmt.Printf("Error opening stream: %v\n", err)
        return
    }
    defer s.Close()

    // Send test message
    message := "Hello from peer!\n"
    s.Write([]byte(message))

    // Read response
    buf := make([]byte, 1024)
    n, _ := s.Read(buf)
    fmt.Printf("Response: %s", buf[:n])
}
```

**Exercise 5**: Test libp2p connection:
```bash
# Terminal 1: Start first peer
cd tutorial/phase4/libp2p-basic
go mod init tutorial-libp2p
go mod tidy
go run main.go

# Copy one of the addresses from the output

# Terminal 2: Connect second peer to first
go run main.go /ip4/127.0.0.1/tcp/XXXX/p2p/PEER_ID
```

This covers the first 4 weeks of learning. The guide would continue with more advanced topics like DHT, NAT traversal, complete VPN implementation, cross-platform considerations, security hardening, etc.

---

## Advanced Topics

### 1. Production Considerations
- Error handling and recovery
- Logging and monitoring
- Performance optimization
- Memory management
- Connection pooling

### 2. Security Hardening
- Forward secrecy
- Peer verification
- Attack mitigation
- Security auditing
- Key rotation

### 3. Cross-Platform Development
- Platform-specific networking code
- Build systems and packaging
- Testing across platforms
- Distribution and updates

### 4. Network Protocols Deep Dive
- QUIC vs TCP trade-offs
- NAT traversal strategies
- Relay optimization
- Quality of Service (QoS)
- Network topology optimization

---

## Learning Resources

### Books
- "Computer Networks" by Andrew Tanenbaum
- "Cryptography Engineering" by Ferguson, Schneier, and Kohno
- "The Go Programming Language" by Donovan and Kernighan

### Online Courses
- [Computer Networking Course (Coursera)](https://www.coursera.org/learn/computer-networking)
- [Practical Networking with Go](https://www.youtube.com/playlist?list=PLy_6D98if3ULEtXtNSY_2qN21VCKgoQAE)

### Documentation
- [libp2p Documentation](https://docs.libp2p.io/)
- [WireGuard Protocol](https://www.wireguard.com/protocol/)
- [Go Concurrency Patterns](https://go.dev/blog/pipelines)

### Tools for Learning
- [Wireshark](https://www.wireshark.org/) - Network packet analyzer
- [tcpdump](https://www.tcpdump.org/) - Command-line packet analyzer
- [netstat/ss](https://linux.die.net/man/8/netstat) - Network statistics

---

## Exercises and Projects

### Beginner Projects
1. **Simple Chat Application**: Build P2P chat using the tutorial concepts
2. **File Transfer Tool**: Send files directly between peers
3. **Network Monitor**: Visualize your network connections

### Intermediate Projects  
1. **Mini VPN**: Implement a simple VPN with 2-3 peers
2. **NAT Traversal Demo**: Show how hole punching works
3. **DHT Implementation**: Build a simple distributed hash table

### Advanced Projects
1. **Full Mesh VPN**: Implement most AWL features
2. **Performance Benchmarks**: Compare different transport protocols
3. **Security Analysis**: Audit the cryptographic implementation

### Contributing to AWL
1. **Documentation**: Improve existing docs
2. **Testing**: Add more test cases
3. **Features**: Implement new functionality
4. **Bug Fixes**: Find and fix issues

This guide provides a structured learning path from basic networking concepts to advanced mesh VPN implementation, with hands-on exercises and real code examples that build understanding progressively.