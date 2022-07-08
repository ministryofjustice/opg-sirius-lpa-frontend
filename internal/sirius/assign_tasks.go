package sirius

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func (c *Client) AssignTasks(ctx Context, assigneeID int, taskIDs []int) error {
	urlIDs := strconv.Itoa(taskIDs[0])
	for _, taskID := range taskIDs[1:] {
		urlIDs += "+" + strconv.Itoa(taskID)
	}

	req, err := c.newRequest(ctx, http.MethodPut, fmt.Sprintf("/api/v1/users/%d/tasks/%s", assigneeID, urlIDs), strings.NewReader(""))
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
