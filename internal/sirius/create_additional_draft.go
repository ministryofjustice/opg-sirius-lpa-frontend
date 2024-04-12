package sirius

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type AdditionalDraft struct {
	CaseType                  []string `json:"types"`
	CorrespondentFirstNames   string   `json:"correspondentFirstNames,omitempty"`
	CorrespondentLastName     string   `json:"correspondentLastName,omitempty"`
	CorrespondentAddress      *Address `json:"correspondentAddress,omitempty"`
	CorrespondenceByWelsh     bool     `json:"correspondenceByWelsh,omitempty"`
	CorrespondenceLargeFormat bool     `json:"correspondenceLargeFormat,omitempty"`
	Source                    string   `json:"source"`
}

func (c *Client) CreateAdditionalDraft(ctx Context, donorID int, lpa AdditionalDraft) (map[string]string, error) {
	data, err := json.Marshal(lpa)
	out := map[string]string{}

	if err != nil {
		return out, err
	}

	req, err := c.newRequest(ctx, http.MethodPost, fmt.Sprintf("/lpa-api/v1/donors/%d/digital-lpas", donorID), bytes.NewReader(data))
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
