package sirius

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type searchDonorsRequest struct {
	Term        string   `json:"term"`
	PersonTypes []string `json:"personTypes"`
}

type searchDonorsResponse struct {
	Results []Person `json:"results"`
}

func (c *Client) SearchDonors(ctx Context, term string) ([]Person, error) {
	if len(term) < 3 {
		return nil, fmt.Errorf("Search term must be at least three characters")
	}

	data, err := json.Marshal(searchDonorsRequest{Term: term, PersonTypes: []string{"Donor"}})
	if err != nil {
		return nil, err
	}

	req, err := c.newRequest(ctx, http.MethodPost, "/lpa-api/v1/search/persons", bytes.NewReader(data))
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

	var v searchDonorsResponse
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return nil, err
	}

	return v.Results, nil
}
