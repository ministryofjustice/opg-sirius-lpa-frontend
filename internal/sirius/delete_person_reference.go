package sirius

import (
	"fmt"
	"net/http"
)

func (c *Client) DeletePersonReference(ctx Context, referenceID int) error {
	req, err := c.newRequest(ctx, http.MethodDelete, fmt.Sprintf("/api/v1/person-references/%d", referenceID), nil)
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
