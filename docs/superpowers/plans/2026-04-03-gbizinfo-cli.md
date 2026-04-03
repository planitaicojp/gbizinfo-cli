# gbizinfo-cli Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build a Go CLI tool that wraps all 17 gBizINFO REST API endpoints with JSON/Table/CSV output, config file support, and Japanese help messages.

**Architecture:** Cobra-based CLI with layered internal packages (api, config, model, output, errors). Single API client since gBizINFO is one service. Commands map 1:1 to API endpoints. Reference pattern: conoha-cli.

**Tech Stack:** Go, spf13/cobra, gopkg.in/yaml.v3, standard library (net/http, text/tabwriter, encoding/csv, encoding/json)

---

## File Structure

```
gbizinfo-cli/
├── main.go                          # Entry point, calls cmd.Execute()
├── go.mod                           # Module: github.com/planitaicojp/gbizinfo-cli
├── Makefile                         # build, test, lint, clean, install
├── .golangci.yml                    # Linter config
├── cmd/
│   ├── root.go                      # Root command, global flags, Execute()
│   ├── version.go                   # version subcommand
│   ├── completion.go                # completion subcommand
│   ├── search.go                    # search subcommand (GET /v1/hojin)
│   ├── get.go                       # get subcommand (GET /v1/hojin/{corp})
│   ├── certification.go             # certification subcommand
│   ├── commendation.go              # commendation subcommand
│   ├── finance.go                   # finance subcommand
│   ├── patent.go                    # patent subcommand
│   ├── procurement.go               # procurement subcommand
│   ├── subsidy.go                   # subsidy subcommand
│   ├── workplace.go                 # workplace subcommand
│   ├── update/
│   │   └── update.go               # update group + 8 subcommands
│   ├── config/
│   │   └── config.go               # config init/set/show subcommands
│   └── cmdutil/
│       └── cmdutil.go              # NewClient(), GetFormat(), ExactArgs()
├── internal/
│   ├── errors/
│   │   └── errors.go              # Error types + exit codes
│   ├── output/
│   │   ├── formatter.go           # Formatter interface + New()
│   │   ├── json.go                # JSON formatter
│   │   ├── table.go               # Table formatter
│   │   └── csv.go                 # CSV formatter
│   ├── config/
│   │   └── config.go             # Config load/save, env vars, defaults
│   ├── model/
│   │   ├── hojin.go              # Hojin, HojinResponse, PageInfo
│   │   ├── certification.go      # Certification model
│   │   ├── commendation.go       # Commendation model
│   │   ├── finance.go            # Finance model
│   │   ├── patent.go             # Patent model
│   │   ├── procurement.go        # Procurement model
│   │   ├── subsidy.go            # Subsidy model
│   │   ├── workplace.go          # Workplace model
│   │   └── update.go             # UpdateResponse model
│   └── api/
│       ├── client.go             # Base HTTP client
│       ├── hojin.go              # 9 hojin methods
│       └── update.go             # 8 updateInfo methods
└── test/
    └── fixtures/                 # Test JSON response files
```

---

### Task 1: Project Scaffolding

**Files:**
- Create: `go.mod`
- Create: `main.go`
- Create: `Makefile`
- Create: `.golangci.yml`

- [ ] **Step 1: Initialize Go module**

Run: `cd /root/dev/planitai/gbizinfo-cli && go mod init github.com/planitaicojp/gbizinfo-cli`

- [ ] **Step 2: Create main.go**

```go
package main

import "github.com/planitaicojp/gbizinfo-cli/cmd"

func main() {
	cmd.Execute()
}
```

- [ ] **Step 3: Create Makefile**

```makefile
BINARY := gbizinfo
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS := -ldflags "-s -w -X github.com/planitaicojp/gbizinfo-cli/cmd.version=$(VERSION)"

.PHONY: build test lint clean install

build:
	go build $(LDFLAGS) -o $(BINARY) .

install:
	go install $(LDFLAGS) .

test:
	go test ./... -v

lint:
	golangci-lint run ./...

clean:
	rm -f $(BINARY)

coverage:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
```

- [ ] **Step 4: Create .golangci.yml**

```yaml
version: "2"
linters:
  default: none
  enable:
    - govet
    - ineffassign
    - staticcheck
    - unused
    - errcheck
  settings:
    errcheck:
      check-type-assertions: false
      check-blank: false
      disable-default-exclusions: false
      exclude-functions:
        - io.Copy
        - (io.Closer).Close
        - (*os.File).Close
        - (net/http.ResponseWriter).Write
        - (*encoding/json.Encoder).Encode
        - (*encoding/json.Decoder).Decode
formatters:
  enable:
    - gofmt
```

- [ ] **Step 5: Commit**

```bash
git add go.mod main.go Makefile .golangci.yml
git commit -m "chore: initialize project scaffolding"
```

---

### Task 2: Error Types and Exit Codes

**Files:**
- Create: `internal/errors/errors.go`
- Test: `internal/errors/errors_test.go`

- [ ] **Step 1: Write the failing test**

```go
// internal/errors/errors_test.go
package errors

import (
	"testing"
)

func TestGetExitCode(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want int
	}{
		{"nil", nil, ExitOK},
		{"auth", &AuthError{Message: "bad token"}, ExitAuth},
		{"not found", &NotFoundError{Resource: "法人", ID: "123"}, ExitNotFound},
		{"validation", &ValidationError{Field: "name", Message: "required"}, ExitValidation},
		{"api", &APIError{StatusCode: 500, Message: "server error"}, ExitAPI},
		{"rate limit", &RateLimitError{Message: "too many requests"}, ExitAPI},
		{"generic", fmt.Errorf("unknown"), ExitGeneral},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetExitCode(tt.err); got != tt.want {
				t.Errorf("GetExitCode() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestErrorMessages(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want string
	}{
		{"auth", &AuthError{Message: "invalid token"}, "認証エラー: invalid token"},
		{"not found", &NotFoundError{Resource: "法人", ID: "123"}, "法人が見つかりません: 123"},
		{"validation", &ValidationError{Field: "name", Message: "required"}, "入力エラー (name): required"},
		{"validation no field", &ValidationError{Message: "bad input"}, "入力エラー: bad input"},
		{"api", &APIError{StatusCode: 500, Message: "fail"}, "APIエラー (HTTP 500): fail"},
		{"api with code", &APIError{StatusCode: 400, Code: "BAD", Message: "fail"}, "APIエラー (HTTP 400, BAD): fail"},
		{"rate limit", &RateLimitError{Message: "limit exceeded"}, "APIレート制限: limit exceeded"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.want {
				t.Errorf("Error() = %q, want %q", got, tt.want)
			}
		})
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `cd /root/dev/planitai/gbizinfo-cli && go test ./internal/errors/ -v`
Expected: FAIL (package does not exist)

- [ ] **Step 3: Write implementation**

```go
// internal/errors/errors.go
package errors

import "fmt"

// Exit codes
const (
	ExitOK         = 0
	ExitGeneral    = 1
	ExitAuth       = 2
	ExitNotFound   = 3
	ExitValidation = 4
	ExitAPI        = 5
)

// ExitCoder is implemented by errors that carry a process exit code.
type ExitCoder interface {
	ExitCode() int
}

// AuthError represents an authentication failure.
type AuthError struct {
	Message string
}

func (e *AuthError) Error() string {
	return fmt.Sprintf("認証エラー: %s", e.Message)
}

func (e *AuthError) ExitCode() int {
	return ExitAuth
}

// NotFoundError indicates that a requested resource was not found.
type NotFoundError struct {
	Resource string
	ID       string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("%sが見つかりません: %s", e.Resource, e.ID)
}

func (e *NotFoundError) ExitCode() int {
	return ExitNotFound
}

// ValidationError represents invalid user input.
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("入力エラー (%s): %s", e.Field, e.Message)
	}
	return fmt.Sprintf("入力エラー: %s", e.Message)
}

func (e *ValidationError) ExitCode() int {
	return ExitValidation
}

// APIError represents an error returned by the gBizINFO API.
type APIError struct {
	StatusCode int
	Code       string
	Message    string
}

func (e *APIError) Error() string {
	if e.Code != "" {
		return fmt.Sprintf("APIエラー (HTTP %d, %s): %s", e.StatusCode, e.Code, e.Message)
	}
	return fmt.Sprintf("APIエラー (HTTP %d): %s", e.StatusCode, e.Message)
}

func (e *APIError) ExitCode() int {
	return ExitAPI
}

// RateLimitError represents a 429 Too Many Requests response.
type RateLimitError struct {
	Message string
}

func (e *RateLimitError) Error() string {
	return fmt.Sprintf("APIレート制限: %s", e.Message)
}

func (e *RateLimitError) ExitCode() int {
	return ExitAPI
}

// GetExitCode returns the exit code for the given error.
func GetExitCode(err error) int {
	if err == nil {
		return ExitOK
	}
	if ec, ok := err.(ExitCoder); ok {
		return ec.ExitCode()
	}
	return ExitGeneral
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `cd /root/dev/planitai/gbizinfo-cli && go test ./internal/errors/ -v`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add internal/errors/
git commit -m "feat: add error types and exit codes"
```

---

### Task 3: Output Formatters

**Files:**
- Create: `internal/output/formatter.go`
- Create: `internal/output/json.go`
- Create: `internal/output/table.go`
- Create: `internal/output/csv.go`
- Test: `internal/output/formatter_test.go`

- [ ] **Step 1: Write the failing test**

```go
// internal/output/formatter_test.go
package output

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

type testItem struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Number  string `json:"number"`
}

func TestJSONFormatter(t *testing.T) {
	items := []testItem{
		{Name: "テスト株式会社", Address: "東京都千代田区", Number: "1234567890123"},
	}
	var buf bytes.Buffer
	f := New("json")
	if err := f.Format(&buf, items); err != nil {
		t.Fatalf("Format() error: %v", err)
	}
	var result []testItem
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}
	if result[0].Name != "テスト株式会社" {
		t.Errorf("Name = %q, want %q", result[0].Name, "テスト株式会社")
	}
}

