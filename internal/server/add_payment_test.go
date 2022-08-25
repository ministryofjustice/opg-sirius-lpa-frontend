package server

import (
	"errors"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

type mockAddPaymentClient struct {
	mock.Mock
}

func (m *mockAddPaymentClient) AddPayment(ctx sirius.Context, caseID int, amount int, source string, paymentDate sirius.DateString) error {
	return m.Called(ctx, caseID, amount, source, paymentDate).Error(0)
}

func (m *mockAddPaymentClient) Case(ctx sirius.Context, id int) (sirius.Case, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(sirius.Case), args.Error(1)
}

func TestGetAddPayment(t *testing.T) {
	caseItem := sirius.Case{
		UID:     "7000-0000-0021",
		SubType: "pfa",
	}

	client := &mockAddPaymentClient{}
	client.
		On("Case", mock.Anything, 4).
		Return(caseItem, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, addPaymentData{
			Case: caseItem,
		}).
		Return(nil)

	r, _ := http.NewRequest(http.MethodGet, "/?id=4", nil)
	w := httptest.NewRecorder()

	err := AddPayment(client, template.Func)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestAddPaymentNoID(t *testing.T) {
	testCases := map[string]string{
		"no-id":  "/",
		"bad-id": "/?id=test",
	}

	for name, testUrl := range testCases {
		t.Run(name, func(t *testing.T) {
			r, _ := http.NewRequest(http.MethodGet, testUrl, nil)
			w := httptest.NewRecorder()

			err := AddPayment(nil, nil)(w, r)

			assert.NotNil(t, err)
		})
	}
}

func TestAddPaymentWhenFailureOnGetCase(t *testing.T) {
	expectedError := errors.New("err")

	client := &mockAddPaymentClient{}
	client.
		On("Case", mock.Anything, 4).
		Return(sirius.Case{}, expectedError)

	r, _ := http.NewRequest(http.MethodGet, "/?id=4", nil)
	w := httptest.NewRecorder()

	err := AddPayment(client, nil)(w, r)

	assert.Equal(t, expectedError, err)
	mock.AssertExpectationsForObjects(t, client)
}

func TestAddPaymentWhenTemplateErrors(t *testing.T) {
	caseItem := sirius.Case{
		UID:     "7000-0000-0021",
		SubType: "pfa",
	}

	client := &mockAddPaymentClient{}
	client.
		On("Case", mock.Anything, 123).
		Return(caseItem, nil)

	expectedError := errors.New("err")

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, addPaymentData{
			Case:        caseItem,
			Amount:      4100,
			Source:      "MAKE",
			PaymentDate: sirius.DateString("2022-01-23"),
		}).
		Return(expectedError)

	r, _ := http.NewRequest(http.MethodGet, "/?id=123", nil)
	w := httptest.NewRecorder()

	err := AddPayment(client, template.Func)(w, r)

	assert.Equal(t, expectedError, err)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestPostAddPayment(t *testing.T) {
	caseitem := sirius.Case{CaseType: "lpa", UID: "700700"}

	client := &mockAddPaymentClient{}
	client.
		On("AddPayment", mock.Anything, 123, 4100, "MAKE", sirius.DateString("2022-01-23")).
		Return(nil)

	client.
		On("Case", mock.Anything, 123).
		Return(caseitem, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, addPaymentData{
			Success:     true,
			Case:        caseitem,
			Amount:      4100,
			Source:      "MAKE",
			PaymentDate: sirius.DateString("2022-01-23"),
		}).
		Return(nil)

	form := url.Values{
		"amount":      {"41.00"},
		"source":      {"MAKE"},
		"paymentDate": {"2022-01-23"},
	}

	r, _ := http.NewRequest(http.MethodPost, "/?id=123", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := AddPayment(client, template.Func)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}
