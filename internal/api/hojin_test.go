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
			t.Errorf("name param = %q", r.URL.Query().Get("name"))
		}
		if r.URL.Query().Get("page") != "1" {
			t.Errorf("page param = %q", r.URL.Query().Get("page"))
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
		t.Errorf("Name = %q", result.Corporations[0].Name)
	}
}

func TestGetHojin(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/hojin/1234567890123" {
			t.Errorf("path = %s", r.URL.Path)
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
		t.Errorf("CorporateNumber = %q", result.Corporations[0].CorporateNumber)
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
