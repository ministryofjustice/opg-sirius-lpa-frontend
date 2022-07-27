package sirius

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type createPersonReferenceRequest struct {
	Reason        string `json:"reason"`
	ReferencedUID string `json:"referencedUid"`
}

func (c *Client) CreatePersonReference(ctx Context, personID int, referencedUID, reason string) error {
	data, err := json.Marshal(createPersonReferenceRequest{
		Reason:        reason,
		ReferencedUID: referencedUID,
	})
	if err != nil {
		return err
	}

	req, err := c.newRequest(ctx, http.MethodPost, fmt.Sprintf("/lpa-api/v1/persons/%d/references", personID), bytes.NewReader(data))
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

	if resp.StatusCode != http.StatusCreated {
		return newStatusError(resp)
	}

	return nil
}