func TestTableFormatter(t *testing.T) {
	items := []testItem{
		{Name: "テスト株式会社", Address: "東京都", Number: "123"},
	}
	var buf bytes.Buffer
	f := New("table")
	if err := f.Format(&buf, items); err != nil {
		t.Fatalf("Format() error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "NAME") {
		t.Errorf("table output missing header NAME, got: %s", out)
	}
	if !strings.Contains(out, "テスト株式会社") {
		t.Errorf("table output missing data, got: %s", out)
	}
}

func TestTableFormatterEmpty(t *testing.T) {
	var items []testItem
	var buf bytes.Buffer
	f := New("table")
	if err := f.Format(&buf, items); err != nil {
		t.Fatalf("Format() error: %v", err)
	}
	if buf.Len() != 0 {
		t.Errorf("expected empty output for empty slice, got: %q", buf.String())
	}
}

func TestCSVFormatter(t *testing.T) {
	items := []testItem{
		{Name: "テスト株式会社", Address: "東京都", Number: "123"},
		{Name: "サンプル合同会社", Address: "大阪府", Number: "456"},
	}
	var buf bytes.Buffer
	f := New("csv")
	if err := f.Format(&buf, items); err != nil {
		t.Fatalf("Format() error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines (header + 2 rows), got %d", len(lines))
	}
	if lines[0] != "name,address,number" {
		t.Errorf("header = %q, want %q", lines[0], "name,address,number")
	}
}

func TestNewDefaultsToJSON(t *testing.T) {
	f := New("unknown")
	if _, ok := f.(*JSONFormatter); !ok {
		t.Errorf("New(unknown) should return JSONFormatter, got %T", f)
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `cd /root/dev/planitai/gbizinfo-cli && go test ./internal/output/ -v`
Expected: FAIL

- [ ] **Step 3: Write formatter.go**

```go
// internal/output/formatter.go
package output

import "io"

// Formatter formats and writes data to a writer.
type Formatter interface {
	Format(w io.Writer, data any) error
}

// New creates a formatter for the given format name.
// Defaults to JSON if the format is unknown.
func New(format string) Formatter {
	switch format {
	case "table":
		return &TableFormatter{}
	case "csv":
		return &CSVFormatter{}
	default:
		return &JSONFormatter{}
	}
}
```

- [ ] **Step 4: Write json.go**

```go
// internal/output/json.go
package output

import (
	"encoding/json"
	"io"
)

// JSONFormatter outputs data as pretty-printed JSON.
type JSONFormatter struct{}

func (f *JSONFormatter) Format(w io.Writer, data any) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	enc.SetEscapeHTML(false)
	return enc.Encode(data)
}
```

- [ ] **Step 5: Write table.go**

```go
// internal/output/table.go
package output

import (
	"fmt"
	"io"
	"reflect"
	"strings"
	"text/tabwriter"
)

// TableFormatter outputs data as an aligned text table.
type TableFormatter struct{}

func (f *TableFormatter) Format(w io.Writer, data any) error {
	val := reflect.ValueOf(data)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Slice {
		_, err := fmt.Fprintf(w, "%v\n", data)
		return err
	}

	if val.Len() == 0 {
		return nil
	}

	tw := tabwriter.NewWriter(w, 0, 4, 2, ' ', 0)

	elem := val.Index(0)
	if elem.Kind() == reflect.Ptr {
		elem = elem.Elem()
	}
	elemType := elem.Type()

	headers := make([]string, elemType.NumField())
	for i := 0; i < elemType.NumField(); i++ {
		field := elemType.Field(i)
		name := field.Tag.Get("json")
		if idx := strings.Index(name, ","); idx != -1 {
			name = name[:idx]
		}
		if name == "" || name == "-" {
			name = field.Name
		}
		headers[i] = strings.ToUpper(name)
	}
	if _, err := fmt.Fprintln(tw, strings.Join(headers, "\t")); err != nil {
		return err
	}

	for i := 0; i < val.Len(); i++ {
		row := val.Index(i)
		if row.Kind() == reflect.Ptr {
			row = row.Elem()
		}
		fields := make([]string, row.NumField())
		for j := 0; j < row.NumField(); j++ {
			fields[j] = fmt.Sprintf("%v", row.Field(j).Interface())
		}
		if _, err := fmt.Fprintln(tw, strings.Join(fields, "\t")); err != nil {
			return err
		}
	}

	return tw.Flush()
}
```

- [ ] **Step 6: Write csv.go**

```go
// internal/output/csv.go
package output

import (
	"encoding/csv"
	"fmt"
	"io"
	"reflect"
	"strings"
)

// CSVFormatter outputs data as CSV.
type CSVFormatter struct{}

func (f *CSVFormatter) Format(w io.Writer, data any) error {
	val := reflect.ValueOf(data)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Slice {
		return fmt.Errorf("CSVフォーマットにはスライスが必要です")
	}
	if val.Len() == 0 {
		return nil
	}

	writer := csv.NewWriter(w)
	defer writer.Flush()

	elem := val.Index(0)
	if elem.Kind() == reflect.Ptr {
		elem = elem.Elem()
	}
	elemType := elem.Type()

	headers := make([]string, elemType.NumField())
	for i := 0; i < elemType.NumField(); i++ {
		field := elemType.Field(i)
		name := field.Tag.Get("json")
		if idx := strings.Index(name, ","); idx != -1 {
			name = name[:idx]
		}
		if name == "" || name == "-" {
			name = field.Name
		}
		headers[i] = name
	}
	if err := writer.Write(headers); err != nil {
		return err
	}

	for i := 0; i < val.Len(); i++ {
		row := val.Index(i)
		if row.Kind() == reflect.Ptr {
			row = row.Elem()
		}
		record := make([]string, row.NumField())
		for j := 0; j < row.NumField(); j++ {
			record[j] = fmt.Sprintf("%v", row.Field(j).Interface())
		}
		if err := writer.Write(record); err != nil {
			return err
		}
	}
	return nil
}
```

- [ ] **Step 7: Run tests to verify they pass**

Run: `cd /root/dev/planitai/gbizinfo-cli && go test ./internal/output/ -v`
Expected: PASS

- [ ] **Step 8: Commit**

```bash
git add internal/output/
git commit -m "feat: add output formatters (JSON, Table, CSV)"
```

---

### Task 4: Config Management

**Files:**
- Create: `internal/config/config.go`
- Test: `internal/config/config_test.go`

- [ ] **Step 1: Write the failing test**

```go
// internal/config/config_test.go
package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadDefault(t *testing.T) {
	// Point to a temp dir with no config file
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

	// Check file permissions
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
```

- [ ] **Step 2: Run test to verify it fails**

Run: `cd /root/dev/planitai/gbizinfo-cli && go test ./internal/config/ -v`
Expected: FAIL

- [ ] **Step 3: Write implementation**

```go
// internal/config/config.go
package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

const (
	DefaultFormat = "json"
	configFile    = "config.yaml"
)

// Environment variable names
const (
	EnvConfigDir = "GBIZINFO_CONFIG_DIR"
	EnvToken     = "GBIZINFO_TOKEN"
	EnvFormat    = "GBIZINFO_FORMAT"
)

// Config represents the CLI configuration.
type Config struct {
	Token    string   `yaml:"token"`
	Defaults Defaults `yaml:"defaults"`
}

// Defaults holds default settings.
type Defaults struct {
	Format string `yaml:"format"`
}

// DefaultConfigDir returns the config directory path.
func DefaultConfigDir() string {
	if d := os.Getenv(EnvConfigDir); d != "" {
		return d
	}
	if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
		return filepath.Join(xdg, "gbizinfo")
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "gbizinfo")
}

// Load reads the config file or returns defaults.
func Load() (*Config, error) {
	dir := DefaultConfigDir()
	path := filepath.Join(dir, configFile)

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return defaultConfig(), nil
		}
		return nil, fmt.Errorf("設定ファイルの読み込みに失敗: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("設定ファイルの解析に失敗: %w", err)
	}
	if cfg.Defaults.Format == "" {
		cfg.Defaults.Format = DefaultFormat
	}
	return &cfg, nil
}

// Save writes the config to disk.
func (c *Config) Save() error {
	dir := DefaultConfigDir()
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("設定ディレクトリの作成に失敗: %w", err)
	}

	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("設定のシリアライズに失敗: %w", err)
	}
	return os.WriteFile(filepath.Join(dir, configFile), data, 0600)
}

// EnvOr returns the environment variable value or the fallback.
func EnvOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// MaskToken masks a token for display.
func MaskToken(token string) string {
	if token == "" {
		return ""
	}
	if len(token) <= 4 {
		return strings.Repeat("*", len(token))
	}
	return token[:4] + strings.Repeat("*", 6)
}

func defaultConfig() *Config {
	return &Config{
		Defaults: Defaults{Format: DefaultFormat},
	}
}
```

- [ ] **Step 4: Add yaml.v3 dependency**

Run: `cd /root/dev/planitai/gbizinfo-cli && go get gopkg.in/yaml.v3`

- [ ] **Step 5: Run tests to verify they pass**

Run: `cd /root/dev/planitai/gbizinfo-cli && go test ./internal/config/ -v`
Expected: PASS

- [ ] **Step 6: Commit**

```bash
git add internal/config/ go.mod go.sum
git commit -m "feat: add config management with YAML persistence"
```

---

### Task 5: API Client Base

**Files:**
- Create: `internal/api/client.go`
- Test: `internal/api/client_test.go`

- [ ] **Step 1: Write the failing test**

```go
// internal/api/client_test.go
package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClientGet(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-hojinInfo-api-token") != "test-token" {
			t.Errorf("missing auth header")
		}
		if r.Header.Get("Accept") != "application/json" {
			t.Errorf("missing accept header")
		}
		if r.Method != http.MethodGet {
			t.Errorf("method = %s, want GET", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-token")
	var result map[string]string
	err := client.Get("/test", &result)
	if err != nil {
		t.Fatalf("Get() error: %v", err)
	}
	if result["status"] != "ok" {
		t.Errorf("status = %q, want %q", result["status"], "ok")
	}
}

func TestClientGet401(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(401)
		w.Write([]byte(`{"message":"Unauthorized"}`))
	}))
	defer server.Close()

	client := NewClient(server.URL, "bad-token")
	var result map[string]string
	err := client.Get("/test", &result)
	if err == nil {
		t.Fatal("expected error for 401")
	}
	if _, ok := err.(*cerrors.AuthError); !ok {
		t.Errorf("expected AuthError, got %T: %v", err, err)
	}
}

func TestClientGet404(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		w.Write([]byte(`{"message":"Not Found"}`))
	}))
	defer server.Close()

	client := NewClient(server.URL, "token")
	var result map[string]string
	err := client.Get("/test", &result)
	if err == nil {
		t.Fatal("expected error for 404")
	}
	if _, ok := err.(*cerrors.NotFoundError); !ok {
		t.Errorf("expected NotFoundError, got %T: %v", err, err)
	}
}

func TestClientGet429(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(429)
		w.Write([]byte(`{"message":"Too Many Requests"}`))
	}))
	defer server.Close()

	client := NewClient(server.URL, "token")
	var result map[string]string
	err := client.Get("/test", &result)
	if err == nil {
		t.Fatal("expected error for 429")
	}
	if _, ok := err.(*cerrors.RateLimitError); !ok {
		t.Errorf("expected RateLimitError, got %T: %v", err, err)
	}
}

func TestClientUserAgent(t *testing.T) {
	var gotUA string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotUA = r.Header.Get("User-Agent")
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{}`))
	}))
	defer server.Close()

	client := NewClient(server.URL, "token")
	client.Get("/test", nil)
	if gotUA != UserAgent {
		t.Errorf("User-Agent = %q, want %q", gotUA, UserAgent)
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `cd /root/dev/planitai/gbizinfo-cli && go test ./internal/api/ -v`
Expected: FAIL

- [ ] **Step 3: Write implementation**

```go
// internal/api/client.go
package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	cerrors "github.com/planitaicojp/gbizinfo-cli/internal/errors"
)

// UserAgent is the User-Agent header sent with all requests.
var UserAgent = "gbizinfo-cli/dev"

const defaultTimeout = 30 * time.Second

// Client is the HTTP client for gBizINFO API.
type Client struct {
	HTTP    *http.Client
	BaseURL string
	Token   string
	Verbose bool
}

// NewClient creates a new API client.
func NewClient(baseURL, token string) *Client {
	return &Client{
		HTTP:    &http.Client{Timeout: defaultTimeout},
		BaseURL: baseURL,
		Token:   token,
	}
}

// Get performs a GET request and decodes the JSON response.
func (c *Client) Get(path string, result any) error {
	url := c.BaseURL + path

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("リクエストの作成に失敗: %w", err)
	}

	resp, err := c.do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if result != nil {
		return json.NewDecoder(resp.Body).Decode(result)
	}
	return nil
}

func (c *Client) do(req *http.Request) (*http.Response, error) {
	req.Header.Set("User-Agent", UserAgent)
	req.Header.Set("Accept", "application/json")
	if c.Token != "" {
		req.Header.Set("X-hojinInfo-api-token", c.Token)
	}

	if c.Verbose {
		debugLogRequest(req)
	}

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, fmt.Errorf("リクエストの送信に失敗: %w", err)
	}

	if c.Verbose {
		debugLogResponse(resp)
	}

	if resp.StatusCode >= 400 {
		return nil, parseAPIError(resp)
	}

	return resp, nil
}

func parseAPIError(resp *http.Response) error {
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	message := string(body)
	var errResp struct {
		Message string `json:"message"`
	}
	if json.Unmarshal(body, &errResp) == nil && errResp.Message != "" {
		message = errResp.Message
	}

	switch resp.StatusCode {
	case 401, 403:
		return &cerrors.AuthError{Message: message}
	case 404:
		return &cerrors.NotFoundError{Resource: "リソース", ID: ""}
	case 429:
		return &cerrors.RateLimitError{Message: message}
	default:
		return &cerrors.APIError{StatusCode: resp.StatusCode, Message: message}
	}
}

func debugLogRequest(req *http.Request) {
	fmt.Fprintf(io.Discard, "")
	// TODO: implement verbose logging to stderr in a later step
}

func debugLogResponse(resp *http.Response) {
	// TODO: implement verbose logging to stderr in a later step
}
```

- [ ] **Step 4: Run tests to verify they pass**

Run: `cd /root/dev/planitai/gbizinfo-cli && go test ./internal/api/ -v`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add internal/api/client.go internal/api/client_test.go
git commit -m "feat: add base API client with error handling"
```

---

### Task 6: Data Models

**Files:**
- Create: `internal/model/hojin.go`
- Create: `internal/model/certification.go`
- Create: `internal/model/commendation.go`
- Create: `internal/model/finance.go`
- Create: `internal/model/patent.go`
- Create: `internal/model/procurement.go`
- Create: `internal/model/subsidy.go`
- Create: `internal/model/workplace.go`
- Create: `internal/model/update.go`

Note: Model field definitions must be verified against the actual gBizINFO API Swagger spec. The models below are based on information gathered from the Ruby gem and blog posts. After the initial scaffolding, run a real API call to compare and adjust fields.

- [ ] **Step 1: Create hojin.go**

```go
// internal/model/hojin.go
package model

// PageInfo holds pagination metadata.
type PageInfo struct {
	TotalCount int `json:"totalCount"`
	TotalPage  int `json:"totalPage"`
	PageNumber int `json:"pageNumber"`
}

// HojinResponse is the response from GET /v1/hojin (search).
type HojinResponse struct {
	PageInfo
	Corporations []Hojin `json:"hojin-infos"`
}

// Hojin represents basic corporate information.
type Hojin struct {
	CorporateNumber string `json:"corporate_number"`
	Name            string `json:"name"`
	NameKana        string `json:"kana"`
	Location        string `json:"location"`
	Status          string `json:"status"`
	UpdateDate      string `json:"update_date"`
	Capital         string `json:"capital_stock"`
	EmployeeNumber  string `json:"employee_number"`
	RepresentName   string `json:"represent_name"`
	CompanyURL      string `json:"company_url"`
	DateOfEstablish string `json:"date_of_establishment"`
	BusinessSummary string `json:"business_summary"`
}

// HojinDetail is a wrapper for single hojin response.
type HojinDetail struct {
	Corporations []Hojin `json:"hojin-infos"`
}
```

- [ ] **Step 2: Create certification.go**

```go
// internal/model/certification.go
package model

// CertificationResponse is the response from GET /v1/hojin/{corp}/certification.
type CertificationResponse struct {
	Corporations []CertificationInfo `json:"hojin-infos"`
}

// CertificationInfo wraps certifications for a corporation.
type CertificationInfo struct {
	CorporateNumber string          `json:"corporate_number"`
	Name            string          `json:"name"`
	Certifications  []Certification `json:"certification"`
}

// Certification represents a single certification/registration.
type Certification struct {
	Title       string `json:"title"`
	DateOfApproval string `json:"date_of_approval"`
	Target      string `json:"target"`
	Category    string `json:"category"`
	EnterpriseScale string `json:"enterprise_scale"`
	GovernmentDepartments string `json:"government_departments"`
}
```

- [ ] **Step 3: Create commendation.go**

```go
// internal/model/commendation.go
package model

// CommendationResponse is the response from GET /v1/hojin/{corp}/commendation.
type CommendationResponse struct {
	Corporations []CommendationInfo `json:"hojin-infos"`
}

// CommendationInfo wraps commendations for a corporation.
type CommendationInfo struct {
	CorporateNumber string         `json:"corporate_number"`
	Name            string         `json:"name"`
	Commendations   []Commendation `json:"commendation"`
}

// Commendation represents a single commendation/award.
type Commendation struct {
	Title       string `json:"title"`
	DateOfCommendation string `json:"date_of_commendation"`
	Target      string `json:"target"`
	Category    string `json:"category"`
	GovernmentDepartments string `json:"government_departments"`
}
```

- [ ] **Step 4: Create finance.go**

```go
// internal/model/finance.go
package model

// FinanceResponse is the response from GET /v1/hojin/{corp}/finance.
type FinanceResponse struct {
	Corporations []FinanceInfo `json:"hojin-infos"`
}

// FinanceInfo wraps financial data for a corporation.
type FinanceInfo struct {
	CorporateNumber string    `json:"corporate_number"`
	Name            string    `json:"name"`
	Finance         []Finance `json:"finance"`
}

// Finance represents financial data for a fiscal period.
type Finance struct {
	AccountingPeriod          string `json:"accounting_period"`
	MajorShareholders         string `json:"major_shareholders"`
	NetSales                  string `json:"net_sales"`
	OperatingRevenue          string `json:"operating_revenue"`
	OrdinaryIncome            string `json:"ordinary_income"`
	Profit                    string `json:"profit"`
	TotalAssets               string `json:"total_assets"`
	NetAssets                 string `json:"net_assets"`
	CapitalStock              string `json:"capital_stock"`
	EmployeeNumber            string `json:"employee_number"`
}
```

- [ ] **Step 5: Create patent.go**

```go
// internal/model/patent.go
package model

// PatentResponse is the response from GET /v1/hojin/{corp}/patent.
type PatentResponse struct {
	Corporations []PatentInfo `json:"hojin-infos"`
}

// PatentInfo wraps patent data for a corporation.
type PatentInfo struct {
	CorporateNumber string   `json:"corporate_number"`
	Name            string   `json:"name"`
	Patents         []Patent `json:"patent"`
}

// Patent represents a single patent.
type Patent struct {
	Title             string `json:"title"`
	DateOfApplication string `json:"date_of_application"`
	PatentNumber      string `json:"patent_number"`
	ApplicationNumber string `json:"application_number"`
	ClassificationPI  string `json:"classification_pi"`
}
```

- [ ] **Step 6: Create procurement.go**

```go
// internal/model/procurement.go
package model

// ProcurementResponse is the response from GET /v1/hojin/{corp}/procurement.
type ProcurementResponse struct {
	Corporations []ProcurementInfo `json:"hojin-infos"`
}

// ProcurementInfo wraps procurement data for a corporation.
type ProcurementInfo struct {
	CorporateNumber string        `json:"corporate_number"`
	Name            string        `json:"name"`
	Procurements    []Procurement `json:"procurement"`
}

// Procurement represents a single procurement record.
type Procurement struct {
	Title              string `json:"title"`
	DateOfOrder        string `json:"date_of_order"`
	Amount             string `json:"amount"`
	GovernmentDepartments string `json:"government_departments"`
	JointSignatures    string `json:"joint_signatures"`
}
```

- [ ] **Step 7: Create subsidy.go**

```go
// internal/model/subsidy.go
package model

// SubsidyResponse is the response from GET /v1/hojin/{corp}/subsidy.
type SubsidyResponse struct {
	Corporations []SubsidyInfo `json:"hojin-infos"`
}

// SubsidyInfo wraps subsidy data for a corporation.
type SubsidyInfo struct {
	CorporateNumber string    `json:"corporate_number"`
	Name            string    `json:"name"`
	Subsidies       []Subsidy `json:"subsidy"`
}

// Subsidy represents a single subsidy record.
type Subsidy struct {
	Title              string `json:"title"`
	DateOfApproval     string `json:"date_of_approval"`
	Amount             string `json:"amount"`
	SubsidyResource    string `json:"subsidy_resource"`
	Target             string `json:"target"`
	GovernmentDepartments string `json:"government_departments"`
	Note               string `json:"note"`
}
```

- [ ] **Step 8: Create workplace.go**

```go
// internal/model/workplace.go
package model

// WorkplaceResponse is the response from GET /v1/hojin/{corp}/workplace.
type WorkplaceResponse struct {
	Corporations []WorkplaceInfo `json:"hojin-infos"`
}

// WorkplaceInfo wraps workplace data for a corporation.
type WorkplaceInfo struct {
	CorporateNumber string      `json:"corporate_number"`
	Name            string      `json:"name"`
	Workplaces      []Workplace `json:"workplace_info"`
}

// Workplace represents workplace information.
type Workplace struct {
	BaseMonth                   string `json:"base_month"`
	EmployeeNumber              string `json:"employee_number"`
	EmployeeNumberRegular       string `json:"employee_number_regular"`
	EmployeeNumberNonRegular    string `json:"employee_number_non_regular"`
	FemaleShareOfManager        string `json:"female_share_of_manager"`
	YearsOfService              string `json:"years_of_service"`
	AnnualSalary                string `json:"annual_salary"`
	AverageContinuousServiceYears string `json:"average_continuous_service_years"`
	AverageAge                  string `json:"average_age"`
	MonthAverageOvertimeHours   string `json:"month_average_overtime_hours"`
	PaidHolidayUsageRate        string `json:"paid_holiday_usage_rate"`
}
```

- [ ] **Step 9: Create update.go**

```go
// internal/model/update.go
package model

// UpdateResponse is the response from GET /v1/hojin/updateInfo/* endpoints.
type UpdateResponse struct {
	PageInfo
	Corporations []Hojin `json:"hojin-infos"`
}

// UpdateParams holds parameters for updateInfo queries.
type UpdateParams struct {
	From string
	To   string
	Page int
}

// SearchParams holds parameters for hojin search queries.
type SearchParams struct {
	Name            string
	Address         string
	CorporateNumber string
	Page            int
	Limit           int
}
```

- [ ] **Step 10: Verify compilation**

Run: `cd /root/dev/planitai/gbizinfo-cli && go build ./internal/model/`
Expected: No errors

- [ ] **Step 11: Commit**

```bash
git add internal/model/
git commit -m "feat: add data models for all API resources"
```

---

### Task 7: API Hojin Methods

**Files:**
- Create: `internal/api/hojin.go`
- Test: `internal/api/hojin_test.go`

- [ ] **Step 1: Write the failing test**

```go
// internal/api/hojin_test.go
package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/planitaicojp/gbizinfo-cli/internal/model"
)

