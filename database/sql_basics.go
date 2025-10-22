package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

// User represents a user in the database
type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Age       int       `json:"age"`
	CreatedAt time.Time `json:"created_at"`
}

// SQLBasics demonstrates basic SQL operations
type SQLBasics struct {
	db *sql.DB
}

// NewSQLBasics creates a new SQLBasics instance
func NewSQLBasics(db *sql.DB) *SQLBasics {
	return &SQLBasics{db: db}
}

// CreateTable creates the users table
func (s *SQLBasics) CreateTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		email VARCHAR(100) UNIQUE NOT NULL,
		age INTEGER,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`

	_, err := s.db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}

	log.Println("Table 'users' created successfully")
	return nil
}

// InsertUser demonstrates INSERT operation
func (s *SQLBasics) InsertUser(name, email string, age int) (*User, error) {
	query := `
	INSERT INTO users (name, email, age) 
	VALUES ($1, $2, $3) 
	RETURNING id, name, email, age, created_at`

	var user User
	err := s.db.QueryRow(query, name, email, age).Scan(
		&user.ID, &user.Name, &user.Email, &user.Age, &user.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to insert user: %w", err)
	}

	log.Printf("User inserted: %+v", user)
	return &user, nil
}

// GetUserByID demonstrates SELECT with WHERE clause
func (s *SQLBasics) GetUserByID(id int) (*User, error) {
	query := `SELECT id, name, email, age, created_at FROM users WHERE id = $1`

	var user User
	err := s.db.QueryRow(query, id).Scan(
		&user.ID, &user.Name, &user.Email, &user.Age, &user.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user with id %d not found", id)
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

// GetAllUsers demonstrates SELECT all records
func (s *SQLBasics) GetAllUsers() ([]User, error) {
	query := `SELECT id, name, email, age, created_at FROM users ORDER BY id`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query users: %w", err)
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Age, &user.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	log.Printf("Retrieved %d users", len(users))
	return users, nil
}

// UpdateUser demonstrates UPDATE operation
func (s *SQLBasics) UpdateUser(id int, name, email string, age int) (*User, error) {
	query := `
	UPDATE users 
	SET name = $1, email = $2, age = $3 
	WHERE id = $4 
	RETURNING id, name, email, age, created_at`

	var user User
	err := s.db.QueryRow(query, name, email, age, id).Scan(
		&user.ID, &user.Name, &user.Email, &user.Age, &user.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user with id %d not found", id)
		}
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	log.Printf("User updated: %+v", user)
	return &user, nil
}

// DeleteUser demonstrates DELETE operation
func (s *SQLBasics) DeleteUser(id int) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := s.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user with id %d not found", id)
	}

	log.Printf("User with id %d deleted successfully", id)
	return nil
}

// SearchUsers demonstrates complex SELECT with LIKE and ORDER BY
func (s *SQLBasics) SearchUsers(searchTerm string) ([]User, error) {
	query := `
	SELECT id, name, email, age, created_at 
	FROM users 
	WHERE name ILIKE $1 OR email ILIKE $1 
	ORDER BY created_at DESC`

	searchPattern := "%" + searchTerm + "%"
	rows, err := s.db.Query(query, searchPattern)
	if err != nil {
		return nil, fmt.Errorf("failed to search users: %w", err)
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Age, &user.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	log.Printf("Found %d users matching '%s'", len(users), searchTerm)
	return users, nil
}

// GetUserCount demonstrates COUNT operation
func (s *SQLBasics) GetUserCount() (int, error) {
	query := `SELECT COUNT(*) FROM users`

	var count int
	err := s.db.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get user count: %w", err)
	}

	log.Printf("Total users: %d", count)
	return count, nil
}

// GetUsersByAgeRange demonstrates range queries
func (s *SQLBasics) GetUsersByAgeRange(minAge, maxAge int) ([]User, error) {
	query := `
	SELECT id, name, email, age, created_at 
	FROM users 
	WHERE age BETWEEN $1 AND $2 
	ORDER BY age ASC`

	rows, err := s.db.Query(query, minAge, maxAge)
	if err != nil {
		return nil, fmt.Errorf("failed to query users by age range: %w", err)
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Age, &user.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	log.Printf("Found %d users between ages %d and %d", len(users), minAge, maxAge)
	return users, nil
}

// CleanupTable removes all records from the table
func (s *SQLBasics) CleanupTable() error {
	query := `DELETE FROM users`

	_, err := s.db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to cleanup table: %w", err)
	}

	log.Println("Table cleaned up successfully")
	return nil
}
