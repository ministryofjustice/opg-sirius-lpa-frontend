package sirius

func (c *Client) SearchDonors(ctx Context, term string) ([]Person, error) {
	resp, _, err := c.Search(ctx, term, 1, []string{"Donor"}, []string{"person"})
	if err != nil {
		return nil, err
	}

	return resp.Results, nil
}
