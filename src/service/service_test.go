package service

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestDetectServiceManager(t *testing.T) {
	got := DetectServiceManager()

	valid := map[ServiceType]bool{
		ServiceUnknown: true,
		ServiceSystemd: true,
		ServiceRunit:   true,
		ServiceLaunchd: true,
		ServiceWindows: true,
		ServiceBSDRC:   true,
	}
	if !valid[got] {
		t.Errorf("DetectServiceManager() = %v, not a recognized ServiceType", got)
	}

	switch runtime.GOOS {
	case "darwin":
		if got != ServiceLaunchd {
			t.Errorf("DetectServiceManager() on darwin = %v, want ServiceLaunchd", got)
		}
	case "windows":
		if got != ServiceWindows {
			t.Errorf("DetectServiceManager() on windows = %v, want ServiceWindows", got)
		}
	case "freebsd", "openbsd", "netbsd":
		if got != ServiceBSDRC {
			t.Errorf("DetectServiceManager() on %s = %v, want ServiceBSDRC", runtime.GOOS, got)
		}
	}
}

func TestGetBinaryPath(t *testing.T) {
	path := GetBinaryPath()
	if path == "" {
		t.Fatal("GetBinaryPath() returned empty string")
	}

	if runtime.GOOS == "windows" {
		if !strings.Contains(path, "jokes.exe") {
			t.Errorf("GetBinaryPath() = %q, want it to contain jokes.exe", path)
		}
		return
	}

	if !filepath.IsAbs(path) {
		t.Errorf("GetBinaryPath() = %q, want an absolute path", path)
	}
	if !strings.HasSuffix(path, "/jokes") {
		t.Errorf("GetBinaryPath() = %q, want it to end with /jokes", path)
	}
}

func TestCopyBinary(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "src-binary")
	dst := filepath.Join(dir, "nested", "dst-binary")
	content := []byte("fake binary contents")

	if err := os.WriteFile(src, content, 0755); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	if err := copyBinary(src, dst); err != nil {
		t.Fatalf("copyBinary() error = %v", err)
	}

	got, err := os.ReadFile(dst)
	if err != nil {
		t.Fatalf("ReadFile(dst) error = %v", err)
	}
	if string(got) != string(content) {
		t.Errorf("copied content = %q, want %q", got, content)
	}

	info, err := os.Stat(dst)
	if err != nil {
		t.Fatalf("Stat(dst) error = %v", err)
	}
	if info.Mode().Perm()&0100 == 0 {
		t.Errorf("copied file mode = %v, want executable bit set", info.Mode())
	}
}

func TestCopyBinaryMissingSource(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "does-not-exist")
	dst := filepath.Join(dir, "dst")

	if err := copyBinary(src, dst); err == nil {
		t.Fatal("copyBinary() error = nil, want error for missing source file")
	}
}

func TestInstallUninstallUnsupportedManager(t *testing.T) {
	if runtime.GOOS != "plan9" {
		t.Skip("Install/Uninstall only return the unsupported-manager error on an OS with no matching case (e.g. plan9); real Install/Uninstall on this OS would mutate system state")
	}

	if err := Install(); err == nil {
		t.Error("Install() error = nil, want error for unsupported service manager")
	}
	if err := Uninstall(); err == nil {
		t.Error("Uninstall() error = nil, want error for unsupported service manager")
	}
}
