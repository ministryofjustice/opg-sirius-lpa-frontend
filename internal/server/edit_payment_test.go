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

type mockEditPaymentClient struct {
	mock.Mock
}

func (m *mockEditPaymentClient) PaymentByID(ctx sirius.Context, id int) (sirius.PaymentDetails, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(sirius.PaymentDetails), args.Error(1)
}

func (m *mockEditPaymentClient) EditPayment(ctx sirius.Context, paymentID int, payment sirius.Payment) error {
	return m.Called(ctx, paymentID, payment).Error(0)
}

func (m *mockEditPaymentClient) Case(ctx sirius.Context, id int) (sirius.Case, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(sirius.Case), args.Error(1)
}

func TestGetEditPayment(t *testing.T) {
	caseItem := sirius.Case{
		UID:     "7000-0000-0021",
		SubType: "pfa",
	}

	paymentDetails := sirius.PaymentDetails{
		CaseId: 4,
		Payment: sirius.Payment{
			ID:          123,
			Amount:      8200,
			Source:      "PHONE",
			PaymentDate: sirius.DateString("2022-07-23"),
		},
	}

	client := &mockEditPaymentClient{}
	client.
		On("Case", mock.Anything, 4).
		Return(caseItem, nil)
	client.
		On("PaymentByID", mock.Anything, 123).
		Return(paymentDetails, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, editPaymentData{
			Case:        caseItem,
			PaymentID:   123,
			Amount:      "82.00",
			Source:      "PHONE",
			PaymentDate: sirius.DateString("2022-07-23"),
		}).
		Return(nil)

	r, _ := http.NewRequest(http.MethodGet, "/?id=4&payment=123", nil)
	w := httptest.NewRecorder()

	err := EditPayment(client, template.Func)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestEditPaymentInvalidURLParams(t *testing.T) {
	testCases := map[string]string{
		"no-params":      "/",
		"no-case-id":     "/?payment=123",
		"no-payment-id":  "/?id=2",
		"bad-case- id":   "/?id=test&payment=123",
		"bad-payment-id": "/?id=2&payment=test",
	}

	for name, testUrl := range testCases {
		t.Run(name, func(t *testing.T) {
			r, _ := http.NewRequest(http.MethodGet, testUrl, nil)
			w := httptest.NewRecorder()

			err := EditPayment(nil, nil)(w, r)

			assert.NotNil(t, err)
		})
	}
}

func TestEditPaymentWhenFailureOnGetPaymentByID(t *testing.T) {
	expectedError := errors.New("err")

	client := &mockEditPaymentClient{}
	client.
		On("PaymentByID", mock.Anything, 123).
		Return(sirius.PaymentDetails{}, expectedError)

	r, _ := http.NewRequest(http.MethodGet, "/?id=4&payment=123", nil)
	w := httptest.NewRecorder()

	err := EditPayment(client, nil)(w, r)

	assert.Equal(t, expectedError, err)
	mock.AssertExpectationsForObjects(t, client)
}

func TestEditPaymentWhenFailureOnGetCase(t *testing.T) {
	expectedError := errors.New("err")

	paymentDetails := sirius.PaymentDetails{
		CaseId: 4,
		Payment: sirius.Payment{
			ID:          123,
			Amount:      8200,
			Source:      "PHONE",
			PaymentDate: sirius.DateString("2022-07-23"),
		},
	}

	client := &mockEditPaymentClient{}
	client.
		On("PaymentByID", mock.Anything, 123).
		Return(paymentDetails, nil)
	client.
		On("Case", mock.Anything, 4).
		Return(sirius.Case{}, expectedError)

	r, _ := http.NewRequest(http.MethodGet, "/?id=4&payment=123", nil)
	w := httptest.NewRecorder()

	err := EditPayment(client, nil)(w, r)

	assert.Equal(t, expectedError, err)
	mock.AssertExpectationsForObjects(t, client)
}

func TestEditPaymentWhenTemplateErrors(t *testing.T) {
	caseItem := sirius.Case{
		UID:     "7000-0000-0021",
		SubType: "pfa",
	}

	paymentDetails := sirius.PaymentDetails{
		CaseId: 4,
		Payment: sirius.Payment{
			ID:          123,
			Amount:      8200,
			Source:      "PHONE",
			PaymentDate: sirius.DateString("2022-07-23"),
		},
	}

	client := &mockEditPaymentClient{}
	client.
		On("PaymentByID", mock.Anything, 123).
		Return(paymentDetails, nil)
	client.
		On("Case", mock.Anything, 4).
		Return(caseItem, nil)

	expectedError := errors.New("err")

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, editPaymentData{
			Case:        caseItem,
			PaymentID:   123,
			Amount:      "82.00",
			Source:      "PHONE",
			PaymentDate: sirius.DateString("2022-07-23"),
		}).
		Return(expectedError)

	r, _ := http.NewRequest(http.MethodGet, "/?id=4&payment=123", nil)
	w := httptest.NewRecorder()

	err := EditPayment(client, template.Func)(w, r)

	assert.Equal(t, expectedError, err)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestPostEditPaymentAmountIncorrectFormat(t *testing.T) {
	for _, amount := range []string{"41", "41.5", "41.555", ".45"} {
		t.Run(amount, func(t *testing.T) {
			caseItem := sirius.Case{CaseType: "lpa", UID: "700700"}

			paymentDetails := sirius.PaymentDetails{
				CaseId: 4,
				Payment: sirius.Payment{
					ID:          123,
					Amount:      8200,
					Source:      "PHONE",
					PaymentDate: sirius.DateString("2022-07-23"),
				},
			}

			client := &mockEditPaymentClient{}
			client.
				On("PaymentByID", mock.Anything, 123).
				Return(paymentDetails, nil)
			client.
				On("Case", mock.Anything, 4).
				Return(caseItem, nil)

			validationError := sirius.ValidationError{
				Field: sirius.FieldErrors{
					"amount": {"reason": "Please enter the amount to 2 decimal places"},
				},
			}

			template := &mockTemplate{}
			template.
				On("Func", mock.Anything, editPaymentData{
					Success:     false,
					Case:        caseItem,
					PaymentID:   123,
					Amount:      amount,
					Source:      "MAKE",
					PaymentDate: sirius.DateString("2022-01-23"),
					Error:       validationError,
				}).
				Return(nil)

			form := url.Values{
				"amount":      {amount},
				"source":      {"MAKE"},
				"paymentDate": {"2022-01-23"},
			}

			r, _ := http.NewRequest(http.MethodPost, "/?id=4&payment=123", strings.NewReader(form.Encode()))
			r.Header.Add("Content-Type", formUrlEncoded)
			w := httptest.NewRecorder()

			err := EditPayment(client, template.Func)(w, r)
			resp := w.Result()

			assert.Nil(t, err)
			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
			mock.AssertExpectationsForObjects(t, client, template)
		})
	}
}

func TestPostEditPayment(t *testing.T) {
	caseItem := sirius.Case{CaseType: "lpa", UID: "700700"}

	paymentDetails := sirius.PaymentDetails{
		CaseId: 4,
		Payment: sirius.Payment{
			ID:          123,
			Amount:      8200,
			Source:      "PHONE",
			PaymentDate: sirius.DateString("2022-02-18"),
		},
	}

	editedPayment := sirius.Payment{
		Amount:      3300,
		Source:      "PHONE",
		PaymentDate: sirius.DateString("2022-02-18"),
	}

	client := &mockEditPaymentClient{}
	client.
		On("PaymentByID", mock.Anything, 123).
		Return(paymentDetails, nil)
	client.
		On("Case", mock.Anything, 4).
		Return(caseItem, nil)
	client.
		On("EditPayment", mock.Anything, 123, editedPayment).
		Return(nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, editPaymentData{
			Success:     true,
			Case:        caseItem,
			PaymentID:   123,
			Amount:      "33.00",
			Source:      "PHONE",
			PaymentDate: sirius.DateString("2022-02-18"),
		}).
		Return(nil)

	form := url.Values{
		"amount":      {"33.00"},
		"source":      {"PHONE"},
		"paymentDate": {"2022-02-18"},
	}

	r, _ := http.NewRequest(http.MethodPost, "/?id=4&payment=123", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := EditPayment(client, template.Func)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}
