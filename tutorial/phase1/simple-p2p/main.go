package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("=== AWL Tutorial: Simple P2P Communication ===")
		fmt.Println()
		fmt.Println("Usage: go run main.go [server|client] [port|address:port]")
		fmt.Println()
		fmt.Println("Examples:")
		fmt.Println("  go run main.go server 8080")
		fmt.Println("  go run main.go client localhost:8080")
		fmt.Println()
		fmt.Println("Start the server first, then connect with the client.")
		fmt.Println("Type messages in the client to see them echoed back.")
		fmt.Println("Type 'quit' to disconnect.")
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
			if strings.TrimSpace(text) == "quit" {
				return
			}
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