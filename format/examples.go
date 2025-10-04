package format

import (
	"fmt"
	"os"
)

type Person struct {
	Name string
	Age  int
}

func (p Person) String() string {
	return fmt.Sprintf("Person{Name: %q, Age: %d}", p.Name, p.Age)
}

func BasicOutput() {
	fmt.Println("=== Basic Output ===")
	fmt.Print("Hello")
	fmt.Print("World")
	fmt.Println()
	fmt.Println("Hello", "World")
	fmt.Printf("%s %s\n", "Hello", "World")
}

func FormattedOutput() {
	fmt.Println("\n=== Formatted Output ===")
	num := 42
	fmt.Printf("Decimal: %d, Binary: %b, Octal: %o, Hex: %x\n", num, num, num, num)

	pi := 3.14159
	fmt.Printf("Default: %f, 2 decimals: %.2f, Scientific: %e\n", pi, pi, pi)

	str := "GoLang"
	fmt.Printf("String: %s, Quoted: %q, Type: %T\n", str, str, str)
	fmt.Printf("Boolean: %t\n", true)
}

func StructFormatting() {
	fmt.Println("\n=== Struct Formatting ===")
	p := Person{"Jerry", 25}
	fmt.Printf("%%v: %v\n", p)
	fmt.Printf("%%+v: %+v\n", p)
	fmt.Printf("%%#v: %#v\n", p)
}

func InputPractice() {
	fmt.Println("\n=== Input Practice ===")
	var name string
	var age int
	fmt.Print("Enter name and age (e.g. Jerry 25): ")
	fmt.Scan(&name, &age)
	fmt.Printf("Result -> Name: %s, Age: %d\n", name, age)
}

func StringFormatting() {
	fmt.Println("\n=== String Formatting ===")
	name, age := "Jerry", 25
	s := fmt.Sprintf("Name: %s, Age: %d", name, age)
	fmt.Println("Sprintf result:", s)
	s2 := fmt.Sprintln("Hello", "World")
	fmt.Println("Sprintln result:", s2)
}

func FileOutput() {
	fmt.Println("\n=== Output to File & Error Stream ===")
	file, err := os.Create("output.txt")
	if err != nil {
		fmt.Println("Failed to create file:", err)
		return
	}
	defer file.Close()
	fmt.Fprint(file, "Hello File\n")
	fmt.Fprintf(os.Stderr, "error: %s\n", "something went wrong")
}

func AdvancedPractice() {
	fmt.Println("\n=== Advanced Practice ===")
	var a, b int
	fmt.Print("Enter two integers (e.g. 8 3): ")
	fmt.Scan(&a, &b)
	fmt.Printf("%d + %d = %d\n", a, b, a+b)
	fmt.Printf("%d - %d = %d\n", a, b, a-b)
	fmt.Printf("%d * %d = %d\n", a, b, a*b)
	if b != 0 {
		fmt.Printf("%d / %d = %.2f\n", a, b, float64(a)/float64(b))
	}
}

func WidthAndPrecision() {
	fmt.Println("\n=== Width and Precision ===")
	number := 123.456789
	fmt.Printf("Default: %f\n", number)
	fmt.Printf("Width 10: %10f\n", number)
	fmt.Printf("Precision 2: %.2f\n", number)
	fmt.Printf("Width 10, Precision 2: %10.2f\n", number)
	fmt.Printf("Left aligned: %-10.2f|\n", number)
	fmt.Printf("Zero padded: %010.2f\n", number)
}

func StringAndRuneFormatting() {
	fmt.Println("\n=== String and Rune Formatting ===")
	text := "Hello, 世界"
	fmt.Printf("String: %s\n", text)
	fmt.Printf("Quoted: %q\n", text)
	fmt.Printf("ASCII only: %+q\n", text)
	fmt.Printf("Unicode: %U\n", '世')
	fmt.Printf("Character: %c\n", '世')
}

func PointerAndInterfaceFormatting() {
	fmt.Println("\n=== Pointer and Interface Formatting ===")
	x := 42
	ptr := &x
	fmt.Printf("Value: %v, Pointer: %p\n", x, ptr)
	fmt.Printf("Pointer value: %v\n", *ptr)

	var empty interface{}
	empty = "Hello"
	fmt.Printf("Interface: %v, Type: %T\n", empty, empty)
}

func CustomFormatting() {
	fmt.Println("\n=== Custom Formatting ===")
	person2 := Person{"Alice", 30}
	fmt.Printf("Person: %v\n", person2)
	fmt.Printf("Person (detailed): %+v\n", person2)
}

func ScanVariations() {
	fmt.Println("\n=== Scan Variations ===")
	fmt.Print("Enter a line of text: ")
	var line string
	fmt.Scanln(&line)
	fmt.Printf("You entered: %q\n", line)

	fmt.Print("Enter multiple words: ")
	var word1, word2, word3 string
	fmt.Scanf("%s %s %s", &word1, &word2, &word3)
	fmt.Printf("Words: %q, %q, %q\n", word1, word2, word3)
}

func RunAllExamples() {
	BasicOutput()
	FormattedOutput()
	StructFormatting()
	InputPractice()
	StringFormatting()
	FileOutput()
	AdvancedPractice()
	WidthAndPrecision()
	StringAndRuneFormatting()
	PointerAndInterfaceFormatting()
	CustomFormatting()
	ScanVariations()
}
