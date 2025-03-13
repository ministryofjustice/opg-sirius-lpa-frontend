package sirius

import (
	"fmt"
)

func (c *Client) AddComplaint(ctx Context, caseID int, caseType CaseType, complaint Complaint) error {
	return c.post(ctx, fmt.Sprintf("/lpa-api/v1/%ss/%d/complaints", caseType, caseID), complaint, nil)
}
