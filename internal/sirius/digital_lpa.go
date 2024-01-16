package sirius

import (
	"encoding/json"
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
	UID          string          `json:"uId"`
	SiriusData   SiriusData      `json:"opg.poas.sirius"`
	LpaStoreData json.RawMessage `json:"opg.poas.lpastore"`
}

type Donor struct {
	ID           int        `json:"id"`
	Firstname    string     `json:"firstname"`
	Surname      string     `json:"surname"`
	DateOfBirth  DateString `json:"dob"`
	AddressLine1 string     `json:"addressLine1"`
	AddressLine2 string     `json:"addressLine2"`
	AddressLine3 string     `json:"addressLine3"`
	Town         string     `json:"town"`
	Postcode     string     `json:"postcode"`
	Country      string     `json:"country"`
	PersonType   string     `json:"personType,omitempty"`
}

func (c *Client) DigitalLpa(ctx Context, uid string) (DigitalLpa, error) {
	var v DigitalLpa
	err := c.get(ctx, fmt.Sprintf("/lpa-api/v1/digital-lpas/%s", uid), &v)

	return v, err
}
