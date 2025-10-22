package main

import (
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql" // MySQL driver
	"github.com/jerrychou/go-practice/database"
	_ "github.com/lib/pq"           // PostgreSQL driver
	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

func main() {
	log.Println("=== Database Operations Demo ===")

	// Create database examples instance
	examples := database.NewDatabaseExamples()

	// Run all examples
	if err := examples.RunAllExamples(); err != nil {
		log.Fatalf("Failed to run database examples: %v", err)
	}

	// Demonstrate with SQLite (if available)
	if err := demonstrateWithSQLite(); err != nil {
		log.Printf("SQLite demonstration failed: %v", err)
		log.Println("This is expected if SQLite driver is not installed")
	}

	// Demonstrate with PostgreSQL (if available)
	if err := demonstrateWithPostgreSQL(); err != nil {
		log.Printf("PostgreSQL demonstration failed: %v", err)
		log.Println("This is expected if PostgreSQL is not running")
	}

	log.Println("=== Database Operations Demo Completed ===")
}

// demonstrateWithSQLite demonstrates database operations with SQLite
func demonstrateWithSQLite() error {
	log.Println("\n--- SQLite Demonstration ---")

	// Create SQLite database
	dbPath := "test.db"
	defer os.Remove(dbPath) // Clean up after demo

	// Connect to SQLite
	drivers := database.NewDatabaseDrivers()
	if err := drivers.ConnectSQLite(dbPath); err != nil {
		return fmt.Errorf("failed to connect to SQLite: %w", err)
	}
	defer drivers.CloseAllConnections()

	db := drivers.GetSQLiteDB()
	if db == nil {
		return fmt.Errorf("SQLite connection not available")
	}

	// Test SQLite connection
	if err := drivers.TestSQLiteConnection(); err != nil {
		return fmt.Errorf("failed to test SQLite connection: %w", err)
	}

	// Run SQL basics examples
	examples := database.NewDatabaseExamples()
	if err := examples.RunSQLBasicsExamples(db); err != nil {
		return fmt.Errorf("failed to run SQL basics examples: %w", err)
	}

	// Run migration examples
	if err := examples.RunMigrationExamples(db); err != nil {
		return fmt.Errorf("failed to run migration examples: %w", err)
	}

	// Run transaction examples
	if err := examples.RunTransactionExamples(db); err != nil {
		return fmt.Errorf("failed to run transaction examples: %w", err)
	}

	log.Println("SQLite demonstration completed successfully")
	return nil
}

// demonstrateWithPostgreSQL demonstrates database operations with PostgreSQL
func demonstrateWithPostgreSQL() error {
	log.Println("\n--- PostgreSQL Demonstration ---")

	// Get PostgreSQL configuration
	config := database.GetDefaultPostgreSQLConfig()

	// Connect to PostgreSQL
	drivers := database.NewDatabaseDrivers()
	if err := drivers.ConnectPostgreSQL(config); err != nil {
		return fmt.Errorf("failed to connect to PostgreSQL: %w", err)
	}
	defer drivers.CloseAllConnections()

	db := drivers.GetPostgreSQLDB()
	if db == nil {
		return fmt.Errorf("PostgreSQL connection not available")
	}

	// Test PostgreSQL connection
	if err := drivers.TestPostgreSQLConnection(); err != nil {
		return fmt.Errorf("failed to test PostgreSQL connection: %w", err)
	}

	// Run connection pool examples
	examples := database.NewDatabaseExamples()
	if err := examples.RunConnectionPoolExamples(db); err != nil {
		return fmt.Errorf("failed to run connection pool examples: %w", err)
	}

	log.Println("PostgreSQL demonstration completed successfully")
	return nil
}

// demonstrateDatabaseDrivers demonstrates different database drivers
func demonstrateDatabaseDrivers() {
	log.Println("\n--- Database Drivers Demonstration ---")

	// Demonstrate PostgreSQL configuration
	postgresConfig := database.GetDefaultPostgreSQLConfig()
	log.Printf("PostgreSQL Config: %+v", postgresConfig)

	// Demonstrate MySQL configuration
	mysqlConfig := database.GetDefaultMySQLConfig()
	log.Printf("MySQL Config: %+v", mysqlConfig)

	// Demonstrate connection pool configurations
	defaultPoolConfig := database.GetDefaultPoolConfig()
	log.Printf("Default Pool Config: %+v", defaultPoolConfig)

	highPerfPoolConfig := database.GetHighPerformancePoolConfig()
	log.Printf("High Performance Pool Config: %+v", highPerfPoolConfig)

	lowResourcePoolConfig := database.GetLowResourcePoolConfig()
	log.Printf("Low Resource Pool Config: %+v", lowResourcePoolConfig)

	log.Println("Database drivers demonstration completed")
}

// demonstrateORMOperations demonstrates ORM operations
func demonstrateORMOperations() {
	log.Println("\n--- ORM Operations Demonstration ---")

	// Note: This would require GORM to be installed
	// go get gorm.io/gorm
	// go get gorm.io/driver/postgres
	// go get gorm.io/driver/mysql
	// go get gorm.io/driver/sqlite

	log.Println("To demonstrate ORM operations:")
	log.Println("1. Install GORM: go get gorm.io/gorm")
	log.Println("2. Install GORM drivers: go get gorm.io/driver/postgres")
	log.Println("3. Set up database connection")
	log.Println("4. Run ORM examples")

	log.Println("ORM operations demonstration completed")
}

// printUsage prints usage information
func printUsage() {
	log.Println("\n=== Database Operations Usage ===")
	log.Println("This demo showcases various database operations in Go:")
	log.Println("")
	log.Println("1. SQL Basics:")
	log.Println("   - Query execution")
	log.Println("   - CRUD operations")
	log.Println("   - Complex queries with WHERE, ORDER BY, etc.")
	log.Println("")
	log.Println("2. Database Drivers:")
	log.Println("   - PostgreSQL driver")
	log.Println("   - MySQL driver")
	log.Println("   - SQLite driver")
	log.Println("   - Connection configuration")
	log.Println("")
	log.Println("3. ORM Basics (GORM):")
	log.Println("   - Model definitions")
	log.Println("   - CRUD operations")
	log.Println("   - Associations (one-to-one, one-to-many)")
	log.Println("   - Migrations")
	log.Println("   - Soft deletes")
	log.Println("")
	log.Println("4. Connection Pooling:")
	log.Println("   - Connection pool management")
	log.Println("   - Health monitoring")
	log.Println("   - Performance benchmarking")
	log.Println("   - Timeout handling")
	log.Println("")
	log.Println("5. Migrations:")
	log.Println("   - Schema management")
	log.Println("   - Version control")
	log.Println("   - Up/Down migrations")
	log.Println("   - Rollback operations")
	log.Println("")
	log.Println("6. Transactions:")
	log.Println("   - ACID operations")
	log.Println("   - Transaction isolation")
	log.Println("   - Rollback handling")
	log.Println("   - Nested transactions")
	log.Println("")
	log.Println("To run with real databases:")
	log.Println("1. Install required drivers")
	log.Println("2. Set up database connections")
	log.Println("3. Update connection strings")
	log.Println("4. Run the examples")
}
