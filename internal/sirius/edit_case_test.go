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

func TestEditCase(t *testing.T) {
	t.Parallel()

	pact, err := newPact()
	assert.NoError(t, err)

	testCases := []struct {
		name          string
		setup         func()
		caseType      CaseType
		expectedError func(int) error
	}{
		{
			name: "OK LPA",
			setup: func() {
				pact.
					AddInteraction().
					Given("I have a pending case assigned").
					UponReceiving("A request to edit the LPA").
					WithCompleteRequest(consumer.Request{
						Method: http.MethodPut,
						Path:   matchers.String("/lpa-api/v1/lpas/800"),
						Body: map[string]interface{}{
							"status":               "Cancelled",
							"expectedPaymentTotal": 8000,
						},
					}).
					WithCompleteResponse(consumer.Response{
						Status:  http.StatusOK,
						Headers: matchers.MapMatcher{"Content-Type": matchers.String("application/json")},
					})
			},
			caseType: CaseTypeLpa,
		},
		{
			name: "OK EPA",
			setup: func() {
				pact.
					AddInteraction().
					Given("I have a pending EPA assigned").
					UponReceiving("A request to edit the EPA").
					WithCompleteRequest(consumer.Request{
						Method: http.MethodPut,
						Path:   matchers.String("/lpa-api/v1/epas/800"),
						Body: map[string]interface{}{
							"status":               "Cancelled",
							"expectedPaymentTotal": 8000,
						},
					}).
					WithCompleteResponse(consumer.Response{
						Status:  http.StatusOK,
						Headers: matchers.MapMatcher{"Content-Type": matchers.String("application/json")},
					})
			},
			caseType: CaseTypeEpa,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.ExecuteTest(t, func(config consumer.MockServerConfig) error {
				client := NewClient(http.DefaultClient, fmt.Sprintf("http://127.0.0.1:%d", config.Port))

				err := client.EditCase(
					Context{Context: context.Background()},
					800, tc.caseType,
					Case{Status: "Cancelled", ExpectedPaymentTotal: 8000},
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
