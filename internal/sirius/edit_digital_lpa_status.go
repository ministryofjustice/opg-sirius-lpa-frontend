package sirius

import (
	"fmt"
)

type CaseStatusData struct {
	Status string `json:"status"`
}

func (c *Client) EditDigitalLPAStatus(ctx Context, caseUID string, caseStatusData CaseStatusData) error {
	return c.put(ctx, fmt.Sprintf("/lpa-api/v1/digital-lpas/%s/update-case-status", caseUID), caseStatusData, nil)
}
