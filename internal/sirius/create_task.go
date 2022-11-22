package sirius

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type TaskRequest struct {
	AssigneeID  int        `json:"assigneeId"`
	Type        string     `json:"type"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	DueDate     DateString `json:"dueDate"`
}

func (c *Client) CreateTask(ctx Context, caseID int, task TaskRequest) error {
	data, err := json.Marshal(task)
	if err != nil {
		return err
	}

	req, err := c.newRequest(ctx, http.MethodPost, fmt.Sprintf("/lpa-api/v1/cases/%d/tasks", caseID), bytes.NewReader(data))
	if err != nil {
		return err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusBadRequest {
		var v ValidationError
		if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
			return err
		}
		return v
	}

	if resp.StatusCode != http.StatusCreated {
		return newStatusError(resp)
	}
	return nil
}
