package sirius

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *Client) AddComplaint(ctx Context, caseID int, caseType CaseType, complaint Complaint) error {
	data, err := json.Marshal(complaint)
	if err != nil {
		return err
	}

	req, err := c.newRequest(ctx, http.MethodPost, fmt.Sprintf("/lpa-api/v1/%ss/%d/complaints", caseType, caseID), bytes.NewReader(data))
	if err != nil {
		return err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close() //#nosec G307 false positive

	if resp.StatusCode == http.StatusBadRequest {
		var v ValidationError
		if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
			return err
		}
		return v
	}

	if resp.StatusCode != http.StatusCreated {
		return newStatusError(resp)
	}

	return nil
}
