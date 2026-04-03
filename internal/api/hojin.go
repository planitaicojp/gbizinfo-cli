package api

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/planitaicojp/gbizinfo-cli/internal/model"
)

func (c *Client) Search(params model.SearchParams) (*model.HojinResponse, error) {
	q := url.Values{}
	if params.Name != "" {
		q.Set("name", params.Name)
	}
	if params.Address != "" {
		q.Set("exist_flg", params.Address)
	}
	if params.CorporateNumber != "" {
		q.Set("corporate_number", params.CorporateNumber)
	}
	if params.Page > 0 {
		q.Set("page", strconv.Itoa(params.Page))
	}
	if params.Limit > 0 {
		q.Set("limit", strconv.Itoa(params.Limit))
	}

	path := "/v1/hojin"
	if len(q) > 0 {
		path += "?" + q.Encode()
	}

	var result model.HojinResponse
	if err := c.Get(path, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) GetHojin(corporateNumber string) (*model.HojinDetail, error) {
	var result model.HojinDetail
	if err := c.Get(fmt.Sprintf("/v1/hojin/%s", corporateNumber), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) GetCertification(corporateNumber string) (*model.CertificationResponse, error) {
	var result model.CertificationResponse
	if err := c.Get(fmt.Sprintf("/v1/hojin/%s/certification", corporateNumber), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) GetCommendation(corporateNumber string) (*model.CommendationResponse, error) {
	var result model.CommendationResponse
	if err := c.Get(fmt.Sprintf("/v1/hojin/%s/commendation", corporateNumber), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) GetFinance(corporateNumber string) (*model.FinanceResponse, error) {
	var result model.FinanceResponse
	if err := c.Get(fmt.Sprintf("/v1/hojin/%s/finance", corporateNumber), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) GetPatent(corporateNumber string) (*model.PatentResponse, error) {
	var result model.PatentResponse
	if err := c.Get(fmt.Sprintf("/v1/hojin/%s/patent", corporateNumber), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) GetProcurement(corporateNumber string) (*model.ProcurementResponse, error) {
	var result model.ProcurementResponse
	if err := c.Get(fmt.Sprintf("/v1/hojin/%s/procurement", corporateNumber), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) GetSubsidy(corporateNumber string) (*model.SubsidyResponse, error) {
	var result model.SubsidyResponse
	if err := c.Get(fmt.Sprintf("/v1/hojin/%s/subsidy", corporateNumber), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) GetWorkplace(corporateNumber string) (*model.WorkplaceResponse, error) {
	var result model.WorkplaceResponse
	if err := c.Get(fmt.Sprintf("/v1/hojin/%s/workplace", corporateNumber), &result); err != nil {
		return nil, err
	}
	return &result, nil
}
