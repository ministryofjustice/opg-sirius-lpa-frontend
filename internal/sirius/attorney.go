package sirius

import (
	"fmt"
	"strings"
)

type Attorney struct {
	Person
	SystemStatus        *bool  `json:"systemStatus,omitempty"`
	RelationshipToDonor string `json:"relationshipToDonor,omitempty"`
}

func (a Attorney) Summary() string {
	return fmt.Sprintf("%s %s", a.Firstname, a.Surname)
}

func (a Attorney) AddressSummary() string {
	i := 0
	s := []string{a.AddressLine1, a.AddressLine2, a.AddressLine3, a.Town, a.County, a.Postcode, a.Country}
	for _, x := range s {
		if x != "" {
			s[i] = x
			i++
		}
	}
	s = s[:i]
	return strings.Join(s, ", ")
}

func (c *Client) CreateAttorney(ctx Context, caseId int, attorney Attorney) error {
	return c.post(ctx, fmt.Sprintf("/lpa-api/v1/epas/%d/attorneys", caseId), attorney, nil)
}
