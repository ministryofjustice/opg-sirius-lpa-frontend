package sirius

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type DraftLpa struct {
	CaseType                  []string `json:"types"`
	CorrespondentFirstNames   string   `json:"correspondentFirstNames,omitempty"`
	CorrespondentLastName     string   `json:"correspondentLastName,omitempty"`
	CorrespondentAddress      *Address `json:"correspondentAddress,omitempty"`
	CorrespondenceByWelsh     bool     `json:"correspondenceByWelsh,omitempty"`
	CorrespondenceLargeFormat bool     `json:"correspondenceLargeFormat,omitempty"`
}

func (c *Client) CreateDraftLpa(ctx Context, donorID int, lpa DraftLpa) (map[string]string, error) {
	data, err := json.Marshal(lpa)
	out := map[string]string{}

	if err != nil {
		return out, err
	}

	req, err := c.newRequest(ctx, http.MethodPost, fmt.Sprintf("/lpa-api/v1/donor/%d/digital-lpas", donorID), bytes.NewReader(data))
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
	return out, nil
}
