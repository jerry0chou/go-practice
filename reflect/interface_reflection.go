package reflect

import (
	"fmt"
	"reflect"
	"strings"
)

// Interface definitions for demonstration
type Writer interface {
	Write(data []byte) (int, error)
	Close() error
}

type Reader interface {
	Read(data []byte) (int, error)
	Close() error
}

type ReadWriter interface {
	Reader
	Writer
}

type Closer interface {
	Close() error
}

type Stringer interface {
	String() string
}

// Concrete implementations
type FileWriter struct {
	filename string
	closed   bool
}

func (fw *FileWriter) Write(data []byte) (int, error) {
	if fw.closed {
		return 0, fmt.Errorf("file is closed")
	}
	fmt.Printf("Writing %d bytes to %s\n", len(data), fw.filename)
	return len(data), nil
}

func (fw *FileWriter) Close() error {
	if fw.closed {
		return fmt.Errorf("file already closed")
	}
	fw.closed = true
	fmt.Printf("Closed file: %s\n", fw.filename)
	return nil
}

func (fw *FileWriter) String() string {
	return fmt.Sprintf("FileWriter{filename: %s, closed: %t}", fw.filename, fw.closed)
}

type BufferReader struct {
	buffer []byte
	pos    int
	closed bool
}

func (br *BufferReader) Read(data []byte) (int, error) {
	if br.closed {
		return 0, fmt.Errorf("buffer is closed")
	}
	if br.pos >= len(br.buffer) {
		return 0, fmt.Errorf("EOF")
	}

	n := copy(data, br.buffer[br.pos:])
	br.pos += n
	return n, nil
}

func (br *BufferReader) Close() error {
	if br.closed {
		return fmt.Errorf("buffer already closed")
	}
	br.closed = true
	fmt.Println("Closed buffer reader")
	return nil
}

func (br *BufferReader) String() string {
	return fmt.Sprintf("BufferReader{buffer: %d bytes, pos: %d, closed: %t}",
		len(br.buffer), br.pos, br.closed)
}

type MemoryReadWriter struct {
	data   []byte
	pos    int
	closed bool
}

func (mrw *MemoryReadWriter) Read(data []byte) (int, error) {
	if mrw.closed {
		return 0, fmt.Errorf("memory is closed")
	}
	if mrw.pos >= len(mrw.data) {
		return 0, fmt.Errorf("EOF")
	}

	n := copy(data, mrw.data[mrw.pos:])
	mrw.pos += n
	return n, nil
}

func (mrw *MemoryReadWriter) Write(data []byte) (int, error) {
	if mrw.closed {
		return 0, fmt.Errorf("memory is closed")
	}

	// Append data to the end
	mrw.data = append(mrw.data, data...)
	return len(data), nil
}

func (mrw *MemoryReadWriter) Close() error {
	if mrw.closed {
		return fmt.Errorf("memory already closed")
	}
	mrw.closed = true
	fmt.Println("Closed memory read-writer")
	return nil
}

func (mrw *MemoryReadWriter) String() string {
	return fmt.Sprintf("MemoryReadWriter{data: %d bytes, pos: %d, closed: %t}",
		len(mrw.data), mrw.pos, mrw.closed)
}

// InterfaceReflection demonstrates interface reflect capabilities
func InterfaceReflection() {
	fmt.Println("🔌 Interface Reflection Examples")
	fmt.Println(strings.Repeat("=", 50))

	// 1. Basic interface inspection
	fmt.Println("\n🔍 1. Basic Interface Inspection:")
	demonstrateInterfaceInspection()

	// 2. Interface type assertions
	fmt.Println("\n✅ 2. Interface Type Assertions:")
	demonstrateInterfaceTypeAssertions()

	// 3. Interface method reflect
	fmt.Println("\n⚙️  3. Interface Method Reflection:")
	demonstrateInterfaceMethodReflection()

	// 4. Interface composition
	fmt.Println("\n🔗 4. Interface Composition:")
	demonstrateInterfaceComposition()

	// 5. Dynamic interface implementation
	fmt.Println("\n🛠️  5. Dynamic Interface Implementation:")
	demonstrateDynamicInterfaceImplementation()

	// 6. Interface value inspection
	fmt.Println("\n🔬 6. Interface Value Inspection:")
	demonstrateInterfaceValueInspection()
}

