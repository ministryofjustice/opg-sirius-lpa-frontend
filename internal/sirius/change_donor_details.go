package sirius

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type ChangeDonorDetailsData struct {
	FirstNames        string     `json:"firstNames"`
	LastName          string     `json:"lastName"`
	OtherNamesKnownBy string     `json:"otherNamesKnownBy"`
	DateOfBirth       DateString `json:"dateOfBirth"`
	Address           Address    `json:"address"`
	Phone             string     `json:"phoneNumber"`
	Email             string     `json:"email"`
	LpaSignedOn       DateString `json:"lpaSignedOn"`
}

func (c *Client) ChangeDonorDetails(ctx Context, caseUID string, donorDetails ChangeDonorDetailsData) error {
	data, err := json.Marshal(donorDetails)
	if err != nil {
		return err
	}

	req, err := c.newRequest(ctx, http.MethodPut, fmt.Sprintf("/lpa-api/v1/digital-lpas/%s/change-donor-details", caseUID), bytes.NewReader(data))
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
