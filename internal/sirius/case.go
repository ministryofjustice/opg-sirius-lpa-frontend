package sirius

import "fmt"

type Case struct {
	UID              string     `json:"uId,omitempty"`
	Status           string     `json:"status"`
	CaseType         string     `json:"caseType,omitempty"`
	CancellationDate DateString `json:"cancellationDate,omitempty"`
	DispatchDate     DateString `json:"dispatchDate,omitempty"`
	DueDate          DateString `json:"dueDate,omitempty"`
	InvalidDate      DateString `json:"invalidDate,omitempty"`
	ReceiptDate      DateString `json:"receiptDate,omitempty"`
	RegistrationDate DateString `json:"registrationDate,omitempty"`
	RejectedDate     DateString `json:"rejectedDate,omitempty"`
	WithdrawnDate    DateString `json:"withdrawnDate,omitempty"`
}

func (c *Client) Case(ctx Context, id int) (Case, error) {
	var v Case
	err := c.get(ctx, fmt.Sprintf("/api/v1/cases/%d", id), &v)

	return v, err
}
