package sirius

func (c *Client) CreateContact(ctx Context, contact Person) (Person, error) {
	var v Person

	err := c.post(ctx, "/lpa-api/v1/non-case-contacts", contact, &v)

	return v, err
}
