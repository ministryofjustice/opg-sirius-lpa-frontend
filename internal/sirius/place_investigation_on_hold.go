package sirius

import (
	"fmt"
)

func (c *Client) PlaceInvestigationOnHold(ctx Context, investigationID int, reason string) error {
	postData := struct {
		Reason string `json:"reason"`
	}{
		Reason: reason,
	}

	return c.post(ctx, fmt.Sprintf("/lpa-api/v1/investigations/%d/hold-periods", investigationID), postData, nil)
}
