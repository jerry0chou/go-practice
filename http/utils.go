package http

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type ResponseData struct {
	StatusCode int                 `json:"status_code"`
	Headers    map[string][]string `json:"headers"`
	Body       string              `json:"body"`
	Duration   time.Duration       `json:"duration"`
}

type RequestOptions struct {
	Method  string
	URL     string
	Headers map[string]string
	Body    any
	Timeout time.Duration
}

func MakeRequest(options RequestOptions) (*ResponseData, error) {
	if options.Timeout == 0 {
		options.Timeout = 10 * time.Second
	}

	client := &http.Client{Timeout: options.Timeout}
	start := time.Now()
	var bodyReader io.Reader
	if options.Body != nil {
		if bodyStr, ok := options.Body.(string); ok {
			bodyReader = strings.NewReader(bodyStr)
		} else {
			jsonBody, err := json.Marshal(options.Body)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal body: %w", err)
			}
			bodyReader = strings.NewReader(string(jsonBody))
		}
	}

	req, err := http.NewRequest(options.Method, options.URL, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	for key, value := range options.Headers {
		req.Header.Set(key, value)
	}
	if req.Header.Get("Content-Type") == "" && options.Body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	duration := time.Since(start)

	return &ResponseData{
		StatusCode: resp.StatusCode,
		Headers:    resp.Header,
		Body:       string(bodyBytes),
		Duration:   duration,
	}, nil
}

func GetJSON(url string, target any) error {
	options := RequestOptions{
		Method: "GET",
		URL:    url,
		Headers: map[string]string{
			"Accept": "application/json",
		},
	}

	resp, err := MakeRequest(options)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("request failed with status: %d", resp.StatusCode)
	}

	return json.Unmarshal([]byte(resp.Body), target)
}

func PostJSON(url string, body, target any) error {
	options := RequestOptions{
		Method: "POST",
		URL:    url,
		Headers: map[string]string{
			"Content-Type": "application/json",
			"Accept":       "application/json",
		},
		Body: body,
	}

	resp, err := MakeRequest(options)
	if err != nil {
		return err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("request failed with status: %d", resp.StatusCode)
	}

	if target != nil {
		return json.Unmarshal([]byte(resp.Body), target)
	}
	return nil
}

func DownloadFileWithOptions(url string) ([]byte, error) {
	options := RequestOptions{
		Method: "GET",
		URL:    url,
	}

	resp, err := MakeRequest(options)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("download failed with status: %d", resp.StatusCode)
	}

	return []byte(resp.Body), nil
}

func CheckURLStatusWithOptions(url string) (int, error) {
	options := RequestOptions{
		Method:  "GET",
		URL:     url,
		Timeout: 5 * time.Second,
	}

	resp, err := MakeRequest(options)
	if err != nil {
		return 0, err
	}

	return resp.StatusCode, nil
}

func BatchRequest(requests []RequestOptions) ([]*ResponseData, error) {
	type result struct {
		index int
		data  *ResponseData
		err   error
	}

	results := make(chan result, len(requests))

	for i, req := range requests {
		go func(index int, request RequestOptions) {
			data, err := MakeRequest(request)
			results <- result{index: index, data: data, err: err}
		}(i, req)
	}

	responses := make([]*ResponseData, len(requests))
	var errors []error

	for i := 0; i < len(requests); i++ {
		res := <-results
		if res.err != nil {
			errors = append(errors, fmt.Errorf("request %d failed: %w", res.index, res.err))
		} else {
			responses[res.index] = res.data
		}
	}

	if len(errors) > 0 {
		return responses, fmt.Errorf("batch request had %d errors: %v", len(errors), errors)
	}

	return responses, nil
}

func RetryRequest(options RequestOptions, maxRetries int, delay time.Duration) (*ResponseData, error) {
	var lastErr error

	for i := 0; i <= maxRetries; i++ {
		resp, err := MakeRequest(options)
		if err == nil {
			return resp, nil
		}

		lastErr = err
		if i < maxRetries {
			time.Sleep(delay)
			delay *= 2
		}
	}

	return nil, fmt.Errorf("request failed after %d retries: %w", maxRetries, lastErr)
}

