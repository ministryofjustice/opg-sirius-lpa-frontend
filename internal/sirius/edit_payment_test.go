package sirius

import (
	"context"
	"fmt"
	"github.com/pact-foundation/pact-go/dsl"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestEditPayment(t *testing.T) {
	t.Parallel()

	pact := newPact()
	defer pact.Teardown()

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
					WithRequest(dsl.Request{
						Method: http.MethodPut,
						Path:   dsl.String("/lpa-api/v1/payments/123"),
						Headers: dsl.MapMatcher{
							"Content-Type": dsl.String("application/json"),
						},
						Body: map[string]interface{}{
							"amount":      dsl.Like(2550),
							"source":      dsl.String("PHONE"),
							"paymentDate": dsl.Term("27/04/2022", `^\d{1,2}/\d{1,2}/\d{4}$`),
						},
					}).
					WillRespondWith(dsl.Response{
						Status:  http.StatusOK,
						Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/json")},
					})
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.Verify(func() error {
				client := NewClient(http.DefaultClient, fmt.Sprintf("http://localhost:%d", pact.Server.Port))

				err := client.EditPayment(Context{Context: context.Background()}, 123, Payment{Amount: 2550, Source: "PHONE", PaymentDate: DateString("2022-04-27")})

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

func TestEditFeeReduction(t *testing.T) {
	t.Parallel()

	pact := newPact()
	defer pact.Teardown()

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
					UponReceiving("A request to edit a fee reduction").
					WithRequest(dsl.Request{
						Method: http.MethodPut,
						Path:   dsl.String("/lpa-api/v1/payments/124"),
						Headers: dsl.MapMatcher{
							"Content-Type": dsl.String("application/json"),
						},
						Body: map[string]interface{}{
							"id":     124,
							"amount": 4100,
							"case": map[string]interface{}{
								"id":     802,
								"status": "",
							},
							"paymentEvidence":  dsl.String("Edited test evidence"),
							"feeReductionType": dsl.String("REMISSION"),
							"paymentDate":      dsl.Term("28/04/2022", `^\d{1,2}/\d{1,2}/\d{4}$`),
							"source":           dsl.String(FeeReductionSource),
						},
					}).
					WillRespondWith(dsl.Response{
						Status:  http.StatusOK,
						Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/json")},
					})
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.Verify(func() error {
				client := NewClient(http.DefaultClient, fmt.Sprintf("http://localhost:%d", pact.Server.Port))

				err := client.EditPayment(Context{Context: context.Background()},
					124,
					Payment{
						ID:               124,
						Amount:           4100,
						Case:             &Case{ID: 802, Status: ""},
						PaymentEvidence:  "Edited test evidence",
						FeeReductionType: "REMISSION",
						PaymentDate:      DateString("2022-04-28"),
						Source:           FeeReductionSource,
					},
				)

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
