package sirius

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Dates struct {
	CancellationDate DateString `json:"cancellationDate,omitempty"`
	DispatchDate     DateString `json:"dispatchDate,omitempty"`
	DueDate          DateString `json:"dueDate,omitempty"`
	InvalidDate      DateString `json:"invalidDate,omitempty"`
	PaymentDate      DateString `json:"paymentDate,omitempty"`
	ReceiptDate      DateString `json:"receiptDate,omitempty"`
	RegistrationDate DateString `json:"registrationDate,omitempty"`
	RejectedDate     DateString `json:"rejectedDate,omitempty"`
	WithdrawnDate    DateString `json:"withdrawnDate,omitempty"`
}

func (c *Client) EditDates(ctx Context, caseID int, caseType CaseType, dates Dates) error {
	data, err := json.Marshal(dates)
	if err != nil {
		return err
	}

	req, err := c.newRequest(ctx, http.MethodPut, fmt.Sprintf("/lpa-api/v1/%ss/%d/edit-dates", caseType, caseID), bytes.NewReader(data))
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

	if resp.StatusCode != http.StatusOK {
		return newStatusError(resp)
	}

	return nil
}
