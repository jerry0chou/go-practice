package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/jerrychou/go-practice/net"
)

func main() {
	mode := flag.String("mode", "demo", "Mode to run: demo, url, network, tcp-server, tcp-client, udp-server, udp-client, chat, broadcast, multicast")
	address := flag.String("address", "localhost", "Server address")
	port := flag.String("port", "8080", "Server port")
	flag.Parse()

	fmt.Println("🌐 Go Network Package Demo")
	fmt.Println(strings.Repeat("=", 50))

	switch *mode {
	case "demo":
		runDemo()
	case "url":
		runURLDemo()
	case "network":
		runNetworkDemo()
	case "tcp-server":
		runTCPServer(*address, *port)
	case "tcp-client":
		runTCPClient(*address, *port)
	case "udp-server":
		runUDPServer(*address, *port)
	case "udp-client":
		runUDPClient(*address, *port)
	case "chat":
		runChatServer(*address, *port)
	case "broadcast":
		runBroadcastServer(*address, *port)
	case "multicast":
		runMulticastServer(*address, *port)
	default:
		fmt.Printf("❌ Unknown mode: %s\n", *mode)
		fmt.Println("Available modes: demo, url, network, tcp-server, tcp-client, udp-server, udp-client, client, chat, broadcast, multicast")
		os.Exit(1)
	}
}

func runDemo() {
	fmt.Println("🎯 Running Complete Demo")
	fmt.Println(strings.Repeat("=", 50))

	fmt.Println("\n" + strings.Repeat("=", 60))
	net.DemonstrateURLOperations()

	fmt.Println("\n" + strings.Repeat("=", 60))
	net.DemonstrateNetworkOperations()

	fmt.Println("\n" + strings.Repeat("=", 60))
	net.DemonstrateTCPOperations()

	fmt.Println("\n" + strings.Repeat("=", 60))
	net.DemonstrateUDPOperations()

	fmt.Println("\n🎉 Demo completed!")
	fmt.Println("\n💡 To run specific demos:")
	fmt.Println("  go run run/net_main.go -mode=url")
	fmt.Println("  go run run/net_main.go -mode=network")
	fmt.Println("  go run run/net_main.go -mode=tcp-server")
	fmt.Println("  go run run/net_main.go -mode=udp-server")
}

func runURLDemo() {
	fmt.Println("🔗 URL Operations Demo")
	fmt.Println(strings.Repeat("=", 50))
	net.DemonstrateURLOperations()
}

func runNetworkDemo() {
	fmt.Println("🌐 Network Operations Demo")
	fmt.Println(strings.Repeat("=", 50))
	net.DemonstrateNetworkOperations()
	net.PrintNetworkInfo()
}

func runTCPServer(address, port string) {
	fmt.Printf("🔌 Starting TCP Server on %s:%s\n", address, port)
	fmt.Println("Press Ctrl+C to stop the server")

	server := net.NewTCPServer(address, port)
	if err := server.Start(); err != nil {
		log.Fatalf("❌ Failed to start TCP server: %v", err)
	}
}

func runTCPClient(address, port string) {
	fmt.Printf("🔗 Starting TCP Client to %s:%s\n", address, port)

	messages := []string{
		"Hello, TCP Server!",
		"How are you?",
		"This is a test message",
		"Goodbye!",
		"quit",
	}

	if err := net.SimpleEchoClient(address, port, messages); err != nil {
		log.Fatalf("❌ TCP Client error: %v", err)
	}
}

func runUDPServer(address, port string) {
	fmt.Printf("📡 Starting UDP Server on %s:%s\n", address, port)
	fmt.Println("Press Ctrl+C to stop the server")

	server := net.NewUDPServer(address, port)
	if err := server.Start(); err != nil {
		log.Fatalf("❌ Failed to start UDP server: %v", err)
	}
}

func runUDPClient(address, port string) {
	fmt.Printf("📤 Starting UDP Client to %s:%s\n", address, port)

	messages := []string{
		"Hello, UDP Server!",
		"This is UDP message 1",
		"This is UDP message 2",
		"UDP is connectionless!",
		"Goodbye UDP!",
	}

	if err := net.SimpleUDPEchoClient(address, port, messages); err != nil {
		log.Fatalf("❌ UDP Client error: %v", err)
	}
}

