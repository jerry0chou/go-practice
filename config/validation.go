package config

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

// ValidationError represents a configuration validation error
type ValidationError struct {
	Field   string
	Value   interface{}
	Message string
}

func (ve ValidationError) Error() string {
	return fmt.Sprintf("validation error for field '%s': %s (value: %v)", ve.Field, ve.Message, ve.Value)
}

// ValidationRule defines a validation rule for a configuration field
type ValidationRule struct {
	Field    string
	Required bool
	Type     reflect.Type
	Min      *float64
	Max      *float64
	Pattern  *regexp.Regexp
	Enum     []interface{}
	Custom   func(interface{}) error
}

// SchemaValidator provides configuration schema validation
type SchemaValidator struct {
	rules map[string]ValidationRule
}

// NewSchemaValidator creates a new schema validator
func NewSchemaValidator() *SchemaValidator {
	return &SchemaValidator{
		rules: make(map[string]ValidationRule),
	}
}

// AddRule adds a validation rule for a field
func (sv *SchemaValidator) AddRule(rule ValidationRule) {
	sv.rules[rule.Field] = rule
}

// Validate validates a configuration struct against the schema
func (sv *SchemaValidator) Validate(config interface{}) error {
	val := reflect.ValueOf(config)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return fmt.Errorf("configuration must be a struct")
	}

	var errors []ValidationError

	// Validate each field
	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		fieldValue := val.Field(i)

		// Get field path (e.g., "app.name", "server.port")
		fieldPath := sv.getFieldPath(field, "")

		// Check if there's a validation rule for this field
		if rule, exists := sv.rules[fieldPath]; exists {
			if err := sv.validateField(fieldValue, rule); err != nil {
				errors = append(errors, ValidationError{
					Field:   fieldPath,
					Value:   fieldValue.Interface(),
					Message: err.Error(),
				})
			}
		}

		// Recursively validate nested structs
		if fieldValue.Kind() == reflect.Struct {
			if err := sv.validateNestedStruct(fieldValue, fieldPath); err != nil {
				errors = append(errors, err...)
			}
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("validation failed: %v", errors)
	}

	return nil
}

// validateField validates a single field against its rule
func (sv *SchemaValidator) validateField(fieldValue reflect.Value, rule ValidationRule) error {
	// Check if field is required and empty
	if rule.Required && sv.isEmpty(fieldValue) {
		return fmt.Errorf("field is required")
	}

	// Skip validation if field is empty and not required
	if sv.isEmpty(fieldValue) {
		return nil
	}

	// Check type
	if rule.Type != nil && fieldValue.Type() != rule.Type {
		return fmt.Errorf("expected type %s, got %s", rule.Type, fieldValue.Type())
	}

	// Check numeric constraints
	if rule.Min != nil || rule.Max != nil {
		if err := sv.validateNumeric(fieldValue, rule); err != nil {
			return err
		}
	}

	// Check pattern (for strings)
	if rule.Pattern != nil {
		if fieldValue.Kind() == reflect.String {
			if !rule.Pattern.MatchString(fieldValue.String()) {
				return fmt.Errorf("value does not match required pattern")
			}
		}
	}

	// Check enum values
	if len(rule.Enum) > 0 {
		if err := sv.validateEnum(fieldValue, rule.Enum); err != nil {
			return err
		}
	}

	// Custom validation
	if rule.Custom != nil {
		if err := rule.Custom(fieldValue.Interface()); err != nil {
			return err
		}
	}

	return nil
}

// validateNestedStruct validates nested struct fields
func (sv *SchemaValidator) validateNestedStruct(structValue reflect.Value, parentPath string) []ValidationError {
	var errors []ValidationError

	for i := 0; i < structValue.NumField(); i++ {
		field := structValue.Type().Field(i)
		fieldValue := structValue.Field(i)

		fieldPath := sv.getFieldPath(field, parentPath)

		if rule, exists := sv.rules[fieldPath]; exists {
			if err := sv.validateField(fieldValue, rule); err != nil {
				errors = append(errors, ValidationError{
					Field:   fieldPath,
					Value:   fieldValue.Interface(),
					Message: err.Error(),
				})
			}
		}

		// Recursively validate nested structs
		if fieldValue.Kind() == reflect.Struct {
			if nestedErrors := sv.validateNestedStruct(fieldValue, fieldPath); len(nestedErrors) > 0 {
				errors = append(errors, nestedErrors...)
			}
		}
	}

	return errors
}

