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

type mockEditDatesClient struct {
	mock.Mock
}

func (m *mockEditDatesClient) Case(ctx sirius.Context, id int) (sirius.Case, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(sirius.Case), args.Error(1)
}

func (m *mockEditDatesClient) EditDates(ctx sirius.Context, caseID int, caseType sirius.CaseType, dates sirius.Dates) error {
	return m.Called(ctx, caseID, caseType, dates).Error(0)
}

func TestGetEditDates(t *testing.T) {
	for _, caseType := range []string{"lpa", "epa"} {
		t.Run(caseType, func(t *testing.T) {
			caseitem := sirius.Case{CaseType: caseType, UID: "700700"}

			client := &mockEditDatesClient{}
			client.
				On("Case", mock.Anything, 123).
				Return(caseitem, nil)

			template := &mockTemplate{}
			template.
				On("Func", mock.Anything, editDatesData{
					Entity: caseType + " 700700",
					Case:   caseitem,
				}).
				Return(nil)

			r, _ := http.NewRequest(http.MethodGet, "/?id=123&case="+caseType, nil)
			w := httptest.NewRecorder()

			err := EditDates(client, template.Func)(w, r)
			resp := w.Result()

			assert.Nil(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			mock.AssertExpectationsForObjects(t, client, template)
		})
	}
}

func TestGetEditDatesNoID(t *testing.T) {
	testCases := map[string]string{
		"no-id":    "/?case=lpa",
		"no-case":  "/?id=123",
		"bad-case": "/?id=123&case=person",
	}

	for name, url := range testCases {
		t.Run(name, func(t *testing.T) {
			r, _ := http.NewRequest(http.MethodGet, url, nil)
			w := httptest.NewRecorder()

			err := EditDates(nil, nil)(w, r)

			assert.NotNil(t, err)
		})
	}
}

func TestGetEditDatesWhenCaseErrors(t *testing.T) {
	expectedError := errors.New("err")
	caseitem := sirius.Case{CaseType: "PFA", UID: "700700"}

	client := &mockEditDatesClient{}
	client.
		On("Case", mock.Anything, 123).
		Return(caseitem, expectedError)

	r, _ := http.NewRequest(http.MethodGet, "/?id=123&case=lpa", nil)
	w := httptest.NewRecorder()

	err := EditDates(client, nil)(w, r)

	assert.Equal(t, expectedError, err)
	mock.AssertExpectationsForObjects(t, client)
}

func TestGetEditDatesWhenTemplateErrors(t *testing.T) {
	expectedError := errors.New("err")
	caseitem := sirius.Case{CaseType: "PFA", UID: "700700"}

	client := &mockEditDatesClient{}
	client.
		On("Case", mock.Anything, 123).
		Return(caseitem, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, editDatesData{
			Entity: "PFA 700700",
			Case:   caseitem,
		}).
		Return(expectedError)

	r, _ := http.NewRequest(http.MethodGet, "/?id=123&case=lpa", nil)
	w := httptest.NewRecorder()

	err := EditDates(client, template.Func)(w, r)

	assert.Equal(t, expectedError, err)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestPostEditDates(t *testing.T) {
	for _, caseType := range []string{"lpa", "epa"} {
		t.Run(caseType, func(t *testing.T) {
			caseitem := sirius.Case{CaseType: caseType, UID: "700700"}

			client := &mockEditDatesClient{}
			client.
				On("EditDates", mock.Anything, 123, sirius.CaseType(caseType), sirius.Dates{
					CancellationDate: sirius.DateString("2022-01-03"),
					DispatchDate:     sirius.DateString("2022-02-03"),
					DueDate:          sirius.DateString("2022-03-05"),
					InvalidDate:      sirius.DateString("2022-04-03"),
					ReceiptDate:      sirius.DateString("2021-11-23"),
					RegistrationDate: sirius.DateString("2022-05-03"),
					RejectedDate:     sirius.DateString("2022-06-03"),
					WithdrawnDate:    sirius.DateString("2022-07-03"),
				}).
				Return(nil)
			client.
				On("Case", mock.Anything, 123).
				Return(caseitem, nil)

			template := &mockTemplate{}
			template.
				On("Func", mock.Anything, editDatesData{
					Success: true,
					Entity:  caseType + " 700700",
					Case:    caseitem,
				}).
				Return(nil)

			form := url.Values{
				"cancellationDate": {"2022-01-03"},
				"dispatchDate":     {"2022-02-03"},
				"dueDate":          {"2022-03-05"},
				"invalidDate":      {"2022-04-03"},
				"receiptDate":      {"2021-11-23"},
				"registrationDate": {"2022-05-03"},
				"rejectedDate":     {"2022-06-03"},
				"withdrawnDate":    {"2022-07-03"},
			}

			r, _ := http.NewRequest(http.MethodPost, "/?id=123&case="+caseType, strings.NewReader(form.Encode()))
			r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
			w := httptest.NewRecorder()

			err := EditDates(client, template.Func)(w, r)
			resp := w.Result()

			assert.Nil(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			mock.AssertExpectationsForObjects(t, client, template)
		})
	}
}

func TestPostEditDatesWhenEditDatesErrors(t *testing.T) {
	expectedError := errors.New("err")

	client := &mockEditDatesClient{}
	client.
		On("EditDates", mock.Anything, 123, sirius.CaseTypeLpa, sirius.Dates{
			RegistrationDate: sirius.DateString("2022-01-03"),
		}).
		Return(expectedError)

	form := url.Values{
		"registrationDate": {"2022-01-03"},
	}

	r, _ := http.NewRequest(http.MethodPost, "/?id=123&case=lpa", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	err := EditDates(client, nil)(w, r)

	assert.Equal(t, expectedError, err)
	mock.AssertExpectationsForObjects(t, client)
}
