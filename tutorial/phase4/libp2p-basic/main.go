package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
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
	fmt.Println("=== AWL Tutorial: libp2p Networking ===\n")
	
	ctx := context.Background()

	// Create libp2p host
	fmt.Println("Creating libp2p host...")
	h, err := libp2p.New(
		libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/0"),
	)
	if err != nil {
		panic(err)
	}
	defer h.Close()

	// Set stream handler for incoming connections
	h.SetStreamHandler(protocol.ID(ProtocolID), handleStream)

	// Print host information
	fmt.Printf("✓ Host created successfully!\n")
	fmt.Printf("Host ID: %s\n", h.ID())
	fmt.Printf("Listening on:\n")
	for _, addr := range h.Addrs() {
		fmt.Printf("  %s/p2p/%s\n", addr, h.ID())
	}

	// If peer address provided, connect to it
	if len(os.Args) > 1 {
		peerAddr := os.Args[1]
		fmt.Printf("\nConnecting to peer: %s\n", peerAddr)
		go connectToPeer(ctx, h, peerAddr)
	} else {
		fmt.Printf("\nTo connect a second peer, run:\n")
		fmt.Printf("go run main.go <multiaddr>\n")
		fmt.Printf("Example: go run main.go /ip4/127.0.0.1/tcp/PORT/p2p/PEER_ID\n")
	}

	fmt.Println("\nWaiting for connections... Press Ctrl+C to exit")

	// Handle graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
	
	fmt.Println("\nShutting down...")
}

func handleStream(s network.Stream) {
	defer s.Close()
	
	remotePeer := s.Conn().RemotePeer()
	fmt.Printf("\n📨 New stream from peer: %s\n", remotePeer.ShortString())
	fmt.Printf("   Remote Address: %s\n", s.Conn().RemoteMultiaddr())
	fmt.Printf("   Protocol: %s\n", s.Protocol())

	// Read messages from the stream
	reader := bufio.NewReader(s)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("   Stream closed by %s\n", remotePeer.ShortString())
			break
		}
		
		message = strings.TrimSpace(message)
		fmt.Printf("   📩 Received: %s\n", message)
		
		// Echo back with a prefix
		response := fmt.Sprintf("Echo from %s: %s\n", s.Conn().LocalPeer().ShortString(), message)
		s.Write([]byte(response))
		
		if message == "quit" {
			fmt.Printf("   Received quit command, closing stream\n")
			break
		}
	}
}

func connectToPeer(ctx context.Context, h host.Host, peerAddr string) {
	// Parse peer address
	maddr, err := multiaddr.NewMultiaddr(peerAddr)
	if err != nil {
		fmt.Printf("❌ Error parsing address: %v\n", err)
		return
	}

	// Extract peer info
	peerInfo, err := peer.AddrInfoFromP2pAddr(maddr)
	if err != nil {
		fmt.Printf("❌ Error extracting peer info: %v\n", err)
		return
	}

	fmt.Printf("Connecting to peer ID: %s\n", peerInfo.ID.ShortString())

	// Connect to peer
	err = h.Connect(ctx, *peerInfo)
	if err != nil {
		fmt.Printf("❌ Error connecting to peer: %v\n", err)
		return
	}

	fmt.Printf("✅ Connected to peer: %s\n", peerInfo.ID.ShortString())

	// Open stream
	s, err := h.NewStream(ctx, peerInfo.ID, protocol.ID(ProtocolID))
	if err != nil {
		fmt.Printf("❌ Error opening stream: %v\n", err)
		return
	}
	defer s.Close()

	fmt.Println("\n💬 Interactive chat started!")
	fmt.Println("Type messages and press Enter. Type 'quit' to exit.")
	fmt.Printf("Chatting with: %s\n\n", peerInfo.ID.ShortString())

	// Start reading responses in a goroutine
	go func() {
		reader := bufio.NewReader(s)
		for {
			response, err := reader.ReadString('\n')
			if err != nil {
				return
			}
			fmt.Printf("📨 %s", response)
		}
	}()

	// Read user input and send messages
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		message := scanner.Text()
		if strings.TrimSpace(message) == "" {
			continue
		}
		
		// Send message
		_, err := s.Write([]byte(message + "\n"))
		if err != nil {
			fmt.Printf("❌ Error sending message: %v\n", err)
			break
		}
		
		if strings.TrimSpace(message) == "quit" {
			break
		}
	}
}