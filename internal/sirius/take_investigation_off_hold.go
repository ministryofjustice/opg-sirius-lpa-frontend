package sirius

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *Client) TakeInvestigationOffHold(ctx Context, holdPeriodId int) error {
	req, err := c.newRequest(ctx, http.MethodDelete, fmt.Sprintf("/lpa-api/v1/hold-periods/%d", holdPeriodId), nil)
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
