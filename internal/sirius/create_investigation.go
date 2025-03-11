package sirius

import (
	"fmt"
)

func (c *Client) CreateInvestigation(ctx Context, caseID int, caseType CaseType, investigation Investigation) error {
	return c.post(ctx, fmt.Sprintf("/lpa-api/v1/%ss/%d/investigations", caseType, caseID), investigation, nil)
}
