package sirius

import "fmt"

type Investigations struct {
	Investigations []Investigation `json:"investigations"`
}

func (c *Client) Investigations(ctx Context, caseType string, id int) ([]Investigation, error) {
	var v []Investigation
	err := c.get(ctx, fmt.Sprintf("/lpa-api/v1/%ss/%d/investigations", caseType, id), &v)

	return v, err
}
