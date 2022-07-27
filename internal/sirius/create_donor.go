package sirius

import (
	"bytes"
	"encoding/json"
	"net/http"
)

func (c *Client) CreateDonor(ctx Context, person Person) (Person, error) {
	data, err := json.Marshal(person)
	if err != nil {
		return Person{}, err
	}

	req, err := c.newRequest(ctx, http.MethodPost, "/lpa-api/v1/donors", bytes.NewReader(data))
	if err != nil {
		return Person{}, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return Person{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusBadRequest {
		var v ValidationError
		if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
			return Person{}, err
		}
		return Person{}, v
	}

	if resp.StatusCode != http.StatusCreated {
		return Person{}, newStatusError(resp)
	}

	var v Person
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return Person{}, err
	}

	return v, nil
}
