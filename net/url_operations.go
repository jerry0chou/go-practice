package net

import (
	"fmt"
	"net/url"
	"strings"
)

type URLInfo struct {
	Scheme   string
	Host     string
	Path     string
	Query    map[string][]string
	Fragment string
	Port     string
	User     *url.Userinfo
}

func ParseURL(rawURL string) (*URLInfo, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %w", err)
	}

	info := &URLInfo{
		Scheme:   parsedURL.Scheme,
		Host:     parsedURL.Host,
		Path:     parsedURL.Path,
		Query:    parsedURL.Query(),
		Fragment: parsedURL.Fragment,
		User:     parsedURL.User,
	}

	if parsedURL.Port() != "" {
		info.Port = parsedURL.Port()
	}

	return info, nil
}

func BuildURL(scheme, host, path string, queryParams map[string]string) string {
	u := &url.URL{
		Scheme: scheme,
		Host:   host,
		Path:   path,
	}

	if len(queryParams) > 0 {
		values := url.Values{}
		for key, value := range queryParams {
			values.Add(key, value)
		}
		u.RawQuery = values.Encode()
	}

	return u.String()
}

func AddQueryParams(baseURL string, params map[string]string) (string, error) {
	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return "", fmt.Errorf("failed to parse base URL: %w", err)
	}

	values := parsedURL.Query()
	for key, value := range params {
		values.Add(key, value)
	}
	parsedURL.RawQuery = values.Encode()

	return parsedURL.String(), nil
}

func EncodeQueryParams(params map[string]string) string {
	values := url.Values{}
	for key, value := range params {
		values.Add(key, value)
	}
	return values.Encode()
}

func DecodeQueryParams(queryString string) (map[string][]string, error) {
	values, err := url.ParseQuery(queryString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse query string: %w", err)
	}
	return values, nil
}

func URLEncode(text string) string {
	return url.QueryEscape(text)
}

func URLDecode(encodedText string) (string, error) {
	decoded, err := url.QueryUnescape(encodedText)
	if err != nil {
		return "", fmt.Errorf("failed to decode URL: %w", err)
	}
	return decoded, nil
}

func IsValidURL(rawURL string) bool {
	_, err := url.ParseRequestURI(rawURL)
	return err == nil
}

func GetDomainFromURL(rawURL string) (string, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "", fmt.Errorf("failed to parse URL: %w", err)
	}

	if parsedURL.Host == "" {
		return "", fmt.Errorf("no host found in URL")
	}

	host := parsedURL.Host
	if colonIndex := strings.Index(host, ":"); colonIndex != -1 {
		host = host[:colonIndex]
	}

	return host, nil
}

func ResolveRelativeURL(baseURL, relativeURL string) (string, error) {
	base, err := url.Parse(baseURL)
	if err != nil {
		return "", fmt.Errorf("failed to parse base URL: %w", err)
	}

	relative, err := url.Parse(relativeURL)
	if err != nil {
		return "", fmt.Errorf("failed to parse relative URL: %w", err)
	}

	resolved := base.ResolveReference(relative)
	return resolved.String(), nil
}

func PrintURLInfo(rawURL string) {
	fmt.Printf("🔗 Analyzing URL: %s\n", rawURL)
	fmt.Println(strings.Repeat("-", 50))

	info, err := ParseURL(rawURL)
	if err != nil {
		fmt.Printf("❌ Error parsing URL: %v\n", err)
		return
	}

	fmt.Printf("📋 URL Components:\n")
	fmt.Printf("   Scheme:   %s\n", info.Scheme)
	fmt.Printf("   Host:     %s\n", info.Host)
	fmt.Printf("   Port:     %s\n", info.Port)
	fmt.Printf("   Path:     %s\n", info.Path)
	fmt.Printf("   Fragment: %s\n", info.Fragment)

	if info.User != nil {
		fmt.Printf("   User:     %s\n", info.User.Username())
		if password, ok := info.User.Password(); ok {
			fmt.Printf("   Password: %s\n", strings.Repeat("*", len(password)))
		}
	}

	if len(info.Query) > 0 {
		fmt.Printf("   Query Parameters:\n")
		for key, values := range info.Query {
			for _, value := range values {
				fmt.Printf("     %s = %s\n", key, value)
			}
		}
	}

	fmt.Println()
}

