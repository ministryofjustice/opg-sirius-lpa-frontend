package sirius

import (
	"encoding/json"
	"net/http"
)

type RefDataItem struct {
	Handle string `json:"handle"`
	Label  string `json:"label"`
}

func (c *Client) WarningTypes(ctx Context) ([]RefDataItem, error) {
	req, err := c.newRequest(ctx, http.MethodGet, "/api/v1/reference-data/warningType", nil)
	if err != nil {
		return nil, err
	}

	res, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, newStatusError(res)
	}

	var v []RefDataItem
	if err := json.NewDecoder(res.Body).Decode(&v); err != nil {
		return nil, err
	}

	return v, nil
}
