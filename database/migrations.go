package database

import (
	"database/sql"
	"fmt"
	"log"
	"sort"
	"time"
)

// Migration represents a database migration
type Migration struct {
	Version   int       `json:"version"`
	Name      string    `json:"name"`
	UpSQL     string    `json:"up_sql"`
	DownSQL   string    `json:"down_sql"`
	AppliedAt time.Time `json:"applied_at,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

// MigrationManager manages database migrations
type MigrationManager struct {
	db         *sql.DB
	migrations []Migration
}

// NewMigrationManager creates a new migration manager
func NewMigrationManager(db *sql.DB) *MigrationManager {
	mm := &MigrationManager{
		db:         db,
		migrations: make([]Migration, 0),
	}

	// Initialize migrations table
	mm.createMigrationsTable()

	// Register default migrations
	mm.registerDefaultMigrations()

	return mm
}

// createMigrationsTable creates the migrations tracking table
func (mm *MigrationManager) createMigrationsTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS schema_migrations (
		version INTEGER PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`

	_, err := mm.db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	log.Println("Migrations table created/verified")
	return nil
}

// registerDefaultMigrations registers default migrations
func (mm *MigrationManager) registerDefaultMigrations() {
	// Migration 1: Create users table
	mm.AddMigration(Migration{
		Version: 1,
		Name:    "create_users_table",
		UpSQL: `
			CREATE TABLE users (
				id SERIAL PRIMARY KEY,
				name VARCHAR(100) NOT NULL,
				email VARCHAR(100) UNIQUE NOT NULL,
				age INTEGER,
				created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
			)`,
		DownSQL:   `DROP TABLE IF EXISTS users`,
		CreatedAt: time.Now(),
	})

	// Migration 2: Create profiles table
	mm.AddMigration(Migration{
		Version: 2,
		Name:    "create_profiles_table",
		UpSQL: `
			CREATE TABLE profiles (
				id SERIAL PRIMARY KEY,
				user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
				bio TEXT,
				website VARCHAR(255),
				location VARCHAR(100),
				created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
			)`,
		DownSQL:   `DROP TABLE IF EXISTS profiles`,
		CreatedAt: time.Now(),
	})

	// Migration 3: Create posts table
	mm.AddMigration(Migration{
		Version: 3,
		Name:    "create_posts_table",
		UpSQL: `
			CREATE TABLE posts (
				id SERIAL PRIMARY KEY,
				user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
				title VARCHAR(200) NOT NULL,
				content TEXT,
				published BOOLEAN DEFAULT FALSE,
				created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
			)`,
		DownSQL:   `DROP TABLE IF EXISTS posts`,
		CreatedAt: time.Now(),
	})

	// Migration 4: Add indexes
	mm.AddMigration(Migration{
		Version: 4,
		Name:    "add_indexes",
		UpSQL: `
			CREATE INDEX idx_users_email ON users(email);
			CREATE INDEX idx_posts_user_id ON posts(user_id);
			CREATE INDEX idx_posts_published ON posts(published);
			CREATE INDEX idx_profiles_user_id ON profiles(user_id)`,
		DownSQL: `
			DROP INDEX IF EXISTS idx_users_email;
			DROP INDEX IF EXISTS idx_posts_user_id;
			DROP INDEX IF EXISTS idx_posts_published;
			DROP INDEX IF EXISTS idx_profiles_user_id`,
		CreatedAt: time.Now(),
	})

	// Migration 5: Add soft delete to users
	mm.AddMigration(Migration{
		Version: 5,
		Name:    "add_soft_delete_to_users",
		UpSQL: `
			ALTER TABLE users ADD COLUMN deleted_at TIMESTAMP;
			CREATE INDEX idx_users_deleted_at ON users(deleted_at)`,
		DownSQL: `
			DROP INDEX IF EXISTS idx_users_deleted_at;
			ALTER TABLE users DROP COLUMN IF EXISTS deleted_at`,
		CreatedAt: time.Now(),
	})
}

// AddMigration adds a migration to the manager
func (mm *MigrationManager) AddMigration(migration Migration) {
	mm.migrations = append(mm.migrations, migration)
	log.Printf("Migration added: %d - %s", migration.Version, migration.Name)
}

// GetAppliedMigrations returns list of applied migrations
func (mm *MigrationManager) GetAppliedMigrations() ([]Migration, error) {
	query := `SELECT version, name, applied_at, created_at FROM schema_migrations ORDER BY version`

	rows, err := mm.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query applied migrations: %w", err)
	}
	defer rows.Close()

	var migrations []Migration
	for rows.Next() {
		var migration Migration
		err := rows.Scan(&migration.Version, &migration.Name, &migration.AppliedAt, &migration.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan migration: %w", err)
		}
		migrations = append(migrations, migration)
	}

	return migrations, nil
}

// GetPendingMigrations returns list of pending migrations
func (mm *MigrationManager) GetPendingMigrations() ([]Migration, error) {
	applied, err := mm.GetAppliedMigrations()
	if err != nil {
		return nil, fmt.Errorf("failed to get applied migrations: %w", err)
	}

	appliedVersions := make(map[int]bool)
	for _, migration := range applied {
		appliedVersions[migration.Version] = true
	}

	var pending []Migration
	for _, migration := range mm.migrations {
		if !appliedVersions[migration.Version] {
			pending = append(pending, migration)
		}
	}

	// Sort by version
	sort.Slice(pending, func(i, j int) bool {
		return pending[i].Version < pending[j].Version
	})

	return pending, nil
}

