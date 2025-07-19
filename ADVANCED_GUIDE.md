# Advanced Topics in Mesh VPN Development

This guide covers advanced concepts for building production-ready mesh VPN systems like AWL. These topics go beyond the basic tutorials and delve into real-world considerations.

## Table of Contents

1. [Production Architecture Patterns](#production-architecture-patterns)
2. [Advanced Networking](#advanced-networking)
3. [Security Hardening](#security-hardening)
4. [Performance Optimization](#performance-optimization)
5. [Cross-Platform Development](#cross-platform-development)
6. [Observability and Monitoring](#observability-and-monitoring)
7. [Deployment and Distribution](#deployment-and-distribution)
8. [Contributing to AWL](#contributing-to-awl)

---

## Production Architecture Patterns

### 1. Event-Driven Architecture

AWL uses an event bus pattern for loose coupling between components:

```go
// Event bus implementation pattern
type EventBus interface {
    Subscribe(eventType reflect.Type, handler func(interface{}))
    Publish(event interface{})
    Unsubscribe(eventType reflect.Type, handler func(interface{}))
}

// Usage in AWL
awlevent.WrapSubscriptionToCallback(a.ctx, func(_ interface{}) {
    a.Tunnel.RefreshPeersList()
}, a.Eventbus, new(awlevent.KnownPeerChanged))
```

**Benefits**:
- Decoupled components
- Easy to test individual parts
- Can add new features without modifying existing code
- Clear data flow

**Implementation Considerations**:
- Event ordering and consistency
- Error handling in event handlers
- Memory leaks from subscriptions
- Performance impact of event routing

### 2. Service-Oriented Design

Each major functionality is encapsulated in a service:

```go
type Application struct {
    // Core services
    P2p        *p2p.P2p           // Networking layer
    Tunnel     *service.Tunnel    // VPN routing
    AuthStatus *service.AuthStatus // Peer authentication
    SOCKS5     *service.SOCKS5    // Proxy service
    Dns        *DNSService        // Domain resolution
    Api        *api.Handler       // Web interface
    
    // Shared infrastructure
    Conf       *config.Config     // Configuration
    Eventbus   awlevent.Bus      // Communication bus
    LogBuffer  *ringbuffer.RingBuffer // Logging
}
```

**Design Principles**:
- Single responsibility per service
- Dependency injection for testability
- Well-defined interfaces between services
- Graceful startup and shutdown

### 3. State Management

**Configuration State**:
```go
type Config struct {
    // Immutable after startup
    P2pNode    P2pNodeConfig    `json:"p2p_node"`
    VPNConfig  VPNConfig        `json:"vpn"`
    
    // Mutable runtime state
    KnownPeers []KnownPeer      `json:"known_peers"`
    mutex      sync.RWMutex     // Protect concurrent access
}

func (c *Config) AddPeer(peer KnownPeer) {
    c.mutex.Lock()
    defer c.mutex.Unlock()
    c.KnownPeers = append(c.KnownPeers, peer)
    c.save() // Persist to disk
}
```

**Runtime State**:
- Connection status per peer
- Network metrics and statistics  
- Active streams and connections
- VPN routing table

---

## Advanced Networking

### 1. NAT Traversal Strategies

**Hole Punching**:
```go
// Simplified hole punching workflow
func attemptHolePunch(localAddr, remoteAddr net.Addr) error {
    // 1. Both peers connect to STUN server
    localPublicAddr := getPublicAddr(stunServer)
    
    // 2. Exchange addresses through signaling
    exchangeAddresses(localPublicAddr, remoteAddr)
    
    // 3. Simultaneously attempt connection
    return connectSimultaneously(localPublicAddr, remoteAddr)
}
```

**TURN Relay Fallback**:
```go
func connectWithRelay(peerID peer.ID) (net.Conn, error) {
    // Direct connection failed, use relay
    relayConn, err := dialThroughRelay(relayServer, peerID)
    if err != nil {
        return nil, fmt.Errorf("relay connection failed: %w", err)
    }
    
    // Still encrypted end-to-end
    return wrapWithTLS(relayConn, peerID), nil
}
```

### 2. Transport Protocol Selection

AWL automatically chooses between QUIC and TCP:

```go
func selectTransport(peerInfo *peer.AddrInfo) transport.Transport {
    // Prefer QUIC for better performance
    if supportsQUIC(peerInfo) {
        return &quicTransport{
            congestionControl: "bbr",
            maxStreams:       100,
        }
    }
    
    // Fallback to TCP for compatibility
    return &tcpTransport{
        keepAlive: 30 * time.Second,
        nodelay:   true,
    }
}
```

**QUIC Advantages**:
- Built-in encryption (TLS 1.3)
- Multiple streams over single connection
- Better congestion control
- Faster connection establishment
- Connection migration support

**TCP Fallback**:
- Universal compatibility
- Well-tested infrastructure
- Firewall/proxy friendly
- Simpler troubleshooting

### 3. Advanced Routing

**Mesh Topology Management**:
```go
type MeshTopology struct {
    nodes map[peer.ID]*Node
    edges map[peer.ID][]peer.ID
    
    // Routing table with multiple paths
    routes map[string][]Route
}

type Route struct {
    NextHop peer.ID
    Cost    int
    Latency time.Duration
    Hops    []peer.ID
}

func (m *MeshTopology) findBestRoute(dstIP string) Route {
    routes := m.routes[dstIP]
    if len(routes) == 0 {
        return Route{} // No route
    }
    
    // Choose based on latency, cost, or load balancing
    return selectRoute(routes, m.currentMetrics())
}
```

**Quality of Service (QoS)**:
```go
type QoSManager struct {
    priorityQueues map[Priority]*PacketQueue
    bandwidthLimit int64
    currentUsage   int64
}

func (q *QoSManager) classifyPacket(packet []byte) Priority {
    // Parse packet to determine priority
    if isVoIP(packet) {
        return HighPriority
    } else if isFileTransfer(packet) {
        return LowPriority
    }
    return NormalPriority
}
```

---

## Security Hardening

### 1. Cryptographic Best Practices

**Forward Secrecy**:
```go
// Rotate session keys periodically
type SessionKeyManager struct {
    currentKey  []byte
    nextKey     []byte
    rotateTime  time.Time
    rotateEvery time.Duration
}

func (skm *SessionKeyManager) rotateIfNeeded() {
    if time.Now().After(skm.rotateTime) {
        skm.currentKey = skm.nextKey
        skm.nextKey = generateNewKey()
        skm.rotateTime = time.Now().Add(skm.rotateEvery)
    }
}
```

**Key Derivation**:
```go
import "golang.org/x/crypto/hkdf"

func deriveKeys(sharedSecret []byte) (encKey, macKey []byte) {
    // Use HKDF for key derivation
    hkdf := hkdf.New(sha256.New, sharedSecret, nil, []byte("awl-vpn-v1"))
    
    encKey = make([]byte, 32) // AES-256
    macKey = make([]byte, 32) // HMAC-SHA256
    
    hkdf.Read(encKey)
    hkdf.Read(macKey)
    
    return encKey, macKey
}
```

### 2. Attack Mitigation

**Rate Limiting**:
```go
type RateLimiter struct {
    requests map[peer.ID]*TokenBucket
    mutex    sync.RWMutex
}

func (rl *RateLimiter) Allow(peerID peer.ID) bool {
    rl.mutex.RLock()
    bucket := rl.requests[peerID]
    rl.mutex.RUnlock()
    
    if bucket == nil {
        bucket = NewTokenBucket(10, time.Second) // 10 requests per second
        rl.mutex.Lock()
        rl.requests[peerID] = bucket
        rl.mutex.Unlock()
    }
    
    return bucket.TryConsume(1)
}
```

**Connection Validation**:
```go
func validateConnection(conn net.Conn, expectedPeer peer.ID) error {
    // Timeout for handshake
    conn.SetDeadline(time.Now().Add(30 * time.Second))
    defer conn.SetDeadline(time.Time{})
    
    // Verify peer identity
    actualPeer, err := performHandshake(conn)
    if err != nil {
        return fmt.Errorf("handshake failed: %w", err)
    }
    
    if actualPeer != expectedPeer {
        return fmt.Errorf("peer identity mismatch")
    }
    
    return nil
}
```

### 3. Network Isolation

**VPN Interface Security**:
```go
func createSecureInterface(name string) error {
    // Create interface with restricted permissions
    iface, err := createTUN(name)
    if err != nil {
        return err
    }
    
    // Configure firewall rules
    rules := []string{
        "DROP all traffic to host network",
        "ALLOW only VPN subnet traffic",
        "LOG suspicious activities",
    }
    
    return applyFirewallRules(iface, rules)
}
```

---

## Performance Optimization

### 1. Memory Management

**Object Pooling**:
```go
var packetPool = sync.Pool{
    New: func() interface{} {
        return make([]byte, 1500) // MTU size
    },
}

func processPacket(data []byte) {
    // Get packet buffer from pool
    packet := packetPool.Get().([]byte)
    defer packetPool.Put(packet) // Return to pool
    
    // Process packet...
    copy(packet, data)
    routePacket(packet[:len(data)])
}
```

**Ring Buffers for Logging**:
```go
type RingBuffer struct {
    buffer []LogEntry
    head   int
    tail   int
    size   int
    mutex  sync.RWMutex
}

func (rb *RingBuffer) Write(entry LogEntry) {
    rb.mutex.Lock()
    defer rb.mutex.Unlock()
    
    rb.buffer[rb.head] = entry
    rb.head = (rb.head + 1) % len(rb.buffer)
    
    if rb.size < len(rb.buffer) {
        rb.size++
    } else {
        rb.tail = (rb.tail + 1) % len(rb.buffer)
    }
}
```

### 2. Concurrent Processing

**Pipeline Pattern**:
```go
func startPacketPipeline(ctx context.Context) {
    packetCh := make(chan []byte, 1000)
    processedCh := make(chan ProcessedPacket, 1000)
    
    // Stage 1: Packet capture
    go capturePackets(ctx, packetCh)
    
    // Stage 2: Packet processing (multiple workers)
    for i := 0; i < runtime.NumCPU(); i++ {
        go processPackets(ctx, packetCh, processedCh)
    }
    
    // Stage 3: Packet routing
    go routePackets(ctx, processedCh)
}
```

**Worker Pool**:
```go
type WorkerPool struct {
    workers   int
    taskQueue chan Task
    wg        sync.WaitGroup
}

func (wp *WorkerPool) Start(ctx context.Context) {
    for i := 0; i < wp.workers; i++ {
        wp.wg.Add(1)
        go wp.worker(ctx)
    }
}

func (wp *WorkerPool) worker(ctx context.Context) {
    defer wp.wg.Done()
    
    for {
        select {
        case task := <-wp.taskQueue:
            task.Execute()
        case <-ctx.Done():
            return
        }
    }
}
```

### 3. Network Optimizations

**Connection Pooling**:
```go
type ConnectionPool struct {
    pools map[peer.ID]*ConnPool
    mutex sync.RWMutex
}

type ConnPool struct {
    connections chan net.Conn
    maxSize     int
}

func (cp *ConnectionPool) GetConnection(peerID peer.ID) net.Conn {
    cp.mutex.RLock()
    pool := cp.pools[peerID]
    cp.mutex.RUnlock()
    
    if pool == nil {
        return cp.createNewConnection(peerID)
    }
    
    select {
    case conn := <-pool.connections:
        return conn
    default:
        return cp.createNewConnection(peerID)
    }
}
```

**Batch Processing**:
```go
func batchPackets(packets <-chan []byte, maxBatch int, timeout time.Duration) <-chan [][]byte {
    batches := make(chan [][]byte)
    
    go func() {
        defer close(batches)
        
        batch := make([][]byte, 0, maxBatch)
        timer := time.NewTimer(timeout)
        
        for {
            select {
            case packet := <-packets:
                batch = append(batch, packet)
                if len(batch) >= maxBatch {
                    batches <- batch
                    batch = make([][]byte, 0, maxBatch)
                    timer.Reset(timeout)
                }
                
            case <-timer.C:
                if len(batch) > 0 {
                    batches <- batch
                    batch = make([][]byte, 0, maxBatch)
                }
                timer.Reset(timeout)
            }
        }
    }()
    
    return batches
}
```

---

## Cross-Platform Development

### 1. Platform Abstraction

**Interface Design**:
```go
// Platform-specific interface implementations
type TUNInterface interface {
    Create(name string, mtu int) error
    Read([]byte) (int, error)
    Write([]byte) (int, error)
    Close() error
}

// Linux implementation
type LinuxTUN struct {
    fd int
    name string
}

// Windows implementation  
type WindowsTUN struct {
    session wintun.Session
    name    string
}

// Factory function
func NewTUNInterface(name string) TUNInterface {
    switch runtime.GOOS {
    case "linux":
        return &LinuxTUN{name: name}
    case "windows":
        return &WindowsTUN{name: name}
    default:
        panic("unsupported platform")
    }
}
```

### 2. Build System

**Cross-Compilation Script**:
```bash
#!/bin/bash
# build.sh excerpt for cross-platform builds

platforms=(
    "linux/amd64"
    "linux/arm64"
    "windows/amd64"
    "darwin/amd64"
    "darwin/arm64"
)

for platform in "${platforms[@]}"; do
    GOOS=${platform%/*}
    GOARCH=${platform#*/}
    
    echo "Building for $GOOS/$GOARCH..."
    
    CGO_ENABLED=0 GOOS=$GOOS GOARCH=$GOARCH \
        go build -ldflags="-w -s" \
        -o "build/awl-$GOOS-$GOARCH" \
        ./cmd/awl
done
```

**Conditional Compilation**:
```go
//go:build linux
// +build linux

package vpn

import "golang.org/x/sys/unix"

func createTUN(name string) (*os.File, error) {
    fd, err := unix.Open("/dev/net/tun", unix.O_RDWR, 0)
    // Linux-specific implementation
}
```

```go
//go:build windows
// +build windows

package vpn

import "golang.zx2c4.com/wireguard/windows/tunnel/winipcfg"

func createTUN(name string) (*wintun.Adapter, error) {
    // Windows-specific implementation using Wintun
}
```

---

## Observability and Monitoring

### 1. Structured Logging

```go
import "go.uber.org/zap"

type Logger struct {
    *zap.Logger
    component string
}

func (l *Logger) LogPeerEvent(event string, peerID peer.ID, fields ...zap.Field) {
    baseFields := []zap.Field{
        zap.String("component", l.component),
        zap.String("event", event),
        zap.String("peer_id", peerID.String()),
        zap.Time("timestamp", time.Now()),
    }
    
    allFields := append(baseFields, fields...)
    l.Info("peer event", allFields...)
}

// Usage
logger.LogPeerEvent("connection_established", peerID,
    zap.Duration("connection_time", connectionTime),
    zap.String("transport", "quic"),
    zap.Int64("bytes_transferred", bytesTransferred),
)
```

### 2. Metrics Collection

```go
type Metrics struct {
    // Connection metrics
    activeConnections prometheus.Gauge
    connectionErrors  prometheus.Counter
    bytesTransferred  prometheus.Counter
    
    // VPN metrics
    packetsRouted     prometheus.Counter
    packetsDropped    prometheus.Counter
    routingLatency    prometheus.Histogram
    
    // P2P metrics
    peerCount         prometheus.Gauge
    dhcpQueries       prometheus.Counter
}

func (m *Metrics) RecordPacketRouted(size int, latency time.Duration) {
    m.packetsRouted.Inc()
    m.bytesTransferred.Add(float64(size))
    m.routingLatency.Observe(latency.Seconds())
}
```

### 3. Health Checks

```go
type HealthChecker struct {
    checks map[string]HealthCheck
    mutex  sync.RWMutex
}

type HealthCheck interface {
    Name() string
    Check(ctx context.Context) error
}

type P2PHealthCheck struct {
    p2p *p2p.P2p
}

func (p *P2PHealthCheck) Check(ctx context.Context) error {
    if len(p.p2p.ConnectedPeers()) == 0 {
        return errors.New("no connected peers")
    }
    
    // Test connectivity to bootstrap nodes
    return p.p2p.PingBootstrapNodes(ctx)
}
```

---

## Deployment and Distribution

### 1. Packaging

**Linux Packages**:
```bash
# Create .deb package
fpm -s dir -t deb \
    --name anywherelan \
    --version $VERSION \
    --description "Mesh VPN software" \
    --depends libcap2-bin \
    --after-install install-scripts/post-install.sh \
    build/awl-linux-amd64=/usr/bin/awl \
    config/awl.service=/etc/systemd/system/
```

**Windows Installer**:
```nsis
; NSIS installer script
Section "Main Application"
    SetOutPath $INSTDIR
    File "awl-windows-amd64.exe"
    File "wintun.dll"
    
    ; Create service
    ExecWait '"$INSTDIR\awl.exe" service install'
SectionEnd
```

### 2. Auto-Updates

```go
type UpdateManager struct {
    currentVersion string
    updateURL      string
    publicKey      ed25519.PublicKey
}

func (um *UpdateManager) CheckForUpdates() (*UpdateInfo, error) {
    resp, err := http.Get(um.updateURL + "/latest")
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    var update UpdateInfo
    if err := json.NewDecoder(resp.Body).Decode(&update); err != nil {
        return nil, err
    }
    
    // Verify signature
    if !um.verifySignature(update) {
        return nil, errors.New("invalid update signature")
    }
    
    return &update, nil
}
```

### 3. Configuration Management

**Environment-Specific Configs**:
```go
type Environment string

const (
    Development Environment = "development"
    Staging     Environment = "staging"
    Production  Environment = "production"
)

func LoadConfig(env Environment) (*Config, error) {
    baseConfig := getBaseConfig()
    
    switch env {
    case Development:
        baseConfig.LogLevel = "debug"
        baseConfig.BootstrapNodes = devBootstrapNodes
    case Production:
        baseConfig.LogLevel = "info"
        baseConfig.BootstrapNodes = prodBootstrapNodes
        baseConfig.EnableMetrics = true
    }
    
    return baseConfig, nil
}
```

---

## Contributing to AWL

### 1. Development Setup

```bash
# Fork the repository on GitHub
git clone https://github.com/your-username/awl.git
cd awl

# Create feature branch
git checkout -b feature/your-feature-name

# Set up development environment
go mod tidy
make dev-setup

# Run tests
make test

# Build all platforms
make build-all
```

### 2. Code Quality

**Linting and Formatting**:
```bash
# Format code
gofmt -w .
goimports -w .

# Run linters
golangci-lint run

# Check for security issues
gosec ./...

# Check for race conditions
go test -race ./...
```

**Testing Guidelines**:
```go
func TestTunnelPacketRouting(t *testing.T) {
    // Table-driven tests
    tests := []struct {
        name     string
        packet   []byte
        expected string
        wantErr  bool
    }{
        {
            name:     "valid packet to known peer",
            packet:   createTestPacket("10.66.0.1", "10.66.0.2"),
            expected: "peer-alice",
            wantErr:  false,
        },
        // More test cases...
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            tunnel := setupTestTunnel(t)
            result, err := tunnel.RoutePacket(tt.packet)
            
            if tt.wantErr {
                assert.Error(t, err)
                return
            }
            
            assert.NoError(t, err)
            assert.Equal(t, tt.expected, result)
        })
    }
}
```

### 3. Documentation

**API Documentation**:
```go
// Package tunnel provides VPN packet routing functionality.
//
// The tunnel service is responsible for routing IP packets between
// peers in the mesh network. It maintains a mapping between IP
// addresses and peer IDs, and forwards packets accordingly.
package tunnel

// RoutePacket routes an IP packet to the appropriate peer.
//
// The function examines the destination IP address in the packet
// header and looks up the corresponding peer in the routing table.
// If a route is found, the packet is forwarded to that peer.
//
// Parameters:
//   - packet: Raw IP packet bytes
//
// Returns:
//   - peerID: ID of the peer the packet was routed to
//   - error: Non-nil if routing failed
//
// Example:
//   peerID, err := tunnel.RoutePacket(ipPacket)
//   if err != nil {
//       log.Printf("Routing failed: %v", err)
//   }
func (t *Tunnel) RoutePacket(packet []byte) (peer.ID, error) {
    // Implementation...
}
```

This advanced guide provides the foundation for understanding and contributing to production mesh VPN systems. The concepts covered here are essential for building robust, secure, and performant networking applications.

Remember that building production networking software requires careful attention to security, performance, and reliability. Always test thoroughly, follow security best practices, and consider the real-world deployment constraints of your target environments.