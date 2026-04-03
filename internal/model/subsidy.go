package model

type SubsidyResponse struct {
	Corporations []SubsidyInfo `json:"hojin-infos"`
}

type SubsidyInfo struct {
	CorporateNumber string    `json:"corporate_number"`
	Name            string    `json:"name"`
	Subsidies       []Subsidy `json:"subsidy"`
}

type Subsidy struct {
	Title                 string `json:"title"`
	DateOfApproval        string `json:"date_of_approval"`
	Amount                string `json:"amount"`
	SubsidyResource       string `json:"subsidy_resource"`
	Target                string `json:"target"`
	GovernmentDepartments string `json:"government_departments"`
	Note                  string `json:"note"`
}
