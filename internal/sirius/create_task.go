package sirius

import (
	"fmt"
)

type TaskRequest struct {
	AssigneeID  int        `json:"assigneeId"`
	Type        string     `json:"type"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	DueDate     DateString `json:"dueDate"`
}

func (c *Client) CreateTask(ctx Context, caseID int, task TaskRequest) error {
	return c.post(ctx, fmt.Sprintf("/lpa-api/v1/cases/%d/tasks", caseID), task, nil)
}
