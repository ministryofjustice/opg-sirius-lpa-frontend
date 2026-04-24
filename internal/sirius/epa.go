package sirius

import "fmt"

type Epa struct {
	Case
	AreAllAttorneysApplyingToRegister *bool      `json:"areAllAttorneysApplyingToRegister,omitempty"`
	DonorHasOtherEpas                 *bool      `json:"donorHasOtherEpas,omitempty"`
	EpaDonorNoticeGivenDate           DateString `json:"epaDonorNoticeGivenDate,omitempty"`
	EpaDonorSignatureDate             DateString `json:"epaDonorSignatureDate,omitempty"`
	HasRelativeToNotice               *bool      `json:"hasRelativeToNotice,omitempty"`
	OtherEpaInfo                      string     `json:"otherEpaInfo,omitempty"`
}

func (c *Client) Epa(ctx Context, id int) (Epa, error) {
	var v Epa
	err := c.get(ctx, fmt.Sprintf("/lpa-api/v1/cases/%d", id), &v)

	if v.Donor != nil {
		if v.Donor.Parent != nil {
			v.Donor = v.Donor.Parent
		}
	}

	return v, err
}

func (c *Client) CreateEpa(ctx Context, donorID int, epa Epa) (Epa, error) {
	var v Epa
	err := c.post(ctx, fmt.Sprintf("/lpa-api/v1/donors/%d/epas", donorID), epa, &v)
	return v, err
}

func (c *Client) UpdateEpa(ctx Context, caseId int, epa Epa) error {
	return c.put(ctx, fmt.Sprintf("/lpa-api/v1/epas/%d", caseId), epa, nil)
}
