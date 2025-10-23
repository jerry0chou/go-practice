package config

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
)

// ConfigReloader handles hot reloading of configuration files
type ConfigReloader struct {
	configPath    string
	watcher       *fsnotify.Watcher
	reloadFunc    func() error
	reloadChannel chan struct{}
	stopChannel   chan struct{}
	mu            sync.RWMutex
	isRunning     bool
	reloadDelay   time.Duration
	lastReload    time.Time
}

// NewConfigReloader creates a new configuration reloader
func NewConfigReloader(configPath string, reloadFunc func() error) (*ConfigReloader, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("failed to create file watcher: %w", err)
	}

	// Watch the directory containing the config file
	configDir := filepath.Dir(configPath)
	if err := watcher.Add(configDir); err != nil {
		watcher.Close()
		return nil, fmt.Errorf("failed to watch directory %s: %w", configDir, err)
	}

	return &ConfigReloader{
		configPath:    configPath,
		watcher:       watcher,
		reloadFunc:    reloadFunc,
		reloadChannel: make(chan struct{}, 1),
		stopChannel:   make(chan struct{}),
		reloadDelay:   1 * time.Second, // Prevent rapid reloads
	}, nil
}

// Start starts the configuration reloader
func (cr *ConfigReloader) Start(ctx context.Context) error {
	cr.mu.Lock()
	defer cr.mu.Unlock()

	if cr.isRunning {
		return fmt.Errorf("config reloader is already running")
	}

	cr.isRunning = true

	go cr.watchLoop(ctx)
	go cr.reloadLoop(ctx)

	return nil
}

// Stop stops the configuration reloader
func (cr *ConfigReloader) Stop() error {
	cr.mu.Lock()
	defer cr.mu.Unlock()

	if !cr.isRunning {
		return fmt.Errorf("config reloader is not running")
	}

	close(cr.stopChannel)
	cr.watcher.Close()
	cr.isRunning = false

	return nil
}

// watchLoop monitors file system events
func (cr *ConfigReloader) watchLoop(ctx context.Context) {
	for {
		select {
		case event, ok := <-cr.watcher.Events:
			if !ok {
				return
			}

			// Check if the event is for our config file
			if filepath.Clean(event.Name) == filepath.Clean(cr.configPath) {
				if event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Create == fsnotify.Create {
					cr.triggerReload()
				}
			}

		case err, ok := <-cr.watcher.Errors:
			if !ok {
				return
			}
			fmt.Printf("File watcher error: %v\n", err)

		case <-ctx.Done():
			return

		case <-cr.stopChannel:
			return
		}
	}
}

// reloadLoop handles configuration reloading
func (cr *ConfigReloader) reloadLoop(ctx context.Context) {
	for {
		select {
		case <-cr.reloadChannel:
			// Debounce rapid reloads
			if time.Since(cr.lastReload) < cr.reloadDelay {
				continue
			}

			if err := cr.reloadConfig(); err != nil {
				fmt.Printf("Failed to reload configuration: %v\n", err)
			} else {
				fmt.Println("Configuration reloaded successfully")
				cr.lastReload = time.Now()
			}

		case <-ctx.Done():
			return

		case <-cr.stopChannel:
			return
		}
	}
}

// triggerReload triggers a configuration reload
func (cr *ConfigReloader) triggerReload() {
	select {
	case cr.reloadChannel <- struct{}{}:
	default:
		// Channel is full, skip this reload
	}
}

// reloadConfig reloads the configuration
func (cr *ConfigReloader) reloadConfig() error {
	// Check if the file still exists
	if _, err := os.Stat(cr.configPath); os.IsNotExist(err) {
		return fmt.Errorf("config file does not exist: %s", cr.configPath)
	}

	// Call the reload function
	return cr.reloadFunc()
}

// IsRunning returns whether the reloader is running
func (cr *ConfigReloader) IsRunning() bool {
	cr.mu.RLock()
	defer cr.mu.RUnlock()
	return cr.isRunning
}

// SetReloadDelay sets the delay between reloads to prevent rapid reloading
func (cr *ConfigReloader) SetReloadDelay(delay time.Duration) {
	cr.mu.Lock()
	defer cr.mu.Unlock()
	cr.reloadDelay = delay
}

// HotReloadManager manages hot reloading for multiple configuration types
type HotReloadManager struct {
	reloaders map[string]*ConfigReloader
	mu        sync.RWMutex
}

// NewHotReloadManager creates a new hot reload manager
func NewHotReloadManager() *HotReloadManager {
	return &HotReloadManager{
		reloaders: make(map[string]*ConfigReloader),
	}
}

// AddConfig adds a configuration file to hot reload
func (hrm *HotReloadManager) AddConfig(name, configPath string, reloadFunc func() error) error {
	hrm.mu.Lock()
	defer hrm.mu.Unlock()

	if _, exists := hrm.reloaders[name]; exists {
		return fmt.Errorf("configuration '%s' is already being watched", name)
	}

	reloader, err := NewConfigReloader(configPath, reloadFunc)
	if err != nil {
		return fmt.Errorf("failed to create reloader for '%s': %w", name, err)
	}

	hrm.reloaders[name] = reloader
	return nil
}

