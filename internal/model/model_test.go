package model_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/planitaicojp/gbizinfo-cli/internal/model"
)

func fixturesDir() string {
	_, f, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(f), "..", "..", "test", "fixtures")
}

func loadFixture(t *testing.T, name string) []byte {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(fixturesDir(), name))
	if err != nil {
		t.Fatalf("fixture %s を読み込めません: %v", name, err)
	}
	return data
}

func TestSearchResponseDeserialization(t *testing.T) {
	var resp model.HojinResponse
	if err := json.Unmarshal(loadFixture(t, "search.json"), &resp); err != nil {
		t.Fatalf("デシリアライズ失敗: %v", err)
	}
	if resp.TotalCount != 2 {
		t.Errorf("TotalCount = %d, want 2", resp.TotalCount)
	}
	if resp.TotalPage != 1 {
		t.Errorf("TotalPage = %d, want 1", resp.TotalPage)
	}
	if len(resp.Corporations) != 2 {
		t.Fatalf("Corporations count = %d, want 2", len(resp.Corporations))
	}
	h := resp.Corporations[0]
	if h.CorporateNumber != "1010001012345" {
		t.Errorf("CorporateNumber = %q, want %q", h.CorporateNumber, "1010001012345")
	}
	if h.Name != "株式会社テスト商事" {
		t.Errorf("Name = %q, want %q", h.Name, "株式会社テスト商事")
	}
	if h.NameKana != "カブシキガイシャテストショウジ" {
		t.Errorf("NameKana = %q, want %q", h.NameKana, "カブシキガイシャテストショウジ")
	}
	if h.Capital != "10000000" {
		t.Errorf("Capital = %q, want %q", h.Capital, "10000000")
	}
}

func TestHojinDetailDeserialization(t *testing.T) {
	var resp model.HojinDetail
	if err := json.Unmarshal(loadFixture(t, "hojin_detail.json"), &resp); err != nil {
		t.Fatalf("デシリアライズ失敗: %v", err)
	}
	if len(resp.Corporations) != 1 {
		t.Fatalf("Corporations count = %d, want 1", len(resp.Corporations))
	}
	h := resp.Corporations[0]
	if h.RepresentName != "山田太郎" {
		t.Errorf("RepresentName = %q, want %q", h.RepresentName, "山田太郎")
	}
	if h.DateOfEstablish != "2010-01-15" {
		t.Errorf("DateOfEstablish = %q, want %q", h.DateOfEstablish, "2010-01-15")
	}
}

func TestCertificationDeserialization(t *testing.T) {
	var resp model.CertificationResponse
	if err := json.Unmarshal(loadFixture(t, "certification.json"), &resp); err != nil {
		t.Fatalf("デシリアライズ失敗: %v", err)
	}
	if len(resp.Corporations) != 1 {
		t.Fatalf("Corporations count = %d, want 1", len(resp.Corporations))
	}
	info := resp.Corporations[0]
	if info.CorporateNumber != "1010001012345" {
		t.Errorf("CorporateNumber = %q", info.CorporateNumber)
	}
	if len(info.Certifications) != 2 {
		t.Fatalf("Certifications count = %d, want 2", len(info.Certifications))
	}
	c := info.Certifications[0]
	if c.Title != "えるぼし認定（3段階目）" {
		t.Errorf("Title = %q", c.Title)
	}
	if c.GovernmentDepartments != "厚生労働省" {
		t.Errorf("GovernmentDepartments = %q", c.GovernmentDepartments)
	}
}

func TestCommendationDeserialization(t *testing.T) {
	var resp model.CommendationResponse
	if err := json.Unmarshal(loadFixture(t, "commendation.json"), &resp); err != nil {
		t.Fatalf("デシリアライズ失敗: %v", err)
	}
	info := resp.Corporations[0]
	if len(info.Commendations) != 1 {
		t.Fatalf("Commendations count = %d, want 1", len(info.Commendations))
	}
	c := info.Commendations[0]
	if c.Title != "健康経営優良法人2024" {
		t.Errorf("Title = %q", c.Title)
	}
	if c.DateOfCommendation != "2024-03-11" {
		t.Errorf("DateOfCommendation = %q", c.DateOfCommendation)
	}
}

