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

type mockCreateInvestigationClient struct {
	mock.Mock
}

func (m *mockCreateInvestigationClient) CreateInvestigation(ctx sirius.Context, caseID int, caseType sirius.CaseType, investigation sirius.Investigation) error {
	args := m.Called(ctx, caseID, caseType, investigation)
	return args.Error(0)
}

func (m *mockCreateInvestigationClient) Case(ctx sirius.Context, id int) (sirius.Case, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(sirius.Case), args.Error(1)
}

func TestGetCreateInvestigation(t *testing.T) {
	for _, caseType := range []string{"lpa", "epa"} {
		t.Run(caseType, func(t *testing.T) {
			caseItem := sirius.Case{CaseType: caseType, UID: "7000"}
			client := &mockCreateInvestigationClient{}
			client.
				On("Case", mock.Anything, 123).
				Return(caseItem, nil)

			template := &mockTemplate{}
			template.
				On("Func", mock.Anything, createInvestigationData{
					Case: caseItem,
				}).
				Return(nil)

			r, _ := http.NewRequest(http.MethodGet, "/?id=123&case="+caseType, nil)
			w := httptest.NewRecorder()

			err := CreateInvestigation(client, template.Func)(w, r)
			resp := w.Result()

			assert.Nil(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			mock.AssertExpectationsForObjects(t, client, template)
		})
	}
}

func TestGetCreateInvestigationBadQuery(t *testing.T) {
	testCases := map[string]string{
		"no-id":    "/?case=lpa",
		"no-case":  "/?id=123",
		"bad-case": "/?id=123&case=person",
	}

	for name, url := range testCases {
		t.Run(name, func(t *testing.T) {
			r, _ := http.NewRequest(http.MethodGet, url, nil)
			w := httptest.NewRecorder()

			err := CreateInvestigation(nil, nil)(w, r)

			assert.NotNil(t, err)
		})
	}
}

func TestGetCreateInvestigationWhenCaseErrors(t *testing.T) {
	client := &mockCreateInvestigationClient{}
	client.
		On("Case", mock.Anything, 123).
		Return(sirius.Case{}, expectedError)

	r, _ := http.NewRequest(http.MethodGet, "/?id=123&case=lpa", nil)
	w := httptest.NewRecorder()

	err := CreateInvestigation(client, nil)(w, r)

	assert.Equal(t, expectedError, err)
	mock.AssertExpectationsForObjects(t, client)
}

func TestGetCreateInvestigationWhenTemplateErrors(t *testing.T) {
	caseItem := sirius.Case{CaseType: "lpa", UID: "7000"}

	client := &mockCreateInvestigationClient{}
	client.
		On("Case", mock.Anything, 123).
		Return(caseItem, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, createInvestigationData{
			Case: caseItem,
		}).
		Return(expectedError)

	r, _ := http.NewRequest(http.MethodGet, "/?id=123&case=lpa", nil)
	w := httptest.NewRecorder()

	err := CreateInvestigation(client, template.Func)(w, r)

	assert.Equal(t, expectedError, err)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestPostCreateInvestigation(t *testing.T) {
	for _, caseType := range []string{"lpa", "epa"} {
		t.Run(caseType, func(t *testing.T) {
			caseItem := sirius.Case{CaseType: caseType, UID: "7000"}
			client := &mockCreateInvestigationClient{}
			client.
				On("Case", mock.Anything, 123).
				Return(caseItem, nil)
			client.
				On("CreateInvestigation", mock.Anything, 123, sirius.CaseType(caseType), sirius.Investigation{
					Title:        "Test Investigation",
					Information:  "This is an investigation",
					Type:         "Priority",
					DateReceived: "2022-04-05",
				}).
				Return(nil)

			template := &mockTemplate{}
			template.
				On("Func", mock.Anything, createInvestigationData{
					Success: true,
					Case:    caseItem,
				}).
				Return(nil)

			form := url.Values{
				"title":        {"Test Investigation"},
				"information":  {"This is an investigation"},
				"type":         {"Priority"},
				"dateReceived": {"2022-04-05"},
			}

			r, _ := http.NewRequest(http.MethodPost, "/?id=123&case="+caseType, strings.NewReader(form.Encode()))
			r.Header.Add("Content-Type", formUrlEncoded)
			w := httptest.NewRecorder()

			err := CreateInvestigation(client, template.Func)(w, r)
			resp := w.Result()

			assert.Nil(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			mock.AssertExpectationsForObjects(t, client, template)
		})
	}
}

func TestPostCreateInvestigationWhenValidationError(t *testing.T) {
	expectedError := sirius.ValidationError{
		Field: sirius.FieldErrors{"field": {"": "problem"}},
	}

	caseItem := sirius.Case{CaseType: "lpa", UID: "7000"}
	investigation := sirius.Investigation{
		Type: "Priority",
	}

	client := &mockCreateInvestigationClient{}
	client.
		On("Case", mock.Anything, 123).
		Return(caseItem, nil)
	client.
		On("CreateInvestigation", mock.Anything, 123, sirius.CaseTypeLpa, investigation).
		Return(expectedError)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, createInvestigationData{
			Success:       false,
			Error:         expectedError,
			Case:          caseItem,
			Investigation: investigation,
		}).
		Return(nil)

	form := url.Values{
		"type": {"Priority"},
	}

	r, _ := http.NewRequest(http.MethodPost, "/?id=123&case=lpa", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := CreateInvestigation(client, template.Func)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestPostCreateInvestigationWhenOtherError(t *testing.T) {
	caseItem := sirius.Case{CaseType: "lpa", UID: "7000"}
	investigation := sirius.Investigation{
		Type: "Priority",
	}

	client := &mockCreateInvestigationClient{}
	client.
		On("Case", mock.Anything, 123).
		Return(caseItem, nil)
	client.
		On("CreateInvestigation", mock.Anything, 123, sirius.CaseTypeLpa, investigation).
		Return(expectedError)

	form := url.Values{
		"type": {"Priority"},
	}

	r, _ := http.NewRequest(http.MethodPost, "/?id=123&case=lpa", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := CreateInvestigation(client, nil)(w, r)

	assert.Equal(t, expectedError, err)
	mock.AssertExpectationsForObjects(t, client)
}
