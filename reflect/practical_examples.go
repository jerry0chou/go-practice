package reflect

import (
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

// PracticalExamples demonstrates real-world uses of reflect
func PracticalExamples() {
	fmt.Println("🛠️  Practical Reflection Examples")
	fmt.Println(strings.Repeat("=", 50))

	// 1. JSON Marshaling/Unmarshaling
	fmt.Println("\n📄 1. JSON Marshaling/Unmarshaling:")
	demonstrateJSONOperations()

	// 2. Struct Validation
	fmt.Println("\n✅ 2. Struct Validation:")
	demonstrateStructValidation()

	// 3. Configuration Loading
	fmt.Println("\n⚙️  3. Configuration Loading:")
	demonstrateConfigurationLoading()

	// 4. Object Cloning
	fmt.Println("\n📋 4. Object Cloning:")
	demonstrateObjectCloning()

	// 5. Generic Utilities
	fmt.Println("\n🔧 5. Generic Utilities:")
	demonstrateGenericUtilities()

	// 6. Plugin System
	fmt.Println("\n🔌 6. Plugin System:")
	demonstratePluginSystem()
}

// Configuration struct for demonstration
type Config struct {
	Database DatabaseConfig `json:"database" validate:"required"`
	Server   ServerConfig   `json:"server" validate:"required"`
	Features FeatureConfig  `json:"features"`
}

type DatabaseConfig struct {
	Host     string `json:"host" validate:"required"`
	Port     int    `json:"port" validate:"required,min=1,max=65535"`
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	SSL      bool   `json:"ssl"`
}

type ServerConfig struct {
	Host string `json:"host" validate:"required"`
	Port int    `json:"port" validate:"required,min=1,max=65535"`
}

type FeatureConfig struct {
	EnableCache    bool `json:"enable_cache"`
	EnableMetrics  bool `json:"enable_metrics"`
	EnableLogging  bool `json:"enable_logging"`
	MaxConnections int  `json:"max_connections" validate:"min=1,max=1000"`
}

// Product struct for cloning demonstration
type Product struct {
	ID          int               `json:"id" validate:"required,min=1"`
	Name        string            `json:"name" validate:"required,min=2"`
	Price       float64           `json:"price" validate:"required,min=0"`
	Description string            `json:"description"`
	Tags        []string          `json:"tags"`
	Metadata    map[string]string `json:"metadata"`
	Active      bool              `json:"active"`
}

func demonstrateJSONOperations() {
	// Create a sample product
	product := Product{
		ID:          1,
		Name:        "Go Programming Book",
		Price:       29.99,
		Description: "Learn Go programming with practical examples",
		Tags:        []string{"programming", "golang", "book"},
		Metadata:    map[string]string{"author": "John Doe", "pages": "300"},
		Active:      true,
	}

	// Custom JSON marshaling using reflect
	fmt.Println("Custom JSON marshaling:")
	jsonData, err := marshalToJSON(product)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("JSON: %s\n", jsonData)
	}

	// Custom JSON unmarshaling using reflect
	fmt.Println("\nCustom JSON unmarshaling:")
	jsonStr := `{"id":2,"name":"Python Book","price":34.99,"description":"Learn Python","tags":["programming","python"],"metadata":{"author":"Jane Smith"},"active":true}`

	var newProduct Product
	err = unmarshalFromJSON(jsonStr, &newProduct)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Unmarshaled product: %+v\n", newProduct)
	}

	// Field-specific JSON operations
	fmt.Println("\nField-specific operations:")
	fieldJSON, err := marshalFieldToJSON(product, "Name")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Name field JSON: %s\n", fieldJSON)
	}
}

