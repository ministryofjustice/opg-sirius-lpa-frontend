package sirius

import "fmt"

type Case struct {
	ID               int        `json:"id"`
	UID              string     `json:"uId"`
	CaseType         string     `json:"caseType"`
	CancellationDate DateString `json:"cancellationDate"`
	DispatchDate     DateString `json:"dispatchDate"`
	DueDate          DateString `json:"dueDate"`
	InvalidDate      DateString `json:"invalidDate"`
	ReceiptDate      DateString `json:"receiptDate"`
	RegistrationDate DateString `json:"registrationDate"`
	RejectedDate     DateString `json:"rejectedDate"`
	WithdrawnDate    DateString `json:"withdrawnDate"`
	Children         []Case     `json:"children"`
}

func (c *Client) Case(ctx Context, id int) (Case, error) {
	var v Case
	err := c.get(ctx, fmt.Sprintf("/api/v1/cases/%d", id), &v)

	return v, err
}
