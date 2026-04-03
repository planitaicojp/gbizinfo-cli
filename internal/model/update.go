package model

type UpdateResponse struct {
	PageInfo
	Corporations []Hojin `json:"hojin-infos"`
}

type UpdateParams struct {
	From string
	To   string
	Page int
}

type SearchParams struct {
	Name            string
	Address         string
	CorporateNumber string
	Page            int
	Limit           int
}
