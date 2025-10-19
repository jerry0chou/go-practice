package reflect

import (
	"fmt"
	"reflect"
	"strings"
)

// FunctionReflection demonstrates function reflect capabilities
func FunctionReflection() {
	fmt.Println("⚙️  Function Reflection Examples")
	fmt.Println(strings.Repeat("=", 50))

	// 1. Basic function inspection
	fmt.Println("\n🔍 1. Basic Function Inspection:")
	demonstrateFunctionInspection()

	// 2. Calling functions dynamically
	fmt.Println("\n📞 2. Calling Functions Dynamically:")
	demonstrateFunctionCalls()

	// 3. Function parameters and return values
	fmt.Println("\n📋 3. Function Parameters and Return Values:")
	demonstrateFunctionSignature()

	// 4. Higher-order functions
	fmt.Println("\n🔄 4. Higher-Order Functions:")
	demonstrateHigherOrderFunctions()

	// 5. Function factories
	fmt.Println("\n🏭 5. Function Factories:")
	demonstrateFunctionFactories()

	// 6. Method vs Function reflect
	fmt.Println("\n🔗 6. Method vs Function Reflection:")
	demonstrateMethodVsFunction()
}

// Sample functions for demonstration
func Add(a, b int) int {
	return a + b
}

func Multiply(a, b int) int {
	return a * b
}

func Greet(name string) string {
	return fmt.Sprintf("Hello, %s!", name)
}

func ProcessNumbers(numbers []int, operation func(int, int) int) int {
	if len(numbers) == 0 {
		return 0
	}
	result := numbers[0]
	for i := 1; i < len(numbers); i++ {
		result = operation(result, numbers[i])
	}
	return result
}

func CreateMultiplier(factor int) func(int) int {
	return func(x int) int {
		return x * factor
	}
}

func VariadicSum(numbers ...int) int {
	sum := 0
	for _, num := range numbers {
		sum += num
	}
	return sum
}

func ReturnMultiple() (int, string, bool) {
	return 42, "hello", true
}

func demonstrateFunctionInspection() {
	// Get function type information
	addFunc := reflect.TypeOf(Add)
	fmt.Printf("Add function type: %v\n", addFunc)
	fmt.Printf("Kind: %v\n", addFunc.Kind())
	fmt.Printf("NumIn: %d\n", addFunc.NumIn())
	fmt.Printf("NumOut: %d\n", addFunc.NumOut())

	// Inspect parameters
	fmt.Println("\nParameters:")
	for i := 0; i < addFunc.NumIn(); i++ {
		param := addFunc.In(i)
		fmt.Printf("  %d. %v\n", i+1, param)
	}

	// Inspect return values
	fmt.Println("\nReturn values:")
	for i := 0; i < addFunc.NumOut(); i++ {
		ret := addFunc.Out(i)
		fmt.Printf("  %d. %v\n", i+1, ret)
	}

	// Check if variadic
	fmt.Printf("\nIs variadic: %v\n", addFunc.IsVariadic())
}

func demonstrateFunctionCalls() {
	// Call Add function
	addFunc := reflect.ValueOf(Add)
	args := []reflect.Value{
		reflect.ValueOf(10),
		reflect.ValueOf(20),
	}
	results := addFunc.Call(args)
	fmt.Printf("Add(10, 20) = %v\n", results[0].Int())

	// Call Greet function
	greetFunc := reflect.ValueOf(Greet)
	args = []reflect.Value{
		reflect.ValueOf("World"),
	}
	results = greetFunc.Call(args)
	fmt.Printf("Greet(\"World\") = %v\n", results[0].String())

	// Call variadic function
	sumFunc := reflect.ValueOf(VariadicSum)
	args = []reflect.Value{
		reflect.ValueOf(1),
		reflect.ValueOf(2),
		reflect.ValueOf(3),
		reflect.ValueOf(4),
		reflect.ValueOf(5),
	}
	results = sumFunc.Call(args)
	fmt.Printf("VariadicSum(1,2,3,4,5) = %v\n", results[0].Int())

	// Call function with multiple return values
	multiFunc := reflect.ValueOf(ReturnMultiple)
	results = multiFunc.Call([]reflect.Value{})
	fmt.Printf("ReturnMultiple() = %v, %v, %v\n",
		results[0].Int(), results[1].String(), results[2].Bool())
}

