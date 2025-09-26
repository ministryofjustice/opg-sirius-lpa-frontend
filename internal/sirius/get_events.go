package sirius

import (
	"fmt"
)

func (c *Client) GetEvents(ctx Context, donorId int, caseId int) (any, error) {
	var v struct {
		Events any `json:"events"`
	}

	err := c.get(ctx, fmt.Sprintf("/lpa-api/v1/persons/%d/events?filter=case:%d&sort=id:desc", donorId, caseId), &v)

	return v.Events, err
}

// GetCombinedEvents Gets combined events from both Sirius and LPA Store for digital LPAs
func (c *Client) GetCombinedEvents(ctx Context, uid string, donorId int, caseId int) (any, error) {
	var v struct {
		Events any `json:"events"`
	}

	err := c.get(ctx, fmt.Sprintf("/lpa-api/v1/digital-lpas/%s/events?donorId=%d&caseId=%d", uid, donorId, caseId), &v)

	return v.Events, err
}
