package sirius

import (
	"bytes"
	"encoding/json"
	"net/http"
)

func (c *Client) LinkPeople(ctx Context, parentId int, childId int) error {
	postData, err := json.Marshal(struct {
		ParentId int `json:"parentId"`
		ChildId  int `json:"childId"`
	}{
		ParentId: parentId,
		ChildId:  childId,
	})

	if err != nil {
		return err
	}

	req, err := c.newRequest(
		ctx,
		http.MethodPost,
		"/api/v1/person-links",
		bytes.NewReader(postData),
	)

	if err != nil {
		return err
	}

	res, err := c.http.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode == http.StatusBadRequest {
		var v ValidationError
		if err := json.NewDecoder(res.Body).Decode(&v); err != nil {
			return err
		}
		return v
	}

	if res.StatusCode != http.StatusNoContent {
		return newStatusError(res)
	}

	return nil
}