// validateNumeric validates numeric constraints
func (sv *SchemaValidator) validateNumeric(fieldValue reflect.Value, rule ValidationRule) error {
	var num float64

	switch fieldValue.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		num = float64(fieldValue.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		num = float64(fieldValue.Uint())
	case reflect.Float32, reflect.Float64:
		num = fieldValue.Float()
	default:
		return fmt.Errorf("field is not numeric")
	}

	if rule.Min != nil && num < *rule.Min {
		return fmt.Errorf("value %v is less than minimum %v", num, *rule.Min)
	}

	if rule.Max != nil && num > *rule.Max {
		return fmt.Errorf("value %v is greater than maximum %v", num, *rule.Max)
	}

	return nil
}

// validateEnum validates enum constraints
func (sv *SchemaValidator) validateEnum(fieldValue reflect.Value, enum []interface{}) error {
	value := fieldValue.Interface()

	for _, validValue := range enum {
		if reflect.DeepEqual(value, validValue) {
			return nil
		}
	}

	return fmt.Errorf("value %v is not in allowed values: %v", value, enum)
}

// isEmpty checks if a field value is empty
func (sv *SchemaValidator) isEmpty(fieldValue reflect.Value) bool {
	switch fieldValue.Kind() {
	case reflect.String:
		return fieldValue.String() == ""
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return fieldValue.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return fieldValue.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return fieldValue.Float() == 0
	case reflect.Bool:
		return !fieldValue.Bool()
	case reflect.Slice, reflect.Map, reflect.Array:
		return fieldValue.Len() == 0
	case reflect.Ptr, reflect.Interface:
		return fieldValue.IsNil()
	}
	return false
}

// getFieldPath constructs the field path for nested structs
func (sv *SchemaValidator) getFieldPath(field reflect.StructField, parentPath string) string {
	fieldName := strings.ToLower(field.Name)

	if parentPath == "" {
		return fieldName
	}

	return parentPath + "." + fieldName
}

// CreateDefaultSchema creates a default validation schema for FileConfig
func CreateDefaultSchema() *SchemaValidator {
	validator := NewSchemaValidator()

	// App configuration rules
	validator.AddRule(ValidationRule{
		Field:    "app.name",
		Required: true,
		Type:     reflect.TypeOf(""),
		Custom: func(value interface{}) error {
			str, ok := value.(string)
			if !ok {
				return fmt.Errorf("must be a string")
			}
			if len(str) < 1 || len(str) > 50 {
				return fmt.Errorf("must be between 1 and 50 characters")
			}
			return nil
		},
	})

	validator.AddRule(ValidationRule{
		Field:    "app.version",
		Required: true,
		Type:     reflect.TypeOf(""),
		Pattern:  regexp.MustCompile(`^\d+\.\d+\.\d+$`),
	})

	validator.AddRule(ValidationRule{
		Field:    "app.environment",
		Required: true,
		Type:     reflect.TypeOf(""),
		Enum:     []interface{}{"development", "staging", "production"},
	})

	// Server configuration rules
	validator.AddRule(ValidationRule{
		Field:    "server.port",
		Required: true,
		Type:     reflect.TypeOf(0),
		Min:      floatPtr(1),
		Max:      floatPtr(65535),
	})

	validator.AddRule(ValidationRule{
		Field:    "server.host",
		Required: true,
		Type:     reflect.TypeOf(""),
		Pattern:  regexp.MustCompile(`^[a-zA-Z0-9.-]+$`),
	})

	// Database configuration rules
	validator.AddRule(ValidationRule{
		Field:    "database.url",
		Required: true,
		Type:     reflect.TypeOf(""),
		Custom: func(value interface{}) error {
			str, ok := value.(string)
			if !ok {
				return fmt.Errorf("must be a string")
			}
			if !strings.HasPrefix(str, "postgres://") && !strings.HasPrefix(str, "mysql://") {
				return fmt.Errorf("must be a valid database URL")
			}
			return nil
		},
	})

	validator.AddRule(ValidationRule{
		Field:    "database.max_connections",
		Required: true,
		Type:     reflect.TypeOf(0),
		Min:      floatPtr(1),
		Max:      floatPtr(100),
	})

	// Logging configuration rules
	validator.AddRule(ValidationRule{
		Field:    "logging.level",
		Required: true,
		Type:     reflect.TypeOf(""),
		Enum:     []interface{}{"debug", "info", "warn", "error", "fatal"},
	})

	validator.AddRule(ValidationRule{
		Field:    "logging.format",
		Required: true,
		Type:     reflect.TypeOf(""),
		Enum:     []interface{}{"json", "text"},
	})

	// Security configuration rules
	validator.AddRule(ValidationRule{
		Field: "security.jwt_secret",
		Type:  reflect.TypeOf(""),
		Custom: func(value interface{}) error {
			str, ok := value.(string)
			if !ok {
				return fmt.Errorf("must be a string")
			}
			if str != "" && len(str) < 32 {
				return fmt.Errorf("must be at least 32 characters long")
			}
			return nil
		},
	})

	return validator
}

