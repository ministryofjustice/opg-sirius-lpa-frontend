package sirius

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pact-foundation/pact-go/dsl"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestAddPayment(t *testing.T) {
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
					Given("I have a pending case assigned").
					UponReceiving("A request to create a payment").
					WithRequest(dsl.Request{
						Method: http.MethodPost,
						Path:   dsl.String("/lpa-api/v1/cases/800/payments"),
						Headers: dsl.MapMatcher{
							"Content-Type": dsl.String("application/json"),
						},
						Body: map[string]interface{}{
							"amount":      4100,
							"source":      "MAKE",
							"paymentDate": "25/04/2022",
						},
					}).
					WillRespondWith(dsl.Response{
						Status:  http.StatusCreated,
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

				err := client.AddPayment(Context{Context: context.Background()}, 800, 4100, "MAKE", DateString("2022-04-25"), "", "", "")

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

func TestAddFeeReduction(t *testing.T) {
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
					Given("I have a pending case assigned").
					UponReceiving("A request to create a fee reduction").
					WithRequest(dsl.Request{
						Method: http.MethodPost,
						Path:   dsl.String("/lpa-api/v1/cases/800/payments"),
						Headers: dsl.MapMatcher{
							"Content-Type": dsl.String("application/json"),
						},
						Body: map[string]interface{}{
							"feeReductionType": "REMISSION",
							"source":           "FEE_REDUCTION",
							"appliedDate":      "25/04/2022",
							"paymentEvidence":  "Test evidence",
						},
					}).
					WillRespondWith(dsl.Response{
						Status:  http.StatusCreated,
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

				err := client.AddPayment(Context{Context: context.Background()}, 800, 0, "FEE_REDUCTION", "", "REMISSION", "Test evidence", DateString("2022-04-25"))

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

func TestGetPostData(t *testing.T) {
	expectedNonReductionPaymentResponse, _ := json.Marshal(struct {
		Amount      int        `json:"amount"`
		Source      string     `json:"source"`
		PaymentDate DateString `json:"paymentDate"`
	}{
		Amount:      100,
		Source:      "MAKE",
		PaymentDate: "02/01/2006",
	})
	expectedFeeReductionPaymentResponse, _ := json.Marshal(struct {
		Source           string     `json:"source"`
		PaymentEvidence  string     `json:"paymentEvidence"`
		FeeReductionType string     `json:"feeReductionType"`
		AppliedDate      DateString `json:"appliedDate"`
	}{
		Source:           "FEE_REDUCTION",
		PaymentEvidence:  "Test",
		FeeReductionType: "REMISSION",
		AppliedDate:      "02/01/2006",
	})

	actualNonReductionPaymentResponse, _ := GetPostData(100, "MAKE", "02/01/2006", "", "", "")
	actualFeeReductionPaymentResponse, _ := GetPostData(0, "FEE_REDUCTION", "", "REMISSION", "Test", "02/01/2006")

	assert.Equal(t, expectedNonReductionPaymentResponse, actualNonReductionPaymentResponse)
	assert.Equal(t, expectedFeeReductionPaymentResponse, actualFeeReductionPaymentResponse)

}