func demonstrateInterfaceInspection() {
	// Create instances
	fileWriter := &FileWriter{filename: "test.txt"}
	bufferReader := &BufferReader{buffer: []byte("Hello, World!")}
	memoryRW := &MemoryReadWriter{data: []byte("Initial data")}

	// Store in interface variables
	var writer Writer = fileWriter
	var reader Reader = bufferReader
	var readWriter ReadWriter = memoryRW
	var closer Closer = fileWriter
	var stringer Stringer = fileWriter

	interfaces := []interface{}{
		writer,
		reader,
		readWriter,
		closer,
		stringer,
	}

	interfaceNames := []string{
		"Writer",
		"Reader",
		"ReadWriter",
		"Closer",
		"Stringer",
	}

	for i, iface := range interfaces {
		fmt.Printf("Interface: %s\n", interfaceNames[i])
		fmt.Printf("  Type: %T\n", iface)
		fmt.Printf("  Value: %v\n", iface)
		fmt.Printf("  Is nil: %t\n", iface == nil)

		// Get interface type
		ifaceType := reflect.TypeOf(iface)
		fmt.Printf("  Interface type: %v\n", ifaceType)
		fmt.Printf("  Kind: %v\n", ifaceType.Kind())
		fmt.Printf("  Num method: %d\n", ifaceType.NumMethod())

		// List methods
		if ifaceType.NumMethod() > 0 {
			fmt.Println("  Methods:")
			for j := 0; j < ifaceType.NumMethod(); j++ {
				method := ifaceType.Method(j)
				fmt.Printf("    %s: %v\n", method.Name, method.Type)
			}
		}
		fmt.Println()
	}
}

func demonstrateInterfaceTypeAssertions() {
	// Create instances and store in interface{}
	var value interface{} = &FileWriter{filename: "example.txt"}

	fmt.Printf("Value: %v\n", value)
	fmt.Printf("Type: %T\n", value)

	// Type assertion to specific interface
	if writer, ok := value.(Writer); ok {
		fmt.Printf("✅ Successfully asserted as Writer: %v\n", writer)
		writer.Write([]byte("test data"))
	} else {
		fmt.Println("❌ Failed to assert as Writer")
	}

	// Type assertion to concrete type
	if fileWriter, ok := value.(*FileWriter); ok {
		fmt.Printf("✅ Successfully asserted as *FileWriter: %v\n", fileWriter)
		fmt.Printf("Filename: %s\n", fileWriter.filename)
	} else {
		fmt.Println("❌ Failed to assert as *FileWriter")
	}

	// Type assertion to different interface
	if stringer, ok := value.(Stringer); ok {
		fmt.Printf("✅ Successfully asserted as Stringer: %v\n", stringer.String())
	} else {
		fmt.Println("❌ Failed to assert as Stringer")
	}

	// Type assertion to wrong type
	if reader, ok := value.(Reader); ok {
		fmt.Printf("✅ Successfully asserted as Reader: %v\n", reader)
	} else {
		fmt.Println("❌ Failed to assert as Reader (expected)")
	}

	// Using reflect for type checking
	fmt.Println("\nUsing reflect for type checking:")
	valueType := reflect.TypeOf(value)
	writerType := reflect.TypeOf((*Writer)(nil)).Elem()
	readerType := reflect.TypeOf((*Reader)(nil)).Elem()

	fmt.Printf("Value type: %v\n", valueType)
	fmt.Printf("Writer interface type: %v\n", writerType)
	fmt.Printf("Reader interface type: %v\n", readerType)
	fmt.Printf("Implements Writer: %t\n", valueType.Implements(writerType))
	fmt.Printf("Implements Reader: %t\n", valueType.Implements(readerType))
}

func demonstrateInterfaceMethodReflection() {
	// Create instance and store in interface
	var writer Writer = &FileWriter{filename: "method_test.txt"}

	// Get interface value
	writerValue := reflect.ValueOf(writer)
	writerType := writerValue.Type()

	fmt.Printf("Writer type: %v\n", writerType)
	fmt.Printf("Number of methods: %d\n", writerType.NumMethod())

	// List and call methods
	for i := 0; i < writerType.NumMethod(); i++ {
		method := writerType.Method(i)
		fmt.Printf("\nMethod %d: %s\n", i+1, method.Name)
		fmt.Printf("  Type: %v\n", method.Type)
		fmt.Printf("  PkgPath: %s\n", method.PkgPath)

		// Call the method
		methodValue := writerValue.Method(i)
		if methodValue.IsValid() {
			fmt.Printf("  Can call: %t\n", methodValue.CanInterface())

			// Call Write method
			if method.Name == "Write" {
				args := []reflect.Value{reflect.ValueOf([]byte("Hello from reflect!"))}
				results := methodValue.Call(args)
				fmt.Printf("  Write result: %v, %v\n", results[0].Int(), results[1].Interface())
			}

			// Call String method
			if method.Name == "String" {
				results := methodValue.Call([]reflect.Value{})
				fmt.Printf("  String result: %v\n", results[0].String())
			}
		}
	}
}

