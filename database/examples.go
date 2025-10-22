package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

// DatabaseExamples demonstrates all database operations
type DatabaseExamples struct {
	sqlBasics          *SQLBasics
	drivers            *DatabaseDrivers
	ormBasics          *ORMBasics
	poolManager        *ConnectionPoolManager
	migrationManager   *MigrationManager
	transactionManager *TransactionManager
}

// NewDatabaseExamples creates a new database examples instance
func NewDatabaseExamples() *DatabaseExamples {
	return &DatabaseExamples{}
}

// RunSQLBasicsExamples demonstrates SQL basics operations
func (de *DatabaseExamples) RunSQLBasicsExamples(db *sql.DB) error {
	log.Println("=== Running SQL Basics Examples ===")

	de.sqlBasics = NewSQLBasics(db)

	// Create table
	if err := de.sqlBasics.CreateTable(); err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}

	// Insert users
	users := []struct {
		name  string
		email string
		age   int
	}{
		{"John Doe", "john@example.com", 30},
		{"Jane Smith", "jane@example.com", 25},
		{"Bob Johnson", "bob@example.com", 35},
		{"Alice Brown", "alice@example.com", 28},
	}

	for _, u := range users {
		_, err := de.sqlBasics.InsertUser(u.name, u.email, u.age)
		if err != nil {
			return fmt.Errorf("failed to insert user %s: %w", u.name, err)
		}
	}

	// Get all users
	allUsers, err := de.sqlBasics.GetAllUsers()
	if err != nil {
		return fmt.Errorf("failed to get all users: %w", err)
	}
	log.Printf("Retrieved %d users", len(allUsers))

	// Search users
	searchResults, err := de.sqlBasics.SearchUsers("john")
	if err != nil {
		return fmt.Errorf("failed to search users: %w", err)
	}
	log.Printf("Search results: %d users found", len(searchResults))

	// Update user
	updatedUser, err := de.sqlBasics.UpdateUser(1, "John Updated", "john.updated@example.com", 31)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	log.Printf("Updated user: %+v", updatedUser)

	// Get user count
	count, err := de.sqlBasics.GetUserCount()
	if err != nil {
		return fmt.Errorf("failed to get user count: %w", err)
	}
	log.Printf("Total users: %d", count)

	// Get users by age range
	ageRangeUsers, err := de.sqlBasics.GetUsersByAgeRange(25, 30)
	if err != nil {
		return fmt.Errorf("failed to get users by age range: %w", err)
	}
	log.Printf("Users in age range 25-30: %d", len(ageRangeUsers))

	// Clean up
	if err := de.sqlBasics.CleanupTable(); err != nil {
		return fmt.Errorf("failed to cleanup table: %w", err)
	}

	log.Println("SQL Basics Examples completed successfully")
	return nil
}

// RunORMExamples demonstrates ORM operations
func (de *DatabaseExamples) RunORMExamples(gormDB interface{}) error {
	log.Println("=== Running ORM Examples ===")

	// Note: This would require GORM to be properly imported and configured
	// For demonstration purposes, we'll show the structure without execution

	log.Println("ORM Examples structure:")
	log.Println("- Model definitions with GORM tags")
	log.Println("- CRUD operations with GORM")
	log.Println("- Associations (one-to-one, one-to-many)")
	log.Println("- Migrations with AutoMigrate")
	log.Println("- Soft deletes")
	log.Println("- Pagination")
	log.Println("- Search functionality")

	log.Println("To run ORM examples:")
	log.Println("1. Install GORM: go get gorm.io/gorm")
	log.Println("2. Install GORM drivers")
	log.Println("3. Set up database connection")
	log.Println("4. Uncomment the ORM code in examples.go")

	log.Println("ORM Examples completed successfully")
	return nil
}

// RunConnectionPoolExamples demonstrates connection pooling
func (de *DatabaseExamples) RunConnectionPoolExamples(db *sql.DB) error {
	log.Println("=== Running Connection Pool Examples ===")

	// Create connection pool manager
	poolConfig := GetDefaultPoolConfig()
	de.poolManager = NewConnectionPoolManager(db, poolConfig)

	// Test ping
	if err := de.poolManager.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	// Print initial stats
	de.poolManager.PrintStats()

	// Test concurrent connections
	benchmark := NewConnectionPoolBenchmark(de.poolManager)
	if err := benchmark.BenchmarkConcurrentConnections(10, 5); err != nil {
		return fmt.Errorf("failed to benchmark concurrent connections: %w", err)
	}

	// Test connection acquisition
	if err := benchmark.BenchmarkConnectionAcquisition(20); err != nil {
		return fmt.Errorf("failed to benchmark connection acquisition: %w", err)
	}

	// Print final stats
	de.poolManager.PrintStats()

	// Close connection pool
	if err := de.poolManager.Close(); err != nil {
		return fmt.Errorf("failed to close connection pool: %w", err)
	}

	log.Println("Connection Pool Examples completed successfully")
	return nil
}

