package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"
)

// TransactionManager manages database transactions
type TransactionManager struct {
	db *sql.DB
}

// NewTransactionManager creates a new transaction manager
func NewTransactionManager(db *sql.DB) *TransactionManager {
	return &TransactionManager{db: db}
}

// TransactionOptions holds transaction options
type TransactionOptions struct {
	IsolationLevel sql.IsolationLevel
	ReadOnly       bool
	Timeout        time.Duration
}

// GetDefaultTransactionOptions returns default transaction options
func GetDefaultTransactionOptions() TransactionOptions {
	return TransactionOptions{
		IsolationLevel: sql.LevelReadCommitted,
		ReadOnly:       false,
		Timeout:        30 * time.Second,
	}
}

// ExecuteTransaction executes a function within a transaction
func (tm *TransactionManager) ExecuteTransaction(fn func(*sql.Tx) error, opts TransactionOptions) error {
	ctx, cancel := context.WithTimeout(context.Background(), opts.Timeout)
	defer cancel()

	tx, err := tm.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: opts.IsolationLevel,
		ReadOnly:  opts.ReadOnly,
	})
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	err = fn(tx)
	return err
}

// TransferMoney demonstrates a money transfer transaction
func (tm *TransactionManager) TransferMoney(fromAccountID, toAccountID int, amount float64) error {
	opts := GetDefaultTransactionOptions()

	return tm.ExecuteTransaction(func(tx *sql.Tx) error {
		// Check if source account has sufficient balance
		var balance float64
		err := tx.QueryRow("SELECT balance FROM accounts WHERE id = $1", fromAccountID).Scan(&balance)
		if err != nil {
			return fmt.Errorf("failed to get source account balance: %w", err)
		}

		if balance < amount {
			return fmt.Errorf("insufficient balance: %.2f < %.2f", balance, amount)
		}

		// Deduct from source account
		_, err = tx.Exec("UPDATE accounts SET balance = balance - $1 WHERE id = $2", amount, fromAccountID)
		if err != nil {
			return fmt.Errorf("failed to deduct from source account: %w", err)
		}

		// Add to destination account
		_, err = tx.Exec("UPDATE accounts SET balance = balance + $1 WHERE id = $2", amount, toAccountID)
		if err != nil {
			return fmt.Errorf("failed to add to destination account: %w", err)
		}

		// Record the transaction
		_, err = tx.Exec(`
			INSERT INTO transactions (from_account_id, to_account_id, amount, type, created_at) 
			VALUES ($1, $2, $3, $4, $5)`,
			fromAccountID, toAccountID, amount, "transfer", time.Now())
		if err != nil {
			return fmt.Errorf("failed to record transaction: %w", err)
		}

		log.Printf("Transfer successful: $%.2f from account %d to account %d", amount, fromAccountID, toAccountID)
		return nil
	}, opts)
}

// CreateUserWithProfile demonstrates creating related records in a transaction
func (tm *TransactionManager) CreateUserWithProfile(name, email string, age int, bio, website, location string) (*User, *Profile, error) {
	var createdUser *User
	var createdProfile *Profile

	opts := GetDefaultTransactionOptions()

	err := tm.ExecuteTransaction(func(tx *sql.Tx) error {
		// Create user
		userQuery := `
			INSERT INTO users (name, email, age, created_at, updated_at) 
			VALUES ($1, $2, $3, $4, $5) 
			RETURNING id, name, email, age, created_at`

		var user User
		err := tx.QueryRow(userQuery, name, email, age, time.Now(), time.Now()).Scan(
			&user.ID, &user.Name, &user.Email, &user.Age, &user.CreatedAt)
		if err != nil {
			return fmt.Errorf("failed to create user: %w", err)
		}
		createdUser = &user

		// Create profile
		profileQuery := `
			INSERT INTO profiles (user_id, bio, website, location, created_at) 
			VALUES ($1, $2, $3, $4, $5) 
			RETURNING id, user_id, bio, website, location`

		var profile Profile
		err = tx.QueryRow(profileQuery, user.ID, bio, website, location, time.Now()).Scan(
			&profile.ID, &profile.UserID, &profile.Bio, &profile.Website, &profile.Location)
		if err != nil {
			return fmt.Errorf("failed to create profile: %w", err)
		}
		createdProfile = &profile

		log.Printf("User and profile created successfully: User ID %d, Profile ID %d", user.ID, profile.ID)
		return nil
	}, opts)

	return createdUser, createdProfile, err
}

