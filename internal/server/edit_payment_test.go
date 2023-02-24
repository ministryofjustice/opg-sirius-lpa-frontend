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

type mockEditPaymentClient struct {
	mock.Mock
}

func (m *mockEditPaymentClient) PaymentByID(ctx sirius.Context, id int) (sirius.Payment, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(sirius.Payment), args.Error(1)
}

func (m *mockEditPaymentClient) EditPayment(ctx sirius.Context, paymentID int, payment sirius.Payment) error {
	return m.Called(ctx, paymentID, payment).Error(0)
}

func (m *mockEditPaymentClient) Case(ctx sirius.Context, id int) (sirius.Case, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(sirius.Case), args.Error(1)
}

func (m *mockEditPaymentClient) RefDataByCategory(ctx sirius.Context, category string) ([]sirius.RefDataItem, error) {
	args := m.Called(ctx, category)
	if args.Get(0) != nil {
		return args.Get(0).([]sirius.RefDataItem), args.Error(1)
	}
	return nil, args.Error(1)
}

func TestGetEditPayment(t *testing.T) {
	caseItem := sirius.Case{
		UID:     "7000-0000-0021",
		SubType: "pfa",
	}

	payment := sirius.Payment{
		ID:          123,
		Amount:      8200,
		Source:      "PHONE",
		PaymentDate: sirius.DateString("2022-07-23"),
		Case:        &sirius.Case{ID: 4},
	}

	paymentSources := []sirius.RefDataItem{
		{
			Handle:         "PHONE",
			Label:          "Paid over the phone",
			UserSelectable: true,
		},
	}

	client := &mockEditPaymentClient{}
	client.
		On("Case", mock.Anything, 4).
		Return(caseItem, nil)
	client.
		On("PaymentByID", mock.Anything, 123).
		Return(payment, nil)
	client.
		On("RefDataByCategory", mock.Anything, sirius.PaymentSourceCategory).
		Return(paymentSources, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, editPaymentData{
			Case:           caseItem,
			PaymentID:      123,
			Amount:         "82.00",
			Source:         "PHONE",
			PaymentDate:    sirius.DateString("2022-07-23"),
			PaymentSources: paymentSources,
		}).
		Return(nil)

	r, _ := http.NewRequest(http.MethodGet, "/?id=123", nil)
	w := httptest.NewRecorder()

	err := EditPayment(client, template.Func)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestEditPaymentWhenFailureOnGetPaymentByID(t *testing.T) {
	paymentSources := []sirius.RefDataItem{
		{
			Handle:         "PHONE",
			Label:          "Paid over the phone",
			UserSelectable: true,
		},
	}
	client := &mockEditPaymentClient{}
	client.
		On("PaymentByID", mock.Anything, 123).
		Return(sirius.Payment{}, expectedError).
		On("RefDataByCategory", mock.Anything, sirius.PaymentSourceCategory).
		Return(paymentSources, nil)

	r, _ := http.NewRequest(http.MethodGet, "/?id=123", nil)
	w := httptest.NewRecorder()

	err := EditPayment(client, nil)(w, r)

	assert.Equal(t, expectedError, err)
	mock.AssertExpectationsForObjects(t, client)
}

func TestEditPaymentWhenFailureOnGetCase(t *testing.T) {
	payment := sirius.Payment{
		ID:          123,
		Amount:      8200,
		Source:      "PHONE",
		PaymentDate: sirius.DateString("2022-07-23"),
		Case:        &sirius.Case{ID: 4},
	}

	paymentSources := []sirius.RefDataItem{
		{
			Handle:         "PHONE",
			Label:          "Paid over the phone",
			UserSelectable: true,
		},
	}

	client := &mockEditPaymentClient{}
	client.
		On("PaymentByID", mock.Anything, 123).
		Return(payment, nil).
		On("RefDataByCategory", mock.Anything, sirius.PaymentSourceCategory).
		Return(paymentSources, nil).
		On("Case", mock.Anything, 4).
		Return(sirius.Case{}, expectedError)

	r, _ := http.NewRequest(http.MethodGet, "/?id=123", nil)
	w := httptest.NewRecorder()

	err := EditPayment(client, nil)(w, r)

	assert.Equal(t, expectedError, err)
	mock.AssertExpectationsForObjects(t, client)
}

func TestEditPaymentWhenFailureOnGetPaymentSourceRefData(t *testing.T) {
	caseItem := sirius.Case{
		UID:     "7000-0000-0021",
		SubType: "pfa",
	}

	payment := sirius.Payment{
		ID:          123,
		Amount:      8200,
		Source:      "PHONE",
		PaymentDate: sirius.DateString("2022-07-23"),
		Case:        &sirius.Case{ID: 4},
	}

	client := &mockEditPaymentClient{}
	client.
		On("Case", mock.Anything, 4).
		Return(caseItem, nil)
	client.
		On("PaymentByID", mock.Anything, 123).
		Return(payment, nil)
	client.
		On("RefDataByCategory", mock.Anything, sirius.PaymentSourceCategory).
		Return([]sirius.RefDataItem{}, expectedError)

	r, _ := http.NewRequest(http.MethodGet, "/?id=123", nil)
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

	payment := sirius.Payment{
		ID:          123,
		Amount:      8200,
		Source:      "PHONE",
		PaymentDate: sirius.DateString("2022-07-23"),
		Case:        &sirius.Case{ID: 4},
	}

	paymentSources := []sirius.RefDataItem{
		{
			Handle:         "PHONE",
			Label:          "Paid over the phone",
			UserSelectable: true,
		},
	}

	client := &mockEditPaymentClient{}
	client.
		On("PaymentByID", mock.Anything, 123).
		Return(payment, nil)
	client.
		On("Case", mock.Anything, 4).
		Return(caseItem, nil)
	client.
		On("RefDataByCategory", mock.Anything, sirius.PaymentSourceCategory).
		Return(paymentSources, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, editPaymentData{
			Case:           caseItem,
			PaymentID:      123,
			Amount:         "82.00",
			Source:         "PHONE",
			PaymentDate:    sirius.DateString("2022-07-23"),
			PaymentSources: paymentSources,
		}).
		Return(expectedError)

	r, _ := http.NewRequest(http.MethodGet, "/?id=123", nil)
	w := httptest.NewRecorder()

	err := EditPayment(client, template.Func)(w, r)

	assert.Equal(t, expectedError, err)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestPostEditPaymentAmountIncorrectFormat(t *testing.T) {
	for _, amount := range []string{"41", "41.5", "41.555", ".45"} {
		t.Run(amount, func(t *testing.T) {
			caseItem := sirius.Case{CaseType: "lpa", UID: "700700"}

			payment := sirius.Payment{
				ID:          123,
				Amount:      8200,
				Source:      "PHONE",
				PaymentDate: sirius.DateString("2022-07-23"),
				Case:        &sirius.Case{ID: 4},
			}

			paymentSources := []sirius.RefDataItem{
				{
					Handle:         "PHONE",
					Label:          "Paid over the phone",
					UserSelectable: true,
				},
			}

			client := &mockEditPaymentClient{}
			client.
				On("PaymentByID", mock.Anything, 123).
				Return(payment, nil)
			client.
				On("Case", mock.Anything, 4).
				Return(caseItem, nil)
			client.
				On("RefDataByCategory", mock.Anything, sirius.PaymentSourceCategory).
				Return(paymentSources, nil)

			validationError := sirius.ValidationError{
				Field: sirius.FieldErrors{
					"amount": {"reason": "Please enter the amount to 2 decimal places"},
				},
			}

			template := &mockTemplate{}
			template.
				On("Func", mock.Anything, editPaymentData{
					Case:           caseItem,
					PaymentID:      123,
					Amount:         amount,
					Source:         "MAKE",
					PaymentDate:    sirius.DateString("2022-01-23"),
					PaymentSources: paymentSources,
					Error:          validationError,
				}).
				Return(nil)

			form := url.Values{
				"amount":      {amount},
				"source":      {"MAKE"},
				"paymentDate": {"2022-01-23"},
			}

			r, _ := http.NewRequest(http.MethodPost, "/?id=123", strings.NewReader(form.Encode()))
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
	caseItem := sirius.Case{CaseType: "lpa", UID: "700700", ID: 4}

	payment := sirius.Payment{
		ID:          123,
		Amount:      8200,
		Source:      "PHONE",
		PaymentDate: sirius.DateString("2022-02-18"),
		Case:        &sirius.Case{ID: 4},
	}

	editedPayment := sirius.Payment{
		Amount:      3300,
		Source:      "PHONE",
		PaymentDate: sirius.DateString("2022-02-18"),
	}

	paymentSources := []sirius.RefDataItem{
		{
			Handle:         "PHONE",
			Label:          "Paid over the phone",
			UserSelectable: true,
		},
	}

	client := &mockEditPaymentClient{}
	client.
		On("PaymentByID", mock.Anything, 123).
		Return(payment, nil)
	client.
		On("Case", mock.Anything, 4).
		Return(caseItem, nil)
	client.
		On("RefDataByCategory", mock.Anything, sirius.PaymentSourceCategory).
		Return(paymentSources, nil)
	client.
		On("EditPayment", mock.Anything, 123, editedPayment).
		Return(nil)

	template := &mockTemplate{}

	form := url.Values{
		"amount":      {"33.00"},
		"source":      {"PHONE"},
		"paymentDate": {"2022-02-18"},
	}

	r, _ := http.NewRequest(http.MethodPost, "/?id=123", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := EditPayment(client, template.Func)(w, r)
	resp := w.Result()

	assert.Equal(t, RedirectError("/payments?id=4"), err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}
