package sirius

import "fmt"

type Task struct {
	ID          int        `json:"id"`
	Status      string     `json:"status"`
	DueDate     DateString `json:"dueDate"`
	Name        string     `json:"name"`
	CaseItems   []Case     `json:"caseItems"`
    Description string     `json:"description"`
	Assignee    User       `json:"assignee"`
}

type taskList struct {
    Tasks []Task `json:"tasks"`
}

func (t Task) Summary() string {
	if len(t.CaseItems) > 0 {
		return fmt.Sprintf("%s: %s", t.CaseItems[0].Summary(), t.Name)
	} else {
		return t.Name
	}
}

func (c *Client) Task(ctx Context, id int) (Task, error) {
	var v Task
	err := c.get(ctx, fmt.Sprintf("/lpa-api/v1/tasks/%d", id), &v)

	return v, err
}

func (c *Client) TasksForCase(ctx Context, caseId int) ([]Task, error) {
	//url := fmt.Sprintf("/lpa-api/v1/cases/%d/tasks?filter=status:Not+started,active:true&limit=99", caseId)
	url := fmt.Sprintf("/lpa-api/v1/cases/%d/tasks", caseId)

	var v taskList
	err := c.get(ctx, url, &v)

	return v.Tasks, err
}
