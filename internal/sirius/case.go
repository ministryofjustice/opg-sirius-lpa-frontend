package sirius

import "fmt"

type Case struct {
	UID      string `json:"uId"`
	CaseType string `json:"caseType"`
}

func (c *Client) Case(ctx Context, id int) (Case, error) {
	var v Case
	err := c.get(ctx, fmt.Sprintf("/api/v1/cases/%d", id), &v)

	return v, err
}
