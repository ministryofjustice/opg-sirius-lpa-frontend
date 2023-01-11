package sirius

import (
	"fmt"
	"net/http"
)

func (c *Client) DeleteDocument(ctx Context, uuid string) error {
	req, err := c.newRequest(ctx, http.MethodDelete, fmt.Sprintf("/lpa-api/v1/documents/%s", uuid), nil)
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
