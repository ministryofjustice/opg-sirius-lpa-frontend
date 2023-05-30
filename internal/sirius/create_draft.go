package sirius

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type Address struct {
	Line1    string `json:"addressLine1"`
	Line2    string `json:"addressLine2,omitempty"`
	Line3    string `json:"addressLine3,omitempty"`
	Town     string `json:"town"`
	Postcode string `json:"postcode"`
	Country  string `json:"country"`
}

type Draft struct {
	CaseType             []string   `json:"types"`
	Source               string     `json:"source"`
	DonorName            string     `json:"donorName"`
	DonorDob             DateString `json:"donorDob"`
	DonorAddress         Address    `json:"donorAddress"`
	CorrespondentName    string     `json:"correspondentName,omitempty"`
	CorrespondentAddress *Address   `json:"correspondentAddress,omitempty"`
	PhoneNumber          string     `json:"donorPhone,omitempty"`
	Email                string     `json:"donorEmail,omitempty"`
}

func (c *Client) CreateDraft(ctx Context, draft Draft) (map[string]string, error) {
	data, err := json.Marshal(draft)
	out := map[string]string{}

	if err != nil {
		return out, err
	}

	req, err := c.newRequest(ctx, http.MethodPost, "/lpa-api/v1/digital-lpas", bytes.NewReader(data))
	if err != nil {
		return out, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return out, err
	}
	defer resp.Body.Close() //#nosec G307 false positive

	if resp.StatusCode == http.StatusBadRequest {
		var v ValidationError
		if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
			return out, err
		}
		return out, v
	}

	if resp.StatusCode != http.StatusCreated {
		return out, newStatusError(resp)
	}

	var v []struct {
		Subtype string `json:"caseSubtype"`
		Uid     string `json:"uId"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return out, err
	}

	for _, lpa := range v {
		out[lpa.Subtype] = lpa.Uid
	}

	return out, nil
}
