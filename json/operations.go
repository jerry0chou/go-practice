package json

import (
	"encoding/json"
	"fmt"
	"time"
)

type Person struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Age       int       `json:"age"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	Tags      []string  `json:"tags,omitempty"`
	Address   *Address  `json:"address,omitempty"`
}

type Address struct {
	Street  string `json:"street"`
	City    string `json:"city"`
	Country string `json:"country"`
	ZipCode string `json:"zip_code"`
}

type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	Description string  `json:"description,omitempty"`
	InStock     bool    `json:"in_stock"`
}

type CustomJSONData struct {
	Value string `json:"value"`
}

func (c CustomJSONData) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"custom_value": c.Value,
		"timestamp":    time.Now().Unix(),
	})
}

func (c *CustomJSONData) UnmarshalJSON(data []byte) error {
	var temp map[string]any
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	if val, ok := temp["custom_value"].(string); ok {
		c.Value = val
	}
	return nil
}

// BasicJSONOperations demonstrates basic JSON operations
func BasicJSONOperations() {
	fmt.Println("=== Basic JSON Operations ===")

	person := Person{
		ID:        1,
		Name:      "John Doe",
		Email:     "john@example.com",
		Age:       30,
		IsActive:  true,
		CreatedAt: time.Now(),
		Tags:      []string{"developer", "golang"},
		Address: &Address{
			Street:  "123 Main St",
			City:    "New York",
			Country: "USA",
			ZipCode: "10001",
		},
	}

	jsonData, err := json.Marshal(person)
	if err != nil {
		fmt.Printf("Error marshaling: %v\n", err)
		return
	}

	fmt.Printf("Marshaled JSON:\n%s\n\n", string(jsonData))

	jsonDataPretty, err := json.MarshalIndent(person, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling with indent: %v\n", err)
		return
	}

	fmt.Printf("Pretty JSON:\n%s\n\n", string(jsonDataPretty))

	var newPerson Person
	err = json.Unmarshal(jsonData, &newPerson)
	if err != nil {
		fmt.Printf("Error unmarshaling: %v\n", err)
		return
	}

	fmt.Printf("Unmarshaled person: %+v\n\n", newPerson)
}

// JSONWithMaps demonstrates working with JSON and maps
func JSONWithMaps() {
	fmt.Println("=== JSON with Maps ===")

	data := map[string]any{
		"name":    "Alice",
		"age":     25,
		"city":    "San Francisco",
		"hobbies": []string{"reading", "swimming", "coding"},
		"address": map[string]string{
			"street": "456 Oak Ave",
			"zip":    "94102",
		},
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Printf("Error marshaling map: %v\n", err)
		return
	}

	fmt.Printf("Map to JSON:\n%s\n\n", string(jsonData))

	var result map[string]any
	err = json.Unmarshal(jsonData, &result)
	if err != nil {
		fmt.Printf("Error unmarshaling to map: %v\n", err)
		return
	}

	fmt.Printf("JSON to map: %+v\n\n", result)

	if name, ok := result["name"].(string); ok {
		fmt.Printf("Name: %s\n", name)
	}

	if hobbies, ok := result["hobbies"].([]any); ok {
		fmt.Printf("Hobbies: ")
		for i, hobby := range hobbies {
			if i > 0 {
				fmt.Print(", ")
			}
			fmt.Print(hobby)
		}
		fmt.Println()
	}
}

// JSONWithSlices demonstrates working with JSON and slices
func JSONWithSlices() {
	fmt.Println("=== JSON with Slices ===")

	products := []Product{
		{ID: 1, Name: "Laptop", Price: 999.99, Description: "High-performance laptop", InStock: true},
		{ID: 2, Name: "Mouse", Price: 29.99, InStock: true},
		{ID: 3, Name: "Keyboard", Price: 79.99, Description: "Mechanical keyboard", InStock: false},
	}

	jsonData, err := json.MarshalIndent(products, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling slice: %v\n", err)
		return
	}

	fmt.Printf("Products JSON:\n%s\n\n", string(jsonData))

	var newProducts []Product
	err = json.Unmarshal(jsonData, &newProducts)
	if err != nil {
		fmt.Printf("Error unmarshaling slice: %v\n", err)
		return
	}

	fmt.Printf("Unmarshaled products:\n")
	for _, product := range newProducts {
		fmt.Printf("- %s: $%.2f (In Stock: %t)\n", product.Name, product.Price, product.InStock)
	}
	fmt.Println()
}

// CustomJSONHandling demonstrates custom JSON marshaling/unmarshaling
func CustomJSONHandling() {
	fmt.Println("=== Custom JSON Handling ===")

	customData := CustomJSONData{Value: "Hello, World!"}

	jsonData, err := json.Marshal(customData)
	if err != nil {
		fmt.Printf("Error marshaling custom data: %v\n", err)
		return
	}

	fmt.Printf("Custom marshaled JSON:\n%s\n\n", string(jsonData))

	var newCustomData CustomJSONData
	err = json.Unmarshal(jsonData, &newCustomData)
	if err != nil {
		fmt.Printf("Error unmarshaling custom data: %v\n", err)
		return
	}

	fmt.Printf("Custom unmarshaled data: %+v\n\n", newCustomData)
}

// JSONStreaming demonstrates streaming JSON operations
func JSONStreaming() {
	fmt.Println("=== JSON Streaming ===")

	people := []Person{
		{ID: 1, Name: "John", Email: "john@example.com", Age: 30, IsActive: true, CreatedAt: time.Now()},
		{ID: 2, Name: "Jane", Email: "jane@example.com", Age: 25, IsActive: true, CreatedAt: time.Now()},
		{ID: 3, Name: "Bob", Email: "bob@example.com", Age: 35, IsActive: false, CreatedAt: time.Now()},
	}

	fmt.Println("Streaming JSON output:")

	for i, person := range people {
		jsonData, err := json.Marshal(person)
		if err != nil {
			fmt.Printf("Error marshaling person %d: %v\n", i+1, err)
			continue
		}

		fmt.Printf("Person %d: %s\n", i+1, string(jsonData))
	}
	fmt.Println()
}

// JSONValidation demonstrates JSON validation and error handling
func JSONValidation() {
	fmt.Println("=== JSON Validation ===")

	validJSON := `{"name": "Test", "age": 30, "active": true}`

	invalidJSON := `{"name": "Test", "age": 30, "active": true`

	var validData map[string]any
	err := json.Unmarshal([]byte(validJSON), &validData)
	if err != nil {
		fmt.Printf("Valid JSON error: %v\n", err)
	} else {
		fmt.Printf("Valid JSON parsed successfully: %+v\n", validData)
	}

	var invalidData map[string]any
	err = json.Unmarshal([]byte(invalidJSON), &invalidData)
	if err != nil {
		fmt.Printf("Invalid JSON error: %v\n", err)
	} else {
		fmt.Printf("Invalid JSON parsed successfully: %+v\n", invalidData)
	}
	fmt.Println()
}

// JSONTags demonstrates different JSON tag options
func JSONTags() {
	fmt.Println("=== JSON Tags ===")

	type TagExample struct {
		ID          int     `json:"id"`
		Name        string  `json:"name"`
		Email       string  `json:"email,omitempty"`
		Description string  `json:"-"`
		Internal    string  `json:"internal,omitempty"`
		Price       float64 `json:"price,string"`
	}

	tagExample := TagExample{
		ID:          1,
		Name:        "Test Product",
		Description: "This will be ignored",
		Internal:    "",
		Price:       99.99,
	}

	jsonData, err := json.MarshalIndent(tagExample, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling tag example: %v\n", err)
		return
	}

	fmt.Printf("Tag example JSON:\n%s\n\n", string(jsonData))
}

// JSONPointer demonstrates working with JSON pointers
func JSONPointer() {
	fmt.Println("=== JSON Pointer ===")

	complexJSON := `{
		"user": {
			"id": 1,
			"profile": {
				"name": "John Doe",
				"settings": {
					"theme": "dark",
					"notifications": true
				}
			}
		},
		"posts": [
			{"id": 1, "title": "First Post"},
			{"id": 2, "title": "Second Post"}
		]
	}`

	var data map[string]any
	err := json.Unmarshal([]byte(complexJSON), &data)
	if err != nil {
		fmt.Printf("Error unmarshaling complex JSON: %v\n", err)
		return
	}

	if user, ok := data["user"].(map[string]any); ok {
		if profile, ok := user["profile"].(map[string]any); ok {
			if name, ok := profile["name"].(string); ok {
				fmt.Printf("User name: %s\n", name)
			}

			if settings, ok := profile["settings"].(map[string]any); ok {
				if theme, ok := settings["theme"].(string); ok {
					fmt.Printf("Theme: %s\n", theme)
				}
			}
		}
	}

	if posts, ok := data["posts"].([]any); ok {
		fmt.Printf("Number of posts: %d\n", len(posts))
		for i, post := range posts {
			if postMap, ok := post.(map[string]any); ok {
				if title, ok := postMap["title"].(string); ok {
					fmt.Printf("Post %d: %s\n", i+1, title)
				}
			}
		}
	}
	fmt.Println()
}

// RunAllJSONExamples runs all JSON examples
func RunAllJSONExamples() {
	fmt.Println("Go JSON Operations with encoding/json")
	fmt.Println("=====================================")
	fmt.Println()

	BasicJSONOperations()
	JSONWithMaps()
	JSONWithSlices()
	CustomJSONHandling()
	JSONStreaming()
	JSONValidation()
	JSONTags()
	JSONPointer()

	fmt.Println("All JSON examples completed!")
}
