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
		description   string
		caseId        int
		request       map[string]string
		response      func() dsl.Response
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
			response: func() dsl.Response {
				return dsl.Response{
					Status:  http.StatusCreated,
					Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/json")},
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
			response: func() dsl.Response {
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

				return dsl.Response{
					Status:  http.StatusBadRequest,
					Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/problem+json")},
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
			request := dsl.Request{
				Method: http.MethodPost,
				Path:   dsl.String("/lpa-api/v1/cases/1/fee-decisions"),
				Headers: dsl.MapMatcher{
					"Content-Type": dsl.String("application/json"),
				},
				Body: tc.request,
			}

			pact.
				AddInteraction().
				Given("a digital LPA exists awaiting fee decision").
				UponReceiving(tc.description).
				WithRequest(request).
				WillRespondWith(tc.response())

			assert.Nil(t, pact.Verify(func() error {
				client := NewClient(http.DefaultClient, fmt.Sprintf("http://localhost:%d", pact.Server.Port))

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
					assert.Equal(t, tc.expectedError(pact.Server.Port), err)
				}
				return nil
			}))
		})
	}
}
