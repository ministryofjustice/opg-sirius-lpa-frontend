package sirius

import "fmt"

type Person struct {
	ID                    int        `json:"id,omitempty"`
	UID                   string     `json:"uId,omitempty"`
	Salutation            string     `json:"salutation"`
	Firstname             string     `json:"firstname"`
	Middlenames           string     `json:"middlenames"`
	Surname               string     `json:"surname"`
	DateOfBirth           DateString `json:"dob"`
	PreviouslyKnownAs     string     `json:"previousNames"`
	AlsoKnownAs           string     `json:"otherNames"`
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
	ResearchOptOut        bool       `json:"researchOptOut"`
	Children              []Person   `json:"children,omitempty"`
}

func (p Person) Summary() string {
	return fmt.Sprintf("%s %s", p.Firstname, p.Surname)
}

func (c *Client) Person(ctx Context, id int) (Person, error) {
	var v Person
	err := c.get(ctx, fmt.Sprintf("/lpa-api/v1/persons/%d", id), &v)

	return v, err
}
