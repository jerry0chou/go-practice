package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	http "github.com/jerrychou/go-practice/http"
)

func main() {
	fmt.Println("🚀 Go HTTP Package Demo")
	fmt.Println("========================")

	for {
		showMenu()
		choice := getUserInput("Please select an option (1-8): ")

		switch choice {
		case "1":
			startHTTPServer()
		case "2":
			demoBasicClient()
		case "3":
			demoAdvancedClient()
		case "4":
			demoGitHubAPI()
		case "5":
			demoHTTPUtils()
		case "6":
			demoJSONUtils()
		case "7":
			demoAllExamples()
		case "8":
			fmt.Println("👋 Goodbye!")
			return
		default:
			fmt.Println("❌ Invalid selection, please try again")
		}

		fmt.Println("\nPress Enter to continue...")
		bufio.NewReader(os.Stdin).ReadLine()
	}
}

func showMenu() {
	fmt.Println("\n📋 Available Operations:")
	fmt.Println("1. 🌐 Start HTTP Server")
	fmt.Println("2. 📡 Basic HTTP Client Examples")
	fmt.Println("3. 🔧 Advanced HTTP Client Examples")
	fmt.Println("4. 🐙 GitHub API Examples")
	fmt.Println("5. 🛠️  HTTP Utility Functions Examples")
	fmt.Println("6. 📄 JSON Utility Functions Examples")
	fmt.Println("7. 🎯 Run All Examples")
	fmt.Println("8. 🚪 Exit")
}

