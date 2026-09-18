package mode

import (
	"errors"
	"testing"
)

func resetMode(t *testing.T) {
	t.Helper()
	orig := Get()
	t.Cleanup(func() {
		mu.Lock()
		currentMode = orig
		mu.Unlock()
	})
}

func TestParseMode(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    Mode
		wantErr bool
	}{
		{"development lowercase", "development", Development, false},
		{"dev shorthand", "dev", Development, false},
		{"production lowercase", "production", Production, false},
		{"prod shorthand", "prod", Production, false},
		{"uppercase normalized", "PRODUCTION", Production, false},
		{"mixed case with whitespace", "  Dev  ", Development, false},
		{"invalid value", "staging", "", true},
		{"empty string", "", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseMode(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ParseMode(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("ParseMode(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestSetAndGet(t *testing.T) {
	resetMode(t)

	if err := Set("development"); err != nil {
		t.Fatalf("Set(\"development\") error = %v", err)
	}
	if Get() != Development {
		t.Errorf("Get() = %q, want %q", Get(), Development)
	}

	if err := Set("invalid"); err == nil {
		t.Error("Set(\"invalid\") should return an error")
	}
	if Get() != Development {
		t.Errorf("Get() after failed Set() = %q, want unchanged %q", Get(), Development)
	}
}

func TestIsDevelopmentIsProduction(t *testing.T) {
	resetMode(t)

	if err := Set("production"); err != nil {
		t.Fatalf("Set() error = %v", err)
	}
	if !IsProduction() || IsDevelopment() {
		t.Errorf("after Set(production): IsProduction()=%v IsDevelopment()=%v", IsProduction(), IsDevelopment())
	}

	if err := Set("development"); err != nil {
		t.Fatalf("Set() error = %v", err)
	}
	if IsProduction() || !IsDevelopment() {
		t.Errorf("after Set(development): IsProduction()=%v IsDevelopment()=%v", IsProduction(), IsDevelopment())
	}
}

func TestInitialize(t *testing.T) {
	resetMode(t)

	t.Run("CLI flag takes priority", func(t *testing.T) {
		t.Setenv("MODE", "production")
		if err := Initialize("development"); err != nil {
			t.Fatalf("Initialize() error = %v", err)
		}
		if Get() != Development {
			t.Errorf("Get() = %q, want %q (CLI flag should win)", Get(), Development)
		}
	})

	t.Run("env var used when no CLI flag", func(t *testing.T) {
		t.Setenv("MODE", "development")
		if err := Initialize(""); err != nil {
			t.Fatalf("Initialize() error = %v", err)
		}
		if Get() != Development {
			t.Errorf("Get() = %q, want %q (env var should apply)", Get(), Development)
		}
	})

	t.Run("invalid CLI flag returns error", func(t *testing.T) {
		if err := Initialize("nonsense"); err == nil {
			t.Error("Initialize(\"nonsense\") should return an error")
		}
	})
}

func TestGetErrorDetail(t *testing.T) {
	resetMode(t)

	if got := GetErrorDetail(nil); got != "" {
		t.Errorf("GetErrorDetail(nil) = %q, want empty string", got)
	}

	testErr := errors.New("boom: internal db failure")

	if err := Set("development"); err != nil {
		t.Fatalf("Set() error = %v", err)
	}
	if got := GetErrorDetail(testErr); got != testErr.Error() {
		t.Errorf("GetErrorDetail() in dev = %q, want %q", got, testErr.Error())
	}

	if err := Set("production"); err != nil {
		t.Fatalf("Set() error = %v", err)
	}
	got := GetErrorDetail(testErr)
	if got == testErr.Error() {
		t.Error("GetErrorDetail() in production should not leak internal error details")
	}
	if got == "" {
		t.Error("GetErrorDetail() in production should return a generic message, not empty")
	}
}

func TestShouldShowDebugEndpoints(t *testing.T) {
	resetMode(t)

	if err := Set("development"); err != nil {
		t.Fatalf("Set() error = %v", err)
	}
	if !ShouldShowDebugEndpoints() {
		t.Error("ShouldShowDebugEndpoints() in development = false, want true")
	}

	if err := Set("production"); err != nil {
		t.Fatalf("Set() error = %v", err)
	}
	if ShouldShowDebugEndpoints() {
		t.Error("ShouldShowDebugEndpoints() in production = true, want false")
	}
}

func TestGetCacheHeaders(t *testing.T) {
	resetMode(t)

	if err := Set("development"); err != nil {
		t.Fatalf("Set() error = %v", err)
	}
	devHeaders := GetCacheHeaders()
	if devHeaders.CacheControl != "no-cache, no-store, must-revalidate" {
		t.Errorf("dev CacheControl = %q, unexpected", devHeaders.CacheControl)
	}

	if err := Set("production"); err != nil {
		t.Fatalf("Set() error = %v", err)
	}
	prodHeaders := GetCacheHeaders()
	if prodHeaders.CacheControl != "public, max-age=31536000, immutable" {
		t.Errorf("prod CacheControl = %q, unexpected", prodHeaders.CacheControl)
	}
	if prodHeaders.Pragma != "" || prodHeaders.Expires != "" {
		t.Errorf("prod headers should not set Pragma/Expires, got %+v", prodHeaders)
	}
}

func TestGetLogLevel(t *testing.T) {
	resetMode(t)

	if err := Set("development"); err != nil {
		t.Fatalf("Set() error = %v", err)
	}
	if GetLogLevel() != "debug" {
		t.Errorf("GetLogLevel() in dev = %q, want debug", GetLogLevel())
	}

	if err := Set("production"); err != nil {
		t.Fatalf("Set() error = %v", err)
	}
	if GetLogLevel() != "info" {
		t.Errorf("GetLogLevel() in prod = %q, want info", GetLogLevel())
	}
}

func TestModeBehaviorFlags(t *testing.T) {
	resetMode(t)

	if err := Set("production"); err != nil {
		t.Fatalf("Set() error = %v", err)
	}
	if !ShouldCacheTemplates() {
		t.Error("ShouldCacheTemplates() in production = false, want true")
	}
	if ShouldEnableAutoReload() {
		t.Error("ShouldEnableAutoReload() in production = true, want false")
	}
	if ShouldEnableProfiling() {
		t.Error("ShouldEnableProfiling() in production = true, want false")
	}
	if GetPanicRecoveryMode() != "graceful" {
		t.Errorf("GetPanicRecoveryMode() in production = %q, want graceful", GetPanicRecoveryMode())
	}

	if err := Set("development"); err != nil {
		t.Fatalf("Set() error = %v", err)
	}
	if ShouldCacheTemplates() {
		t.Error("ShouldCacheTemplates() in development = true, want false")
	}
	if !ShouldEnableAutoReload() {
		t.Error("ShouldEnableAutoReload() in development = false, want true")
	}
	if !ShouldEnableProfiling() {
		t.Error("ShouldEnableProfiling() in development = false, want true")
	}
	if GetPanicRecoveryMode() != "verbose" {
		t.Errorf("GetPanicRecoveryMode() in development = %q, want verbose", GetPanicRecoveryMode())
	}
}

func TestModeStringAndValidate(t *testing.T) {
	if Production.String() != "production" {
		t.Errorf("Production.String() = %q, want production", Production.String())
	}
	if Development.String() != "development" {
		t.Errorf("Development.String() = %q, want development", Development.String())
	}

	if err := Production.Validate(); err != nil {
		t.Errorf("Production.Validate() error = %v", err)
	}
	if err := Development.Validate(); err != nil {
		t.Errorf("Development.Validate() error = %v", err)
	}
	if err := Mode("bogus").Validate(); err == nil {
		t.Error("Mode(\"bogus\").Validate() should return an error")
	}
}
