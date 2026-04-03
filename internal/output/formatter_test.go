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

// Nested struct types for testing flatten behavior
type nestedChild struct {
	Title string `json:"title"`
	Date  string `json:"date"`
}

type nestedParent struct {
	CorporateNumber string        `json:"corporate_number"`
	Name            string        `json:"name"`
	Children        []nestedChild `json:"children"`
}

func TestTableFormatterNested(t *testing.T) {
	items := []nestedParent{
		{
			CorporateNumber: "1234567890123",
			Name:            "テスト株式会社",
			Children: []nestedChild{
				{Title: "認定A", Date: "2024-01-01"},
				{Title: "認定B", Date: "2024-06-01"},
			},
		},
	}
	var buf bytes.Buffer
	f := New("table")
	if err := f.Format(&buf, items); err != nil {
		t.Fatalf("Format() error: %v", err)
	}
	out := buf.String()
	// Headers should include parent + child fields, NOT a "CHILDREN" column
	if strings.Contains(out, "CHILDREN") {
		t.Errorf("should not contain raw CHILDREN header, got:\n%s", out)
	}
	if !strings.Contains(out, "CORPORATE_NUMBER") {
		t.Errorf("missing CORPORATE_NUMBER header, got:\n%s", out)
	}
	if !strings.Contains(out, "TITLE") {
		t.Errorf("missing TITLE header, got:\n%s", out)
	}
	// Should have 2 data rows (one per child)
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != 3 { // header + 2 rows
		t.Errorf("expected 3 lines (header + 2 rows), got %d:\n%s", len(lines), out)
	}
	if !strings.Contains(out, "認定A") || !strings.Contains(out, "認定B") {
		t.Errorf("missing nested data, got:\n%s", out)
	}
}

func TestCSVFormatterNested(t *testing.T) {
	items := []nestedParent{
		{
			CorporateNumber: "1234567890123",
			Name:            "テスト株式会社",
			Children: []nestedChild{
				{Title: "補助金X", Date: "2024-03-15"},
			},
		},
		{
			CorporateNumber: "9876543210987",
			Name:            "サンプル合同会社",
			Children: []nestedChild{
				{Title: "補助金Y", Date: "2024-07-01"},
				{Title: "補助金Z", Date: "2024-12-01"},
			},
		},
	}
	var buf bytes.Buffer
	f := New("csv")
	if err := f.Format(&buf, items); err != nil {
		t.Fatalf("Format() error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	// header + 1 row (first corp) + 2 rows (second corp) = 4
	if len(lines) != 4 {
		t.Fatalf("expected 4 lines, got %d:\n%s", len(lines), buf.String())
	}
	if lines[0] != "corporate_number,name,title,date" {
		t.Errorf("header = %q, want %q", lines[0], "corporate_number,name,title,date")
	}
}

func TestTableFormatterNestedEmpty(t *testing.T) {
	items := []nestedParent{
		{
			CorporateNumber: "1234567890123",
			Name:            "テスト株式会社",
			Children:        nil,
		},
	}
	var buf bytes.Buffer
	f := New("table")
	if err := f.Format(&buf, items); err != nil {
		t.Fatalf("Format() error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "1234567890123") {
		t.Errorf("should still show parent fields when children empty, got:\n%s", out)
	}
}

// Test []*Struct slice fields (pointer elements)
type ptrChild struct {
	Label string `json:"label"`
}

type ptrParent struct {
	ID       string      `json:"id"`
	Children []*ptrChild `json:"children"`
}

func TestTableFormatterPointerSlice(t *testing.T) {
	items := []ptrParent{
		{
			ID: "001",
			Children: []*ptrChild{
				{Label: "child-a"},
				{Label: "child-b"},
			},
		},
	}
	var buf bytes.Buffer
	f := New("table")
	if err := f.Format(&buf, items); err != nil {
		t.Fatalf("Format() error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "LABEL") {
		t.Errorf("should expand []*Struct fields, got:\n%s", out)
	}
	if !strings.Contains(out, "child-a") || !strings.Contains(out, "child-b") {
		t.Errorf("missing pointer child data, got:\n%s", out)
	}
}

// Test []string fields are joined with "; "
type taggedItem struct {
	Name string   `json:"name"`
	Tags []string `json:"tags"`
}

func TestTableFormatterStringSlice(t *testing.T) {
	items := []taggedItem{
		{Name: "item1", Tags: []string{"go", "cli", "api"}},
	}
	var buf bytes.Buffer
	f := New("table")
	if err := f.Format(&buf, items); err != nil {
		t.Fatalf("Format() error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "go; cli; api") {
		t.Errorf("[]string should be joined with '; ', got:\n%s", out)
	}
}

// Test struct with multiple []struct fields: only first is expanded
type multiSliceParent struct {
	ID     string        `json:"id"`
	First  []nestedChild `json:"first"`
	Second []nestedChild `json:"second"`
}

func TestTableFormatterMultipleSliceFields(t *testing.T) {
	items := []multiSliceParent{
		{
			ID:     "001",
			First:  []nestedChild{{Title: "A", Date: "2024-01-01"}},
			Second: []nestedChild{{Title: "B", Date: "2024-06-01"}},
		},
	}
	var buf bytes.Buffer
	f := New("table")
	if err := f.Format(&buf, items); err != nil {
		t.Fatalf("Format() error: %v", err)
	}
	out := buf.String()
	// First slice should be expanded
	if !strings.Contains(out, "TITLE") {
		t.Errorf("first []struct should be expanded, got:\n%s", out)
	}
	if !strings.Contains(out, "A") {
		t.Errorf("first slice data missing, got:\n%s", out)
	}
	// Second slice should appear as a scalar column (SECOND header)
	if !strings.Contains(out, "SECOND") {
		t.Errorf("second []struct should appear as scalar column, got:\n%s", out)
	}
}
