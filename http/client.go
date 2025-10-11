package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type HTTPClient struct {
	client  *http.Client
	baseURL string
	headers map[string]string
}

func NewHTTPClient(baseURL string) *HTTPClient {
	return &HTTPClient{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		baseURL: baseURL,
		headers: make(map[string]string),
	}
}

func NewHTTPClientWithTimeout(baseURL string, timeout time.Duration) *HTTPClient {
	return &HTTPClient{
		client: &http.Client{
			Timeout: timeout,
		},
		baseURL: baseURL,
		headers: make(map[string]string),
	}
}

func (c *HTTPClient) SetHeader(key, value string) {
	c.headers[key] = value
}

func (c *HTTPClient) SetHeaders(headers map[string]string) {
	for k, v := range headers {
		c.headers[k] = v
	}
}

func (c *HTTPClient) Get(path string) (*http.Response, error) {
	return c.request("GET", path, nil)
}

func (c *HTTPClient) Post(path string, body any) (*http.Response, error) {
	return c.requestWithJSON("POST", path, body)
}

func (c *HTTPClient) Put(path string, body any) (*http.Response, error) {
	return c.requestWithJSON("PUT", path, body)
}

func (c *HTTPClient) Delete(path string) (*http.Response, error) {
	return c.request("DELETE", path, nil)
}

func (c *HTTPClient) request(method, path string, body io.Reader) (*http.Response, error) {
	url := c.buildURL(path)

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	for key, value := range c.headers {
		req.Header.Set(key, value)
	}

	// Set default Content-Type if not set
	if req.Header.Get("Content-Type") == "" && body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return c.client.Do(req)
}

func (c *HTTPClient) requestWithJSON(method, path string, body any) (*http.Response, error) {
	var jsonBody []byte
	var err error

	if body != nil {
		jsonBody, err = json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal JSON: %w", err)
		}
	}

	return c.request(method, path, bytes.NewBuffer(jsonBody))
}

func (c *HTTPClient) buildURL(path string) string {
	if strings.HasPrefix(path, "http") {
		return path
	}

	baseURL := strings.TrimSuffix(c.baseURL, "/")
	path = strings.TrimPrefix(path, "/")

	return fmt.Sprintf("%s/%s", baseURL, path)
}

func (c *HTTPClient) GetJSON(path string, target any) error {
	resp, err := c.Get(path)
	if err != nil {
		return fmt.Errorf("GET request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("request failed with status: %d", resp.StatusCode)
	}

	return json.NewDecoder(resp.Body).Decode(target)
}

func (c *HTTPClient) PostJSON(path string, body, target any) error {
	resp, err := c.Post(path, body)
	if err != nil {
		return fmt.Errorf("POST request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("request failed with status: %d", resp.StatusCode)
	}

	if target != nil {
		return json.NewDecoder(resp.Body).Decode(target)
	}
	return nil
}

func SimpleGet(url string) (*http.Response, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	return client.Get(url)
}

func SimplePost(url string, body any) (*http.Response, error) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON: %w", err)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	return client.Post(url, "application/json", bytes.NewBuffer(jsonBody))
}

func SimplePostForm(url string, data map[string]string) (*http.Response, error) {
	formData := make([]string, 0, len(data))
	for k, v := range data {
		formData = append(formData, fmt.Sprintf("%s=%s", k, v))
	}

	client := &http.Client{Timeout: 10 * time.Second}
	return client.Post(url, "application/x-www-form-urlencoded",
		strings.NewReader(strings.Join(formData, "&")))
}

func DownloadFile(url string) ([]byte, error) {
	resp, err := SimpleGet(url)
	if err != nil {
		return nil, fmt.Errorf("failed to download file: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("download failed with status: %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}

func CheckURLStatus(url string) (int, error) {
	resp, err := SimpleGet(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	return resp.StatusCode, nil
}

func ExampleBasicRequests() {
	fmt.Println("=== Basic HTTP Client Examples ===")

	// Simple GET request
	fmt.Println("\n1. Simple GET request:")
	resp, err := SimpleGet("https://httpbin.org/get")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Status: %s\n", resp.Status)
		resp.Body.Close()
	}

	fmt.Println("\n2. Simple POST request:")
	data := map[string]any{
		"name": "John Doe",
		"age":  25,
	}
	resp, err = SimplePost("https://httpbin.org/post", data)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Status: %s\n", resp.Status)
		resp.Body.Close()
	}

	// Simple POST form
	fmt.Println("\n3. Simple POST form:")
	formData := map[string]string{
		"username": "testuser",
		"password": "testpass",
	}
	resp, err = SimplePostForm("https://httpbin.org/post", formData)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Status: %s\n", resp.Status)
		resp.Body.Close()
	}
}

func ExampleAdvancedClient() {
	fmt.Println("\n=== Advanced HTTP Client Examples ===")

	// Create client with base URL
	client := NewHTTPClient("https://httpbin.org")

	// Set headers
	client.SetHeaders(map[string]string{
		"User-Agent": "Go-HTTP-Client/1.0",
		"Accept":     "application/json",
	})

	// GET request
	fmt.Println("\n1. GET request with custom headers:")
	resp, err := client.Get("/get")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Status: %s\n", resp.Status)
		resp.Body.Close()
	}

	fmt.Println("\n2. POST request with JSON:")
	userData := map[string]any{
		"name":  "Jane Smith",
		"email": "jane@example.com",
	}
	resp, err = client.Post("/post", userData)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Status: %s\n", resp.Status)
		resp.Body.Close()
	}

	fmt.Println("\n3. GET request with JSON parsing:")
	var result map[string]any
	err = client.GetJSON("/json", &result)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Parsed JSON: %+v\n", result)
	}
}