// StartAll starts all configuration reloaders
func (hrm *HotReloadManager) StartAll(ctx context.Context) error {
	hrm.mu.RLock()
	defer hrm.mu.RUnlock()

	for name, reloader := range hrm.reloaders {
		if err := reloader.Start(ctx); err != nil {
			return fmt.Errorf("failed to start reloader for '%s': %w", name, err)
		}
	}

	return nil
}

// StopAll stops all configuration reloaders
func (hrm *HotReloadManager) StopAll() error {
	hrm.mu.RLock()
	defer hrm.mu.RUnlock()

	var errors []string

	for name, reloader := range hrm.reloaders {
		if err := reloader.Stop(); err != nil {
			errors = append(errors, fmt.Sprintf("%s: %v", name, err))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("failed to stop some reloaders: %s", errors)
	}

	return nil
}

// StopConfig stops a specific configuration reloader
func (hrm *HotReloadManager) StopConfig(name string) error {
	hrm.mu.Lock()
	defer hrm.mu.Unlock()

	reloader, exists := hrm.reloaders[name]
	if !exists {
		return fmt.Errorf("configuration '%s' is not being watched", name)
	}

	if err := reloader.Stop(); err != nil {
		return fmt.Errorf("failed to stop reloader for '%s': %w", name, err)
	}

	delete(hrm.reloaders, name)
	return nil
}

// GetStatus returns the status of all reloaders
func (hrm *HotReloadManager) GetStatus() map[string]bool {
	hrm.mu.RLock()
	defer hrm.mu.RUnlock()

	status := make(map[string]bool)
	for name, reloader := range hrm.reloaders {
		status[name] = reloader.IsRunning()
	}

	return status
}

// ConfigReloadCallback defines a callback function for configuration reloads
type ConfigReloadCallback func(config interface{}) error

// ReloadableConfig represents a configuration that can be hot reloaded
type ReloadableConfig struct {
	config     interface{}
	loader     *ConfigLoader
	validator  *SchemaValidator
	callbacks  []ConfigReloadCallback
	mu         sync.RWMutex
	reloadTime time.Time
}

// NewReloadableConfig creates a new reloadable configuration
func NewReloadableConfig(configPath string, config interface{}, validator *SchemaValidator) (*ReloadableConfig, error) {
	loader := NewConfigLoader(configPath)

	return &ReloadableConfig{
		config:    config,
		loader:    loader,
		validator: validator,
		callbacks: make([]ConfigReloadCallback, 0),
	}, nil
}

// AddCallback adds a callback function to be called when configuration is reloaded
func (rc *ReloadableConfig) AddCallback(callback ConfigReloadCallback) {
	rc.mu.Lock()
	defer rc.mu.Unlock()
	rc.callbacks = append(rc.callbacks, callback)
}

// Reload reloads the configuration from file
func (rc *ReloadableConfig) Reload() error {
	rc.mu.Lock()
	defer rc.mu.Unlock()

	// Load new configuration
	newConfig, err := rc.loader.Load()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Validate new configuration
	if rc.validator != nil {
		if err := rc.validator.Validate(newConfig); err != nil {
			return fmt.Errorf("configuration validation failed: %w", err)
		}
	}

	// Update configuration
	rc.config = newConfig
	rc.reloadTime = time.Now()

	// Call all callbacks
	for _, callback := range rc.callbacks {
		if err := callback(rc.config); err != nil {
			fmt.Printf("Warning: callback failed during config reload: %v\n", err)
		}
	}

	return nil
}

// GetConfig returns the current configuration
func (rc *ReloadableConfig) GetConfig() interface{} {
	rc.mu.RLock()
	defer rc.mu.RUnlock()
	return rc.config
}

// GetReloadTime returns the time when configuration was last reloaded
func (rc *ReloadableConfig) GetReloadTime() time.Time {
	rc.mu.RLock()
	defer rc.mu.RUnlock()
	return rc.reloadTime
}

// CreateHotReloadExample demonstrates hot reloading functionality
func CreateHotReloadExample(configPath string) error {
	// Create a default configuration file
	if err := CreateDefaultConfig(configPath); err != nil {
		return fmt.Errorf("failed to create default config: %w", err)
	}

	// Create reloadable configuration
	validator := CreateDefaultSchema()
	reloadableConfig, err := NewReloadableConfig(configPath, &FileConfig{}, validator)
	if err != nil {
		return fmt.Errorf("failed to create reloadable config: %w", err)
	}

	// Add a callback to handle configuration changes
	reloadableConfig.AddCallback(func(config interface{}) error {
		fmt.Println("Configuration reloaded!")
		if fc, ok := config.(*FileConfig); ok {
			fc.PrintConfig()
		}
		return nil
	})

	// Initial load
	if err := reloadableConfig.Reload(); err != nil {
		return fmt.Errorf("failed to load initial configuration: %w", err)
	}

	// Create hot reload manager
	manager := NewHotReloadManager()

	// Add configuration to hot reload
	if err := manager.AddConfig("main", configPath, reloadableConfig.Reload); err != nil {
		return fmt.Errorf("failed to add config to hot reload: %w", err)
	}

	// Start hot reloading
	ctx := context.Background()
	if err := manager.StartAll(ctx); err != nil {
		return fmt.Errorf("failed to start hot reload: %w", err)
	}

	fmt.Println("Hot reload is active. Modify the configuration file to see changes.")
	fmt.Println("Press Ctrl+C to stop.")

	// Keep running until interrupted
	select {}
}
