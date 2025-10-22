package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"
)

// PoolConfig holds connection pool configuration
type PoolConfig struct {
	MaxOpenConns     int           `json:"max_open_conns"`
	MaxIdleConns     int           `json:"max_idle_conns"`
	ConnMaxLifetime  time.Duration `json:"conn_max_lifetime"`
	ConnMaxIdleTime  time.Duration `json:"conn_max_idle_time"`
	HealthCheckDelay time.Duration `json:"health_check_delay"`
}

// ConnectionPoolManager manages database connection pools
type ConnectionPoolManager struct {
	db     *sql.DB
	config PoolConfig
	ctx    context.Context
	cancel context.CancelFunc
	mu     sync.RWMutex
	stats  *PoolStats
}

// PoolStats holds connection pool statistics
type PoolStats struct {
	OpenConnections   int           `json:"open_connections"`
	InUse             int           `json:"in_use"`
	Idle              int           `json:"idle"`
	WaitCount         int64         `json:"wait_count"`
	WaitDuration      time.Duration `json:"wait_duration"`
	MaxIdleClosed     int64         `json:"max_idle_closed"`
	MaxIdleTimeClosed int64         `json:"max_idle_time_closed"`
	MaxLifetimeClosed int64         `json:"max_lifetime_closed"`
	LastUpdated       time.Time     `json:"last_updated"`
}

// NewConnectionPoolManager creates a new connection pool manager
func NewConnectionPoolManager(db *sql.DB, config PoolConfig) *ConnectionPoolManager {
	ctx, cancel := context.WithCancel(context.Background())

	manager := &ConnectionPoolManager{
		db:     db,
		config: config,
		ctx:    ctx,
		cancel: cancel,
		stats:  &PoolStats{},
	}

	// Configure the database connection pool
	manager.configurePool()

	// Start health monitoring
	go manager.startHealthMonitoring()

	return manager
}

// configurePool configures the database connection pool
func (cpm *ConnectionPoolManager) configurePool() {
	cpm.db.SetMaxOpenConns(cpm.config.MaxOpenConns)
	cpm.db.SetMaxIdleConns(cpm.config.MaxIdleConns)
	cpm.db.SetConnMaxLifetime(cpm.config.ConnMaxLifetime)
	cpm.db.SetConnMaxIdleTime(cpm.config.ConnMaxIdleTime)

	log.Printf("Connection pool configured: MaxOpen=%d, MaxIdle=%d, MaxLifetime=%v, MaxIdleTime=%v",
		cpm.config.MaxOpenConns, cpm.config.MaxIdleConns, cpm.config.ConnMaxLifetime, cpm.config.ConnMaxIdleTime)
}

// startHealthMonitoring starts monitoring the connection pool health
func (cpm *ConnectionPoolManager) startHealthMonitoring() {
	ticker := time.NewTicker(cpm.config.HealthCheckDelay)
	defer ticker.Stop()

	for {
		select {
		case <-cpm.ctx.Done():
			log.Println("Health monitoring stopped")
			return
		case <-ticker.C:
			cpm.updateStats()
			cpm.checkHealth()
		}
	}
}

// updateStats updates the connection pool statistics
func (cpm *ConnectionPoolManager) updateStats() {
	cpm.mu.Lock()
	defer cpm.mu.Unlock()

	stats := cpm.db.Stats()
	cpm.stats = &PoolStats{
		OpenConnections:   stats.OpenConnections,
		InUse:             stats.InUse,
		Idle:              stats.Idle,
		WaitCount:         stats.WaitCount,
		WaitDuration:      stats.WaitDuration,
		MaxIdleClosed:     stats.MaxIdleClosed,
		MaxIdleTimeClosed: stats.MaxIdleTimeClosed,
		MaxLifetimeClosed: stats.MaxLifetimeClosed,
		LastUpdated:       time.Now(),
	}
}

// checkHealth checks the health of the connection pool
func (cpm *ConnectionPoolManager) checkHealth() {
	cpm.mu.RLock()
	stats := *cpm.stats
	cpm.mu.RUnlock()

	// Check for potential issues
	if stats.WaitCount > 100 {
		log.Printf("WARNING: High wait count detected: %d", stats.WaitCount)
	}

	if stats.WaitDuration > 5*time.Second {
		log.Printf("WARNING: High wait duration detected: %v", stats.WaitDuration)
	}

	if stats.OpenConnections >= cpm.config.MaxOpenConns {
		log.Printf("WARNING: Connection pool at maximum capacity: %d/%d",
			stats.OpenConnections, cpm.config.MaxOpenConns)
	}
}

// GetStats returns current connection pool statistics
func (cpm *ConnectionPoolManager) GetStats() PoolStats {
	cpm.mu.RLock()
	defer cpm.mu.RUnlock()
	return *cpm.stats
}

// GetDB returns the database connection
func (cpm *ConnectionPoolManager) GetDB() *sql.DB {
	return cpm.db
}

// Ping tests the database connection
func (cpm *ConnectionPoolManager) Ping() error {
	ctx, cancel := context.WithTimeout(cpm.ctx, 5*time.Second)
	defer cancel()

	return cpm.db.PingContext(ctx)
}

// PingWithRetry tests the database connection with retry logic
func (cpm *ConnectionPoolManager) PingWithRetry(maxRetries int, retryDelay time.Duration) error {
	var lastErr error

	for i := 0; i < maxRetries; i++ {
		if err := cpm.Ping(); err != nil {
			lastErr = err
			log.Printf("Ping attempt %d/%d failed: %v", i+1, maxRetries, err)
			if i < maxRetries-1 {
				time.Sleep(retryDelay)
			}
		} else {
			log.Printf("Ping successful on attempt %d", i+1)
			return nil
		}
	}

	return fmt.Errorf("ping failed after %d attempts: %w", maxRetries, lastErr)
}

