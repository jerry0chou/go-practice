# Go Reflection Package

A comprehensive collection of Go reflection examples and utilities designed to make reflection concepts easy to understand and practical to use.

## 📚 Overview

This package demonstrates various aspects of Go's `reflect` package through practical examples and utility classes. It covers everything from basic type inspection to advanced use cases like dynamic function calls and plugin systems.

## 🚀 Quick Start

```bash
# Run all examples
go run run/reflect_main.go -mode=all

# Run specific examples
go run run/reflect_main.go -mode=basic
go run run/reflect_main.go -mode=struct
go run run/reflect_main.go -mode=function
go run run/reflect_main.go -mode=interface
go run run/reflect_main.go -mode=practical
go run run/reflect_main.go -mode=utilities
```

## 📖 Contents

### 1. Basic Reflection (`basic_reflection.go`)

**Core Concepts:**
- Type and Value inspection
- Kind vs Type differences
- Zero values and type creation
- Type conversions

**Key Functions:**
- `BasicReflection()` - Main demonstration function
- `DemonstrateTypeChecker()` - Utility class demonstration

**Example:**
```go
// Get type information
value := reflect.ValueOf(42)
fmt.Printf("Type: %v, Kind: %v\n", value.Type(), value.Kind())

// Create new values
newInt := reflect.New(reflect.TypeOf(int(0))).Elem()
newInt.SetInt(100)
```

### 2. Struct Reflection (`struct_reflection.go`)

**Core Concepts:**
- Field access and modification
- Struct tags parsing
- Method reflection
- Anonymous fields and embedding
- Dynamic struct creation

**Key Functions:**
- `StructReflection()` - Main demonstration function
- `DemonstrateStructAnalyzer()` - Utility class demonstration

**Example:**
```go
// Access struct fields
user := &User{Name: "John", Age: 30}
value := reflect.ValueOf(user).Elem()
nameField := value.FieldByName("Name")
fmt.Printf("Name: %v\n", nameField.String())

// Parse struct tags
field := reflect.TypeOf(user).Elem().Field(0)
jsonTag := field.Tag.Get("json")
```

### 3. Function Reflection (`function_reflection.go`)

**Core Concepts:**
- Function type inspection
- Dynamic function calls
- Parameter and return value analysis
- Higher-order functions
- Function factories

**Key Functions:**
- `FunctionReflection()` - Main demonstration function
- `DemonstrateFunctionRegistry()` - Utility class demonstration

**Example:**
```go
// Call function dynamically
addFunc := reflect.ValueOf(Add)
args := []reflect.Value{reflect.ValueOf(10), reflect.ValueOf(20)}
results := addFunc.Call(args)
fmt.Printf("Result: %v\n", results[0].Int())
```

### 4. Interface Reflection (`interface_reflection.go`)

**Core Concepts:**
- Interface type assertions
- Method reflection on interfaces
- Interface composition
- Dynamic interface implementation
- Interface value inspection

**Key Functions:**
- `InterfaceReflection()` - Main demonstration function
- `DemonstrateInterfaceAnalyzer()` - Utility class demonstration

**Example:**
```go
// Check interface implementation
var writer Writer = &FileWriter{}
writerType := reflect.TypeOf(writer)
readerType := reflect.TypeOf((*Reader)(nil)).Elem()
implements := writerType.Implements(readerType)
```

### 5. Practical Examples (`practical_examples.go`)

**Real-world Use Cases:**
- JSON marshaling/unmarshaling
- Struct validation
- Configuration loading
- Object cloning
- Generic utilities
- Plugin system

**Key Functions:**
- `PracticalExamples()` - Main demonstration function

**Example:**
```go
// Custom JSON marshaling
jsonData, err := marshalToJSON(product)

// Struct validation
errors := validateStruct(config)

// Deep cloning
cloned, err := deepClone(original)
```

## 🔧 Utility Classes

### TypeChecker
Provides utility functions for type checking and analysis.

```go
checker := &TypeChecker{}
isNumeric := checker.IsNumeric(42)
info := checker.GetTypeInfo(value)
```

### StructAnalyzer
Comprehensive struct analysis and inspection tools.

```go
analyzer := &StructAnalyzer{}
analysis := analyzer.AnalyzeStruct(user)
```

### FunctionRegistry
Dynamic function management and execution.

```go
registry := NewFunctionRegistry()
registry.Register("add", Add)
result, err := registry.Call("add", 10, 20)
```

### InterfaceAnalyzer
Interface analysis and implementation checking.

```go
analyzer := &InterfaceAnalyzer{}
analysis := analyzer.AnalyzeInterface(writer)
```

## 🎯 Key Concepts

### reflect.Type vs reflect.Value
- **Type**: Information about the type itself
- **Value**: Information about a specific value of that type

### Kind vs Type
- **Kind**: The underlying type (int, string, struct, etc.)
- **Type**: More specific information including package path, method set, etc.

### Type Assertions vs Reflection
- **Type Assertions**: Compile-time safe, faster
- **Reflection**: Runtime inspection, more flexible but slower

## ⚠️ Best Practices

1. **Performance**: Reflection is slower than direct type operations
2. **Caching**: Cache `reflect.Type` and `reflect.Value` when possible
3. **Safety**: Always check `CanSet()` before modifying values
4. **Error Handling**: Handle panics from reflection operations
5. **Alternatives**: Prefer type assertions over reflection when possible
6. **Interfaces**: Use interfaces to reduce reflection needs

## 🚨 Common Pitfalls

1. **Nil Pointers**: Always check if values are valid
2. **Type Mismatches**: Ensure types are compatible before conversion
3. **Unaddressable Values**: Use pointers for settable values
4. **Method Resolution**: Be careful with embedded types and method resolution
5. **Memory Leaks**: Be mindful of circular references in deep cloning

## 📝 Examples by Use Case

### JSON Processing
```go
// Custom JSON marshaling
jsonData, err := marshalToJSON(structValue)

// Custom JSON unmarshaling
err := unmarshalFromJSON(jsonString, &structValue)
```

### Validation
```go
// Struct validation with tags
errors := validateStruct(config)
```

### Configuration
```go
// Load configuration from map
err := loadConfigFromMap(configData, &config)

// Load with defaults
err := loadConfigWithDefaults(configData, &config)
```

### Cloning
```go
// Deep clone (independent copies)
cloned, err := deepClone(original)

// Shallow clone (shared references)
shallow := shallowClone(original)
```

### Generic Operations
```go
// Get/set field values
value, err := getFieldValue(struct, "FieldName")
err := setFieldValue(&struct, "FieldName", newValue)

// Copy specific fields
err := copyFields(source, &target, "Name", "Age")
```

## 🔍 Debugging Tips

1. **Print Types**: Use `fmt.Printf("%T", value)` to see types
2. **Inspect Values**: Use `reflect.ValueOf(value)` to inspect values
3. **Check Validity**: Always check `IsValid()` before using values
4. **Type Information**: Use `reflect.TypeOf(value)` for type details
5. **Method Lists**: Use `NumMethod()` and `Method(i)` to list methods

## 📚 Further Reading

- [Go Reflection Documentation](https://golang.org/pkg/reflect/)
- [The Go Blog: The Laws of Reflection](https://blog.golang.org/laws-of-reflection)
- [Effective Go: Reflection](https://golang.org/doc/effective_go.html#reflection)

## 🤝 Contributing

Feel free to add more examples or improve existing ones. The goal is to make reflection concepts accessible and practical for Go developers.

## 📄 License

This package is part of the go-practice project and follows the same license terms.