func TestSearch(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/hojin" {
			t.Errorf("path = %s, want /v1/hojin", r.URL.Path)
		}
		if r.URL.Query().Get("name") != "テスト" {
			t.Errorf("name param = %q, want %q", r.URL.Query().Get("name"), "テスト")
		}
		if r.URL.Query().Get("page") != "1" {
			t.Errorf("page param = %q, want %q", r.URL.Query().Get("page"), "1")
		}
		resp := model.HojinResponse{
			PageInfo:     model.PageInfo{TotalCount: 1, TotalPage: 1, PageNumber: 1},
			Corporations: []model.Hojin{{CorporateNumber: "1234567890123", Name: "テスト株式会社"}},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL, "token")
	result, err := client.Search(model.SearchParams{Name: "テスト", Page: 1})
	if err != nil {
		t.Fatalf("Search() error: %v", err)
	}
	if result.TotalCount != 1 {
		t.Errorf("TotalCount = %d, want 1", result.TotalCount)
	}
	if result.Corporations[0].Name != "テスト株式会社" {
		t.Errorf("Name = %q, want %q", result.Corporations[0].Name, "テスト株式会社")
	}
}

func TestGetHojin(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/hojin/1234567890123" {
			t.Errorf("path = %s, want /v1/hojin/1234567890123", r.URL.Path)
		}
		resp := model.HojinDetail{
			Corporations: []model.Hojin{{CorporateNumber: "1234567890123", Name: "テスト株式会社"}},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL, "token")
	result, err := client.GetHojin("1234567890123")
	if err != nil {
		t.Fatalf("GetHojin() error: %v", err)
	}
	if result.Corporations[0].CorporateNumber != "1234567890123" {
		t.Errorf("CorporateNumber = %q, want %q", result.Corporations[0].CorporateNumber, "1234567890123")
	}
}

