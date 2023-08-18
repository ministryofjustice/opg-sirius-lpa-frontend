package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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

func (m *mockDeletePaymentClient) RefDataByCategory(ctx sirius.Context, category string) ([]sirius.RefDataItem, error) {
	args := m.Called(ctx, category)
	if args.Get(0) != nil {
		return args.Get(0).([]sirius.RefDataItem), args.Error(1)
	}
	return nil, args.Error(1)
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

	feeReductionTypes := []sirius.RefDataItem{
		{
			Handle: "REMISSION",
			Label:  "Remission",
		},
	}

	client := &mockDeletePaymentClient{}
	client.
		On("PaymentByID", mock.Anything, 123).
		Return(payment, nil)
	client.
		On("Case", mock.Anything, 4).
		Return(caseItem, nil)
	client.
		On("RefDataByCategory", mock.Anything, sirius.FeeReductionTypeCategory).
		Return(feeReductionTypes, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, deletePaymentData{
			Case:              caseItem,
			Payment:           payment,
			FeeReductionTypes: feeReductionTypes,
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
	payment := sirius.Payment{
		ID:          123,
		Amount:      8200,
		Source:      "PHONE",
		PaymentDate: sirius.DateString("2022-07-23"),
		Case:        &sirius.Case{ID: 4},
	}

	feeReductionTypes := []sirius.RefDataItem{
		{
			Handle: "REMISSION",
			Label:  "Remission",
		},
	}

	client := &mockDeletePaymentClient{}
	client.
		On("PaymentByID", mock.Anything, 123).
		Return(payment, nil)
	client.
		On("Case", mock.Anything, 4).
		Return(sirius.Case{}, expectedError)
	client.
		On("RefDataByCategory", mock.Anything, sirius.FeeReductionTypeCategory).
		Return(feeReductionTypes, nil)

	r, _ := http.NewRequest(http.MethodGet, "/?id=123", nil)
	w := httptest.NewRecorder()

	err := DeletePayment(client, nil)(w, r)

	assert.Equal(t, expectedError, err)
	mock.AssertExpectationsForObjects(t, client)
}

func TestDeletePaymentWhenFailureOnGetFeeReductionTypes(t *testing.T) {
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
	client.
		On("RefDataByCategory", mock.Anything, sirius.FeeReductionTypeCategory).
		Return([]sirius.RefDataItem{}, expectedError)

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

	feeReductionTypes := []sirius.RefDataItem{
		{
			Handle: "REMISSION",
			Label:  "Remission",
		},
	}

	client := &mockDeletePaymentClient{}
	client.
		On("PaymentByID", mock.Anything, 123).
		Return(payment, nil)
	client.
		On("Case", mock.Anything, 4).
		Return(caseItem, nil)
	client.
		On("RefDataByCategory", mock.Anything, sirius.FeeReductionTypeCategory).
		Return(feeReductionTypes, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, deletePaymentData{
			Case:              caseItem,
			Payment:           payment,
			FeeReductionTypes: feeReductionTypes,
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

	feeReductionTypes := []sirius.RefDataItem{
		{
			Handle: "REMISSION",
			Label:  "Remission",
		},
	}

	client := &mockDeletePaymentClient{}
	client.
		On("PaymentByID", mock.Anything, 123).
		Return(payment, nil)
	client.
		On("Case", mock.Anything, 4).
		Return(caseItem, nil)
	client.
		On("RefDataByCategory", mock.Anything, sirius.FeeReductionTypeCategory).
		Return(feeReductionTypes, nil)
	client.
		On("DeletePayment", mock.Anything, 123).
		Return(nil)

	template := &mockTemplate{}

	r, _ := http.NewRequest(http.MethodPost, "/?id=123", nil)
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := DeletePayment(client, template.Func)(w, r)
	resp := w.Result()

	assert.Equal(t, RedirectError("/payments/4"), err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}
