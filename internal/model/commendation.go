package model

type CommendationResponse struct {
	Corporations []CommendationInfo `json:"hojin-infos"`
}

type CommendationInfo struct {
	CorporateNumber string         `json:"corporate_number"`
	Name            string         `json:"name"`
	Commendations   []Commendation `json:"commendation"`
}

type Commendation struct {
	Title                 string `json:"title"`
	DateOfCommendation    string `json:"date_of_commendation"`
	Target                string `json:"target"`
	Category              string `json:"category"`
	GovernmentDepartments string `json:"government_departments"`
}
