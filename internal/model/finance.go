package model

type FinanceResponse struct {
	Corporations []FinanceInfo `json:"hojin-infos"`
}

type FinanceInfo struct {
	CorporateNumber string    `json:"corporate_number"`
	Name            string    `json:"name"`
	Finance         []Finance `json:"finance"`
}

type Finance struct {
	AccountingPeriod  string `json:"accounting_period"`
	MajorShareholders string `json:"major_shareholders"`
	NetSales          string `json:"net_sales"`
	OperatingRevenue  string `json:"operating_revenue"`
	OrdinaryIncome    string `json:"ordinary_income"`
	Profit            string `json:"profit"`
	TotalAssets       string `json:"total_assets"`
	NetAssets         string `json:"net_assets"`
	CapitalStock      string `json:"capital_stock"`
	EmployeeNumber    string `json:"employee_number"`
}
