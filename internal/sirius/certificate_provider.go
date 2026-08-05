package sirius

func (c *Client) CreateCertificateProvider(ctx Context, caseId int, certificateProvider Person) error {
	certificateProvider.CaseId = caseId
	certificateProvider.PersonType = "CertificateProvider"
	return c.post(ctx, "/lpa-api/v1/persons", []Person{certificateProvider}, nil)
}
