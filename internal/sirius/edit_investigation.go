package sirius

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode == http.StatusBadRequest {
		var v ValidationError
		var x FlexibleFields
		if err := json.Unmarshal(body, &x); err != nil {
			return err
		}
		formattedFieldErrors, err := formatToFieldErrors(x)
		if err != nil {
			return err
		}
		v.Field = formattedFieldErrors
		v.Detail = "Payload failed validation"
		return v
	}

	if resp.StatusCode != http.StatusOK {
		return newStatusError(resp)
	}

	return nil
}
