package server

import (
	"net/http"
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

func (m *mockGetPayments) CaseSummary(ctx sirius.Context, uid string) (sirius.CaseSummary, error) {
	args := m.Called(ctx, uid)
	return args.Get(0).(sirius.CaseSummary), args.Error(1)
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

	expectedPaymentTotal := 8000

	caseItem := sirius.Case{
		UID:                  "7000-0000-0021",
		SubType:              "pfa",
		ExpectedPaymentTotal: expectedPaymentTotal,
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
		On("Payments", mock.Anything, 901).
		Return(allPayments, nil).
		On("Case", mock.Anything, 901).
		Return(caseItem, nil).
		On("RefDataByCategory", mock.Anything, sirius.PaymentSourceCategory).
		Return(paymentSources, nil).
		On("RefDataByCategory", mock.Anything, sirius.PaymentReferenceType).
		Return(referenceTypes, nil).
		On("RefDataByCategory", mock.Anything, sirius.FeeReductionTypeCategory).
		Return(feeReductionTypes, nil).
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
			OutstandingFee:    expectedPaymentTotal - 5538,
		}).
		Return(nil)

	server := newMockServer("/payments/{id}", GetPayments(client, template.Func))

	req, _ := http.NewRequest(http.MethodGet, "/payments/901", nil)
	resp, err := server.serve(req)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.Code)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestGetPaymentsBadID(t *testing.T) {
	client := &mockGetPayments{}
	template := &mockTemplate{}

	server := newMockServer("/payments/{id}", GetPayments(client, template.Func))

	req, _ := http.NewRequest(http.MethodGet, "/payments/bad-id", nil)
	_, err := server.serve(req)

	assert.NotNil(t, err)
}

func TestGetPaymentsWhenFailureOnGetCase(t *testing.T) {
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
		On("Case", mock.Anything, 8).
		Return(sirius.Case{}, errExample).
		On("RefDataByCategory", mock.Anything, sirius.PaymentSourceCategory).
		Return(paymentSources, nil).
		On("RefDataByCategory", mock.Anything, sirius.PaymentReferenceType).
		Return(referenceTypes, nil).
		On("RefDataByCategory", mock.Anything, sirius.FeeReductionTypeCategory).
		Return(feeReductionTypes, nil).
		On("GetUserDetails", mock.Anything).
		Return(user, nil)

	client.AssertNotCalled(t, "Payments")

	server := newMockServer("/payments/{id}", GetPayments(client, nil))

	req, _ := http.NewRequest(http.MethodGet, "/payments/8", nil)
	_, err := server.serve(req)

	assert.Equal(t, errExample, err)
	mock.AssertExpectationsForObjects(t, client)
}

func TestGetPaymentsWhenFailureOnGetCaseSummary(t *testing.T) {
	client := &mockGetPayments{}
	client.
		On("CaseSummary", mock.Anything, "M-QQQQ-WWWW-EEEE").
		Return(sirius.CaseSummary{}, errExample)

	server := newMockServer("/lpa/{uid}/payments", GetPayments(client, nil))

	req, _ := http.NewRequest(http.MethodGet, "/lpa/M-QQQQ-WWWW-EEEE/payments", nil)
	_, err := server.serve(req)

	assert.Equal(t, errExample, err)
	mock.AssertExpectationsForObjects(t, client)
}

