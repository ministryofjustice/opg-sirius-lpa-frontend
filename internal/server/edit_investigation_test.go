package server

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockEditInvestigationClient struct {
	mock.Mock
}

func (m *mockEditInvestigationClient) EditInvestigation(ctx sirius.Context, investigationID int, investigation sirius.Investigation) error {
	args := m.Called(ctx, investigationID, investigation)
	return args.Error(0)
}

func (m *mockEditInvestigationClient) Investigation(ctx sirius.Context, id int) (sirius.Investigation, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(sirius.Investigation), args.Error(1)
}

func TestGetEditInvestigation(t *testing.T) {
	investigation := sirius.Investigation{
		ID: 123,
	}

	client := &mockEditInvestigationClient{}
	client.
		On("Investigation", mock.Anything, 123).
		Return(investigation, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, editInvestigationData{
			Investigation:        investigation,
			ApprovalOutcomeTypes: approvalOutcomeTypes,
		}).
		Return(nil)

	r, _ := http.NewRequest(http.MethodGet, "/?id=123", nil)
	w := httptest.NewRecorder()

	err := EditInvestigation(client, template.Func)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestGetEditInvestigationNoID(t *testing.T) {
	r, _ := http.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	err := EditInvestigation(nil, nil)(w, r)

	assert.NotNil(t, err)
}

func TestGetEditInvestigationWhenInvestigationErrors(t *testing.T) {
	client := &mockEditInvestigationClient{}
	client.
		On("Investigation", mock.Anything, 123).
		Return(sirius.Investigation{}, expectedError)

	r, _ := http.NewRequest(http.MethodGet, "/?id=123", nil)
	w := httptest.NewRecorder()

	err := EditInvestigation(client, nil)(w, r)

	assert.Equal(t, expectedError, err)
	mock.AssertExpectationsForObjects(t, client)
}

func TestGetEditInvestigationWhenTemplateErrors(t *testing.T) {
	client := &mockEditInvestigationClient{}
	client.
		On("Investigation", mock.Anything, 123).
		Return(sirius.Investigation{}, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, editInvestigationData{
			ApprovalOutcomeTypes: approvalOutcomeTypes,
		}).
		Return(expectedError)

	r, _ := http.NewRequest(http.MethodGet, "/?id=123", nil)
	w := httptest.NewRecorder()

	err := EditInvestigation(client, template.Func)(w, r)

	assert.Equal(t, expectedError, err)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestPostEditInvestigation(t *testing.T) {
	investigation := sirius.Investigation{
		Title:                    "Test Investigation",
		Information:              "This is an investigation",
		Type:                     "Priority",
		DateReceived:             sirius.DateString("2022-04-05"),
		ApprovalDate:             sirius.DateString("2022-04-05"),
		ApprovalOutcome:          "Court Application",
		InvestigationClosureDate: sirius.DateString("2022-04-05"),
	}

	client := &mockEditInvestigationClient{}
	client.
		On("Investigation", mock.Anything, 123).
		Return(investigation, nil)
	client.
		On("EditInvestigation", mock.Anything, 123, investigation).
		Return(nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, editInvestigationData{
			Success:              true,
			ApprovalOutcomeTypes: approvalOutcomeTypes,
			Investigation:        investigation,
		}).
		Return(nil)

	form := url.Values{
		"title":                    {"Test Investigation"},
		"information":              {"This is an investigation"},
		"type":                     {"Priority"},
		"dateReceived":             {"2022-04-05"},
		"approvalDate":             {"2022-04-05"},
		"approvalOutcome":          {"Court Application"},
		"investigationClosureDate": {"2022-04-05"},
	}

	r, _ := http.NewRequest(http.MethodPost, "/?id=123", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := EditInvestigation(client, template.Func)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestPostEditInvestigationWhenEditInvestigationValidationError(t *testing.T) {
	expectedError := sirius.ValidationError{
		Field: sirius.FieldErrors{"field": {"": "problem"}},
	}

	investigation := sirius.Investigation{Title: "Test Investigation"}

	client := &mockEditInvestigationClient{}
	client.
		On("Investigation", mock.Anything, 123).
		Return(investigation, nil)
	client.
		On("EditInvestigation", mock.Anything, 123, investigation).
		Return(expectedError)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, editInvestigationData{
			Success:              false,
			Error:                expectedError,
			Investigation:        investigation,
			ApprovalOutcomeTypes: approvalOutcomeTypes,
		}).
		Return(nil)

	form := url.Values{
		"title": {"Test Investigation"},
	}

	r, _ := http.NewRequest(http.MethodPost, "/?id=123", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := EditInvestigation(client, template.Func)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestPostEditInvestigationWhenEditInvestigationOtherError(t *testing.T) {
	investigation := sirius.Investigation{Title: "Test Investigation"}

	client := &mockEditInvestigationClient{}
	client.
		On("EditInvestigation", mock.Anything, 123, investigation).
		Return(expectedError)

	form := url.Values{
		"title": {"Test Investigation"},
	}

	r, _ := http.NewRequest(http.MethodPost, "/?id=123", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := EditInvestigation(client, nil)(w, r)

	assert.Equal(t, expectedError, err)
	mock.AssertExpectationsForObjects(t, client)
}
