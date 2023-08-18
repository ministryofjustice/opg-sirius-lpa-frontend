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

type mockEditFeeReductionClient struct {
	mock.Mock
}

func (m *mockEditFeeReductionClient) PaymentByID(ctx sirius.Context, id int) (sirius.Payment, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(sirius.Payment), args.Error(1)
}

func (m *mockEditFeeReductionClient) EditPayment(ctx sirius.Context, feeReductionID int, feeReduction sirius.Payment) error {
	return m.Called(ctx, feeReductionID, feeReduction).Error(0)
}

func (m *mockEditFeeReductionClient) Case(ctx sirius.Context, id int) (sirius.Case, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(sirius.Case), args.Error(1)
}

func (m *mockEditFeeReductionClient) RefDataByCategory(ctx sirius.Context, category string) ([]sirius.RefDataItem, error) {
	args := m.Called(ctx, category)
	if args.Get(0) != nil {
		return args.Get(0).([]sirius.RefDataItem), args.Error(1)
	}
	return nil, args.Error(1)
}

func TestGetEditFeeReduction(t *testing.T) {
	caseItem := sirius.Case{
		UID:     "7000-0000-0021",
		SubType: "pfa",
	}

	feeReduction := sirius.Payment{
		ID:               123,
		PaymentEvidence:  "Test evidence",
		FeeReductionType: "REMISSION",
		Source:           sirius.FeeReductionSource,
		PaymentDate:      sirius.DateString("2022-07-23"),
		Case:             &sirius.Case{ID: 4},
	}

	feeReductionTypes := []sirius.RefDataItem{
		{
			Handle: "REMISSION",
			Label:  "Remission",
		},
	}

	client := &mockEditFeeReductionClient{}
	client.
		On("Case", mock.Anything, 4).
		Return(caseItem, nil).
		On("PaymentByID", mock.Anything, 123).
		Return(feeReduction, nil).
		On("RefDataByCategory", mock.Anything, sirius.FeeReductionTypeCategory).
		Return(feeReductionTypes, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, editFeeReductionData{
			Case:              caseItem,
			PaymentID:         123,
			FeeReduction:      feeReduction,
			FeeReductionTypes: feeReductionTypes,
		}).
		Return(nil)

	r, _ := http.NewRequest(http.MethodGet, "/?id=123", nil)
	w := httptest.NewRecorder()

	err := EditFeeReduction(client, template.Func)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestEditFeeReductionWhenFailureOnGetPaymentByID(t *testing.T) {
	feeReductionTypes := []sirius.RefDataItem{
		{
			Handle: "REMISSION",
			Label:  "Remission",
		},
	}

	client := &mockEditFeeReductionClient{}
	client.
		On("PaymentByID", mock.Anything, 123).
		Return(sirius.Payment{}, expectedError).
		On("RefDataByCategory", mock.Anything, sirius.FeeReductionTypeCategory).
		Return(feeReductionTypes, nil)

	r, _ := http.NewRequest(http.MethodGet, "/?id=123", nil)
	w := httptest.NewRecorder()

	err := EditFeeReduction(client, nil)(w, r)

	assert.Equal(t, expectedError, err)
	mock.AssertExpectationsForObjects(t, client)
}

func TestEditFeeReductionWhenFailureOnGetCase(t *testing.T) {
	feeReduction := sirius.Payment{
		ID:               123,
		PaymentEvidence:  "Test evidence",
		FeeReductionType: "REMISSION",
		Source:           sirius.FeeReductionSource,
		PaymentDate:      sirius.DateString("2022-07-23"),
		Case:             &sirius.Case{ID: 4},
	}

	feeReductionTypes := []sirius.RefDataItem{
		{
			Handle: "REMISSION",
			Label:  "Remission",
		},
	}

	client := &mockEditFeeReductionClient{}
	client.
		On("PaymentByID", mock.Anything, 123).
		Return(feeReduction, nil).
		On("RefDataByCategory", mock.Anything, sirius.FeeReductionTypeCategory).
		Return(feeReductionTypes, nil).
		On("Case", mock.Anything, 4).
		Return(sirius.Case{}, expectedError)

	r, _ := http.NewRequest(http.MethodGet, "/?id=123", nil)
	w := httptest.NewRecorder()

	err := EditFeeReduction(client, nil)(w, r)

	assert.Equal(t, expectedError, err)
	mock.AssertExpectationsForObjects(t, client)
}

func TestEditFeeReductionWhenFailureOnGetPaymentSourceRefData(t *testing.T) {
	caseItem := sirius.Case{
		UID:     "7000-0000-0021",
		SubType: "pfa",
	}

	feeReduction := sirius.Payment{
		ID:               123,
		PaymentEvidence:  "Test evidence",
		FeeReductionType: "REMISSION",
		Source:           sirius.FeeReductionSource,
		PaymentDate:      sirius.DateString("2022-07-23"),
		Case:             &sirius.Case{ID: 4},
	}

	client := &mockEditFeeReductionClient{}
	client.
		On("PaymentByID", mock.Anything, 123).
		Return(feeReduction, nil).
		On("Case", mock.Anything, 4).
		Return(caseItem, nil).
		On("RefDataByCategory", mock.Anything, sirius.FeeReductionTypeCategory).
		Return([]sirius.RefDataItem{}, expectedError)

	r, _ := http.NewRequest(http.MethodGet, "/?id=123", nil)
	w := httptest.NewRecorder()

	err := EditFeeReduction(client, nil)(w, r)

	assert.Equal(t, expectedError, err)
	mock.AssertExpectationsForObjects(t, client)
}

func TestEditFeeReductionWhenTemplateErrors(t *testing.T) {
	caseItem := sirius.Case{
		UID:     "7000-0000-0021",
		SubType: "pfa",
	}

	feeReduction := sirius.Payment{
		ID:               123,
		PaymentEvidence:  "Test evidence",
		FeeReductionType: "REMISSION",
		Source:           sirius.FeeReductionSource,
		PaymentDate:      sirius.DateString("2022-07-23"),
		Case:             &sirius.Case{ID: 4},
	}

	feeReductionTypes := []sirius.RefDataItem{
		{
			Handle: "REMISSION",
			Label:  "Remission",
		},
	}

	client := &mockEditFeeReductionClient{}
	client.
		On("PaymentByID", mock.Anything, 123).
		Return(feeReduction, nil).
		On("Case", mock.Anything, 4).
		Return(caseItem, nil).
		On("RefDataByCategory", mock.Anything, sirius.FeeReductionTypeCategory).
		Return(feeReductionTypes, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, editFeeReductionData{
			Case:              caseItem,
			PaymentID:         123,
			FeeReduction:      feeReduction,
			FeeReductionTypes: feeReductionTypes,
		}).
		Return(expectedError)

	r, _ := http.NewRequest(http.MethodGet, "/?id=123", nil)
	w := httptest.NewRecorder()

	err := EditFeeReduction(client, template.Func)(w, r)

	assert.Equal(t, expectedError, err)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestPostEditFeeReductionValidationError(t *testing.T) {
	caseItem := sirius.Case{ID: 4, CaseType: "lpa", UID: "700700"}

	feeReduction := sirius.Payment{
		ID:               123,
		PaymentEvidence:  "",
		FeeReductionType: "REMISSION",
		PaymentDate:      sirius.DateString("2022-07-23"),
		Source:           sirius.FeeReductionSource,
		Case:             &caseItem,
	}

	feeReductionTypes := []sirius.RefDataItem{
		{
			Handle: "REMISSION",
			Label:  "Remission",
		},
	}

	validationError := sirius.ValidationError{
		Field: sirius.FieldErrors{
			"paymentEvidence": {"reason": "Value is required and cannot be empty"},
		},
	}

	client := &mockEditFeeReductionClient{}
	client.
		On("PaymentByID", mock.Anything, 123).
		Return(feeReduction, nil).
		On("Case", mock.Anything, 4).
		Return(caseItem, nil).
		On("RefDataByCategory", mock.Anything, sirius.FeeReductionTypeCategory).
		Return(feeReductionTypes, nil).
		On("EditPayment", mock.Anything, 123, feeReduction).
		Return(validationError)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, editFeeReductionData{
			Case:              caseItem,
			PaymentID:         123,
			FeeReduction:      feeReduction,
			FeeReductionTypes: feeReductionTypes,
			Error:             validationError,
		}).
		Return(nil)

	form := url.Values{
		"id":               {"123"},
		"source":           {sirius.FeeReductionSource},
		"paymentEvidence":  {""},
		"paymentDate":      {"2022-07-23"},
		"feeReductionType": {"REMISSION"},
	}

	r, _ := http.NewRequest(http.MethodPost, "/?id=123", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := EditFeeReduction(client, template.Func)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestPostEditFeeReduction(t *testing.T) {
	caseItem := sirius.Case{ID: 4, CaseType: "lpa", UID: "700700"}

	feeReduction := sirius.Payment{
		ID:               123,
		PaymentEvidence:  "Test evidence",
		FeeReductionType: "REMISSION",
		Source:           sirius.FeeReductionSource,
		PaymentDate:      sirius.DateString("2022-07-23"),
		Case:             &caseItem,
	}

	editedFeeReduction := sirius.Payment{
		ID:               123,
		PaymentEvidence:  "Edited evidence",
		FeeReductionType: "REMISSION",
		PaymentDate:      sirius.DateString("2022-07-23"),
		Source:           sirius.FeeReductionSource,
		Case:             &caseItem,
	}

	feeReductionTypes := []sirius.RefDataItem{
		{
			Handle: "REMISSION",
			Label:  "Remission",
		},
	}

	client := &mockEditFeeReductionClient{}
	client.
		On("PaymentByID", mock.Anything, 123).
		Return(feeReduction, nil).
		On("Case", mock.Anything, 4).
		Return(caseItem, nil).
		On("RefDataByCategory", mock.Anything, sirius.FeeReductionTypeCategory).
		Return(feeReductionTypes, nil).
		On("EditPayment", mock.Anything, 123, editedFeeReduction).
		Return(nil)

	template := &mockTemplate{}

	form := url.Values{
		"id":               {"123"},
		"source":           {sirius.FeeReductionSource},
		"paymentEvidence":  {"Edited evidence"},
		"paymentDate":      {"2022-07-23"},
		"feeReductionType": {"REMISSION"},
	}

	r, _ := http.NewRequest(http.MethodPost, "/?id=123", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := EditFeeReduction(client, template.Func)(w, r)
	resp := w.Result()

	assert.Equal(t, RedirectError("/payments/4"), err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}