func demonstrateFunctionSignature() {
	functions := []interface{}{
		Add,
		Greet,
		ProcessNumbers,
		CreateMultiplier,
		VariadicSum,
		ReturnMultiple,
	}

	for i, fn := range functions {
		fnType := reflect.TypeOf(fn)
		fmt.Printf("Function %d: %v\n", i+1, fnType)
		fmt.Printf("  Kind: %v\n", fnType.Kind())
		fmt.Printf("  Parameters: %d\n", fnType.NumIn())
		fmt.Printf("  Return values: %d\n", fnType.NumOut())
		fmt.Printf("  Is variadic: %v\n", fnType.IsVariadic())

		// Show parameter details
		if fnType.NumIn() > 0 {
			fmt.Println("  Parameters:")
			for j := 0; j < fnType.NumIn(); j++ {
				param := fnType.In(j)
				fmt.Printf("    %d. %v\n", j+1, param)
			}
		}

		// Show return value details
		if fnType.NumOut() > 0 {
			fmt.Println("  Return values:")
			for j := 0; j < fnType.NumOut(); j++ {
				ret := fnType.Out(j)
				fmt.Printf("    %d. %v\n", j+1, ret)
			}
		}
		fmt.Println()
	}
}

func demonstrateHigherOrderFunctions() {
	// Create a function that takes another function as parameter
	numbers := []int{1, 2, 3, 4, 5}

	// Get ProcessNumbers function
	processFunc := reflect.ValueOf(ProcessNumbers)

	// Create different operation functions
	addFunc := reflect.ValueOf(Add)
	multiplyFunc := reflect.ValueOf(Multiply)

	// Call ProcessNumbers with Add
	args := []reflect.Value{
		reflect.ValueOf(numbers),
		addFunc,
	}
	results := processFunc.Call(args)
	fmt.Printf("ProcessNumbers with Add: %v\n", results[0].Int())

	// Call ProcessNumbers with Multiply
	args = []reflect.Value{
		reflect.ValueOf(numbers),
		multiplyFunc,
	}
	results = processFunc.Call(args)
	fmt.Printf("ProcessNumbers with Multiply: %v\n", results[0].Int())

	// Demonstrate function composition
	fmt.Println("\nFunction composition:")
	composeFunc := reflect.ValueOf(createComposedFunction)
	args = []reflect.Value{
		reflect.ValueOf(10),
		reflect.ValueOf(5),
	}
	results = composeFunc.Call(args)
	fmt.Printf("Composed function result: %v\n", results[0].Int())
}

func createComposedFunction(x, y int) int {
	// This function composes Add and Multiply
	// First add, then multiply by 2
	addResult := Add(x, y)
	multiplyResult := Multiply(addResult, 2)
	return multiplyResult
}

func demonstrateFunctionFactories() {
	// Create multiplier functions using reflect
	factoryFunc := reflect.ValueOf(CreateMultiplier)

	// Create different multipliers
	multipliers := []int{2, 3, 5, 10}

	for _, factor := range multipliers {
		// Call CreateMultiplier with the factor
		args := []reflect.Value{reflect.ValueOf(factor)}
		multiplierFunc := factoryFunc.Call(args)[0]

		// Test the created multiplier
		testValue := 7
		testArgs := []reflect.Value{reflect.ValueOf(testValue)}
		result := multiplierFunc.Call(testArgs)[0]
		fmt.Printf("Multiply by %d: %d * %d = %v\n",
			factor, testValue, factor, result.Int())
	}

	// Create a function that returns functions
	fmt.Println("\nDynamic function creation:")
	dynamicFunc := createDynamicFunction("add")
	args := []reflect.Value{
		reflect.ValueOf(10),
		reflect.ValueOf(20),
	}
	results := dynamicFunc.Call(args)
	fmt.Printf("Dynamic add function: %v\n", results[0].Int())

	dynamicFunc = createDynamicFunction("multiply")
	results = dynamicFunc.Call(args)
	fmt.Printf("Dynamic multiply function: %v\n", results[0].Int())
}

func createDynamicFunction(operation string) reflect.Value {
	switch operation {
	case "add":
		return reflect.ValueOf(Add)
	case "multiply":
		return reflect.ValueOf(Multiply)
	default:
		return reflect.ValueOf(func(a, b int) int { return 0 })
	}
}

