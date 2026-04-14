package sirius

import "fmt"

func (c *Client) CreateEpa(ctx Context, donorID int, epa Case) error {
	return c.post(ctx, fmt.Sprintf("/lpa-api/v1/donors/%d/epas", donorID), epa, nil)
}
