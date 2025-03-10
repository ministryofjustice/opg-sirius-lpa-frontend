package sirius

import (
	"fmt"
)

func (c *Client) DeletePersonReference(ctx Context, referenceID int) error {
	return c.delete(ctx, fmt.Sprintf("/lpa-api/v1/person-references/%d", referenceID))
}
