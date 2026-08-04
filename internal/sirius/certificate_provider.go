package sirius

import (
	"fmt"
)

func (c *Client) CreateCertificateProvider(ctx Context, caseId int, certificateProvider Person) error {
	certificateProvider.CaseId = caseId
	certificateProvider.PersonType = "CertificateProvider"
	return c.post(ctx, "/lpa-api/v1/persons", []Person{certificateProvider}, nil)
}

func (c *Client) UpdateCertificateProvider(ctx Context, certificateProviderId int, certificateProvider Person) error {
	return c.put(ctx, fmt.Sprintf("/lpa-api/v1/persons/%d", certificateProviderId), certificateProvider, nil)
}