func demonstrateStructValidation() {
	// Create various configs for validation
	validConfig := Config{
		Database: DatabaseConfig{
			Host:     "localhost",
			Port:     5432,
			Username: "user",
			Password: "password",
			SSL:      true,
		},
		Server: ServerConfig{
			Host: "0.0.0.0",
			Port: 8080,
		},
		Features: FeatureConfig{
			EnableCache:    true,
			EnableMetrics:  false,
			EnableLogging:  true,
			MaxConnections: 100,
		},
	}

	invalidConfig := Config{
		Database: DatabaseConfig{
			Host:     "", // Invalid: empty host
			Port:     0,  // Invalid: port 0
			Username: "user",
			Password: "password",
		},
		Server: ServerConfig{
			Host: "localhost",
			Port: 99999, // Invalid: port too high
		},
		Features: FeatureConfig{
			MaxConnections: 2000, // Invalid: too many connections
		},
	}

	configs := []Config{validConfig, invalidConfig}
	configNames := []string{"Valid Config", "Invalid Config"}

	for i, config := range configs {
		fmt.Printf("Validating %s:\n", configNames[i])
		errors := validateStruct(config)
		if len(errors) == 0 {
			fmt.Println("  ✅ All validations passed")
		} else {
			fmt.Println("  ❌ Validation errors:")
			for _, err := range errors {
				fmt.Printf("    - %s\n", err)
			}
		}
		fmt.Println()
	}
}

func demonstrateConfigurationLoading() {
	// Simulate loading configuration from different sources
	configData := map[string]interface{}{
		"database": map[string]interface{}{
			"host":     "db.example.com",
			"port":     3306,
			"username": "admin",
			"password": "secret123",
			"ssl":      true,
		},
		"server": map[string]interface{}{
			"host": "api.example.com",
			"port": 443,
		},
		"features": map[string]interface{}{
			"enable_cache":    true,
			"enable_metrics":  true,
			"enable_logging":  false,
			"max_connections": 500,
		},
	}

	fmt.Println("Loading configuration from map:")
	var config Config
	err := loadConfigFromMap(configData, &config)
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
	} else {
		fmt.Printf("Loaded config: %+v\n", config)
	}

	// Load configuration with defaults
	fmt.Println("\nLoading configuration with defaults:")
	var configWithDefaults Config
	err = loadConfigWithDefaults(configData, &configWithDefaults)
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
	} else {
		fmt.Printf("Config with defaults: %+v\n", configWithDefaults)
	}
}

func demonstrateObjectCloning() {
	// Create original product
	original := Product{
		ID:          1,
		Name:        "Original Product",
		Price:       99.99,
		Description: "This is the original product",
		Tags:        []string{"original", "test"},
		Metadata:    map[string]string{"version": "1.0", "category": "electronics"},
		Active:      true,
	}

	fmt.Printf("Original product: %+v\n", original)

	// Deep clone the product
	cloned, err := deepClone(original)
	if err != nil {
		fmt.Printf("Error cloning: %v\n", err)
		return
	}

	fmt.Printf("Cloned product: %+v\n", cloned)

	// Modify the cloned product
	clonedProduct := cloned.(Product)
	clonedProduct.Name = "Modified Product"
	clonedProduct.Price = 149.99
	clonedProduct.Tags = append(clonedProduct.Tags, "modified")
	clonedProduct.Metadata["version"] = "2.0"

	fmt.Printf("Modified cloned product: %+v\n", clonedProduct)
	fmt.Printf("Original product (should be unchanged): %+v\n", original)

	// Shallow clone demonstration
	fmt.Println("\nShallow clone demonstration:")
	shallowCloned := shallowClone(original)
	shallowClonedProduct := shallowCloned.(Product)
	shallowClonedProduct.Name = "Shallow Modified"
	shallowClonedProduct.Tags[0] = "shallow"
	shallowClonedProduct.Metadata["version"] = "shallow"

	fmt.Printf("Shallow cloned product: %+v\n", shallowClonedProduct)
	fmt.Printf("Original product (may be affected): %+v\n", original)
}

