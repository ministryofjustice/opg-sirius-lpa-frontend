package sirius

import (
	"fmt"
)

type Complaints struct {
	Complaints []Complaint `json:"complaint"`
}

func (c *Client) Complaints(ctx Context, caseType string, caseId int) ([]Complaint, error) {
	var v []Complaint

	err := c.get(ctx, fmt.Sprintf("/lpa-api/v1/%ss/%d/complaints", caseType, caseId), &v)
	return v, err
}
