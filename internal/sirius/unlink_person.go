package sirius

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *Client) UnlinkPerson(ctx Context, parentId int, childId int) error {
	data, err := json.Marshal(struct {
		ChildIds []int `json:"childIds"`
	}{
		ChildIds: []int{childId},
	})

	if err != nil {
		return err
	}

	req, err := c.newRequest(ctx, http.MethodPatch, fmt.Sprintf("/lpa-api/v1/person-links/%d", parentId), bytes.NewReader(data))
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

	if resp.StatusCode != http.StatusNoContent {
		return newStatusError(resp)
	}

	return nil
}
