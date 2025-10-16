package net

import (
	"fmt"
	"net"
	"strings"
	"time"
)

type UDPServer struct {
	Address string
	Port    string
	conn    *net.UDPConn
}

func NewUDPServer(address, port string) *UDPServer {
	return &UDPServer{
		Address: address,
		Port:    port,
	}
}

func (s *UDPServer) Start() error {
	address := net.JoinHostPort(s.Address, s.Port)
	udpAddr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return fmt.Errorf("failed to resolve UDP address: %w", err)
	}

	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		return fmt.Errorf("failed to start UDP server: %w", err)
	}

	s.conn = conn
	fmt.Printf("🚀 UDP Server started on %s\n", address)

	buffer := make([]byte, 1024)

	for {
		n, clientAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Printf("❌ Error reading from UDP: %v\n", err)
			continue
		}

		message := string(buffer[:n])
		fmt.Printf("📨 Received from %s: %s\n", clientAddr.String(), message)

		response := fmt.Sprintf("Echo: %s", message)
		_, err = conn.WriteToUDP([]byte(response), clientAddr)
		if err != nil {
			fmt.Printf("❌ Error writing to UDP: %v\n", err)
		} else {
			fmt.Printf("📤 Sent to %s: %s\n", clientAddr.String(), response)
		}
	}
}

func (s *UDPServer) Stop() error {
	if s.conn != nil {
		return s.conn.Close()
	}
	return nil
}

type UDPClient struct {
	Address string
	Port    string
	conn    *net.UDPConn
}

func NewUDPClient(address, port string) *UDPClient {
	return &UDPClient{
		Address: address,
		Port:    port,
	}
}

func (c *UDPClient) Connect() error {
	localAddr, err := net.ResolveUDPAddr("udp", ":0")
	if err != nil {
		return fmt.Errorf("failed to resolve local UDP address: %w", err)
	}

	conn, err := net.ListenUDP("udp", localAddr)
	if err != nil {
		return fmt.Errorf("failed to create UDP connection: %w", err)
	}

	c.conn = conn
	fmt.Printf("🔗 UDP Client connected from %s\n", conn.LocalAddr().String())
	return nil
}

func (c *UDPClient) SendMessage(message string) error {
	if c.conn == nil {
		return fmt.Errorf("not connected to server")
	}

	serverAddr, err := net.ResolveUDPAddr("udp", net.JoinHostPort(c.Address, c.Port))
	if err != nil {
		return fmt.Errorf("failed to resolve server address: %w", err)
	}

	_, err = c.conn.WriteToUDP([]byte(message), serverAddr)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}

func (c *UDPClient) ReadResponse() (string, net.Addr, error) {
	if c.conn == nil {
		return "", nil, fmt.Errorf("not connected to server")
	}

	c.conn.SetReadDeadline(time.Now().Add(5 * time.Second))

	buffer := make([]byte, 1024)
	n, addr, err := c.conn.ReadFromUDP(buffer)
	if err != nil {
		return "", nil, fmt.Errorf("failed to read response: %w", err)
	}

	return string(buffer[:n]), addr, nil
}

func (c *UDPClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

func SimpleUDPEchoServer(address, port string) error {
	server := NewUDPServer(address, port)
	return server.Start()
}

func SimpleUDPEchoClient(address, port string, messages []string) error {
	client := NewUDPClient(address, port)

	if err := client.Connect(); err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}
	defer client.Close()

	for _, message := range messages {
		fmt.Printf("📤 Sending: %s\n", message)

		if err := client.SendMessage(message); err != nil {
			return fmt.Errorf("failed to send message: %w", err)
		}

		response, addr, err := client.ReadResponse()
		if err != nil {
			return fmt.Errorf("failed to read response: %w", err)
		}

		fmt.Printf("📥 Received from %s: %s\n", addr.String(), response)
		time.Sleep(1 * time.Second)
	}

	return nil
}

type BroadcastServer struct {
	Address string
	Port    string
	conn    *net.UDPConn
}

func NewBroadcastServer(address, port string) *BroadcastServer {
	return &BroadcastServer{
		Address: address,
		Port:    port,
	}
}

func (bs *BroadcastServer) Start() error {
	address := net.JoinHostPort(bs.Address, bs.Port)
	udpAddr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return fmt.Errorf("failed to resolve UDP address: %w", err)
	}

	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		return fmt.Errorf("failed to start broadcast server: %w", err)
	}

	bs.conn = conn
	fmt.Printf("📡 Broadcast Server started on %s\n", address)

	conn.SetWriteBuffer(1024)

	go bs.broadcastMessages()

	buffer := make([]byte, 1024)
	for {
		n, clientAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Printf("❌ Error reading from UDP: %v\n", err)
			continue
		}

		message := string(buffer[:n])
		fmt.Printf("📨 Received from %s: %s\n", clientAddr.String(), message)
	}
}

