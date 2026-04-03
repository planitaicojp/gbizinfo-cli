package model

type CertificationResponse struct {
	Corporations []CertificationInfo `json:"hojin-infos"`
}

type CertificationInfo struct {
	CorporateNumber string          `json:"corporate_number"`
	Name            string          `json:"name"`
	Certifications  []Certification `json:"certification"`
}

type Certification struct {
	Title                 string `json:"title"`
	DateOfApproval        string `json:"date_of_approval"`
	Target                string `json:"target"`
	Category              string `json:"category"`
	EnterpriseScale       string `json:"enterprise_scale"`
	GovernmentDepartments string `json:"government_departments"`
}
