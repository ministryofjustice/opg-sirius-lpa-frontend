package sirius

import (
	"fmt"
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
	out := map[string]string{}
	var v []struct {
		Subtype string `json:"caseSubtype"`
		Uid     string `json:"uId"`
	}

	err := c.post(ctx, fmt.Sprintf("/lpa-api/v1/donors/%d/digital-lpas", donorID), lpa, &v)

	if err != nil {
		return out, nil
	}

	for _, lpa := range v {
		out[lpa.Subtype] = lpa.Uid
	}

	return out, nil
}
