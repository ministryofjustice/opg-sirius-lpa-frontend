package sirius

import (
	"fmt"
	"slices"
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
	address := []string{a.AddressLine1, a.AddressLine2, a.AddressLine3, a.Town, a.County, a.Postcode, a.Country}
	filteredAddress := slices.DeleteFunc(address, func(x string) bool { return x == "" })
	return strings.Join(filteredAddress, ", ")
}

func (c *Client) CreateAttorney(ctx Context, caseId int, attorney Attorney) error {
	return c.post(ctx, fmt.Sprintf("/lpa-api/v1/epas/%d/attorneys", caseId), attorney, nil)
}

func (c *Client) UpdateAttorney(ctx Context, attorneyId int, attorney Attorney) error {
	return c.put(ctx, fmt.Sprintf("/lpa-api/v1/attorneys/%d", attorneyId), attorney, nil)
}
