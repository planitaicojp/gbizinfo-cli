package api

import (
	"net/url"
	"strconv"

	"github.com/planitaicojp/gbizinfo-cli/internal/model"
)

func (c *Client) getUpdate(path string, params model.UpdateParams) (*model.UpdateResponse, error) {
	q := url.Values{}
	if params.From != "" {
		q.Set("from", params.From)
	}
	if params.To != "" {
		q.Set("to", params.To)
	}
	if params.Page > 0 {
		q.Set("page", strconv.Itoa(params.Page))
	}
	if len(q) > 0 {
		path += "?" + q.Encode()
	}

	var result model.UpdateResponse
	if err := c.Get(path, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) GetUpdateInfo(params model.UpdateParams) (*model.UpdateResponse, error) {
	return c.getUpdate("/v1/hojin/updateInfo", params)
}

func (c *Client) GetUpdateCertification(params model.UpdateParams) (*model.UpdateResponse, error) {
	return c.getUpdate("/v1/hojin/updateInfo/certification", params)
}

func (c *Client) GetUpdateCommendation(params model.UpdateParams) (*model.UpdateResponse, error) {
	return c.getUpdate("/v1/hojin/updateInfo/commendation", params)
}

func (c *Client) GetUpdateFinance(params model.UpdateParams) (*model.UpdateResponse, error) {
	return c.getUpdate("/v1/hojin/updateInfo/finance", params)
}

func (c *Client) GetUpdatePatent(params model.UpdateParams) (*model.UpdateResponse, error) {
	return c.getUpdate("/v1/hojin/updateInfo/patent", params)
}

func (c *Client) GetUpdateProcurement(params model.UpdateParams) (*model.UpdateResponse, error) {
	return c.getUpdate("/v1/hojin/updateInfo/procurement", params)
}

func (c *Client) GetUpdateSubsidy(params model.UpdateParams) (*model.UpdateResponse, error) {
	return c.getUpdate("/v1/hojin/updateInfo/subsidy", params)
}

func (c *Client) GetUpdateWorkplace(params model.UpdateParams) (*model.UpdateResponse, error) {
	return c.getUpdate("/v1/hojin/updateInfo/workplace", params)
}
