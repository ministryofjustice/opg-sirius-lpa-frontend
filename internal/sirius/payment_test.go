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
		cookies          []*http.Cookie
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
						Path:   dsl.String("/lpa-api/v1/cases/9/payments"),
					}).
					WillRespondWith(dsl.Response{
						Status: http.StatusOK,
						Body: dsl.EachLike(map[string]interface{}{
							"id":          dsl.Like(2),
							"source":      dsl.Like("MAKE"),
							"amount":      dsl.Like(4100),
							"paymentDate": dsl.String("23/01/2022"),
						}, 1),
						Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/json")},
					})
			},
			expectedResponse: []Payment{
				{
					ID:          2,
					Source:      "MAKE",
					Amount:      4100,
					PaymentDate: DateString("2022-01-23"),
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.Verify(func() error {
				client := NewClient(http.DefaultClient, fmt.Sprintf("http://localhost:%d", pact.Server.Port))

				payments, err := client.Payments(Context{Context: context.Background()}, 9)

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
						}),
						Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/json")},
					})
			},
			expectedResponse: Payment{
				ID:          123,
				Source:      "PHONE",
				Amount:      4100,
				PaymentDate: DateString("2022-01-23"),
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
