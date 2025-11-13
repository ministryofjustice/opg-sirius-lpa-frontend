package sirius

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/pact-foundation/pact-go/v2/consumer"
	"github.com/pact-foundation/pact-go/v2/matchers"
	"github.com/stretchr/testify/assert"
)

func TestPayment(t *testing.T) {
	t.Parallel()

	pact, err := newPact()
	assert.NoError(t, err)

	testCases := []struct {
		name             string
		setup            func()
		expectedError    func(int) error
		expectedResponse []Payment
	}{
		{
			name: "OK",
			setup: func() {
				pact.
					AddInteraction().
					Given("I have an lpa which has been paid for").
					UponReceiving("A request for the payments by case").
					WithCompleteRequest(consumer.Request{
						Method: http.MethodGet,
						Path:   matchers.String("/lpa-api/v1/cases/800/payments"),
					}).
					WithCompleteResponse(consumer.Response{
						Status: http.StatusOK,
						Body: matchers.EachLike(map[string]interface{}{
							"id":          matchers.Like(2),
							"source":      matchers.Like("MAKE"),
							"amount":      matchers.Like(4100),
							"paymentDate": matchers.String("23/01/2022"),
							"case": matchers.Like(map[string]interface{}{
								"id": matchers.Like(800),
							}),
						}, 1),
						Headers: matchers.MapMatcher{"Content-Type": matchers.String("application/json")},
					})
			},
			expectedResponse: []Payment{
				{
					ID:          2,
					Source:      "MAKE",
					Amount:      4100,
					PaymentDate: DateString("2022-01-23"),
					Case:        &Case{ID: 800},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.ExecuteTest(t, func(config consumer.MockServerConfig) error {
				client := NewClient(http.DefaultClient, fmt.Sprintf("http://127.0.0.1:%d", config.Port))

				payments, err := client.Payments(Context{Context: context.Background()}, 800)

				assert.Equal(t, tc.expectedResponse, payments)
				if tc.expectedError == nil {
					assert.Nil(t, err)
				} else {
					assert.Equal(t, tc.expectedError(config.Port), err)
				}
				return nil
			}))
		})
	}
}

func TestFeeReductionOnCase(t *testing.T) {
	t.Parallel()

	pact, err := newPact()
	assert.NoError(t, err)

	testCases := []struct {
		name             string
		setup            func()
		expectedError    func(int) error
		expectedResponse []Payment
	}{
		{
			name: "OK",
			setup: func() {
				pact.
					AddInteraction().
					Given("I have an lpa which has a fee reduction").
					UponReceiving("A request for the fee reduction by case").
					WithCompleteRequest(consumer.Request{
						Method: http.MethodGet,
						Path:   matchers.String("/lpa-api/v1/cases/802/payments"),
					}).
					WithCompleteResponse(consumer.Response{
						Status: http.StatusOK,
						Body: matchers.EachLike(map[string]interface{}{
							"id":               matchers.Like(3),
							"source":           matchers.Like(FeeReductionSource),
							"feeReductionType": matchers.Like("REMISSION"),
							"paymentEvidence":  matchers.Like("Test\nmultiple\nline evidence"),
							"paymentDate":      matchers.String("24/01/2022"),
							"case": matchers.Like(map[string]interface{}{
								"id": matchers.Like(802),
							}),
							"amount": matchers.Like(4100),
						}, 1),
						Headers: matchers.MapMatcher{"Content-Type": matchers.String("application/json")},
					})
			},
			expectedResponse: []Payment{
				{
					ID:               3,
					Source:           FeeReductionSource,
					FeeReductionType: "REMISSION",
					PaymentEvidence:  "Test\nmultiple\nline evidence",
					PaymentDate:      DateString("2022-01-24"),
					Case:             &Case{ID: 802},
					Amount:           4100,
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.ExecuteTest(t, func(config consumer.MockServerConfig) error {
				client := NewClient(http.DefaultClient, fmt.Sprintf("http://127.0.0.1:%d", config.Port))

				payments, err := client.Payments(Context{Context: context.Background()}, 802)

				assert.Equal(t, tc.expectedResponse, payments)
				if tc.expectedError == nil {
					assert.Nil(t, err)
				} else {
					assert.Equal(t, tc.expectedError(config.Port), err)
				}
				return nil
			}))
		})
	}
}

