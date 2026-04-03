package cmdutil

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"

	cerrors "github.com/planitaicojp/gbizinfo-cli/internal/errors"
)

// helper: create a cobra command with common flags matching root.go
func newTestCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "test"}
	cmd.Flags().StringP("token", "t", "", "APIトークン")
	cmd.Flags().StringP("format", "f", "", "出力形式")
	cmd.Flags().Bool("verbose", false, "詳細出力")
	cmd.Flags().StringP("corporate-number", "c", "", "法人番号")
	return cmd
}

// --- CorporateNumberArg ---

func TestCorporateNumberArg_ValidPositional(t *testing.T) {
	cmd := newTestCmd()
	cn, err := CorporateNumberArg(cmd, []string{"1234567890123"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cn != "1234567890123" {
		t.Errorf("got %q, want %q", cn, "1234567890123")
	}
}

func TestCorporateNumberArg_ValidFlag(t *testing.T) {
	cmd := newTestCmd()
	_ = cmd.Flags().Set("corporate-number", "9876543210987")
	cn, err := CorporateNumberArg(cmd, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cn != "9876543210987" {
		t.Errorf("got %q, want %q", cn, "9876543210987")
	}
}

func TestCorporateNumberArg_InvalidTooShort(t *testing.T) {
	cmd := newTestCmd()
	_, err := CorporateNumberArg(cmd, []string{"12345"})
	if err == nil {
		t.Fatal("expected error for short corporate number")
	}
	var ve *cerrors.ValidationError
	if !isValidationError(err, &ve) {
		t.Errorf("expected ValidationError, got %T", err)
	}
}

func TestCorporateNumberArg_InvalidNonDigit(t *testing.T) {
	cmd := newTestCmd()
	_, err := CorporateNumberArg(cmd, []string{"123456789012a"})
	if err == nil {
		t.Fatal("expected error for non-digit corporate number")
	}
}

func TestCorporateNumberArg_Missing(t *testing.T) {
	cmd := newTestCmd()
	_, err := CorporateNumberArg(cmd, nil)
	if err == nil {
		t.Fatal("expected error for missing corporate number")
	}
	var ve *cerrors.ValidationError
	if !isValidationError(err, &ve) {
		t.Errorf("expected ValidationError, got %T", err)
	}
}

func TestCorporateNumberArg_InvalidFlag(t *testing.T) {
	cmd := newTestCmd()
	_ = cmd.Flags().Set("corporate-number", "abc")
	_, err := CorporateNumberArg(cmd, nil)
	if err == nil {
		t.Fatal("expected error for invalid flag value")
	}
}

// --- ExactArgs ---

func TestExactArgs_Correct(t *testing.T) {
	cmd := newTestCmd()
	fn := ExactArgs(2)
	if err := fn(cmd, []string{"a", "b"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestExactArgs_TooFew(t *testing.T) {
	cmd := newTestCmd()
	fn := ExactArgs(2)
	err := fn(cmd, []string{"a"})
	if err == nil {
		t.Fatal("expected error for too few args")
	}
	var ve *cerrors.ValidationError
	if !isValidationError(err, &ve) {
		t.Errorf("expected ValidationError, got %T", err)
	}
}

func TestExactArgs_TooMany(t *testing.T) {
	cmd := newTestCmd()
	fn := ExactArgs(1)
	err := fn(cmd, []string{"a", "b", "c"})
	if err == nil {
		t.Fatal("expected error for too many args")
	}
}

// --- GetFormat ---

func TestGetFormat_Flag(t *testing.T) {
	cmd := newTestCmd()
	_ = cmd.Flags().Set("format", "table")
	if got := GetFormat(cmd); got != "table" {
		t.Errorf("got %q, want %q", got, "table")
	}
}

func TestGetFormat_Env(t *testing.T) {
	cmd := newTestCmd()
	t.Setenv("GBIZINFO_FORMAT", "csv")
	if got := GetFormat(cmd); got != "csv" {
		t.Errorf("got %q, want %q", got, "csv")
	}
}

func TestGetFormat_FlagOverridesEnv(t *testing.T) {
	cmd := newTestCmd()
	t.Setenv("GBIZINFO_FORMAT", "csv")
	_ = cmd.Flags().Set("format", "table")
	if got := GetFormat(cmd); got != "table" {
		t.Errorf("flag should override env: got %q, want %q", got, "table")
	}
}

func TestGetFormat_Default(t *testing.T) {
	cmd := newTestCmd()
	t.Setenv("GBIZINFO_FORMAT", "")
	// point config dir to temp so no config file is found
	t.Setenv("GBIZINFO_CONFIG_DIR", t.TempDir())
	if got := GetFormat(cmd); got != "json" {
		t.Errorf("got %q, want %q", got, "json")
	}
}

func TestGetFormat_ConfigFile(t *testing.T) {
	cmd := newTestCmd()
	t.Setenv("GBIZINFO_FORMAT", "")
	dir := t.TempDir()
	t.Setenv("GBIZINFO_CONFIG_DIR", dir)
	// write a config file with format=table
	if err := os.WriteFile(filepath.Join(dir, "config.yaml"), []byte("defaults:\n  format: table\n"), 0600); err != nil {
		t.Fatal(err)
	}
	if got := GetFormat(cmd); got != "table" {
		t.Errorf("got %q, want %q", got, "table")
	}
}

// --- NewClient ---

func TestNewClient_TokenFromFlag(t *testing.T) {
	cmd := newTestCmd()
	_ = cmd.Flags().Set("token", "test-token-123")
	client, err := NewClient(cmd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if client.Token != "test-token-123" {
		t.Errorf("Token = %q, want %q", client.Token, "test-token-123")
	}
}

func TestNewClient_TokenFromEnv(t *testing.T) {
	cmd := newTestCmd()
	t.Setenv("GBIZINFO_TOKEN", "env-token-456")
	client, err := NewClient(cmd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if client.Token != "env-token-456" {
		t.Errorf("Token = %q, want %q", client.Token, "env-token-456")
	}
}

func TestNewClient_NoToken(t *testing.T) {
	cmd := newTestCmd()
	t.Setenv("GBIZINFO_TOKEN", "")
	t.Setenv("GBIZINFO_CONFIG_DIR", t.TempDir())
	_, err := NewClient(cmd)
	if err == nil {
		t.Fatal("expected AuthError when no token is set")
	}
	var ae *cerrors.AuthError
	if !isAuthError(err, &ae) {
		t.Errorf("expected AuthError, got %T: %v", err, err)
	}
}

func TestNewClient_FlagOverridesEnv(t *testing.T) {
	cmd := newTestCmd()
	t.Setenv("GBIZINFO_TOKEN", "env-token")
	_ = cmd.Flags().Set("token", "flag-token")
	client, err := NewClient(cmd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if client.Token != "flag-token" {
		t.Errorf("flag should override env: Token = %q, want %q", client.Token, "flag-token")
	}
}

// helpers for error type assertions
func isValidationError(err error, target **cerrors.ValidationError) bool {
	ve, ok := err.(*cerrors.ValidationError)
	if ok && target != nil {
		*target = ve
	}
	return ok
}

func isAuthError(err error, target **cerrors.AuthError) bool {
	ae, ok := err.(*cerrors.AuthError)
	if ok && target != nil {
		*target = ae
	}
	return ok
}
