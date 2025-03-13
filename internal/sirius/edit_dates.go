package sirius

import (
	"fmt"
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
	return c.put(ctx, fmt.Sprintf("/lpa-api/v1/%ss/%d/edit-dates", caseType, caseID), dates, nil)
}
