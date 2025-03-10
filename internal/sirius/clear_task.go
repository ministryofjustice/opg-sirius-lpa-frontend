package sirius

import (
	"fmt"
)

func (c *Client) ClearTask(ctx Context, taskID int) error {
	var v interface{}
	err := c.put(ctx, fmt.Sprintf("/lpa-api/v1/tasks/%d/mark-as-completed", taskID), nil, &v)

	return err
}
