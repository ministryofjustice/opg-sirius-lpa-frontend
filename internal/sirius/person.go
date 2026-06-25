package sirius

import (
	"fmt"
	"strings"
)

type Recipient interface {
	Summary() string
	AddressSummary() string
}

type Person struct {
	AddressLine1          string     `json:"addressLine1"`
	AddressLine2          string     `json:"addressLine2"`
	AddressLine3          string     `json:"addressLine3"`
	AlsoKnownAs           string     `json:"otherNames"`
	Cases                 []*Case    `json:"cases,omitempty"`
	Children              []Person   `json:"children,omitempty"`
	CompanyName           string     `json:"companyName"`
	CompanyReference      string     `json:"companyReference"`
	CorrespondenceByEmail bool       `json:"correspondenceByEmail"`
	CorrespondenceByPhone bool       `json:"correspondenceByPhone"`
	CorrespondenceByPost  bool       `json:"correspondenceByPost"`
	CorrespondenceByWelsh bool       `json:"correspondenceByWelsh"`
	Country               string     `json:"country"`
	County                string     `json:"county"`
	DateOfBirth           DateString `json:"dob"`
	DateOfDeath           DateString `json:"dateOfDeath"`
	DisplayName           string     `json:"displayName,omitempty"`
	Email                 string     `json:"email"`
	Firstname             string     `json:"firstname"`
	ID                    int        `json:"id,omitempty"`
	IsAirmailRequired     bool       `json:"isAirmailRequired"`
	Middlenames           string     `json:"middlenames"`
	Parent                *Person    `json:"parent,omitempty"`
	PersonType            string     `json:"personType,omitempty"`
	PhoneNumber           string     `json:"phoneNumber"`
	Postcode              string     `json:"postcode"`
	PreviouslyKnownAs     string     `json:"previousNames"`
	ResearchOptOut        bool       `json:"researchOptOut"`
	SageId                string     `json:"sageId"`
	Salutation            string     `json:"salutation"`
	Surname               string     `json:"surname"`
	Town                  string     `json:"town"`
	UID                   string     `json:"uId,omitempty"`
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
