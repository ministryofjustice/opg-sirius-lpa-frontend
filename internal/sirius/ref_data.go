package sirius

import (
	"fmt"
)

const (
	AttorneyRemovedReasonCategory string = "attorneyRemovedReason"
	CaseStatusChangeReason        string = "caseChangeReason"
	CompensationType              string = "compensationType"
	ComplainantCategory           string = "complainantCategory"
	ComplaintCategory             string = "complaintCategory"
	ComplaintOrigin               string = "complaintOrigin"
	CountryCategory               string = "country"
	DocumentTemplateIdCategory    string = "documentTemplateId"
	FeeDecisionTypeCategory       string = "feeDecisionType"
	FeeReductionTypeCategory      string = "feeReductionType"
	PaymentReferenceType          string = "paymentReferenceType"
	PaymentSourceCategory         string = "paymentSource"
	RelationshipToDonorCategory   string = "relationshipToDonor"
	WarningTypeCategory           string = "warningType"
)

type RefDataItem struct {
	Handle         string        `json:"handle"`
	Label          string        `json:"label"`
	UserSelectable bool          `json:"userSelectable"`
	Subcategories  []RefDataItem `json:"subcategories"`
	ParentSources  []string      `json:"parentSources"`
	ValidSubTypes  []string      `json:"validSubTypes"`
}

func (c *Client) RefDataByCategory(ctx Context, category string) ([]RefDataItem, error) {
	var v []RefDataItem

	if cached, ok := getCached(category); ok {
		return cached, nil
	}

	err := c.get(ctx, fmt.Sprintf("/lpa-api/v1/reference-data/%s", category), &v)

	setCached(category, v)

	return v, err
}