func runChatServer(address, port string) {
	fmt.Printf("💬 Starting Chat Server on %s:%s\n", address, port)
	fmt.Println("Multiple clients can connect and chat with each other")
	fmt.Println("Press Ctrl+C to stop the server")

	chatServer := net.NewChatServer(address, port)
	if err := chatServer.Start(); err != nil {
		log.Fatalf("❌ Failed to start chat server: %v", err)
	}
}

func runBroadcastServer(address, port string) {
	fmt.Printf("📡 Starting Broadcast Server on %s:%s\n", address, port)
	fmt.Println("Server will broadcast messages to all network interfaces")
	fmt.Println("Press Ctrl+C to stop the server")

	broadcastServer := net.NewBroadcastServer(address, port)
	if err := broadcastServer.Start(); err != nil {
		log.Fatalf("❌ Failed to start broadcast server: %v", err)
	}
}

func runMulticastServer(address, port string) {
	fmt.Printf("📺 Starting Multicast Server on %s:%s\n", address, port)
	fmt.Println("Server will send multicast messages")
	fmt.Println("Press Ctrl+C to stop the server")

	if address == "localhost" {
		address = "224.0.0.1"
	}

	multicastServer := net.NewMulticastServer(address, port)
	if err := multicastServer.Start(); err != nil {
		log.Fatalf("❌ Failed to start multicast server: %v", err)
	}
}

func exampleUsage() {
	fmt.Println("📚 Example Usage:")
	fmt.Println(strings.Repeat("=", 50))

	fmt.Println("\n🔗 URL Operations:")
	fmt.Println("  // Parse a URL")
	fmt.Println("  info, err := net.ParseURL(\"https://example.com/path?param=value\")")
	fmt.Println("  if err != nil { log.Fatal(err) }")
	fmt.Println("  fmt.Printf(\"Host: %s\\n\", info.Host)")

	fmt.Println("\n  // Build a URL")
	fmt.Println("  url := net.BuildURL(\"https\", \"api.example.com\", \"/v1/users\", map[string]string{\"page\": \"1\"})")
	fmt.Println("  fmt.Println(url)")

	fmt.Println("\n🌐 Network Operations:")
	fmt.Println("  // Resolve hostname")
	fmt.Println("  ips, err := net.ResolveHostname(\"google.com\")")
	fmt.Println("  if err != nil { log.Fatal(err) }")
	fmt.Println("  fmt.Printf(\"IPs: %v\\n\", ips)")

	fmt.Println("\n  // Check if IP is private")
	fmt.Println("  isPrivate := net.IsPrivateIP(\"192.168.1.1\")")
	fmt.Println("  fmt.Printf(\"Is private: %t\\n\", isPrivate)")

	fmt.Println("\n🔌 TCP Operations:")
	fmt.Println("  // Start TCP server")
	fmt.Println("  server := net.NewTCPServer(\"localhost\", \"8080\")")
	fmt.Println("  go server.Start()")
	fmt.Println("  time.Sleep(1 * time.Second)")

	fmt.Println("\n  // Connect TCP client")
	fmt.Println("  client := net.NewTCPClient(\"localhost\", \"8080\")")
	fmt.Println("  client.Connect()")
	fmt.Println("  client.SendMessage(\"Hello Server!\")")
	fmt.Println("  response, _ := client.ReadResponse()")
	fmt.Println("  fmt.Println(response)")

	fmt.Println("\n📡 UDP Operations:")
	fmt.Println("  // Start UDP server")
	fmt.Println("  server := net.NewUDPServer(\"localhost\", \"8080\")")
	fmt.Println("  go server.Start()")
	fmt.Println("  time.Sleep(1 * time.Second)")

	fmt.Println("\n  // Connect UDP client")
	fmt.Println("  client := net.NewUDPClient(\"localhost\", \"8080\")")
	fmt.Println("  client.Connect()")
	fmt.Println("  client.SendMessage(\"Hello UDP Server!\")")
	fmt.Println("  response, addr, _ := client.ReadResponse()")
	fmt.Println("  fmt.Printf(\"Response from %s: %s\\n\", addr, response)")
}
