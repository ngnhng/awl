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
	fmt.Println("=== AWL Tutorial: TUN Interface ===\n")
	
	// Create TUN interface
	fmt.Println("Creating TUN interface...")
	tunDevice, err := tun.CreateTUN("awl-tutorial", 1500)
	if err != nil {
		fmt.Printf("Error creating TUN: %v\n", err)
		fmt.Println("\nNote: This requires admin/root privileges")
		fmt.Println("On Linux/macOS: sudo go run main.go")
		fmt.Println("On Windows: Run as Administrator")
		os.Exit(1)
	}
	defer tunDevice.Close()

	fmt.Println("✓ Created TUN interface: awl-tutorial")
	fmt.Println("\nTo configure the interface, run these commands:")
	fmt.Println("  Linux/macOS:")
	fmt.Println("    sudo ip addr add 10.66.0.1/24 dev awl-tutorial")
	fmt.Println("    sudo ip link set awl-tutorial up")
	fmt.Println("  Windows:")
	fmt.Println("    netsh interface ip set address \"awl-tutorial\" static 10.66.0.1 255.255.255.0")
	
	fmt.Println("\nTest commands after configuration:")
	fmt.Println("  ping 10.66.0.2   # Should see packet captured")
	fmt.Println("  ping 8.8.8.8     # Should see packet if routing configured")
	
	fmt.Println("\nPress Ctrl+C to stop...\n")

	// Handle graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	// Read packets in a goroutine
	go readPackets(tunDevice)

	// Wait for signal
	<-sigCh
	fmt.Println("\nShutting down TUN interface...")
}

func readPackets(tunDevice tun.Device) {
	packet := make([]byte, 1500)
	packetCount := 0
	
	fmt.Println("Listening for packets...")
	
	for {
		n, err := tunDevice.Read(packet[:], 0)
		if err != nil {
			fmt.Printf("Error reading packet: %v\n", err)
			continue
		}

		if n > 0 {
			packetCount++
			analyzePacket(packet[:n], packetCount)
		}
	}
}

func analyzePacket(packet []byte, count int) {
	if len(packet) < 20 {
		fmt.Printf("Packet #%d: Too short (%d bytes)\n", count, len(packet))
		return
	}

	// Parse basic IP header (simplified)
	version := packet[0] >> 4
	headerLen := int(packet[0]&0x0F) * 4
	protocol := packet[9]
	srcIP := net.IP(packet[12:16])
	dstIP := net.IP(packet[16:20])

	protocolName := getProtocolName(protocol)
	
	fmt.Printf("Packet #%d: IPv%d %s -> %s (%s) %d bytes\n",
		count, version, srcIP, dstIP, protocolName, len(packet))

	// Show packet details for first few packets
	if count <= 3 {
		fmt.Printf("  Header Length: %d bytes\n", headerLen)
		fmt.Printf("  Protocol: %d (%s)\n", protocol, protocolName)
		fmt.Printf("  Raw Header: %x...\n", packet[:min(20, len(packet))])
		
		// For ICMP packets, show type
		if protocol == 1 && len(packet) > 20 {
			icmpType := packet[20]
			fmt.Printf("  ICMP Type: %d (%s)\n", icmpType, getICMPType(icmpType))
		}
		fmt.Println()
	}

	// In a real mesh VPN, we would:
	fmt.Printf("  -> Route to peer responsible for %s\n", dstIP)
	fmt.Printf("  -> Or drop if no route available\n\n")
}

func getProtocolName(protocol byte) string {
	switch protocol {
	case 1:
		return "ICMP"
	case 6:
		return "TCP"
	case 17:
		return "UDP"
	default:
		return fmt.Sprintf("Unknown(%d)", protocol)
	}
}

func getICMPType(icmpType byte) string {
	switch icmpType {
	case 0:
		return "Echo Reply"
	case 8:
		return "Echo Request"
	case 3:
		return "Destination Unreachable"
	case 11:
		return "Time Exceeded"
	default:
		return fmt.Sprintf("Type %d", icmpType)
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}