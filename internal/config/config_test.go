package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadDefault(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("GBIZINFO_CONFIG_DIR", tmp)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if cfg.Defaults.Format != DefaultFormat {
		t.Errorf("Format = %q, want %q", cfg.Defaults.Format, DefaultFormat)
	}
}

func TestSaveAndLoad(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("GBIZINFO_CONFIG_DIR", tmp)

	cfg := &Config{
		Token:    "test-token-123",
		Defaults: Defaults{Format: "table"},
	}
	if err := cfg.Save(); err != nil {
		t.Fatalf("Save() error: %v", err)
	}

	info, err := os.Stat(filepath.Join(tmp, "config.yaml"))
	if err != nil {
		t.Fatalf("Stat() error: %v", err)
	}
	if perm := info.Mode().Perm(); perm != 0600 {
		t.Errorf("file perm = %o, want 0600", perm)
	}

	loaded, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if loaded.Token != "test-token-123" {
		t.Errorf("Token = %q, want %q", loaded.Token, "test-token-123")
	}
	if loaded.Defaults.Format != "table" {
		t.Errorf("Format = %q, want %q", loaded.Defaults.Format, "table")
	}
}

func TestEnvOverride(t *testing.T) {
	t.Setenv("GBIZINFO_TOKEN", "env-token")
	if got := EnvOr(EnvToken, "fallback"); got != "env-token" {
		t.Errorf("EnvOr() = %q, want %q", got, "env-token")
	}
}

func TestEnvOrFallback(t *testing.T) {
	if got := EnvOr("GBIZINFO_NONEXISTENT", "fallback"); got != "fallback" {
		t.Errorf("EnvOr() = %q, want %q", got, "fallback")
	}
}

func TestMaskToken(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"abcdefghij", "abcd******"},
		{"abc", "***"},
		{"", ""},
	}
	for _, tt := range tests {
		if got := MaskToken(tt.input); got != tt.want {
			t.Errorf("MaskToken(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}
