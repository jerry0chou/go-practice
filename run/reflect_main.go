package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/jerrychou/go-practice/reflect"
)

func main() {
	mode := flag.String("mode", "all", "Mode to run: all, basic, struct, function, interface, practical, utilities")
	flag.Parse()

	fmt.Println("🔍 Go Reflection Package Demo")
	fmt.Println(strings.Repeat("=", 50))

	switch *mode {
	case "all":
		runAllExamples()
	case "basic":
		runBasicExamples()
	case "struct":
		runStructExamples()
	case "function":
		runFunctionExamples()
	case "interface":
		runInterfaceExamples()
	case "practical":
		runPracticalExamples()
	case "utilities":
		runUtilityExamples()
	default:
		fmt.Printf("❌ Unknown mode: %s\n", *mode)
		fmt.Println("Available modes: all, basic, struct, function, interface, practical, utilities")
		os.Exit(1)
	}
}

func runAllExamples() {
	fmt.Println("🎯 Running All Reflection Examples")
	fmt.Println(strings.Repeat("=", 50))

	fmt.Println("\n" + strings.Repeat("=", 60))
	reflect.BasicReflection()

	fmt.Println("\n" + strings.Repeat("=", 60))
	reflect.StructReflection()

	fmt.Println("\n" + strings.Repeat("=", 60))
	reflect.FunctionReflection()

	fmt.Println("\n" + strings.Repeat("=", 60))
	reflect.InterfaceReflection()

	fmt.Println("\n" + strings.Repeat("=", 60))
	reflect.PracticalExamples()

	fmt.Println("\n" + strings.Repeat("=", 60))
	runUtilityExamples()

	fmt.Println("\n🎉 All examples completed!")
	fmt.Println("\n💡 To run specific examples:")
	fmt.Println("  go run run/reflect_main.go -mode=basic")
	fmt.Println("  go run run/reflect_main.go -mode=struct")
	fmt.Println("  go run run/reflect_main.go -mode=function")
	fmt.Println("  go run run/reflect_main.go -mode=interface")
	fmt.Println("  go run run/reflect_main.go -mode=practical")
	fmt.Println("  go run run/reflect_main.go -mode=utilities")
}

func runBasicExamples() {
	fmt.Println("🔍 Basic Reflection Examples")
	fmt.Println(strings.Repeat("=", 50))
	reflect.BasicReflection()
	reflect.DemonstrateTypeChecker()
}

func runStructExamples() {
	fmt.Println("🏗️  Struct Reflection Examples")
	fmt.Println(strings.Repeat("=", 50))
	reflect.StructReflection()
	reflect.DemonstrateStructAnalyzer()
}

func runFunctionExamples() {
	fmt.Println("⚙️  Function Reflection Examples")
	fmt.Println(strings.Repeat("=", 50))
	reflect.FunctionReflection()
	reflect.DemonstrateFunctionRegistry()
}

func runInterfaceExamples() {
	fmt.Println("🔌 Interface Reflection Examples")
	fmt.Println(strings.Repeat("=", 50))
	reflect.InterfaceReflection()
	reflect.DemonstrateInterfaceAnalyzer()
}

func runPracticalExamples() {
	fmt.Println("🛠️  Practical Reflection Examples")
	fmt.Println(strings.Repeat("=", 50))
	reflect.PracticalExamples()
}

func runUtilityExamples() {
	fmt.Println("🔧 Reflection Utility Examples")
	fmt.Println(strings.Repeat("=", 50))

	fmt.Println("\n📊 TypeChecker Utility:")
	reflect.DemonstrateTypeChecker()

	fmt.Println("\n🏗️  StructAnalyzer Utility:")
	reflect.DemonstrateStructAnalyzer()

	fmt.Println("\n⚙️  FunctionRegistry Utility:")
	reflect.DemonstrateFunctionRegistry()

	fmt.Println("\n🔌 InterfaceAnalyzer Utility:")
	reflect.DemonstrateInterfaceAnalyzer()
}
