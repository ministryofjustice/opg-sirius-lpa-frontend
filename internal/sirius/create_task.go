package sirius

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"
)

type TaskRequest struct {
	CaseID      int        `json:"caseId"`
	AssigneeID  int        `json:"assigneeId"`
	Type        string     `json:"type"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	DueDate     DateString `json:"dueDate"`
}

func (c *Client) CreateTask(ctx Context, task TaskRequest) error {
	// In the old forms description is checked to enable the save button. We can't
	// validate it server-side as Supervision allow it to be empty.
	if strings.TrimSpace(task.Description) == "" {
		return ValidationError{
			Field: FieldErrors{
				"description": {"": "Value can't be empty"},
			},
		}
	}

	data, err := json.Marshal(task)
	if err != nil {
		return err
	}

	req, err := c.newRequest(ctx, http.MethodPost, "/api/v1/tasks", bytes.NewReader(data))
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