func TestGetFinance(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/hojin/1234567890123/finance" {
			t.Errorf("path = %s", r.URL.Path)
		}
		resp := model.FinanceResponse{
			Corporations: []model.FinanceInfo{{CorporateNumber: "1234567890123", Name: "テスト"}},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL, "token")
	result, err := client.GetFinance("1234567890123")
	if err != nil {
		t.Fatalf("GetFinance() error: %v", err)
	}
	if result.Corporations[0].CorporateNumber != "1234567890123" {
		t.Errorf("CorporateNumber = %q", result.Corporations[0].CorporateNumber)
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `cd /root/dev/planitai/gbizinfo-cli && go test ./internal/api/ -v -run TestSearch`
Expected: FAIL

- [ ] **Step 3: Write implementation**

```go
// internal/api/hojin.go
package api

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/planitaicojp/gbizinfo-cli/internal/model"
)

// Search searches for corporations.
func (c *Client) Search(params model.SearchParams) (*model.HojinResponse, error) {
	q := url.Values{}
	if params.Name != "" {
		q.Set("name", params.Name)
	}
	if params.Address != "" {
		q.Set("exist_flg", params.Address)
	}
	if params.CorporateNumber != "" {
		q.Set("corporate_number", params.CorporateNumber)
	}
	if params.Page > 0 {
		q.Set("page", strconv.Itoa(params.Page))
	}
	if params.Limit > 0 {
		q.Set("limit", strconv.Itoa(params.Limit))
	}

	path := "/v1/hojin"
	if len(q) > 0 {
		path += "?" + q.Encode()
	}

	var result model.HojinResponse
	if err := c.Get(path, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetHojin retrieves basic corporate information.
func (c *Client) GetHojin(corporateNumber string) (*model.HojinDetail, error) {
	var result model.HojinDetail
	if err := c.Get(fmt.Sprintf("/v1/hojin/%s", corporateNumber), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetCertification retrieves certification data.
func (c *Client) GetCertification(corporateNumber string) (*model.CertificationResponse, error) {
	var result model.CertificationResponse
	if err := c.Get(fmt.Sprintf("/v1/hojin/%s/certification", corporateNumber), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetCommendation retrieves commendation data.
func (c *Client) GetCommendation(corporateNumber string) (*model.CommendationResponse, error) {
	var result model.CommendationResponse
	if err := c.Get(fmt.Sprintf("/v1/hojin/%s/commendation", corporateNumber), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetFinance retrieves financial data.
func (c *Client) GetFinance(corporateNumber string) (*model.FinanceResponse, error) {
	var result model.FinanceResponse
	if err := c.Get(fmt.Sprintf("/v1/hojin/%s/finance", corporateNumber), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetPatent retrieves patent data.
func (c *Client) GetPatent(corporateNumber string) (*model.PatentResponse, error) {
	var result model.PatentResponse
	if err := c.Get(fmt.Sprintf("/v1/hojin/%s/patent", corporateNumber), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetProcurement retrieves procurement data.
func (c *Client) GetProcurement(corporateNumber string) (*model.ProcurementResponse, error) {
	var result model.ProcurementResponse
	if err := c.Get(fmt.Sprintf("/v1/hojin/%s/procurement", corporateNumber), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetSubsidy retrieves subsidy data.
func (c *Client) GetSubsidy(corporateNumber string) (*model.SubsidyResponse, error) {
	var result model.SubsidyResponse
	if err := c.Get(fmt.Sprintf("/v1/hojin/%s/subsidy", corporateNumber), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetWorkplace retrieves workplace data.
func (c *Client) GetWorkplace(corporateNumber string) (*model.WorkplaceResponse, error) {
	var result model.WorkplaceResponse
	if err := c.Get(fmt.Sprintf("/v1/hojin/%s/workplace", corporateNumber), &result); err != nil {
		return nil, err
	}
	return &result, nil
}
```

- [ ] **Step 4: Run tests to verify they pass**

Run: `cd /root/dev/planitai/gbizinfo-cli && go test ./internal/api/ -v`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add internal/api/hojin.go internal/api/hojin_test.go
git commit -m "feat: add API methods for hojin endpoints"
```

---

### Task 8: API Update Methods

**Files:**
- Create: `internal/api/update.go`
- Test: `internal/api/update_test.go`

- [ ] **Step 1: Write the failing test**

```go
// internal/api/update_test.go
package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/planitaicojp/gbizinfo-cli/internal/model"
)

