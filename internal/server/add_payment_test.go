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

func (m *mockAddPaymentClient) AddPayment(ctx sirius.Context, caseID int, amount int, source string, paymentDate sirius.DateString, feeReductionType string, paymentEvidence string, appliedDate sirius.DateString) error {
	return m.Called(ctx, caseID, amount, source, paymentDate, feeReductionType, paymentEvidence, appliedDate).Error(0)
}

func (m *mockAddPaymentClient) Case(ctx sirius.Context, id int) (sirius.Case, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(sirius.Case), args.Error(1)
}

func (m *mockAddPaymentClient) RefDataByCategory(ctx sirius.Context, category string) ([]sirius.RefDataItem, error) {
	args := m.Called(ctx, category)
	if args.Get(0) != nil {
		return args.Get(0).([]sirius.RefDataItem), args.Error(1)
	}
	return nil, args.Error(1)
}

func TestGetAddPayment(t *testing.T) {
	caseItem := sirius.Case{
		UID:     "7000-0000-0021",
		SubType: "pfa",
	}

	paymentSources := []sirius.RefDataItem{
		{
			Handle:         "PHONE",
			Label:          "Paid over the phone",
			UserSelectable: true,
		},
	}

	feeReductionTypes := []sirius.RefDataItem{
		{
			Handle: "REMISSION",
			Label:  "Remission",
		},
	}

	client := &mockAddPaymentClient{}
	client.
		On("Case", mock.Anything, 4).
		Return(caseItem, nil)
	client.
		On("RefDataByCategory", mock.Anything, sirius.PaymentSourceCategory).
		Return(paymentSources, nil)
	client.
		On("RefDataByCategory", mock.Anything, sirius.FeeReductionTypeCategory).
		Return(feeReductionTypes, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, addPaymentData{
			Case:              caseItem,
			PaymentSources:    paymentSources,
			FeeReductionTypes: feeReductionTypes,
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

	paymentSources := []sirius.RefDataItem{
		{
			Handle:         "PHONE",
			Label:          "Paid over the phone",
			UserSelectable: true,
		},
	}

	feeReductionTypes := []sirius.RefDataItem{
		{
			Handle: "REMISSION",
			Label:  "Remission",
		},
	}

	client := &mockAddPaymentClient{}
	client.
		On("Case", mock.Anything, 123).
		Return(caseItem, nil)
	client.
		On("RefDataByCategory", mock.Anything, sirius.PaymentSourceCategory).
		Return(paymentSources, nil)
	client.
		On("RefDataByCategory", mock.Anything, sirius.FeeReductionTypeCategory).
		Return(feeReductionTypes, nil)

	expectedError := errors.New("err")

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, addPaymentData{
			Case:              caseItem,
			PaymentSources:    paymentSources,
			FeeReductionTypes: feeReductionTypes,
		}).
		Return(expectedError)

	r, _ := http.NewRequest(http.MethodGet, "/?id=123", nil)
	w := httptest.NewRecorder()

	err := AddPayment(client, template.Func)(w, r)

	assert.Equal(t, expectedError, err)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestAddPaymentWhenFailureOnGetPaymentSourceRefData(t *testing.T) {
	caseItem := sirius.Case{
		UID:     "7000-0000-0021",
		SubType: "pfa",
	}

	expectedError := errors.New("err")

	client := &mockAddPaymentClient{}
	client.
		On("Case", mock.Anything, 123).
		Return(caseItem, nil)
	client.
		On("RefDataByCategory", mock.Anything, sirius.PaymentSourceCategory).
		Return([]sirius.RefDataItem{}, expectedError)

	r, _ := http.NewRequest(http.MethodGet, "/?id=123", nil)
	w := httptest.NewRecorder()

	err := AddPayment(client, nil)(w, r)

	assert.Equal(t, expectedError, err)
	mock.AssertExpectationsForObjects(t, client)
}

func TestAddPaymentWhenFailureOnGetFeeReductionTypesRefData(t *testing.T) {
	caseItem := sirius.Case{
		UID:     "7000-0000-0021",
		SubType: "pfa",
	}

	paymentSources := []sirius.RefDataItem{
		{
			Handle:         "PHONE",
			Label:          "Paid over the phone",
			UserSelectable: true,
		},
	}

	expectedError := errors.New("err")

	client := &mockAddPaymentClient{}
	client.
		On("Case", mock.Anything, 123).
		Return(caseItem, nil)
	client.
		On("RefDataByCategory", mock.Anything, sirius.PaymentSourceCategory).
		Return(paymentSources, nil)
	client.
		On("RefDataByCategory", mock.Anything, sirius.FeeReductionTypeCategory).
		Return([]sirius.RefDataItem{}, expectedError)

	r, _ := http.NewRequest(http.MethodGet, "/?id=123", nil)
	w := httptest.NewRecorder()

	err := AddPayment(client, nil)(w, r)

	assert.Equal(t, expectedError, err)
	mock.AssertExpectationsForObjects(t, client)
}

func TestPostAddPayment(t *testing.T) {
	caseitem := sirius.Case{CaseType: "lpa", UID: "700700"}

	paymentSources := []sirius.RefDataItem{
		{
			Handle:         "PHONE",
			Label:          "Paid over the phone",
			UserSelectable: true,
		},
	}

	feeReductionTypes := []sirius.RefDataItem{
		{
			Handle: "REMISSION",
			Label:  "Remission",
		},
	}

	client := &mockAddPaymentClient{}
	client.
		On("AddPayment", mock.Anything, 123, 4100, "MAKE", sirius.DateString("2022-01-23"), "", "", sirius.DateString("")).
		Return(nil)
	client.
		On("Case", mock.Anything, 123).
		Return(caseitem, nil)
	client.
		On("RefDataByCategory", mock.Anything, sirius.PaymentSourceCategory).
		Return(paymentSources, nil)
	client.
		On("RefDataByCategory", mock.Anything, sirius.FeeReductionTypeCategory).
		Return(feeReductionTypes, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, addPaymentData{
			Success:           true,
			Case:              caseitem,
			Amount:            "41.00",
			Source:            "MAKE",
			PaymentDate:       sirius.DateString("2022-01-23"),
			PaymentSources:    paymentSources,
			FeeReductionTypes: feeReductionTypes,
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

func TestPostAddPaymentAmountIncorrectFormat(t *testing.T) {
	for _, amount := range []string{"41", "41.5", "41.555", ".45"} {
		t.Run(amount, func(t *testing.T) {
			caseitem := sirius.Case{CaseType: "lpa", UID: "700700"}

			paymentSources := []sirius.RefDataItem{
				{
					Handle:         "PHONE",
					Label:          "Paid over the phone",
					UserSelectable: true,
				},
			}

			feeReductionTypes := []sirius.RefDataItem{
				{
					Handle: "REMISSION",
					Label:  "Remission",
				},
			}

			client := &mockAddPaymentClient{}
			client.
				On("Case", mock.Anything, 123).
				Return(caseitem, nil)
			client.
				On("RefDataByCategory", mock.Anything, sirius.PaymentSourceCategory).
				Return(paymentSources, nil)
			client.
				On("RefDataByCategory", mock.Anything, sirius.FeeReductionTypeCategory).
				Return(feeReductionTypes, nil)

			validationError := sirius.ValidationError{
				Field: sirius.FieldErrors{
					"amount": {"reason": "Please enter the amount to 2 decimal places"},
				},
			}

			template := &mockTemplate{}
			template.
				On("Func", mock.Anything, addPaymentData{
					Success:           false,
					Case:              caseitem,
					Amount:            amount,
					Source:            "MAKE",
					PaymentDate:       sirius.DateString("2022-01-23"),
					Error:             validationError,
					PaymentSources:    paymentSources,
					FeeReductionTypes: feeReductionTypes,
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

			err := AddPayment(client, template.Func)(w, r)
			resp := w.Result()

			assert.Nil(t, err)
			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
			mock.AssertExpectationsForObjects(t, client, template)
		})
	}
}

func TestPostFeeReduction(t *testing.T) {
	caseitem := sirius.Case{CaseType: "lpa", UID: "700700"}

	paymentSources := []sirius.RefDataItem{
		{
			Handle:         "PHONE",
			Label:          "Paid over the phone",
			UserSelectable: true,
		},
	}

	feeReductionTypes := []sirius.RefDataItem{
		{
			Handle: "REMISSION",
			Label:  "Remission",
		},
	}

	client := &mockAddPaymentClient{}
	client.
		On("AddPayment", mock.Anything, 123, 0, "FEE_REDUCTION", sirius.DateString(""), "REMISSION", "Test evidence", sirius.DateString("2022-01-23")).
		Return(nil)
	client.
		On("Case", mock.Anything, 123).
		Return(caseitem, nil)
	client.
		On("RefDataByCategory", mock.Anything, sirius.PaymentSourceCategory).
		Return(paymentSources, nil)
	client.
		On("RefDataByCategory", mock.Anything, sirius.FeeReductionTypeCategory).
		Return(feeReductionTypes, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, addPaymentData{
			Success:           true,
			Case:              caseitem,
			FeeReductionType:  "REMISSION",
			Source:            "FEE_REDUCTION",
			AppliedDate:       "2022-01-23",
			PaymentEvidence:   "Test evidence",
			PaymentSources:    paymentSources,
			FeeReductionTypes: feeReductionTypes,
		}).
		Return(nil)

	form := url.Values{
		"feeReductionType": {"REMISSION"},
		"source":           {"FEE_REDUCTION"},
		"appliedDate":      {"2022-01-23"},
		"paymentEvidence":  {"Test evidence"},
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