func TestNoPaymentOnCase(t *testing.T) {
	t.Parallel()

	pact, err := newPact()
	assert.NoError(t, err)

	testCases := []struct {
		name             string
		setup            func()
		expectedError    func(int) error
		expectedResponse []Payment
	}{
		{
			name: "OK",
			setup: func() {
				pact.
					AddInteraction().
					Given("I have a pending case with no payment assigned").
					UponReceiving("A request for the payments by case").
					WithCompleteRequest(consumer.Request{
						Method: http.MethodGet,
						Path:   matchers.String("/lpa-api/v1/cases/801/payments"),
					}).
					WithCompleteResponse(consumer.Response{
						Status:  http.StatusOK,
						Body:    matchers.Like([]Payment{}),
						Headers: matchers.MapMatcher{"Content-Type": matchers.String("application/json")},
					})
			},
			expectedResponse: []Payment{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.ExecuteTest(t, func(config consumer.MockServerConfig) error {
				client := NewClient(http.DefaultClient, fmt.Sprintf("http://127.0.0.1:%d", config.Port))

				payments, err := client.Payments(Context{Context: context.Background()}, 801)

				assert.Equal(t, tc.expectedResponse, payments)
				if tc.expectedError == nil {
					assert.Nil(t, err)
				} else {
					assert.Equal(t, tc.expectedError(config.Port), err)
				}
				return nil
			}))
		})
	}
}

func TestPaymentByID(t *testing.T) {
	t.Parallel()

	pact, err := newPact()
	assert.NoError(t, err)

	testCases := []struct {
		name             string
		setup            func()
		cookies          []*http.Cookie
		expectedError    func(int) error
		expectedResponse Payment
	}{
		{
			name: "OK",
			setup: func() {
				pact.
					AddInteraction().
					Given("I have an lpa which has been paid for").
					UponReceiving("A request for that payment by ID").
					WithCompleteRequest(consumer.Request{
						Method: http.MethodGet,
						Path:   matchers.String("/lpa-api/v1/payments/123"),
					}).
					WithCompleteResponse(consumer.Response{
						Status: http.StatusOK,
						Body: matchers.Like(map[string]interface{}{
							"id":          matchers.Like(123),
							"source":      matchers.Like("PHONE"),
							"amount":      matchers.Like(4100),
							"paymentDate": matchers.String("23/01/2022"),
							"case": matchers.Like(map[string]interface{}{
								"id": matchers.Like(800),
							}),
						}),
						Headers: matchers.MapMatcher{"Content-Type": matchers.String("application/json")},
					})
			},
			expectedResponse: Payment{
				ID:          123,
				Source:      "PHONE",
				Amount:      4100,
				PaymentDate: DateString("2022-01-23"),
				Case:        &Case{ID: 800},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.ExecuteTest(t, func(config consumer.MockServerConfig) error {
				client := NewClient(http.DefaultClient, fmt.Sprintf("http://127.0.0.1:%d", config.Port))

				payment, err := client.PaymentByID(Context{Context: context.Background()}, 123)

				assert.Equal(t, tc.expectedResponse, payment)
				if tc.expectedError == nil {
					assert.Nil(t, err)
				} else {
					assert.Equal(t, tc.expectedError(config.Port), err)
				}
				return nil
			}))
		})
	}
}

func TestFeeReductionByID(t *testing.T) {
	t.Parallel()

	pact, err := newPact()
	assert.NoError(t, err)

	testCases := []struct {
		name             string
		setup            func()
		cookies          []*http.Cookie
		expectedError    func(int) error
		expectedResponse Payment
	}{
		{
			name: "OK",
			setup: func() {
				pact.
					AddInteraction().
					Given("I have an lpa which has a fee reduction").
					UponReceiving("A request for that fee reduction by ID").
					WithCompleteRequest(consumer.Request{
						Method: http.MethodGet,
						Path:   matchers.String("/lpa-api/v1/payments/124"),
					}).
					WithCompleteResponse(consumer.Response{
						Status: http.StatusOK,
						Body: matchers.Like(map[string]interface{}{
							"id":               matchers.Like(124),
							"source":           matchers.Like("FEE_REDUCTION"),
							"paymentEvidence":  matchers.String("Test evidence"),
							"feeReductionType": matchers.String("REMISSION"),
							"paymentDate":      matchers.String("23/01/2022"),
							"case": matchers.Like(map[string]interface{}{
								"id": matchers.Like(802),
							}),
							"amount": matchers.Like(4100),
						}),
						Headers: matchers.MapMatcher{"Content-Type": matchers.String("application/json")},
					})
			},
			expectedResponse: Payment{
				ID:               124,
				Source:           "FEE_REDUCTION",
				PaymentEvidence:  "Test evidence",
				FeeReductionType: "REMISSION",
				PaymentDate:      DateString("2022-01-23"),
				Case:             &Case{ID: 802},
				Amount:           4100,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.ExecuteTest(t, func(config consumer.MockServerConfig) error {
				client := NewClient(http.DefaultClient, fmt.Sprintf("http://127.0.0.1:%d", config.Port))

				payment, err := client.PaymentByID(Context{Context: context.Background()}, 124)

				assert.Equal(t, tc.expectedResponse, payment)
				if tc.expectedError == nil {
					assert.Nil(t, err)
				} else {
					assert.Equal(t, tc.expectedError(config.Port), err)
				}
				return nil
			}))
		})
	}
}