// MigrateUp applies all pending migrations
func (mm *MigrationManager) MigrateUp() error {
	pending, err := mm.GetPendingMigrations()
	if err != nil {
		return fmt.Errorf("failed to get pending migrations: %w", err)
	}

	if len(pending) == 0 {
		log.Println("No pending migrations")
		return nil
	}

	log.Printf("Applying %d pending migrations", len(pending))

	for _, migration := range pending {
		if err := mm.applyMigration(migration); err != nil {
			return fmt.Errorf("failed to apply migration %d: %w", migration.Version, err)
		}
	}

	log.Println("All migrations applied successfully")
	return nil
}

// MigrateDown rolls back the last migration
func (mm *MigrationManager) MigrateDown() error {
	applied, err := mm.GetAppliedMigrations()
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	if len(applied) == 0 {
		log.Println("No migrations to rollback")
		return nil
	}

	// Get the last applied migration
	lastMigration := applied[len(applied)-1]

	// Find the migration definition
	var migrationDef Migration
	for _, m := range mm.migrations {
		if m.Version == lastMigration.Version {
			migrationDef = m
			break
		}
	}

	if migrationDef.Version == 0 {
		return fmt.Errorf("migration definition not found for version %d", lastMigration.Version)
	}

	log.Printf("Rolling back migration %d: %s", migrationDef.Version, migrationDef.Name)

	// Execute down migration
	if _, err := mm.db.Exec(migrationDef.DownSQL); err != nil {
		return fmt.Errorf("failed to execute down migration: %w", err)
	}

	// Remove from applied migrations
	if err := mm.removeAppliedMigration(lastMigration.Version); err != nil {
		return fmt.Errorf("failed to remove applied migration: %w", err)
	}

	log.Printf("Migration %d rolled back successfully", migrationDef.Version)
	return nil
}

// applyMigration applies a single migration
func (mm *MigrationManager) applyMigration(migration Migration) error {
	log.Printf("Applying migration %d: %s", migration.Version, migration.Name)

	// Start transaction
	tx, err := mm.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Execute migration SQL
	if _, err := tx.Exec(migration.UpSQL); err != nil {
		return fmt.Errorf("failed to execute migration SQL: %w", err)
	}

	// Record migration as applied
	insertQuery := `INSERT INTO schema_migrations (version, name, applied_at, created_at) VALUES ($1, $2, $3, $4)`
	_, err = tx.Exec(insertQuery, migration.Version, migration.Name, time.Now(), migration.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to record migration: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit migration: %w", err)
	}

	log.Printf("Migration %d applied successfully", migration.Version)
	return nil
}

// removeAppliedMigration removes a migration from applied migrations
func (mm *MigrationManager) removeAppliedMigration(version int) error {
	query := `DELETE FROM schema_migrations WHERE version = $1`
	_, err := mm.db.Exec(query, version)
	if err != nil {
		return fmt.Errorf("failed to remove applied migration: %w", err)
	}

	return nil
}

// GetMigrationStatus returns the current migration status
func (mm *MigrationManager) GetMigrationStatus() error {
	applied, err := mm.GetAppliedMigrations()
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	pending, err := mm.GetPendingMigrations()
	if err != nil {
		return fmt.Errorf("failed to get pending migrations: %w", err)
	}

	log.Println("=== Migration Status ===")
	log.Printf("Applied migrations: %d", len(applied))
	for _, migration := range applied {
		log.Printf("  ✓ %d - %s (applied at %v)", migration.Version, migration.Name, migration.AppliedAt)
	}

	log.Printf("Pending migrations: %d", len(pending))
	for _, migration := range pending {
		log.Printf("  ○ %d - %s", migration.Version, migration.Name)
	}
	log.Println("========================")

	return nil
}

// ResetMigrations removes all applied migrations (dangerous!)
func (mm *MigrationManager) ResetMigrations() error {
	log.Println("WARNING: Resetting all migrations - this will remove all migration records!")

	query := `DELETE FROM schema_migrations`
	_, err := mm.db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to reset migrations: %w", err)
	}

	log.Println("All migration records removed")
	return nil
}

// ValidateMigrations validates that all migrations are properly defined
func (mm *MigrationManager) ValidateMigrations() error {
	log.Println("Validating migrations...")

	// Check for duplicate versions
	versions := make(map[int]bool)
	for _, migration := range mm.migrations {
		if versions[migration.Version] {
			return fmt.Errorf("duplicate migration version: %d", migration.Version)
		}
		versions[migration.Version] = true
	}

	// Check for missing up/down SQL
	for _, migration := range mm.migrations {
		if migration.UpSQL == "" {
			return fmt.Errorf("migration %d (%s) missing UpSQL", migration.Version, migration.Name)
		}
		if migration.DownSQL == "" {
			return fmt.Errorf("migration %d (%s) missing DownSQL", migration.Version, migration.Name)
		}
	}

	log.Printf("Validation passed: %d migrations are valid", len(mm.migrations))
	return nil
}

// CreateCustomMigration creates a custom migration
func (mm *MigrationManager) CreateCustomMigration(version int, name, upSQL, downSQL string) {
	migration := Migration{
		Version:   version,
		Name:      name,
		UpSQL:     upSQL,
		DownSQL:   downSQL,
		CreatedAt: time.Now(),
	}

	mm.AddMigration(migration)
	log.Printf("Custom migration created: %d - %s", version, name)
}
