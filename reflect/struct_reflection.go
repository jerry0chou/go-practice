package reflect

import (
	"fmt"
	"reflect"
	"strings"
)

// User represents a user with various field types and tags
type User struct {
	ID       int                    `json:"id" db:"user_id" validate:"required,min=1"`
	Name     string                 `json:"name" db:"user_name" validate:"required,min=2"`
	Email    string                 `json:"email" db:"email" validate:"required,email"`
	Age      int                    `json:"age" db:"age" validate:"min=0,max=120"`
	Active   bool                   `json:"active" db:"is_active"`
	Tags     []string               `json:"tags" db:"tags"`
	Metadata map[string]interface{} `json:"metadata" db:"metadata"`
}

// Admin extends User with additional fields
type Admin struct {
	User
	Permissions []string `json:"permissions" db:"permissions"`
	Level       int      `json:"level" db:"admin_level" validate:"min=1,max=10"`
}

// Methods for User
func (u *User) GetFullInfo() string {
	return fmt.Sprintf("User: %s (%s), Age: %d, Active: %t", u.Name, u.Email, u.Age, u.Active)
}

func (u *User) IsAdult() bool {
	return u.Age >= 18
}

func (u *User) SetActive(active bool) {
	u.Active = active
}

// Methods for Admin
func (a *Admin) GetAdminInfo() string {
	return fmt.Sprintf("Admin: %s, Level: %d, Permissions: %v", a.Name, a.Level, a.Permissions)
}

// StructReflection demonstrates struct reflect capabilities
func StructReflection() {
	fmt.Println("🏗️  Struct Reflection Examples")
	fmt.Println(strings.Repeat("=", 50))

	// 1. Basic struct inspection
	fmt.Println("\n🔍 1. Basic Struct Inspection:")
	demonstrateStructInspection()

	// 2. Field access and modification
	fmt.Println("\n✏️  2. Field Access and Modification:")
	demonstrateFieldAccess()

	// 3. Struct tags
	fmt.Println("\n🏷️  3. Struct Tags:")
	demonstrateStructTags()

	// 4. Method reflect
	fmt.Println("\n⚙️  4. Method Reflection:")
	demonstrateMethodReflection()

	// 5. Anonymous fields and embedding
	fmt.Println("\n🔗 5. Anonymous Fields and Embedding:")
	demonstrateAnonymousFields()

	// 6. Creating structs dynamically
	fmt.Println("\n🛠️  6. Creating Structs Dynamically:")
	demonstrateDynamicStructCreation()
}

func demonstrateStructInspection() {
	user := User{
		ID:       1,
		Name:     "John Doe",
		Email:    "john@example.com",
		Age:      30,
		Active:   true,
		Tags:     []string{"developer", "golang"},
		Metadata: map[string]interface{}{"department": "engineering", "location": "remote"},
	}

	// Get struct type
	structType := reflect.TypeOf(user)
	fmt.Printf("Struct type: %v\n", structType)
	fmt.Printf("Struct kind: %v\n", structType.Kind())
	fmt.Printf("Number of fields: %d\n", structType.NumField())

	// Iterate through fields
	fmt.Println("\nFields:")
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		fmt.Printf("  %d. %s (%s) - %s\n", i+1, field.Name, field.Type, field.Tag)
	}
}

func demonstrateFieldAccess() {
	user := &User{
		ID:     1,
		Name:   "Jane Smith",
		Email:  "jane@example.com",
		Age:    25,
		Active: false,
	}

	// Get struct value (must be pointer to modify)
	value := reflect.ValueOf(user).Elem()
	structType := value.Type()

	fmt.Printf("Original user: %+v\n", user)

	// Access and modify fields
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		fieldValue := value.Field(i)

		fmt.Printf("Field: %s\n", field.Name)
		fmt.Printf("  Type: %v\n", field.Type)
		fmt.Printf("  Value: %v\n", fieldValue.Interface())
		fmt.Printf("  Can Set: %v\n", fieldValue.CanSet())

		// Modify string fields
		if field.Type.Kind() == reflect.String && fieldValue.CanSet() {
			oldValue := fieldValue.String()
			fieldValue.SetString("Modified: " + oldValue)
			fmt.Printf("  Modified to: %v\n", fieldValue.Interface())
		}

		// Modify int fields
		if field.Type.Kind() == reflect.Int && fieldValue.CanSet() {
			oldValue := fieldValue.Int()
			fieldValue.SetInt(oldValue * 2)
			fmt.Printf("  Modified to: %v\n", fieldValue.Interface())
		}

		fmt.Println()
	}

	fmt.Printf("Modified user: %+v\n", user)
}

