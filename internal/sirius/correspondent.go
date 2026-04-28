package sirius

import (
	"fmt"
	"slices"
	"strings"
)

type Correspondent struct {
	Person
	SystemStatus        *bool  `json:"systemStatus,omitempty"`
	RelationshipToDonor string `json:"relationshipToDonor,omitempty"`
	CaseId              int    `json:"caseId,omitempty"`
}

func (c Correspondent) Summary() string {
	return fmt.Sprintf("%s %s", c.Firstname, c.Surname)
}

func (c Correspondent) AddressSummary() string {
	address := []string{c.AddressLine1, c.AddressLine2, c.AddressLine3, c.Town, c.County, c.Postcode, c.Country}
	filteredAddress := slices.DeleteFunc(address, func(x string) bool { return x == "" })
	return strings.Join(filteredAddress, ", ")
}

func (c *Client) CreateCorrespondent(ctx Context, caseId int, correspondent Correspondent) error {
	correspondent.CaseId = caseId
	correspondent.PersonType = "Correspondent"
	return c.post(ctx, "/lpa-api/v1/persons", []Correspondent{correspondent}, nil)
}

func (c *Client) UpdateCorrespondent(ctx Context, correspondentId int, correspondent Correspondent) error {
	return c.put(ctx, fmt.Sprintf("/lpa-api/v1/persons/%d", correspondentId), correspondent, nil)
}
