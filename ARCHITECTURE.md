# AWL (Anywherelan) High-Level Architecture

## Overview

AWL (Anywherelan) is a mesh VPN application written in Go that enables secure peer-to-peer networking between devices without requiring centralized coordination servers. It creates a virtual private network overlay that allows devices to communicate directly, similar to Tailscale, Tinc, or WireGuard, but with a fully decentralized approach.

## Key Architectural Principles

1. **Peer-to-Peer First**: No central servers required for operation
2. **Event-Driven Design**: Components communicate through a central event bus
3. **Modular Architecture**: Clear separation of concerns with well-defined interfaces
4. **Cross-Platform Support**: Single codebase supporting multiple operating systems
5. **Transport Agnostic**: Multiple transport protocols (QUIC, TCP) with automatic selection

## Core Components

### Application Entry Points

The application provides multiple interfaces for different use cases:

- **`cmd/awl/`** - Headless server version for Linux servers and automation
- **`cmd/awl-tray/`** - Desktop application with system tray integration
- **`cmd/gomobile-lib/`** - Android library for mobile application

### Central Application Structure

The main `Application` struct (`application.go`) orchestrates all components:

```go
type Application struct {
    LogBuffer  *ringbuffer.RingBuffer  // Centralized logging
    Conf       *config.Config          // Configuration management
    Eventbus   awlevent.Bus           // Event-driven communication
    vpnDevice  *vpn.Device            // Virtual network interface
    P2p        *p2p.P2p               // Peer-to-peer networking
    Api        *api.Handler           // REST API for web interface
    AuthStatus *service.AuthStatus    // Peer authentication
    Tunnel     *service.Tunnel        // VPN packet routing
    SOCKS5     *service.SOCKS5        // Proxy service
    Dns        *DNSService            // DNS resolution for .awl domains
}
```

## Detailed Component Architecture

### 1. P2P Networking Layer (`p2p/`)

**Purpose**: Handles all peer-to-peer networking, discovery, and connection management.

**Key Technologies**:
- **libp2p**: IPFS networking stack for robust P2P networking
- **DHT (Distributed Hash Table)**: For decentralized peer discovery
- **QUIC/TCP**: Multiple transport protocols with TLS encryption
- **NAT Traversal**: Automatic hole punching and relay fallback

**Responsibilities**:
- Peer discovery through bootstrap nodes
- Connection establishment and maintenance
- Transport protocol selection and management
- Network metrics and connection quality monitoring

**Key Files**:
- `p2p/p2p.go` - Main P2P service implementation
- `p2p/metrics.go` - Network performance monitoring

### 2. VPN Layer (`vpn/`)

**Purpose**: Manages the virtual network interface and packet processing.

**Key Technologies**:
- **TUN/TAP interfaces**: Layer 3 virtual networking
- **Platform-specific implementations**: Linux, macOS, Windows, Android
- **WireGuard libraries**: For some platforms (Windows uses Wintun)

**Responsibilities**:
- Creating and managing virtual network interfaces
- Packet capture and injection
- IP address management and routing
- Platform-specific networking optimizations

**Key Files**:
- `vpn/vpn.go` - Core VPN device implementation
- `vpn/iface_*.go` - Platform-specific interface implementations

### 3. Service Layer (`service/`)

**Purpose**: High-level services that coordinate between P2P and VPN layers.

#### Tunnel Service (`service/tunnel.go`)
- Routes VPN packets between peers
- Maintains peer-to-IP mappings
- Handles packet forwarding and routing decisions

#### SOCKS5 Service (`service/socks5.go`)
- Provides SOCKS5 proxy functionality
- Allows routing traffic through remote peers
- Supports both TCP and UDP proxying

#### Authentication Service (`service/auth_status.go`)
- Manages peer authentication and trust relationships
- Handles friend requests and peer approval
- Maintains peer status and connectivity information

### 4. API Layer (`api/`)

**Purpose**: Provides REST API for web interface and external integration.

