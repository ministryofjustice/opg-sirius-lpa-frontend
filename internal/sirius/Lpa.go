package sirius

import "fmt"

type Lpa struct {
	Case
	AdditionalInfo                            string     `json:"additionalInfo,omitempty"`
	AnyOtherInfo                              *bool      `json:"anyOtherInfo,omitempty"`
	ApplicantIds                              []int      `json:"applicantIds,omitempty"`
	ApplicantSignatureDate                    DateString `json:"applicantSignatureDate,omitempty"`
	ApplicantType                             string     `json:"applicantType,omitempty"`
	ApplicationHasGuidance                    *bool      `json:"applicationHasGuidance,omitempty"`
	ApplicationHasRestrictions                *bool      `json:"applicationHasRestrictions,omitempty"`
	AttorneyActDecisions                      string     `json:"attorneyActDecisions,omitempty"`
	CardPaymentContact                        string     `json:"cardPaymentContact,omitempty"`
	CertificateProviderSignature              *bool      `json:"certificateProviderSignature,omitempty"`
	CertificateProviderSignatureDate          DateString `json:"certificateProviderSignatureDate,omitempty"`
	DonorSignatureWitnessed                   *bool      `json:"donorSignatureWitnessed,omitempty"`
	HaveAppliedForFeeRemission                *bool      `json:"haveAppliedForFeeRemission,omitempty"`
	LifeSustainingTreatmentSignedAndWitnessed *bool      `json:"lifeSustainingTreatmentSignedAndWitnessed,omitempty"`
	LpaDonorSignature                         *bool      `json:"lpaDonorSignature,omitempty"`
	OnlineLpaId                               string     `json:"onlineLpaId,omitempty"`
	PaymentByDebitCreditCard                  *bool      `json:"paymentByDebitCreditCard,omitempty"`
	PaymentRemission                          *bool      `json:"paymentRemission,omitempty"`
	RepeatApplication                         *bool      `json:"repeatApplication,omitempty"`
	RepeatApplicationReference                string     `json:"repeatApplicationReference,omitempty"`
}

func (c *Client) Lpa(ctx Context, id int) (Lpa, error) {
	var v Lpa
	err := c.get(ctx, fmt.Sprintf("/lpa-api/v1/cases/%d", id), &v)

	if v.Donor != nil {
		if v.Donor.Parent != nil {
			v.Donor = v.Donor.Parent
		}
	}

	return v, err
}

func (c *Client) CreateLpa(ctx Context, donorID int, lpa Lpa) (Lpa, error) {
	var v Lpa
	err := c.post(ctx, fmt.Sprintf("/lpa-api/v1/donors/%d/lpas", donorID), lpa, &v)
	return v, err
}

func (c *Client) UpdateLpa(ctx Context, caseId int, lpa Lpa) error {
	return c.put(ctx, fmt.Sprintf("/lpa-api/v1/lpas/%d", caseId), lpa, nil)
}
