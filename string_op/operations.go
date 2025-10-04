package string_op

import (
	"encoding/base64"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

func BasicOperations() {
	fmt.Println("=== Basic String Operations ===")

	text := "  Hello, World!  "
	fmt.Printf("Original: '%s'\n", text)
	fmt.Printf("Length: %d characters\n", len(text))
	fmt.Printf("UTF-8 length: %d runes\n", utf8.RuneCountInString(text))

	trimmed := strings.TrimSpace(text)
	fmt.Printf("Trimmed: '%s'\n", trimmed)

	upper := strings.ToUpper(text)
	lower := strings.ToLower(text)
	title := strings.Title(text)
	fmt.Printf("Uppercase: '%s'\n", upper)
	fmt.Printf("Lowercase: '%s'\n", lower)
	fmt.Printf("Title case: '%s'\n", title)

	repeated := strings.Repeat("Go ", 3)
	fmt.Printf("Repeated: '%s'\n", repeated)

	hasPrefix := strings.HasPrefix(text, "  ")
	hasSuffix := strings.HasSuffix(text, "  ")
	fmt.Printf("Has prefix '  ': %t\n", hasPrefix)
	fmt.Printf("Has suffix '  ': %t\n", hasSuffix)
}

func SearchOperations() {
	fmt.Println("\n=== String Search Operations ===")

	text := "Go is awesome! Go is fast!"
	fmt.Printf("Text: '%s'\n", text)

	contains := strings.Contains(text, "Go")
	fmt.Printf("Contains 'Go': %t\n", contains)

	index := strings.Index(text, "Go")
	lastIndex := strings.LastIndex(text, "Go")
	fmt.Printf("First 'Go' at index: %d\n", index)
	fmt.Printf("Last 'Go' at index: %d\n", lastIndex)

	count := strings.Count(text, "Go")
	fmt.Printf("Count of 'Go': %d\n", count)

	fields := strings.Fields(text)
	fmt.Printf("Fields: %v\n", fields)

	words := strings.Split(text, " ")
	fmt.Printf("Split by space: %v\n", words)

	joined := strings.Join(words, "-")
	fmt.Printf("Joined with '-': %s\n", joined)
}

func ReplaceOperations() {
	fmt.Println("\n=== String Replace Operations ===")

	text := "Hello World! Hello Go!"
	fmt.Printf("Original: '%s'\n", text)

	replaced := strings.Replace(text, "Hello", "Hi", 1)
	fmt.Printf("Replace first 'Hello': '%s'\n", replaced)

	replacedAll := strings.ReplaceAll(text, "Hello", "Hi")
	fmt.Printf("Replace all 'Hello': '%s'\n", replacedAll)

	replacer := strings.NewReplacer("Hello", "Hi", "World", "Universe")
	replacedMultiple := replacer.Replace(text)
	fmt.Printf("Multiple replacements: '%s'\n", replacedMultiple)

	trimmed := strings.Trim(text, "!")
	fmt.Printf("Trim '!': '%s'\n", trimmed)

	trimmedLeft := strings.TrimLeft(text, "H")
	trimmedRight := strings.TrimRight(text, "!")
	fmt.Printf("Trim left 'H': '%s'\n", trimmedLeft)
	fmt.Printf("Trim right '!': '%s'\n", trimmedRight)
}

func ValidationOperations() {
	fmt.Println("\n=== String Validation Operations ===")

	texts := []string{"Hello123", "hello", "HELLO", "123", "Hello World", "hello@world.com"}

	for _, text := range texts {
		fmt.Printf("\nText: '%s'\n", text)
		fmt.Printf("  Is empty: %t\n", text == "")
		fmt.Printf("  Is numeric: %t\n", isNumeric(text))
		fmt.Printf("  Is alphabetic: %t\n", isAlphabetic(text))
		fmt.Printf("  Is alphanumeric: %t\n", isAlphanumeric(text))
		fmt.Printf("  Is email: %t\n", isValidEmail(text))
		fmt.Printf("  Has uppercase: %t\n", hasUppercase(text))
		fmt.Printf("  Has lowercase: %t\n", hasLowercase(text))
		fmt.Printf("  Has digit: %t\n", hasDigit(text))
	}
}

func ParsingOperations() {
	fmt.Println("\n=== String Parsing Operations ===")

	numbers := []string{"123", "45.67", "true", "false", "invalid"}
	for _, num := range numbers {
		if intVal, err := strconv.Atoi(num); err == nil {
			fmt.Printf("'%s' -> int: %d\n", num, intVal)
		}
		if floatVal, err := strconv.ParseFloat(num, 64); err == nil {
			fmt.Printf("'%s' -> float: %.2f\n", num, floatVal)
		}
		if boolVal, err := strconv.ParseBool(num); err == nil {
			fmt.Printf("'%s' -> bool: %t\n", num, boolVal)
		}
	}

	text := "name=John,age=30,city=New York"
	pairs := strings.Split(text, ",")
	fmt.Printf("\nParsing key-value pairs from: '%s'\n", text)
	for _, pair := range pairs {
		parts := strings.Split(pair, "=")
		if len(parts) == 2 {
			fmt.Printf("  %s: %s\n", parts[0], parts[1])
		}
	}
}

func EncodingOperations() {
	fmt.Println("\n=== String Encoding Operations ===")

	text := "Hello, 世界!"
	fmt.Printf("Original: '%s'\n", text)

	base64Encoded := base64.StdEncoding.EncodeToString([]byte(text))
	fmt.Printf("Base64 encoded: %s\n", base64Encoded)

	base64Decoded, err := base64.StdEncoding.DecodeString(base64Encoded)
	if err == nil {
		fmt.Printf("Base64 decoded: '%s'\n", string(base64Decoded))
	}

	urlEncoded := strings.ReplaceAll(text, " ", "%20")
	fmt.Printf("URL encoded: %s\n", urlEncoded)

	htmlEscaped := strings.ReplaceAll(text, "<", "&lt;")
	htmlEscaped = strings.ReplaceAll(htmlEscaped, ">", "&gt;")
	fmt.Printf("HTML escaped: %s\n", htmlEscaped)
}

func AdvancedOperations() {
	fmt.Println("\n=== Advanced String Operations ===")

	text := "Go is awesome! Go is fast! Go is simple!"
	fmt.Printf("Text: '%s'\n", text)

	words := strings.Fields(text)
	wordCount := make(map[string]int)
	for _, word := range words {
		word = strings.Trim(word, "!")
		wordCount[word]++
	}
	fmt.Printf("Word frequency: %v\n", wordCount)

	longest := findLongestWord(text)
	fmt.Printf("Longest word: '%s'\n", longest)

	reversed := reverseString(text)
	fmt.Printf("Reversed: '%s'\n", reversed)

	palindrome := "racecar"
	fmt.Printf("'%s' is palindrome: %t\n", palindrome, isPalindrome(palindrome))

	anagram1, anagram2 := "listen", "silent"
	fmt.Printf("'%s' and '%s' are anagrams: %t\n", anagram1, anagram2, areAnagrams(anagram1, anagram2))
}

func RegularExpressionOperations() {
	fmt.Println("\n=== Regular Expression Operations ===")

	text := "Contact us at support@example.com or call +1-555-123-4567"
	fmt.Printf("Text: '%s'\n", text)

	emailPattern := `\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Z|a-z]{2,}\b`
	emailRegex, _ := regexp.Compile(emailPattern)
	emails := emailRegex.FindAllString(text, -1)
	fmt.Printf("Found emails: %v\n", emails)

	phonePattern := `\+?1?[-.\s]?\(?([0-9]{3})\)?[-.\s]?([0-9]{3})[-.\s]?([0-9]{4})`
	phoneRegex, _ := regexp.Compile(phonePattern)
	phones := phoneRegex.FindAllString(text, -1)
	fmt.Printf("Found phones: %v\n", phones)

	replaced := emailRegex.ReplaceAllString(text, "[EMAIL]")
	fmt.Printf("Replaced emails: '%s'\n", replaced)
}

func UtilityOperations() {
	fmt.Println("\n=== Utility String Operations ===")

	text := "  Hello, World!  "
	fmt.Printf("Original: '%s'\n", text)

	normalized := normalizeString(text)
	fmt.Printf("Normalized: '%s'\n", normalized)

	slug := createSlug("Hello, World! This is a test.")
	fmt.Printf("Slug: '%s'\n", slug)

	truncated := truncateString(text, 10)
	fmt.Printf("Truncated to 10 chars: '%s'\n", truncated)

	padded := padString("Go", 10, " ")
	fmt.Printf("Padded to 10 chars: '%s'\n", padded)

	indented := indentString("Line 1\nLine 2\nLine 3", "  ")
	fmt.Printf("Indented:\n%s\n", indented)
}

func isNumeric(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}

func isAlphabetic(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) {
			return false
		}
	}
	return len(s) > 0
}

