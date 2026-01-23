package server

import (
	"net/http"
	"testing"

	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockDeleteDocumentClient struct {
	mock.Mock
}

func (m *mockDeleteDocumentClient) GetUserDetails(ctx sirius.Context) (sirius.User, error) {
	args := m.Called(ctx)
	return args.Get(0).(sirius.User), args.Error(1)
}

func (m *mockDeleteDocumentClient) DeleteDocument(ctx sirius.Context, uuid string) error {
	return m.Called(ctx, uuid).Error(0)
}

func (m *mockDeleteDocumentClient) DocumentByUUID(ctx sirius.Context, uuid string) (sirius.Document, error) {
	args := m.Called(ctx, uuid)
	return args.Get(0).(sirius.Document), args.Error(1)
}

func TestGetDeleteDocument(t *testing.T) {
	user := sirius.User{ID: 66, DisplayName: "Me", Roles: []string{"System Admin"}}

	for _, caseType := range []string{"lpa", "epa"} {
		t.Run(caseType, func(t *testing.T) {
			donor := sirius.Person{ID: 1}
			caseItem := sirius.Case{CaseType: caseType, UID: "7000", Donor: &donor}
			document := sirius.Document{
				ID:         1,
				UUID:       "dfef6714-b4fe-44c2-b26e-90dfe3663e95",
				SystemType: "LP-LETTER",
				Type:       sirius.TypeSave,
				CaseItems:  []sirius.Case{caseItem},
			}

			client := &mockDeleteDocumentClient{}
			client.
				On("DocumentByUUID", mock.Anything, document.UUID).
				Return(document, nil)
			client.
				On("GetUserDetails", mock.Anything).
				Return(user, nil)

			template := &mockTemplate{}
			templateData := deleteDocumentData{
				Document:       document,
				IsSysAdminUser: true,
				DonorId:        1,
			}

			template.
				On("Func", mock.Anything, templateData).
				Return(nil)

			server := newMockServer("/view-document/{uuid}", DeleteDocument(client, template.Func))

			req, _ := http.NewRequest(http.MethodGet, "/view-document/dfef6714-b4fe-44c2-b26e-90dfe3663e95", nil)
			_, err := server.serve(req)

			assert.Nil(t, err)
			mock.AssertExpectationsForObjects(t, client, template)
		})
	}
}

func TestGetDeleteDocumentWhenCaseErrors(t *testing.T) {
	client := &mockDeleteDocumentClient{}

	client.
		On("DocumentByUUID", mock.Anything, "dfef6714-b4fe-44c2-b26e-90dfe3663e95").
		Return(sirius.Document{}, errExample)

	server := newMockServer("/delete-document/{uuid}", DeleteDocument(client, nil))

	req, _ := http.NewRequest(http.MethodGet, "/delete-document/dfef6714-b4fe-44c2-b26e-90dfe3663e95", nil)
	_, err := server.serve(req)

	assert.Equal(t, errExample, err)
}

func TestPostDeletingADocument(t *testing.T) {
	donor := sirius.Person{ID: 1}
	caseItem := sirius.Case{CaseType: "lpa", UID: "7000", Donor: &donor}
	document := sirius.Document{
		ID:                  1,
		UUID:                "dfef6714-b4fe-44c2-b26e-90dfe3663e95",
		SystemType:          "LP-LETTER",
		Type:                sirius.TypeSave,
		CaseItems:           []sirius.Case{caseItem},
		FriendlyDescription: "Docfriendly",
		CreatedDate:         "02/07/2025",
	}
	user := sirius.User{ID: 66, DisplayName: "Me", Roles: []string{"System Admin"}}

	client := &mockDeleteDocumentClient{}
	client.
		On("DocumentByUUID", mock.Anything, document.UUID).
		Return(document, nil)
	client.
		On("GetUserDetails", mock.Anything).
		Return(user, nil)
	client.
		On("DeleteDocument", mock.Anything, document.UUID).
		Return(nil)

	template := &mockTemplate{}

	server := newMockServer("/delete-document/{uuid}", DeleteDocument(client, template.Func))
	r, _ := http.NewRequest(http.MethodPost, "/delete-document/dfef6714-b4fe-44c2-b26e-90dfe3663e95", nil)
	r.Header.Add("Content-Type", formUrlEncoded)

	_, err := server.serve(r)

	assert.Equal(t, RedirectError("/donor/1/documents?success=true&documentFriendlyName=Docfriendly&documentCreatedTime=02/07/2025&uid[]=7000"), err)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestGetDeleteDocumentWhenGetUserDetailsErrors(t *testing.T) {
	client := &mockDeleteDocumentClient{}
	document := sirius.Document{
		ID:         1,
		UUID:       "dfef6714-b4fe-44c2-b26e-90dfe3663e95",
		SystemType: "LP-LETTER",
		Type:       sirius.TypeSave,
	}

	client.
		On("DocumentByUUID", mock.Anything, document.UUID).
		Return(document, nil)
	client.
		On("GetUserDetails", mock.Anything).
		Return(sirius.User{}, errExample)

	server := newMockServer("/view-document/{uuid}", DeleteDocument(client, nil))

	req, _ := http.NewRequest(http.MethodGet, "/view-document/dfef6714-b4fe-44c2-b26e-90dfe3663e95", nil)
	_, err := server.serve(req)

	assert.Equal(t, errExample, err)
}

func TestPostDeleteDocumentWhenDeleteErrors(t *testing.T) {
	donor := sirius.Person{ID: 1}
	caseItem := sirius.Case{CaseType: "lpa", UID: "7000", Donor: &donor}
	document := sirius.Document{
		ID:                  1,
		UUID:                "dfef6714-b4fe-44c2-b26e-90dfe3663e95",
		SystemType:          "LP-LETTER",
		Type:                sirius.TypeSave,
		CaseItems:           []sirius.Case{caseItem},
		FriendlyDescription: "Docfriendly",
		CreatedDate:         "02/07/2025",
	}
	user := sirius.User{ID: 66, DisplayName: "Me", Roles: []string{"System Admin"}}

	client := &mockDeleteDocumentClient{}
	client.
		On("DocumentByUUID", mock.Anything, document.UUID).
		Return(document, nil)
	client.
		On("GetUserDetails", mock.Anything).
		Return(user, nil)
	client.
		On("DeleteDocument", mock.Anything, document.UUID).
		Return(errExample)

	server := newMockServer("/delete-document/{uuid}", DeleteDocument(client, nil))
	r, _ := http.NewRequest(http.MethodPost, "/delete-document/dfef6714-b4fe-44c2-b26e-90dfe3663e95", nil)
	r.Header.Add("Content-Type", formUrlEncoded)

	_, err := server.serve(r)

	assert.Equal(t, errExample, err)
	mock.AssertExpectationsForObjects(t, client)
}
