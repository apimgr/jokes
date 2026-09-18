package config

import (
	"net"
	"os"
	"path/filepath"
	"testing"
)

func TestGetConfigPath(t *testing.T) {
	if os.Geteuid() == 0 {
		want := filepath.Join("/etc", "apimgr", "jokes", "server.yml")
		if got := getConfigPath(); got != want {
			t.Errorf("getConfigPath() as root = %q, want %q", got, want)
		}
		if got := GetConfigPath(); got != want {
			t.Errorf("GetConfigPath() as root = %q, want %q", got, want)
		}
		return
	}

	home := t.TempDir()
	t.Setenv("HOME", home)

	want := filepath.Join(home, ".config", "apimgr", "jokes", "server.yml")
	if got := getConfigPath(); got != want {
		t.Errorf("getConfigPath() = %q, want %q", got, want)
	}
	if got := GetConfigPath(); got != want {
		t.Errorf("GetConfigPath() = %q, want %q", got, want)
	}
}

func TestApplyDefaults(t *testing.T) {
	t.Run("fills every empty field", func(t *testing.T) {
		cfg := &Config{}
		applyDefaults(cfg)

		if cfg.Server.Address != "[::]" {
			t.Errorf("Server.Address = %q, want [::]", cfg.Server.Address)
		}
		if cfg.Server.Port == 0 {
			t.Error("Server.Port should be assigned a non-zero default")
		}
		if cfg.Server.FQDN == "" {
			t.Error("Server.FQDN should default to hostname or localhost")
		}
		if cfg.Server.Mode != "production" {
			t.Errorf("Server.Mode = %q, want production", cfg.Server.Mode)
		}
		if cfg.Server.UpdateBranch != "stable" {
			t.Errorf("Server.UpdateBranch = %q, want stable", cfg.Server.UpdateBranch)
		}
		if cfg.Server.Admin.Username != "administrator" {
			t.Errorf("Server.Admin.Username = %q, want administrator", cfg.Server.Admin.Username)
		}
		if cfg.Server.Admin.Email != "admin@"+cfg.Server.FQDN {
			t.Errorf("Server.Admin.Email = %q, want admin@%s", cfg.Server.Admin.Email, cfg.Server.FQDN)
		}
		if cfg.Server.RateLimit.Requests != 120 {
			t.Errorf("Server.RateLimit.Requests = %d, want 120", cfg.Server.RateLimit.Requests)
		}
		if cfg.Server.RateLimit.Window != 60 {
			t.Errorf("Server.RateLimit.Window = %d, want 60", cfg.Server.RateLimit.Window)
		}
		if cfg.Server.Logs.Level != "info" {
			t.Errorf("Server.Logs.Level = %q, want info", cfg.Server.Logs.Level)
		}
		if cfg.Server.Logs.Access != "access.log" {
			t.Errorf("Server.Logs.Access = %q, want access.log", cfg.Server.Logs.Access)
		}
		if cfg.Server.Logs.Server != "server.log" {
			t.Errorf("Server.Logs.Server = %q, want server.log", cfg.Server.Logs.Server)
		}
		if cfg.Server.SSL.CertPath == "" {
			t.Error("Server.SSL.CertPath should not be empty")
		}
		if cfg.Server.SSL.LetsEncrypt.Email != cfg.Server.Admin.Email {
			t.Errorf("Server.SSL.LetsEncrypt.Email = %q, want %q", cfg.Server.SSL.LetsEncrypt.Email, cfg.Server.Admin.Email)
		}
		if cfg.Server.SSL.LetsEncrypt.Challenge != "http-01" {
			t.Errorf("Server.SSL.LetsEncrypt.Challenge = %q, want http-01", cfg.Server.SSL.LetsEncrypt.Challenge)
		}
		if cfg.Server.Schedule.CertRenewal != "daily" {
			t.Errorf("Server.Schedule.CertRenewal = %q, want daily", cfg.Server.Schedule.CertRenewal)
		}
		if cfg.Server.Schedule.Notifications != "hourly" {
			t.Errorf("Server.Schedule.Notifications = %q, want hourly", cfg.Server.Schedule.Notifications)
		}
		if cfg.Server.Schedule.Cleanup != "weekly" {
			t.Errorf("Server.Schedule.Cleanup = %q, want weekly", cfg.Server.Schedule.Cleanup)
		}
		if cfg.Server.Metrics.Endpoint != "/metrics" {
			t.Errorf("Server.Metrics.Endpoint = %q, want /metrics", cfg.Server.Metrics.Endpoint)
		}
		if cfg.Web.Theme != "dark" {
			t.Errorf("Web.Theme = %q, want dark", cfg.Web.Theme)
		}
		if cfg.Web.CORS != "*" {
			t.Errorf("Web.CORS = %q, want *", cfg.Web.CORS)
		}
	})

	t.Run("does not override explicitly set values", func(t *testing.T) {
		cfg := &Config{
			Server: ServerConfig{
				Address: "127.0.0.1",
				Port:    8080,
				FQDN:    "example.com",
				Mode:    "development",
			},
		}
		applyDefaults(cfg)

		if cfg.Server.Address != "127.0.0.1" {
			t.Errorf("Server.Address was overridden: got %q", cfg.Server.Address)
		}
		if cfg.Server.Port != 8080 {
			t.Errorf("Server.Port was overridden: got %d", cfg.Server.Port)
		}
		if cfg.Server.FQDN != "example.com" {
			t.Errorf("Server.FQDN was overridden: got %q", cfg.Server.FQDN)
		}
		if cfg.Server.Mode != "development" {
			t.Errorf("Server.Mode was overridden: got %q", cfg.Server.Mode)
		}
	})
}

