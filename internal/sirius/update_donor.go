package sirius

import (
	"fmt"
)

func (c *Client) UpdateDonor(ctx Context, personID int, person Person) error {
	return c.put(ctx, fmt.Sprintf("/lpa-api/v1/donors/%d", personID), person, nil)
}