func TestAddPayment(t *testing.T) {
	t.Parallel()

	pact, err := newPact()
	assert.NoError(t, err)

	testCases := []struct {
		name          string
		setup         func()
		expectedError func(int) error
	}{
		{
			name: "OK",
			setup: func() {
				pact.
					AddInteraction().
					Given("I have a pending case assigned").
					UponReceiving("A request to create a payment").
					WithCompleteRequest(consumer.Request{
						Method: http.MethodPost,
						Path:   matchers.String("/lpa-api/v1/cases/800/payments"),
						Headers: matchers.MapMatcher{
							"Content-Type": matchers.String("application/json"),
						},
						Body: map[string]interface{}{
							"amount":      4100,
							"source":      "PHONE",
							"paymentDate": "25/04/2022",
						},
					}).
					WithCompleteResponse(consumer.Response{
						Status:  http.StatusCreated,
						Headers: matchers.MapMatcher{"Content-Type": matchers.String("application/json")},
					})
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.ExecuteTest(t, func(config consumer.MockServerConfig) error {
				client := NewClient(http.DefaultClient, fmt.Sprintf("http://127.0.0.1:%d", config.Port))

				err := client.AddPayment(Context{Context: context.Background()}, 800, 4100, "PHONE", DateString("2022-04-25"))

				if tc.expectedError == nil {
					assert.Nil(t, err)
				} else {
					assert.Equal(t, tc.expectedError(config.Port), err)
				}
				return nil
			}))
		})
	}
}

func TestEditPayment(t *testing.T) {
	t.Parallel()

	pact, err := newPact()
	assert.NoError(t, err)

	testCases := []struct {
		name          string
		setup         func()
		expectedError func(int) error
	}{
		{
			name: "OK",
			setup: func() {
				pact.
					AddInteraction().
					Given("I have an lpa which has been paid for").
					UponReceiving("A request to edit a payment").
					WithCompleteRequest(consumer.Request{
						Method: http.MethodPut,
						Path:   matchers.String("/lpa-api/v1/payments/123"),
						Headers: matchers.MapMatcher{
							"Content-Type": matchers.String("application/json"),
						},
						Body: map[string]interface{}{
							"amount":      matchers.Like(2550),
							"source":      matchers.String("PHONE"),
							"paymentDate": matchers.Term("27/04/2022", `^\d{1,2}/\d{1,2}/\d{4}$`),
						},
					}).
					WithCompleteResponse(consumer.Response{
						Status:  http.StatusOK,
						Headers: matchers.MapMatcher{"Content-Type": matchers.String("application/json")},
					})
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.ExecuteTest(t, func(config consumer.MockServerConfig) error {
				client := NewClient(http.DefaultClient, fmt.Sprintf("http://127.0.0.1:%d", config.Port))

				err := client.EditPayment(Context{Context: context.Background()}, 123, Payment{Amount: 2550, Source: "PHONE", PaymentDate: DateString("2022-04-27")})

				if tc.expectedError == nil {
					assert.Nil(t, err)
				} else {
					assert.Equal(t, tc.expectedError(config.Port), err)
				}
				return nil
			}))
		})
	}
}

