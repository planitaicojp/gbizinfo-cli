package model

type PageInfo struct {
	TotalCount int `json:"totalCount"`
	TotalPage  int `json:"totalPage"`
	PageNumber int `json:"pageNumber"`
}

type HojinResponse struct {
	PageInfo
	Corporations []Hojin `json:"hojin-infos"`
}

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

type HojinDetail struct {
	Corporations []Hojin `json:"hojin-infos"`
}
