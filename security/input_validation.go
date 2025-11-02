package security

import (
	"encoding/json"
	"fmt"
	"html"
	"net/url"
	"regexp"
	"strings"
	"unicode"
)

// ValidationRule represents a validation rule
type ValidationRule struct {
	Required bool
	MinLen   int
	MaxLen   int
	Pattern  string
	Type     string
}

// ValidationResult represents the result of validation
type ValidationResult struct {
	Valid     bool
	Errors    []string
	Sanitized string
}

// InputValidator handles input validation and sanitization
type InputValidator struct {
	rules map[string]ValidationRule
}

// NewInputValidator creates a new input validator
func NewInputValidator() *InputValidator {
	return &InputValidator{
		rules: make(map[string]ValidationRule),
	}
}

// AddRule adds a validation rule for a field
func (v *InputValidator) AddRule(field string, rule ValidationRule) {
	v.rules[field] = rule
}

// ValidateString validates a string input
func (v *InputValidator) ValidateString(field, value string) ValidationResult {
	result := ValidationResult{
		Valid:     true,
		Errors:    []string{},
		Sanitized: value,
	}

	rule, exists := v.rules[field]
	if !exists {
		return result
	}

	// Check required
	if rule.Required && strings.TrimSpace(value) == "" {
		result.Valid = false
		result.Errors = append(result.Errors, fmt.Sprintf("%s is required", field))
		return result
	}

	// Skip further validation if empty and not required
	if strings.TrimSpace(value) == "" {
		return result
	}

	// Sanitize input
	sanitized := v.SanitizeString(value)
	result.Sanitized = sanitized

	// Check length
	if rule.MinLen > 0 && len(sanitized) < rule.MinLen {
		result.Valid = false
		result.Errors = append(result.Errors, fmt.Sprintf("%s must be at least %d characters", field, rule.MinLen))
	}

	if rule.MaxLen > 0 && len(sanitized) > rule.MaxLen {
		result.Valid = false
		result.Errors = append(result.Errors, fmt.Sprintf("%s must be at most %d characters", field, rule.MaxLen))
	}

	// Check pattern
	if rule.Pattern != "" {
		matched, err := regexp.MatchString(rule.Pattern, sanitized)
		if err != nil {
			result.Valid = false
			result.Errors = append(result.Errors, fmt.Sprintf("invalid pattern for %s", field))
		} else if !matched {
			result.Valid = false
			result.Errors = append(result.Errors, fmt.Sprintf("%s does not match required pattern", field))
		}
	}

	// Type-specific validation
	switch rule.Type {
	case "email":
		if !v.IsValidEmail(sanitized) {
			result.Valid = false
			result.Errors = append(result.Errors, fmt.Sprintf("%s must be a valid email address", field))
		}
	case "url":
		if !v.IsValidURL(sanitized) {
			result.Valid = false
			result.Errors = append(result.Errors, fmt.Sprintf("%s must be a valid URL", field))
		}
	case "alphanumeric":
		if !v.IsAlphanumeric(sanitized) {
			result.Valid = false
			result.Errors = append(result.Errors, fmt.Sprintf("%s must contain only alphanumeric characters", field))
		}
	}

	return result
}

// SanitizeString sanitizes a string to prevent XSS attacks
func (v *InputValidator) SanitizeString(input string) string {
	// HTML escape to prevent XSS
	sanitized := html.EscapeString(input)

	// Remove null bytes
	sanitized = strings.ReplaceAll(sanitized, "\x00", "")

	// Trim whitespace
	sanitized = strings.TrimSpace(sanitized)

	return sanitized
}

// SanitizeHTML sanitizes HTML content while preserving safe tags
func (v *InputValidator) SanitizeHTML(input string) string {
	// Remove script tags and their content
	scriptRegex := regexp.MustCompile(`(?i)<script[^>]*>.*?</script>`)
	sanitized := scriptRegex.ReplaceAllString(input, "")

	// Remove javascript: protocols
	jsRegex := regexp.MustCompile(`(?i)javascript:`)
	sanitized = jsRegex.ReplaceAllString(sanitized, "")

	// Remove on* event handlers
	eventRegex := regexp.MustCompile(`(?i)\s+on\w+\s*=\s*["'][^"']*["']`)
	sanitized = eventRegex.ReplaceAllString(sanitized, "")

	return sanitized
}