// ValidateFileConfig validates a FileConfig using the default schema
func ValidateFileConfig(config *FileConfig) error {
	validator := CreateDefaultSchema()
	return validator.Validate(config)
}

// ValidateEnvConfig validates an EnvConfig
func ValidateEnvConfig(config *EnvConfig) error {
	var errors []string

	// Validate app name
	if config.AppName == "" {
		errors = append(errors, "APP_NAME is required")
	}

	// Validate app version format
	versionPattern := regexp.MustCompile(`^\d+\.\d+\.\d+$`)
	if !versionPattern.MatchString(config.AppVersion) {
		errors = append(errors, "APP_VERSION must be in format x.y.z")
	}

	// Validate environment
	validEnvs := []string{"development", "staging", "production"}
	if !contains(validEnvs, config.AppEnvironment) {
		errors = append(errors, fmt.Sprintf("APP_ENV must be one of: %s", strings.Join(validEnvs, ", ")))
	}

	// Validate server port
	if config.ServerPort < 1 || config.ServerPort > 65535 {
		errors = append(errors, "SERVER_PORT must be between 1 and 65535")
	}

	// Validate log level
	validLogLevels := []string{"debug", "info", "warn", "error", "fatal"}
	if !contains(validLogLevels, config.LogLevel) {
		errors = append(errors, fmt.Sprintf("LOG_LEVEL must be one of: %s", strings.Join(validLogLevels, ", ")))
	}

	// Validate log format
	validLogFormats := []string{"json", "text"}
	if !contains(validLogFormats, config.LogFormat) {
		errors = append(errors, fmt.Sprintf("LOG_FORMAT must be one of: %s", strings.Join(validLogFormats, ", ")))
	}

	// Validate database URL
	if config.DatabaseURL == "" {
		errors = append(errors, "DATABASE_URL is required")
	}

	// Validate database max connections
	if config.DatabaseMaxConns < 1 || config.DatabaseMaxConns > 100 {
		errors = append(errors, "DATABASE_MAX_CONNS must be between 1 and 100")
	}

	// Production-specific validations
	if config.AppEnvironment == "production" {
		if config.JWTSecret == "" {
			errors = append(errors, "JWT_SECRET is required in production")
		}
		if len(config.JWTSecret) < 32 {
			errors = append(errors, "JWT_SECRET must be at least 32 characters long")
		}
		if config.SessionSecret == "" {
			errors = append(errors, "SESSION_SECRET is required in production")
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("validation failed: %s", strings.Join(errors, "; "))
	}

	return nil
}

// Helper function to create float pointer
func floatPtr(f float64) *float64 {
	return &f
}
