package sirius

import (
	"bytes"
	"encoding/json"
	"net/http"
)

func (c *Client) CreateWarning(ctx Context, personId int, warningType string, warningNote string, caseIDs []int) error {
	type WarningRequest struct {
		PersonID    int    `json:"personId,omitempty"`
		CaseIDs     []int  `json:"caseIds,omitempty"`
		WarningType string `json:"warningType"`
		WarningText string `json:"warningText"`
	}

	Data := WarningRequest{
		PersonID:    personId,
		CaseIDs:     caseIDs,
		WarningType: warningType,
		WarningText: warningNote,
	}

	postData, err := json.Marshal(Data)

	if err != nil {
		return err
	}

	req, err := c.newRequest(
		ctx, http.MethodPost,
		"/lpa-api/v1/warnings",
		bytes.NewReader(postData),
	)

	if err != nil {
		return err
	}

	res, err := c.http.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close() //#nosec G307 false positive

	if res.StatusCode == http.StatusBadRequest {
		var v ValidationError
		if err := json.NewDecoder(res.Body).Decode(&v); err != nil {
			return err
		}
		return v
	}

	if res.StatusCode != http.StatusCreated {
		return newStatusError(res)
	}
	return nil
}