// BatchInsertUsers demonstrates batch operations in a transaction
func (tm *TransactionManager) BatchInsertUsers(users []User) error {
	opts := GetDefaultTransactionOptions()

	return tm.ExecuteTransaction(func(tx *sql.Tx) error {
		stmt, err := tx.Prepare(`
			INSERT INTO users (name, email, age, created_at, updated_at) 
			VALUES ($1, $2, $3, $4, $5)`)
		if err != nil {
			return fmt.Errorf("failed to prepare statement: %w", err)
		}
		defer stmt.Close()

		for _, user := range users {
			_, err = stmt.Exec(user.Name, user.Email, user.Age, time.Now(), time.Now())
			if err != nil {
				return fmt.Errorf("failed to insert user %s: %w", user.Name, err)
			}
		}

		log.Printf("Batch inserted %d users successfully", len(users))
		return nil
	}, opts)
}

// UpdateUserWithPosts demonstrates updating related records
func (tm *TransactionManager) UpdateUserWithPosts(userID int, newName, newEmail string, postsToAdd []Post) error {
	opts := GetDefaultTransactionOptions()

	return tm.ExecuteTransaction(func(tx *sql.Tx) error {
		// Update user
		_, err := tx.Exec(`
			UPDATE users 
			SET name = $1, email = $2, updated_at = $3 
			WHERE id = $4`,
			newName, newEmail, time.Now(), userID)
		if err != nil {
			return fmt.Errorf("failed to update user: %w", err)
		}

		// Add new posts
		for _, post := range postsToAdd {
			_, err = tx.Exec(`
				INSERT INTO posts (user_id, title, content, published, created_at, updated_at) 
				VALUES ($1, $2, $3, $4, $5, $6)`,
				userID, post.Title, post.Content, post.Published, time.Now(), time.Now())
			if err != nil {
				return fmt.Errorf("failed to insert post %s: %w", post.Title, err)
			}
		}

		log.Printf("Updated user %d and added %d posts", userID, len(postsToAdd))
		return nil
	}, opts)
}

// DeleteUserCascade demonstrates cascade delete in a transaction
func (tm *TransactionManager) DeleteUserCascade(userID int) error {
	opts := GetDefaultTransactionOptions()

	return tm.ExecuteTransaction(func(tx *sql.Tx) error {
		// Delete posts first (due to foreign key constraints)
		_, err := tx.Exec("DELETE FROM posts WHERE user_id = $1", userID)
		if err != nil {
			return fmt.Errorf("failed to delete posts: %w", err)
		}

		// Delete profile
		_, err = tx.Exec("DELETE FROM profiles WHERE user_id = $1", userID)
		if err != nil {
			return fmt.Errorf("failed to delete profile: %w", err)
		}

		// Delete user
		result, err := tx.Exec("DELETE FROM users WHERE id = $1", userID)
		if err != nil {
			return fmt.Errorf("failed to delete user: %w", err)
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return fmt.Errorf("failed to get rows affected: %w", err)
		}

		if rowsAffected == 0 {
			return fmt.Errorf("user with id %d not found", userID)
		}

		log.Printf("User %d and all related records deleted successfully", userID)
		return nil
	}, opts)
}

// ReadOnlyTransaction demonstrates read-only transactions
func (tm *TransactionManager) ReadOnlyTransaction() ([]User, error) {
	opts := TransactionOptions{
		IsolationLevel: sql.LevelReadCommitted,
		ReadOnly:       true,
		Timeout:        10 * time.Second,
	}

	var users []User

	err := tm.ExecuteTransaction(func(tx *sql.Tx) error {
		rows, err := tx.Query("SELECT id, name, email, age, created_at FROM users ORDER BY id")
		if err != nil {
			return fmt.Errorf("failed to query users: %w", err)
		}
		defer rows.Close()

		for rows.Next() {
			var user User
			err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Age, &user.CreatedAt)
			if err != nil {
				return fmt.Errorf("failed to scan user: %w", err)
			}
			users = append(users, user)
		}

		return rows.Err()
	}, opts)

	return users, err
}

