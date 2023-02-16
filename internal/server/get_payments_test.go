package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockGetPayments struct {
	mock.Mock
}

func (m *mockGetPayments) Payments(ctx sirius.Context, id int) ([]sirius.Payment, error) {
	args := m.Called(ctx, id)
	return args.Get(0).([]sirius.Payment), args.Error(1)
}

func (m *mockGetPayments) Case(ctx sirius.Context, id int) (sirius.Case, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(sirius.Case), args.Error(1)
}

func (m *mockGetPayments) RefDataByCategory(ctx sirius.Context, category string) ([]sirius.RefDataItem, error) {
	args := m.Called(ctx, category)
	if args.Get(0) != nil {
		return args.Get(0).([]sirius.RefDataItem), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockGetPayments) GetUserDetails(ctx sirius.Context) (sirius.User, error) {
	args := m.Called(ctx)
	return args.Get(0).(sirius.User), args.Error(1)
}

func TestGetPayments(t *testing.T) {
	allPayments := []sirius.Payment{
		{
			ID:     2,
			Amount: 4100,
		},
		{
			ID:     3,
			Amount: 1438,
		},
	}

	nonReductionPayments := []sirius.Payment{
		{
			ID:     2,
			Amount: 4100,
		},
		{
			ID:     3,
			Amount: 1438,
		},
	}

	caseItem := sirius.Case{
		UID:     "7000-0000-0021",
		SubType: "pfa",
	}

	paymentSources := []sirius.RefDataItem{
		{
			Handle: "PHONE",
			Label:  "Paid over the phone",
		},
	}

	feeReductionTypes := []sirius.RefDataItem{
		{
			Handle: "REMISSION",
			Label:  "Remission",
		},
	}

	user := sirius.User{ID: 1, DisplayName: "Test User", Roles: []string{"OPG User", "Reduced Fees User"}}

	referenceTypes := []sirius.RefDataItem{
		{
			Handle: "GOVUK",
			Label:  "GOV.UK Pay",
		},
	}

	client := &mockGetPayments{}
	client.
		On("Payments", mock.Anything, 4).
		Return(allPayments, nil)
	client.
		On("Case", mock.Anything, 4).
		Return(caseItem, nil)
	client.
		On("RefDataByCategory", mock.Anything, sirius.PaymentSourceCategory).
		Return(paymentSources, nil)
	client.
		On("RefDataByCategory", mock.Anything, sirius.PaymentReferenceType).
		Return(referenceTypes, nil)
	client.
		On("RefDataByCategory", mock.Anything, sirius.FeeReductionTypeCategory).
		Return(feeReductionTypes, nil)
	client.
		On("GetUserDetails", mock.Anything).
		Return(user, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, getPaymentsData{
			PaymentSources:    paymentSources,
			ReferenceTypes:    referenceTypes,
			Payments:          nonReductionPayments,
			FeeReductionTypes: feeReductionTypes,
			Case:              caseItem,
			TotalPaid:         5538,
			IsReducedFeesUser: true,
			OutstandingFee:    2662,
		}).
		Return(nil)

	r, _ := http.NewRequest(http.MethodGet, "/?id=4", nil)
	w := httptest.NewRecorder()

	err := GetPayments(client, template.Func)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestGetPaymentsNoID(t *testing.T) {
	testCases := map[string]string{
		"no-id":  "/",
		"bad-id": "/?id=test",
	}

	for name, testUrl := range testCases {
		t.Run(name, func(t *testing.T) {
			r, _ := http.NewRequest(http.MethodGet, testUrl, nil)
			w := httptest.NewRecorder()

			err := GetPayments(nil, nil)(w, r)

			assert.NotNil(t, err)
		})
	}
}

func TestGetPaymentsWhenFailureOnGetCase(t *testing.T) {
	client := &mockGetPayments{}
	client.
		On("Case", mock.Anything, 4).
		Return(sirius.Case{}, expectedError)

	r, _ := http.NewRequest(http.MethodGet, "/?id=4", nil)
	w := httptest.NewRecorder()

	err := GetPayments(client, nil)(w, r)

	assert.Equal(t, expectedError, err)
	mock.AssertExpectationsForObjects(t, client)
}

func TestGetPaymentsWhenFailureOnGetPayments(t *testing.T) {
	caseItem := sirius.Case{
		UID:     "7000-0000-0021",
		SubType: "pfa",
	}

	client := &mockGetPayments{}
	client.
		On("Case", mock.Anything, 4).
		Return(caseItem, nil)

	client.
		On("Payments", mock.Anything, 4).
		Return([]sirius.Payment{}, expectedError)

	r, _ := http.NewRequest(http.MethodGet, "/?id=4", nil)
	w := httptest.NewRecorder()

	err := GetPayments(client, nil)(w, r)

	assert.Equal(t, expectedError, err)
	mock.AssertExpectationsForObjects(t, client)
}

func TestGetPaymentsWhenFailureOnGetPaymentSourceRefData(t *testing.T) {
	payments := []sirius.Payment{
		{
			ID:          2,
			Source:      "PHONE",
			Amount:      4100,
			PaymentDate: sirius.DateString("2022-08-23T14:55:20+00:00"),
		},
	}

	caseItem := sirius.Case{
		UID:     "7000-0000-0021",
		SubType: "pfa",
	}

	client := &mockGetPayments{}
	client.
		On("Case", mock.Anything, 4).
		Return(caseItem, nil)
	client.
		On("Payments", mock.Anything, 4).
		Return(payments, nil)
	client.
		On("RefDataByCategory", mock.Anything, sirius.PaymentSourceCategory).
		Return([]sirius.RefDataItem{}, expectedError)

	r, _ := http.NewRequest(http.MethodGet, "/?id=4", nil)
	w := httptest.NewRecorder()

	err := GetPayments(client, nil)(w, r)

	assert.Equal(t, expectedError, err)
	mock.AssertExpectationsForObjects(t, client)
}

func TestGetPaymentsWhenFailureOnGetReferenceTypeRefData(t *testing.T) {
	payments := []sirius.Payment{
		{
			ID:          2,
			Source:      "PHONE",
			Amount:      4100,
			PaymentDate: sirius.DateString("2022-08-23T14:55:20+00:00"),
		},
	}

	caseItem := sirius.Case{
		UID:     "7000-0000-0021",
		SubType: "pfa",
	}

	paymentSources := []sirius.RefDataItem{
		{
			Handle: "PHONE",
			Label:  "Paid over the phone",
		},
	}

	client := &mockGetPayments{}
	client.
		On("Case", mock.Anything, 4).
		Return(caseItem, nil)
	client.
		On("Payments", mock.Anything, 4).
		Return(payments, nil)
	client.
		On("RefDataByCategory", mock.Anything, sirius.PaymentSourceCategory).
		Return(paymentSources, nil).
		On("RefDataByCategory", mock.Anything, sirius.PaymentReferenceType).
		Return([]sirius.RefDataItem{}, expectedError)

	r, _ := http.NewRequest(http.MethodGet, "/?id=4", nil)
	w := httptest.NewRecorder()

	err := GetPayments(client, nil)(w, r)

	assert.Equal(t, expectedError, err)
	mock.AssertExpectationsForObjects(t, client)
}

func TestGetPaymentsWhenFailureOnFeeReductionTypesRefData(t *testing.T) {
	payments := []sirius.Payment{
		{
			ID:          2,
			Source:      "PHONE",
			Amount:      4100,
			PaymentDate: sirius.DateString("2022-08-23T14:55:20+00:00"),
		},
	}

	caseItem := sirius.Case{
		UID:     "7000-0000-0021",
		SubType: "pfa",
	}

	paymentSources := []sirius.RefDataItem{
		{
			Handle: "PHONE",
			Label:  "Paid over the phone",
		},
	}

	referenceTypes := []sirius.RefDataItem{
		{
			Handle: "GOVUK",
			Label:  "GOV.UK Pay",
		},
	}

	client := &mockGetPayments{}
	client.
		On("Case", mock.Anything, 4).
		Return(caseItem, nil)
	client.
		On("Payments", mock.Anything, 4).
		Return(payments, nil)
	client.
		On("RefDataByCategory", mock.Anything, sirius.PaymentSourceCategory).
		Return(paymentSources, nil).
		On("RefDataByCategory", mock.Anything, sirius.PaymentReferenceType).
		Return(referenceTypes, nil).
		On("RefDataByCategory", mock.Anything, sirius.FeeReductionTypeCategory).
		Return([]sirius.RefDataItem{}, expectedError)

	r, _ := http.NewRequest(http.MethodGet, "/?id=4", nil)
	w := httptest.NewRecorder()

	err := GetPayments(client, nil)(w, r)

	assert.Equal(t, expectedError, err)
	mock.AssertExpectationsForObjects(t, client)
}

func TestGetPaymentsWhenTemplateErrors(t *testing.T) {
	payments := []sirius.Payment{
		{
			ID:          2,
			Source:      "PHONE",
			Amount:      4100,
			PaymentDate: sirius.DateString("2022-08-23T14:55:20+00:00"),
		},
	}

	caseItem := sirius.Case{
		UID:     "7000-0000-0021",
		SubType: "pfa",
	}

	paymentSources := []sirius.RefDataItem{
		{
			Handle: "PHONE",
			Label:  "Paid over the phone",
		},
	}

	feeReductionTypes := []sirius.RefDataItem{
		{
			Handle: "REMISSION",
			Label:  "Remission",
		},
	}

	user := sirius.User{ID: 1, DisplayName: "Test User", Roles: []string{"OPG User", "Case Manager"}}

	referenceTypes := []sirius.RefDataItem{
		{
			Handle: "GOVUK",
			Label:  "GOV.UK Pay",
		},
	}

	client := &mockGetPayments{}
	client.
		On("Payments", mock.Anything, 4).
		Return(payments, nil)
	client.
		On("Case", mock.Anything, 4).
		Return(caseItem, nil)
	client.
		On("RefDataByCategory", mock.Anything, sirius.PaymentSourceCategory).
		Return(paymentSources, nil)

	client.
		On("RefDataByCategory", mock.Anything, sirius.PaymentReferenceType).
		Return(referenceTypes, nil)

	client.
		On("RefDataByCategory", mock.Anything, sirius.FeeReductionTypeCategory).
		Return(feeReductionTypes, nil)
	client.
		On("GetUserDetails", mock.Anything).
		Return(user, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, getPaymentsData{
			Payments:          payments,
			PaymentSources:    paymentSources,
			ReferenceTypes:    referenceTypes,
			Case:              caseItem,
			TotalPaid:         4100,
			IsReducedFeesUser: false,
			FeeReductionTypes: feeReductionTypes,
			OutstandingFee:    4100,
		}).
		Return(expectedError)

	r, _ := http.NewRequest(http.MethodGet, "/?id=4", nil)
	w := httptest.NewRecorder()

	err := GetPayments(client, template.Func)(w, r)

	assert.Equal(t, expectedError, err)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestGetPaymentWhenRefundDue(t *testing.T) {
	allPayments := []sirius.Payment{
		{
			ID:     2,
			Amount: 5000,
		},
		{
			ID:               4,
			Source:           sirius.FeeReductionSource,
			FeeReductionType: "REMISSION",
			PaymentEvidence:  "Test",
			PaymentDate:      "2022-04-05",
			Amount:           4100,
		},
	}

	nonReductionPayments := []sirius.Payment{
		{
			ID:     2,
			Amount: 5000,
		},
	}

	feeReductions := []sirius.Payment{
		{
			ID:               4,
			Source:           sirius.FeeReductionSource,
			FeeReductionType: "REMISSION",
			PaymentEvidence:  "Test",
			PaymentDate:      "2022-04-05",
			Amount:           4100,
		},
	}

	caseItem := sirius.Case{
		UID:     "7000-0000-0021",
		SubType: "pfa",
	}

	paymentSources := []sirius.RefDataItem{
		{
			Handle: "PHONE",
			Label:  "Paid over the phone",
		},
	}

	feeReductionTypes := []sirius.RefDataItem{
		{
			Handle: "REMISSION",
			Label:  "Remission",
		},
	}

	user := sirius.User{ID: 1, DisplayName: "Test User", Roles: []string{"OPG User", "Reduced Fees User"}}

	referenceTypes := []sirius.RefDataItem{
		{
			Handle: "GOVUK",
			Label:  "GOV.UK Pay",
		},
	}

	client := &mockGetPayments{}
	client.
		On("Payments", mock.Anything, 4).
		Return(allPayments, nil)
	client.
		On("Case", mock.Anything, 4).
		Return(caseItem, nil)
	client.
		On("RefDataByCategory", mock.Anything, sirius.PaymentSourceCategory).
		Return(paymentSources, nil)
	client.
		On("RefDataByCategory", mock.Anything, sirius.PaymentReferenceType).
		Return(referenceTypes, nil)
	client.
		On("RefDataByCategory", mock.Anything, sirius.FeeReductionTypeCategory).
		Return(feeReductionTypes, nil)
	client.
		On("GetUserDetails", mock.Anything).
		Return(user, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, getPaymentsData{
			PaymentSources:    paymentSources,
			ReferenceTypes:    referenceTypes,
			Payments:          nonReductionPayments,
			FeeReductions:     feeReductions,
			FeeReductionTypes: feeReductionTypes,
			Case:              caseItem,
			TotalPaid:         5000,
			IsReducedFeesUser: true,
			RefundAmount:      900,
		}).
		Return(nil)

	r, _ := http.NewRequest(http.MethodGet, "/?id=4", nil)
	w := httptest.NewRecorder()

	err := GetPayments(client, template.Func)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestGetPaymentsCalculations(t *testing.T) {
	testCases := []struct {
		allPayments    []sirius.Payment
		payments       []sirius.Payment
		feeReductions  []sirius.Payment
		refunds        []sirius.Payment
		totalPaid      int
		totalRefunds   int
		outstandingFee int
		refundAmount   int
	}{
		{
			[]sirius.Payment{{ID: 2, Amount: 4100}, {ID: 3, Amount: 1500}, {ID: 4, Amount: -4100}},
			[]sirius.Payment{{ID: 2, Amount: 4100}, {ID: 3, Amount: 1500}},
			[]sirius.Payment(nil),
			[]sirius.Payment{{ID: 4, Amount: -4100}},
			5600,
			4100,
			6700,
			0,
		}, // underpaid, outstanding fee due
		{
			[]sirius.Payment{
				{ID: 2, Amount: 4100},
				{ID: 7, Amount: 4100},
				{ID: 3, Amount: -2500},
				{ID: 4, Amount: -1600},
				{ID: 8, Amount: 8200, Source: "FEE_REDUCTION", FeeReductionType: "EXEMPTION"},
			},
			[]sirius.Payment{{ID: 2, Amount: 4100}, {ID: 7, Amount: 4100}},
			[]sirius.Payment{{ID: 8, Amount: 8200, Source: "FEE_REDUCTION", FeeReductionType: "EXEMPTION"}},
			[]sirius.Payment{{ID: 3, Amount: -2500}, {ID: 4, Amount: -1600}},
			8200,
			4100,
			0,
			4100,
		}, // overpaid, refund due
		{
			[]sirius.Payment{
				{ID: 2, Amount: 4100},
				{ID: 7, Amount: 4100},
				{ID: 8, Amount: 4100, Source: "FEE_REDUCTION", FeeReductionType: "REMISSION"},
				{ID: 3, Amount: -4100},
			},
			[]sirius.Payment{{ID: 2, Amount: 4100}, {ID: 7, Amount: 4100}},
			[]sirius.Payment{{ID: 8, Amount: 4100, Source: "FEE_REDUCTION", FeeReductionType: "REMISSION"}},
			[]sirius.Payment{{ID: 3, Amount: -4100}},
			8200,
			4100,
			0,
			0,
		}, // paid in full
	}
	for _, tc := range testCases {
		caseItem := sirius.Case{
			UID:     "7000-0000-0021",
			SubType: "pfa",
		}

		paymentSources := []sirius.RefDataItem{
			{
				Handle: "PHONE",
				Label:  "Paid over the phone",
			},
		}

		feeReductionTypes := []sirius.RefDataItem{
			{
				Handle: "REMISSION",
				Label:  "Remission",
			},
		}

		user := sirius.User{ID: 1, DisplayName: "Test User", Roles: []string{"OPG User", "Reduced Fees User"}}

		referenceTypes := []sirius.RefDataItem{
			{
				Handle: "GOVUK",
				Label:  "GOV.UK Pay",
			},
		}

		client := &mockGetPayments{}
		client.
			On("Payments", mock.Anything, 4).
			Return(tc.allPayments, nil)
		client.
			On("Case", mock.Anything, 4).
			Return(caseItem, nil)
		client.
			On("RefDataByCategory", mock.Anything, sirius.PaymentSourceCategory).
			Return(paymentSources, nil)
		client.
			On("RefDataByCategory", mock.Anything, sirius.PaymentReferenceType).
			Return(referenceTypes, nil)
		client.
			On("RefDataByCategory", mock.Anything, sirius.FeeReductionTypeCategory).
			Return(feeReductionTypes, nil)
		client.
			On("GetUserDetails", mock.Anything).
			Return(user, nil)

		template := &mockTemplate{}
		template.
			On("Func", mock.Anything, getPaymentsData{
				PaymentSources:    paymentSources,
				ReferenceTypes:    referenceTypes,
				FeeReductionTypes: feeReductionTypes,
				Payments:          tc.payments,
				FeeReductions:     tc.feeReductions,
				Refunds:           tc.refunds,
				Case:              caseItem,
				TotalPaid:         tc.totalPaid,
				TotalRefunds:      tc.totalRefunds,
				OutstandingFee:    tc.outstandingFee,
				RefundAmount:      tc.refundAmount,
				IsReducedFeesUser: true,
			}).
			Return(nil)

		r, _ := http.NewRequest(http.MethodGet, "/?id=4", nil)
		w := httptest.NewRecorder()

		err := GetPayments(client, template.Func)(w, r)
		resp := w.Result()

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		mock.AssertExpectationsForObjects(t, client, template)
	}
}