func TestGetDefaultConfig(t *testing.T) {
	cfg := getDefaultConfig()

	if cfg.Server.Port == 0 {
		t.Error("default config Port should not be zero")
	}
	if cfg.Server.Address != "[::]" {
		t.Errorf("default config Address = %q, want [::]", cfg.Server.Address)
	}
	if !cfg.Server.PIDFile {
		t.Error("default config PIDFile should be true")
	}
	if cfg.Server.Admin.Username != "administrator" {
		t.Errorf("default config Admin.Username = %q, want administrator", cfg.Server.Admin.Username)
	}
	if !cfg.Server.RateLimit.Enabled {
		t.Error("default config RateLimit.Enabled should be true")
	}
	if cfg.Server.SSL.Enabled {
		t.Error("default config SSL.Enabled should be false")
	}
	if !cfg.Server.Schedule.Enabled {
		t.Error("default config Schedule.Enabled should be true")
	}
	if cfg.Server.Metrics.Enabled {
		t.Error("default config Metrics.Enabled should be false")
	}
	if !cfg.Server.Notifications.Enabled {
		t.Error("default config Notifications.Enabled should be true")
	}
	if cfg.Web.Theme != "dark" {
		t.Errorf("default config Web.Theme = %q, want dark", cfg.Web.Theme)
	}
}

// configTestDir returns the directory LoadConfig/getConfigPath will use for
// this process, and registers cleanup. When running as root (always true in
// the Docker CI container) that directory is the fixed system path, and the
// cleanup removes only entries this subtest itself created there.
func configTestDir(t *testing.T) string {
	t.Helper()
	if os.Geteuid() == 0 {
		dir := filepath.Dir(getConfigPath())
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("setup failed: %v", err)
		}
		entries, _ := os.ReadDir(dir)
		existing := make(map[string]bool, len(entries))
		for _, e := range entries {
			existing[e.Name()] = true
		}
		t.Cleanup(func() {
			entries, _ := os.ReadDir(dir)
			for _, e := range entries {
				if !existing[e.Name()] {
					os.RemoveAll(filepath.Join(dir, e.Name()))
				}
			}
		})
		return dir
	}

	home := t.TempDir()
	t.Setenv("HOME", home)
	return filepath.Join(home, ".config", "apimgr", "jokes")
}

