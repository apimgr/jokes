package config

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

// AdminConfig holds admin panel credentials
type AdminConfig struct {
	Email    string `yaml:"email"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Token    string `yaml:"token"`
}

// RateLimitConfig holds rate limiting settings
type RateLimitConfig struct {
	Enabled  bool `yaml:"enabled"`
	Requests int  `yaml:"requests"`
	Window   int  `yaml:"window"`
}

// LogsConfig holds logging configuration
type LogsConfig struct {
	Level  string `yaml:"level"`
	Access string `yaml:"access"`
	Server string `yaml:"server"`
}

// LetsEncryptConfig holds Let's Encrypt settings
type LetsEncryptConfig struct {
	Enabled         bool   `yaml:"enabled"`
	Email           string `yaml:"email"`
	Challenge       string `yaml:"challenge"`
	DNSProviderType string `yaml:"dns_provider_type"`
	DNSProviderKey  string `yaml:"dns_provider_key"`
	RFC2136Server   string `yaml:"rfc2136_server"`
	RFC2136Name     string `yaml:"rfc2136_name"`
	RFC2136Algo     string `yaml:"rfc2136_algorithm"`
}

// SSLConfig holds SSL/TLS settings
type SSLConfig struct {
	Enabled     bool              `yaml:"enabled"`
	CertPath    string            `yaml:"cert_path"`
	LetsEncrypt LetsEncryptConfig `yaml:"letsencrypt"`
}

// ScheduleConfig holds scheduler settings
type ScheduleConfig struct {
	Enabled       bool   `yaml:"enabled"`
	CertRenewal   string `yaml:"cert_renewal"`
	Notifications string `yaml:"notifications"`
	Cleanup       string `yaml:"cleanup"`
}

// MetricsConfig holds metrics/Prometheus settings
type MetricsConfig struct {
	Enabled       bool   `yaml:"enabled"`
	Endpoint      string `yaml:"endpoint"`
	IncludeSystem bool   `yaml:"include_system"`
	Token         string `yaml:"token"`
}

// NotificationsConfig holds notification settings
type NotificationsConfig struct {
	Enabled bool `yaml:"enabled"`
	Email   bool `yaml:"email"`
	Bell    bool `yaml:"bell"`
}

// ServerConfig holds all server settings
type ServerConfig struct {
	Port          int                 `yaml:"port"`
	Address       string              `yaml:"address"`
	FQDN          string              `yaml:"fqdn"`
	Title         string              `yaml:"title"`
	Description   string              `yaml:"description"`
	PIDFile       bool                `yaml:"pidfile"`
	Admin         AdminConfig         `yaml:"admin"`
	RateLimit     RateLimitConfig     `yaml:"rate_limit"`
	Logs          LogsConfig          `yaml:"logs"`
	SSL           SSLConfig           `yaml:"ssl"`
	Schedule      ScheduleConfig      `yaml:"schedule"`
	Metrics       MetricsConfig       `yaml:"metrics"`
	Notifications NotificationsConfig `yaml:"notifications"`
}

// FooterConfig holds footer customization
type FooterConfig struct {
	TrackingID string `yaml:"tracking_id"`
	CustomHTML string `yaml:"custom_html"`
}

// WebConfig holds web frontend settings
type WebConfig struct {
	Theme   string       `yaml:"theme"`
	Logo    string       `yaml:"logo"`
	Favicon string       `yaml:"favicon"`
	CORS    string       `yaml:"cors"`
	Footer  FooterConfig `yaml:"footer"`
}

// Config is the main configuration structure
type Config struct {
	Server ServerConfig `yaml:"server"`
	Web    WebConfig    `yaml:"web"`
}

// LoadConfig loads configuration from YAML file
// Priority: {configdir}/apimgr/jokes/server.yml
// If root/escalated: /etc/apimgr/jokes/server.yml
// Otherwise: ~/.config/apimgr/jokes/server.yml
// Auto-migrates from .yaml to .yml if needed
func LoadConfig() (*Config, error) {
	configPath := getConfigPath()

	// Check for .yaml and migrate to .yml if needed
	migrateYamlToYml(configPath)

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
	applyDefaults(&cfg)

	return &cfg, nil
}

// migrateYamlToYml migrates config from .yaml to .yml extension
func migrateYamlToYml(ymlPath string) {
	yamlPath := ymlPath[:len(ymlPath)-4] + ".yaml" // Replace .yml with .yaml

	// Check if .yaml exists and .yml doesn't
	if _, err := os.Stat(yamlPath); err == nil {
		if _, err := os.Stat(ymlPath); os.IsNotExist(err) {
			// Migrate .yaml to .yml
			if err := os.Rename(yamlPath, ymlPath); err != nil {
				log.Printf("Warning: failed to migrate %s to %s: %v", yamlPath, ymlPath, err)
			} else {
				log.Printf("Migrated config from %s to %s", yamlPath, ymlPath)
			}
		}
	}
}

// applyDefaults applies default values for missing config fields
func applyDefaults(cfg *Config) {
	// Server defaults
	if cfg.Server.Address == "" {
		cfg.Server.Address = "[::]"
	}
	if cfg.Server.Port == 0 {
		cfg.Server.Port = findAvailablePort()
	}
	if cfg.Server.FQDN == "" {
		if hostname, err := os.Hostname(); err == nil {
			cfg.Server.FQDN = hostname
		} else {
			cfg.Server.FQDN = "localhost"
		}
	}

	// Admin defaults
	if cfg.Server.Admin.Username == "" {
		cfg.Server.Admin.Username = "administrator"
	}
	if cfg.Server.Admin.Email == "" {
		cfg.Server.Admin.Email = "admin@" + cfg.Server.FQDN
	}

	// Rate limit defaults
	if cfg.Server.RateLimit.Requests == 0 {
		cfg.Server.RateLimit.Requests = 120
	}
	if cfg.Server.RateLimit.Window == 0 {
		cfg.Server.RateLimit.Window = 60
	}

	// Logs defaults
	if cfg.Server.Logs.Level == "" {
		cfg.Server.Logs.Level = "info"
	}
	if cfg.Server.Logs.Access == "" {
		cfg.Server.Logs.Access = "access.log"
	}
	if cfg.Server.Logs.Server == "" {
		cfg.Server.Logs.Server = "server.log"
	}

	// SSL defaults
	if cfg.Server.SSL.CertPath == "" {
		if os.Geteuid() == 0 {
			cfg.Server.SSL.CertPath = "/etc/apimgr/jokes/ssl/certs"
		} else {
			home, _ := os.UserHomeDir()
			cfg.Server.SSL.CertPath = filepath.Join(home, ".config", "apimgr", "jokes", "ssl", "certs")
		}
	}
	if cfg.Server.SSL.LetsEncrypt.Email == "" {
		cfg.Server.SSL.LetsEncrypt.Email = cfg.Server.Admin.Email
	}
	if cfg.Server.SSL.LetsEncrypt.Challenge == "" {
		cfg.Server.SSL.LetsEncrypt.Challenge = "http-01"
	}

	// Schedule defaults
	if cfg.Server.Schedule.CertRenewal == "" {
		cfg.Server.Schedule.CertRenewal = "daily"
	}
	if cfg.Server.Schedule.Notifications == "" {
		cfg.Server.Schedule.Notifications = "hourly"
	}
	if cfg.Server.Schedule.Cleanup == "" {
		cfg.Server.Schedule.Cleanup = "weekly"
	}

	// Metrics defaults
	if cfg.Server.Metrics.Endpoint == "" {
		cfg.Server.Metrics.Endpoint = "/metrics"
	}

	// Web defaults
	if cfg.Web.Theme == "" {
		cfg.Web.Theme = "dark"
	}
	if cfg.Web.CORS == "" {
		cfg.Web.CORS = "*"
	}
}

func getConfigPath() string {
	// Check if running as root (UID 0) or escalated
	if os.Geteuid() == 0 {
		return "/etc/apimgr/jokes/server.yml"
	}

	// Use ~/.config for regular users
	homeDir, err := os.UserHomeDir()
	if err != nil {
		// Fallback to current directory
		return "server.yml"
	}

	return filepath.Join(homeDir, ".config", "apimgr", "jokes", "server.yml")
}

func getDefaultConfig() *Config {
	hostname, _ := os.Hostname()
	if hostname == "" {
		hostname = "localhost"
	}

	certPath := "/etc/apimgr/jokes/ssl/certs"
	if os.Geteuid() != 0 {
		home, _ := os.UserHomeDir()
		certPath = filepath.Join(home, ".config", "apimgr", "jokes", "ssl", "certs")
	}

	cfg := &Config{
		Server: ServerConfig{
			Port:    findAvailablePort(),
			Address: "[::]",
			FQDN:    hostname,
			PIDFile: true,
			Admin: AdminConfig{
				Email:    "admin@" + hostname,
				Username: "administrator",
				Password: "", // Auto-generated on first run
				Token:    "", // Auto-generated on first run
			},
			RateLimit: RateLimitConfig{
				Enabled:  true,
				Requests: 120,
				Window:   60,
			},
			Logs: LogsConfig{
				Level:  "info",
				Access: "access.log",
				Server: "server.log",
			},
			SSL: SSLConfig{
				Enabled:  false,
				CertPath: certPath,
				LetsEncrypt: LetsEncryptConfig{
					Enabled:   false,
					Email:     "admin@" + hostname,
					Challenge: "http-01",
				},
			},
			Schedule: ScheduleConfig{
				Enabled:       true,
				CertRenewal:   "daily",
				Notifications: "hourly",
				Cleanup:       "weekly",
			},
			Metrics: MetricsConfig{
				Enabled:       false,
				Endpoint:      "/metrics",
				IncludeSystem: true,
			},
			Notifications: NotificationsConfig{
				Enabled: true,
				Email:   true,
				Bell:    true,
			},
		},
		Web: WebConfig{
			Theme:   "dark",
			Logo:    "",
			Favicon: "",
			CORS:    "*",
		},
	}
	return cfg
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
