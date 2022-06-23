package sirius

import (
	"fmt"
)

func (c *Client) GetAllowedStatuses(ctx Context, caseId int, caseType CaseType) ([]string, error) {
	var v []string
	err := c.get(ctx, fmt.Sprintf("/api/v1/%ss/%d/available-statuses", caseType, caseId), &v)

	return v, err
}
