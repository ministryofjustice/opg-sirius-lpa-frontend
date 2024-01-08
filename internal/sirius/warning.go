package sirius

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Warning struct {
	ID          int    `json:"id"`
	DateAdded   string `json:"dateAdded"`
	WarningType string `json:"warningType"`
	WarningText string `json:"warningText"`
	CaseItems   []Case `json:"caseItems"`
}

func (c *Client) WarningsForCase(ctx Context, caseId int) ([]Warning, error) {
	path := fmt.Sprintf("/lpa-api/v1/cases/%d/warnings", caseId)

	req, err := c.newRequest(ctx, http.MethodGet, path, nil)

	if err != nil {
		return nil, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close() //#nosec G307 false positive

	if resp.StatusCode != http.StatusOK {
		return nil, newStatusError(resp)
	}

	var warningList []Warning
	err = json.NewDecoder(resp.Body).Decode(&warningList)
	if err != nil {
		return nil, err
	}

	return warningList, nil
}
