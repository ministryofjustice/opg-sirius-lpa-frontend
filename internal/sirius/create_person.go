package sirius

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type PersonType string

const (
	PersonTypeDonor = PersonType("donor")
)

func (c *Client) CreatePerson(ctx Context, personType PersonType, person Person) (Person, error) {
	data, err := json.Marshal(person)
	if err != nil {
		return Person{}, err
	}

	endpoint := "/lpa-api/v1/persons"
	if personType == "donor" {
		endpoint = "/lpa-api/v1/donors"
	}

	req, err := c.newRequest(ctx, http.MethodPost, endpoint, bytes.NewReader(data))
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
