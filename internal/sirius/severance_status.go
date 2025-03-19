package sirius

import (
	"fmt"
)

type SeveranceStatusData struct {
	Status string `json:"status"`
}

func (c *Client) SeveranceRequired(ctx Context, caseUID string, severanceStatusData SeveranceStatusData) error {
	err := c.put(ctx, fmt.Sprintf("/lpa-api/v1/digital-lpas/%s/severance-status", caseUID), severanceStatusData, nil)

	return err
}
