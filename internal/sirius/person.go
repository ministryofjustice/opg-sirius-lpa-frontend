package sirius

import "fmt"

type Person struct {
	ID                    int        `json:"id,omitempty"`
	UID                   string     `json:"uId,omitempty"`
	Salutation            string     `json:"salutation,omitempty"`
	Firstname             string     `json:"firstname"`
	Middlenames           string     `json:"middlenames,omitempty"`
	Surname               string     `json:"surname"`
	DateOfBirth           DateString `json:"dob,omitempty"`
	PreviouslyKnownAs     string     `json:"previousNames,omitempty"`
	AlsoKnownAs           string     `json:"otherNames,omitempty"`
	AddressLine1          string     `json:"addressLine1,omitempty"`
	AddressLine2          string     `json:"addressLine2,omitempty"`
	AddressLine3          string     `json:"addressLine3,omitempty"`
	Town                  string     `json:"town,omitempty"`
	County                string     `json:"county,omitempty"`
	Postcode              string     `json:"postcode,omitempty"`
	Country               string     `json:"country,omitempty"`
	IsAirmailRequired     bool       `json:"isAirmailRequired,omitempty"`
	PhoneNumber           string     `json:"phoneNumber,omitempty"`
	Email                 string     `json:"email,omitempty"`
	SageId                string     `json:"sageId,omitempty"`
	CorrespondenceByPost  bool       `json:"correspondenceByPost,omitempty"`
	CorrespondenceByEmail bool       `json:"correspondenceByEmail,omitempty"`
	CorrespondenceByPhone bool       `json:"correspondenceByPhone,omitempty"`
	CorrespondenceByWelsh bool       `json:"correspondenceByWelsh,omitempty"`
	ResearchOptOut        bool       `json:"researchOptOut,omitempty"`
	Children              []Person   `json:"children,omitempty"`
}

func (p Person) Summary() string {
	return fmt.Sprintf("%s %s", p.Firstname, p.Surname)
}

func (c *Client) Person(ctx Context, id int) (Person, error) {
	var v Person
	err := c.get(ctx, fmt.Sprintf("/api/v1/persons/%d", id), &v)

	return v, err
}
