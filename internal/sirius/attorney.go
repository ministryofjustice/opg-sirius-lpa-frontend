package sirius

import "fmt"

func (c *Client) CreateAttorney(ctx Context, caseId int, attorney Person) error {
	return c.post(ctx, fmt.Sprintf("/lpa-api/v1/epas/%d/attorneys", caseId), attorney, nil)
}
