package sirius

import (
	"encoding/json"
	"net/http"
)

func (c *Client) NoteTypes(ctx Context) ([]string, error) {
	req, err := c.newRequest(ctx, http.MethodGet, "/api/v1/note-types/lpa", nil)
	if err != nil {
		return nil, err
	}

	res, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, newStatusError(res)
	}

	var v []string
	if err := json.NewDecoder(res.Body).Decode(&v); err != nil {
		return nil, err
	}

	return v, nil
}
