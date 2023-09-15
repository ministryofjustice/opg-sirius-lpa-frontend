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

func (m *mockCreateDocumentDigitalLpaClient) DigitalLpa(ctx sirius.Context, uid string) (sirius.DigitalLpa, error) {
	args := m.Called(ctx, uid)
	return args.Get(0).(sirius.DigitalLpa), args.Error(1)
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
	digitalLpa := sirius.DigitalLpa{
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
	}

	templateData := []sirius.DocumentTemplateData{{TemplateId: "DL-EXAMPLE", Label: "Example DL Form"}}

	client := &mockCreateDocumentDigitalLpaClient{}
	client.
		On("DigitalLpa", mock.Anything, "M-TWGJ-CDDJ-4NTL").
		Return(digitalLpa, nil)
	client.
		On("DocumentTemplates", mock.Anything, sirius.CaseTypeDigitalLpa).
		Return(templateData, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, createDocumentDigitalLpaData{
			Lpa:                   digitalLpa,
			DocumentTemplates:     sortDocumentData(templateData),
			ComponentDocumentData: buildComponentDocumentData(templateData),
			Recipients: []sirius.Person{{
				ID:           -1,
				Firstname:    "Zackary",
				Surname:      "Lemmonds",
				PersonType:   "Donor",
				AddressLine1: "9 Mount Pleasant Drive",
				Town:         "East Harling",
				Postcode:     "NR16 2GB",
				Country:      "UK",
			}},
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
	templateData := []sirius.DocumentTemplateData{{TemplateId: "DL-EXAMPLE", Label: "Example DL Form"}}

	client := &mockCreateDocumentDigitalLpaClient{}
	client.
		On("DigitalLpa", mock.Anything, "M-TWGJ-CDDJ-4NTL").
		Return(sirius.DigitalLpa{}, expectedError)
	client.
		On("DocumentTemplates", mock.Anything, sirius.CaseTypeDigitalLpa).
		Return(templateData, nil)

	template := &mockTemplate{}

	server := newMockServer("/lpa/{uid}/documents/new", CreateDocumentDigitalLpa(client, template.Func))

	req, _ := http.NewRequest(http.MethodGet, "/lpa/M-TWGJ-CDDJ-4NTL/documents/new", nil)
	_, err := server.serve(req)

	assert.Equal(t, expectedError, err)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestPostCreateDocumentDigitalLpa(t *testing.T) {
	digitalLpa := sirius.DigitalLpa{
		ID: 1234,
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
		On("DigitalLpa", mock.Anything, "M-TWGJ-CDDJ-4NTL").
		Return(digitalLpa, nil)
	client.
		On("DocumentTemplates", mock.Anything, sirius.CaseTypeDigitalLpa).
		Return(templateData, nil)
	client.
		On("CreateContact", mock.Anything, sirius.Person{
			ID:           -1,
			Firstname:    "Zackary",
			Surname:      "Lemmonds",
			PersonType:   "Donor",
			AddressLine1: "9 Mount Pleasant Drive",
			Town:         "East Harling",
			Postcode:     "NR16 2GB",
			Country:      "UK",
		}).
		Return(sirius.Person{ID: 4829}, nil)
	client.
		On("CreateDocument", mock.Anything, 1234, 4829, "DL-EXAMPLE", []string{"DL_INS_01", "DL_INS_02"}).
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

	assert.Equal(t, RedirectError("/edit-document?id=1234&case=digital_lpa"), err)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestPostCreateDocumentDigitalLpaInvalid(t *testing.T) {
	digitalLpa := sirius.DigitalLpa{
		ID: 1234,
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
		On("DigitalLpa", mock.Anything, "M-TWGJ-CDDJ-4NTL").
		Return(digitalLpa, nil)
	client.
		On("DocumentTemplates", mock.Anything, sirius.CaseTypeDigitalLpa).
		Return(templateData, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, createDocumentDigitalLpaData{
			Lpa:                   digitalLpa,
			DocumentTemplates:     sortDocumentData(templateData),
			ComponentDocumentData: buildComponentDocumentData(templateData),
			Recipients: []sirius.Person{{
				ID:           -1,
				Firstname:    "Zackary",
				Surname:      "Lemmonds",
				PersonType:   "Donor",
				AddressLine1: "9 Mount Pleasant Drive",
				Town:         "East Harling",
				Postcode:     "NR16 2GB",
				Country:      "UK",
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
