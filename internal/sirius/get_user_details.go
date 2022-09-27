package sirius

func (c *Client) GetUserDetails(ctx Context) (User, error) {
	var v User
	err := c.get(ctx, "/lpa-api/v1/users/current", &v)

	return v, err
}

func (u User) HasRole(searchRole string) bool {
	for _, userRole := range u.Roles {
		if userRole == searchRole {
			return true
		}
	}

	return false
}
