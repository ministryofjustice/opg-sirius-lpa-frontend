package sirius

func (c *Client) CreateCorrespondent(ctx Context, persons []Person) error {
	var v []Person
	return c.post(ctx, "/lpa-api/v1/persons", persons, &v)
}
