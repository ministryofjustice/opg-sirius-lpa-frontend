package sirius

import "fmt"

type assignee struct {
	DisplayName string `json:"displayName"`
}

type taskSummary struct {
	ID int `json:"id"`
	Description string `json:"description"`
	Assignee assignee `json:"assignee"`
}

type TaskList struct {
	Tasks []taskSummary `json:"tasks"`
}

func (c *Client) TasksForCase(ctx Context, caseId int) (TaskList, error) {
	//url := fmt.Sprintf("/lpa-api/v1/cases/%d/tasks?filter=status:Not started,active:true&limit=99", caseId)
	url := fmt.Sprintf("/lpa-api/v1/cases/%d/tasks", caseId)

	var v TaskList
	err := c.get(ctx, url, &v)

	return v, err
}
