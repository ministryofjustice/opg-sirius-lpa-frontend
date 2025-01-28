package server

import (
	"errors"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockCreateDocumentDigitalLpaClient struct {
	mock.Mock
}

func (m *mockCreateDocumentDigitalLpaClient) CaseSummary(ctx sirius.Context, uid string) (sirius.CaseSummary, error) {
	args := m.Called(ctx, uid)
	return args.Get(0).(sirius.CaseSummary), args.Error(1)
}

func (m *mockCreateDocumentDigitalLpaClient) TasksForCase(ctx sirius.Context, id int) ([]sirius.Task, error) {
	args := m.Called(ctx, id)
	return args.Get(0).([]sirius.Task), args.Error(1)
}

func (m *mockCreateDocumentDigitalLpaClient) DocumentTemplates(ctx sirius.Context, caseType sirius.CaseType) ([]sirius.DocumentTemplateData, error) {
	args := m.Called(ctx, caseType)
	return args.Get(0).([]sirius.DocumentTemplateData), args.Error(1)
}

func (m *mockCreateDocumentDigitalLpaClient) CreateContact(ctx sirius.Context, contact sirius.Person) (sirius.Person, error) {
	args := m.Called(ctx, contact)
	return args.Get(0).(sirius.Person), args.Error(1)
}

func (m *mockCreateDocumentDigitalLpaClient) CreateDocument(ctx sirius.Context, caseID, correspondentID int, templateID string, inserts []string) (sirius.Document, error) {
	args := m.Called(ctx, caseID, correspondentID, templateID, inserts)
	return args.Get(0).(sirius.Document), args.Error(1)
}

func TestGetCreateDocumentDigitalLpa(t *testing.T) {
	caseSummary := sirius.CaseSummary{
		DigitalLpa: sirius.DigitalLpa{
			SiriusData: sirius.SiriusData{
				ID: 15,
				Application: sirius.Draft{
					DonorFirstNames: "Zackary",
					DonorLastName:   "Lemmonds",
					DonorAddress: sirius.Address{
						Line1:    "9 Mount Pleasant Drive",
						Town:     "East Harling",
						Postcode: "NR16 2GB",
						Country:  "UK",
					},
				},
			},
			LpaStoreData: sirius.LpaStoreData{
				Donor: sirius.LpaStoreDonor{
					LpaStorePerson: sirius.LpaStorePerson{
						Uid:        "302b05c7-896c-4290-904e-2005e4f1e81e",
						FirstNames: "Zackary",
						LastName:   "Lemmonds",
						Address: sirius.LpaStoreAddress{
							Line1:    "9 Mount Pleasant Drive",
							Town:     "East Harling",
							Postcode: "NR16 2GB",
							Country:  "UK",
						},
					},
					DateOfBirth: "18/04/1965",
				},
				CertificateProvider: sirius.LpaStoreCertificateProvider{
					LpaStorePerson: sirius.LpaStorePerson{
						Uid:        "787ea57e-4db0-4578-8b34-5911e333db03",
						FirstNames: "Stevie",
						LastName:   "Bergnaum",
						Address: sirius.LpaStoreAddress{
							Line1:    "2 Chris Court",
							Town:     "North Beahan",
							Postcode: "TZ08 3QW",
							Country:  "UK",
						},
					},
				},
				Attorneys: []sirius.LpaStoreAttorney{
					{
						LpaStorePerson: sirius.LpaStorePerson{
							Uid:        "0b7976ad-f862-4d51-a6a1-6e25e6386c8d",
							FirstNames: "Zella",
							LastName:   "Zulauf-Jakubowski",
							Address: sirius.LpaStoreAddress{
								Line1:    "86 Labadie Path",
								Line2:    "Gibson-on-Crona",
								Town:     "Gwent",
								Postcode: "KJ94 8SB",
								Country:  "UK",
							},
						},
						Status:          "active",
						AppointmentType: "original",
					},
					{
						LpaStorePerson: sirius.LpaStorePerson{
							Uid:        "af64bee2-4fed-4327-b14f-0ba663927435",
							FirstNames: "Trinity",
							LastName:   "Greenfelder",
							Address: sirius.LpaStoreAddress{
								Line1:    "16 Nolan Gate",
								Line2:    "Wittingborough",
								Line3:    "Old Mann Place",
								Town:     "Central",
								Postcode: "FO7 0ZY",
								Country:  "UK",
							},
						},
						Status:          "removed",
						AppointmentType: "original",
					},
					{
						LpaStorePerson: sirius.LpaStorePerson{
							Uid:        "89a35eb7-6f9a-4108-b2e7-6c6238217728",
							FirstNames: "Devan",
							LastName:   "Purdy-Conn",
							Address: sirius.LpaStoreAddress{
								Line1:    "1 The Drive",
								Line2:    "North Vandervort",
								Line3:    "St. Beer",
								Town:     "Herefordshire",
								Postcode: "YW08 0UQ",
								Country:  "UK",
							},
						},
						Status:          "inactive",
						AppointmentType: "replacement",
					},
				},
			},
		},
		TaskList: []sirius.Task{},
	}

	templateData := []sirius.DocumentTemplateData{{TemplateId: "DL-EXAMPLE", Label: "Example DL Form"}}

	client := &mockCreateDocumentDigitalLpaClient{}
	client.
		On("CaseSummary", mock.Anything, "M-TWGJ-CDDJ-4NTL").
		Return(caseSummary, nil)
	client.
		On("DocumentTemplates", mock.Anything, sirius.CaseTypeDigitalLpa).
		Return(templateData, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, createDocumentDigitalLpaData{
			CaseSummary:           caseSummary,
			DocumentTemplates:     sortDocumentData(templateData),
			ComponentDocumentData: buildComponentDocumentData(templateData),
			Recipients: []sirius.Person{
				{
					ID:           -1,
					UID:          "302b05c7-896c-4290-904e-2005e4f1e81e",
					Firstname:    "Zackary",
					Surname:      "Lemmonds",
					DateOfBirth:  "18/04/1965",
					AddressLine1: "9 Mount Pleasant Drive",
					Town:         "East Harling",
					Postcode:     "NR16 2GB",
					Country:      "UK",
					PersonType:   "Donor",
				},
				{
					ID:           -2,
					UID:          "787ea57e-4db0-4578-8b34-5911e333db03",
					Firstname:    "Stevie",
					Surname:      "Bergnaum",
					AddressLine1: "2 Chris Court",
					Town:         "North Beahan",
					Postcode:     "TZ08 3QW",
					Country:      "UK",
					PersonType:   "Certificate Provider",
				},
				{
					ID:           -3,
					UID:          "0b7976ad-f862-4d51-a6a1-6e25e6386c8d",
					Firstname:    "Zella",
					Surname:      "Zulauf-Jakubowski",
					AddressLine1: "86 Labadie Path",
					AddressLine2: "Gibson-on-Crona",
					Town:         "Gwent",
					Postcode:     "KJ94 8SB",
					Country:      "UK",
					PersonType:   "Attorney",
				},
				{
					ID:           -4,
					UID:          "89a35eb7-6f9a-4108-b2e7-6c6238217728",
					Firstname:    "Devan",
					Surname:      "Purdy-Conn",
					AddressLine1: "1 The Drive",
					AddressLine2: "North Vandervort",
					AddressLine3: "St. Beer",
					Town:         "Herefordshire",
					Postcode:     "YW08 0UQ",
					Country:      "UK",
					PersonType:   "Replacement Attorney",
				},
			},
		}).
		Return(nil)

	server := newMockServer("/lpa/{uid}/documents/new", CreateDocumentDigitalLpa(client, template.Func))

	req, _ := http.NewRequest(http.MethodGet, "/lpa/M-TWGJ-CDDJ-4NTL/documents/new", nil)
	_, err := server.serve(req)

	assert.Nil(t, err)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestGetCreateDocumentDigitalLpaError(t *testing.T) {
	expectedError := errors.New("expected error")

	client := &mockCreateDocumentDigitalLpaClient{}
	client.
		On("CaseSummary", mock.Anything, "M-TWGJ-CDDJ-4NTL").
		Return(sirius.CaseSummary{}, expectedError)

	template := &mockTemplate{}

	server := newMockServer("/lpa/{uid}/documents/new", CreateDocumentDigitalLpa(client, template.Func))

	req, _ := http.NewRequest(http.MethodGet, "/lpa/M-TWGJ-CDDJ-4NTL/documents/new", nil)
	_, err := server.serve(req)

	assert.Equal(t, expectedError, err)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestPostCreateDocumentDigitalLpa(t *testing.T) {
	caseSummary := sirius.CaseSummary{
		DigitalLpa: sirius.DigitalLpa{
			SiriusData: sirius.SiriusData{
				ID: 1344,
				Application: sirius.Draft{
					DonorFirstNames: "Zackary",
					DonorLastName:   "Lemmonds",
					DonorAddress: sirius.Address{
						Line1:    "9 Mount Pleasant Drive",
						Town:     "East Harling",
						Postcode: "NR16 2GB",
						Country:  "UK",
					},
				},
			},
			LpaStoreData: sirius.LpaStoreData{
				Donor: sirius.LpaStoreDonor{
					LpaStorePerson: sirius.LpaStorePerson{
						Uid:        "302b05c7-896c-4290-904e-2005e4f1e81e",
						FirstNames: "Zackary",
						LastName:   "Lemmonds",
						Address: sirius.LpaStoreAddress{
							Line1:    "9 Mount Pleasant Drive",
							Town:     "East Harling",
							Postcode: "NR16 2GB",
							Country:  "UK",
						},
					},
					DateOfBirth: "18/04/1965",
				},
			},
		},
		TaskList: []sirius.Task{},
	}

	templateData := []sirius.DocumentTemplateData{
		{
			TemplateId: "DL-EXAMPLE",
			Label:      "Example DL Form",
			Inserts: []sirius.Insert{
				{InsertId: "DL_INS_01", Label: "Example Insert 1"},
				{InsertId: "DL_INS_02", Label: "Example Insert 2"},
			},
		},
	}

	donorRecipient := sirius.Person{
		ID:           -1,
		UID:          "302b05c7-896c-4290-904e-2005e4f1e81e",
		Firstname:    "Zackary",
		Surname:      "Lemmonds",
		DateOfBirth:  "18/04/1965",
		AddressLine1: "9 Mount Pleasant Drive",
		Town:         "East Harling",
		Postcode:     "NR16 2GB",
		Country:      "UK",
		PersonType:   "Donor",
	}

	contactData := sirius.Person{
		ID:           12,
		Firstname:    "Zackary",
		Surname:      "Lemmonds",
		DateOfBirth:  "18/04/1965",
		AddressLine1: "9 Mount Pleasant Drive",
		Town:         "East Harling",
		Postcode:     "NR16 2GB",
		Country:      "UK",
		PersonType:   "Donor",
	}

	client := &mockCreateDocumentDigitalLpaClient{}
	client.
		On("CaseSummary", mock.Anything, "M-TWGJ-CDDJ-4NTL").
		Return(caseSummary, nil)
	client.
		On("DocumentTemplates", mock.Anything, sirius.CaseTypeDigitalLpa).
		Return(templateData, nil)
	client.
		On("CreateContact", mock.Anything, donorRecipient).
		Return(contactData, nil)
	client.
		On("CreateDocument", mock.Anything, 1344, 12, "DL-EXAMPLE", []string{"DL_INS_01", "DL_INS_02"}).
		Return(sirius.Document{}, nil)

	template := &mockTemplate{}

	server := newMockServer("/lpa/{uid}/documents/new", CreateDocumentDigitalLpa(client, template.Func))

	form := url.Values{
		"templateId":       {"DL-EXAMPLE"},
		"insert":           {"DL_INS_01", "DL_INS_02"},
		"selectRecipients": {"-1"},
	}

	req, _ := http.NewRequest(http.MethodPost, "/lpa/M-TWGJ-CDDJ-4NTL/documents/new", strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", formUrlEncoded)
	_, err := server.serve(req)

	assert.Equal(t, RedirectError("/edit-document?id=1344&case=digital_lpa"), err)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestPostCreateDocumentDigitalLpaInvalid(t *testing.T) {
	caseSummary := sirius.CaseSummary{
		DigitalLpa: sirius.DigitalLpa{
			SiriusData: sirius.SiriusData{
				ID: 1666,
				Application: sirius.Draft{
					DonorFirstNames: "Zackary",
					DonorLastName:   "Lemmonds",
					DonorAddress: sirius.Address{
						Line1:    "9 Mount Pleasant Drive",
						Town:     "East Harling",
						Postcode: "NR16 2GB",
						Country:  "UK",
					},
				},
			},
			LpaStoreData: sirius.LpaStoreData{
				Donor: sirius.LpaStoreDonor{
					LpaStorePerson: sirius.LpaStorePerson{
						Uid:        "302b05c7-896c-4290-904e-2005e4f1e81e",
						FirstNames: "Zackary",
						LastName:   "Lemmonds",
						Address: sirius.LpaStoreAddress{
							Line1:    "9 Mount Pleasant Drive",
							Town:     "East Harling",
							Postcode: "NR16 2GB",
							Country:  "UK",
						},
					},
					DateOfBirth: "18/04/1965",
				},
			},
		},
		TaskList: []sirius.Task{},
	}

	templateData := []sirius.DocumentTemplateData{
		{
			TemplateId: "DL-EXAMPLE",
			Label:      "Example DL Form",
			Inserts: []sirius.Insert{
				{InsertId: "DL_INS_01", Label: "Example Insert 1"},
				{InsertId: "DL_INS_02", Label: "Example Insert 2"},
			},
		},
	}

	client := &mockCreateDocumentDigitalLpaClient{}
	client.
		On("CaseSummary", mock.Anything, "M-TWGJ-CDDJ-4NTL").
		Return(caseSummary, nil)
	client.
		On("DocumentTemplates", mock.Anything, sirius.CaseTypeDigitalLpa).
		Return(templateData, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, createDocumentDigitalLpaData{
			CaseSummary:           caseSummary,
			DocumentTemplates:     sortDocumentData(templateData),
			ComponentDocumentData: buildComponentDocumentData(templateData),
			Recipients: []sirius.Person{{
				ID:           -1,
				UID:          "302b05c7-896c-4290-904e-2005e4f1e81e",
				Firstname:    "Zackary",
				Surname:      "Lemmonds",
				DateOfBirth:  "18/04/1965",
				AddressLine1: "9 Mount Pleasant Drive",
				Town:         "East Harling",
				Postcode:     "NR16 2GB",
				Country:      "UK",
				PersonType:   "Donor",
			}},

			Error: sirius.ValidationError{
				Field: sirius.FieldErrors{
					"templateId":      {"reason": "Please select a document template to continue"},
					"selectRecipient": {"reason": "Please select a recipient to continue"},
				},
			},
		}).
		Return(nil)

	server := newMockServer("/lpa/{uid}/documents/new", CreateDocumentDigitalLpa(client, template.Func))

	form := url.Values{}

	req, _ := http.NewRequest(http.MethodPost, "/lpa/M-TWGJ-CDDJ-4NTL/documents/new", strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", formUrlEncoded)
	_, err := server.serve(req)

	assert.Nil(t, err)
	mock.AssertExpectationsForObjects(t, client, template)
}