// ExecuteWithTimeout executes a query with timeout
func (cpm *ConnectionPoolManager) ExecuteWithTimeout(query string, timeout time.Duration, args ...interface{}) (sql.Result, error) {
	ctx, cancel := context.WithTimeout(cpm.ctx, timeout)
	defer cancel()

	return cpm.db.ExecContext(ctx, query, args...)
}

// QueryWithTimeout executes a query with timeout and returns rows
func (cpm *ConnectionPoolManager) QueryWithTimeout(query string, timeout time.Duration, args ...interface{}) (*sql.Rows, error) {
	ctx, cancel := context.WithTimeout(cpm.ctx, timeout)
	defer cancel()

	return cpm.db.QueryContext(ctx, query, args...)
}

// QueryRowWithTimeout executes a query with timeout and returns a single row
func (cpm *ConnectionPoolManager) QueryRowWithTimeout(query string, timeout time.Duration, args ...interface{}) *sql.Row {
	ctx, cancel := context.WithTimeout(cpm.ctx, timeout)
	defer cancel()

	return cpm.db.QueryRowContext(ctx, query, args...)
}

// GetConnectionWithTimeout gets a connection from the pool with timeout
func (cpm *ConnectionPoolManager) GetConnectionWithTimeout(timeout time.Duration) (*sql.Conn, error) {
	ctx, cancel := context.WithTimeout(cpm.ctx, timeout)
	defer cancel()

	return cpm.db.Conn(ctx)
}

// Close closes the connection pool manager
func (cpm *ConnectionPoolManager) Close() error {
	cpm.cancel()
	return cpm.db.Close()
}

// PrintStats prints current connection pool statistics
func (cpm *ConnectionPoolManager) PrintStats() {
	stats := cpm.GetStats()

	log.Println("=== Connection Pool Statistics ===")
	log.Printf("Open Connections: %d", stats.OpenConnections)
	log.Printf("In Use: %d", stats.InUse)
	log.Printf("Idle: %d", stats.Idle)
	log.Printf("Wait Count: %d", stats.WaitCount)
	log.Printf("Wait Duration: %v", stats.WaitDuration)
	log.Printf("Max Idle Closed: %d", stats.MaxIdleClosed)
	log.Printf("Max Idle Time Closed: %d", stats.MaxIdleTimeClosed)
	log.Printf("Max Lifetime Closed: %d", stats.MaxLifetimeClosed)
	log.Printf("Last Updated: %v", stats.LastUpdated)
	log.Println("==================================")
}

// GetDefaultPoolConfig returns default connection pool configuration
func GetDefaultPoolConfig() PoolConfig {
	return PoolConfig{
		MaxOpenConns:     25,
		MaxIdleConns:     10,
		ConnMaxLifetime:  time.Hour,
		ConnMaxIdleTime:  30 * time.Minute,
		HealthCheckDelay: 30 * time.Second,
	}
}

// GetHighPerformancePoolConfig returns high-performance connection pool configuration
func GetHighPerformancePoolConfig() PoolConfig {
	return PoolConfig{
		MaxOpenConns:     100,
		MaxIdleConns:     50,
		ConnMaxLifetime:  2 * time.Hour,
		ConnMaxIdleTime:  15 * time.Minute,
		HealthCheckDelay: 10 * time.Second,
	}
}

// GetLowResourcePoolConfig returns low-resource connection pool configuration
func GetLowResourcePoolConfig() PoolConfig {
	return PoolConfig{
		MaxOpenConns:     5,
		MaxIdleConns:     2,
		ConnMaxLifetime:  30 * time.Minute,
		ConnMaxIdleTime:  5 * time.Minute,
		HealthCheckDelay: 60 * time.Second,
	}
}

// ConnectionPoolBenchmark benchmarks connection pool performance
type ConnectionPoolBenchmark struct {
	manager *ConnectionPoolManager
}

// NewConnectionPoolBenchmark creates a new benchmark instance
func NewConnectionPoolBenchmark(manager *ConnectionPoolManager) *ConnectionPoolBenchmark {
	return &ConnectionPoolBenchmark{manager: manager}
}

// BenchmarkConcurrentConnections benchmarks concurrent connection usage
func (cpb *ConnectionPoolBenchmark) BenchmarkConcurrentConnections(numGoroutines, queriesPerGoroutine int) error {
	var wg sync.WaitGroup
	start := time.Now()

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(goroutineID int) {
			defer wg.Done()

			for j := 0; j < queriesPerGoroutine; j++ {
				query := "SELECT 1"
				_, err := cpb.manager.QueryWithTimeout(query, 5*time.Second)
				if err != nil {
					log.Printf("Goroutine %d, Query %d failed: %v", goroutineID, j, err)
					return
				}
			}
		}(i)
	}

	wg.Wait()
	duration := time.Since(start)

	log.Printf("Benchmark completed: %d goroutines, %d queries each, took %v",
		numGoroutines, queriesPerGoroutine, duration)

	cpb.manager.PrintStats()
	return nil
}

// BenchmarkConnectionAcquisition benchmarks connection acquisition time
func (cpb *ConnectionPoolBenchmark) BenchmarkConnectionAcquisition(numConnections int) error {
	start := time.Now()

	for i := 0; i < numConnections; i++ {
		conn, err := cpb.manager.GetConnectionWithTimeout(5 * time.Second)
		if err != nil {
			return fmt.Errorf("failed to get connection %d: %w", i, err)
		}
		conn.Close()
	}

	duration := time.Since(start)
	avgTime := duration / time.Duration(numConnections)

	log.Printf("Connection acquisition benchmark: %d connections, avg time: %v",
		numConnections, avgTime)

	return nil
}
