package sirius

import (
	"context"
	"fmt"
	"github.com/pact-foundation/pact-go/dsl"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestPayment(t *testing.T) {
	t.Parallel()

	pact := newPact()
	defer pact.Teardown()

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
					WithRequest(dsl.Request{
						Method: http.MethodGet,
						Path:   dsl.String("/lpa-api/v1/cases/800/payments"),
					}).
					WillRespondWith(dsl.Response{
						Status: http.StatusOK,
						Body: []map[string]interface{}{{
							"id":          dsl.Like(2),
							"source":      dsl.Like("MAKE"),
							"amount":      dsl.Like(4100),
							"paymentDate": dsl.String("23/01/2022"),
							"case": dsl.Like(map[string]interface{}{
								"id": dsl.Like(800),
							}),
						}, {
							"id":               dsl.Like(3),
							"source":           dsl.Like(FeeReductionSource),
							"feeReductionType": dsl.Like("REMISSION"),
							"paymentEvidence":  dsl.Like("Test\nmultiple\nline evidence"),
							"paymentDate":      dsl.String("24/01/2022"),
							"case": dsl.Like(map[string]interface{}{
								"id": dsl.Like(800),
							}),
						},
						},
						Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/json")},
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
				{
					ID:               3,
					Source:           FeeReductionSource,
					FeeReductionType: "REMISSION",
					PaymentEvidence:  "Test\nmultiple\nline evidence",
					PaymentDate:      DateString("2022-01-24"),
					Case:             &Case{ID: 800},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.Verify(func() error {
				client := NewClient(http.DefaultClient, fmt.Sprintf("http://localhost:%d", pact.Server.Port))

				payments, err := client.Payments(Context{Context: context.Background()}, 800)

				assert.Equal(t, tc.expectedResponse, payments)
				if tc.expectedError == nil {
					assert.Nil(t, err)
				} else {
					assert.Equal(t, tc.expectedError(pact.Server.Port), err)
				}
				return nil
			}))
		})
	}
}

func TestNoPaymentOnCase(t *testing.T) {
	t.Parallel()

	pact := newPact()
	defer pact.Teardown()

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
					WithRequest(dsl.Request{
						Method: http.MethodGet,
						Path:   dsl.String("/lpa-api/v1/cases/801/payments"),
					}).
					WillRespondWith(dsl.Response{
						Status:  http.StatusOK,
						Body:    dsl.EachLike(map[string]interface{}{}, 0),
						Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/json")},
					})
			},
			expectedResponse: []Payment{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.Verify(func() error {
				client := NewClient(http.DefaultClient, fmt.Sprintf("http://localhost:%d", pact.Server.Port))

				payments, err := client.Payments(Context{Context: context.Background()}, 801)

				assert.Equal(t, tc.expectedResponse, payments)
				if tc.expectedError == nil {
					assert.Nil(t, err)
				} else {
					assert.Equal(t, tc.expectedError(pact.Server.Port), err)
				}
				return nil
			}))
		})
	}
}

func TestPaymentByID(t *testing.T) {
	t.Parallel()

	pact := newPact()
	defer pact.Teardown()

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
					WithRequest(dsl.Request{
						Method: http.MethodGet,
						Path:   dsl.String("/lpa-api/v1/payments/123"),
					}).
					WillRespondWith(dsl.Response{
						Status: http.StatusOK,
						Body: dsl.Like(map[string]interface{}{
							"id":          dsl.Like(123),
							"source":      dsl.Like("PHONE"),
							"amount":      dsl.Like(4100),
							"paymentDate": dsl.String("23/01/2022"),
							"case": dsl.Like(map[string]interface{}{
								"id": dsl.Like(800),
							}),
						}),
						Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/json")},
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

			assert.Nil(t, pact.Verify(func() error {
				client := NewClient(http.DefaultClient, fmt.Sprintf("http://localhost:%d", pact.Server.Port))

				payment, err := client.PaymentByID(Context{Context: context.Background()}, 123)

				assert.Equal(t, tc.expectedResponse, payment)
				if tc.expectedError == nil {
					assert.Nil(t, err)
				} else {
					assert.Equal(t, tc.expectedError(pact.Server.Port), err)
				}
				return nil
			}))
		})
	}
}
