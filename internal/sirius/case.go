package sirius

import (
	"fmt"

	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/shared"
)

type Case struct {
	Applicants                                []Person          `json:"applicants,omitempty"`
	Assignee                                  *Person           `json:"assignee,omitempty"`
	Attorneys                                 []Person          `json:"attorneys,omitempty"`
	CancellationDate                          DateString        `json:"cancellationDate,omitempty"`
	CaseAttorneyJointly                       bool              `json:"caseAttorneyJointly,omitempty"`
	CaseAttorneyJointlyAndJointlyAndSeverally bool              `json:"caseAttorneyJointlyAndJointlyAndSeverally,omitempty"`
	CaseAttorneyJointlyAndSeverally           bool              `json:"caseAttorneyJointlyAndSeverally,omitempty"`
	CaseAttorneySingular                      bool              `json:"caseAttorneySingular,omitempty"`
	CaseType                                  string            `json:"caseType,omitempty"`
	Complaints                                []interface{}     `json:"complaints,omitempty"`
	Correspondent                             *Person           `json:"correspondent,omitempty"`
	DispatchDate                              DateString        `json:"dispatchDate,omitempty"`
	Documents                                 []interface{}     `json:"documents,omitempty"`
	Donor                                     *Person           `json:"donor,omitempty"`
	DonorHasOtherEpas                         bool              `json:"donorHasOtherEpas,omitempty"`
	DueDate                                   DateString        `json:"dueDate,omitempty"`
	ExpectedPaymentTotal                      int               `json:"expectedPaymentTotal"`
	FilingDate                                DateString        `json:"filingDate,omitempty"`
	ID                                        int               `json:"id,omitempty"`
	InvalidDate                               DateString        `json:"invalidDate,omitempty"`
	Investigations                            []interface{}     `json:"investigations,omitempty"`
	NotifiedPersons                           []Person          `json:"notifiedPersons,omitempty"`
	Notes                                     []interface{}     `json:"notes,omitempty"`
	OtherEpaInfo                              string            `json:"otherEpaInfo,omitempty"`
	PaymentDate                               DateString        `json:"paymentDate,omitempty"`
	RagRating                                 int               `json:"ragRating,omitempty"`
	ReceiptDate                               DateString        `json:"receiptDate,omitempty"`
	RegistrationDate                          DateString        `json:"registrationDate,omitempty"`
	RejectedDate                              DateString        `json:"rejectedDate,omitempty"`
	ReplacementAttorneys                      []Person          `json:"replacementAttorneys,omitempty"`
	RevokedDate                               DateString        `json:"revokedDate,omitempty"`
	Status                                    shared.CaseStatus `json:"status"`
	StatusDate                                DateString        `json:"statusDate,omitempty"`
	SubType                                   string            `json:"caseSubtype,omitempty"`
	Tasks                                     []interface{}     `json:"tasks,omitempty"`
	TrustCorporations                         []Person          `json:"trustCorporations,omitempty"`
	UID                                       string            `json:"uId,omitempty"`
	ValidationChecks                          []interface{}     `json:"validationChecks,omitempty"`
	WithdrawnDate                             DateString        `json:"withdrawnDate,omitempty"`
}

func (c Case) Summary() string {
	return fmt.Sprintf("%s %s", c.CaseType, c.UID)
}

func (c *Client) Case(ctx Context, id int) (Case, error) {
	var v Case
	err := c.get(ctx, fmt.Sprintf("/lpa-api/v1/cases/%d", id), &v)

	if v.Donor != nil {
		if v.Donor.Parent != nil {
			v.Donor = v.Donor.Parent
		}
	}

	return v, err
}

func (c Case) FilterInactiveAttorneys() Case {
	var activeAttorneys []Person
	var activeTrustCorps []Person

	for _, attorney := range c.Attorneys {
		if attorney.SystemStatus {
			activeAttorneys = append(activeAttorneys, attorney)
		}
	}

	for _, trustCorp := range c.TrustCorporations {
		if trustCorp.SystemStatus {
			activeTrustCorps = append(activeTrustCorps, trustCorp)
		}
	}

	c.Attorneys = activeAttorneys
	c.TrustCorporations = activeTrustCorps
	return c
}
