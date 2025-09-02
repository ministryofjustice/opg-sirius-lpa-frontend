package sirius

import (
	"fmt"

	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/shared"
)

type DigitalLpa struct {
	UID          string       `json:"uId"`
	SiriusData   SiriusData   `json:"opg.poas.sirius"`
	LpaStoreData LpaStoreData `json:"opg.poas.lpastore"`
}

type SiriusData struct {
	ID                 int               `json:"id"`
	UID                string            `json:"uId"`
	Application        Draft             `json:"application"`
	Subtype            string            `json:"caseSubtype"`
	CreatedDate        DateString        `json:"createdDate"`
	Status             shared.CaseStatus `json:"status"`
	ComplaintCount     int               `json:"complaintCount"`
	InvestigationCount int               `json:"investigationCount"`
	TaskCount          int               `json:"taskCount"`
	WarningCount       int               `json:"warningCount"`
	ObjectionCount     int               `json:"objectionCount"`
	LinkedCases        []SiriusData      `json:"linkedDigitalLpas"`
	Donor              Donor             `json:"donor"`
	DueDate            DateString        `json:"dueDate"`
	StatusColour       string
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
	Phone        string     `json:"phone,omitempty"`
	Email        string     `json:"email,omitempty"`
}

type LpaStoreData struct {
	Donor                                       LpaStoreDonor               `json:"donor"`
	Channel                                     string                      `json:"channel"`
	Status                                      shared.CaseStatus           `json:"status"`
	Attorneys                                   []LpaStoreAttorney          `json:"attorneys"`
	TrustCorporations                           []LpaStoreTrustCorporation  `json:"trustCorporations"`
	CertificateProvider                         LpaStoreCertificateProvider `json:"certificateProvider"`
	CertificateProviderNotRelatedConfirmedAt    string                      `json:"certificateProviderNotRelatedConfirmedAt"`
	PeopleToNotify                              []LpaStorePersonToNotify    `json:"peopleToNotify"`
	HowAttorneysMakeDecisions                   string                      `json:"howAttorneysMakeDecisions"`
	HowAttorneysMakeDecisionsDetails            string                      `json:"howAttorneysMakeDecisionsDetails"`
	WhenTheLpaCanBeUsed                         string                      `json:"whenTheLpaCanBeUsed"`
	HowReplacementAttorneysMakeDecisions        string                      `json:"howReplacementAttorneysMakeDecisions"`
	HowReplacementAttorneysMakeDecisionsDetails string                      `json:"howReplacementAttorneysMakeDecisionsDetails"`
	HowReplacementAttorneysStepIn               string                      `json:"howReplacementAttorneysStepIn"`
	HowReplacementAttorneysStepInDetails        string                      `json:"howReplacementAttorneysStepInDetails"`
	LifeSustainingTreatmentOption               string                      `json:"lifeSustainingTreatmentOption"`
	RestrictionsAndConditions                   string                      `json:"restrictionsAndConditions"`
	RestrictionsAndConditionsImages             []LpaStoreImage             `json:"restrictionsAndConditionsImages"`
	SignedAt                                    string                      `json:"signedAt"`
}

type ActorIdentityCheck struct {
	CheckedAt string `json:"checkedAt"`
	Type      string `json:"type"`
}

type LpaStoreDonor struct {
	LpaStorePerson
	DateOfBirth               string              `json:"dateOfBirth"`
	OtherNamesKnownBy         string              `json:"otherNamesKnownBy"`
	ContactLanguagePreference string              `json:"contactLanguagePreference"`
	IdentityCheck             *ActorIdentityCheck `json:"identityCheck"`
}

type LpaStoreAttorney struct {
	LpaStorePerson
	DateOfBirth               string `json:"dateOfBirth"`
	Status                    string `json:"status"`
	AppointmentType           string `json:"appointmentType"`
	Mobile                    string `json:"mobile"`
	ContactLanguagePreference string `json:"contactLanguagePreference"`
	SignedAt                  string `json:"signedAt"`
	Email                     string `json:"email"`
	Decisions                 bool   `json:"cannotMakeJointDecisions,omitempty"`
}

type LpaStoreTrustCorporation struct {
	LpaStoreAttorney
	Name          string      `json:"name"`
	CompanyNumber string      `json:"companyNumber"`
	Signatories   []Signatory `json:"signatories,omitempty"`
}

type Signatory struct {
	FirstNames        string `json:"firstNames"`
	LastName          string `json:"lastName"`
	ProfessionalTitle string `json:"professionalTitle"`
	SignedAt          string `json:"signedAt"`
}

type LpaStoreCertificateProvider struct {
	LpaStorePerson
	Phone                     string `json:"phone"`
	Channel                   string `json:"channel"`
	ContactLanguagePreference string `json:"contactLanguagePreference"`
	SignedAt                  string `json:"signedAt"`
}

type LpaStorePersonToNotify struct {
	LpaStorePerson
}

type LpaStorePerson struct {
	Uid        string          `json:"uid"`
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

type LpaStoreImage struct {
	Path string `json:"path"`
}

func (c *Client) DigitalLpa(ctx Context, uid string, presignImages bool) (DigitalLpa, error) {
	var v DigitalLpa
	url := fmt.Sprintf("/lpa-api/v1/digital-lpas/%s", uid)
	if presignImages {
		url += "?presignImages"
	}
	err := c.get(ctx, url, &v)
	return v, err
}
