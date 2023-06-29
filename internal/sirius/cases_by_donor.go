package sirius

import (
	"fmt"
)

func (c *Client) CasesByDonor(ctx Context, id int) ([]Case, error) {
	var v struct {
		Cases []Case `json:"cases"`
	}

	err := c.get(ctx, fmt.Sprintf("/lpa-api/v1/persons/%d/cases", id), &v)

	return v.Cases, err
}
