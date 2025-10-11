package http

import (
	"fmt"
	"log"
	"time"
)

func RunAllExamples() {
	fmt.Println("🚀 Go HTTP Package - All Examples")
	fmt.Println("=================================")

	fmt.Println("\n1. 📡 Basic HTTP Client Examples")
	fmt.Println("--------------------------------")
	ExampleBasicRequests()

	time.Sleep(1 * time.Second)

	fmt.Println("\n2. 🔧 Advanced HTTP Client Examples")
	fmt.Println("-----------------------------------")
	ExampleAdvancedClient()

	time.Sleep(1 * time.Second)

	fmt.Println("\n3. 🐙 GitHub API Examples")
	fmt.Println("-------------------------")
	ExampleGitHubAPI()

	time.Sleep(1 * time.Second)

	fmt.Println("\n4. 🛠️  HTTP Utils Examples")
	fmt.Println("--------------------------")
	ExampleHTTPUtils()

	time.Sleep(1 * time.Second)

	fmt.Println("\n5. 📄 JSON Utils Examples")
	fmt.Println("-------------------------")
	ExampleJSONUtils()

	fmt.Println("\n✅ All examples completed!")
}

func QuickStartServer(port string) {
	if port == "" {
		port = "8080"
	}

	fmt.Printf("🚀 Starting HTTP server on port %s\n", port)
	fmt.Printf("📋 Available endpoints:\n")
	fmt.Printf("   http://localhost:%s/           - Home page\n", port)
	fmt.Printf("   http://localhost:%s/health     - Health check\n", port)
	fmt.Printf("   http://localhost:%s/time       - Current time\n", port)
	fmt.Printf("   http://localhost:%s/users      - Users list\n", port)
	fmt.Printf("   http://localhost:%s/api/users  - Users API\n", port)
	fmt.Printf("\nPress Ctrl+C to stop the server\n\n")

	StartServer(port)
}

func QuickTestClient() {
	fmt.Println("🧪 Quick HTTP Client Test")
	fmt.Println("=========================")

	fmt.Println("\n1. Testing basic GET request...")
	resp, err := SimpleGet("https://httpbin.org/get")
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		fmt.Printf("✅ GET request successful: %s\n", resp.Status)
		resp.Body.Close()
	}

	fmt.Println("\n2. Testing basic POST request...")
	data := map[string]any{
		"message":   "Hello from Go HTTP client!",
		"timestamp": time.Now().Format(time.RFC3339),
	}
	resp, err = SimplePost("https://httpbin.org/post", data)
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		fmt.Printf("✅ POST request successful: %s\n", resp.Status)
		resp.Body.Close()
	}

	fmt.Println("\n3. Testing GitHub API...")
	client := NewGitHubClientWithoutAuth()
	user, err := client.GetUser("octocat")
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		fmt.Printf("✅ GitHub API successful: %s (%s)\n", user.Name, user.Login)
	}

	fmt.Println("\n🎉 Quick test completed!")
}

func BenchmarkExample() {
	fmt.Println("⚡ HTTP Performance Benchmark")
	fmt.Println("=============================")

	urls := []string{
		"https://httpbin.org/get",
		"https://httpbin.org/headers",
		"https://httpbin.org/user-agent",
		"https://httpbin.org/ip",
		"https://httpbin.org/json",
	}

	fmt.Println("\n1. Sequential requests...")
	start := time.Now()
	for i, url := range urls {
		resp, err := SimpleGet(url)
		if err != nil {
			log.Printf("Request %d failed: %v", i+1, err)
		} else {
			resp.Body.Close()
		}
	}
	sequentialTime := time.Since(start)
	fmt.Printf("Sequential time: %s\n", FormatDuration(sequentialTime))

	fmt.Println("\n2. Concurrent requests...")
	start = time.Now()
	var requests []RequestOptions
	for _, url := range urls {
		requests = append(requests, RequestOptions{
			Method: "GET",
			URL:    url,
		})
	}
	responses, err := BatchRequest(requests)
	concurrentTime := time.Since(start)

	if err != nil {
		log.Printf("Concurrent requests failed: %v", err)
	} else {
		fmt.Printf("Concurrent time: %s\n", FormatDuration(concurrentTime))
		fmt.Printf("Speedup: %.2fx\n", float64(sequentialTime)/float64(concurrentTime))
		fmt.Printf("Received %d responses\n", len(responses))
	}
}

