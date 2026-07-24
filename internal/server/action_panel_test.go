package server

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockActionPanelClient struct {
	mock.Mock
}

func (m *mockActionPanelClient) GetPersonDocuments(ctx sirius.Context, personID int, caseIDs []string) (sirius.DocumentList, error) {
	args := m.Called(ctx, personID, caseIDs)
	return args.Get(0).(sirius.DocumentList), args.Error(1)
}

func (m *mockActionPanelClient) CasesByDonor(ctx sirius.Context, id int) ([]sirius.Case, error) {
	args := m.Called(ctx, id)
	return args.Get(0).([]sirius.Case), args.Error(1)
}

func (m *mockActionPanelClient) GetDraftCount(ctx sirius.Context, caseType string, caseId int) (sirius.DocumentDraftCount, error) {
	args := m.Called(ctx, caseType, caseId)
	return args.Get(0).(sirius.DocumentDraftCount), args.Error(1)
}

func (m *mockActionPanelClient) Person(ctx sirius.Context, id int) (sirius.Person, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(sirius.Person), args.Error(1)
}

func (m *mockActionPanelClient) PersonReferences(ctx sirius.Context, id int) ([]sirius.PersonReference, error) {
	args := m.Called(ctx, id)
	return args.Get(0).([]sirius.PersonReference), args.Error(1)
}

func (m *mockActionPanelClient) TasksForCase(ctx sirius.Context, caseId int) ([]sirius.Task, error) {
	args := m.Called(ctx, caseId)
	return args.Get(0).([]sirius.Task), args.Error(1)
}

func (m *mockActionPanelClient) GetUserPermissions(ctx sirius.Context) (sirius.Permissions, error) {
	args := m.Called(ctx)
	return args.Get(0).(sirius.Permissions), args.Error(1)
}

var actionPanelPermissions = sirius.Permissions{
	"reporting":                sirius.PermissionType{Permissions: []string{"GET"}},
	"v1-cases-tasks-post":      sirius.PermissionType{Permissions: []string{"POST"}},
	"v1-donors":                sirius.PermissionType{Permissions: []string{"POST", "PUT"}},
	"v1-donors-epas":           sirius.PermissionType{Permissions: []string{"POST"}},
	"v1-donors-lpas":           sirius.PermissionType{Permissions: []string{"POST"}},
	"v1-lpas":                  sirius.PermissionType{Permissions: []string{"PUT"}},
	"v1-lpas-documents-draft":  sirius.PermissionType{Permissions: []string{"POST"}},
	"v1-lpas-investigations":   sirius.PermissionType{Permissions: []string{"POST"}},
	"v1-notes":                 sirius.PermissionType{Permissions: []string{"POST"}},
	"v1-payments":              sirius.PermissionType{Permissions: []string{"GET"}},
	"v1-person-links":          sirius.PermissionType{Permissions: []string{"POST", "PATCH"}},
	"v1-person-references":     sirius.PermissionType{Permissions: []string{"DELETE"}},
	"v1-persons":               sirius.PermissionType{Permissions: []string{"GET"}},
	"v1-persons-cases":         sirius.PermissionType{Permissions: []string{"GET"}},
	"v1-persons-references":    sirius.PermissionType{Permissions: []string{"POST"}},
	"v1-poa-tasks":             sirius.PermissionType{Permissions: []string{"PUT"}},
	"v1-users-updateusercases": sirius.PermissionType{Permissions: []string{"PUT"}},
	"v1-warnings":              sirius.PermissionType{Permissions: []string{"POST"}},
}