func TestGetPaymentsWhenFailureOnGetPayments(t *testing.T) {
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
		On("Case", mock.Anything, 9).
		Return(caseItem, nil).
		On("RefDataByCategory", mock.Anything, sirius.PaymentSourceCategory).
		Return(paymentSources, nil).
		On("RefDataByCategory", mock.Anything, sirius.PaymentReferenceType).
		Return(referenceTypes, nil).
		On("RefDataByCategory", mock.Anything, sirius.FeeReductionTypeCategory).
		Return(feeReductionTypes, nil).
		On("GetUserDetails", mock.Anything).
		Return(user, nil).
		On("Payments", mock.Anything, 9).
		Return([]sirius.Payment{}, errExample)

	server := newMockServer("/payments/{id}", GetPayments(client, nil))

	req, _ := http.NewRequest(http.MethodGet, "/payments/9", nil)
	_, err := server.serve(req)

	assert.Equal(t, errExample, err)
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
		On("Case", mock.Anything, 111).
		Return(caseItem, nil).
		On("Payments", mock.Anything, 111).
		Return(payments, nil).
		On("RefDataByCategory", mock.Anything, sirius.PaymentReferenceType).
		Return(referenceTypes, nil).
		On("GetUserDetails", mock.Anything).
		Return(user, nil).
		On("RefDataByCategory", mock.Anything, sirius.FeeReductionTypeCategory).
		Return(feeReductionTypes, nil).
		On("RefDataByCategory", mock.Anything, sirius.PaymentSourceCategory).
		Return([]sirius.RefDataItem{}, errExample)

	server := newMockServer("/payments/{id}", GetPayments(client, nil))

	req, _ := http.NewRequest(http.MethodGet, "/payments/111", nil)
	_, err := server.serve(req)

	assert.Equal(t, errExample, err)
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

	feeReductionTypes := []sirius.RefDataItem{
		{
			Handle: "REMISSION",
			Label:  "Remission",
		},
	}

	user := sirius.User{ID: 1, DisplayName: "Test User", Roles: []string{"OPG User", "Reduced Fees User"}}

	client := &mockGetPayments{}
	client.
		On("Case", mock.Anything, 4).
		Return(caseItem, nil).
		On("Payments", mock.Anything, 4).
		Return(payments, nil).
		On("RefDataByCategory", mock.Anything, sirius.PaymentSourceCategory).
		Return(paymentSources, nil).
		On("RefDataByCategory", mock.Anything, sirius.FeeReductionTypeCategory).
		Return(feeReductionTypes, nil).
		On("GetUserDetails", mock.Anything).
		Return(user, nil).
		On("RefDataByCategory", mock.Anything, sirius.PaymentReferenceType).
		Return([]sirius.RefDataItem{}, errExample)

	server := newMockServer("/payments/{id}", GetPayments(client, nil))

	req, _ := http.NewRequest(http.MethodGet, "/payments/4", nil)
	_, err := server.serve(req)

	assert.Equal(t, errExample, err)
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

	user := sirius.User{ID: 1, DisplayName: "Test User", Roles: []string{"OPG User", "Reduced Fees User"}}

	client := &mockGetPayments{}
	client.
		On("Case", mock.Anything, 876).
		Return(caseItem, nil).
		On("Payments", mock.Anything, 876).
		Return(payments, nil).
		On("RefDataByCategory", mock.Anything, sirius.PaymentSourceCategory).
		Return(paymentSources, nil).
		On("RefDataByCategory", mock.Anything, sirius.PaymentReferenceType).
		Return(referenceTypes, nil).
		On("GetUserDetails", mock.Anything).
		Return(user, nil).
		On("RefDataByCategory", mock.Anything, sirius.FeeReductionTypeCategory).
		Return([]sirius.RefDataItem{}, errExample)

	server := newMockServer("/payments/{id}", GetPayments(client, nil))

	req, _ := http.NewRequest(http.MethodGet, "/payments/876", nil)
	_, err := server.serve(req)

	assert.Equal(t, errExample, err)
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

	expectedPaymentTotal := 8000

	caseItem := sirius.Case{
		UID:                  "7000-0000-0021",
		SubType:              "pfa",
		ExpectedPaymentTotal: expectedPaymentTotal,
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
		On("Payments", mock.Anything, 554).
		Return(payments, nil).
		On("Case", mock.Anything, 554).
		Return(caseItem, nil).
		On("RefDataByCategory", mock.Anything, sirius.PaymentSourceCategory).
		Return(paymentSources, nil).
		On("RefDataByCategory", mock.Anything, sirius.PaymentReferenceType).
		Return(referenceTypes, nil).
		On("RefDataByCategory", mock.Anything, sirius.FeeReductionTypeCategory).
		Return(feeReductionTypes, nil).
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
			OutstandingFee:    expectedPaymentTotal - 4100,
		}).
		Return(errExample)

	server := newMockServer("/payments/{id}", GetPayments(client, template.Func))

	req, _ := http.NewRequest(http.MethodGet, "/payments/554", nil)
	_, err := server.serve(req)

	assert.Equal(t, errExample, err)
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

	expectedPaymentTotal := 8000

	caseItem := sirius.Case{
		ID:                   742,
		UID:                  "7000-0000-0021",
		SubType:              "pfa",
		ExpectedPaymentTotal: expectedPaymentTotal,
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

	// get by case ID
	client1 := &mockGetPayments{}
	client1.
		On("Payments", mock.Anything, 742).
		Return(allPayments, nil).
		On("Case", mock.Anything, 742).
		Return(caseItem, nil).
		On("RefDataByCategory", mock.Anything, sirius.PaymentSourceCategory).
		Return(paymentSources, nil).
		On("RefDataByCategory", mock.Anything, sirius.PaymentReferenceType).
		Return(referenceTypes, nil).
		On("RefDataByCategory", mock.Anything, sirius.FeeReductionTypeCategory).
		Return(feeReductionTypes, nil).
		On("GetUserDetails", mock.Anything).
		Return(user, nil)

	template1 := &mockTemplate{}
	template1.
		On("Func", mock.Anything, getPaymentsData{
			PaymentSources:    paymentSources,
			ReferenceTypes:    referenceTypes,
			Payments:          nonReductionPayments,
			FeeReductions:     feeReductions,
			FeeReductionTypes: feeReductionTypes,
			Case:              caseItem,
			TotalPaid:         5000,
			IsReducedFeesUser: true,
			RefundAmount:      (expectedPaymentTotal - 4100 - 5000) * -1,
		}).
		Return(nil)

	server := newMockServer("/payments/{id}", GetPayments(client1, template1.Func))
	req, _ := http.NewRequest(http.MethodGet, "/payments/742", nil)
	_, err := server.serve(req)

	assert.Nil(t, err)
	mock.AssertExpectationsForObjects(t, client1, template1)

	// get by digital LPA UID
	caseSummary := sirius.CaseSummary{
		DigitalLpa: sirius.DigitalLpa{
			SiriusData: sirius.SiriusData{
				ID: 742,
			},
		},
		TaskList: []sirius.Task{},
	}

	client2 := &mockGetPayments{}
	client2.
		On("Payments", mock.Anything, 742).
		Return(allPayments, nil).
		On("Case", mock.Anything, 742).
		Return(caseItem, nil).
		On("CaseSummary", mock.Anything, "M-QQQQ-EEEE-YYYY").
		Return(caseSummary, nil).
		On("RefDataByCategory", mock.Anything, sirius.PaymentSourceCategory).
		Return(paymentSources, nil).
		On("RefDataByCategory", mock.Anything, sirius.PaymentReferenceType).
		Return(referenceTypes, nil).
		On("RefDataByCategory", mock.Anything, sirius.FeeReductionTypeCategory).
		Return(feeReductionTypes, nil).
		On("GetUserDetails", mock.Anything).
		Return(user, nil)

	template2 := &mockTemplate{}
	template2.
		On("Func", mock.Anything, getPaymentsData{
			PaymentSources:    paymentSources,
			ReferenceTypes:    referenceTypes,
			Payments:          nonReductionPayments,
			FeeReductions:     feeReductions,
			FeeReductionTypes: feeReductionTypes,
			Case:              caseItem,
			CaseSummary:       caseSummary,
			TotalPaid:         5000,
			IsReducedFeesUser: true,
			RefundAmount:      (expectedPaymentTotal - 4100 - 5000) * -1,
		}).
		Return(nil)

	server = newMockServer("/lpa/{uid}/payments", GetPayments(client2, template2.Func))
	req, _ = http.NewRequest(http.MethodGet, "/lpa/M-QQQQ-EEEE-YYYY/payments", nil)
	_, err = server.serve(req)

	assert.Nil(t, err)
	mock.AssertExpectationsForObjects(t, client2, template2)
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
			UID:                  "7000-0000-0021",
			SubType:              "pfa",
			ExpectedPaymentTotal: 8200,
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
			On("Payments", mock.Anything, 929).
			Return(tc.allPayments, nil).
			On("Case", mock.Anything, 929).
			Return(caseItem, nil).
			On("RefDataByCategory", mock.Anything, sirius.PaymentSourceCategory).
			Return(paymentSources, nil).
			On("RefDataByCategory", mock.Anything, sirius.PaymentReferenceType).
			Return(referenceTypes, nil).
			On("RefDataByCategory", mock.Anything, sirius.FeeReductionTypeCategory).
			Return(feeReductionTypes, nil).
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

		server := newMockServer("/payments/{id}", GetPayments(client, template.Func))

		req, _ := http.NewRequest(http.MethodGet, "/payments/929", nil)
		resp, err := server.serve(req)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, resp.Code)
		mock.AssertExpectationsForObjects(t, client, template)
	}
}
