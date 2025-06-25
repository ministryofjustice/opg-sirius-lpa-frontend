package sirius

type Address struct {
	Line1    string `json:"addressLine1"`
	Line2    string `json:"addressLine2,omitempty"`
	Line3    string `json:"addressLine3,omitempty"`
	Town     string `json:"town"`
	Postcode string `json:"postcode"`
	Country  string `json:"country"`
}

type DonorIdentityCheck struct {
	State     string `json:"state,omitempty"`
	CheckedAt string `json:"checkedAt,omitempty"`
	Reference string `json:"reference,omitempty"`
}

type Draft struct {
	CaseType                  []string              `json:"types"`
	Source                    string                `json:"source"`
	DonorFirstNames           string                `json:"donorFirstNames"`
	DonorLastName             string                `json:"donorLastName"`
	DonorDob                  DateString            `json:"donorDob"`
	DonorAddress              Address               `json:"donorAddress"`
	CorrespondentFirstNames   string                `json:"correspondentFirstNames,omitempty"`
	CorrespondentLastName     string                `json:"correspondentLastName,omitempty"`
	CorrespondentAddress      *Address              `json:"correspondentAddress,omitempty"`
	PhoneNumber               string                `json:"donorPhone,omitempty"`
	Email                     string                `json:"donorEmail,omitempty"`
	CorrespondenceByWelsh     bool                  `json:"correspondenceByWelsh,omitempty"`
	CorrespondenceLargeFormat bool                  `json:"correspondenceLargeFormat,omitempty"`
	SeveranceStatus           string                `json:"severanceStatus,omitempty"`
	SeveranceApplication      *SeveranceApplication `json:"severanceApplication,omitempty"`
	DonorIdentityCheck        *DonorIdentityCheck   `json:"donorIdentityCheck,omitempty"`
}

func (c *Client) CreateDraft(ctx Context, draft Draft) (map[string]string, error) {
	var v []struct {
		Subtype string `json:"caseSubtype"`
		Uid     string `json:"uId"`
	}
	err := c.post(ctx, "/lpa-api/v1/digital-lpas", draft, &v)

	out := map[string]string{}
	for _, lpa := range v {
		out[lpa.Subtype] = lpa.Uid
	}

	return out, err
}