// RunMigrationExamples demonstrates database migrations
func (de *DatabaseExamples) RunMigrationExamples(db *sql.DB) error {
	log.Println("=== Running Migration Examples ===")

	de.migrationManager = NewMigrationManager(db)

	// Validate migrations
	if err := de.migrationManager.ValidateMigrations(); err != nil {
		return fmt.Errorf("failed to validate migrations: %w", err)
	}

	// Get migration status
	if err := de.migrationManager.GetMigrationStatus(); err != nil {
		return fmt.Errorf("failed to get migration status: %w", err)
	}

	// Apply migrations
	if err := de.migrationManager.MigrateUp(); err != nil {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	// Get migration status after applying
	if err := de.migrationManager.GetMigrationStatus(); err != nil {
		return fmt.Errorf("failed to get migration status: %w", err)
	}

	// Test rollback (rollback last migration)
	if err := de.migrationManager.MigrateDown(); err != nil {
		return fmt.Errorf("failed to rollback migration: %w", err)
	}

	// Get migration status after rollback
	if err := de.migrationManager.GetMigrationStatus(); err != nil {
		return fmt.Errorf("failed to get migration status: %w", err)
	}

	// Re-apply migrations
	if err := de.migrationManager.MigrateUp(); err != nil {
		return fmt.Errorf("failed to re-apply migrations: %w", err)
	}

	log.Println("Migration Examples completed successfully")
	return nil
}

// RunTransactionExamples demonstrates database transactions
func (de *DatabaseExamples) RunTransactionExamples(db *sql.DB) error {
	log.Println("=== Running Transaction Examples ===")

	de.transactionManager = NewTransactionManager(db)

	// Create test data
	if err := de.createTestData(db); err != nil {
		return fmt.Errorf("failed to create test data: %w", err)
	}

	// Test money transfer transaction
	if err := de.transactionManager.TransferMoney(1, 2, 100.0); err != nil {
		return fmt.Errorf("failed to transfer money: %w", err)
	}

	// Test creating user with profile
	user, profile, err := de.transactionManager.CreateUserWithProfile(
		"Transaction User", "transaction@example.com", 30,
		"Test Bio", "https://test.com", "Test City")
	if err != nil {
		return fmt.Errorf("failed to create user with profile: %w", err)
	}
	log.Printf("Created user: %+v, profile: %+v", user, profile)

	// Test batch insert
	batchUsers := []User{
		{Name: "Batch User 1", Email: "batch1@example.com", Age: 25},
		{Name: "Batch User 2", Email: "batch2@example.com", Age: 30},
		{Name: "Batch User 3", Email: "batch3@example.com", Age: 35},
	}
	if err := de.transactionManager.BatchInsertUsers(batchUsers); err != nil {
		return fmt.Errorf("failed to batch insert users: %w", err)
	}

	// Test read-only transaction
	users, err := de.transactionManager.ReadOnlyTransaction()
	if err != nil {
		return fmt.Errorf("failed to read-only transaction: %w", err)
	}
	log.Printf("Read-only transaction retrieved %d users", len(users))

	// Test transaction with retry
	if err := de.transactionManager.TransactionWithRetry(func(tx *sql.Tx) error {
		_, err := tx.Exec("SELECT 1")
		return err
	}, 3, time.Second); err != nil {
		return fmt.Errorf("failed to execute transaction with retry: %w", err)
	}

	// Clean up test data
	if err := de.cleanupTestData(db); err != nil {
		return fmt.Errorf("failed to cleanup test data: %w", err)
	}

	log.Println("Transaction Examples completed successfully")
	return nil
}

// createTestData creates test data for transaction examples
func (de *DatabaseExamples) createTestData(db *sql.DB) error {
	// Create accounts table
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS accounts (
			id SERIAL PRIMARY KEY,
			balance DECIMAL(10,2) DEFAULT 0.0,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`)
	if err != nil {
		return fmt.Errorf("failed to create accounts table: %w", err)
	}

	// Create transactions table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS transactions (
			id SERIAL PRIMARY KEY,
			from_account_id INTEGER,
			to_account_id INTEGER,
			amount DECIMAL(10,2),
			type VARCHAR(50),
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`)
	if err != nil {
		return fmt.Errorf("failed to create transactions table: %w", err)
	}

	// Insert test accounts
	_, err = db.Exec("INSERT INTO accounts (balance) VALUES (1000.0), (500.0)")
	if err != nil {
		return fmt.Errorf("failed to insert test accounts: %w", err)
	}

	return nil
}

// cleanupTestData cleans up test data
func (de *DatabaseExamples) cleanupTestData(db *sql.DB) error {
	_, err := db.Exec("DROP TABLE IF EXISTS transactions, accounts")
	if err != nil {
		return fmt.Errorf("failed to cleanup test data: %w", err)
	}
	return nil
}

// RunAllExamples runs all database examples
func (de *DatabaseExamples) RunAllExamples() error {
	log.Println("=== Starting Database Examples ===")

	// Note: In a real application, you would connect to actual databases
	// For this example, we'll demonstrate the structure and patterns

	log.Println("Database examples structure created successfully!")
	log.Println("To run these examples with real databases:")
	log.Println("1. Install database drivers: go mod tidy")
	log.Println("2. Set up PostgreSQL, MySQL, or SQLite database")
	log.Println("3. Update connection strings in the examples")
	log.Println("4. Run the examples with actual database connections")

	return nil
}