// NestedTransaction demonstrates nested transaction handling
func (tm *TransactionManager) NestedTransaction() error {
	opts := GetDefaultTransactionOptions()

	return tm.ExecuteTransaction(func(outerTx *sql.Tx) error {
		log.Println("Outer transaction started")

		// Create a user in outer transaction
		_, err := outerTx.Exec(`
			INSERT INTO users (name, email, age, created_at, updated_at) 
			VALUES ($1, $2, $3, $4, $5)`,
			"Outer User", "outer@example.com", 30, time.Now(), time.Now())
		if err != nil {
			return fmt.Errorf("failed to create outer user: %w", err)
		}

		// Simulate nested transaction (in real scenario, this would be a savepoint)
		log.Println("Simulating nested transaction...")

		// This would typically use savepoints in a real implementation
		// For demonstration, we'll use a separate transaction
		return tm.ExecuteTransaction(func(innerTx *sql.Tx) error {
			log.Println("Inner transaction started")

			// Create a user in inner transaction
			_, err := innerTx.Exec(`
				INSERT INTO users (name, email, age, created_at, updated_at) 
				VALUES ($1, $2, $3, $4, $5)`,
				"Inner User", "inner@example.com", 25, time.Now(), time.Now())
			if err != nil {
				return fmt.Errorf("failed to create inner user: %w", err)
			}

			log.Println("Inner transaction completed")
			return nil
		}, opts)
	}, opts)
}

// TransactionWithRetry demonstrates transaction with retry logic
func (tm *TransactionManager) TransactionWithRetry(fn func(*sql.Tx) error, maxRetries int, retryDelay time.Duration) error {
	var lastErr error

	for i := 0; i < maxRetries; i++ {
		opts := GetDefaultTransactionOptions()

		err := tm.ExecuteTransaction(fn, opts)
		if err == nil {
			log.Printf("Transaction succeeded on attempt %d", i+1)
			return nil
		}

		lastErr = err
		log.Printf("Transaction attempt %d failed: %v", i+1, err)

		if i < maxRetries-1 {
			time.Sleep(retryDelay)
		}
	}

	return fmt.Errorf("transaction failed after %d attempts: %w", maxRetries, lastErr)
}

// GetTransactionIsolationLevel returns the current transaction isolation level
func (tm *TransactionManager) GetTransactionIsolationLevel() (string, error) {
	var level string
	err := tm.db.QueryRow("SELECT current_setting('transaction_isolation')").Scan(&level)
	if err != nil {
		return "", fmt.Errorf("failed to get isolation level: %w", err)
	}

	return level, nil
}

// SetTransactionIsolationLevel sets the transaction isolation level
func (tm *TransactionManager) SetTransactionIsolationLevel(level string) error {
	_, err := tm.db.Exec(fmt.Sprintf("SET TRANSACTION ISOLATION LEVEL %s", level))
	if err != nil {
		return fmt.Errorf("failed to set isolation level: %w", err)
	}

	log.Printf("Transaction isolation level set to: %s", level)
	return nil
}

// TransactionBenchmark benchmarks transaction performance
type TransactionBenchmark struct {
	manager *TransactionManager
}

// NewTransactionBenchmark creates a new transaction benchmark
func NewTransactionBenchmark(manager *TransactionManager) *TransactionBenchmark {
	return &TransactionBenchmark{manager: manager}
}

// BenchmarkConcurrentTransactions benchmarks concurrent transaction execution
func (tb *TransactionBenchmark) BenchmarkConcurrentTransactions(numTransactions int) error {
	start := time.Now()

	// Create a channel to collect results
	results := make(chan error, numTransactions)

	// Launch concurrent transactions
	for i := 0; i < numTransactions; i++ {
		go func(transactionID int) {
			opts := GetDefaultTransactionOptions()
			err := tb.manager.ExecuteTransaction(func(tx *sql.Tx) error {
				// Simulate some work
				_, err := tx.Exec("SELECT 1")
				return err
			}, opts)
			results <- err
		}(i)
	}

	// Collect results
	var errors []error
	for i := 0; i < numTransactions; i++ {
		if err := <-results; err != nil {
			errors = append(errors, err)
		}
	}

	duration := time.Since(start)
	avgTime := duration / time.Duration(numTransactions)

	log.Printf("Transaction benchmark completed: %d transactions, avg time: %v", numTransactions, avgTime)

	if len(errors) > 0 {
		log.Printf("Errors occurred: %d", len(errors))
		return fmt.Errorf("benchmark completed with %d errors", len(errors))
	}

	return nil
}