func ErrorHandlingExample() {
	fmt.Println("❌ HTTP Error Handling Example")
	fmt.Println("==============================")

	testCases := []struct {
		name string
		url  string
	}{
		{"Valid URL", "https://httpbin.org/get"},
		{"Invalid URL", "https://invalid-url-that-does-not-exist.com"},
		{"404 Not Found", "https://httpbin.org/status/404"},
		{"500 Server Error", "https://httpbin.org/status/500"},
		{"Timeout", "https://httpbin.org/delay/10"},
	}

	for _, tc := range testCases {
		fmt.Printf("\nTesting: %s\n", tc.name)

		options := RequestOptions{
			Method:  "GET",
			URL:     tc.url,
			Timeout: 5 * time.Second, // Short timeout for demo
		}

		resp, err := MakeRequest(options)
		if err != nil {
			fmt.Printf("❌ Error: %v\n", err)
		} else {
			fmt.Printf("✅ Status: %d\n", resp.StatusCode)
		}
	}
}

func MiddlewareExample() {
	fmt.Println("🔧 Middleware Example")
	fmt.Println("=====================")

	fmt.Println("Available middleware:")
	fmt.Println("  - loggingMiddleware: Logs all requests")
	fmt.Println("  - corsMiddleware: Adds CORS headers")
	fmt.Println("  - authMiddleware: Basic authentication")
	fmt.Println("  - rateLimitMiddleware: Rate limiting")
	fmt.Println("  - recoveryMiddleware: Panic recovery")
	fmt.Println("  - customHeadersMiddleware: Custom headers")

	fmt.Println("\nTo use middleware, create a server with:")
	fmt.Println("  handler := ChainMiddleware(mux, middleware1, middleware2, ...)")
	fmt.Println("  server := &http.Server{Handler: handler}")
}

func RealWorldExample() {
	fmt.Println("🌍 Real-World HTTP Usage Example")
	fmt.Println("================================")

	fmt.Println("\n1. API Health Check")
	status, err := CheckURLStatus("https://api.github.com")
	if err != nil {
		fmt.Printf("❌ GitHub API is down: %v\n", err)
	} else {
		fmt.Printf("✅ GitHub API is up: %d\n", status)
	}

	fmt.Println("\n2. Download and Parse JSON")
	var jsonData map[string]any
	err = GetJSON("https://httpbin.org/json", &jsonData)
	if err != nil {
		fmt.Printf("❌ Failed to get JSON: %v\n", err)
	} else {
		fmt.Printf("✅ JSON downloaded: %+v\n", jsonData)
	}

	fmt.Println("\n3. Retry with Exponential Backoff")
	options := RequestOptions{
		Method: "GET",
		URL:    "https://httpbin.org/status/500",
	}

	resp, err := RetryRequest(options, 3, time.Second)
	if err != nil {
		fmt.Printf("❌ Retry failed: %v\n", err)
	} else {
		fmt.Printf("✅ Retry succeeded: %d\n", resp.StatusCode)
	}

	fmt.Println("\n4. Batch Processing")
	urls := []string{
		"https://httpbin.org/get",
		"https://httpbin.org/headers",
		"https://httpbin.org/user-agent",
	}

	var requests []RequestOptions
	for _, url := range urls {
		requests = append(requests, RequestOptions{
			Method: "GET",
			URL:    url,
		})
	}

	responses, err := BatchRequest(requests)
	if err != nil {
		fmt.Printf("❌ Batch request failed: %v\n", err)
	} else {
		fmt.Printf("✅ Batch request succeeded: %d responses\n", len(responses))
	}
}