func TestGetActionPanel(t *testing.T) {
	cases := []sirius.Case{
		{ID: 1, UID: "7000-0000-0001", CaseType: "LPA"},
		{ID: 2, UID: "7000-0000-0002", CaseType: "LPA"},
	}

	client := &mockActionPanelClient{}
	client.
		On("CasesByDonor", mock.Anything, 123).
		Return(cases, nil)
	client.
		On("PersonReferences", mock.Anything, 123).
		Return([]sirius.PersonReference{{ID: 987}}, nil)
	client.
		On("Person", mock.Anything, 123).
		Return(sirius.Person{ID: 123}, nil)
	client.
		On("GetUserPermissions", mock.Anything).
		Return(actionPanelPermissions, nil)
	client.
		On("GetPersonDocuments", mock.Anything, 123, []string(nil)).
		Return(sirius.DocumentList{}, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, ActionPanelData{
			ActionPanelButtons: []ActionPanelButton{
				{
					Label:    "Create warning",
					URL:      "/create-warning?id=123&entity=person",
					IconName: "aw-create-warning",
					Disabled: false,
				},
				{
					Label:    "Create event",
					URL:      "/create-event?id=123&entity=person",
					IconName: "aw-new-event",
					Disabled: false,
				},
				{
					Label:    "Add complaint",
					URL:      "",
					IconName: "aw-log-complaint",
					Disabled: true,
				},
				{
					Label:    "Create document",
					URL:      "",
					IconName: "aw-new-template",
					Disabled: true,
				},
				{
					Label:    "Retrieve draft",
					URL:      "",
					IconName: "aw-new-template",
					Disabled: true,
				},
				{
					Label:    "Change status",
					URL:      "",
					IconName: "aw-change-status",
					Disabled: true,
				},
				{
					Label:    "Fees",
					URL:      "",
					IconName: "aw-fees",
					Disabled: true,
				},
				{
					Label:    "New task",
					URL:      "",
					IconName: "aw-new-task",
					Disabled: true,
				},
				{
					Label:    "Assign task",
					URL:      "",
					IconName: "aw-assign-task",
					Disabled: true,
				},
				{
					Label:    "Create donor",
					URL:      "/create-donor?id=123&entity=person",
					IconName: "aw-create-person",
					Disabled: false,
				},
				{
					Label:    "Edit donor",
					URL:      "/edit-donor?id=123&entity=person",
					IconName: "aw-edit-person",
					Disabled: false,
				},
				{
					Label:    "Edit dates",
					URL:      "",
					IconName: "calendar-open",
					Disabled: true,
				},
				{
					Label:    "MI reporting",
					URL:      "/mi-reporting?donorId=123",
					IconName: "aw-mi",
					Disabled: false,
				},
				{
					Label:    "Allocate Case",
					URL:      "/allocate-cases?id=1&id=2&entity=lpa",
					IconName: "aw-allocate-case",
					Disabled: false,
				},
				{
					Label:    "Link record",
					URL:      "/link-person?id=123",
					IconName: "aw-link",
					Disabled: false,
				},
				{
					Label:    "Unlink record",
					URL:      "/unlink-person?id=123",
					IconName: "aw-unlink",
					Disabled: true,
				},
				{
					Label:    "Delete relationship",
					URL:      "/delete-relationship?id=123",
					IconName: "icon-minus",
					Disabled: false,
				},
				{
					Label:    "Create relationship",
					URL:      "/create-relationship?id=123&entity=person",
					IconName: "aw-relationship",
					Disabled: false,
				},
				{
					Label:    "Create epa case",
					URL:      "/create-epa?id=123",
					IconName: "aw-create-case",
					Disabled: false,
				},
				{
					Label:    "Create lpa case",
					URL:      "/create-lpa?id=123",
					IconName: "aw-create-case",
					Disabled: false,
				},
				{
					Label:    "Edit case",
					URL:      "",
					IconName: "aw-edit-case",
					Disabled: true,
				},
				{
					Label:    "Add investigation",
					URL:      "",
					IconName: "icon-investigation",
					Disabled: true,
				},
			},
		}).
		Return(nil)

	r, _ := http.NewRequest(http.MethodGet, "/?donorId=123&entity=lpa", nil)
	w := httptest.NewRecorder()

	err := ActionPanel(client, template.Func)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
	client.AssertNotCalled(t, "TasksForCase")
}

