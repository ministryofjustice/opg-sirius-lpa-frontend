package sirius

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Case struct {
	UID      string `json:"uId"`
	CaseType string `json:"caseType"`
}

func (c *Client) Case(ctx Context, id int) (Case, error) {
	var v Case

	req, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf("/api/v1/cases/%d", id), nil)
	if err != nil {
		return v, err
	}

	res, err := c.http.Do(req)
	if err != nil {
		return v, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return v, newStatusError(res)
	}

	if err := json.NewDecoder(res.Body).Decode(&v); err != nil {
		return v, err
	}

	return v, nil
}
