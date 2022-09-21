package sirius

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const FeeReductionSource = "FEE_REDUCTION"

func (c *Client) AddPayment(ctx Context, caseID int, amount int, source string, paymentDate DateString, feeReductionType string, paymentEvidence string, appliedDate DateString) error {
	postData, err := GetPostData(amount, source, paymentDate, feeReductionType, paymentEvidence, appliedDate)

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

func GetPostData(amount int, source string, paymentDate DateString, feeReductionType string, paymentEvidence string, appliedDate DateString) ([]byte, error) {
	if source == FeeReductionSource {
		return json.Marshal(struct {
			Source           string     `json:"source"`
			PaymentEvidence  string     `json:"paymentEvidence"`
			FeeReductionType string     `json:"feeReductionType"`
			AppliedDate      DateString `json:"appliedDate"`
		}{
			Source:           source,
			PaymentEvidence:  paymentEvidence,
			FeeReductionType: feeReductionType,
			AppliedDate:      appliedDate,
		})
	} else {
		return json.Marshal(struct {
			Amount      int        `json:"amount"`
			Source      string     `json:"source"`
			PaymentDate DateString `json:"paymentDate"`
		}{
			Amount:      amount,
			Source:      source,
			PaymentDate: paymentDate,
		})
	}
}
