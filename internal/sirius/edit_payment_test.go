package sirius

import (
	"context"
	"fmt"
	"github.com/pact-foundation/pact-go/v2/consumer"
	"github.com/pact-foundation/pact-go/v2/matchers"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestEditPayment(t *testing.T) {
	t.Parallel()

	pact, err := newPact2()
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

func TestEditFeeReduction(t *testing.T) {
	t.Parallel()

	pact, err := newPact2()
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
					UponReceiving("A request to edit a fee reduction").
					WithCompleteRequest(consumer.Request{
						Method: http.MethodPut,
						Path:   matchers.String("/lpa-api/v1/payments/124"),
						Headers: matchers.MapMatcher{
							"Content-Type": matchers.String("application/json"),
						},
						Body: map[string]interface{}{
							"id":     124,
							"amount": 4100,
							"case": map[string]interface{}{
								"id":     802,
								"status": "",
							},
							"paymentEvidence":  matchers.String("Edited test evidence"),
							"feeReductionType": matchers.String("REMISSION"),
							"paymentDate":      matchers.Term("28/04/2022", `^\d{1,2}/\d{1,2}/\d{4}$`),
							"source":           matchers.String(FeeReductionSource),
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
					assert.Equal(t, tc.expectedError(config.Port), err)
				}
				return nil
			}))
		})
	}
}
