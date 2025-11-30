package config

import (
	"fmt"
	"math/rand"
	"net"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

type ServerConfig struct {
	Host      string `yaml:"host"`
	Port      int    `yaml:"port"`
	RateLimit int    `yaml:"rate_limit"`
}

type Config struct {
	Server ServerConfig `yaml:"server"`
}

// LoadConfig loads configuration from YAML file
// Priority: {configdir}/apimgr/jokes/server.yaml
// If root/escalated: /etc/apimgr/jokes/server.yaml
// Otherwise: ~/.config/apimgr/jokes/server.yaml
// If config file doesn't exist, creates it with defaults
func LoadConfig() (*Config, error) {
	configPath := getConfigPath()

	// If config file doesn't exist, create it with defaults
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		cfg := getDefaultConfig()
		if err := createDefaultConfig(configPath, cfg); err != nil {
			// If we can't create config, just use defaults without error
			// (might not have write permissions)
			return cfg, nil
		}
		return cfg, nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Apply defaults for missing values
	if cfg.Server.Host == "" {
		cfg.Server.Host = "0.0.0.0"
	}
	if cfg.Server.Port == 0 {
		cfg.Server.Port = findAvailablePort()
	}
	if cfg.Server.RateLimit == 0 {
		cfg.Server.RateLimit = 2000
	}

	return &cfg, nil
}

func getConfigPath() string {
	// Check if running as root (UID 0) or escalated
	if os.Geteuid() == 0 {
		return "/etc/apimgr/jokes/server.yaml"
	}

	// Use ~/.config for regular users
	homeDir, err := os.UserHomeDir()
	if err != nil {
		// Fallback to current directory
		return "server.yaml"
	}

	return filepath.Join(homeDir, ".config", "apimgr", "jokes", "server.yaml")
}

func getDefaultConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Host:      "0.0.0.0",
			Port:      findAvailablePort(),
			RateLimit: 2000,
		},
	}
}

// GetConfigPath returns the path where config file should be located
func GetConfigPath() string {
	return getConfigPath()
}

// createDefaultConfig creates a default configuration file
func createDefaultConfig(configPath string, cfg *Config) error {
	// Create directory if it doesn't exist
	dir := filepath.Dir(configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Marshal config to YAML
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Write config file
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// findAvailablePort finds an available port in the 64xxx range
func findAvailablePort() int {
	rand.Seed(time.Now().UnixNano())

	// Try to find an available port in 64000-64999 range
	for i := 0; i < 100; i++ {
		port := 64000 + rand.Intn(1000)
		if isPortAvailable(port) {
			return port
		}
	}

	// Fallback to 3009 if no random port available
	return 3009
}

// isPortAvailable checks if a port is available
func isPortAvailable(port int) bool {
	address := fmt.Sprintf(":%d", port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return false
	}
	listener.Close()
	return true
}
