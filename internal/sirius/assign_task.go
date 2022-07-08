package sirius

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type assignTaskRequest struct {
	AssigneeID int `json:"assigneeId"`
	ID         int `json:"id"`
}

func (c *Client) AssignTask(ctx Context, assigneeID, taskID int) error {
	data, err := json.Marshal(assignTaskRequest{AssigneeID: assigneeID, ID: taskID})
	if err != nil {
		return err
	}

	req, err := c.newRequest(ctx, http.MethodPut, fmt.Sprintf("/api/v1/tasks/%d", taskID), bytes.NewReader(data))
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

	if resp.StatusCode != http.StatusOK {
		return newStatusError(resp)
	}

	return nil
}
