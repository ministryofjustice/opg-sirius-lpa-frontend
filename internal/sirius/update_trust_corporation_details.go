package sirius

import (
	"fmt"
)

type UpdateTrustCorporationDetails struct {
	Name          string  `json:"name"`
	Address       Address `json:"address"`
	Phone         string  `json:"phoneNumber"`
	Email         string  `json:"email"`
	CompanyNumber string  `json:"companyNumber"`
}

func (c *Client) UpdateTrustCorporationDetails(ctx Context, caseUID string, trustCorpUID string, trustCorpDetails UpdateTrustCorporationDetails) error {
	path := fmt.Sprintf("/lpa-api/v1/digital-lpas/%s/trust-corporation/%s/change-details", caseUID, trustCorpUID)

	return c.put(ctx, path, trustCorpDetails, nil)
}
