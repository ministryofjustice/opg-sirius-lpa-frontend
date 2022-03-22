package sirius

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Person struct {
	Firstname string `json:"firstname"`
	Surname   string `json:"surname"`
}

func (c *Client) Person(ctx Context, id int) (Person, error) {
	var v Person

	req, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf("/api/v1/persons/%d", id), nil)
	if err != nil {
		return v, err
	}

	res, err := c.http.Do(req)
	if err != nil {
		return v, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return v, newStatusError(res)
	}

	if err := json.NewDecoder(res.Body).Decode(&v); err != nil {
		return v, err
	}

	return v, nil
}