func demonstrateMethodVsFunction() {
	// Create a user instance
	user := &User{
		ID:     1,
		Name:   "Test User",
		Email:  "test@example.com",
		Age:    25,
		Active: true,
	}

	// Get method vs function
	method := reflect.ValueOf(user).MethodByName("GetFullInfo")
	function := reflect.ValueOf(Greet)

	fmt.Println("Method vs Function comparison:")
	fmt.Printf("Method type: %v\n", method.Type())
	fmt.Printf("Function type: %v\n", function.Type())

	// Call method
	if method.IsValid() {
		results := method.Call([]reflect.Value{})
		fmt.Printf("Method call result: %v\n", results[0].String())
	}

	// Call function
	results := function.Call([]reflect.Value{reflect.ValueOf("Method")})
	fmt.Printf("Function call result: %v\n", results[0].String())

	// Show the difference in receiver
	fmt.Println("\nMethod receiver analysis:")
	userType := reflect.TypeOf(user)
	for i := 0; i < userType.NumMethod(); i++ {
		method := userType.Method(i)
		fmt.Printf("  Method: %s\n", method.Name)
		fmt.Printf("    Type: %v\n", method.Type)
		fmt.Printf("    First parameter (receiver): %v\n", method.Type.In(0))
	}
}

// FunctionRegistry provides a registry for dynamic function management
type FunctionRegistry struct {
	functions map[string]reflect.Value
}

// NewFunctionRegistry creates a new function registry
func NewFunctionRegistry() *FunctionRegistry {
	return &FunctionRegistry{
		functions: make(map[string]reflect.Value),
	}
}

// Register adds a function to the registry
func (fr *FunctionRegistry) Register(name string, fn interface{}) {
	fr.functions[name] = reflect.ValueOf(fn)
}

// Call calls a registered function by name
func (fr *FunctionRegistry) Call(name string, args ...interface{}) ([]interface{}, error) {
	fn, exists := fr.functions[name]
	if !exists {
		return nil, fmt.Errorf("function %s not found", name)
	}

	// Convert args to reflect.Value slice
	reflectArgs := make([]reflect.Value, len(args))
	for i, arg := range args {
		reflectArgs[i] = reflect.ValueOf(arg)
	}

	// Call the function
	results := fn.Call(reflectArgs)

	// Convert results back to interface{} slice
	interfaceResults := make([]interface{}, len(results))
	for i, result := range results {
		interfaceResults[i] = result.Interface()
	}

	return interfaceResults, nil
}

// ListFunctions returns all registered function names
func (fr *FunctionRegistry) ListFunctions() []string {
	names := make([]string, 0, len(fr.functions))
	for name := range fr.functions {
		names = append(names, name)
	}
	return names
}

// GetFunctionInfo returns information about a registered function
func (fr *FunctionRegistry) GetFunctionInfo(name string) (map[string]interface{}, error) {
	fn, exists := fr.functions[name]
	if !exists {
		return nil, fmt.Errorf("function %s not found", name)
	}

	fnType := fn.Type()
	info := map[string]interface{}{
		"name":       name,
		"type":       fnType.String(),
		"numIn":      fnType.NumIn(),
		"numOut":     fnType.NumOut(),
		"isVariadic": fnType.IsVariadic(),
	}

	// Parameter types
	paramTypes := make([]string, fnType.NumIn())
	for i := 0; i < fnType.NumIn(); i++ {
		paramTypes[i] = fnType.In(i).String()
	}
	info["paramTypes"] = paramTypes

	// Return types
	returnTypes := make([]string, fnType.NumOut())
	for i := 0; i < fnType.NumOut(); i++ {
		returnTypes[i] = fnType.Out(i).String()
	}
	info["returnTypes"] = returnTypes

	return info, nil
}

// DemonstrateFunctionRegistry shows how to use the FunctionRegistry
func DemonstrateFunctionRegistry() {
	fmt.Println("\n🔧 FunctionRegistry Utility:")
	fmt.Println(strings.Repeat("-", 30))

	registry := NewFunctionRegistry()

	// Register functions
	registry.Register("add", Add)
	registry.Register("multiply", Multiply)
	registry.Register("greet", Greet)
	registry.Register("variadicSum", VariadicSum)

	// List registered functions
	fmt.Println("Registered functions:")
	for _, name := range registry.ListFunctions() {
		info, _ := registry.GetFunctionInfo(name)
		fmt.Printf("  %s: %s\n", name, info["type"])
	}

	// Call functions dynamically
	fmt.Println("\nCalling functions:")

	// Call add
	results, err := registry.Call("add", 10, 20)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("add(10, 20) = %v\n", results[0])
	}

	// Call greet
	results, err = registry.Call("greet", "Reflection")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("greet(\"Reflection\") = %v\n", results[0])
	}

	// Call variadic function
	results, err = registry.Call("variadicSum", 1, 2, 3, 4, 5)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("variadicSum(1,2,3,4,5) = %v\n", results[0])
	}
}
