package sirius

import (
	"fmt"
)

func (c *Client) EditInvestigation(ctx Context, investigationID int, investigation Investigation) error {
	return c.put(ctx, fmt.Sprintf("/lpa-api/v1/investigations/%d", investigationID), investigation, nil)
}
