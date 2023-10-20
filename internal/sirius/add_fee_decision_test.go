package sirius

import (
	"context"
	"fmt"
	"github.com/pact-foundation/pact-go/dsl"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestAddFeeDecisionBadJSON(t *testing.T) {
	client := NewClient(http.DefaultClient, "http://not/real/server")

	err := client.AddFeeDecision(
		Context{Context: context.Background()},
		0,
		"DECLINED_REMISSION",
		"None",
		"00/11/2999",
	)

	assert.ErrorContains(t, err, "failed to format non-date")
}

func TestAddFeeDecision(t *testing.T) {
	t.Parallel()

	pact := newPact()
	defer pact.Teardown()

	testCases := []struct {
		name          string
		caseId        int
		setup         func()
		expectedError func(int) error
	}{
		{
			name: "OK",
			caseId: 999999,
			setup: func() {
				pact.
					AddInteraction().
					Given("I have case 999999 assigned").
					UponReceiving("A request to add a fee decision for case 999999").
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
		{
			name: "InternalServerError",
			caseId: 111,
			setup: func() {
				pact.
					AddInteraction().
					Given("I have case 111 assigned").
					UponReceiving("A request to add a fee decision for case 111").
					WithRequest(dsl.Request{
						Method: http.MethodPost,
						Path:   dsl.String("/lpa-api/v1/cases/111/fee-decisions"),
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
						Status:  http.StatusInternalServerError,
						Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/json")},
					})
			},
			expectedError: func(pactPort int) error {
				return StatusError{
					Code: http.StatusInternalServerError,
					URL: fmt.Sprintf("http://localhost:%d/lpa-api/v1/cases/111/fee-decisions", pactPort),
					Method: http.MethodPost,
					CorrelationId: "",
				}
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
					tc.caseId,
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
