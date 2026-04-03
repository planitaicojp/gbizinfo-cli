package model

type ProcurementResponse struct {
	Corporations []ProcurementInfo `json:"hojin-infos"`
}

type ProcurementInfo struct {
	CorporateNumber string        `json:"corporate_number"`
	Name            string        `json:"name"`
	Procurements    []Procurement `json:"procurement"`
}

type Procurement struct {
	Title                 string `json:"title"`
	DateOfOrder           string `json:"date_of_order"`
	Amount                string `json:"amount"`
	GovernmentDepartments string `json:"government_departments"`
	JointSignatures       string `json:"joint_signatures"`
}