func demonstrateGenericUtilities() {
	// Generic field getter/setter
	product := Product{
		ID:     1,
		Name:   "Test Product",
		Price:  19.99,
		Active: true,
	}

	fmt.Println("Generic field operations:")

	// Get field values
	if name, err := getFieldValue(product, "Name"); err == nil {
		fmt.Printf("Name field: %v\n", name)
	}

	if price, err := getFieldValue(product, "Price"); err == nil {
		fmt.Printf("Price field: %v\n", price)
	}

	// Set field values
	fmt.Println("\nSetting field values:")
	if err := setFieldValue(&product, "Name", "Updated Product"); err == nil {
		fmt.Printf("Updated name: %v\n", product.Name)
	}

	if err := setFieldValue(&product, "Price", 29.99); err == nil {
		fmt.Printf("Updated price: %v\n", product.Price)
	}

	// Generic comparison
	fmt.Println("\nGeneric comparison:")
	product1 := Product{ID: 1, Name: "Product 1", Price: 10.0}
	product2 := Product{ID: 1, Name: "Product 1", Price: 10.0}
	product3 := Product{ID: 2, Name: "Product 2", Price: 20.0}

	fmt.Printf("product1 == product2: %t\n", deepEqual(product1, product2))
	fmt.Printf("product1 == product3: %t\n", deepEqual(product1, product3))

	// Generic field copying
	fmt.Println("\nGeneric field copying:")
	source := Product{ID: 1, Name: "Source", Price: 100.0, Active: true}
	target := Product{ID: 2, Name: "Target", Price: 200.0, Active: false}

	fmt.Printf("Before copy - Source: %+v\n", source)
	fmt.Printf("Before copy - Target: %+v\n", target)

	copyFields(source, &target, "Name", "Price")

	fmt.Printf("After copy - Source: %+v\n", source)
	fmt.Printf("After copy - Target: %+v\n", target)
}

func demonstratePluginSystem() {
	// Create a plugin registry
	registry := NewPluginRegistry()

	// Register plugins
	registry.RegisterPlugin("logger", &LoggerPlugin{})
	registry.RegisterPlugin("cache", &CachePlugin{})
	registry.RegisterPlugin("metrics", &MetricsPlugin{})

	// List registered plugins
	fmt.Println("Registered plugins:")
	plugins := registry.ListPlugins()
	for _, name := range plugins {
		fmt.Printf("  - %s\n", name)
	}

	// Execute plugins
	fmt.Println("\nExecuting plugins:")
	context := map[string]interface{}{
		"message": "Hello from plugin system",
		"level":   "info",
	}

	for _, name := range plugins {
		fmt.Printf("Executing %s plugin:\n", name)
		result, err := registry.ExecutePlugin(name, context)
		if err != nil {
			fmt.Printf("  Error: %v\n", err)
		} else {
			fmt.Printf("  Result: %v\n", result)
		}
	}
}

// JSON Operations
func marshalToJSON(v interface{}) (string, error) {
	value := reflect.ValueOf(v)
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}

	if value.Kind() != reflect.Struct {
		return "", fmt.Errorf("value must be a struct")
	}

	result := make(map[string]interface{})
	structType := value.Type()

	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		fieldValue := value.Field(i)

		// Get JSON tag
		jsonTag := field.Tag.Get("json")
		if jsonTag == "" || jsonTag == "-" {
			continue
		}

		// Remove options from tag (e.g., "name,omitempty" -> "name")
		jsonName := strings.Split(jsonTag, ",")[0]

		result[jsonName] = fieldValue.Interface()
	}

	jsonBytes, err := json.Marshal(result)
	return string(jsonBytes), err
}

func unmarshalFromJSON(jsonStr string, v interface{}) error {
	value := reflect.ValueOf(v)
	if value.Kind() != reflect.Ptr || value.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("v must be a pointer to a struct")
	}

	// Parse JSON
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &data); err != nil {
		return err
	}

	// Set struct fields
	structValue := value.Elem()
	structType := structValue.Type()

	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		fieldValue := structValue.Field(i)

		// Get JSON tag
		jsonTag := field.Tag.Get("json")
		if jsonTag == "" || jsonTag == "-" {
			continue
		}

		jsonName := strings.Split(jsonTag, ",")[0]

		// Find corresponding value in JSON data
		if jsonValue, exists := data[jsonName]; exists {
			if err := setValueFromInterface(fieldValue, jsonValue); err != nil {
				return fmt.Errorf("error setting field %s: %v", field.Name, err)
			}
		}
	}

	return nil
}

func marshalFieldToJSON(v interface{}, fieldName string) (string, error) {
	value := reflect.ValueOf(v)
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}

	fieldValue := value.FieldByName(fieldName)
	if !fieldValue.IsValid() {
		return "", fmt.Errorf("field %s not found", fieldName)
	}

	jsonBytes, err := json.Marshal(fieldValue.Interface())
	return string(jsonBytes), err
}

