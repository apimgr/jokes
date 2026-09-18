package ssl

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// writeSelfSignedCert writes a throwaway self-signed cert/key pair to the
// given paths for exercising certificate-loading code paths in tests.
func writeSelfSignedCert(t *testing.T, certPath, keyPath string) {
	t.Helper()

	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("rsa.GenerateKey() error = %v", err)
	}

	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject:      pkix.Name{CommonName: "example.com"},
		NotBefore:    time.Now().Add(-time.Hour),
		NotAfter:     time.Now().Add(time.Hour),
		KeyUsage:     x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
	}

	der, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	if err != nil {
		t.Fatalf("x509.CreateCertificate() error = %v", err)
	}

	certOut, err := os.Create(certPath)
	if err != nil {
		t.Fatalf("os.Create(cert) error = %v", err)
	}
	defer certOut.Close()
	if err := pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: der}); err != nil {
		t.Fatalf("pem.Encode(cert) error = %v", err)
	}

	keyOut, err := os.Create(keyPath)
	if err != nil {
		t.Fatalf("os.Create(key) error = %v", err)
	}
	defer keyOut.Close()
	if err := pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)}); err != nil {
		t.Fatalf("pem.Encode(key) error = %v", err)
	}
}

func TestFileExists(t *testing.T) {
	dir := t.TempDir()
	existing := filepath.Join(dir, "present.txt")
	if err := os.WriteFile(existing, []byte("x"), 0644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	if !fileExists(existing) {
		t.Error("fileExists() = false for existing file, want true")
	}
	if fileExists(filepath.Join(dir, "missing.txt")) {
		t.Error("fileExists() = true for missing file, want false")
	}
}

func TestFindManualCertsCrtKeyFormat(t *testing.T) {
	dir := t.TempDir()
	domain := "example.com"
	if err := os.WriteFile(filepath.Join(dir, domain+".crt"), []byte("cert"), 0644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, domain+".key"), []byte("key"), 0644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	m := NewManager(Config{CertPath: dir})
	cert, key := m.findManualCerts([]string{domain})

	wantCert := filepath.Join(dir, domain+".crt")
	wantKey := filepath.Join(dir, domain+".key")
	if cert != wantCert || key != wantKey {
		t.Errorf("findManualCerts() = (%q, %q), want (%q, %q)", cert, key, wantCert, wantKey)
	}
}

func TestFindManualCertsFullchainFormat(t *testing.T) {
	dir := t.TempDir()
	domain := "example.com"
	domainDir := filepath.Join(dir, domain)
	if err := os.MkdirAll(domainDir, 0755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}
	if err := os.WriteFile(filepath.Join(domainDir, "fullchain.pem"), []byte("cert"), 0644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	if err := os.WriteFile(filepath.Join(domainDir, "privkey.pem"), []byte("key"), 0644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	m := NewManager(Config{CertPath: dir})
	cert, key := m.findManualCerts([]string{domain})

	wantCert := filepath.Join(domainDir, "fullchain.pem")
	wantKey := filepath.Join(domainDir, "privkey.pem")
	if cert != wantCert || key != wantKey {
		t.Errorf("findManualCerts() = (%q, %q), want (%q, %q)", cert, key, wantCert, wantKey)
	}
}

func TestFindManualCertsNotFound(t *testing.T) {
	dir := t.TempDir()
	m := NewManager(Config{CertPath: dir})

	cert, key := m.findManualCerts([]string{"example.com"})
	if cert != "" || key != "" {
		t.Errorf("findManualCerts() = (%q, %q), want (\"\", \"\")", cert, key)
	}
}

func TestFindManualCertsEmptyCertPath(t *testing.T) {
	m := NewManager(Config{})
	cert, key := m.findManualCerts([]string{"example.com"})
	if cert != "" || key != "" {
		t.Errorf("findManualCerts() with empty CertPath = (%q, %q), want (\"\", \"\")", cert, key)
	}
}

func TestFindExistingCertsNotFound(t *testing.T) {
	m := NewManager(Config{})
	cert, key := m.findExistingCerts([]string{"nonexistent-domain-for-test.invalid"})
	if cert != "" || key != "" {
		t.Errorf("findExistingCerts() = (%q, %q), want (\"\", \"\")", cert, key)
	}
}

func TestGetTLSConfigDisabled(t *testing.T) {
	m := NewManager(Config{Enabled: false})
	cfg, err := m.GetTLSConfig([]string{"example.com"})
	if err != nil {
		t.Fatalf("GetTLSConfig() error = %v, want nil", err)
	}
	if cfg != nil {
		t.Errorf("GetTLSConfig() = %v, want nil when disabled", cfg)
	}
}

func TestGetTLSConfigNoCertsAvailable(t *testing.T) {
	dir := t.TempDir()
	m := NewManager(Config{Enabled: true, CertPath: dir})
	cfg, err := m.GetTLSConfig([]string{"example.com"})
	if err == nil {
		t.Fatal("GetTLSConfig() error = nil, want error when no certs and Let's Encrypt disabled")
	}
	if cfg != nil {
		t.Errorf("GetTLSConfig() = %v, want nil", cfg)
	}
}

func TestGetTLSConfigManualCerts(t *testing.T) {
	dir := t.TempDir()
	domain := "example.com"
	writeSelfSignedCert(t, filepath.Join(dir, domain+".crt"), filepath.Join(dir, domain+".key"))

	m := NewManager(Config{Enabled: true, CertPath: dir})
	cfg, err := m.GetTLSConfig([]string{domain})
	if err != nil {
		t.Fatalf("GetTLSConfig() error = %v, want nil", err)
	}
	if cfg == nil || len(cfg.Certificates) != 1 {
		t.Fatalf("GetTLSConfig() = %v, want config with one certificate", cfg)
	}
}

func TestGetTLSConfigLetsEncrypt(t *testing.T) {
	dir := t.TempDir()
	m := NewManager(Config{
		Enabled:  true,
		CertPath: dir,
		LetsEncrypt: LetsEncryptConfig{
			Enabled: true,
			Email:   "admin@example.com",
		},
	})

	cfg, err := m.GetTLSConfig([]string{"example.com"})
	if err != nil {
		t.Fatalf("GetTLSConfig() error = %v, want nil", err)
	}
	if cfg == nil {
		t.Fatal("GetTLSConfig() = nil, want a Let's Encrypt-backed config")
	}
	if _, err := os.Stat(filepath.Join(dir, "autocert")); err != nil {
		t.Errorf("autocert cache dir not created: %v", err)
	}
}

func TestGetHTTPHandlerFallback(t *testing.T) {
	m := NewManager(Config{})
	called := false
	fallback := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	})

	handler := m.GetHTTPHandler(fallback)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if !called {
		t.Error("GetHTTPHandler() did not delegate to fallback when certManager is nil")
	}
	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusOK)
	}
}