func getUserInput(prompt string) string {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func startHTTPServer() {
	fmt.Println("\n🌐 Starting HTTP Server...")

	port := getUserInput("Please enter port number (default 8080): ")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server will start at http://localhost:%s\n", port)
	fmt.Println("Press Ctrl+C to stop the server")

	http.StartServer(port)
}

func demoBasicClient() {
	fmt.Println("\n📡 Basic HTTP Client Examples")
	fmt.Println("=============================")

	http.ExampleBasicRequests()
}

func demoAdvancedClient() {
	fmt.Println("\n🔧 Advanced HTTP Client Examples")
	fmt.Println("=================================")

	http.ExampleAdvancedClient()
}

func demoGitHubAPI() {
	fmt.Println("\n🐙 GitHub API Examples")
	fmt.Println("======================")

	useAuth := getUserInput("Use GitHub authentication? (y/n): ")

	if strings.ToLower(useAuth) == "y" {
		token := getUserInput("Please enter GitHub Personal Access Token: ")
		if token != "" {
			http.ExampleGitHubWithAuth(token)
		} else {
			fmt.Println("No token provided, using unauthenticated mode")
			http.ExampleGitHubAPI()
		}
	} else {
		http.ExampleGitHubAPI()
	}

	username := getUserInput("\nEnter GitHub username for detailed info (or press Enter to skip): ")
	if username != "" {
		http.SimpleGitHubUserInfo(username)
	}
}

func demoHTTPUtils() {
	fmt.Println("\n🛠️  HTTP Utility Functions Examples")
	fmt.Println("===================================")

	http.ExampleHTTPUtils()
}

func demoJSONUtils() {
	fmt.Println("\n📄 JSON Utility Functions Examples")
	fmt.Println("===================================")

	http.ExampleJSONUtils()
}

func demoAllExamples() {
	fmt.Println("\n🎯 Running All Examples")
	fmt.Println("=======================")

	fmt.Println("1. Basic HTTP Client Examples")
	http.ExampleBasicRequests()

	time.Sleep(2 * time.Second)

	fmt.Println("\n2. Advanced HTTP Client Examples")
	http.ExampleAdvancedClient()

	time.Sleep(2 * time.Second)

	fmt.Println("\n3. GitHub API Examples")
	http.ExampleGitHubAPI()

	time.Sleep(2 * time.Second)

	fmt.Println("\n4. HTTP Utility Functions Examples")
	http.ExampleHTTPUtils()

	time.Sleep(2 * time.Second)

	fmt.Println("\n5. JSON Utility Functions Examples")
	http.ExampleJSONUtils()

	fmt.Println("\n✅ All examples completed!")
}

func interactiveServerDemo() {
	fmt.Println("\n🌐 Interactive Server Demo")
	fmt.Println("==========================")

	port := getUserInput("Please enter port number (default 8080): ")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server will start at http://localhost:%s\n", port)
	fmt.Println("Available endpoints:")
	fmt.Println("  GET  /           - Home page")
	fmt.Println("  GET  /health     - Health check")
	fmt.Println("  GET  /time       - Current time")
	fmt.Println("  GET  /users      - Users list (HTML)")
	fmt.Println("  GET  /users/{id} - User details (HTML)")
	fmt.Println("  GET  /api/users  - Users list (JSON)")
	fmt.Println("  GET  /api/users/{id} - User details (JSON)")

	start := getUserInput("Press Enter to start server...")
	if start == "" {
		http.StartServer(port)
	}
}

func customRequestDemo() {
	fmt.Println("\n🔧 Custom Request Demo")
	fmt.Println("======================")

	method := getUserInput("Request method (GET/POST/PUT/DELETE): ")
	url := getUserInput("Request URL: ")

	if method == "" || url == "" {
		fmt.Println("❌ Method and URL cannot be empty")
		return
	}

	options := http.RequestOptions{
		Method: strings.ToUpper(method),
		URL:    url,
		Headers: map[string]string{
			"User-Agent": "Go-HTTP-Demo/1.0",
		},
	}

	if method == "POST" || method == "PUT" {
		body := getUserInput("Request body (JSON format): ")
		if body != "" {
			options.Body = body
			options.Headers["Content-Type"] = "application/json"
		}
	}

	fmt.Println("\nSending request...")
	resp, err := http.MakeRequest(options)
	if err != nil {
		fmt.Printf("❌ Request failed: %v\n", err)
		return
	}

	fmt.Println("\n📥 Response result:")
	http.PrintResponse(resp)
}

func batchRequestDemo() {
	fmt.Println("\n📦 Batch Request Demo")
	fmt.Println("=====================")

	urls := []string{
		"https://httpbin.org/get",
		"https://httpbin.org/headers",
		"https://httpbin.org/user-agent",
		"https://httpbin.org/ip",
	}

	fmt.Printf("Will request %d URLs simultaneously:\n", len(urls))
	for i, url := range urls {
		fmt.Printf("  %d. %s\n", i+1, url)
	}

	start := getUserInput("\nPress Enter to start batch requests...")
	if start != "" {
		return
	}

	var requests []http.RequestOptions
	for _, url := range urls {
		requests = append(requests, http.RequestOptions{
			Method: "GET",
			URL:    url,
			Headers: map[string]string{
				"User-Agent": "Go-Batch-Demo/1.0",
			},
		})
	}

	fmt.Println("\nSending batch requests...")
	startTime := time.Now()
	responses, err := http.BatchRequest(requests)
	duration := time.Since(startTime)

	if err != nil {
		fmt.Printf("❌ Batch request failed: %v\n", err)
		return
	}

	fmt.Printf("\n✅ Batch request completed! Total time: %s\n", http.FormatDuration(duration))
	for i, resp := range responses {
		fmt.Printf("Request %d: Status %d, Time %s\n",
			i+1, resp.StatusCode, http.FormatDuration(resp.Duration))
	}
}

func retryRequestDemo() {
	fmt.Println("\n🔄 Retry Request Demo")
	fmt.Println("=====================")

	url := getUserInput("Request URL (suggest using a URL that will fail to demonstrate retry): ")
	if url == "" {
		url = "https://httpbin.org/status/500"
	}

	maxRetries := getUserInput("Maximum retry count (default 3): ")
	retries, err := strconv.Atoi(maxRetries)
	if err != nil || retries <= 0 {
		retries = 3
	}

	delay := getUserInput("Retry delay (seconds, default 1): ")
	delaySec, err := strconv.Atoi(delay)
	if err != nil || delaySec <= 0 {
		delaySec = 1
	}

	options := http.RequestOptions{
		Method: "GET",
		URL:    url,
		Headers: map[string]string{
			"User-Agent": "Go-Retry-Demo/1.0",
		},
	}

	fmt.Printf("\nStarting retry request (max %d times, delay %d seconds)...\n", retries, delaySec)

	startTime := time.Now()
	resp, err := http.RetryRequest(options, retries, time.Duration(delaySec)*time.Second)
	duration := time.Since(startTime)

	if err != nil {
		fmt.Printf("❌ Retry request finally failed: %v\n", err)
		fmt.Printf("Total time: %s\n", http.FormatDuration(duration))
	} else {
		fmt.Printf("✅ Retry request succeeded!\n")
		fmt.Printf("Status code: %d\n", resp.StatusCode)
		fmt.Printf("Total time: %s\n", http.FormatDuration(duration))
	}
}
