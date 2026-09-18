package paths

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGetDataDir(t *testing.T) {
	dir := GetDataDir()
	if dir == "" {
		t.Fatal("GetDataDir returned empty string")
	}
	if os.Geteuid() == 0 {
		want := filepath.Join("/var/lib", orgName, appName)
		if dir != want {
			t.Errorf("GetDataDir() as root = %q, want %q", dir, want)
		}
		return
	}
	if !strings.Contains(dir, filepath.Join(".local/share", orgName, appName)) {
		t.Errorf("GetDataDir() = %q, want to contain .local/share/%s/%s", dir, orgName, appName)
	}
}

func TestGetConfigDir(t *testing.T) {
	dir := GetConfigDir()
	if dir == "" {
		t.Fatal("GetConfigDir returned empty string")
	}
	if os.Geteuid() == 0 {
		want := filepath.Join("/etc", orgName, appName)
		if dir != want {
			t.Errorf("GetConfigDir() as root = %q, want %q", dir, want)
		}
		return
	}
	if !strings.Contains(dir, filepath.Join(".config", orgName, appName)) {
		t.Errorf("GetConfigDir() = %q, want to contain .config/%s/%s", dir, orgName, appName)
	}
}

func TestGetLogDir(t *testing.T) {
	dir := GetLogDir()
	if dir == "" {
		t.Fatal("GetLogDir returned empty string")
	}
	if os.Geteuid() == 0 {
		want := filepath.Join("/var/log", orgName, appName)
		if dir != want {
			t.Errorf("GetLogDir() as root = %q, want %q", dir, want)
		}
		return
	}
	if !strings.HasSuffix(dir, "logs") {
		t.Errorf("GetLogDir() = %q, want to end with 'logs'", dir)
	}
}

func TestGetPIDFilePath(t *testing.T) {
	path := GetPIDFilePath()
	if path == "" {
		t.Fatal("GetPIDFilePath returned empty string")
	}
	if !strings.HasSuffix(path, appName+".pid") {
		t.Errorf("GetPIDFilePath() = %q, want suffix %q", path, appName+".pid")
	}
}

func TestGetBackupDir(t *testing.T) {
	dir := GetBackupDir()
	if dir == "" {
		t.Fatal("GetBackupDir returned empty string")
	}
	if os.Geteuid() == 0 {
		want := filepath.Join("/mnt/Backups", orgName, appName)
		if dir != want {
			t.Errorf("GetBackupDir() as root = %q, want %q", dir, want)
		}
		return
	}
	if !strings.Contains(dir, filepath.Join(orgName, appName)) {
		t.Errorf("GetBackupDir() = %q, want to contain %s/%s", dir, orgName, appName)
	}
}

func TestEnsureDir(t *testing.T) {
	base := t.TempDir()

	t.Run("creates new directory", func(t *testing.T) {
		target := filepath.Join(base, "sub", "nested")
		if err := EnsureDir(target); err != nil {
			t.Fatalf("EnsureDir() error = %v", err)
		}
		info, err := os.Stat(target)
		if err != nil {
			t.Fatalf("directory was not created: %v", err)
		}
		if !info.IsDir() {
			t.Errorf("target exists but is not a directory")
		}
	})

	t.Run("idempotent on existing directory", func(t *testing.T) {
		target := filepath.Join(base, "already-exists")
		if err := os.MkdirAll(target, 0755); err != nil {
			t.Fatalf("setup failed: %v", err)
		}
		if err := EnsureDir(target); err != nil {
			t.Errorf("EnsureDir() on existing dir returned error: %v", err)
		}
	})
}

func TestWriteAndRemovePIDFile(t *testing.T) {
	if os.Geteuid() != 0 {
		home := t.TempDir()
		t.Setenv("HOME", home)
	}
	t.Cleanup(RemovePIDFile)

	if err := WritePIDFile(); err != nil {
		t.Fatalf("WritePIDFile() error = %v", err)
	}

	pidPath := GetPIDFilePath()
	data, err := os.ReadFile(pidPath)
	if err != nil {
		t.Fatalf("PID file was not written: %v", err)
	}
	if len(data) == 0 {
		t.Error("PID file is empty")
	}

	RemovePIDFile()
	if _, err := os.Stat(pidPath); !os.IsNotExist(err) {
		t.Errorf("PID file still exists after RemovePIDFile(), stat err = %v", err)
	}
}

func TestWritePIDFileDetectsRunningInstance(t *testing.T) {
	if os.Geteuid() != 0 {
		home := t.TempDir()
		t.Setenv("HOME", home)
	}

	if err := WritePIDFile(); err != nil {
		t.Fatalf("first WritePIDFile() error = %v", err)
	}
	defer RemovePIDFile()

	// NOTE: WritePIDFile's liveness check calls process.Signal(os.Signal(nil)),
	// which always errors regardless of whether the process is alive, so the
	// existing PID file is always treated as stale and silently overwritten.
	// This asserts the actual (buggy) behavior rather than the intended one;
	// see the test-writer report for the bug this documents.
	if err := WritePIDFile(); err != nil {
		t.Errorf("second WritePIDFile() = %v, want nil (liveness check is currently a no-op due to a source bug)", err)
	}
}
