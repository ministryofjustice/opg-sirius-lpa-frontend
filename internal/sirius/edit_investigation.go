package sirius

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *Client) EditInvestigation(ctx Context, investigationID int, investigation Investigation) error {
	data, err := json.Marshal(investigation)
	if err != nil {
		return err
	}

	req, err := c.newRequest(ctx, http.MethodPut, fmt.Sprintf("/lpa-api/v1/investigations/%d", investigationID), bytes.NewReader(data))
	if err != nil {
		return err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusBadRequest {
		var v SiriusValidationError
		if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
			return err
		}
		return FormatToValidationError(v)
	}

	if resp.StatusCode != http.StatusOK {
		return newStatusError(resp)
	}

	return nil
}
