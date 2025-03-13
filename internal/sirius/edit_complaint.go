package sirius

import (
	"fmt"
)

func (c *Client) EditComplaint(ctx Context, id int, complaint Complaint) error {
	return c.put(ctx, fmt.Sprintf("/lpa-api/v1/complaints/%d", id), complaint, nil)
}
