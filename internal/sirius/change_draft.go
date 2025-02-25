package sirius

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type ChangeDraft struct {
	FirstNames string `json:"firstNames"`
	LastName   string `json:"lastName"`
	//DateOfBirth       DateString `json:"dateOfBirth"`
	Address Address `json:"address"`
	Phone   string  `json:"phoneNumber"`
	Email   string  `json:"email"`
}

func (c *Client) ChangeDraft(ctx Context, caseUID string, draftDetails ChangeDraft) error {
	data, err := json.Marshal(draftDetails)
	if err != nil {
		return err
	}

	req, err := c.newRequest(ctx, http.MethodPut, fmt.Sprintf("/lpa-api/v1/digital-lpas/%s/change-draft", caseUID), bytes.NewReader(data))
	if err != nil {
		return err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

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
