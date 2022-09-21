package server

import (
	"errors"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockDeletePaymentClient struct {
	mock.Mock
}

func (m *mockDeletePaymentClient) PaymentByID(ctx sirius.Context, id int) (sirius.Payment, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(sirius.Payment), args.Error(1)
}

func (m *mockDeletePaymentClient) Case(ctx sirius.Context, id int) (sirius.Case, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(sirius.Case), args.Error(1)
}

func (m *mockDeletePaymentClient) DeletePayment(ctx sirius.Context, paymentID int) error {
	return m.Called(ctx, paymentID).Error(0)
}

func TestGetDeletePayment(t *testing.T) {
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

	client := &mockDeletePaymentClient{}
	client.
		On("PaymentByID", mock.Anything, 123).
		Return(payment, nil)
	client.
		On("Case", mock.Anything, 4).
		Return(caseItem, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, deletePaymentData{
			Case:    caseItem,
			Payment: payment,
		}).
		Return(nil)

	r, _ := http.NewRequest(http.MethodGet, "/?id=123", nil)
	w := httptest.NewRecorder()

	err := DeletePayment(client, template.Func)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestDeletePaymentWhenFailureOnGetPaymentByID(t *testing.T) {
	expectedError := errors.New("err")

	client := &mockDeletePaymentClient{}
	client.
		On("PaymentByID", mock.Anything, 123).
		Return(sirius.Payment{}, expectedError)

	r, _ := http.NewRequest(http.MethodGet, "/?id=123", nil)
	w := httptest.NewRecorder()

	err := DeletePayment(client, nil)(w, r)

	assert.Equal(t, expectedError, err)
	mock.AssertExpectationsForObjects(t, client)
}

func TestDeletePaymentWhenFailureOnGetCase(t *testing.T) {
	expectedError := errors.New("err")

	payment := sirius.Payment{
		ID:          123,
		Amount:      8200,
		Source:      "PHONE",
		PaymentDate: sirius.DateString("2022-07-23"),
		Case:        &sirius.Case{ID: 4},
	}

	client := &mockDeletePaymentClient{}
	client.
		On("PaymentByID", mock.Anything, 123).
		Return(payment, nil)
	client.
		On("Case", mock.Anything, 4).
		Return(sirius.Case{}, expectedError)

	r, _ := http.NewRequest(http.MethodGet, "/?id=123", nil)
	w := httptest.NewRecorder()

	err := DeletePayment(client, nil)(w, r)

	assert.Equal(t, expectedError, err)
	mock.AssertExpectationsForObjects(t, client)
}

func TestDeletePaymentWhenTemplateErrors(t *testing.T) {
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

	client := &mockDeletePaymentClient{}
	client.
		On("PaymentByID", mock.Anything, 123).
		Return(payment, nil)
	client.
		On("Case", mock.Anything, 4).
		Return(caseItem, nil)

	expectedError := errors.New("err")

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, deletePaymentData{
			Case:    caseItem,
			Payment: payment,
		}).
		Return(expectedError)

	r, _ := http.NewRequest(http.MethodGet, "/?id=123", nil)
	w := httptest.NewRecorder()

	err := DeletePayment(client, template.Func)(w, r)

	assert.Equal(t, expectedError, err)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestPostDeletePayment(t *testing.T) {
	caseItem := sirius.Case{CaseType: "lpa", UID: "700700"}

	payment := sirius.Payment{
		ID:          123,
		Amount:      8200,
		Source:      "PHONE",
		PaymentDate: sirius.DateString("2022-02-18"),
		Case:        &sirius.Case{ID: 4},
	}

	client := &mockDeletePaymentClient{}
	client.
		On("PaymentByID", mock.Anything, 123).
		Return(payment, nil)
	client.
		On("Case", mock.Anything, 4).
		Return(caseItem, nil)

	client.
		On("DeletePayment", mock.Anything, 123).
		Return(nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, deletePaymentData{
			Success: true,
			Case:    caseItem,
			Payment: payment,
		}).
		Return(nil)

	r, _ := http.NewRequest(http.MethodPost, "/?id=123", nil)
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := DeletePayment(client, template.Func)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}
