package sirius

import (
	"fmt"
)

func (c *Client) DeleteDocument(ctx Context, uuid string) error {
	return c.delete(ctx, fmt.Sprintf("/lpa-api/v1/documents/%s", uuid))
}