func ParseJSONResponse(body string) (map[string]any, error) {
	var result map[string]any
	err := json.Unmarshal([]byte(body), &result)
	return result, err
}

func FormatDuration(d time.Duration) string {
	if d < time.Microsecond {
		return fmt.Sprintf("%.0fns", float64(d.Nanoseconds()))
	} else if d < time.Millisecond {
		return fmt.Sprintf("%.2fμs", float64(d.Nanoseconds())/1000)
	} else if d < time.Second {
		return fmt.Sprintf("%.2fms", float64(d.Nanoseconds())/1000000)
	} else {
		return fmt.Sprintf("%.2fs", d.Seconds())
	}
}

func PrintResponse(resp *ResponseData) {
	fmt.Printf("Status: %d\n", resp.StatusCode)
	fmt.Printf("Duration: %s\n", FormatDuration(resp.Duration))
	fmt.Printf("Headers:\n")
	for key, values := range resp.Headers {
		fmt.Printf("  %s: %s\n", key, strings.Join(values, ", "))
	}
	fmt.Printf("Body:\n%s\n", resp.Body)
}

func ExampleHTTPUtils() {
	fmt.Println("=== HTTP Utils Examples ===")

	fmt.Println("\n1. Simple GET request:")
	options := RequestOptions{
		Method: "GET",
		URL:    "https://httpbin.org/get",
		Headers: map[string]string{
			"User-Agent": "Go-HTTP-Utils/1.0",
		},
	}

	resp, err := MakeRequest(options)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Status: %d, Duration: %s\n", resp.StatusCode, FormatDuration(resp.Duration))
	}

	fmt.Println("\n2. POST request with JSON:")
	postOptions := RequestOptions{
		Method: "POST",
		URL:    "https://httpbin.org/post",
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: map[string]any{
			"name": "Test User",
			"age":  25,
		},
	}

	resp, err = MakeRequest(postOptions)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Status: %d, Duration: %s\n", resp.StatusCode, FormatDuration(resp.Duration))
	}

	fmt.Println("\n3. Batch requests:")
	batchOptions := []RequestOptions{
		{Method: "GET", URL: "https://httpbin.org/get"},
		{Method: "GET", URL: "https://httpbin.org/headers"},
		{Method: "GET", URL: "https://httpbin.org/user-agent"},
	}

	responses, err := BatchRequest(batchOptions)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Batch completed: %d requests\n", len(responses))
		for i, resp := range responses {
			fmt.Printf("  Request %d: Status %d, Duration %s\n",
				i+1, resp.StatusCode, FormatDuration(resp.Duration))
		}
	}

	fmt.Println("\n4. Retry request (with invalid URL to demonstrate retry):")
	retryOptions := RequestOptions{
		Method: "GET",
		URL:    "https://httpbin.org/status/500", // This will return 500 status
	}

	resp, err = RetryRequest(retryOptions, 3, time.Second)
	if err != nil {
		fmt.Printf("Retry failed: %v\n", err)
	} else {
		fmt.Printf("Retry succeeded: Status %d\n", resp.StatusCode)
	}
}

func ExampleJSONUtils() {
	fmt.Println("\n=== JSON Utils Examples ===")

	fmt.Println("\n1. Get JSON data:")
	var jsonData map[string]any
	err := GetJSON("https://httpbin.org/json", &jsonData)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("JSON data: %+v\n", jsonData)
	}

	fmt.Println("\n2. Post JSON data:")
	postData := map[string]any{
		"message":   "Hello from Go!",
		"timestamp": time.Now().Format(time.RFC3339),
	}

	var response map[string]any
	err = PostJSON("https://httpbin.org/post", postData, &response)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Response received: %+v\n", response)
	}
}
