package model

type PatentResponse struct {
	Corporations []PatentInfo `json:"hojin-infos"`
}

type PatentInfo struct {
	CorporateNumber string   `json:"corporate_number"`
	Name            string   `json:"name"`
	Patents         []Patent `json:"patent"`
}

type Patent struct {
	Title             string `json:"title"`
	DateOfApplication string `json:"date_of_application"`
	PatentNumber      string `json:"patent_number"`
	ApplicationNumber string `json:"application_number"`
	ClassificationPI  string `json:"classification_pi"`
}
