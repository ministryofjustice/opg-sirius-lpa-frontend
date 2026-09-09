package sirius

import (
	"fmt"
)

type TrustCorporation struct {
	Attorney
	TrustCorporationAppointedAs string `json:"trustCorporationAppointedAs,omitempty"`
	IsReplacementAttorney       bool   `json:"isReplacementAttorney,omitempty"`
}

func (c *Client) CreateTrustCorporation(ctx Context, caseId int, trustCorporation TrustCorporation) error {
	trustCorporation.CaseId = caseId

	return c.post(ctx, "/lpa-api/v1/trust-corporation", trustCorporation, nil)
}

func (c *Client) UpdateTrustCorporation(ctx Context, trustCorporationId int, trustCorporation TrustCorporation) error {
	return c.put(ctx, fmt.Sprintf("/lpa-api/v1/trust-corporation/%d", trustCorporationId), trustCorporation, nil)
}