func TestFinanceDeserialization(t *testing.T) {
	var resp model.FinanceResponse
	if err := json.Unmarshal(loadFixture(t, "finance.json"), &resp); err != nil {
		t.Fatalf("デシリアライズ失敗: %v", err)
	}
	info := resp.Corporations[0]
	if len(info.Finance) != 2 {
		t.Fatalf("Finance count = %d, want 2", len(info.Finance))
	}
	f := info.Finance[0]
	if f.AccountingPeriod != "2023年度" {
		t.Errorf("AccountingPeriod = %q", f.AccountingPeriod)
	}
	if f.NetSales != "500000000" {
		t.Errorf("NetSales = %q", f.NetSales)
	}
	if f.Profit != "30000000" {
		t.Errorf("Profit = %q", f.Profit)
	}
}

func TestPatentDeserialization(t *testing.T) {
	var resp model.PatentResponse
	if err := json.Unmarshal(loadFixture(t, "patent.json"), &resp); err != nil {
		t.Fatalf("デシリアライズ失敗: %v", err)
	}
	info := resp.Corporations[0]
	if len(info.Patents) != 1 {
		t.Fatalf("Patents count = %d, want 1", len(info.Patents))
	}
	p := info.Patents[0]
	if p.PatentNumber != "特許第7000001号" {
		t.Errorf("PatentNumber = %q", p.PatentNumber)
	}
	if p.ClassificationPI != "G06F 16/00" {
		t.Errorf("ClassificationPI = %q", p.ClassificationPI)
	}
}

func TestProcurementDeserialization(t *testing.T) {
	var resp model.ProcurementResponse
	if err := json.Unmarshal(loadFixture(t, "procurement.json"), &resp); err != nil {
		t.Fatalf("デシリアライズ失敗: %v", err)
	}
	info := resp.Corporations[0]
	if len(info.Procurements) != 2 {
		t.Fatalf("Procurements count = %d, want 2", len(info.Procurements))
	}
	p := info.Procurements[0]
	if p.Amount != "15000000" {
		t.Errorf("Amount = %q", p.Amount)
	}
	if p.GovernmentDepartments != "デジタル庁" {
		t.Errorf("GovernmentDepartments = %q", p.GovernmentDepartments)
	}
}

func TestSubsidyDeserialization(t *testing.T) {
	var resp model.SubsidyResponse
	if err := json.Unmarshal(loadFixture(t, "subsidy.json"), &resp); err != nil {
		t.Fatalf("デシリアライズ失敗: %v", err)
	}
	info := resp.Corporations[0]
	if len(info.Subsidies) != 1 {
		t.Fatalf("Subsidies count = %d, want 1", len(info.Subsidies))
	}
	s := info.Subsidies[0]
	if s.Title != "IT導入補助金2024" {
		t.Errorf("Title = %q", s.Title)
	}
	if s.SubsidyResource != "中小企業庁" {
		t.Errorf("SubsidyResource = %q", s.SubsidyResource)
	}
}

func TestWorkplaceDeserialization(t *testing.T) {
	var resp model.WorkplaceResponse
	if err := json.Unmarshal(loadFixture(t, "workplace.json"), &resp); err != nil {
		t.Fatalf("デシリアライズ失敗: %v", err)
	}
	info := resp.Corporations[0]
	if len(info.Workplaces) != 1 {
		t.Fatalf("Workplaces count = %d, want 1", len(info.Workplaces))
	}
	w := info.Workplaces[0]
	if w.BaseMonth != "2024-06" {
		t.Errorf("BaseMonth = %q", w.BaseMonth)
	}
	if w.FemaleShareOfManager != "30.0" {
		t.Errorf("FemaleShareOfManager = %q", w.FemaleShareOfManager)
	}
	if w.PaidHolidayUsageRate != "65.0" {
		t.Errorf("PaidHolidayUsageRate = %q", w.PaidHolidayUsageRate)
	}
}
