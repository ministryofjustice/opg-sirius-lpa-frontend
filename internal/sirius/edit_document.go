package sirius

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *Client) EditDocument(ctx Context, uuid string, content string) error {
	postData, err := json.Marshal(struct {
		Content string `json:"content"`
	}{
		Content: content,
	})

	if err != nil {
		return err
	}

	req, err := c.newRequest(
		ctx,
		http.MethodPut,
		fmt.Sprintf("/lpa-api/v1/documents/%s", uuid),
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

	if res.StatusCode != http.StatusOK {
		return newStatusError(res)
	}
	return nil
}
