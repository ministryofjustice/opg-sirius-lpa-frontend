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

type mockUpdateInvestigationClient struct {
	mock.Mock
}

func (m *mockUpdateInvestigationClient) UpdateInvestigation(ctx sirius.Context, investigationID int, investigation sirius.Investigation) error {
	args := m.Called(ctx, investigationID, investigation)
	return args.Error(0)
}

func (m *mockUpdateInvestigationClient) Investigation(ctx sirius.Context, id int) (sirius.Investigation, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(sirius.Investigation), args.Error(1)
}

func TestGetUpdateInvestigation(t *testing.T) {
	investigation := sirius.Investigation{
		ID: 123,
	}

	client := &mockUpdateInvestigationClient{}
	client.
		On("Investigation", mock.Anything, 123).
		Return(investigation, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, updateInvestigationData{
			Investigation:        investigation,
			ApprovalOutcomeTypes: approvalOutcomeTypes,
		}).
		Return(nil)

	r, _ := http.NewRequest(http.MethodGet, "/?id=123", nil)
	w := httptest.NewRecorder()

	err := UpdateInvestigation(client, template.Func)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestGetUpdateInvestigationNoID(t *testing.T) {
	r, _ := http.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	err := UpdateInvestigation(nil, nil)(w, r)

	assert.NotNil(t, err)
}

func TestGetUpdateInvestigationWhenInvestigationErrors(t *testing.T) {
	client := &mockUpdateInvestigationClient{}
	client.
		On("Investigation", mock.Anything, 123).
		Return(sirius.Investigation{}, errExample)

	r, _ := http.NewRequest(http.MethodGet, "/?id=123", nil)
	w := httptest.NewRecorder()

	err := UpdateInvestigation(client, nil)(w, r)

	assert.Equal(t, errExample, err)
	mock.AssertExpectationsForObjects(t, client)
}

func TestGetUpdateInvestigationWhenTemplateErrors(t *testing.T) {
	client := &mockUpdateInvestigationClient{}
	client.
		On("Investigation", mock.Anything, 123).
		Return(sirius.Investigation{}, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, updateInvestigationData{
			ApprovalOutcomeTypes: approvalOutcomeTypes,
		}).
		Return(errExample)

	r, _ := http.NewRequest(http.MethodGet, "/?id=123", nil)
	w := httptest.NewRecorder()

	err := UpdateInvestigation(client, template.Func)(w, r)

	assert.Equal(t, errExample, err)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestPostUpdateInvestigation(t *testing.T) {
	investigation := sirius.Investigation{
		Title:                    "Test Investigation",
		Information:              "This is an investigation",
		Type:                     "Priority",
		DateReceived:             sirius.DateString("2022-04-05"),
		ApprovalDate:             sirius.DateString("2022-04-05"),
		RiskAssessmentDate:       sirius.DateString("2022-04-05"),
		ApprovalOutcome:          "Court Application",
		InvestigationClosureDate: sirius.DateString("2022-04-05"),
	}

	client := &mockUpdateInvestigationClient{}
	client.
		On("Investigation", mock.Anything, 123).
		Return(investigation, nil)
	client.
		On("UpdateInvestigation", mock.Anything, 123, investigation).
		Return(nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, updateInvestigationData{
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
		"riskAssessmentDate":       {"2022-04-05"},
		"approvalOutcome":          {"Court Application"},
		"investigationClosureDate": {"2022-04-05"},
	}

	r, _ := http.NewRequest(http.MethodPost, "/?id=123", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := UpdateInvestigation(client, template.Func)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestPostUpdateInvestigationWhenUpdateInvestigationValidationError(t *testing.T) {
	expectedError := sirius.ValidationError{
		Field: sirius.FieldErrors{"field": {"": "problem"}},
	}

	investigation := sirius.Investigation{Title: "Test Investigation"}

	client := &mockUpdateInvestigationClient{}
	client.
		On("Investigation", mock.Anything, 123).
		Return(investigation, nil)
	client.
		On("UpdateInvestigation", mock.Anything, 123, investigation).
		Return(expectedError)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, updateInvestigationData{
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

	err := UpdateInvestigation(client, template.Func)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestPostUpdateInvestigationWhenUpdateInvestigationOtherError(t *testing.T) {
	investigation := sirius.Investigation{Title: "Test Investigation"}

	client := &mockUpdateInvestigationClient{}
	client.
		On("UpdateInvestigation", mock.Anything, 123, investigation).
		Return(errExample)

	form := url.Values{
		"title": {"Test Investigation"},
	}

	r, _ := http.NewRequest(http.MethodPost, "/?id=123", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := UpdateInvestigation(client, nil)(w, r)

	assert.Equal(t, errExample, err)
	mock.AssertExpectationsForObjects(t, client)
}