func demonstrateInterfaceComposition() {
	// Create a ReadWriter instance
	memoryRW := &MemoryReadWriter{data: []byte("Composition test")}
	var readWriter ReadWriter = memoryRW

	fmt.Printf("ReadWriter: %v\n", readWriter)

	// Check what interfaces it implements
	interfaces := []reflect.Type{
		reflect.TypeOf((*Reader)(nil)).Elem(),
		reflect.TypeOf((*Writer)(nil)).Elem(),
		reflect.TypeOf((*ReadWriter)(nil)).Elem(),
		reflect.TypeOf((*Closer)(nil)).Elem(),
		reflect.TypeOf((*Stringer)(nil)).Elem(),
	}

	interfaceNames := []string{"Reader", "Writer", "ReadWriter", "Closer", "Stringer"}
	readWriterType := reflect.TypeOf(memoryRW)

	fmt.Println("\nInterface implementation check:")
	for i, ifaceType := range interfaces {
		implements := readWriterType.Implements(ifaceType)
		fmt.Printf("  Implements %s: %t\n", interfaceNames[i], implements)
	}

	// Demonstrate method resolution
	fmt.Println("\nMethod resolution:")
	readWriterValue := reflect.ValueOf(readWriter)
	readWriterReflectType := readWriterValue.Type()

	for i := 0; i < readWriterReflectType.NumMethod(); i++ {
		method := readWriterReflectType.Method(i)
		fmt.Printf("  Method: %s (%v)\n", method.Name, method.Type)
	}
}

func demonstrateDynamicInterfaceImplementation() {
	// Create a dynamic implementation of Writer interface
	dynamicWriter := &DynamicWriter{
		name: "Dynamic Writer",
		data: make([]byte, 0),
	}

	var writer Writer = dynamicWriter

	fmt.Printf("Dynamic writer: %v\n", writer)

	// Test the dynamic implementation
	testData := []byte("Dynamic test data")
	n, err := writer.Write(testData)
	fmt.Printf("Write result: %d bytes, error: %v\n", n, err)

	// Close the writer
	err = writer.Close()
	fmt.Printf("Close result: error: %v\n", err)

	// Test String method
	if stringer, ok := writer.(Stringer); ok {
		fmt.Printf("String representation: %v\n", stringer.String())
	}
}

// DynamicWriter is a dynamic implementation of Writer interface
type DynamicWriter struct {
	name   string
	data   []byte
	closed bool
}

func (dw *DynamicWriter) Write(data []byte) (int, error) {
	if dw.closed {
		return 0, fmt.Errorf("writer is closed")
	}
	dw.data = append(dw.data, data...)
	fmt.Printf("DynamicWriter[%s] wrote %d bytes\n", dw.name, len(data))
	return len(data), nil
}

func (dw *DynamicWriter) Close() error {
	if dw.closed {
		return fmt.Errorf("writer already closed")
	}
	dw.closed = true
	fmt.Printf("DynamicWriter[%s] closed\n", dw.name)
	return nil
}

func (dw *DynamicWriter) String() string {
	return fmt.Sprintf("DynamicWriter{name: %s, data: %d bytes, closed: %t}",
		dw.name, len(dw.data), dw.closed)
}

func demonstrateInterfaceValueInspection() {
	// Create various interface values
	var nilInterface interface{}
	var writer Writer = &FileWriter{filename: "inspection.txt"}
	var reader Reader = &BufferReader{buffer: []byte("inspection data")}
	var emptyInterface interface{} = 42

	values := []interface{}{
		nilInterface,
		writer,
		reader,
		emptyInterface,
	}

	valueNames := []string{
		"nil interface",
		"Writer interface",
		"Reader interface",
		"empty interface with int",
	}

	for i, value := range values {
		fmt.Printf("Value %d: %s\n", i+1, valueNames[i])
		fmt.Printf("  Value: %v\n", value)
		fmt.Printf("  Type: %T\n", value)
		fmt.Printf("  Is nil: %t\n", value == nil)

		// Reflection analysis
		valueReflect := reflect.ValueOf(value)
		fmt.Printf("  Reflect value: %v\n", valueReflect)
		fmt.Printf("  Is valid: %t\n", valueReflect.IsValid())

		if valueReflect.IsValid() {
			fmt.Printf("  Reflect type: %v\n", valueReflect.Type())
			fmt.Printf("  Reflect kind: %v\n", valueReflect.Kind())
			fmt.Printf("  Is zero: %t\n", valueReflect.IsZero())
			fmt.Printf("  Can interface: %t\n", valueReflect.CanInterface())

			// Interface-specific analysis
			if valueReflect.Kind() == reflect.Interface {
				fmt.Printf("  Interface value: %v\n", valueReflect.Interface())
				fmt.Printf("  Interface type: %v\n", valueReflect.Type())

				// Get the underlying value
				if valueReflect.Elem().IsValid() {
					fmt.Printf("  Underlying value: %v\n", valueReflect.Elem())
					fmt.Printf("  Underlying type: %v\n", valueReflect.Elem().Type())
				}
			}
		} else {
			fmt.Printf("  Reflect type: <invalid>\n")
			fmt.Printf("  Reflect kind: <invalid>\n")
			fmt.Printf("  Is zero: true\n")
			fmt.Printf("  Can interface: false\n")
		}
		fmt.Println()
	}
}

