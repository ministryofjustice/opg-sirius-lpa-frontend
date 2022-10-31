package sirius

import "fmt"

type HoldPeriod struct {
	ID            int           `json:"id,omitempty"`
	Investigation Investigation `json:"investigation"`
	Reason        string        `json:"reason"`
}

func (c *Client) HoldPeriod(ctx Context, id int) (HoldPeriod, error) {
	var v HoldPeriod
	err := c.get(ctx, fmt.Sprintf("/lpa-api/v1/hold-periods/%d", id), &v)

	return v, err
}
