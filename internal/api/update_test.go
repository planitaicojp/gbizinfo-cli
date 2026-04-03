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