// InterfaceAnalyzer provides utility functions for interface analysis
type InterfaceAnalyzer struct{}

// AnalyzeInterface provides comprehensive interface analysis
func (ia *InterfaceAnalyzer) AnalyzeInterface(value interface{}) map[string]interface{} {
	valueReflect := reflect.ValueOf(value)

	analysis := map[string]interface{}{
		"value":      value,
		"isNil":      value == nil,
		"isValid":    valueReflect.IsValid(),
		"methods":    []map[string]interface{}{},
		"implements": []string{},
	}

	// Handle nil values
	if value == nil || !valueReflect.IsValid() {
		analysis["type"] = "<nil>"
		analysis["kind"] = "invalid"
		analysis["isZero"] = true
		analysis["canInterface"] = false
		return analysis
	}

	typeReflect := valueReflect.Type()
	analysis["type"] = typeReflect.String()
	analysis["kind"] = typeReflect.Kind().String()
	analysis["isZero"] = valueReflect.IsZero()
	analysis["canInterface"] = valueReflect.CanInterface()

	// Analyze methods
	if typeReflect.Kind() == reflect.Interface || typeReflect.NumMethod() > 0 {
		methods := []map[string]interface{}{}
		for i := 0; i < typeReflect.NumMethod(); i++ {
			method := typeReflect.Method(i)
			methodInfo := map[string]interface{}{
				"name":    method.Name,
				"type":    method.Type.String(),
				"pkgPath": method.PkgPath,
			}
			methods = append(methods, methodInfo)
		}
		analysis["methods"] = methods
	}

	// Check common interface implementations
	commonInterfaces := []struct {
		name string
		typ  reflect.Type
	}{
		{"Stringer", reflect.TypeOf((*Stringer)(nil)).Elem()},
		{"Writer", reflect.TypeOf((*Writer)(nil)).Elem()},
		{"Reader", reflect.TypeOf((*Reader)(nil)).Elem()},
		{"Closer", reflect.TypeOf((*Closer)(nil)).Elem()},
	}

	implements := []string{}
	for _, iface := range commonInterfaces {
		if typeReflect.Implements(iface.typ) {
			implements = append(implements, iface.name)
		}
	}
	analysis["implements"] = implements

	return analysis
}

// DemonstrateInterfaceAnalyzer shows how to use the InterfaceAnalyzer utility
func DemonstrateInterfaceAnalyzer() {
	fmt.Println("\n🔧 InterfaceAnalyzer Utility:")
	fmt.Println(strings.Repeat("-", 30))

	analyzer := &InterfaceAnalyzer{}

	// Test various interface values
	values := []interface{}{
		nil,
		&FileWriter{filename: "analyzer.txt"},
		&BufferReader{buffer: []byte("analyzer data")},
		&MemoryReadWriter{data: []byte("analyzer rw")},
		42,
		"hello",
	}

	for i, value := range values {
		fmt.Printf("Analysis %d: %T\n", i+1, value)
		analysis := analyzer.AnalyzeInterface(value)

		for key, val := range analysis {
			if key == "methods" {
				methods := val.([]map[string]interface{})
				fmt.Printf("  %s: %d methods\n", key, len(methods))
				for j, method := range methods {
					fmt.Printf("    %d. %s\n", j+1, method["name"])
				}
			} else if key == "implements" {
				implements := val.([]string)
				fmt.Printf("  %s: %v\n", key, implements)
			} else {
				fmt.Printf("  %s: %v\n", key, val)
			}
		}
		fmt.Println()
	}
}