func demonstrateStructTags() {
	user := User{}
	structType := reflect.TypeOf(user)

	fmt.Println("Struct tags analysis:")
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		tag := field.Tag

		fmt.Printf("Field: %s\n", field.Name)
		fmt.Printf("  All tags: %s\n", tag)

		// Parse specific tags
		if jsonTag := tag.Get("json"); jsonTag != "" {
			fmt.Printf("  JSON tag: %s\n", jsonTag)
		}
		if dbTag := tag.Get("db"); dbTag != "" {
			fmt.Printf("  DB tag: %s\n", dbTag)
		}
		if validateTag := tag.Get("validate"); validateTag != "" {
			fmt.Printf("  Validate tag: %s\n", validateTag)
		}
		fmt.Println()
	}
}

func demonstrateMethodReflection() {
	user := &User{
		ID:     1,
		Name:   "Alice Johnson",
		Email:  "alice@example.com",
		Age:    28,
		Active: true,
	}

	// Get struct type
	structType := reflect.TypeOf(user)
	fmt.Printf("Methods for %s:\n", structType.Elem().Name())

	// List all methods
	for i := 0; i < structType.NumMethod(); i++ {
		method := structType.Method(i)
		fmt.Printf("  %d. %s\n", i+1, method.Name)
		fmt.Printf("     Type: %v\n", method.Type)
		fmt.Printf("     PkgPath: %s\n", method.PkgPath)
	}

	// Call methods dynamically
	fmt.Println("\nCalling methods dynamically:")

	// Call GetFullInfo
	if method, found := structType.MethodByName("GetFullInfo"); found {
		results := method.Func.Call([]reflect.Value{reflect.ValueOf(user)})
		fmt.Printf("GetFullInfo(): %s\n", results[0].String())
	}

	// Call IsAdult
	if method, found := structType.MethodByName("IsAdult"); found {
		results := method.Func.Call([]reflect.Value{reflect.ValueOf(user)})
		fmt.Printf("IsAdult(): %t\n", results[0].Bool())
	}

	// Call SetActive with parameter
	if method, found := structType.MethodByName("SetActive"); found {
		method.Func.Call([]reflect.Value{reflect.ValueOf(user), reflect.ValueOf(false)})
		fmt.Printf("After SetActive(false): Active = %t\n", user.Active)
	}
}

func demonstrateAnonymousFields() {
	admin := &Admin{
		User: User{
			ID:     1,
			Name:   "Bob Admin",
			Email:  "bob@admin.com",
			Age:    35,
			Active: true,
		},
		Permissions: []string{"read", "write", "delete"},
		Level:       5,
	}

	// Get struct type
	structType := reflect.TypeOf(admin).Elem()
	fmt.Printf("Admin struct type: %v\n", structType)
	fmt.Printf("Number of fields: %d\n", structType.NumField())

	// Iterate through fields (including embedded ones)
	fmt.Println("\nAll fields (including embedded):")
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		fmt.Printf("  %d. %s (%s)\n", i+1, field.Name, field.Type)
		fmt.Printf("     Anonymous: %t\n", field.Anonymous)
		fmt.Printf("     PkgPath: %s\n", field.PkgPath)
	}

	// Access embedded field
	value := reflect.ValueOf(admin).Elem()
	userField := value.FieldByName("User")
	if userField.IsValid() {
		fmt.Printf("\nEmbedded User field: %v\n", userField.Interface())
	}
}

func demonstrateDynamicStructCreation() {
	// Create a new User struct dynamically
	userType := reflect.TypeOf(User{})
	newUser := reflect.New(userType).Elem()

	// Set field values
	if nameField := newUser.FieldByName("Name"); nameField.IsValid() && nameField.CanSet() {
		nameField.SetString("Dynamic User")
	}
	if emailField := newUser.FieldByName("Email"); emailField.IsValid() && emailField.CanSet() {
		emailField.SetString("dynamic@example.com")
	}
	if ageField := newUser.FieldByName("Age"); ageField.IsValid() && ageField.CanSet() {
		ageField.SetInt(25)
	}
	if activeField := newUser.FieldByName("Active"); activeField.IsValid() && activeField.CanSet() {
		activeField.SetBool(true)
	}

	// Create slice and map fields
	if tagsField := newUser.FieldByName("Tags"); tagsField.IsValid() && tagsField.CanSet() {
		tags := reflect.MakeSlice(reflect.TypeOf([]string{}), 2, 2)
		tags.Index(0).SetString("dynamic")
		tags.Index(1).SetString("reflect")
		tagsField.Set(tags)
	}

	if metadataField := newUser.FieldByName("Metadata"); metadataField.IsValid() && metadataField.CanSet() {
		metadata := reflect.MakeMap(reflect.TypeOf(map[string]interface{}{}))
		metadata.SetMapIndex(reflect.ValueOf("created_by"), reflect.ValueOf("reflect"))
		metadata.SetMapIndex(reflect.ValueOf("version"), reflect.ValueOf(1.0))
		metadataField.Set(metadata)
	}

	fmt.Printf("Dynamically created user: %+v\n", newUser.Interface())
}

