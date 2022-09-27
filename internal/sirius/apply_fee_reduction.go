package sirius

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const FeeReductionSource = "FEE_REDUCTION"

func (c *Client) ApplyFeeReduction(ctx Context, caseID int, feeReductionType string, paymentEvidence string, paymentDate DateString) error {
	postData, err := json.Marshal(struct {
		Source           string     `json:"source"`
		PaymentEvidence  string     `json:"paymentEvidence"`
		FeeReductionType string     `json:"feeReductionType"`
		PaymentDate      DateString `json:"paymentDate"`
	}{
		Source:           FeeReductionSource,
		PaymentEvidence:  paymentEvidence,
		FeeReductionType: feeReductionType,
		PaymentDate:      paymentDate,
	})

	if err != nil {
		return err
	}

	req, err := c.newRequest(ctx, http.MethodPost, fmt.Sprintf("/lpa-api/v1/cases/%d/payments", caseID), bytes.NewReader(postData))
	if err != nil {
		return err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusBadRequest {
		var v ValidationError
		if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
			return err
		}
		return v
	}

	if resp.StatusCode != http.StatusCreated {
		return newStatusError(resp)
	}

	return nil
}