func TestGetUpdateInfo(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/hojin/updateInfo" {
			t.Errorf("path = %s, want /v1/hojin/updateInfo", r.URL.Path)
		}
		if r.URL.Query().Get("from") != "2024-01-01" {
			t.Errorf("from = %q", r.URL.Query().Get("from"))
		}
		if r.URL.Query().Get("to") != "2024-01-31" {
			t.Errorf("to = %q", r.URL.Query().Get("to"))
		}
		resp := model.UpdateResponse{
			PageInfo:     model.PageInfo{TotalCount: 5, TotalPage: 1, PageNumber: 1},
			Corporations: []model.Hojin{{CorporateNumber: "111", Name: "更新法人"}},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL, "token")
	params := model.UpdateParams{From: "2024-01-01", To: "2024-01-31"}
	result, err := client.GetUpdateInfo(params)
	if err != nil {
		t.Fatalf("GetUpdateInfo() error: %v", err)
	}
	if result.TotalCount != 5 {
		t.Errorf("TotalCount = %d, want 5", result.TotalCount)
	}
}

func TestGetUpdateFinance(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/hojin/updateInfo/finance" {
			t.Errorf("path = %s", r.URL.Path)
		}
		resp := model.UpdateResponse{
			PageInfo: model.PageInfo{TotalCount: 1, TotalPage: 1, PageNumber: 1},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL, "token")
	result, err := client.GetUpdateFinance(model.UpdateParams{From: "2024-01-01", To: "2024-12-31"})
	if err != nil {
		t.Fatalf("GetUpdateFinance() error: %v", err)
	}
	if result.TotalCount != 1 {
		t.Errorf("TotalCount = %d, want 1", result.TotalCount)
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `cd /root/dev/planitai/gbizinfo-cli && go test ./internal/api/ -v -run TestGetUpdate`
Expected: FAIL

- [ ] **Step 3: Write implementation**

```go
// internal/api/update.go
package api

import (
	"net/url"
	"strconv"

	"github.com/planitaicojp/gbizinfo-cli/internal/model"
)

func (c *Client) getUpdate(path string, params model.UpdateParams) (*model.UpdateResponse, error) {
	q := url.Values{}
	if params.From != "" {
		q.Set("from", params.From)
	}
	if params.To != "" {
		q.Set("to", params.To)
	}
	if params.Page > 0 {
		q.Set("page", strconv.Itoa(params.Page))
	}
	if len(q) > 0 {
		path += "?" + q.Encode()
	}

	var result model.UpdateResponse
	if err := c.Get(path, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetUpdateInfo retrieves corporations updated within a period.
func (c *Client) GetUpdateInfo(params model.UpdateParams) (*model.UpdateResponse, error) {
	return c.getUpdate("/v1/hojin/updateInfo", params)
}

// GetUpdateCertification retrieves certifications updated within a period.
func (c *Client) GetUpdateCertification(params model.UpdateParams) (*model.UpdateResponse, error) {
	return c.getUpdate("/v1/hojin/updateInfo/certification", params)
}

// GetUpdateCommendation retrieves commendations updated within a period.
func (c *Client) GetUpdateCommendation(params model.UpdateParams) (*model.UpdateResponse, error) {
	return c.getUpdate("/v1/hojin/updateInfo/commendation", params)
}

// GetUpdateFinance retrieves financial data updated within a period.
func (c *Client) GetUpdateFinance(params model.UpdateParams) (*model.UpdateResponse, error) {
	return c.getUpdate("/v1/hojin/updateInfo/finance", params)
}

// GetUpdatePatent retrieves patents updated within a period.
func (c *Client) GetUpdatePatent(params model.UpdateParams) (*model.UpdateResponse, error) {
	return c.getUpdate("/v1/hojin/updateInfo/patent", params)
}

// GetUpdateProcurement retrieves procurement data updated within a period.
func (c *Client) GetUpdateProcurement(params model.UpdateParams) (*model.UpdateResponse, error) {
	return c.getUpdate("/v1/hojin/updateInfo/procurement", params)
}

// GetUpdateSubsidy retrieves subsidies updated within a period.
func (c *Client) GetUpdateSubsidy(params model.UpdateParams) (*model.UpdateResponse, error) {
	return c.getUpdate("/v1/hojin/updateInfo/subsidy", params)
}

// GetUpdateWorkplace retrieves workplace data updated within a period.
func (c *Client) GetUpdateWorkplace(params model.UpdateParams) (*model.UpdateResponse, error) {
	return c.getUpdate("/v1/hojin/updateInfo/workplace", params)
}
```

- [ ] **Step 4: Run tests to verify they pass**

Run: `cd /root/dev/planitai/gbizinfo-cli && go test ./internal/api/ -v`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add internal/api/update.go internal/api/update_test.go
git commit -m "feat: add API methods for updateInfo endpoints"
```

---

### Task 9: Command Utilities (cmdutil)

**Files:**
- Create: `cmd/cmdutil/cmdutil.go`

- [ ] **Step 1: Write cmdutil.go**

```go
// cmd/cmdutil/cmdutil.go
package cmdutil

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/planitaicojp/gbizinfo-cli/internal/api"
	"github.com/planitaicojp/gbizinfo-cli/internal/config"
	cerrors "github.com/planitaicojp/gbizinfo-cli/internal/errors"
)

const defaultBaseURL = "https://info.gbiz.go.jp/hojin"

// NewClient creates an API client from the cobra command context.
func NewClient(cmd *cobra.Command) (*api.Client, error) {
	token, _ := cmd.Flags().GetString("token")
	if token == "" {
		token = config.EnvOr(config.EnvToken, "")
	}
	if token == "" {
		cfg, err := config.Load()
		if err != nil {
			return nil, err
		}
		token = cfg.Token
	}
	if token == "" {
		return nil, &cerrors.AuthError{Message: "APIトークンが設定されていません。gbizinfo config init を実行してください"}
	}

	verbose, _ := cmd.Flags().GetBool("verbose")
	client := api.NewClient(defaultBaseURL, token)
	client.Verbose = verbose
	return client, nil
}

// GetFormat returns the output format from flags, env, or config.
func GetFormat(cmd *cobra.Command) string {
	format, _ := cmd.Flags().GetString("format")
	if format != "" {
		return format
	}
	if f := config.EnvOr(config.EnvFormat, ""); f != "" {
		return f
	}
	cfg, err := config.Load()
	if err != nil {
		return config.DefaultFormat
	}
	if cfg.Defaults.Format != "" {
		return cfg.Defaults.Format
	}
	return config.DefaultFormat
}

// ExactArgs returns a PositionalArgs that reports usage on mismatch.
func ExactArgs(n int) cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		if len(args) != n {
			return fmt.Errorf("%d個の引数が必要です（%d個指定されました）\n\n使い方:\n  %s", n, len(args), cmd.UseLine())
		}
		return nil
	}
}

// CorporateNumberArg extracts the corporate number from args or --corporate-number flag.
func CorporateNumberArg(cmd *cobra.Command, args []string) (string, error) {
	if len(args) > 0 {
		return args[0], nil
	}
	cn, _ := cmd.Flags().GetString("corporate-number")
	if cn != "" {
		return cn, nil
	}
	return "", &cerrors.ValidationError{Field: "corporate_number", Message: "法人番号を指定してください"}
}
```

- [ ] **Step 2: Verify compilation**

Run: `cd /root/dev/planitai/gbizinfo-cli && go get github.com/spf13/cobra && go build ./cmd/cmdutil/`
Expected: No errors

- [ ] **Step 3: Commit**

```bash
git add cmd/cmdutil/ go.mod go.sum
git commit -m "feat: add command utilities (client creation, format, args)"
```

---

### Task 10: Root, Version, Completion Commands

**Files:**
- Create: `cmd/root.go`
- Create: `cmd/version.go`
- Create: `cmd/completion.go`

- [ ] **Step 1: Write root.go**

```go
// cmd/root.go
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	cmdconfig "github.com/planitaicojp/gbizinfo-cli/cmd/config"
	"github.com/planitaicojp/gbizinfo-cli/cmd/update"
	cerrors "github.com/planitaicojp/gbizinfo-cli/internal/errors"
)

var version = "dev"

var rootCmd = &cobra.Command{
	Use:           "gbizinfo",
	Short:         "gBizINFO REST API CLI",
	Long:          "gBizINFO（経済産業省 法人情報API）を操作するCLIツール",
	SilenceUsage:  true,
	SilenceErrors: true,
}

func init() {
	rootCmd.PersistentFlags().StringP("format", "f", "", "出力形式 (json/table/csv)")
	rootCmd.PersistentFlags().StringP("token", "t", "", "APIトークン")
	rootCmd.PersistentFlags().Bool("no-color", false, "色出力を無効化")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "詳細出力")

	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(completionCmd)
	rootCmd.AddCommand(cmdconfig.Cmd)
	rootCmd.AddCommand(update.Cmd)
	rootCmd.AddCommand(searchCmd)
	rootCmd.AddCommand(getCmd)
	rootCmd.AddCommand(certificationCmd)
	rootCmd.AddCommand(commendationCmd)
	rootCmd.AddCommand(financeCmd)
	rootCmd.AddCommand(patentCmd)
	rootCmd.AddCommand(procurementCmd)
	rootCmd.AddCommand(subsidyCmd)
	rootCmd.AddCommand(workplaceCmd)
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(cerrors.GetExitCode(err))
	}
}
```

- [ ] **Step 2: Write version.go**

```go
// cmd/version.go
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "バージョン情報を表示",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("gbizinfo version %s\n", version)
	},
}
```

- [ ] **Step 3: Write completion.go**

```go
// cmd/completion.go
package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/planitaicojp/gbizinfo-cli/cmd/cmdutil"
)

