package sirius

import (
	"context"
	"fmt"
	"github.com/pact-foundation/pact-go/dsl"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestAddFeeDecision(t *testing.T) {
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
					Given("I have a case assigned").
					UponReceiving("A request to add a fee decision").
					WithRequest(dsl.Request{
						Method: http.MethodPost,
						Path:   dsl.String("/lpa-api/v1/cases/999999/fee-decisions"),
						Headers: dsl.MapMatcher{
							"Content-Type": dsl.String("application/json"),
						},
						Body: map[string]interface{}{
							"decisionType": "DECLINED_REMISSION",
							"decisionReason": "Insufficient evidence",
							"decisionDate": "18/10/2023",
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

				// ctx Context, caseID int, decisionType string, decisionReason string, decisionDate DateString
				err := client.AddFeeDecision(
					Context{Context: context.Background()},
					999999,
					"DECLINED_REMISSION",
					"Insufficient evidence",
					DateString("2023-10-18"),
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
