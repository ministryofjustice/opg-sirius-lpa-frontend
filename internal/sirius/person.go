package sirius

import (
	"fmt"
	"strings"
)

type Person struct {
	ID                    int        `json:"id,omitempty"`
	UID                   string     `json:"uId,omitempty"`
	Salutation            string     `json:"salutation"`
	Firstname             string     `json:"firstname"`
	Middlenames           string     `json:"middlenames"`
	Surname               string     `json:"surname"`
	DateOfBirth           DateString `json:"dob,omitempty"`
	PreviouslyKnownAs     string     `json:"previousNames,omitempty"`
	AlsoKnownAs           string     `json:"otherNames,omitempty"`
	AddressLine1          string     `json:"addressLine1"`
	AddressLine2          string     `json:"addressLine2"`
	AddressLine3          string     `json:"addressLine3"`
	Town                  string     `json:"town"`
	County                string     `json:"county"`
	Postcode              string     `json:"postcode"`
	Country               string     `json:"country"`
	IsAirmailRequired     bool       `json:"isAirmailRequired"`
	PhoneNumber           string     `json:"phoneNumber"`
	Email                 string     `json:"email"`
	SageId                string     `json:"sageId"`
	CorrespondenceByPost  bool       `json:"correspondenceByPost"`
	CorrespondenceByEmail bool       `json:"correspondenceByEmail"`
	CorrespondenceByPhone bool       `json:"correspondenceByPhone"`
	CorrespondenceByWelsh bool       `json:"correspondenceByWelsh"`
	ResearchOptOut        bool       `json:"researchOptOut,omitempty"`
	Children              []Person   `json:"children,omitempty"`
	CompanyName           string     `json:"companyName,omitempty"`
	CompanyReference      string     `json:"companyReference,omitempty"`
	PersonType            string     `json:"personType,omitempty"`
	Cases                 []*Case    `json:"cases,omitempty"`
}

func (p Person) Summary() string {
	return fmt.Sprintf("%s %s", p.Firstname, p.Surname)
}

func (p Person) AddressSummary() string {
	i := 0
	s := []string{p.AddressLine1, p.AddressLine2, p.AddressLine3, p.Town, p.County, p.Postcode, p.Country}
	for _, x := range s {
		if x != "" {
			s[i] = x
			i++
		}
	}
	s = s[:i]
	return strings.Join(s, ", ")
}

func (c *Client) Person(ctx Context, id int) (Person, error) {
	var v Person
	err := c.get(ctx, fmt.Sprintf("/lpa-api/v1/persons/%d", id), &v)

	return v, err
}