var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish|powershell]",
	Short: "シェル補完スクリプトを生成",
	Long: `シェル補完スクリプトを生成します。

使用例:
  # Bash
  gbizinfo completion bash > /etc/bash_completion.d/gbizinfo

  # Zsh
  gbizinfo completion zsh > "${fpath[1]}/_gbizinfo"

  # Fish
  gbizinfo completion fish > ~/.config/fish/completions/gbizinfo.fish`,
	ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
	Args:                  cmdutil.ExactArgs(1),
	DisableFlagsInUseLine: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		switch args[0] {
		case "bash":
			return rootCmd.GenBashCompletion(os.Stdout)
		case "zsh":
			return rootCmd.GenZshCompletion(os.Stdout)
		case "fish":
			return rootCmd.GenFishCompletion(os.Stdout, true)
		case "powershell":
			return rootCmd.GenPowerShellCompletionWithDesc(os.Stdout)
		default:
			return nil
		}
	},
}
```

Note: These files reference commands (searchCmd, getCmd, etc.) that don't exist yet. They will be created in the next tasks. The project won't compile until all commands are added.

- [ ] **Step 4: Commit**

```bash
git add cmd/root.go cmd/version.go cmd/completion.go
git commit -m "feat: add root, version, and completion commands"
```

---

### Task 11: Config Subcommands

**Files:**
- Create: `cmd/config/config.go`

- [ ] **Step 1: Write config.go**

```go
// cmd/config/config.go
package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/planitaicojp/gbizinfo-cli/cmd/cmdutil"
	iconfig "github.com/planitaicojp/gbizinfo-cli/internal/config"
)

// Cmd is the config command group.
var Cmd = &cobra.Command{
	Use:   "config",
	Short: "設定を管理",
}

func init() {
	Cmd.AddCommand(initCmd)
	Cmd.AddCommand(setCmd)
	Cmd.AddCommand(showCmd)
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "初期設定を行う",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Fprint(os.Stderr, "gBizINFO APIトークンを入力してください: ")
		reader := bufio.NewReader(os.Stdin)
		token, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("入力の読み取りに失敗: %w", err)
		}
		token = strings.TrimSpace(token)
		if token == "" {
			return fmt.Errorf("トークンが入力されませんでした")
		}

		cfg, err := iconfig.Load()
		if err != nil {
			return err
		}
		cfg.Token = token
		if err := cfg.Save(); err != nil {
			return err
		}

		fmt.Fprintf(os.Stderr, "設定を保存しました: %s/config.yaml\n", iconfig.DefaultConfigDir())
		return nil
	},
}

var setCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "設定値を変更",
	Long:  "設定値を変更します。キー: token, format",
	Args:  cmdutil.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		key, value := args[0], args[1]

		cfg, err := iconfig.Load()
		if err != nil {
			return err
		}

		switch key {
		case "token":
			cfg.Token = value
		case "format":
			cfg.Defaults.Format = value
		default:
			return fmt.Errorf("不明な設定キー: %s (使用可能: token, format)", key)
		}

		return cfg.Save()
	},
}

var showCmd = &cobra.Command{
	Use:   "show",
	Short: "現在の設定を表示",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := iconfig.Load()
		if err != nil {
			return err
		}

		fmt.Printf("設定ディレクトリ: %s\n", iconfig.DefaultConfigDir())
		fmt.Printf("トークン:         %s\n", iconfig.MaskToken(cfg.Token))
		fmt.Printf("出力形式:         %s\n", cfg.Defaults.Format)
		return nil
	},
}
```

- [ ] **Step 2: Verify compilation**

Run: `cd /root/dev/planitai/gbizinfo-cli && go build ./cmd/config/`
Expected: No errors

- [ ] **Step 3: Commit**

```bash
git add cmd/config/
git commit -m "feat: add config init/set/show subcommands"
```

---

### Task 12: Search Command

**Files:**
- Create: `cmd/search.go`

- [ ] **Step 1: Write search.go**

```go
// cmd/search.go
package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/planitaicojp/gbizinfo-cli/cmd/cmdutil"
	"github.com/planitaicojp/gbizinfo-cli/internal/model"
	"github.com/planitaicojp/gbizinfo-cli/internal/output"
)

var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "法人を検索",
	Long:  "gBizINFOに登録された法人を検索します。",
	Example: `  gbizinfo search --name トヨタ
  gbizinfo search --name トヨタ --address 愛知県
  gbizinfo search -c 1234567890123`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := cmdutil.NewClient(cmd)
		if err != nil {
			return err
		}

		name, _ := cmd.Flags().GetString("name")
		address, _ := cmd.Flags().GetString("address")
		cn, _ := cmd.Flags().GetString("corporate-number")
		page, _ := cmd.Flags().GetInt("page")
		limit, _ := cmd.Flags().GetInt("limit")

		result, err := client.Search(model.SearchParams{
			Name:            name,
			Address:         address,
			CorporateNumber: cn,
			Page:            page,
			Limit:           limit,
		})
		if err != nil {
			return err
		}

		format := cmdutil.GetFormat(cmd)
		return output.New(format).Format(os.Stdout, result.Corporations)
	},
}

func init() {
	searchCmd.Flags().StringP("name", "n", "", "法人名")
	searchCmd.Flags().String("address", "", "所在地")
	searchCmd.Flags().StringP("corporate-number", "c", "", "法人番号")
	searchCmd.Flags().IntP("page", "p", 1, "ページ番号")
	searchCmd.Flags().IntP("limit", "l", 0, "表示件数")
}
```

- [ ] **Step 2: Commit**

```bash
git add cmd/search.go
git commit -m "feat: add search command"
```

---

### Task 13: Get Command

**Files:**
- Create: `cmd/get.go`

- [ ] **Step 1: Write get.go**

```go
// cmd/get.go
package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/planitaicojp/gbizinfo-cli/cmd/cmdutil"
	"github.com/planitaicojp/gbizinfo-cli/internal/output"
)

var getCmd = &cobra.Command{
	Use:   "get [法人番号]",
	Short: "法人基本情報を取得",
	Long:  "指定した法人番号の基本情報を取得します。",
	Example: `  gbizinfo get 1234567890123
  gbizinfo get -c 1234567890123 -f table`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := cmdutil.NewClient(cmd)
		if err != nil {
			return err
		}

		cn, err := cmdutil.CorporateNumberArg(cmd, args)
		if err != nil {
			return err
		}

		result, err := client.GetHojin(cn)
		if err != nil {
			return err
		}

		format := cmdutil.GetFormat(cmd)
		return output.New(format).Format(os.Stdout, result.Corporations)
	},
}

func init() {
	getCmd.Flags().StringP("corporate-number", "c", "", "法人番号")
}
```

- [ ] **Step 2: Commit**

```bash
git add cmd/get.go
git commit -m "feat: add get command"
```

---

### Task 14: Resource Detail Commands (certification, commendation, finance, patent, procurement, subsidy, workplace)

**Files:**
- Create: `cmd/certification.go`
- Create: `cmd/commendation.go`
- Create: `cmd/finance.go`
- Create: `cmd/patent.go`
- Create: `cmd/procurement.go`
- Create: `cmd/subsidy.go`
- Create: `cmd/workplace.go`

- [ ] **Step 1: Write certification.go**

```go
// cmd/certification.go
package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/planitaicojp/gbizinfo-cli/cmd/cmdutil"
	"github.com/planitaicojp/gbizinfo-cli/internal/output"
)

var certificationCmd = &cobra.Command{
	Use:     "certification [法人番号]",
	Short:   "届出・認定情報を取得",
	Long:    "指定した法人番号の届出・認定情報を取得します。",
	Example: "  gbizinfo certification 1234567890123",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := cmdutil.NewClient(cmd)
		if err != nil {
			return err
		}
		cn, err := cmdutil.CorporateNumberArg(cmd, args)
		if err != nil {
			return err
		}
		result, err := client.GetCertification(cn)
		if err != nil {
			return err
		}
		format := cmdutil.GetFormat(cmd)
		return output.New(format).Format(os.Stdout, result.Corporations)
	},
}

func init() {
	certificationCmd.Flags().StringP("corporate-number", "c", "", "法人番号")
}
```

- [ ] **Step 2: Write commendation.go**

```go
// cmd/commendation.go
package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/planitaicojp/gbizinfo-cli/cmd/cmdutil"
	"github.com/planitaicojp/gbizinfo-cli/internal/output"
)

var commendationCmd = &cobra.Command{
	Use:     "commendation [法人番号]",
	Short:   "表彰情報を取得",
	Long:    "指定した法人番号の表彰情報を取得します。",
	Example: "  gbizinfo commendation 1234567890123",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := cmdutil.NewClient(cmd)
		if err != nil {
			return err
		}
		cn, err := cmdutil.CorporateNumberArg(cmd, args)
		if err != nil {
			return err
		}
		result, err := client.GetCommendation(cn)
		if err != nil {
			return err
		}
		format := cmdutil.GetFormat(cmd)
		return output.New(format).Format(os.Stdout, result.Corporations)
	},
}

func init() {
	commendationCmd.Flags().StringP("corporate-number", "c", "", "法人番号")
}
```

- [ ] **Step 3: Write finance.go**

```go
// cmd/finance.go
package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/planitaicojp/gbizinfo-cli/cmd/cmdutil"
	"github.com/planitaicojp/gbizinfo-cli/internal/output"
)

var financeCmd = &cobra.Command{
	Use:     "finance [法人番号]",
	Short:   "財務情報を取得",
	Long:    "指定した法人番号の財務情報を取得します。",
	Example: "  gbizinfo finance 1234567890123",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := cmdutil.NewClient(cmd)
		if err != nil {
			return err
		}
		cn, err := cmdutil.CorporateNumberArg(cmd, args)
		if err != nil {
			return err
		}
		result, err := client.GetFinance(cn)
		if err != nil {
			return err
		}
		format := cmdutil.GetFormat(cmd)
		return output.New(format).Format(os.Stdout, result.Corporations)
	},
}

