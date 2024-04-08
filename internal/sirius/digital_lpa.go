package sirius

import (
	"fmt"
)

type DigitalLpa struct {
	UID          string       `json:"uId"`
	SiriusData   SiriusData   `json:"opg.poas.sirius"`
	LpaStoreData LpaStoreData `json:"opg.poas.lpastore"`
}

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
	LinkedCases        []SiriusData `json:"linkedDigitalLpas"`
	Donor              Donor        `json:"donor"`
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

type LpaStoreData struct {
	Donor                                       LpaStoreDonor               `json:"donor"`
	Attorneys                                   []LpaStoreAttorney          `json:"attorneys"`
	CertificateProvider                         LpaStoreCertificateProvider `json:"certificateProvider"`
	PeopleToNotify                              []LpaStorePersonToNotify    `json:"peopleToNotify"`
	HowAttorneysMakeDecisions                   string                      `json:"howAttorneysMakeDecisions"`
	HowAttorneysMakeDecisionsDetails            string                      `json:"howAttorneysMakeDecisionsDetails"`
	WhenTheLpaCanBeUsed                         string                      `json:"whenTheLpaCanBeUsed"`
	HowReplacementAttorneysMakeDecisions        string                      `json:"howReplacementAttorneysMakeDecisions"`
	HowReplacementAttorneysMakeDecisionsDetails string                      `json:"howReplacementAttorneysMakeDecisionsDetails"`
	HowReplacementAttorneysStepIn               string                      `json:"howReplacementAttorneysStepIn"`
	HowReplacementAttorneysStepInDetails        string                      `json:"howReplacementAttorneysStepInDetails"`
	RestrictionsAndConditions                   string                      `json:"restrictionsAndConditions"`
	SignedAt                                    string                      `json:"signedAt"`
}

type LpaStoreDonor struct {
	LpaStorePerson
	DateOfBirth       string `json:"dateOfBirth"`
	OtherNamesKnownBy string `json:"otherNamesKnownBy"`
}

type LpaStoreAttorney struct {
	LpaStorePerson
	DateOfBirth string `json:"dateOfBirth"`
	Status      string `json:"status"`
	Mobile      string `json:"mobile"`
}

type LpaStoreCertificateProvider struct {
	LpaStorePerson
	Channel                   string `json:"channel"`
	ContactLanguagePreference string `json:"contactLanguagePreference"`
}

type LpaStorePersonToNotify struct {
	LpaStorePerson
}

type LpaStorePerson struct {
	FirstNames string          `json:"firstNames"`
	LastName   string          `json:"lastName"`
	Address    LpaStoreAddress `json:"address"`
	Email      string          `json:"email"`
}

type LpaStoreAddress struct {
	Line1    string `json:"line1"`
	Line2    string `json:"line2"`
	Line3    string `json:"line3"`
	Town     string `json:"town"`
	Postcode string `json:"postcode"`
	Country  string `json:"country"`
}

func (c *Client) DigitalLpa(ctx Context, uid string) (DigitalLpa, error) {
	var v DigitalLpa
	err := c.get(ctx, fmt.Sprintf("/lpa-api/v1/digital-lpas/%s", uid), &v)
	return v, err
}
