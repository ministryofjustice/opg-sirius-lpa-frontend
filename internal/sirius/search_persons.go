package sirius

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type searchPersonsRequest struct {
	Term string `json:"term"`
}

type searchPersonsResponse struct {
	Results []Person `json:"results"`
}

func (c *Client) SearchPersons(ctx Context, term string) ([]Person, error) {
	if len(term) < 3 {
		return nil, fmt.Errorf("Search term must be at least three characters")
	}

	data, err := json.Marshal(searchPersonsRequest{Term: term})
	if err != nil {
		return nil, err
	}

	req, err := c.newRequest(ctx, http.MethodPost, "/api/v1/search/persons", bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, newStatusError(resp)
	}

	var v searchPersonsResponse
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return nil, err
	}

	return v.Results, nil
}