// IsValidEmail validates email format
func (v *InputValidator) IsValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// IsValidURL validates URL format
func (v *InputValidator) IsValidURL(urlStr string) bool {
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return false
	}
	return parsedURL.Scheme != "" && parsedURL.Host != ""
}

// IsAlphanumeric checks if string contains only alphanumeric characters
func (v *InputValidator) IsAlphanumeric(input string) bool {
	for _, char := range input {
		if !unicode.IsLetter(char) && !unicode.IsDigit(char) {
			return false
		}
	}
	return true
}

// ValidateJSON validates JSON input
func (v *InputValidator) ValidateJSON(jsonStr string) ValidationResult {
	result := ValidationResult{
		Valid:     true,
		Errors:    []string{},
		Sanitized: jsonStr,
	}

	// Check if it's valid JSON
	var jsonData interface{}
	if err := json.Unmarshal([]byte(jsonStr), &jsonData); err != nil {
		result.Valid = false
		result.Errors = append(result.Errors, "Invalid JSON format")
		return result
	}

	// Sanitize JSON by re-marshaling
	sanitizedBytes, err := json.Marshal(jsonData)
	if err != nil {
		result.Valid = false
		result.Errors = append(result.Errors, "Failed to sanitize JSON")
		return result
	}

	result.Sanitized = string(sanitizedBytes)
	return result
}

// PreventSQLInjection sanitizes input to prevent SQL injection
func (v *InputValidator) PreventSQLInjection(input string) string {
	// Remove or escape dangerous SQL characters
	dangerousChars := []string{
		"'", "\"", ";", "--", "/*", "*/", "xp_", "sp_",
		"exec", "execute", "select", "insert", "update", "delete",
		"drop", "create", "alter", "union", "or", "and",
	}

	sanitized := input
	for _, char := range dangerousChars {
		sanitized = strings.ReplaceAll(sanitized, char, "")
	}

	return sanitized
}

// ValidatePasswordStrength validates password strength
func (v *InputValidator) ValidatePasswordStrength(password string) ValidationResult {
	result := ValidationResult{
		Valid:     true,
		Errors:    []string{},
		Sanitized: password,
	}

	if len(password) < 8 {
		result.Valid = false
		result.Errors = append(result.Errors, "Password must be at least 8 characters long")
	}

	hasUpper := false
	hasLower := false
	hasDigit := false
	hasSpecial := false

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasDigit = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if !hasUpper {
		result.Valid = false
		result.Errors = append(result.Errors, "Password must contain at least one uppercase letter")
	}
	if !hasLower {
		result.Valid = false
		result.Errors = append(result.Errors, "Password must contain at least one lowercase letter")
	}
	if !hasDigit {
		result.Valid = false
		result.Errors = append(result.Errors, "Password must contain at least one digit")
	}
	if !hasSpecial {
		result.Valid = false
		result.Errors = append(result.Errors, "Password must contain at least one special character")
	}

	return result
}

// ValidateAndSanitizeMap validates and sanitizes a map of inputs
func (v *InputValidator) ValidateAndSanitizeMap(inputs map[string]string) (map[string]string, []string) {
	sanitized := make(map[string]string)
	var allErrors []string

	for field, value := range inputs {
		result := v.ValidateString(field, value)
		sanitized[field] = result.Sanitized
		allErrors = append(allErrors, result.Errors...)
	}

	return sanitized, allErrors
}

// CreateCommonRules creates common validation rules
func (v *InputValidator) CreateCommonRules() {
	v.AddRule("username", ValidationRule{
		Required: true,
		MinLen:   3,
		MaxLen:   20,
		Pattern:  `^[a-zA-Z0-9_]+$`,
		Type:     "alphanumeric",
	})

	v.AddRule("email", ValidationRule{
		Required: true,
		MinLen:   5,
		MaxLen:   100,
		Type:     "email",
	})

	v.AddRule("password", ValidationRule{
		Required: true,
		MinLen:   8,
		MaxLen:   128,
	})

	v.AddRule("name", ValidationRule{
		Required: true,
		MinLen:   1,
		MaxLen:   50,
		Pattern:  `^[a-zA-Z\s]+$`,
	})

	v.AddRule("url", ValidationRule{
		Required: false,
		MinLen:   0,
		MaxLen:   500,
		Type:     "url",
	})
}
