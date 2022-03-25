package sirius

type Team struct {
	ID          int    `json:"id"`
	DisplayName string `json:"displayName"`
}

func (c *Client) Teams(ctx Context) ([]Team, error) {
	var v []Team
	err := c.get(ctx, "/api/v1/teams", &v)

	return v, err
}
