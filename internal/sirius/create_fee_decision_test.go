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

	pact, err := newPact()
	assert.NoError(t, err)

	testCases := []struct {
		name          string
		description   string
		caseId        int
		request       map[string]string
		response      func() consumer.Response
		expectedError func(int) error
	}{
		{
			name:        "OK",
			description: "A valid fee decision request",
			request: map[string]string{
				"decisionType":   "DECLINED_REMISSION",
				"decisionReason": "Insufficient evidence",
				"decisionDate":   "18/10/2023",
			},
			response: func() consumer.Response {
				return consumer.Response{
					Status:  http.StatusCreated,
					Headers: matchers.MapMatcher{"Content-Type": matchers.String("application/json")},
				}
			},
		},
		{
			name:        "ValidationError",
			description: "A fee decision request with invalid data",
			request: map[string]string{
				"decisionType":   "",
				"decisionReason": "Some reason",
				"decisionDate":   "18/10/2023",
			},
			response: func() consumer.Response {
				body := map[string]interface{}{
					"validation_errors": map[string]interface{}{
						"decisionType": map[string]string{
							"isEmpty": "Value is required and can't be empty",
						},
					},
					"detail": "Payload failed validation",
					"type":   "http://www.w3.org/Protocols/rfc2616/rfc2616-sec10.html",
					"status": 400,
					"title":  "Bad Request",
				}

				return consumer.Response{
					Status:  http.StatusBadRequest,
					Headers: matchers.MapMatcher{"Content-Type": matchers.String("application/problem+json")},
					Body:    body,
				}
			},
			expectedError: func(pactPort int) error {
				return ValidationError{
					Field: FieldErrors{
						"decisionType": map[string]string{
							"isEmpty": "Value is required and can't be empty",
						},
					},
					Detail: "Payload failed validation",
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			request := consumer.Request{
				Method: http.MethodPost,
				Path:   matchers.String("/lpa-api/v1/cases/1/fee-decisions"),
				Headers: matchers.MapMatcher{
					"Content-Type": matchers.String("application/json"),
				},
				Body: tc.request,
			}

			pact.
				AddInteraction().
				Given("a digital LPA exists awaiting fee decision").
				UponReceiving(tc.description).
				WithCompleteRequest(request).
				WithCompleteResponse(tc.response())

			assert.Nil(t, pact.ExecuteTest(t, func(config consumer.MockServerConfig) error {
				client := NewClient(http.DefaultClient, fmt.Sprintf("http://127.0.0.1:%d", config.Port))

				// ctx Context, caseID int, decisionType string, decisionReason string, decisionDate DateString
				err := client.AddFeeDecision(
					Context{Context: context.Background()},
					1,
					tc.request["decisionType"],
					tc.request["decisionReason"],
					DateString("2023-10-18"),
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
