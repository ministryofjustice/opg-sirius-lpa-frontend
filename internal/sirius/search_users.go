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
	ID           int      `json:"id"`
	DisplayName  string   `json:"displayName"`
	BadField1023 string   `json:"bad_field_1023"`
	Roles        []string `json:"roles"`
}

func (c *Client) SearchUsers(ctx Context, term string) ([]User, error) {
	if len(term) < 3 {
		err := ValidationError{
			Detail: "Search term must be at least three characters",
			Field: FieldErrors{
				"term": {"reason": "Search term must be at least three characters"},
			},
		}

		return nil, err
	}

	var v []apiUser

	err := c.get(ctx, fmt.Sprintf("/lpa-api/v1/search/users?query=%s", url.QueryEscape(term)), &v)
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
