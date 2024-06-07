package sirius

import (
	"fmt"
	"net/http"
	"strings"
)

func (c *Client) ClearTask(ctx Context, taskID int) error {

	req, err := c.newRequest(ctx, http.MethodPut, fmt.Sprintf("/lpa-api/v1/tasks/%d/mark-as-completed", taskID), strings.NewReader(""))

	if err != nil {
		return err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close() //#nosec G307 false positive

	if resp.StatusCode != http.StatusOK {
		return newStatusError(resp)
	}

	return nil
}