func setValueFromInterface(fieldValue reflect.Value, jsonValue interface{}) error {
	if !fieldValue.CanSet() {
		return fmt.Errorf("field cannot be set")
	}

	jsonValueType := reflect.TypeOf(jsonValue)
	fieldType := fieldValue.Type()

	// Direct assignment if types match
	if jsonValueType.AssignableTo(fieldType) {
		fieldValue.Set(reflect.ValueOf(jsonValue))
		return nil
	}

	// Type conversion
	if jsonValueType.ConvertibleTo(fieldType) {
		fieldValue.Set(reflect.ValueOf(jsonValue).Convert(fieldType))
		return nil
	}

	return fmt.Errorf("cannot convert %v to %v", jsonValueType, fieldType)
}

// Validation
func validateStruct(v interface{}) []string {
	var errors []string
	value := reflect.ValueOf(v)
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}

	if value.Kind() != reflect.Struct {
		return []string{"value must be a struct"}
	}

	structType := value.Type()
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		fieldValue := value.Field(i)

		// Get validation tag
		validateTag := field.Tag.Get("validate")
		if validateTag == "" {
			continue
		}

		// Parse validation rules
		rules := strings.Split(validateTag, ",")
		fieldErrors := validateField(field.Name, fieldValue, rules)
		errors = append(errors, fieldErrors...)
	}

	return errors
}

func validateField(fieldName string, fieldValue reflect.Value, rules []string) []string {
	var errors []string

	for _, rule := range rules {
		rule = strings.TrimSpace(rule)

		switch rule {
		case "required":
			if fieldValue.IsZero() {
				errors = append(errors, fmt.Sprintf("%s is required", fieldName))
			}
		case "email":
			if fieldValue.Kind() == reflect.String {
				email := fieldValue.String()
				if !isValidEmail(email) {
					errors = append(errors, fmt.Sprintf("%s must be a valid email", fieldName))
				}
			}
		default:
			if strings.HasPrefix(rule, "min=") {
				minStr := strings.TrimPrefix(rule, "min=")
				if min, err := strconv.Atoi(minStr); err == nil {
					if !validateMin(fieldValue, min) {
						errors = append(errors, fmt.Sprintf("%s must be at least %d", fieldName, min))
					}
				}
			} else if strings.HasPrefix(rule, "max=") {
				maxStr := strings.TrimPrefix(rule, "max=")
				if max, err := strconv.Atoi(maxStr); err == nil {
					if !validateMax(fieldValue, max) {
						errors = append(errors, fmt.Sprintf("%s must be at most %d", fieldName, max))
					}
				}
			}
		}
	}

	return errors
}

func isValidEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(pattern, email)
	return matched
}

func validateMin(value reflect.Value, min int) bool {
	switch value.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return value.Int() >= int64(min)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return value.Uint() >= uint64(min)
	case reflect.Float32, reflect.Float64:
		return value.Float() >= float64(min)
	case reflect.String:
		return len(value.String()) >= min
	case reflect.Slice, reflect.Array:
		return value.Len() >= min
	default:
		return true
	}
}

func validateMax(value reflect.Value, max int) bool {
	switch value.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return value.Int() <= int64(max)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return value.Uint() <= uint64(max)
	case reflect.Float32, reflect.Float64:
		return value.Float() <= float64(max)
	case reflect.String:
		return len(value.String()) <= max
	case reflect.Slice, reflect.Array:
		return value.Len() <= max
	default:
		return true
	}
}

// Configuration Loading
func loadConfigFromMap(data map[string]interface{}, config interface{}) error {
	value := reflect.ValueOf(config)
	if value.Kind() != reflect.Ptr || value.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("config must be a pointer to a struct")
	}

	return setStructFromMap(value.Elem(), data)
}

func setStructFromMap(structValue reflect.Value, data map[string]interface{}) error {
	structType := structValue.Type()

	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		fieldValue := structValue.Field(i)

		// Get JSON tag for field name
		jsonTag := field.Tag.Get("json")
		if jsonTag == "" {
			jsonTag = strings.ToLower(field.Name)
		} else {
			jsonTag = strings.Split(jsonTag, ",")[0]
		}

		// Find value in data
		if value, exists := data[jsonTag]; exists {
			if fieldValue.Kind() == reflect.Struct {
				if subMap, ok := value.(map[string]interface{}); ok {
					if err := setStructFromMap(fieldValue, subMap); err != nil {
						return err
					}
				}
			} else {
				if err := setValueFromInterface(fieldValue, value); err != nil {
					return fmt.Errorf("error setting field %s: %v", field.Name, err)
				}
			}
		}
	}

	return nil
}

