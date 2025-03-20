package sirius

import (
	"fmt"
)

type SeveranceStatusData struct {
	SeveranceStatus string `json:"severanceStatus"`
}

func (c *Client) UpdateSeveranceStatus(ctx Context, caseUID string, severanceStatusData SeveranceStatusData) error {
	err := c.put(ctx, fmt.Sprintf("/lpa-api/v1/digital-lpas/%s/severance-status", caseUID), severanceStatusData, nil)

	return err
}
