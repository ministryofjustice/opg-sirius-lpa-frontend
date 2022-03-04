package sirius

import (
	"encoding/json"
	"net/http"
)

func (c *Client) WarningTypes(ctx Context) ([]RefDataItem, error) {
	var v RefData

	req, err := c.newRequest(
		ctx,
		http.MethodGet,
		"/api/v1/reference-data?filter=warningType",
		nil,
	)

	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, newStatusError(res)
	}

	if err := json.NewDecoder(res.Body).Decode(&v); err != nil {
		return nil, err
	}

	return v.WarningTypes, nil
}

type RefDataItem struct {
	Handle string `json:"handle"`
	Label  string `json:"label"`
}

type RefData struct {
	WarningTypes []RefDataItem `json:"warningType"`
}
