package sirius

import (
	"fmt"
)

func (c *Client) GetEvents(ctx Context, id int, caseId string) (any, error) {
	var v struct {
		Events any `json:"events"`
	}

	err := c.get(ctx, fmt.Sprintf("/lpa-api/v1/persons/%d/events", id, caseId), &v)

	return v.Events, err
}