//func TestEditFeeReduction(t *testing.T) {
//	t.Parallel()
//
//	pact, err := newPact()
//	assert.NoError(t, err)
//
//	testCases := []struct {
//		name          string
//		setup         func()
//		expectedError func(int) error
//	}{
//		{
//			name: "OK",
//			setup: func() {
//				pact.
//					AddInteraction().
//					Given("I have an lpa which has a fee reduction").
//					UponReceiving("A request to edit a fee reduction").
//					WithCompleteRequest(consumer.Request{
//						Method: http.MethodPut,
//						Path:   matchers.String("/lpa-api/v1/payments/124/reduction"),
//						Headers: matchers.MapMatcher{
//							"Content-Type": matchers.String("application/json"),
//						},
//						Body: map[string]interface{}{
//							"id":     124,
//							"amount": 4100,
//							"case": map[string]interface{}{
//								"id":                   802,
//								"status":               "Draft",
//								"expectedPaymentTotal": 8000,
//							},
//							"paymentEvidence":  matchers.String("Edited test evidence"),
//							"feeReductionType": matchers.String("REMISSION"),
//							"paymentDate":      matchers.Term("28/04/2022", `^\d{1,2}/\d{1,2}/\d{4}$`),
//							"source":           matchers.String(FeeReductionSource),
//						},
//					}).
//					WithCompleteResponse(consumer.Response{
//						Status:  http.StatusOK,
//						Headers: matchers.MapMatcher{"Content-Type": matchers.String("application/json")},
//					})
//			},
//		},
//	}
//
//	for _, tc := range testCases {
//		t.Run(tc.name, func(t *testing.T) {
//			tc.setup()
//
//			assert.Nil(t, pact.ExecuteTest(t, func(config consumer.MockServerConfig) error {
//				client := NewClient(http.DefaultClient, fmt.Sprintf("http://127.0.0.1:%d", config.Port))
//
//				err := client.EditPayment(Context{Context: context.Background()},
//					124,
//					Payment{
//						ID:               124,
//						Amount:           4100,
//						Case:             &Case{ID: 802, Status: shared.CaseStatusTypeDraft, ExpectedPaymentTotal: 8000},
//						PaymentEvidence:  "Edited test evidence",
//						FeeReductionType: "REMISSION",
//						PaymentDate:      "2022-04-28",
//						Source:           FeeReductionSource,
//					},
//				)
//
//				if tc.expectedError == nil {
//					assert.Nil(t, err)
//				} else {
//					assert.Equal(t, tc.expectedError(config.Port), err)
//				}
//				return nil
//			}))
//		})
//	}
//}

func TestDeletePayment(t *testing.T) {
	t.Parallel()

	pact, err := newPact()
	assert.NoError(t, err)

	testCases := []struct {
		name          string
		setup         func()
		expectedError func(int) error
	}{
		{
			name: "OK",
			setup: func() {
				pact.
					AddInteraction().
					Given("I have an lpa which has been paid for").
					UponReceiving("A request to delete a payment").
					WithCompleteRequest(consumer.Request{
						Method: http.MethodDelete,
						Path:   matchers.String("/lpa-api/v1/payments/123"),
					}).
					WithCompleteResponse(consumer.Response{
						Status: http.StatusNoContent,
					})
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.ExecuteTest(t, func(config consumer.MockServerConfig) error {
				client := NewClient(http.DefaultClient, fmt.Sprintf("http://127.0.0.1:%d", config.Port))

				err := client.DeletePayment(Context{Context: context.Background()}, 123)

				if tc.expectedError == nil {
					assert.Nil(t, err)
				} else {
					assert.Equal(t, tc.expectedError(config.Port), err)
				}
				return nil
			}))
		})
	}
}

func TestDeleteFeeReduction(t *testing.T) {
	t.Parallel()

	pact, err := newPact()
	assert.NoError(t, err)

	testCases := []struct {
		name          string
		setup         func()
		expectedError func(int) error
	}{
		{
			name: "OK",
			setup: func() {
				pact.
					AddInteraction().
					Given("I have an lpa which has a fee reduction").
					UponReceiving("A request to delete a fee reduction").
					WithCompleteRequest(consumer.Request{
						Method: http.MethodDelete,
						Path:   matchers.String("/lpa-api/v1/payments/124"),
					}).
					WithCompleteResponse(consumer.Response{
						Status: http.StatusNoContent,
					})
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.ExecuteTest(t, func(config consumer.MockServerConfig) error {
				client := NewClient(http.DefaultClient, fmt.Sprintf("http://127.0.0.1:%d", config.Port))

				err := client.DeletePayment(Context{Context: context.Background()}, 124)

				if tc.expectedError == nil {
					assert.Nil(t, err)
				} else {
					assert.Equal(t, tc.expectedError(config.Port), err)
				}
				return nil
			}))
		})
	}
}
