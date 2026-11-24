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
	Description string     `json:"description"`
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
	querystring.Set("sort", "duedate:ASC")

	req, err := c.newRequestWithQuery(ctx, http.MethodGet, path, querystring, nil)

	if err != nil {
		return nil, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close() //nolint:errcheck // no need to check error when closing body

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

func (c *Client) TasksForDonor(ctx Context, donorId int) ([]Task, error) {
	fmt.Printf("TasksForDonor: fetching tasks for donor ID %d\n", donorId)

	path := fmt.Sprintf("/lpa-api/v1/persons/%d/tasks", donorId)
	fmt.Printf("TasksForDonor: constructed path: %s\n", path)

	querystring := url.Values{}
	querystring.Set("filter", "status:Not started,active:true")
	querystring.Set("limit", "99")
	querystring.Set("sort", "duedate:ASC")
	fmt.Printf("TasksForDonor: querystring: %s\n", querystring.Encode())

	req, err := c.newRequestWithQuery(ctx, http.MethodGet, path, querystring, nil)

	if err != nil {
		fmt.Printf("TasksForDonor: error creating request: %v\n", err)
		return nil, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		fmt.Printf("TasksForDonor: error executing request: %v\n", err)
		return nil, err
	}

	defer resp.Body.Close() //nolint:errcheck // no need to check error when closing body

	fmt.Printf("TasksForDonor: received response with status code: %d\n", resp.StatusCode)

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("TasksForDonor: unexpected status code, returning error\n")
		return nil, newStatusError(resp)
	}

	var v taskList
	err = json.NewDecoder(resp.Body).Decode(&v)
	if err != nil {
		fmt.Printf("TasksForDonor: error decoding response body: %v\n", err)
		return nil, err
	}

	fmt.Printf("TasksForDonor: successfully retrieved %d tasks for donor ID %d\n", len(v.Tasks), donorId)
	return v.Tasks, err
}
