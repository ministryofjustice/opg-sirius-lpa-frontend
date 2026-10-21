package server

import (
	"net/http"
	"testing"

	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockCasePageClient struct {
	mock.Mock
}

func (m *mockCasePageClient) GetPersonDocuments(ctx sirius.Context, personID int, caseIDs []string) (sirius.DocumentList, error) {
	args := m.Called(ctx, personID, caseIDs)
	return args.Get(0).(sirius.DocumentList), args.Error(1)
}

func (m *mockCasePageClient) CasesByDonor(ctx sirius.Context, id int) ([]sirius.Case, error) {
	args := m.Called(ctx, id)
	return args.Get(0).([]sirius.Case), args.Error(1)
}

func (m *mockCasePageClient) GetDraftCount(ctx sirius.Context, caseType string, caseId int) (sirius.DocumentDraftCount, error) {
	args := m.Called(ctx, caseType, caseId)
	return args.Get(0).(sirius.DocumentDraftCount), args.Error(1)
}

func (m *mockCasePageClient) Person(ctx sirius.Context, id int) (sirius.Person, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(sirius.Person), args.Error(1)
}

func (m *mockCasePageClient) PersonReferences(ctx sirius.Context, id int) ([]sirius.PersonReference, error) {
	args := m.Called(ctx, id)
	return args.Get(0).([]sirius.PersonReference), args.Error(1)
}

func (m *mockCasePageClient) TasksForCase(ctx sirius.Context, caseId int) ([]sirius.Task, error) {
	args := m.Called(ctx, caseId)
	return args.Get(0).([]sirius.Task), args.Error(1)
}

func (m *mockCasePageClient) GetUserPermissions(ctx sirius.Context) (sirius.Permissions, error) {
	args := m.Called(ctx)
	return args.Get(0).(sirius.Permissions), args.Error(1)
}

func (m *mockCasePageClient) FeatureToggles(ctx sirius.Context) (sirius.FeatureToggles, error) {
	args := m.Called(ctx)
	return args.Get(0).(sirius.FeatureToggles), args.Error(1)
}

func TestGetCasePage(t *testing.T) {
	cases := []sirius.Case{
		{ID: 1, UID: "7000-0000-0001", CaseType: "LPA"},
	}
	donor := sirius.Person{
		ID: 123,
	}

	client := &mockCasePageClient{}
	client.
		On("CasesByDonor", mock.Anything, 123).
		Return(cases, nil)
	client.
		On("PersonReferences", mock.Anything, 123).
		Return([]sirius.PersonReference{}, nil)
	client.
		On("Person", mock.Anything, 123).
		Return(sirius.Person{ID: 123}, nil)
	client.
		On("GetUserPermissions", mock.Anything).
		Return(sirius.Permissions{}, nil)
	client.
		On("GetPersonDocuments", mock.Anything, 123, []string{"1"}).
		Return(sirius.DocumentList{}, nil)
	client.
		On("TasksForCase", mock.Anything, 1).
		Return([]sirius.Task{}, nil)
	client.
		On("GetDraftCount", mock.Anything, "lpa", 1).
		Return(sirius.DocumentDraftCount{}, nil)
	client.
		On("FeatureToggles", mock.Anything).
		Return(sirius.FeatureToggles{"poasNewCasePage": true}, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything,
			casePageData{
				SelectedCases:                  cases,
				Person:                         donor,
				Donor:                          donor,
				DonorID:                        123,
				CaseUids:                       "&uid[]=7000-0000-0001",
				ActionPanelButtons:             GetActionPanelButtons(cases, 123, "&uid[]=7000-0000-0001", false, false, false, []int{}, sirius.Permissions{}),
				HasV1PersonsGetPermission:      false,
				HasV1PersonsCasesGetPermission: false,
				HeaderButtons: SiriusHeaderButtons{
					Calendar:  true,
					Documents: true,
				},
			},
		).
		Return(nil)

	server := newMockServer("/donor/{id}", CasePage(client, template.Func))

	r, _ := http.NewRequest(http.MethodGet, "/donor/123?uid[]=7000-0000-0001", nil)
	_, err := server.serve(r)

	assert.Nil(t, err)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestGetCasePageWhenFeatureTogglesErrors(t *testing.T) {
	client := &mockCasePageClient{}
	client.
		On("FeatureToggles", mock.Anything).
		Return(sirius.FeatureToggles{}, sirius.StatusError{Code: http.StatusForbidden})

	server := newMockServer("/donor/{id}", CasePage(client, nil))

	r, _ := http.NewRequest(http.MethodGet, "/donor/123?uid[]=7000-0000-0001", nil)
	_, err := server.serve(r)

	assert.Equal(t, sirius.StatusError{Code: http.StatusForbidden}, err)
}

func TestGetCasePageWhenFeatureTogglesNotEnabled(t *testing.T) {
	client := &mockCasePageClient{}
	client.
		On("FeatureToggles", mock.Anything).
		Return(sirius.FeatureToggles{"poasNewCasePage": false}, nil)

	server := newMockServer("/donor/{id}", CasePage(client, nil))

	r, _ := http.NewRequest(http.MethodGet, "/donor/123?uid[]=7000-0000-0001", nil)
	_, err := server.serve(r)

	assert.Equal(t, sirius.StatusError{Code: http.StatusForbidden}, err)
}

func TestGetCasePageWhenNoDonor(t *testing.T) {
	client := &mockCasePageClient{}
	client.
		On("GetUserPermissions", mock.Anything).
		Return(sirius.Permissions{}, nil)
	client.
		On("TasksForCase", mock.Anything, 1).
		Return([]sirius.Task{}, nil)
	client.
		On("GetDraftCount", mock.Anything, "lpa", 1).
		Return(sirius.DocumentDraftCount{}, nil)
	client.
		On("FeatureToggles", mock.Anything).
		Return(sirius.FeatureToggles{"poasNewCasePage": true}, nil)

	server := newMockServer("/donor/{id}", CasePage(client, nil))

	r, _ := http.NewRequest(http.MethodGet, "/donor/0?uid[]=7000-0000-0001", nil)
	_, err := server.serve(r)

	assert.Equal(t, sirius.StatusError{Code: http.StatusNotFound}, err)
}