func DemonstrateURLOperations() {
	fmt.Println("🌐 URL Operations Demo")
	fmt.Println(strings.Repeat("=", 60))

	urls := []string{
		"https://www.example.com:8080/path/to/resource?param1=value1&param2=value2#section",
		"http://user:pass@example.com:3000/api/users?id=123&name=john",
		"ftp://files.example.com/downloads/file.txt",
		"mailto:user@example.com?subject=Hello&body=World",
		"/relative/path?query=value",
		"invalid-url",
	}

	for _, rawURL := range urls {
		PrintURLInfo(rawURL)
	}

	fmt.Println("🔨 URL Building Examples:")
	fmt.Println(strings.Repeat("-", 30))

	builtURL := BuildURL("https", "api.example.com", "/v1/users", map[string]string{
		"page":  "1",
		"limit": "10",
		"sort":  "name",
	})
	fmt.Printf("Built URL: %s\n", builtURL)

	baseURL := "https://api.example.com/v1/users"
	modifiedURL, err := AddQueryParams(baseURL, map[string]string{
		"filter": "active",
		"role":   "admin",
	})
	if err != nil {
		fmt.Printf("Error adding query params: %v\n", err)
	} else {
		fmt.Printf("Modified URL: %s\n", modifiedURL)
	}

	fmt.Println("\n🔐 URL Encoding/Decoding:")
	fmt.Println(strings.Repeat("-", 30))

	text := "Hello World! @#$%"
	encoded := URLEncode(text)
	fmt.Printf("Original: %s\n", text)
	fmt.Printf("Encoded:  %s\n", encoded)

	decoded, err := URLDecode(encoded)
	if err != nil {
		fmt.Printf("Error decoding: %v\n", err)
	} else {
		fmt.Printf("Decoded:  %s\n", decoded)
	}

	fmt.Println("\n📝 Query Parameter Operations:")
	fmt.Println(strings.Repeat("-", 30))

	params := map[string]string{
		"search": "golang tutorial",
		"type":   "programming",
		"level":  "beginner",
	}
	encodedQuery := EncodeQueryParams(params)
	fmt.Printf("Encoded Query: %s\n", encodedQuery)

	decodedParams, err := DecodeQueryParams(encodedQuery)
	if err != nil {
		fmt.Printf("Error decoding query: %v\n", err)
	} else {
		fmt.Printf("Decoded Query:\n")
		for key, values := range decodedParams {
			for _, value := range values {
				fmt.Printf("  %s = %s\n", key, value)
			}
		}
	}

	fmt.Println("\n✅ URL Validation:")
	fmt.Println(strings.Repeat("-", 30))

	testURLs := []string{
		"https://www.google.com",
		"http://localhost:8080",
		"ftp://files.example.com",
		"not-a-url",
		"",
	}

	for _, testURL := range testURLs {
		valid := IsValidURL(testURL)
		status := "❌ Invalid"
		if valid {
			status = "✅ Valid"
		}
		fmt.Printf("%s: %s\n", status, testURL)
	}

	fmt.Println("\n🌍 Domain Extraction:")
	fmt.Println(strings.Repeat("-", 30))

	domainURLs := []string{
		"https://www.example.com:8080/path",
		"http://subdomain.example.org",
		"https://api.service.co.uk/v1",
	}

	for _, domainURL := range domainURLs {
		domain, err := GetDomainFromURL(domainURL)
		if err != nil {
			fmt.Printf("Error extracting domain from %s: %v\n", domainURL, err)
		} else {
			fmt.Printf("Domain from %s: %s\n", domainURL, domain)
		}
	}

	fmt.Println("\n🔗 Relative URL Resolution:")
	fmt.Println(strings.Repeat("-", 30))

	baseURLs := []string{
		"https://www.example.com/api/v1/",
		"http://localhost:8080/",
	}

	relativeURLs := []string{
		"users",
		"../v2/users",
		"/absolute/path",
		"../../other/path",
	}

	for _, base := range baseURLs {
		fmt.Printf("Base URL: %s\n", base)
		for _, relative := range relativeURLs {
			resolved, err := ResolveRelativeURL(base, relative)
			if err != nil {
				fmt.Printf("  Error resolving %s: %v\n", relative, err)
			} else {
				fmt.Printf("  %s -> %s\n", relative, resolved)
			}
		}
		fmt.Println()
	}
}
