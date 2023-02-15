package sirius

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *Client) AddPayment(ctx Context, caseID int, amount int, source string, paymentDate DateString) error {
	postData, err := json.Marshal(struct {
		Amount      int        `json:"amount"`
		Source      string     `json:"source"`
		PaymentDate DateString `json:"paymentDate"`
	}{
		Amount:      amount,
		Source:      source,
		PaymentDate: paymentDate,
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
	defer resp.Body.Close() //#nosec G307 false positive

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