// StructAnalyzer provides utility functions for struct analysis
type StructAnalyzer struct{}

// AnalyzeStruct provides comprehensive struct analysis
func (sa *StructAnalyzer) AnalyzeStruct(s interface{}) map[string]interface{} {
	structType := reflect.TypeOf(s)
	structValue := reflect.ValueOf(s)

	// Handle pointers
	if structType.Kind() == reflect.Ptr {
		structType = structType.Elem()
		structValue = structValue.Elem()
	}

	analysis := map[string]interface{}{
		"name":       structType.Name(),
		"package":    structType.PkgPath(),
		"kind":       structType.Kind().String(),
		"numFields":  structType.NumField(),
		"numMethods": structType.NumMethod(),
		"fields":     []map[string]interface{}{},
		"methods":    []map[string]interface{}{},
		"embedded":   []string{},
	}

	// Analyze fields
	fields := []map[string]interface{}{}
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		fieldValue := structValue.Field(i)

		fieldInfo := map[string]interface{}{
			"name":      field.Name,
			"type":      field.Type.String(),
			"kind":      field.Type.Kind().String(),
			"anonymous": field.Anonymous,
			"tag":       string(field.Tag),
			"value":     fieldValue.Interface(),
			"canSet":    fieldValue.CanSet(),
			"isZero":    fieldValue.IsZero(),
		}

		// Parse tags
		tags := map[string]string{}
		if jsonTag := field.Tag.Get("json"); jsonTag != "" {
			tags["json"] = jsonTag
		}
		if dbTag := field.Tag.Get("db"); dbTag != "" {
			tags["db"] = dbTag
		}
		if validateTag := field.Tag.Get("validate"); validateTag != "" {
			tags["validate"] = validateTag
		}
		fieldInfo["parsedTags"] = tags

		fields = append(fields, fieldInfo)

		if field.Anonymous {
			analysis["embedded"] = append(analysis["embedded"].([]string), field.Name)
		}
	}
	analysis["fields"] = fields

	// Analyze methods
	methods := []map[string]interface{}{}
	for i := 0; i < structType.NumMethod(); i++ {
		method := structType.Method(i)
		methodInfo := map[string]interface{}{
			"name":    method.Name,
			"type":    method.Type.String(),
			"pkgPath": method.PkgPath,
		}
		methods = append(methods, methodInfo)
	}
	analysis["methods"] = methods

	return analysis
}

// DemonstrateStructAnalyzer shows how to use the StructAnalyzer utility
func DemonstrateStructAnalyzer() {
	fmt.Println("\n🔧 StructAnalyzer Utility:")
	fmt.Println(strings.Repeat("-", 30))

	analyzer := &StructAnalyzer{}

	user := User{
		ID:       1,
		Name:     "Test User",
		Email:    "test@example.com",
		Age:      30,
		Active:   true,
		Tags:     []string{"test"},
		Metadata: map[string]interface{}{"test": true},
	}

	admin := Admin{
		User:        user,
		Permissions: []string{"admin"},
		Level:       3,
	}

	structs := []interface{}{user, admin}

	for i, s := range structs {
		fmt.Printf("Analysis %d: %T\n", i+1, s)
		analysis := analyzer.AnalyzeStruct(s)

		fmt.Printf("  Name: %s\n", analysis["name"])
		fmt.Printf("  Package: %s\n", analysis["package"])
		fmt.Printf("  Fields: %d\n", analysis["numFields"])
		fmt.Printf("  Methods: %d\n", analysis["numMethods"])

		if embedded := analysis["embedded"].([]string); len(embedded) > 0 {
			fmt.Printf("  Embedded: %v\n", embedded)
		}

		fmt.Println()
	}
}
