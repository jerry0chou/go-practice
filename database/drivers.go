package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql" // MySQL driver
	_ "github.com/lib/pq"              // PostgreSQL driver
	_ "github.com/mattn/go-sqlite3"    // SQLite driver
)

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
	MaxConns int
	MaxIdle  int
}

// DatabaseDrivers demonstrates different database drivers
type DatabaseDrivers struct {
	postgresDB *sql.DB
	mysqlDB    *sql.DB
	sqliteDB   *sql.DB
}

// NewDatabaseDrivers creates a new DatabaseDrivers instance
func NewDatabaseDrivers() *DatabaseDrivers {
	return &DatabaseDrivers{}
}

// ConnectPostgreSQL demonstrates PostgreSQL connection
func (d *DatabaseDrivers) ConnectPostgreSQL(config DatabaseConfig) error {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.User, config.Password, config.DBName, config.SSLMode)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("failed to open PostgreSQL connection: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(config.MaxConns)
	db.SetMaxIdleConns(config.MaxIdle)
	db.SetConnMaxLifetime(time.Hour)

	// Test connection
	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping PostgreSQL: %w", err)
	}

	d.postgresDB = db
	log.Println("PostgreSQL connection established successfully")
	return nil
}

// ConnectMySQL demonstrates MySQL connection
func (d *DatabaseDrivers) ConnectMySQL(config DatabaseConfig) error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.User, config.Password, config.Host, config.Port, config.DBName)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("failed to open MySQL connection: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(config.MaxConns)
	db.SetMaxIdleConns(config.MaxIdle)
	db.SetConnMaxLifetime(time.Hour)

	// Test connection
	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping MySQL: %w", err)
	}

	d.mysqlDB = db
	log.Println("MySQL connection established successfully")
	return nil
}

// ConnectSQLite demonstrates SQLite connection
func (d *DatabaseDrivers) ConnectSQLite(dbPath string) error {
	dsn := fmt.Sprintf("file:%s?cache=shared&mode=rwc", dbPath)

	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return fmt.Errorf("failed to open SQLite connection: %w", err)
	}

	// Configure connection pool for SQLite
	db.SetMaxOpenConns(1) // SQLite doesn't support concurrent writes
	db.SetMaxIdleConns(1)
	db.SetConnMaxLifetime(time.Hour)

	// Test connection
	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping SQLite: %w", err)
	}

	d.sqliteDB = db
	log.Println("SQLite connection established successfully")
	return nil
}

// GetPostgreSQLDB returns PostgreSQL database connection
func (d *DatabaseDrivers) GetPostgreSQLDB() *sql.DB {
	return d.postgresDB
}

// GetMySQLDB returns MySQL database connection
func (d *DatabaseDrivers) GetMySQLDB() *sql.DB {
	return d.mysqlDB
}

// GetSQLiteDB returns SQLite database connection
func (d *DatabaseDrivers) GetSQLiteDB() *sql.DB {
	return d.sqliteDB
}

// TestPostgreSQLConnection tests PostgreSQL specific features
func (d *DatabaseDrivers) TestPostgreSQLConnection() error {
	if d.postgresDB == nil {
		return fmt.Errorf("PostgreSQL connection not established")
	}

	// Test PostgreSQL specific features
	query := `SELECT version()`
	var version string
	err := d.postgresDB.QueryRow(query).Scan(&version)
	if err != nil {
		return fmt.Errorf("failed to get PostgreSQL version: %w", err)
	}

	log.Printf("PostgreSQL version: %s", version)
	return nil
}

// TestMySQLConnection tests MySQL specific features
func (d *DatabaseDrivers) TestMySQLConnection() error {
	if d.mysqlDB == nil {
		return fmt.Errorf("MySQL connection not established")
	}

	// Test MySQL specific features
	query := `SELECT VERSION()`
	var version string
	err := d.mysqlDB.QueryRow(query).Scan(&version)
	if err != nil {
		return fmt.Errorf("failed to get MySQL version: %w", err)
	}

	log.Printf("MySQL version: %s", version)
	return nil
}

// TestSQLiteConnection tests SQLite specific features
func (d *DatabaseDrivers) TestSQLiteConnection() error {
	if d.sqliteDB == nil {
		return fmt.Errorf("SQLite connection not established")
	}

	// Test SQLite specific features
	query := `SELECT sqlite_version()`
	var version string
	err := d.sqliteDB.QueryRow(query).Scan(&version)
	if err != nil {
		return fmt.Errorf("failed to get SQLite version: %w", err)
	}

	log.Printf("SQLite version: %s", version)
	return nil
}

// GetConnectionStats returns connection pool statistics
func (d *DatabaseDrivers) GetConnectionStats(db *sql.DB, dbType string) {
	stats := db.Stats()
	log.Printf("%s Connection Stats:", dbType)
	log.Printf("  Open Connections: %d", stats.OpenConnections)
	log.Printf("  In Use: %d", stats.InUse)
	log.Printf("  Idle: %d", stats.Idle)
	log.Printf("  Wait Count: %d", stats.WaitCount)
	log.Printf("  Wait Duration: %v", stats.WaitDuration)
	log.Printf("  Max Idle Closed: %d", stats.MaxIdleClosed)
	log.Printf("  Max Idle Time Closed: %d", stats.MaxIdleTimeClosed)
	log.Printf("  Max Lifetime Closed: %d", stats.MaxLifetimeClosed)
}

// CloseAllConnections closes all database connections
func (d *DatabaseDrivers) CloseAllConnections() error {
	var errors []error

	if d.postgresDB != nil {
		if err := d.postgresDB.Close(); err != nil {
			errors = append(errors, fmt.Errorf("failed to close PostgreSQL: %w", err))
		} else {
			log.Println("PostgreSQL connection closed")
		}
	}

	if d.mysqlDB != nil {
		if err := d.mysqlDB.Close(); err != nil {
			errors = append(errors, fmt.Errorf("failed to close MySQL: %w", err))
		} else {
			log.Println("MySQL connection closed")
		}
	}

	if d.sqliteDB != nil {
		if err := d.sqliteDB.Close(); err != nil {
			errors = append(errors, fmt.Errorf("failed to close SQLite: %w", err))
		} else {
			log.Println("SQLite connection closed")
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("errors closing connections: %v", errors)
	}

	return nil
}

// GetDefaultPostgreSQLConfig returns default PostgreSQL configuration
func GetDefaultPostgreSQLConfig() DatabaseConfig {
	return DatabaseConfig{
		Host:     "localhost",
		Port:     5432,
		User:     "postgres",
		Password: "password",
		DBName:   "testdb",
		SSLMode:  "disable",
		MaxConns: 10,
		MaxIdle:  5,
	}
}

// GetDefaultMySQLConfig returns default MySQL configuration
func GetDefaultMySQLConfig() DatabaseConfig {
	return DatabaseConfig{
		Host:     "localhost",
		Port:     3306,
		User:     "root",
		Password: "password",
		DBName:   "testdb",
		MaxConns: 10,
		MaxIdle:  5,
	}
}
