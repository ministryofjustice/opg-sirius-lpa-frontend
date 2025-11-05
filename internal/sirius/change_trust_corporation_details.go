package sirius

import (
	"fmt"
)

type ChangeTrustCorporationDetails struct {
	Name          string  `json:"name"`
	Address       Address `json:"address"`
	Phone         string  `json:"phoneNumber"`
	Email         string  `json:"email"`
	CompanyNumber string  `json:"companyNumber"`
}

func (c *Client) ChangeTrustCorporationDetails(ctx Context, caseUID string, trustCorpUID string, trustCorpDetails ChangeTrustCorporationDetails) error {
	path := fmt.Sprintf("/lpa-api/v1/digital-lpas/%s/trust-corporation/%s/change-details", caseUID, trustCorpUID)

	return c.put(ctx, path, trustCorpDetails, nil)
}
