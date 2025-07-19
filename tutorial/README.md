# AWL Tutorial Code Examples

This directory contains hands-on code examples for learning how to build a mesh VPN like AWL.

## Prerequisites

- Go 1.24 or later
- Basic understanding of networking concepts
- On Linux/macOS: `sudo` privileges for TUN interface examples
- On Windows: Administrator privileges for TUN interface examples

## Tutorial Structure

### Phase 1: Basic P2P Networking
Learn fundamental peer-to-peer networking concepts:
- **simple-p2p**: Basic TCP client/server communication
- **discovery**: Bootstrap server for peer discovery
- **peer**: Peer client that registers and discovers other peers

```bash
# Try the simple P2P example
cd phase1/simple-p2p
go run main.go server 8080        # Terminal 1
go run main.go client localhost:8080  # Terminal 2
```

### Phase 2: Cryptography and Security
Add cryptographic identity and authentication:
- **identity**: Generate Ed25519 key pairs for peer identity
- **auth**: Challenge-response authentication protocol

```bash
# Learn about cryptographic identity
cd phase2/identity
go run main.go

# Try the authentication protocol
cd phase2/auth
go run main.go
```

### Phase 3: VPN Functionality
Create virtual network interfaces and packet routing:
- **tun**: Create and use TUN interfaces (requires admin privileges)
- **routing**: Implement packet routing logic

```bash
# Create a TUN interface (needs sudo/admin)
cd phase3/tun
sudo go run main.go

# Learn about packet routing
cd phase3/routing
go run main.go
```

### Phase 4: libp2p Integration
Use professional P2P networking library:
- **libp2p-basic**: Basic libp2p host with stream communication

```bash
# Try libp2p networking
cd phase4/libp2p-basic
go mod tidy
go run main.go                    # Terminal 1 (copy the address)
go run main.go <address>          # Terminal 2 (paste address)
```

## Learning Objectives

Each phase builds on the previous one:

1. **Phase 1**: Understand basic networking, TCP connections, peer discovery
2. **Phase 2**: Learn cryptographic concepts, digital signatures, authentication
3. **Phase 3**: Understand virtual networking, packet processing, routing
4. **Phase 4**: Use production-grade P2P libraries, advanced networking

## Running the Examples

### Phase 1 Examples
```bash
# Bootstrap discovery demo
cd phase1/discovery
go run main.go &                  # Start bootstrap server

cd ../peer
go run main.go peer1 192.168.1.100:9001 "Alice"  # Terminal 1
go run main.go peer2 192.168.1.101:9001 "Bob"    # Terminal 2
```

### Phase 2 Examples
```bash
# All Phase 2 examples are self-contained
cd phase2/identity && go run main.go
cd phase2/auth && go run main.go
```

### Phase 3 Examples
```bash
# TUN interface (requires privileges)
cd phase3/tun
sudo go run main.go

# In another terminal, configure the interface:
sudo ip addr add 10.66.0.1/24 dev awl-tutorial
sudo ip link set awl-tutorial up
ping 10.66.0.2  # Should see packets in the first terminal

# Routing demo (no privileges needed)
cd phase3/routing
go run main.go
```

### Phase 4 Examples
```bash
cd phase4/libp2p-basic
go mod tidy
go run main.go  # Copy one of the addresses shown

# In another terminal:
go run main.go /ip4/127.0.0.1/tcp/XXXX/p2p/PEER_ID
# Then type messages to chat between peers
```

## Key Concepts Demonstrated

### Networking Concepts
- TCP client/server programming
- Peer-to-peer discovery and connection
- Virtual network interfaces (TUN)
- Packet parsing and routing
- libp2p networking stack

### Security Concepts
- Ed25519 cryptographic keys
- Digital signatures
- Challenge-response authentication
- Preventing replay attacks
- Secure peer identity

### Go Programming Patterns
- Goroutines and channels
- Context for cancellation
- Interfaces for abstraction
- Mutex for thread safety
- Error handling

## Troubleshooting

### Permission Issues
- TUN examples require admin/root privileges
- On Linux/macOS: use `sudo`
- On Windows: Run as Administrator

### Network Issues
- Ensure firewall allows connections
- Use `netstat` or `ss` to check listening ports
- Try `127.0.0.1` if external IPs don't work

### Build Issues
- Ensure Go 1.24+ is installed
- Run `go mod tidy` in directories with go.mod files
- Check that all imports are available

## Next Steps

After completing these tutorials:

1. **Study the main AWL codebase** to see production implementation
2. **Read the LEARNING_GUIDE.md** for deeper theoretical knowledge
3. **Try building your own mini VPN** using the concepts learned
4. **Contribute to AWL** by adding features or improving documentation

## Additional Resources

- [AWL Architecture](../ARCHITECTURE.md) - High-level system design
- [AWL Technical Analysis](../TECHNICAL_ANALYSIS.md) - Detailed implementation
- [libp2p Documentation](https://docs.libp2p.io/) - P2P networking library
- [WireGuard Protocol](https://www.wireguard.com/protocol/) - Modern VPN design

Happy learning! 🚀