package sirius

import "fmt"

type Person struct {
	Firstname string `json:"firstname"`
	Surname   string `json:"surname"`
}

func (c *Client) Person(ctx Context, id int) (Person, error) {
	var v Person
	err := c.get(ctx, fmt.Sprintf("/api/v1/persons/%d", id), &v)

	return v, err
}
