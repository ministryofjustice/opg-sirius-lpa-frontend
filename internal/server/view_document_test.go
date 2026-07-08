package server

import (
	"net/http"
	"strings"
	"testing"

	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockViewDocumentClient struct {
	mock.Mock
}

func (m *mockViewDocumentClient) DocumentByUUID(ctx sirius.Context, uuid string) (sirius.Document, error) {
	args := m.Called(ctx, uuid)
	return args.Get(0).(sirius.Document), args.Error(1)
}

func (m *mockViewDocumentClient) GetUserDetails(ctx sirius.Context) (sirius.User, error) {
	args := m.Called(ctx)
	return args.Get(0).(sirius.User), args.Error(1)
}

func (m *mockViewDocumentClient) Person(ctx sirius.Context, id int) (sirius.Person, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(sirius.Person), args.Error(1)
}

func (m *mockViewDocumentClient) GetUserPermissions(ctx sirius.Context) (sirius.Permissions, error) {
	args := m.Called(ctx)
	return args.Get(0).(sirius.Permissions), args.Error(1)
}

func (m *mockViewDocumentClient) Case(ctx sirius.Context, id int) (sirius.Case, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(sirius.Case), args.Error(1)
}

func (m *mockViewDocumentClient) GetDraftCount(ctx sirius.Context, caseType string, caseId int) (sirius.DocumentDraftCount, error) {
	args := m.Called(ctx, caseType, caseId)
	return args.Get(0).(sirius.DocumentDraftCount), args.Error(1)
}

func (m *mockViewDocumentClient) PersonReferences(ctx sirius.Context, id int) ([]sirius.PersonReference, error) {
	args := m.Called(ctx, id)
	return args.Get(0).([]sirius.PersonReference), args.Error(1)
}

func TestGetViewDocument(t *testing.T) {
	user := sirius.User{ID: 66, DisplayName: "Me", Roles: []string{"System Admin"}}
	person := sirius.Person{ID: 33}
	for _, caseType := range []string{"lpa", "epa"} {
		t.Run(caseType, func(t *testing.T) {
			document := sirius.Document{
				ID:         1,
				UUID:       "dfef6714-b4fe-44c2-b26e-90dfe3663e95",
				SystemType: "LP-LETTER",
				Type:       sirius.TypeSave,
			}
			caseData := sirius.Case{
				ID:       34,
				UID:      "7000-1234-1234",
				CaseType: strings.ToUpper(caseType),
			}
			draftCount := sirius.DocumentDraftCount{DraftCount: 0}

			client := &mockViewDocumentClient{}
			client.
				On("DocumentByUUID", mock.Anything, document.UUID).
				Return(document, nil)
			client.
				On("GetUserDetails", mock.Anything).
				Return(user, nil)
			client.
				On("Person", mock.Anything, 33).
				Return(person, nil)
			client.
				On("Case", mock.Anything, 34).
				Return(caseData, nil)
			client.
				On("GetDraftCount", mock.Anything, strings.ToLower(caseType), 34).
				Return(draftCount, nil)
			client.
				On("GetUserPermissions", mock.Anything).
				Return(sirius.Permissions{}, nil)
			client.
				On("PersonReferences", mock.Anything, 33).
				Return([]sirius.PersonReference{{ID: 987}}, nil)

			template := &mockTemplate{}
			templateData := viewDocumentData{
				Document:        document,
				IsSysAdminUser:  true,
				Pane:            1,
				DonorID:         33,
				SelectedCaseIds: "34",
				Person:          person,
				CaseUids:        "&uid[]=7000-1234-1234",
				SelectedCases:   []sirius.Case{caseData},
				HeaderButtons: SiriusHeaderButtons{
					BackToTimeline: true,
					Calendar:       true,
					CaseInfo:       true,
					PersonInfo:     true,
				},
			}

			template.
				On("Func", mock.Anything, mock.MatchedBy(func(data viewDocumentData) bool {
					return data.Document.UUID == templateData.Document.UUID &&
						data.IsSysAdminUser == templateData.IsSysAdminUser &&
						data.Pane == templateData.Pane &&
						data.DonorID == templateData.DonorID &&
						data.SelectedCaseIds == templateData.SelectedCaseIds &&
						data.CaseUids == templateData.CaseUids &&
						data.HeaderButtons == templateData.HeaderButtons
				})).
				Return(nil)

			server := newMockServer("/view-document/{uuid}/{donorId}", ViewDocument(client, template.Func))

			req, _ := http.NewRequest(http.MethodGet, "/view-document/dfef6714-b4fe-44c2-b26e-90dfe3663e95/33?case=34", nil)
			_, err := server.serve(req)

			assert.Nil(t, err)
			mock.AssertExpectationsForObjects(t, client, template)
		})
	}
}

func TestGetViewDocumentWhenCaseErrors(t *testing.T) {
	client := &mockViewDocumentClient{}

	client.
		On("Person", mock.Anything, 33).
		Return(sirius.Person{ID: 33}, nil)
	client.
		On("DocumentByUUID", mock.Anything, "dfef6714-b4fe-44c2-b26e-90dfe3663e95").
		Return(sirius.Document{}, errExample)

	server := newMockServer("/view-document/{uuid}/{donorId}", ViewDocument(client, nil))

	req, _ := http.NewRequest(http.MethodGet, "/view-document/dfef6714-b4fe-44c2-b26e-90dfe3663e95/33", nil)
	_, err := server.serve(req)

	assert.Equal(t, errExample, err)
}

func TestGetViewDocumentWhenPermissionsErrors(t *testing.T) {
	client := &mockViewDocumentClient{}

	document := sirius.Document{
		ID:         1,
		UUID:       "dfef6714-b4fe-44c2-b26e-90dfe3663e95",
		SystemType: "LP-LETTER",
		Type:       sirius.TypeSave,
	}

	client.
		On("Person", mock.Anything, 33).
		Return(sirius.Person{ID: 33}, nil)
	client.
		On("DocumentByUUID", mock.Anything, "dfef6714-b4fe-44c2-b26e-90dfe3663e95").
		Return(document, nil)
	client.
		On("GetUserDetails", mock.Anything).
		Return(sirius.User{}, nil)
	client.
		On("PersonReferences", mock.Anything, 33).
		Return([]sirius.PersonReference{{ID: 987}}, nil)
	client.
		On("Case", mock.Anything, 34).
		Return(sirius.Case{ID: 34, UID: "7000-1234-1234", CaseType: "lpa"}, nil)
	client.
		On("GetDraftCount", mock.Anything, "lpa", 34).
		Return(sirius.DocumentDraftCount{DraftCount: 0}, nil)
	client.
		On("GetUserPermissions", mock.Anything).
		Return(sirius.Permissions{}, errExample)

	server := newMockServer("/view-document/{uuid}/{donorId}", ViewDocument(client, nil))

	req, _ := http.NewRequest(http.MethodGet, "/view-document/dfef6714-b4fe-44c2-b26e-90dfe3663e95/33?case=34", nil)
	_, err := server.serve(req)

	assert.Equal(t, errExample, err)
}

func TestGetViewDocumentWhenGetUserDetailsErrors(t *testing.T) {
	client := &mockViewDocumentClient{}
	document := sirius.Document{
		ID:         1,
		UUID:       "dfef6714-b4fe-44c2-b26e-90dfe3663e95",
		SystemType: "LP-LETTER",
		Type:       sirius.TypeSave,
	}

	client.
		On("Person", mock.Anything, 33).
		Return(sirius.Person{ID: 33}, nil)
	client.
		On("DocumentByUUID", mock.Anything, document.UUID).
		Return(document, nil)
	client.
		On("GetUserDetails", mock.Anything).
		Return(sirius.User{}, errExample)

	server := newMockServer("/view-document/{uuid}/{donorId}", ViewDocument(client, nil))

	req, _ := http.NewRequest(http.MethodGet, "/view-document/dfef6714-b4fe-44c2-b26e-90dfe3663e95/33", nil)
	_, err := server.serve(req)

	assert.Equal(t, errExample, err)
}