func init() {
	financeCmd.Flags().StringP("corporate-number", "c", "", "法人番号")
}
```

- [ ] **Step 4: Write patent.go**

```go
// cmd/patent.go
package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/planitaicojp/gbizinfo-cli/cmd/cmdutil"
	"github.com/planitaicojp/gbizinfo-cli/internal/output"
)

var patentCmd = &cobra.Command{
	Use:     "patent [法人番号]",
	Short:   "特許情報を取得",
	Long:    "指定した法人番号の特許情報を取得します。",
	Example: "  gbizinfo patent 1234567890123",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := cmdutil.NewClient(cmd)
		if err != nil {
			return err
		}
		cn, err := cmdutil.CorporateNumberArg(cmd, args)
		if err != nil {
			return err
		}
		result, err := client.GetPatent(cn)
		if err != nil {
			return err
		}
		format := cmdutil.GetFormat(cmd)
		return output.New(format).Format(os.Stdout, result.Corporations)
	},
}

func init() {
	patentCmd.Flags().StringP("corporate-number", "c", "", "法人番号")
}
```

- [ ] **Step 5: Write procurement.go**

```go
// cmd/procurement.go
package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/planitaicojp/gbizinfo-cli/cmd/cmdutil"
	"github.com/planitaicojp/gbizinfo-cli/internal/output"
)

var procurementCmd = &cobra.Command{
	Use:     "procurement [法人番号]",
	Short:   "調達情報を取得",
	Long:    "指定した法人番号の調達情報（入札・契約実績）を取得します。",
	Example: "  gbizinfo procurement 1234567890123",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := cmdutil.NewClient(cmd)
		if err != nil {
			return err
		}
		cn, err := cmdutil.CorporateNumberArg(cmd, args)
		if err != nil {
			return err
		}
		result, err := client.GetProcurement(cn)
		if err != nil {
			return err
		}
		format := cmdutil.GetFormat(cmd)
		return output.New(format).Format(os.Stdout, result.Corporations)
	},
}

func init() {
	procurementCmd.Flags().StringP("corporate-number", "c", "", "法人番号")
}
```

- [ ] **Step 6: Write subsidy.go**

```go
// cmd/subsidy.go
package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/planitaicojp/gbizinfo-cli/cmd/cmdutil"
	"github.com/planitaicojp/gbizinfo-cli/internal/output"
)

var subsidyCmd = &cobra.Command{
	Use:     "subsidy [法人番号]",
	Short:   "補助金情報を取得",
	Long:    "指定した法人番号の補助金情報を取得します。",
	Example: "  gbizinfo subsidy 1234567890123",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := cmdutil.NewClient(cmd)
		if err != nil {
			return err
		}
		cn, err := cmdutil.CorporateNumberArg(cmd, args)
		if err != nil {
			return err
		}
		result, err := client.GetSubsidy(cn)
		if err != nil {
			return err
		}
		format := cmdutil.GetFormat(cmd)
		return output.New(format).Format(os.Stdout, result.Corporations)
	},
}

func init() {
	subsidyCmd.Flags().StringP("corporate-number", "c", "", "法人番号")
}
```

- [ ] **Step 7: Write workplace.go**

```go
// cmd/workplace.go
package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/planitaicojp/gbizinfo-cli/cmd/cmdutil"
	"github.com/planitaicojp/gbizinfo-cli/internal/output"
)

var workplaceCmd = &cobra.Command{
	Use:     "workplace [法人番号]",
	Short:   "職場情報を取得",
	Long:    "指定した法人番号の職場情報（従業員数、女性比率等）を取得します。",
	Example: "  gbizinfo workplace 1234567890123",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := cmdutil.NewClient(cmd)
		if err != nil {
			return err
		}
		cn, err := cmdutil.CorporateNumberArg(cmd, args)
		if err != nil {
			return err
		}
		result, err := client.GetWorkplace(cn)
		if err != nil {
			return err
		}
		format := cmdutil.GetFormat(cmd)
		return output.New(format).Format(os.Stdout, result.Corporations)
	},
}

func init() {
	workplaceCmd.Flags().StringP("corporate-number", "c", "", "法人番号")
}
```

- [ ] **Step 8: Commit**

```bash
git add cmd/certification.go cmd/commendation.go cmd/finance.go cmd/patent.go cmd/procurement.go cmd/subsidy.go cmd/workplace.go
git commit -m "feat: add resource detail commands (certification through workplace)"
```

---

### Task 15: Update Command Group

**Files:**
- Create: `cmd/update/update.go`

- [ ] **Step 1: Write update.go**

```go
// cmd/update/update.go
package update

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/planitaicojp/gbizinfo-cli/cmd/cmdutil"
	"github.com/planitaicojp/gbizinfo-cli/internal/model"
	"github.com/planitaicojp/gbizinfo-cli/internal/output"
)

// Cmd is the update command group.
var Cmd = &cobra.Command{
	Use:   "update",
	Short: "期間指定で更新情報を取得",
	Long:  "指定した期間内に更新された法人情報を取得します。",
}

func init() {
	Cmd.AddCommand(hojinCmd)
	Cmd.AddCommand(certificationCmd)
	Cmd.AddCommand(commendationCmd)
	Cmd.AddCommand(financeCmd)
	Cmd.AddCommand(patentCmd)
	Cmd.AddCommand(procurementCmd)
	Cmd.AddCommand(subsidyCmd)
	Cmd.AddCommand(workplaceCmd)
}

func updateParams(cmd *cobra.Command) model.UpdateParams {
	from, _ := cmd.Flags().GetString("from")
	to, _ := cmd.Flags().GetString("to")
	page, _ := cmd.Flags().GetInt("page")
	return model.UpdateParams{From: from, To: to, Page: page}
}

func addUpdateFlags(cmd *cobra.Command) {
	cmd.Flags().String("from", "", "開始日 (YYYY-MM-DD)")
	cmd.Flags().String("to", "", "終了日 (YYYY-MM-DD)")
	cmd.Flags().IntP("page", "p", 1, "ページ番号")
}

var hojinCmd = &cobra.Command{
	Use:     "hojin",
	Short:   "期間内に更新された法人情報を取得",
	Example: "  gbizinfo update hojin --from 2024-01-01 --to 2024-01-31",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := cmdutil.NewClient(cmd)
		if err != nil {
			return err
		}
		result, err := client.GetUpdateInfo(updateParams(cmd))
		if err != nil {
			return err
		}
		format := cmdutil.GetFormat(cmd)
		return output.New(format).Format(os.Stdout, result.Corporations)
	},
}

var certificationCmd = &cobra.Command{
	Use:     "certification",
	Short:   "期間内に更新された認定情報を取得",
	Example: "  gbizinfo update certification --from 2024-01-01 --to 2024-01-31",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := cmdutil.NewClient(cmd)
		if err != nil {
			return err
		}
		result, err := client.GetUpdateCertification(updateParams(cmd))
		if err != nil {
			return err
		}
		format := cmdutil.GetFormat(cmd)
		return output.New(format).Format(os.Stdout, result.Corporations)
	},
}

var commendationCmd = &cobra.Command{
	Use:     "commendation",
	Short:   "期間内に更新された表彰情報を取得",
	Example: "  gbizinfo update commendation --from 2024-01-01 --to 2024-01-31",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := cmdutil.NewClient(cmd)
		if err != nil {
			return err
		}
		result, err := client.GetUpdateCommendation(updateParams(cmd))
		if err != nil {
			return err
		}
		format := cmdutil.GetFormat(cmd)
		return output.New(format).Format(os.Stdout, result.Corporations)
	},
}

var financeCmd = &cobra.Command{
	Use:     "finance",
	Short:   "期間内に更新された財務情報を取得",
	Example: "  gbizinfo update finance --from 2024-01-01 --to 2024-01-31",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := cmdutil.NewClient(cmd)
		if err != nil {
			return err
		}
		result, err := client.GetUpdateFinance(updateParams(cmd))
		if err != nil {
			return err
		}
		format := cmdutil.GetFormat(cmd)
		return output.New(format).Format(os.Stdout, result.Corporations)
	},
}

var patentCmd = &cobra.Command{
	Use:     "patent",
	Short:   "期間内に更新された特許情報を取得",
	Example: "  gbizinfo update patent --from 2024-01-01 --to 2024-01-31",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := cmdutil.NewClient(cmd)
		if err != nil {
			return err
		}
		result, err := client.GetUpdatePatent(updateParams(cmd))
		if err != nil {
			return err
		}
		format := cmdutil.GetFormat(cmd)
		return output.New(format).Format(os.Stdout, result.Corporations)
	},
}

var procurementCmd = &cobra.Command{
	Use:     "procurement",
	Short:   "期間内に更新された調達情報を取得",
	Example: "  gbizinfo update procurement --from 2024-01-01 --to 2024-01-31",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := cmdutil.NewClient(cmd)
		if err != nil {
			return err
		}
		result, err := client.GetUpdateProcurement(updateParams(cmd))
		if err != nil {
			return err
		}
		format := cmdutil.GetFormat(cmd)
		return output.New(format).Format(os.Stdout, result.Corporations)
	},
}

var subsidyCmd = &cobra.Command{
	Use:     "subsidy",
	Short:   "期間内に更新された補助金情報を取得",
	Example: "  gbizinfo update subsidy --from 2024-01-01 --to 2024-01-31",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := cmdutil.NewClient(cmd)
		if err != nil {
			return err
		}
		result, err := client.GetUpdateSubsidy(updateParams(cmd))
		if err != nil {
			return err
		}
		format := cmdutil.GetFormat(cmd)
		return output.New(format).Format(os.Stdout, result.Corporations)
	},
}

