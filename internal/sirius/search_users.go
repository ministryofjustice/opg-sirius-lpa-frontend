package sirius

import (
	"fmt"
	"net/url"
)

type apiUser struct {
	ID          int           `json:"id"`
	DisplayName string        `json:"displayName"`
	Teams       []apiUserTeam `json:"teams"`
}

type apiUserTeam struct {
	DisplayName string `json:"displayName"`
}

type User struct {
	ID          int    `json:"id"`
	DisplayName string `json:"displayName"`
}

func (c *Client) SearchUsers(ctx Context, term string) ([]User, error) {
	if len(term) < 3 {
		return nil, fmt.Errorf("Search term must be at least three characters")
	}

	var v []apiUser

	err := c.get(ctx, fmt.Sprintf("/api/v1/search/users?query=%s", url.QueryEscape(term)), &v)
	if err != nil {
		return nil, err
	}

	var users []User
	for _, u := range v {
		user := User{
			ID:          u.ID,
			DisplayName: u.DisplayName,
		}

		if len(u.Teams) > 0 {
			user.DisplayName += fmt.Sprintf(" (%s)", u.Teams[0].DisplayName)
		}

		users = append(users, user)
	}

	return users, nil
}
