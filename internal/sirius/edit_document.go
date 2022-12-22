package sirius

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *Client) EditDocument(ctx Context, uuid string, content string) (Document, error) {
	postData, err := json.Marshal(struct {
		Content string `json:"content"`
	}{
		Content: content,
	})
	if err != nil {
		return Document{}, err
	}

	req, err := c.newRequest(ctx, http.MethodPut, fmt.Sprintf("/lpa-api/v1/documents/%s", uuid), bytes.NewReader(postData))
	if err != nil {
		return Document{}, err
	}

	res, err := c.http.Do(req)
	if err != nil {
		return Document{}, err
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusBadRequest {
		var v ValidationError
		if err := json.NewDecoder(res.Body).Decode(&v); err != nil {
			return Document{}, err
		}
		return Document{}, v
	}

	if res.StatusCode != http.StatusOK {
		return Document{}, newStatusError(res)
	}

	var d Document
	if err := json.NewDecoder(res.Body).Decode(&d); err != nil {
		return Document{}, err
	}
	return d, nil
}
