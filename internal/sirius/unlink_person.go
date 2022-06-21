package sirius

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *Client) UnlinkPerson(ctx Context, parentId int, childIds []int) error {
	data, err := json.Marshal(struct {
		ChildIds []int `json:"childIds"`
	}{
		ChildIds: childIds,
	})

	if err != nil {
		return err
	}

	req, err := c.newRequest(ctx, http.MethodDelete, fmt.Sprintf("/api/v1/person-links/%d", parentId), bytes.NewReader(data))
	if err != nil {
		return err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return newStatusError(resp)
	}

	return nil
}