func TestChallengeServerSetGetClear(t *testing.T) {
	cs := NewChallengeServer()
	token := "abc123"
	auth := "abc123.keyauth"

	cs.SetToken(token, auth)

	req := httptest.NewRequest(http.MethodGet, "/.well-known/acme-challenge/"+token, nil)
	rec := httptest.NewRecorder()
	handled := cs.ServeHTTP(rec, req)

	if !handled {
		t.Fatal("ServeHTTP() = false for known token, want true")
	}
	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if got := rec.Body.String(); got != auth {
		t.Errorf("body = %q, want %q", got, auth)
	}

	cs.ClearToken(token)
	req2 := httptest.NewRequest(http.MethodGet, "/.well-known/acme-challenge/"+token, nil)
	rec2 := httptest.NewRecorder()
	handled2 := cs.ServeHTTP(rec2, req2)

	if !handled2 {
		t.Fatal("ServeHTTP() = false after clear for challenge path, want true (still handled as 404)")
	}
	if rec2.Code != http.StatusNotFound {
		t.Errorf("status after clear = %d, want %d", rec2.Code, http.StatusNotFound)
	}
}

func TestChallengeServerUnknownToken(t *testing.T) {
	cs := NewChallengeServer()
	req := httptest.NewRequest(http.MethodGet, "/.well-known/acme-challenge/unknown", nil)
	rec := httptest.NewRecorder()

	handled := cs.ServeHTTP(rec, req)
	if !handled {
		t.Fatal("ServeHTTP() = false for challenge path with unknown token, want true")
	}
	if rec.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusNotFound)
	}
}

func TestChallengeServerNonChallengePath(t *testing.T) {
	cs := NewChallengeServer()
	req := httptest.NewRequest(http.MethodGet, "/some/other/path", nil)
	rec := httptest.NewRecorder()

	handled := cs.ServeHTTP(rec, req)
	if handled {
		t.Error("ServeHTTP() = true for non-challenge path, want false")
	}
}

func TestParseChallenge(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"http-01 exact", "http-01", "http-01"},
		{"http01 no dash", "http01", "http-01"},
		{"http short", "http", "http-01"},
		{"tls-alpn-01 exact", "tls-alpn-01", "tls-alpn-01"},
		{"tlsalpn01 no dash", "tlsalpn01", "tls-alpn-01"},
		{"tls-alpn short", "tls-alpn", "tls-alpn-01"},
		{"tls short", "tls", "tls-alpn-01"},
		{"dns-01 exact", "dns-01", "dns-01"},
		{"dns01 no dash", "dns01", "dns-01"},
		{"dns short", "dns", "dns-01"},
		{"uppercase", "HTTP-01", "http-01"},
		{"whitespace", "  http-01  ", "http-01"},
		{"empty defaults to http-01", "", "http-01"},
		{"unknown defaults to http-01", "carrier-pigeon-01", "http-01"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseChallenge(tt.input); got != tt.want {
				t.Errorf("ParseChallenge(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
