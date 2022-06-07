package sirius

import (
	"fmt"
)

type PersonReference struct {
	ReferenceID int    `json:"referenceId"`
	ID          int    `json:"id"`
	UID         int    `json:"uid"`
	DisplayName string `json:"displayName"`
	Reason      string `json:"reason"`
}

func (c *Client) PersonReferences(ctx Context, personID int) ([]PersonReference, error) {
	var v []PersonReference
	err := c.get(ctx, fmt.Sprintf("/api/v1/persons/%d/references", personID), &v)

	return v, err
}
