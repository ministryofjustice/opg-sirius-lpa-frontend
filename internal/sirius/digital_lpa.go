package sirius

import (
	"fmt"
)

type SiriusData struct {
	ID                 int          `json:"id"`
	UID                string       `json:"uId"`
	Application        Draft        `json:"application"`
	Subtype            string       `json:"caseSubtype"`
	CreatedDate        DateString   `json:"createdDate"`
	Status             string       `json:"status"`
	ComplaintCount     int          `json:"complaintCount"`
	InvestigationCount int          `json:"investigationCount"`
	TaskCount          int          `json:"taskCount"`
	WarningCount       int          `json:"warningCount"`
	ObjectionCount     int          `json:"objectionCount"`
	LinkedCases        []DigitalLpa `json:"linkedDigitalLpas"`
	Donor              Donor        `json:"donor"`
}

type DigitalLpa struct {
	UID        string     `json:"uId"`
	SiriusData SiriusData `json:"opg.poas.sirius"`
}

type Donor struct {
	ID int `json:"id"`
}

func (c *Client) DigitalLpa(ctx Context, uid string) (DigitalLpa, error) {
	var v DigitalLpa
	err := c.get(ctx, fmt.Sprintf("/lpa-api/v1/digital-lpas/%s", uid), &v)

	return v, err
}
