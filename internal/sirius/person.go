package sirius

import "fmt"

type Person struct {
	ID           int        `json:"id"`
	UID          string     `json:"uId"`
	Salutation   string     `json:"salutation"`
	Firstname    string     `json:"firstname"`
	Surname      string     `json:"surname"`
	DateOfBirth  DateString `json:"dob"`
	AddressLine1 string     `json:"addressLine1"`
	Children     []Person   `json:"children"`
}

func (c *Client) Person(ctx Context, id int) (Person, error) {
	var v Person
	err := c.get(ctx, fmt.Sprintf("/api/v1/persons/%d", id), &v)

	return v, err
}
