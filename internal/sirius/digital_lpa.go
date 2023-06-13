package sirius

import (
	"fmt"
	"log"
)

type DigitalLpa struct {
	UID                string     `json:"uId"`
	Application        Draft      `json:"application"`
	Subtype            string     `json:"caseSubtype"`
	CreatedDate        DateString `json:"createdDate"`
	Status             string     `json:"status"`
	ComplaintCount     int        `json:"complaintCount"`
	InvestigationCount int        `json:"investigationCount"`
	TaskCount          int        `json:"taskCount"`
	WarningCount       int        `json:"warningCount"`
}

func (c *Client) DigitalLpa(ctx Context, uid string) (DigitalLpa, error) {
	var v DigitalLpa
	err := c.get(ctx, fmt.Sprintf("/lpa-api/v1/digital-lpas/%s", uid), &v)

	log.Print(v)

	return v, err
}