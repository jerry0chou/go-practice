package reflect

import (
	"fmt"
	"reflect"
	"strings"
)

// BasicReflection demonstrates fundamental reflect concepts
func BasicReflection() {
	fmt.Println("🔍 Basic Reflection Examples")
	fmt.Println(strings.Repeat("=", 50))

	// 1. Getting Type and Value
	fmt.Println("\n📋 1. Type and Value Information:")
	demonstrateTypeAndValue()

	// 2. Kind vs Type
	fmt.Println("\n🏷️  2. Kind vs Type:")
	demonstrateKindVsType()

	// 3. Zero Values
	fmt.Println("\n⚪ 3. Zero Values:")
	demonstrateZeroValues()

	// 4. Creating Values
	fmt.Println("\n🛠️  4. Creating Values:")
	demonstrateCreatingValues()

	// 5. Converting Between Types
	fmt.Println("\n🔄 5. Type Conversions:")
	demonstrateTypeConversions()
}

func demonstrateTypeAndValue() {
	// Different ways to get type and value information
	var num int = 42
	var str string = "Hello, Reflection!"
	var slice []int = []int{1, 2, 3}

	// Using reflect.TypeOf() - gets the type
	fmt.Printf("Type of num: %v\n", reflect.TypeOf(num))
	fmt.Printf("Type of str: %v\n", reflect.TypeOf(str))
	fmt.Printf("Type of slice: %v\n", reflect.TypeOf(slice))

	// Using reflect.ValueOf() - gets the value
	fmt.Printf("Value of num: %v\n", reflect.ValueOf(num))
	fmt.Printf("Value of str: %v\n", reflect.ValueOf(str))
	fmt.Printf("Value of slice: %v\n", reflect.ValueOf(slice))

	// Getting both type and value
	value := reflect.ValueOf(num)
	fmt.Printf("Value: %v, Type: %v, Kind: %v\n",
		value, value.Type(), value.Kind())
}

func demonstrateKindVsType() {
	// Kind is the underlying type, Type includes more specific information
	var num int = 42
	var numPtr *int = &num
	var numSlice []int = []int{1, 2, 3}

	values := []interface{}{num, numPtr, numSlice}

	for i, v := range values {
		value := reflect.ValueOf(v)
		fmt.Printf("Value %d: %v\n", i+1, v)
		fmt.Printf("  Type: %v\n", value.Type())
		fmt.Printf("  Kind: %v\n", value.Kind())
		fmt.Printf("  Can Set: %v\n", value.CanSet())
		fmt.Println()
	}
}

func demonstrateZeroValues() {
	// Getting zero values for different types
	types := []reflect.Type{
		reflect.TypeOf(int(0)),
		reflect.TypeOf(string("")),
		reflect.TypeOf(bool(false)),
		reflect.TypeOf([]int(nil)),
		reflect.TypeOf(map[string]int(nil)),
	}

	fmt.Println("Zero values for different types:")
	for _, t := range types {
		zeroValue := reflect.Zero(t)
		fmt.Printf("  %v: %v\n", t, zeroValue)
	}
}

func demonstrateCreatingValues() {
	// Creating new values using reflect
	fmt.Println("Creating new values:")

	// Create a new int
	intType := reflect.TypeOf(int(0))
	newInt := reflect.New(intType).Elem()
	newInt.SetInt(100)
	fmt.Printf("  New int: %v (type: %v)\n", newInt.Int(), newInt.Type())

	// Create a new string
	stringType := reflect.TypeOf(string(""))
	newString := reflect.New(stringType).Elem()
	newString.SetString("Created with reflect!")
	fmt.Printf("  New string: %v (type: %v)\n", newString.String(), newString.Type())

	// Create a new slice
	sliceType := reflect.TypeOf([]int{})
	newSlice := reflect.MakeSlice(sliceType, 3, 3)
	for i := 0; i < 3; i++ {
		newSlice.Index(i).SetInt(int64(i * 10))
	}
	fmt.Printf("  New slice: %v (type: %v)\n", newSlice.Interface(), newSlice.Type())
}

func demonstrateTypeConversions() {
	// Converting between types using reflect
	fmt.Println("Type conversions:")

	// Convert int to float64
	intValue := reflect.ValueOf(42)
	if intValue.CanConvert(reflect.TypeOf(float64(0))) {
		floatValue := intValue.Convert(reflect.TypeOf(float64(0)))
		fmt.Printf("  int %v -> float64 %v\n", intValue.Int(), floatValue.Float())
	}

	// Convert string to []byte
	strValue := reflect.ValueOf("Hello")
	if strValue.CanConvert(reflect.TypeOf([]byte{})) {
		byteValue := strValue.Convert(reflect.TypeOf([]byte{}))
		fmt.Printf("  string %v -> []byte %v\n", strValue.String(), byteValue.Bytes())
	}

	// Interface conversion
	var i interface{} = 42
	value := reflect.ValueOf(i)
	if value.CanConvert(reflect.TypeOf(string(""))) {
		strValue := value.Convert(reflect.TypeOf(string("")))
		fmt.Printf("  interface{} %v -> string %v\n", value.Interface(), strValue.String())
	}
}

// TypeChecker provides utility functions for type checking
type TypeChecker struct{}

// IsNumeric checks if a value is a numeric type
func (tc *TypeChecker) IsNumeric(value interface{}) bool {
	kind := reflect.TypeOf(value).Kind()
	return kind >= reflect.Int && kind <= reflect.Complex128
}

// IsCollection checks if a value is a collection type (slice, array, map)
func (tc *TypeChecker) IsCollection(value interface{}) bool {
	kind := reflect.TypeOf(value).Kind()
	return kind == reflect.Slice || kind == reflect.Array || kind == reflect.Map
}

// GetTypeInfo returns detailed type information
func (tc *TypeChecker) GetTypeInfo(value interface{}) map[string]interface{} {
	v := reflect.ValueOf(value)
	t := v.Type()

	return map[string]interface{}{
		"value":        v.Interface(),
		"type":         t.String(),
		"kind":         t.Kind().String(),
		"canSet":       v.CanSet(),
		"isZero":       v.IsZero(),
		"isValid":      v.IsValid(),
		"isNumeric":    tc.IsNumeric(value),
		"isCollection": tc.IsCollection(value),
	}
}

// DemonstrateTypeChecker shows how to use the TypeChecker utility
func DemonstrateTypeChecker() {
	fmt.Println("\n🔧 TypeChecker Utility:")
	fmt.Println(strings.Repeat("-", 30))

	checker := &TypeChecker{}

	values := []interface{}{
		42,
		"hello",
		[]int{1, 2, 3},
		map[string]int{"a": 1},
		true,
		3.14,
	}

	for _, value := range values {
		info := checker.GetTypeInfo(value)
		fmt.Printf("Value: %v\n", value)
		for key, val := range info {
			fmt.Printf("  %s: %v\n", key, val)
		}
		fmt.Println()
	}
}
