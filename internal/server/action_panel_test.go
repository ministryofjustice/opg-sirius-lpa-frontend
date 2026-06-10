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

func (m *mockActionPanelClient) CasesByDonor(ctx sirius.Context, id int) ([]sirius.Case, error) {
	args := m.Called(ctx, id)
	return args.Get(0).([]sirius.Case), args.Error(1)
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

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, ActionPanelData{
			DonorID:       123,
			SelectedCases: cases,
			CaseType:      "lpa",
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

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, ActionPanelData{
			DonorID:       123,
			SelectedCases: []sirius.Case{cases[0]},
			CaseUids:      "&uid[]=7000-0000-0001",
			CaseType:      "lpa",
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

func TestGetActionPanelNoDonorID(t *testing.T) {
	client := &mockActionPanelClient{}

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, ActionPanelData{
			CaseType: "lpa",
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
}

func TestGetActionPanelInvalidEntityType(t *testing.T) {
	r, _ := http.NewRequest(http.MethodGet, "/?entity=invalid", nil)
	w := httptest.NewRecorder()

	err := ActionPanel(nil, nil)(w, r)

	assert.NotNil(t, err)
}

func TestGetActionPanelWhenCasesByDonorErrors(t *testing.T) {
	expectedError := errors.New("cases by donor error")

	client := &mockActionPanelClient{}
	client.
		On("CasesByDonor", mock.Anything, 123).
		Return([]sirius.Case{}, expectedError)

	r, _ := http.NewRequest(http.MethodGet, "/?donorId=123&entity=lpa", nil)
	w := httptest.NewRecorder()

	err := ActionPanel(client, nil)(w, r)

	assert.Equal(t, expectedError, err)
}
