package sirius

import "fmt"

func (c *Client) PersonByUid(ctx Context, uid string) (Person, error) {
	var v Person
	err := c.get(ctx, fmt.Sprintf("/api/v1/persons/by-uid/%s", uid), &v)

	return v, err
}
