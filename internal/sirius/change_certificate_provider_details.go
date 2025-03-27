package sirius

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type ChangeCertificateProviderDetails struct {
	FirstNames string     `json:"firstNames"`
	LastName   string     `json:"lastName"`
	Address    Address    `json:"address"`
	Phone      string     `json:"phoneNumber"`
	Email      string     `json:"email"`
	SignedAt   DateString `json:"signedAt"`
}

func (c *Client) ChangeCertificateProviderDetails(ctx Context, caseUID string, certificateProviderDetails ChangeCertificateProviderDetails) error {
	data, err := json.Marshal(certificateProviderDetails)
	if err != nil {
		return err
	}

	req, err := c.newRequest(
		ctx,
		http.MethodPut,
		fmt.Sprintf("/lpa-api/v1/digital-lpas/%s/change-certificate-provider-details", caseUID),
		bytes.NewReader(data),
	)
	if err != nil {
		return err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close() //nolint:errcheck // no need to check error when closing body

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
