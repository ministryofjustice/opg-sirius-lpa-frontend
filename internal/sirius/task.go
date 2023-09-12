package sirius

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type Task struct {
	ID          int        `json:"id"`
	Status      string     `json:"status"`
	DueDate     DateString `json:"dueDate"`
	Name        string     `json:"name"`
	CaseItems   []Case     `json:"caseItems"`
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
	path := fmt.Sprintf("/lpa-api/v1/cases/%d/tasks", caseId)

	querystring := url.Values{}
	querystring.Set("filter", "status:Not started,active:true")
	querystring.Set("limit", "99")

	req, err := c.newRequestWithQuery(ctx, http.MethodGet, path, querystring, nil)

	if err != nil {
		return nil, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close() //#nosec G307 false positive

	if resp.StatusCode != http.StatusOK {
		return nil, newStatusError(resp)
	}

	var v taskList
	err = json.NewDecoder(resp.Body).Decode(&v)
	if err != nil {
		return nil, err
	}

	return v.Tasks, err
}
