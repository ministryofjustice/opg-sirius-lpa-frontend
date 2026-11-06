package sirius

import (
	"fmt"
)

func (c *Client) UpdateDocument(ctx Context, uuid string, content string) (Document, error) {
	var document Document
	data := struct {
		Content string `json:"content"`
	}{
		Content: content,
	}

	err := c.put(ctx, fmt.Sprintf("/lpa-api/v1/documents/%s", uuid), data, &document)

	return document, err
}
