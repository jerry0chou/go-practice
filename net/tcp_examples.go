package net

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"time"
)

type TCPServer struct {
	Address string
	Port    string
	ln      net.Listener
}

func NewTCPServer(address, port string) *TCPServer {
	return &TCPServer{
		Address: address,
		Port:    port,
	}
}

func (s *TCPServer) Start() error {
	address := net.JoinHostPort(s.Address, s.Port)
	ln, err := net.Listen("tcp", address)
	if err != nil {
		return fmt.Errorf("failed to start TCP server: %w", err)
	}

	s.ln = ln
	fmt.Printf("🚀 TCP Server started on %s\n", address)

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Printf("❌ Error accepting connection: %v\n", err)
			continue
		}

		go s.handleConnection(conn)
	}
}

func (s *TCPServer) Stop() error {
	if s.ln != nil {
		return s.ln.Close()
	}
	return nil
}

func (s *TCPServer) handleConnection(conn net.Conn) {
	defer conn.Close()

	clientAddr := conn.RemoteAddr().String()
	fmt.Printf("📞 New connection from %s\n", clientAddr)

	conn.SetReadDeadline(time.Now().Add(30 * time.Second))

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		message := strings.TrimSpace(scanner.Text())
		fmt.Printf("📨 Received from %s: %s\n", clientAddr, message)

		response := fmt.Sprintf("Echo: %s\n", message)
		_, err := conn.Write([]byte(response))
		if err != nil {
			fmt.Printf("❌ Error writing to %s: %v\n", clientAddr, err)
			break
		}

		if strings.ToLower(message) == "quit" {
			fmt.Printf("👋 Client %s requested to quit\n", clientAddr)
			break
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("❌ Error reading from %s: %v\n", clientAddr, err)
	}

	fmt.Printf("🔌 Connection from %s closed\n", clientAddr)
}

type TCPClient struct {
	Address string
	Port    string
	conn    net.Conn
}

func NewTCPClient(address, port string) *TCPClient {
	return &TCPClient{
		Address: address,
		Port:    port,
	}
}

func (c *TCPClient) Connect() error {
	address := net.JoinHostPort(c.Address, c.Port)
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return fmt.Errorf("failed to connect to TCP server: %w", err)
	}

	c.conn = conn
	fmt.Printf("🔗 Connected to TCP server at %s\n", address)
	return nil
}

func (c *TCPClient) SendMessage(message string) error {
	if c.conn == nil {
		return fmt.Errorf("not connected to server")
	}

	_, err := c.conn.Write([]byte(message + "\n"))
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}

func (c *TCPClient) ReadResponse() (string, error) {
	if c.conn == nil {
		return "", fmt.Errorf("not connected to server")
	}

	c.conn.SetReadDeadline(time.Now().Add(5 * time.Second))

	response := make([]byte, 1024)
	n, err := c.conn.Read(response)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	return strings.TrimSpace(string(response[:n])), nil
}

func (c *TCPClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

func SimpleEchoServer(address, port string) error {
	server := NewTCPServer(address, port)
	return server.Start()
}

func SimpleEchoClient(address, port string, messages []string) error {
	client := NewTCPClient(address, port)

	if err := client.Connect(); err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}
	defer client.Close()

	for _, message := range messages {
		fmt.Printf("📤 Sending: %s\n", message)

		if err := client.SendMessage(message); err != nil {
			return fmt.Errorf("failed to send message: %w", err)
		}

		response, err := client.ReadResponse()
		if err != nil {
			return fmt.Errorf("failed to read response: %w", err)
		}

		fmt.Printf("📥 Received: %s\n", response)
		time.Sleep(1 * time.Second)
	}

	return nil
}

type ChatServer struct {
	Address    string
	Port       string
	clients    map[net.Conn]string
	broadcast  chan string
	register   chan net.Conn
	unregister chan net.Conn
	ln         net.Listener
}