func loadConfigWithDefaults(data map[string]interface{}, config interface{}) error {
	// First set defaults
	setDefaults(config)

	// Then load from data
	return loadConfigFromMap(data, config)
}

func setDefaults(config interface{}) {
	value := reflect.ValueOf(config)
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}

	if value.Kind() != reflect.Struct {
		return
	}

	structType := value.Type()
	for i := 0; i < structType.NumField(); i++ {
		fieldValue := value.Field(i)

		if fieldValue.CanSet() && fieldValue.IsZero() {
			switch fieldValue.Kind() {
			case reflect.String:
				fieldValue.SetString("")
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				fieldValue.SetInt(0)
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				fieldValue.SetUint(0)
			case reflect.Float32, reflect.Float64:
				fieldValue.SetFloat(0)
			case reflect.Bool:
				fieldValue.SetBool(false)
			case reflect.Slice:
				fieldValue.Set(reflect.MakeSlice(fieldValue.Type(), 0, 0))
			case reflect.Map:
				fieldValue.Set(reflect.MakeMap(fieldValue.Type()))
			case reflect.Struct:
				setDefaults(fieldValue.Addr().Interface())
			}
		}
	}
}

// Cloning
func deepClone(original interface{}) (interface{}, error) {
	originalValue := reflect.ValueOf(original)
	originalType := originalValue.Type()

	// Create new instance
	newValue := reflect.New(originalType).Elem()

	// Copy all fields
	for i := 0; i < originalType.NumField(); i++ {
		originalField := originalValue.Field(i)
		newField := newValue.Field(i)

		if newField.CanSet() {
			switch originalField.Kind() {
			case reflect.Slice:
				if originalField.IsNil() {
					newField.Set(reflect.Zero(originalField.Type()))
				} else {
					newSlice := reflect.MakeSlice(originalField.Type(), originalField.Len(), originalField.Cap())
					reflect.Copy(newSlice, originalField)
					newField.Set(newSlice)
				}
			case reflect.Map:
				if originalField.IsNil() {
					newField.Set(reflect.Zero(originalField.Type()))
				} else {
					newMap := reflect.MakeMap(originalField.Type())
					for _, key := range originalField.MapKeys() {
						value := originalField.MapIndex(key)
						newMap.SetMapIndex(key, value)
					}
					newField.Set(newMap)
				}
			case reflect.Struct:
				clonedStruct, err := deepClone(originalField.Interface())
				if err != nil {
					return nil, err
				}
				newField.Set(reflect.ValueOf(clonedStruct))
			case reflect.Ptr:
				if originalField.IsNil() {
					newField.Set(reflect.Zero(originalField.Type()))
				} else {
					clonedPtr, err := deepClone(originalField.Elem().Interface())
					if err != nil {
						return nil, err
					}
					newPtr := reflect.New(originalField.Type().Elem())
					newPtr.Elem().Set(reflect.ValueOf(clonedPtr))
					newField.Set(newPtr)
				}
			default:
				newField.Set(originalField)
			}
		}
	}

	return newValue.Interface(), nil
}

func shallowClone(original interface{}) interface{} {
	originalValue := reflect.ValueOf(original)
	originalType := originalValue.Type()

	// Create new instance
	newValue := reflect.New(originalType).Elem()

	// Copy all fields (shallow copy)
	for i := 0; i < originalType.NumField(); i++ {
		originalField := originalValue.Field(i)
		newField := newValue.Field(i)

		if newField.CanSet() {
			newField.Set(originalField)
		}
	}

	return newValue.Interface()
}

