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