func isAlphanumeric(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			return false
		}
	}
	return len(s) > 0
}

func isValidEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(pattern, email)
	return matched
}

func hasUppercase(s string) bool {
	for _, r := range s {
		if unicode.IsUpper(r) {
			return true
		}
	}
	return false
}

func hasLowercase(s string) bool {
	for _, r := range s {
		if unicode.IsLower(r) {
			return true
		}
	}
	return false
}

func hasDigit(s string) bool {
	for _, r := range s {
		if unicode.IsDigit(r) {
			return true
		}
	}
	return false
}

func findLongestWord(text string) string {
	words := strings.Fields(text)
	longest := ""
	for _, word := range words {
		word = strings.Trim(word, "!.,?")
		if len(word) > len(longest) {
			longest = word
		}
	}
	return longest
}

func reverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func isPalindrome(s string) bool {
	s = strings.ToLower(strings.ReplaceAll(s, " ", ""))
	return s == reverseString(s)
}

func areAnagrams(s1, s2 string) bool {
	if len(s1) != len(s2) {
		return false
	}

	count1 := make(map[rune]int)
	count2 := make(map[rune]int)

	for _, r := range strings.ToLower(s1) {
		count1[r]++
	}
	for _, r := range strings.ToLower(s2) {
		count2[r]++
	}

	for r, count := range count1 {
		if count2[r] != count {
			return false
		}
	}
	return true
}

func normalizeString(s string) string {
	s = strings.TrimSpace(s)
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, "  ", " ")
	return s
}

func createSlug(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", "-")
	s = strings.ReplaceAll(s, ",", "")
	s = strings.ReplaceAll(s, "!", "")
	s = strings.ReplaceAll(s, ".", "")
	return s
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

func padString(s string, length int, pad string) string {
	if len(s) >= length {
		return s
	}
	padding := strings.Repeat(pad, length-len(s))
	return s + padding
}

func indentString(s, indent string) string {
	lines := strings.Split(s, "\n")
	for i, line := range lines {
		lines[i] = indent + line
	}
	return strings.Join(lines, "\n")
}

func RunAllStringExamples() {
	fmt.Println("🎯 Go String Operations Examples")
	fmt.Println("=================================")

	BasicOperations()
	SearchOperations()
	ReplaceOperations()
	ValidationOperations()
	ParsingOperations()
	EncodingOperations()
	AdvancedOperations()
	RegularExpressionOperations()
	UtilityOperations()

	fmt.Println("\n✅ All string operations completed!")
}
