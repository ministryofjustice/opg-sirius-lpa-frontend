package sirius

import (
	"fmt"
)

type DeletedCase struct {
	UID         string     `json:"uId"`
	OnlineLpaId string     `json:"onlineLpaId"`
	Type        string     `json:"type"`
	Status      string     `json:"status"`
	DeletedAt   DateString `json:"deletedAt"`
	Reason      string     `json:"deletionReason"`
}

func (c *Client) DeletedCases(ctx Context, uid string) ([]DeletedCase, error) {
	var v []DeletedCase
	err := c.get(ctx, fmt.Sprintf("/lpa-api/v1/deleted-cases?uid=%s", uid), &v)

	return v, err
}