// Generic Utilities
func getFieldValue(v interface{}, fieldName string) (interface{}, error) {
	value := reflect.ValueOf(v)
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}

	if value.Kind() != reflect.Struct {
		return nil, fmt.Errorf("value must be a struct")
	}

	fieldValue := value.FieldByName(fieldName)
	if !fieldValue.IsValid() {
		return nil, fmt.Errorf("field %s not found", fieldName)
	}

	return fieldValue.Interface(), nil
}

func setFieldValue(v interface{}, fieldName string, newValue interface{}) error {
	value := reflect.ValueOf(v)
	if value.Kind() != reflect.Ptr || value.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("v must be a pointer to a struct")
	}

	value = value.Elem()
	fieldValue := value.FieldByName(fieldName)
	if !fieldValue.IsValid() {
		return fmt.Errorf("field %s not found", fieldName)
	}

	if !fieldValue.CanSet() {
		return fmt.Errorf("field %s cannot be set", fieldName)
	}

	return setValueFromInterface(fieldValue, newValue)
}

func deepEqual(a, b interface{}) bool {
	return reflect.DeepEqual(a, b)
}

func copyFields(source, target interface{}, fieldNames ...string) error {
	sourceValue := reflect.ValueOf(source)
	targetValue := reflect.ValueOf(target)

	if sourceValue.Kind() == reflect.Ptr {
		sourceValue = sourceValue.Elem()
	}
	if targetValue.Kind() == reflect.Ptr {
		targetValue = targetValue.Elem()
	}

	if sourceValue.Kind() != reflect.Struct || targetValue.Kind() != reflect.Struct {
		return fmt.Errorf("both source and target must be structs")
	}

	for _, fieldName := range fieldNames {
		sourceField := sourceValue.FieldByName(fieldName)
		targetField := targetValue.FieldByName(fieldName)

		if !sourceField.IsValid() || !targetField.IsValid() {
			continue
		}

		if !targetField.CanSet() {
			continue
		}

		if sourceField.Type() == targetField.Type() {
			targetField.Set(sourceField)
		}
	}

	return nil
}

// Plugin System
type Plugin interface {
	Execute(context map[string]interface{}) (interface{}, error)
}

type LoggerPlugin struct{}

func (p *LoggerPlugin) Execute(context map[string]interface{}) (interface{}, error) {
	message, _ := context["message"].(string)
	level, _ := context["level"].(string)
	return fmt.Sprintf("LOG[%s]: %s", level, message), nil
}

type CachePlugin struct{}

func (p *CachePlugin) Execute(context map[string]interface{}) (interface{}, error) {
	message, _ := context["message"].(string)
	return fmt.Sprintf("CACHE: Stored '%s'", message), nil
}

type MetricsPlugin struct{}

func (p *MetricsPlugin) Execute(context map[string]interface{}) (interface{}, error) {
	message, _ := context["message"].(string)
	return fmt.Sprintf("METRICS: Recorded event for '%s'", message), nil
}

type PluginRegistry struct {
	plugins map[string]Plugin
}

func NewPluginRegistry() *PluginRegistry {
	return &PluginRegistry{
		plugins: make(map[string]Plugin),
	}
}

func (pr *PluginRegistry) RegisterPlugin(name string, plugin Plugin) {
	pr.plugins[name] = plugin
}

func (pr *PluginRegistry) ExecutePlugin(name string, context map[string]interface{}) (interface{}, error) {
	plugin, exists := pr.plugins[name]
	if !exists {
		return nil, fmt.Errorf("plugin %s not found", name)
	}

	// Use reflect to call Execute method
	pluginValue := reflect.ValueOf(plugin)
	executeMethod := pluginValue.MethodByName("Execute")

	if !executeMethod.IsValid() {
		return nil, fmt.Errorf("plugin %s does not implement Execute method", name)
	}

	// Call Execute method
	args := []reflect.Value{reflect.ValueOf(context)}
	results := executeMethod.Call(args)

	// Handle results
	if len(results) != 2 {
		return nil, fmt.Errorf("plugin %s Execute method must return (interface{}, error)", name)
	}

	result := results[0].Interface()
	var err error
	if !results[1].IsNil() {
		err = results[1].Interface().(error)
	}

	return result, err
}

func (pr *PluginRegistry) ListPlugins() []string {
	names := make([]string, 0, len(pr.plugins))
	for name := range pr.plugins {
		names = append(names, name)
	}
	return names
}