func TestGetActionPanelWithUIDFilter(t *testing.T) {
	cases := []sirius.Case{
		{ID: 1, UID: "7000-0000-0001", CaseType: "LPA"},
		{ID: 2, UID: "7000-0000-0002", CaseType: "LPA"},
	}

	client := &mockActionPanelClient{}
	client.
		On("CasesByDonor", mock.Anything, 123).
		Return(cases, nil)
	client.
		On("GetDraftCount", mock.Anything, "lpa", 1).
		Return(sirius.DocumentDraftCount{DraftCount: 1}, nil)
	client.
		On("TasksForCase", mock.Anything, 1).
		Return([]sirius.Task{{ID: 990}}, nil)
	client.
		On("PersonReferences", mock.Anything, 123).
		Return([]sirius.PersonReference{{ID: 987}}, nil)
	client.
		On("Person", mock.Anything, 123).
		Return(sirius.Person{ID: 123, Children: []sirius.Person{{ID: 456}}}, nil)
	client.
		On("GetUserPermissions", mock.Anything).
		Return(actionPanelPermissions, nil)
	client.
		On("GetPersonDocuments", mock.Anything, 123, []string{"1"}).
		Return(sirius.DocumentList{}, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, ActionPanelData{
			ActionPanelButtons: []ActionPanelButton{
				{
					Label:    "Create warning",
					URL:      "/create-warning?id=123&entity=lpa&uid[]=7000-0000-0001",
					IconName: "aw-create-warning",
					Disabled: false,
				},
				{
					Label:    "Create event",
					URL:      "/create-event?id=123&entity=person&uid[]=7000-0000-0001",
					IconName: "aw-new-event",
					Disabled: false,
				},
				{
					Label:    "Add complaint",
					URL:      "/add-complaint?id=1&case=lpa",
					IconName: "aw-log-complaint",
					Disabled: false,
				},
				{
					Label:    "Create document",
					URL:      "/create-document?id=1&case=lpa",
					IconName: "aw-new-template",
					Disabled: false,
				},
				{
					Label:    "Retrieve draft",
					URL:      "/edit-document?id=1&case=lpa",
					IconName: "aw-new-template",
					Disabled: false,
				},
				{
					Label:    "Change status",
					URL:      "/change-status?id=1&case=lpa&donorId=123&uid[]=7000-0000-0001",
					IconName: "aw-change-status",
					Disabled: false,
				},
				{
					Label:    "Fees",
					URL:      "/payments/1",
					IconName: "aw-fees",
					Disabled: false,
				},
				{
					Label:    "New task",
					URL:      "/create-task?id=1&entity=lpa&uid[]=7000-0000-0001",
					IconName: "aw-new-task",
					Disabled: false,
				},
				{
					Label:    "Assign task",
					URL:      "/assign-task?id=990&donorId=123&uid[]=7000-0000-0001",
					IconName: "aw-assign-task",
					Disabled: false,
				},
				{
					Label:    "Create donor",
					URL:      "/create-donor?id=123&entity=person&uid[]=7000-0000-0001",
					IconName: "aw-create-person",
					Disabled: false,
				},
				{
					Label:    "Edit donor",
					URL:      "/edit-donor?id=123&entity=person&uid[]=7000-0000-0001",
					IconName: "aw-edit-person",
					Disabled: false,
				},
				{
					Label:    "Edit dates",
					URL:      "/edit-dates?id=1&case=lpa",
					IconName: "calendar-open",
					Disabled: false,
				},
				{
					Label:    "MI reporting",
					URL:      "/mi-reporting?donorId=123&uid[]=7000-0000-0001",
					IconName: "aw-mi",
					Disabled: false,
				},
				{
					Label:    "Allocate Case",
					URL:      "/allocate-cases?id=1&entity=lpa&uid[]=7000-0000-0001",
					IconName: "aw-allocate-case",
					Disabled: false,
				},
				{
					Label:    "Link record",
					URL:      "/link-person?id=123&uid[]=7000-0000-0001",
					IconName: "aw-link",
					Disabled: false,
				},
				{
					Label:    "Unlink record",
					URL:      "/unlink-person?id=123&uid[]=7000-0000-0001",
					IconName: "aw-unlink",
					Disabled: false,
				},
				{
					Label:    "Delete relationship",
					URL:      "/delete-relationship?id=123&uid[]=7000-0000-0001",
					IconName: "icon-minus",
					Disabled: false,
				},
				{
					Label:    "Create relationship",
					URL:      "/create-relationship?id=123&entity=person&uid[]=7000-0000-0001",
					IconName: "aw-relationship",
					Disabled: false,
				},
				{
					Label:    "Create epa case",
					URL:      "/create-epa?id=123",
					IconName: "aw-create-case",
					Disabled: true,
				},
				{
					Label:    "Create lpa case",
					URL:      "/create-lpa?id=123",
					IconName: "aw-create-case",
					Disabled: true,
				},
				{
					Label:    "Edit case",
					URL:      "/create-lpa?id=123&caseId=1",
					IconName: "aw-edit-case",
					Disabled: false,
				},
				{
					Label:    "Add investigation",
					URL:      "/create-investigation?id=1&case=lpa&uid[]=7000-0000-0001",
					IconName: "icon-investigation",
					Disabled: false,
				},
			},
		}).
		Return(nil)

	r, _ := http.NewRequest(http.MethodGet, "/?donorId=123&entity=lpa&uid[]=7000-0000-0001", nil)
	w := httptest.NewRecorder()

	err := ActionPanel(client, template.Func)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestGetActionPanelNoOutstandingTasks(t *testing.T) {
	cases := []sirius.Case{
		{ID: 1, UID: "7000-0000-0001", CaseType: "LPA"},
	}

	client := &mockActionPanelClient{}
	client.
		On("CasesByDonor", mock.Anything, 123).
		Return(cases, nil)
	client.
		On("GetDraftCount", mock.Anything, "lpa", 1).
		Return(sirius.DocumentDraftCount{}, nil)
	client.
		On("TasksForCase", mock.Anything, 1).
		Return([]sirius.Task{}, nil)
	client.
		On("PersonReferences", mock.Anything, 123).
		Return([]sirius.PersonReference{}, nil)
	client.
		On("Person", mock.Anything, 123).
		Return(sirius.Person{}, nil)
	client.
		On("GetUserPermissions", mock.Anything).
		Return(sirius.Permissions{}, nil)
	client.
		On("GetPersonDocuments", mock.Anything, 123, []string(nil)).
		Return(sirius.DocumentList{}, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, mock.MatchedBy(func(data ActionPanelData) bool {
			for _, btn := range data.ActionPanelButtons {
				if btn.Label == "Assign task" {
					return assert.Equal(t, "", btn.URL) && assert.True(t, btn.Disabled)
				}
			}
			return false
		})).
		Return(nil)

	r, _ := http.NewRequest(http.MethodGet, "/?donorId=123&entity=lpa", nil)
	w := httptest.NewRecorder()

	err := ActionPanel(client, template.Func)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestGetActionPanelNoDonorID(t *testing.T) {
	client := &mockActionPanelClient{}
	client.
		On("GetUserPermissions", mock.Anything).
		Return(actionPanelPermissions, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, ActionPanelData{
			ActionPanelButtons: []ActionPanelButton{
				{
					Label:    "Create warning",
					URL:      "/create-warning?id=0&entity=person",
					IconName: "aw-create-warning",
					Disabled: false,
				},
				{
					Label:    "Create event",
					URL:      "/create-event?id=0&entity=person",
					IconName: "aw-new-event",
					Disabled: false,
				},
				{
					Label:    "Add complaint",
					URL:      "",
					IconName: "aw-log-complaint",
					Disabled: true,
				},
				{
					Label:    "Create document",
					URL:      "",
					IconName: "aw-new-template",
					Disabled: true,
				},
				{
					Label:    "Retrieve draft",
					URL:      "",
					IconName: "aw-new-template",
					Disabled: true,
				},
				{
					Label:    "Change status",
					URL:      "",
					IconName: "aw-change-status",
					Disabled: true,
				},
				{
					Label:    "Fees",
					URL:      "",
					IconName: "aw-fees",
					Disabled: true,
				},
				{
					Label:    "New task",
					URL:      "",
					IconName: "aw-new-task",
					Disabled: true,
				},
				{
					Label:    "Assign task",
					URL:      "",
					IconName: "aw-assign-task",
					Disabled: true,
				},
				{
					Label:    "Create donor",
					URL:      "/create-donor?id=0&entity=person",
					IconName: "aw-create-person",
					Disabled: false,
				},
				{
					Label:    "Edit donor",
					URL:      "/edit-donor?id=0&entity=person",
					IconName: "aw-edit-person",
					Disabled: false,
				},
				{
					Label:    "Edit dates",
					URL:      "",
					IconName: "calendar-open",
					Disabled: true,
				},
				{
					Label:    "MI reporting",
					URL:      "/mi-reporting?donorId=0",
					IconName: "aw-mi",
					Disabled: false,
				},
				{
					Label:    "Allocate Case",
					URL:      "",
					IconName: "aw-allocate-case",
					Disabled: true,
				},
				{
					Label:    "Link record",
					URL:      "/link-person?id=0",
					IconName: "aw-link",
					Disabled: true,
				},
				{
					Label:    "Unlink record",
					URL:      "/unlink-person?id=0",
					IconName: "aw-unlink",
					Disabled: true,
				},
				{
					Label:    "Delete relationship",
					URL:      "/delete-relationship?id=0",
					IconName: "icon-minus",
					Disabled: true,
				},
				{
					Label:    "Create relationship",
					URL:      "/create-relationship?id=0&entity=person",
					IconName: "aw-relationship",
					Disabled: false,
				},
				{
					Label:    "Create epa case",
					URL:      "/create-epa?id=0",
					IconName: "aw-create-case",
					Disabled: false,
				},
				{
					Label:    "Create lpa case",
					URL:      "/create-lpa?id=0",
					IconName: "aw-create-case",
					Disabled: false,
				},
				{
					Label:    "Edit case",
					URL:      "",
					IconName: "aw-edit-case",
					Disabled: true,
				},
				{
					Label:    "Add investigation",
					URL:      "",
					IconName: "icon-investigation",
					Disabled: true,
				},
			},
		}).
		Return(nil)

	r, _ := http.NewRequest(http.MethodGet, "/?entity=lpa", nil)
	w := httptest.NewRecorder()

	err := ActionPanel(client, template.Func)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
	client.AssertNotCalled(t, "CasesByDonor")
	client.AssertNotCalled(t, "TasksForCase")
	client.AssertNotCalled(t, "Person")
}

func TestGetActionPanelEditEpaOnlyEnabledWhenSingleEpaCaseSelected(t *testing.T) {
	cases := []sirius.Case{
		{ID: 1, UID: "7000-0000-0001", CaseType: "LPA"},
		{ID: 2, UID: "7000-0000-0002", CaseType: "LPA"},
		{ID: 3, UID: "7000-0000-0003", CaseType: "EPA"},
	}

	client := &mockActionPanelClient{}
	client.
		On("CasesByDonor", mock.Anything, 123).
		Return(cases, nil)
	client.
		On("GetDraftCount", mock.Anything, "epa", 3).
		Return(sirius.DocumentDraftCount{DraftCount: 1}, nil)
	client.
		On("Person", mock.Anything, 123).
		Return(sirius.Person{}, nil)
	client.
		On("TasksForCase", mock.Anything, 3).
		Return([]sirius.Task{}, nil)
	client.
		On("PersonReferences", mock.Anything, 123).
		Return([]sirius.PersonReference{{ID: 987}}, nil)
	client.
		On("GetPersonDocuments", mock.Anything, 123, []string{"3"}).
		Return(sirius.DocumentList{}, nil)
	client.
		On("GetUserPermissions", mock.Anything).
		Return(actionPanelPermissions, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, mock.MatchedBy(func(data ActionPanelData) bool {
			for _, btn := range data.ActionPanelButtons {
				if btn.Label == "Edit case" {
					return btn.URL == "/create-epa?id=123&caseId=3" && btn.Disabled == false
				}
			}
			return false
		})).
		Return(nil)

	r, _ := http.NewRequest(http.MethodGet, "/?donorId=123&entity=epa&uid[]=7000-0000-0003", nil)
	w := httptest.NewRecorder()

	err := ActionPanel(client, template.Func)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
	client.AssertNotCalled(t, "CasesByDonor")
	client.AssertNotCalled(t, "Person")
}

func TestGetActionPanelEditLpaOnlyEnabledWhenSingleLpaCaseSelected(t *testing.T) {
	cases := []sirius.Case{
		{ID: 1, UID: "7000-0000-0001", CaseType: "LPA"},
		{ID: 2, UID: "7000-0000-0002", CaseType: "LPA"},
		{ID: 3, UID: "7000-0000-0003", CaseType: "EPA"},
	}

	client := &mockActionPanelClient{}
	client.
		On("CasesByDonor", mock.Anything, 123).
		Return(cases, nil)
	client.
		On("GetDraftCount", mock.Anything, "lpa", 2).
		Return(sirius.DocumentDraftCount{DraftCount: 1}, nil)
	client.
		On("Person", mock.Anything, 123).
		Return(sirius.Person{}, nil)
	client.
		On("TasksForCase", mock.Anything, 2).
		Return([]sirius.Task{}, nil)
	client.
		On("PersonReferences", mock.Anything, 123).
		Return([]sirius.PersonReference{{ID: 987}}, nil)
	client.
		On("GetUserPermissions", mock.Anything).
		Return(actionPanelPermissions, nil)
	client.
		On("GetPersonDocuments", mock.Anything, 123, []string{"2"}).
		Return(sirius.DocumentList{}, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, mock.MatchedBy(func(data ActionPanelData) bool {
			for _, btn := range data.ActionPanelButtons {
				if btn.Label == "Edit case" {
					return btn.URL == "/create-lpa?id=123&caseId=2" && btn.Disabled == false
				}
			}
			return false
		})).
		Return(nil)

	r, _ := http.NewRequest(http.MethodGet, "/?donorId=123&entity=lpa&uid[]=7000-0000-0002", nil)
	w := httptest.NewRecorder()

	err := ActionPanel(client, template.Func)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
	client.AssertNotCalled(t, "CasesByDonor")
	client.AssertNotCalled(t, "Person")
}

func TestGetActionPanelWhenCasesByDonorErrors(t *testing.T) {
	expectedError := errors.New("cases by donor error")

	client := &mockActionPanelClient{}
	client.
		On("CasesByDonor", mock.Anything, 123).
		Return([]sirius.Case{}, expectedError)
	client.
		On("Person", mock.Anything, 123).
		Return(sirius.Person{}, nil)
	client.
		On("PersonReferences", mock.Anything, 123).
		Return([]sirius.PersonReference{{ID: 987}}, nil)
	client.
		On("GetUserPermissions", mock.Anything).
		Return(actionPanelPermissions, nil)

	r, _ := http.NewRequest(http.MethodGet, "/?donorId=123&entity=lpa", nil)
	w := httptest.NewRecorder()

	err := ActionPanel(client, nil)(w, r)

	assert.Equal(t, expectedError, err)
}

func TestGetActionPanelWhenPermissionsErrors(t *testing.T) {
	client := &mockActionPanelClient{}
	client.
		On("CasesByDonor", mock.Anything, 123).
		Return([]sirius.Case{{ID: 1, UID: "7000-0000-0001", CaseType: "LPA"}}, nil)
	client.
		On("GetDraftCount", mock.Anything, "lpa", 1).
		Return(sirius.DocumentDraftCount{DraftCount: 0}, nil)
	client.
		On("TasksForCase", mock.Anything, 1).
		Return([]sirius.Task{}, nil)
	client.
		On("PersonReferences", mock.Anything, 123).
		Return([]sirius.PersonReference{{ID: 987}}, nil)
	client.
		On("Person", mock.Anything, 123).
		Return(sirius.Person{}, nil)
	client.
		On("GetPersonDocuments", mock.Anything, 123, []string(nil)).
		Return(sirius.DocumentList{}, nil)
	client.
		On("GetUserPermissions", mock.Anything).
		Return(actionPanelPermissions, errExample)

	r, _ := http.NewRequest(http.MethodGet, "/?donorId=123&entity=lpa", nil)
	w := httptest.NewRecorder()

	err := ActionPanel(client, nil)(w, r)

	assert.Equal(t, errExample, err)
}

func TestGetActionPanelWhenGetDraftCountErrors(t *testing.T) {
	expectedError := errors.New("get draft count error")

	cases := []sirius.Case{
		{ID: 1, UID: "7000-0000-0001", CaseType: "LPA"},
		{ID: 2, UID: "7000-0000-0002", CaseType: "LPA"},
	}

	client := &mockActionPanelClient{}
	client.
		On("CasesByDonor", mock.Anything, 123).
		Return(cases, nil)
	client.
		On("GetDraftCount", mock.Anything, "lpa", 1).
		Return(sirius.DocumentDraftCount{}, expectedError)
	client.
		On("TasksForCase", mock.Anything, 1).
		Return([]sirius.Task{}, nil)
	client.
		On("Person", mock.Anything, 123).
		Return(sirius.Person{}, nil)
	client.
		On("PersonReferences", mock.Anything, 123).
		Return([]sirius.PersonReference{{ID: 987}}, nil)
	client.
		On("GetPersonDocuments", mock.Anything, 123, []string{"1"}).
		Return(sirius.DocumentList{}, nil)
	client.
		On("GetUserPermissions", mock.Anything).
		Return(actionPanelPermissions, nil)

	r, _ := http.NewRequest(http.MethodGet, "/?donorId=123&entity=lpa&uid[]=7000-0000-0001", nil)
	w := httptest.NewRecorder()

	err := ActionPanel(client, nil)(w, r)

	assert.Equal(t, expectedError, err)
	client.AssertNotCalled(t, "TasksForCase")
}

func TestGetActionPanelWhenTasksForCaseErrors(t *testing.T) {
	expectedError := errors.New("tasks for case error")

	cases := []sirius.Case{
		{ID: 1, UID: "7000-0000-0001", CaseType: "LPA"},
		{ID: 2, UID: "7000-0000-0002", CaseType: "LPA"},
	}

	client := &mockActionPanelClient{}
	client.
		On("CasesByDonor", mock.Anything, 123).
		Return(cases, nil)
	client.
		On("GetDraftCount", mock.Anything, "lpa", 1).
		Return(sirius.DocumentDraftCount{DraftCount: 1}, nil)
	client.
		On("TasksForCase", mock.Anything, 1).
		Return([]sirius.Task{}, expectedError)
	client.
		On("GetPersonDocuments", mock.Anything, 123, []string{"1"}).
		Return(sirius.DocumentList{}, nil)
	client.
		On("PersonReferences", mock.Anything, 123).
		Return([]sirius.PersonReference{{ID: 987}}, nil)
	client.
		On("Person", mock.Anything, 123).
		Return(sirius.Person{}, nil)
	client.
		On("GetUserPermissions", mock.Anything).
		Return(sirius.Permissions{}, nil)

	r, _ := http.NewRequest(http.MethodGet, "/?donorId=123&entity=lpa&uid[]=7000-0000-0001", nil)
	w := httptest.NewRecorder()

	err := ActionPanel(client, nil)(w, r)

	assert.Equal(t, expectedError, err)
}

func TestGetActionPanelWhenPersonReferencesErrors(t *testing.T) {
	expectedError := errors.New("get person references error")

	cases := []sirius.Case{
		{ID: 1, UID: "7000-0000-0001", CaseType: "LPA"},
		{ID: 2, UID: "7000-0000-0002", CaseType: "LPA"},
	}

	client := &mockActionPanelClient{}
	client.
		On("CasesByDonor", mock.Anything, 123).
		Return(cases, nil)
	client.
		On("GetDraftCount", mock.Anything, "lpa", 1).
		Return(sirius.DocumentDraftCount{}, nil)
	client.
		On("Person", mock.Anything, 123).
		Return(sirius.Person{}, nil)
	client.
		On("TasksForCase", mock.Anything, 1).
		Return([]sirius.Task{}, nil)
	client.
		On("GetPersonDocuments", mock.Anything, 123, []string{"1"}).
		Return(sirius.DocumentList{}, nil)
	client.
		On("PersonReferences", mock.Anything, 123).
		Return([]sirius.PersonReference{}, expectedError)
	client.
		On("GetUserPermissions", mock.Anything).
		Return(actionPanelPermissions, nil)

	r, _ := http.NewRequest(http.MethodGet, "/?donorId=123&entity=lpa&uid[]=7000-0000-0001", nil)
	w := httptest.NewRecorder()

	err := ActionPanel(client, nil)(w, r)

	assert.Equal(t, expectedError, err)
}

func TestGetActionPanelWhenPersonErrors(t *testing.T) {
	expectedError := errors.New("get person error")

	cases := []sirius.Case{
		{ID: 1, UID: "7000-0000-0001", CaseType: "LPA"},
		{ID: 2, UID: "7000-0000-0002", CaseType: "LPA"},
	}

	client := &mockActionPanelClient{}
	client.
		On("CasesByDonor", mock.Anything, 123).
		Return(cases, nil)
	client.
		On("GetDraftCount", mock.Anything, "lpa", 1).
		Return(sirius.DocumentDraftCount{}, nil)
	client.
		On("TasksForCase", mock.Anything, 1).
		Return([]sirius.Task{}, nil)
	client.
		On("GetPersonDocuments", mock.Anything, 123, []string{"1"}).
		Return(sirius.DocumentList{}, nil)
	client.
		On("PersonReferences", mock.Anything, 123).
		Return([]sirius.PersonReference{{ID: 987}}, nil)
	client.
		On("Person", mock.Anything, 123).
		Return(sirius.Person{}, expectedError)
	client.
		On("GetUserPermissions", mock.Anything).
		Return(actionPanelPermissions, nil)

	r, _ := http.NewRequest(http.MethodGet, "/?donorId=123&entity=lpa&uid[]=7000-0000-0001", nil)
	w := httptest.NewRecorder()

	err := ActionPanel(client, nil)(w, r)

	assert.Equal(t, expectedError, err)
}