func TestLoadConfig(t *testing.T) {
	t.Run("creates default config when file is missing", func(t *testing.T) {
		dir := configTestDir(t)
		path := filepath.Join(dir, "server.yml")
		os.Remove(path)

		cfg, err := LoadConfig()
		if err != nil {
			t.Fatalf("LoadConfig() error = %v", err)
		}
		if cfg == nil {
			t.Fatal("LoadConfig() returned nil config")
		}

		if _, err := os.Stat(path); err != nil {
			t.Errorf("expected default config to be written at %q: %v", path, err)
		}
	})

	t.Run("loads and applies defaults on top of existing partial config", func(t *testing.T) {
		dir := configTestDir(t)
		path := filepath.Join(dir, "server.yml")
		yamlContent := "server:\n  address: \"10.0.0.1\"\n  port: 9999\n"
		if err := os.WriteFile(path, []byte(yamlContent), 0644); err != nil {
			t.Fatalf("setup failed: %v", err)
		}

		cfg, err := LoadConfig()
		if err != nil {
			t.Fatalf("LoadConfig() error = %v", err)
		}
		if cfg.Server.Address != "10.0.0.1" {
			t.Errorf("Server.Address = %q, want 10.0.0.1", cfg.Server.Address)
		}
		if cfg.Server.Port != 9999 {
			t.Errorf("Server.Port = %d, want 9999", cfg.Server.Port)
		}
		if cfg.Server.Mode != "production" {
			t.Errorf("Server.Mode = %q, want default production", cfg.Server.Mode)
		}
	})

	t.Run("malformed YAML returns error", func(t *testing.T) {
		dir := configTestDir(t)
		path := filepath.Join(dir, "server.yml")
		if err := os.WriteFile(path, []byte(": not valid yaml : :"), 0644); err != nil {
			t.Fatalf("setup failed: %v", err)
		}

		if _, err := LoadConfig(); err == nil {
			t.Error("LoadConfig() with malformed YAML should return an error")
		}
	})

	t.Run("migrates .yaml file to .yml", func(t *testing.T) {
		dir := configTestDir(t)
		path := filepath.Join(dir, "server.yml")
		os.Remove(path)
		yamlPath := filepath.Join(dir, "server.yaml")
		if err := os.WriteFile(yamlPath, []byte("server:\n  port: 7777\n"), 0644); err != nil {
			t.Fatalf("setup failed: %v", err)
		}

		cfg, err := LoadConfig()
		if err != nil {
			t.Fatalf("LoadConfig() error = %v", err)
		}
		if cfg.Server.Port != 7777 {
			t.Errorf("Server.Port after migration = %d, want 7777", cfg.Server.Port)
		}

		ymlPath := filepath.Join(dir, "server.yml")
		if _, err := os.Stat(ymlPath); err != nil {
			t.Errorf("expected migrated file at %q: %v", ymlPath, err)
		}
		if _, err := os.Stat(yamlPath); !os.IsNotExist(err) {
			t.Errorf("expected .yaml file to be removed after migration, stat err = %v", err)
		}
	})
}

func TestIsPortAvailable(t *testing.T) {
	t.Run("free port reports available", func(t *testing.T) {
		ln, err := net.Listen("tcp", ":0")
		if err != nil {
			t.Fatalf("failed to find a free port: %v", err)
		}
		port := ln.Addr().(*net.TCPAddr).Port
		if err := ln.Close(); err != nil {
			t.Fatalf("failed to close listener: %v", err)
		}

		if !isPortAvailable(port) {
			t.Errorf("isPortAvailable(%d) = false, want true after closing listener", port)
		}
	})

	t.Run("occupied port reports unavailable", func(t *testing.T) {
		ln, err := net.Listen("tcp", ":0")
		if err != nil {
			t.Fatalf("failed to bind a port: %v", err)
		}
		defer ln.Close()
		port := ln.Addr().(*net.TCPAddr).Port

		if isPortAvailable(port) {
			t.Errorf("isPortAvailable(%d) = true, want false while listener is active", port)
		}
	})
}

func TestFindAvailablePort(t *testing.T) {
	port := findAvailablePort()
	if port < 1 || port > 65535 {
		t.Errorf("findAvailablePort() = %d, want a valid TCP port", port)
	}
}