func (bs *BroadcastServer) broadcastMessages() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	messageCount := 0
	for range ticker.C {
		messageCount++
		message := fmt.Sprintf("Broadcast message #%d from server", messageCount)

		interfaces, err := net.Interfaces()
		if err != nil {
			fmt.Printf("❌ Error getting interfaces: %v\n", err)
			continue
		}

		for _, iface := range interfaces {
			if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
				continue
			}

			addrs, err := iface.Addrs()
			if err != nil {
				continue
			}

			for _, addr := range addrs {
				if ipnet, ok := addr.(*net.IPNet); ok && ipnet.IP.To4() != nil {
					broadcastAddr := getBroadcastAddress(ipnet)
					if broadcastAddr != nil {
						udpAddr := &net.UDPAddr{
							IP:   broadcastAddr,
							Port: 8080,
						}

						_, err := bs.conn.WriteToUDP([]byte(message), udpAddr)
						if err != nil {
							fmt.Printf("❌ Error broadcasting to %s: %v\n", udpAddr.String(), err)
						} else {
							fmt.Printf("📡 Broadcasted to %s: %s\n", udpAddr.String(), message)
						}
					}
				}
			}
		}
	}
}

func getBroadcastAddress(ipnet *net.IPNet) net.IP {
	if ipnet.IP.To4() == nil {
		return nil
	}

	ip := ipnet.IP.To4()
	mask := ipnet.Mask

	broadcast := make(net.IP, 4)
	for i := 0; i < 4; i++ {
		broadcast[i] = ip[i] | (mask[i] ^ 0xff)
	}

	return broadcast
}

func (bs *BroadcastServer) Stop() error {
	if bs.conn != nil {
		return bs.conn.Close()
	}
	return nil
}

type MulticastServer struct {
	Address string
	Port    string
	conn    *net.UDPConn
}

func NewMulticastServer(address, port string) *MulticastServer {
	return &MulticastServer{
		Address: address,
		Port:    port,
	}
}

func (ms *MulticastServer) Start() error {
	address := net.JoinHostPort(ms.Address, ms.Port)
	udpAddr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return fmt.Errorf("failed to resolve UDP address: %w", err)
	}

	conn, err := net.ListenMulticastUDP("udp", nil, udpAddr)
	if err != nil {
		return fmt.Errorf("failed to start multicast server: %w", err)
	}

	ms.conn = conn
	fmt.Printf("📺 Multicast Server started on %s\n", address)

	go ms.multicastMessages()

	buffer := make([]byte, 1024)
	for {
		n, clientAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Printf("❌ Error reading from UDP: %v\n", err)
			continue
		}

		message := string(buffer[:n])
		fmt.Printf("📨 Received from %s: %s\n", clientAddr.String(), message)
	}
}

func (ms *MulticastServer) multicastMessages() {
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	messageCount := 0
	for range ticker.C {
		messageCount++
		message := fmt.Sprintf("Multicast message #%d from server", messageCount)

		address := net.JoinHostPort(ms.Address, ms.Port)
		udpAddr, err := net.ResolveUDPAddr("udp", address)
		if err != nil {
			fmt.Printf("❌ Error resolving multicast address: %v\n", err)
			continue
		}

		_, err = ms.conn.WriteToUDP([]byte(message), udpAddr)
		if err != nil {
			fmt.Printf("❌ Error sending multicast: %v\n", err)
		} else {
			fmt.Printf("📺 Multicasted: %s\n", message)
		}
	}
}

func (ms *MulticastServer) Stop() error {
	if ms.conn != nil {
		return ms.conn.Close()
	}
	return nil
}

func DemonstrateUDPOperations() {
	fmt.Println("📡 UDP Operations Demo")
	fmt.Println(strings.Repeat("=", 60))

	fmt.Println("📝 UDP Examples Available:")
	fmt.Println("  1. Simple UDP Echo Server/Client")
	fmt.Println("  2. UDP Broadcast Server")
	fmt.Println("  3. UDP Multicast Server")
	fmt.Println("  4. Connectionless Communication")
	fmt.Println("  5. Fire-and-Forget Messaging")

	fmt.Println("\n💡 To test UDP operations:")
	fmt.Println("  1. Start a UDP server: go run run/net_main.go -mode=udp-server")
	fmt.Println("  2. Start a UDP client: go run run/net_main.go -mode=udp-client")
	fmt.Println("  3. Start a broadcast server: go run run/net_main.go -mode=broadcast")
	fmt.Println("  4. Start a multicast server: go run run/net_main.go -mode=multicast")

	fmt.Println("\n🔧 Available Functions:")
	fmt.Println("  - SimpleUDPEchoServer(address, port)")
	fmt.Println("  - SimpleUDPEchoClient(address, port, messages)")
	fmt.Println("  - NewBroadcastServer(address, port)")
	fmt.Println("  - NewMulticastServer(address, port)")
	fmt.Println("  - NewUDPServer(address, port)")
	fmt.Println("  - NewUDPClient(address, port)")

	fmt.Println("\n📊 UDP vs TCP Comparison:")
	fmt.Println("  UDP:")
	fmt.Println("    ✅ Faster (no connection overhead)")
	fmt.Println("    ✅ Lower latency")
	fmt.Println("    ✅ Good for real-time applications")
	fmt.Println("    ❌ No delivery guarantee")
	fmt.Println("    ❌ No ordering guarantee")
	fmt.Println("    ❌ No congestion control")
	fmt.Println()
	fmt.Println("  TCP:")
	fmt.Println("    ✅ Reliable delivery")
	fmt.Println("    ✅ Ordered delivery")
	fmt.Println("    ✅ Congestion control")
	fmt.Println("    ❌ Higher overhead")
	fmt.Println("    ❌ Higher latency")
	fmt.Println("    ❌ Connection state required")
}