**Responsibilities**:
- Web-based configuration interface
- Peer management endpoints
- Real-time status and metrics
- Configuration updates and validation

**Key Files**:
- `api/api.go` - Main API handler and routing
- `api/peers.go` - Peer management endpoints
- `api/settings.go` - Configuration management endpoints

### 5. Configuration Management (`config/`)

**Purpose**: Handles all application configuration and state persistence.

**Features**:
- JSON-based configuration files
- Cryptographic identity management
- Network configuration (IP allocation, subnets)
- Peer relationship storage
- Cross-platform config directory handling

### 6. DNS Service (`awldns/`)

**Purpose**: Provides DNS resolution for `.awl` domain names.

**Features**:
- Maps peer names to VPN IP addresses
- Integrates with system DNS configuration
- Supports both IPv4 and IPv6 resolution
- Platform-specific DNS integration

## Data Flow and Component Interactions

### 1. Application Startup Flow

```
1. Load Configuration → 2. Initialize P2P Host → 3. Create VPN Interface
                    ↓
4. Start Services → 5. Register Protocol Handlers → 6. Begin Peer Discovery
```

### 2. Peer Connection Flow

```
Peer Discovery (DHT) → Authentication Exchange → Direct Connection Attempt
                                              ↓
                        Relay Connection (if direct fails) → Stream Establishment
```

### 3. VPN Packet Flow

```
Application → TUN Interface → VPN Device → Tunnel Service → P2P Stream → Remote Peer
```

### 4. SOCKS5 Proxy Flow

```
Local App → SOCKS5 Service → P2P Connection → Remote Peer → Internet
```

## Event-Driven Architecture

The application uses a central event bus (`awlevent/`) for loose coupling between components:

**Key Events**:
- `KnownPeerChanged` - Triggers tunnel peer list refresh
- Peer connection status changes
- Configuration updates
- Network interface changes

**Benefits**:
- Decoupled component communication
- Easy testing and mocking
- Clear data flow understanding
- Extensible event handling

## Security Model

### Cryptographic Foundation
- **Ed25519** public/private key pairs for peer identity
- **TLS** encryption for all peer-to-peer communication
- **Peer verification** through cryptographic signatures

### Trust Model
- **Explicit trust**: Peers must be manually added (friend requests)
- **No automatic trust**: New peers require approval
- **Revocation support**: Peers can be blocked or removed

### Network Isolation
- **Private IP ranges**: Uses 10.66.0.0/16 by default
- **Interface isolation**: VPN traffic separated from regular network
- **DNS isolation**: .awl domains resolved internally

## Platform-Specific Considerations

### Linux
- Uses standard TUN interfaces
- Integrates with systemd for service management
- Supports various architectures (amd64, arm64, mips, etc.)

### Windows
- Uses Wintun driver for better performance
- Requires administrator privileges for interface creation
- System tray integration for user experience

### macOS
- Standard TUN interface support
- System tray integration
- Requires admin privileges for interface creation

### Android
- Uses VpnService API for packet capture
- gomobile bindings for Go integration
- Flutter-based user interface

## Build and Distribution

The project uses a comprehensive build system (`build.sh`) that:
- Cross-compiles for multiple platforms and architectures
- Bundles platform-specific dependencies (like Wintun)
- Creates distributable packages for each platform
- Supports both server and desktop variants

## Performance Characteristics

### Network Performance
- **MTU**: 3500 bytes default for better performance
- **Zero-copy**: Where possible, minimizes packet copying
- **Connection pooling**: Reuses P2P connections efficiently
- **Transport selection**: Automatically chooses best transport (QUIC vs TCP)

### Resource Usage
- **Memory efficient**: Ring buffers and object pooling
- **CPU optimized**: Concurrent packet processing
- **Minimal overhead**: Direct peer connections reduce latency

This architecture provides a robust, scalable, and secure foundation for mesh VPN networking with excellent cross-platform support and user experience.