package sirius

import (
	"fmt"
	"strconv"
)

func (c *Client) AssignTasks(ctx Context, assigneeID int, taskIDs []int) error {
	urlIDs := strconv.Itoa(taskIDs[0])
	for _, taskID := range taskIDs[1:] {
		urlIDs += "+" + strconv.Itoa(taskID)
	}

	return c.put(ctx, fmt.Sprintf("/lpa-api/v1/users/%d/tasks/%s", assigneeID, urlIDs), nil, nil)
}
