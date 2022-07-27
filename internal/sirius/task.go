package sirius

import "fmt"

type Task struct {
	ID        int        `json:"id"`
	Status    string     `json:"status"`
	DueDate   DateString `json:"dueDate"`
	Name      string     `json:"name"`
	CaseItems []Case     `json:"caseItems"`
}

func (t Task) Summary() string {
	return fmt.Sprintf("%s: %s", t.CaseItems[0].Summary(), t.Name)
}

func (c *Client) Task(ctx Context, id int) (Task, error) {
	var v Task
	err := c.get(ctx, fmt.Sprintf("/lpa-api/v1/tasks/%d", id), &v)

	return v, err
}