var workplaceCmd = &cobra.Command{
	Use:     "workplace",
	Short:   "期間内に更新された職場情報を取得",
	Example: "  gbizinfo update workplace --from 2024-01-01 --to 2024-01-31",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := cmdutil.NewClient(cmd)
		if err != nil {
			return err
		}
		result, err := client.GetUpdateWorkplace(updateParams(cmd))
		if err != nil {
			return err
		}
		format := cmdutil.GetFormat(cmd)
		return output.New(format).Format(os.Stdout, result.Corporations)
	},
}

func init() {
	for _, cmd := range []*cobra.Command{
		hojinCmd, certificationCmd, commendationCmd, financeCmd,
		patentCmd, procurementCmd, subsidyCmd, workplaceCmd,
	} {
		addUpdateFlags(cmd)
	}
}
```

- [ ] **Step 2: Commit**

```bash
git add cmd/update/
git commit -m "feat: add update command group with 8 subcommands"
```

---

### Task 16: Build and Smoke Test

**Files:**
- Modify: all (compilation check)

- [ ] **Step 1: Build the binary**

Run: `cd /root/dev/planitai/gbizinfo-cli && go mod tidy && make build`
Expected: Binary `gbizinfo` created with no errors

- [ ] **Step 2: Verify version command**

Run: `./gbizinfo version`
Expected: `gbizinfo version dev`

- [ ] **Step 3: Verify help output**

Run: `./gbizinfo --help`
Expected: Shows all subcommands (search, get, certification, commendation, finance, patent, procurement, subsidy, workplace, update, config, version, completion) with Japanese descriptions

- [ ] **Step 4: Verify subcommand help**

Run: `./gbizinfo search --help`
Expected: Shows search flags (--name, --address, --corporate-number, --page, --limit) with Japanese descriptions

- [ ] **Step 5: Verify update subcommand help**

Run: `./gbizinfo update --help`
Expected: Shows 8 subcommands (hojin, certification, commendation, finance, patent, procurement, subsidy, workplace)

- [ ] **Step 6: Verify config show without config file**

Run: `./gbizinfo config show`
Expected: Shows defaults (format: json, empty masked token)

- [ ] **Step 7: Run all tests**

Run: `cd /root/dev/planitai/gbizinfo-cli && make test`
Expected: All tests pass

- [ ] **Step 8: Commit go.sum and any tidied files**

```bash
git add go.mod go.sum
git commit -m "chore: tidy dependencies and verify build"
```

---

### Task 17: Verbose Debug Logging

**Files:**
- Modify: `internal/api/client.go`

- [ ] **Step 1: Replace debug stub functions in client.go**

Replace the placeholder `debugLogRequest` and `debugLogResponse` functions:

```go
func debugLogRequest(req *http.Request) {
	fmt.Fprintf(os.Stderr, "> %s %s\n", req.Method, req.URL.String())
	for key, vals := range req.Header {
		for _, v := range vals {
			if key == "X-Hojininfo-Api-Token" {
				v = v[:4] + "******"
			}
			fmt.Fprintf(os.Stderr, "> %s: %s\n", key, v)
		}
	}
	fmt.Fprintln(os.Stderr)
}

func debugLogResponse(resp *http.Response) {
	fmt.Fprintf(os.Stderr, "< %s\n", resp.Status)
	for key, vals := range resp.Header {
		for _, v := range vals {
			fmt.Fprintf(os.Stderr, "< %s: %s\n", key, v)
		}
	}
	fmt.Fprintln(os.Stderr)
}
```

Add `"os"` to the imports if not already present.

- [ ] **Step 2: Verify build**

Run: `cd /root/dev/planitai/gbizinfo-cli && make build`
Expected: No errors

- [ ] **Step 3: Commit**

```bash
git add internal/api/client.go
git commit -m "feat: add verbose debug logging for API requests"
```

---

### Task 18: GoReleaser + Homebrew Tap + Scoop Bucket

**Files:**
- Create: `.goreleaser.yaml`
- Create: `.github/workflows/release.yml`

- [ ] **Step 1: Create .goreleaser.yaml**

```yaml
version: 2

builds:
  - binary: gbizinfo
    ldflags:
      - -s -w -X github.com/planitaicojp/gbizinfo-cli/cmd.version={{.Version}}
    goos:
      - linux
      - darwin
      - windows
    goarch:
      - amd64
      - arm64

archives:
  - format: tar.gz
    name_template: "{{ .ProjectName }}_{{ .Version }}_{{ .Os }}_{{ .Arch }}"
    format_overrides:
      - goos: windows
        format: zip

checksum:
  name_template: "checksums.txt"

changelog:
  sort: asc
  filters:
    exclude:
      - "^docs:"
      - "^test:"
      - "^chore:"

brews:
  - repository:
      owner: planitaicojp
      name: homebrew-tap
      token: "{{ .Env.HOMEBREW_TAP_TOKEN }}"
    homepage: "https://github.com/planitaicojp/gbizinfo-cli"
    description: "gBizINFO REST API CLI - 経済産業省 法人情報API CLIツール"
    license: "MIT"
    install: |
      bin.install "gbizinfo"
    test: |
      system "#{bin}/gbizinfo", "version"

scoops:
  - repository:
      owner: planitaicojp
      name: scoop-bucket
      token: "{{ .Env.SCOOP_BUCKET_TOKEN }}"
    homepage: "https://github.com/planitaicojp/gbizinfo-cli"
    description: "gBizINFO REST API CLI - 経済産業省 法人情報API CLIツール"
    license: "MIT"
```

- [ ] **Step 2: Create GitHub Actions release workflow**

```yaml
# .github/workflows/release.yml
name: Release

on:
  push:
    tags:
      - "v*"

permissions:
  contents: write

jobs:
  goreleaser:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0

      - uses: actions/setup-go@v5
        with:
          go-version-file: go.mod

      - uses: goreleaser/goreleaser-action@v6
        with:
          version: "~> v2"
          args: release --clean
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
          HOMEBREW_TAP_TOKEN: ${{ secrets.HOMEBREW_TAP_TOKEN }}
          SCOOP_BUCKET_TOKEN: ${{ secrets.SCOOP_BUCKET_TOKEN }}
```

- [ ] **Step 3: Commit**

```bash
git add .goreleaser.yaml .github/workflows/release.yml
git commit -m "ci: add goreleaser config with Homebrew tap and Scoop bucket"
```

---

### Task 19: README

**Files:**
- Create: `README.md`

- [ ] **Step 1: Write README.md**

````markdown
# gbizinfo-cli

gBizINFO（経済産業省 法人情報API）を操作するCLIツール。

## インストール

### Homebrew (macOS / Linux)

```bash
brew install planitaicojp/tap/gbizinfo-cli
```

### Scoop (Windows)

```powershell
scoop bucket add planitaicojp https://github.com/planitaicojp/scoop-bucket.git
scoop install gbizinfo-cli
```

### Go

```bash
go install github.com/planitaicojp/gbizinfo-cli@latest
```

### リリースバイナリ

[GitHub Releases](https://github.com/planitaicojp/gbizinfo-cli/releases) からダウンロード:

```bash
# Linux (amd64)
curl -sL https://github.com/planitaicojp/gbizinfo-cli/releases/latest/download/gbizinfo-cli_Linux_amd64.tar.gz | tar xz
sudo mv gbizinfo /usr/local/bin/

# macOS (Apple Silicon)
curl -sL https://github.com/planitaicojp/gbizinfo-cli/releases/latest/download/gbizinfo-cli_Darwin_arm64.tar.gz | tar xz
sudo mv gbizinfo /usr/local/bin/

# Windows (amd64) — PowerShell
Invoke-WebRequest -Uri https://github.com/planitaicojp/gbizinfo-cli/releases/latest/download/gbizinfo-cli_Windows_amd64.zip -OutFile gbizinfo.zip
Expand-Archive gbizinfo.zip -DestinationPath .
```

## 初期設定

[gBizINFO](https://info.gbiz.go.jp/) でアカウント登録後、APIトークンを取得してください。

```bash
gbizinfo config init
```

環境変数でも設定できます:

```bash
export GBIZINFO_TOKEN=your-token-here
```

## 使い方

### 法人検索

```bash
gbizinfo search --name トヨタ
gbizinfo search --name トヨタ --address 愛知県
gbizinfo search -c 1234567890123
```

### 法人情報取得

```bash
gbizinfo get 1234567890123          # 基本情報
gbizinfo finance 1234567890123      # 財務情報
gbizinfo subsidy 1234567890123      # 補助金情報
gbizinfo patent 1234567890123       # 特許情報
gbizinfo procurement 1234567890123  # 調達情報
gbizinfo certification 1234567890123 # 届出・認定情報
gbizinfo commendation 1234567890123 # 表彰情報
gbizinfo workplace 1234567890123    # 職場情報
```

### 期間指定更新情報

```bash
gbizinfo update hojin --from 2024-01-01 --to 2024-01-31
gbizinfo update finance --from 2024-01-01 --to 2024-12-31
```

### 出力形式

```bash
gbizinfo search --name トヨタ -f table
gbizinfo search --name トヨタ -f csv
gbizinfo search --name トヨタ -f json   # デフォルト
```

## 設定

設定ファイル: `~/.config/gbizinfo/config.yaml`

```yaml
token: "your-api-token"
defaults:
  format: json
```

### 環境変数

| 変数名 | 説明 |
|--------|------|
| `GBIZINFO_TOKEN` | APIトークン |
| `GBIZINFO_FORMAT` | デフォルト出力形式 (json/table/csv) |
| `GBIZINFO_CONFIG_DIR` | 設定ディレクトリのパス |

### 優先順位

CLIフラグ > 環境変数 > 設定ファイル > デフォルト値

## 対象API

[gBizINFO REST API](https://info.gbiz.go.jp/api/index.html) — 経済産業省が提供する法人情報API

## ライセンス

MIT
````

- [ ] **Step 2: Commit**

```bash
git add README.md
git commit -m "docs: add README with installation and usage instructions"
```
