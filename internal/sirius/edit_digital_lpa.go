package sirius

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type CaseStatusData struct {
	Status string
}

func (c *Client) EditDigitalLPA(ctx Context, caseUID string, caseStatusData CaseStatusData) error {
	data, err := json.Marshal(caseStatusData)
	if err != nil {
		return err
	}

	req, err := c.newRequest(ctx, http.MethodPut, fmt.Sprintf("/lpa-api/v1/digital-lpas/%s/update-case-status", caseUID), bytes.NewReader(data))
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

	if resp.StatusCode != http.StatusOK {
		return newStatusError(resp)
	}

	return nil
}
