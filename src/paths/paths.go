package paths

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

const (
	appName = "jokes"
	orgName = "apimgr"
)

// GetDataDir returns the data directory path
func GetDataDir() string {
	if os.Geteuid() == 0 {
		return filepath.Join("/var/lib", orgName, appName)
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".local/share", orgName, appName)
}

// GetConfigDir returns the config directory path
func GetConfigDir() string {
	if os.Geteuid() == 0 {
		return filepath.Join("/etc", orgName, appName)
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", orgName, appName)
}

// GetLogDir returns the log directory path
func GetLogDir() string {
	if os.Geteuid() == 0 {
		return filepath.Join("/var/log", orgName, appName)
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".local/share", orgName, appName, "logs")
}

// GetPIDFilePath returns the PID file path
func GetPIDFilePath() string {
	if os.Geteuid() == 0 {
		return filepath.Join("/var/run", orgName, appName+".pid")
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".local/share", orgName, appName, appName+".pid")
}

// EnsureDir creates a directory if it doesn't exist
func EnsureDir(path string) error {
	return os.MkdirAll(path, 0755)
}

// WritePIDFile writes the current process PID to file
func WritePIDFile() error {
	pidPath := GetPIDFilePath()

	// Ensure directory exists
	if err := EnsureDir(filepath.Dir(pidPath)); err != nil {
		return fmt.Errorf("failed to create PID directory: %w", err)
	}

	// Check for stale PID file
	if data, err := os.ReadFile(pidPath); err == nil {
		if oldPID, err := strconv.Atoi(string(data)); err == nil {
			// Check if process exists
			if process, err := os.FindProcess(oldPID); err == nil {
				// On Unix, FindProcess always succeeds, so we need to send signal 0
				if err := process.Signal(os.Signal(nil)); err == nil {
					return fmt.Errorf("another instance is running (PID %d)", oldPID)
				}
			}
		}
		// Stale PID file, remove it
		os.Remove(pidPath)
	}

	// Write new PID file
	pid := os.Getpid()
	if err := os.WriteFile(pidPath, []byte(strconv.Itoa(pid)), 0644); err != nil {
		return fmt.Errorf("failed to write PID file: %w", err)
	}

	return nil
}

// RemovePIDFile removes the PID file
func RemovePIDFile() {
	os.Remove(GetPIDFilePath())
}

// GetBackupDir returns the backup directory path
func GetBackupDir() string {
	if os.Geteuid() == 0 {
		return filepath.Join("/mnt/Backups", orgName, appName)
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".local/backups", orgName, appName)
}
