package sirius

import (
	"fmt"
)

const (
	PaymentSourceCategory      string = "paymentSource"
	WarningTypeCategory        string = "warningType"
	FeeReductionTypeCategory   string = "feeReductionType"
	PaymentReferenceType       string = "paymentReferenceType"
	DocumentTemplateIdCategory string = "documentTemplateId"
	ComplainantCategory        string = "complainantCategory"
	ComplaintOrigin            string = "complaintOrigin"
	CompensationType           string = "compensationType"
	ComplaintCategory          string = "complaintCategory"
)

type RefDataItem struct {
	Handle         string        `json:"handle"`
	Label          string        `json:"label"`
	UserSelectable bool          `json:"userSelectable"`
	Subcategories  []RefDataItem `json:"subcategories"`
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
