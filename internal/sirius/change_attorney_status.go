package sirius

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type AttorneyUpdatedStatus struct {
	UID    string `json:"uid"`
	Status string `json:"status"`
}

type attorneyStatusRequest struct {
	AttorneyStatuses []AttorneyUpdatedStatus `json:"attorneyStatuses"`
}

func (c *Client) ChangeAttorneyStatus(ctx Context, caseUID string, attorneyUpdatedStatus []AttorneyUpdatedStatus) error {
	data, err := json.Marshal(attorneyStatusRequest{AttorneyStatuses: attorneyUpdatedStatus})
	if err != nil {
		return err
	}

	req, err := c.newRequest(ctx, http.MethodPut, fmt.Sprintf("/lpa-api/v1/digital-lpas/%s/attorney-status", caseUID), bytes.NewReader(data))
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

	if resp.StatusCode != http.StatusNoContent {
		return newStatusError(resp)
	}

	return nil
}
