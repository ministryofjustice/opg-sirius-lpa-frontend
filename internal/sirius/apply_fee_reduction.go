package sirius

import (
	"fmt"
)

const FeeReductionSource = "FEE_REDUCTION"

func (c *Client) ApplyFeeReduction(ctx Context, caseID int, feeReductionType string, paymentEvidence string, paymentDate DateString) error {
	data := struct {
		Source           string     `json:"source"`
		PaymentEvidence  string     `json:"paymentEvidence"`
		FeeReductionType string     `json:"feeReductionType"`
		PaymentDate      DateString `json:"paymentDate"`
	}{
		Source:           FeeReductionSource,
		PaymentEvidence:  paymentEvidence,
		FeeReductionType: feeReductionType,
		PaymentDate:      paymentDate,
	}

	return c.post(ctx, fmt.Sprintf("/lpa-api/v1/cases/%d/payments", caseID), data, nil)
}