func NewChatServer(address, port string) *ChatServer {
	return &ChatServer{
		Address:    address,
		Port:       port,
		clients:    make(map[net.Conn]string),
		broadcast:  make(chan string),
		register:   make(chan net.Conn),
		unregister: make(chan net.Conn),
	}
}

func (cs *ChatServer) Start() error {
	address := net.JoinHostPort(cs.Address, cs.Port)
	ln, err := net.Listen("tcp", address)
	if err != nil {
		return fmt.Errorf("failed to start chat server: %w", err)
	}

	cs.ln = ln
	fmt.Printf("💬 Chat Server started on %s\n", address)

	go cs.broadcaster()

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Printf("❌ Error accepting connection: %v\n", err)
			continue
		}

		cs.register <- conn
		go cs.handleChatConnection(conn)
	}
}

func (cs *ChatServer) Stop() error {
	if cs.ln != nil {
		return cs.ln.Close()
	}
	return nil
}

func (cs *ChatServer) broadcaster() {
	for {
		select {
		case conn := <-cs.register:
			clientAddr := conn.RemoteAddr().String()
			cs.clients[conn] = clientAddr
			fmt.Printf("👤 Client %s joined the chat\n", clientAddr)

			message := fmt.Sprintf("👤 %s joined the chat\n", clientAddr)
			cs.broadcastMessage(message, conn)

		case conn := <-cs.unregister:
			if clientAddr, ok := cs.clients[conn]; ok {
				delete(cs.clients, conn)
				conn.Close()
				fmt.Printf("👋 Client %s left the chat\n", clientAddr)

				message := fmt.Sprintf("👋 %s left the chat\n", clientAddr)
				cs.broadcastMessage(message, nil)
			}

		case message := <-cs.broadcast:
			cs.broadcastMessage(message, nil)
		}
	}
}

func (cs *ChatServer) broadcastMessage(message string, exclude net.Conn) {
	for conn := range cs.clients {
		if conn != exclude {
			_, err := conn.Write([]byte(message))
			if err != nil {
				cs.unregister <- conn
			}
		}
	}
}

func (cs *ChatServer) handleChatConnection(conn net.Conn) {
	defer func() {
		cs.unregister <- conn
	}()

	clientAddr := conn.RemoteAddr().String()
	conn.SetReadDeadline(time.Now().Add(60 * time.Second))

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		message := strings.TrimSpace(scanner.Text())
		if message == "" {
			continue
		}

		formattedMessage := fmt.Sprintf("[%s] %s: %s\n",
			time.Now().Format("15:04:05"), clientAddr, message)

		fmt.Printf("💬 %s", formattedMessage)
		cs.broadcast <- formattedMessage

		if strings.ToLower(message) == "quit" {
			break
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("❌ Error reading from %s: %v\n", clientAddr, err)
	}
}

func DemonstrateTCPOperations() {
	fmt.Println("🔌 TCP Operations Demo")
	fmt.Println(strings.Repeat("=", 60))

	fmt.Println("📝 TCP Examples Available:")
	fmt.Println("  1. Simple Echo Server/Client")
	fmt.Println("  2. Chat Server with Multiple Clients")
	fmt.Println("  3. Connection Handling with Timeouts")
	fmt.Println("  4. Error Handling and Recovery")

	fmt.Println("\n💡 To test TCP operations:")
	fmt.Println("  1. Start a server: go run run/net_main.go -mode=tcp-server")
	fmt.Println("  2. Start a client: go run run/net_main.go -mode=tcp-client")
	fmt.Println("  3. Start a chat server: go run run/net_main.go -mode=chat")

	fmt.Println("\n🔧 Available Functions:")
	fmt.Println("  - SimpleEchoServer(address, port)")
	fmt.Println("  - SimpleEchoClient(address, port, messages)")
	fmt.Println("  - NewChatServer(address, port)")
	fmt.Println("  - NewTCPServer(address, port)")
	fmt.Println("  - NewTCPClient(address, port)")
}
