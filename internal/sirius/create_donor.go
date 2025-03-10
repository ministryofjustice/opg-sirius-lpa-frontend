package sirius

func (c *Client) CreateDonor(ctx Context, person Person) (Person, error) {
	var v Person
	err := c.post(ctx, "/lpa-api/v1/donors", person, &v)

	return v, err
}
