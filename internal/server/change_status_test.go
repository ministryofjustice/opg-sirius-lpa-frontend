package server

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockChangeStatusClient struct {
	mock.Mock
}

func (m *mockChangeStatusClient) Case(ctx sirius.Context, id int) (sirius.Case, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(sirius.Case), args.Error(1)
}

func (m *mockChangeStatusClient) EditCase(ctx sirius.Context, caseID int, caseType sirius.CaseType, caseData sirius.Case) error {
	return m.Called(ctx, caseID, caseType, caseData).Error(0)
}

func (m *mockChangeStatusClient) AvailableStatuses(ctx sirius.Context, caseID int, caseType sirius.CaseType) ([]string, error) {
	args := m.Called(ctx, caseID, caseType)
	return args.Get(0).([]string), args.Error(1)
}

func TestGetChangeStatus(t *testing.T) {
	for _, caseType := range []string{"lpa", "epa"} {
		t.Run(caseType, func(t *testing.T) {
			caseitem := sirius.Case{CaseType: caseType, UID: "700700"}

			client := &mockChangeStatusClient{}
			client.
				On("Case", mock.Anything, 123).
				Return(caseitem, nil)

			client.
				On("AvailableStatuses", mock.Anything, 123, sirius.CaseType(caseType)).
				Return([]string{"Cancelled", "Withdrawn"}, nil)

			template := &mockTemplate{}
			template.
				On("Func", mock.Anything, changeStatusData{
					Entity:            caseType + " 700700",
					AvailableStatuses: []string{"Cancelled", "Withdrawn"},
				}).
				Return(nil)

			r, _ := http.NewRequest(http.MethodGet, "/?id=123&case="+caseType, nil)
			w := httptest.NewRecorder()

			err := ChangeStatus(client, template.Func)(w, r)
			resp := w.Result()

			assert.Nil(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			mock.AssertExpectationsForObjects(t, client, template)
		})
	}
}

func TestGetChangeStatusNoID(t *testing.T) {
	testCases := map[string]string{
		"no-id":    "/?case=lpa",
		"no-case":  "/?id=123",
		"bad-case": "/?id=123&case=person",
	}

	for name, url := range testCases {
		t.Run(name, func(t *testing.T) {
			r, _ := http.NewRequest(http.MethodGet, url, nil)
			w := httptest.NewRecorder()

			err := ChangeStatus(nil, nil)(w, r)

			assert.NotNil(t, err)
		})
	}
}

func TestGetChangeStatusWhenCaseErrors(t *testing.T) {
	expectedError := errors.New("err")
	caseitem := sirius.Case{CaseType: "PFA", UID: "700700"}

	client := &mockChangeStatusClient{}
	client.
		On("Case", mock.Anything, 123).
		Return(caseitem, expectedError)

	r, _ := http.NewRequest(http.MethodGet, "/?id=123&case=lpa", nil)
	w := httptest.NewRecorder()

	err := ChangeStatus(client, nil)(w, r)

	assert.Equal(t, expectedError, err)
	mock.AssertExpectationsForObjects(t, client)
}

func TestGetChangeStatusWhenAvailableStatusesErrors(t *testing.T) {
	expectedError := errors.New("err")
	caseitem := sirius.Case{CaseType: "PFA", UID: "700700"}

	client := &mockChangeStatusClient{}
	client.
		On("Case", mock.Anything, 123).
		Return(caseitem, nil)

	client.
		On("AvailableStatuses", mock.Anything, 123, sirius.CaseTypeLpa).
		Return([]string{}, expectedError)

	r, _ := http.NewRequest(http.MethodGet, "/?id=123&case=lpa", nil)
	w := httptest.NewRecorder()

	err := ChangeStatus(client, nil)(w, r)

	assert.Equal(t, expectedError, err)
	mock.AssertExpectationsForObjects(t, client)
}

func TestGetChangeStatusWhenTemplateErrors(t *testing.T) {
	expectedError := errors.New("err")
	caseitem := sirius.Case{CaseType: "PFA", UID: "700700"}

	client := &mockChangeStatusClient{}
	client.
		On("Case", mock.Anything, 123).
		Return(caseitem, nil)

	client.
		On("AvailableStatuses", mock.Anything, 123, sirius.CaseTypeLpa).
		Return([]string{"Cancelled", "Withdrawn"}, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, changeStatusData{
			Entity:            "PFA 700700",
			AvailableStatuses: []string{"Cancelled", "Withdrawn"},
		}).
		Return(expectedError)

	r, _ := http.NewRequest(http.MethodGet, "/?id=123&case=lpa", nil)
	w := httptest.NewRecorder()

	err := ChangeStatus(client, template.Func)(w, r)

	assert.Equal(t, expectedError, err)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestPostChangeStatus(t *testing.T) {
	for _, caseType := range []string{"lpa", "epa"} {
		t.Run(caseType, func(t *testing.T) {
			caseitem := sirius.Case{CaseType: caseType, UID: "700700"}

			client := &mockChangeStatusClient{}
			client.
				On("EditCase", mock.Anything, 123, sirius.CaseType(caseType), sirius.Case{
					Status: "Withdrawn",
				}).
				Return(nil)

			client.
				On("Case", mock.Anything, 123).
				Return(caseitem, nil)

			client.
				On("AvailableStatuses", mock.Anything, 123, sirius.CaseType(caseType)).
				Return([]string{"Cancelled", "Withdrawn"}, nil)

			template := &mockTemplate{}
			template.
				On("Func", mock.Anything, changeStatusData{
					Success:           true,
					Entity:            caseType + " 700700",
					AvailableStatuses: []string{"Cancelled", "Withdrawn"},
					NewStatus:         "Withdrawn",
				}).
				Return(nil)

			form := url.Values{
				"status": {"Withdrawn"},
			}

			r, _ := http.NewRequest(http.MethodPost, "/?id=123&case="+caseType, strings.NewReader(form.Encode()))
			r.Header.Add("Content-Type", formUrlEncoded)
			w := httptest.NewRecorder()

			err := ChangeStatus(client, template.Func)(w, r)
			resp := w.Result()

			assert.Nil(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			mock.AssertExpectationsForObjects(t, client, template)
		})
	}
}

func TestPostChangeStatusWhenChangeStatusErrors(t *testing.T) {
	expectedError := errors.New("err")
	caseitem := sirius.Case{CaseType: "lpa", UID: "700700"}

	client := &mockChangeStatusClient{}
	client.
		On("EditCase", mock.Anything, 123, sirius.CaseTypeLpa, sirius.Case{
			Status: "Withdrawn",
		}).
		Return(expectedError)

	client.
		On("Case", mock.Anything, 123).
		Return(caseitem, nil)

	client.
		On("AvailableStatuses", mock.Anything, 123, sirius.CaseTypeLpa).
		Return([]string{"Cancelled", "Withdrawn"}, nil)

	form := url.Values{
		"status": {"Withdrawn"},
	}

	r, _ := http.NewRequest(http.MethodPost, "/?id=123&case=lpa", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := ChangeStatus(client, nil)(w, r)

	assert.Equal(t, expectedError, err)
	mock.AssertExpectationsForObjects(t, client)
}
