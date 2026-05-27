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
