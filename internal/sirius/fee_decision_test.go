package sirius

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/shared"
	"github.com/pact-foundation/pact-go/v2/consumer"
	"github.com/pact-foundation/pact-go/v2/matchers"
	"github.com/stretchr/testify/assert"
)

func TestCreateFeeDecisionBadJSON(t *testing.T) {
	client := NewClient(http.DefaultClient, "http://not/real/server")

	err := client.CreateFeeDecision(
		Context{Context: context.Background()},
		0,
		"DECLINED_REMISSION",
		"None",
		"00/11/2999",
	)

	assert.ErrorContains(t, err, "failed to format non-date")
}

func TestCreateFeeDecision(t *testing.T) {
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
				err := client.CreateFeeDecision(
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

func TestDeleteFeeDecision(t *testing.T) {
	t.Parallel()

	pact, err := newPact()
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
					Given("I have a case with declined fee reductions").
					UponReceiving("A request to delete a fee decision").
					WithCompleteRequest(consumer.Request{
						Method: http.MethodDelete,
						Path:   matchers.String("/lpa-api/v1/fee-decisions/1245"),
					}).
					WithCompleteResponse(consumer.Response{
						Status: http.StatusNoContent,
					})
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.ExecuteTest(t, func(config consumer.MockServerConfig) error {
				client := NewClient(http.DefaultClient, fmt.Sprintf("http://127.0.0.1:%d", config.Port))

				err := client.DeleteFeeDecision(Context{Context: context.Background()}, 1245)

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

func TestGetFeeDecision(t *testing.T) {
	t.Parallel()

	pact, err := newPact()
	assert.NoError(t, err)

	testCases := []struct {
		name             string
		setup            func()
		expectedError    func(int) error
		expectedResponse FeeDecision
	}{
		{
			name: "OK",
			setup: func() {
				pact.
					AddInteraction().
					Given("I have a case with declined fee reductions").
					UponReceiving("A request to get a declined fee decision").
					WithCompleteRequest(consumer.Request{
						Method: http.MethodGet,
						Path:   matchers.String("/lpa-api/v1/fee-decisions/1245"),
					}).
					WithCompleteResponse(consumer.Response{
						Status: http.StatusOK,
						Body: matchers.Like(map[string]interface{}{
							"id":             matchers.Like(1245),
							"decisionType":   matchers.String("DECLINED_EXEMPTION"),
							"decisionReason": matchers.String("Test"),
							"decisionDate":   matchers.String("01/03/2026"),
						}),
						Headers: matchers.MapMatcher{"Content-Type": matchers.String("application/json")},
					})
			},
			expectedResponse: FeeDecision{
				ID:             1245,
				DecisionType:   shared.DecisionTypeExemption,
				DecisionReason: "Test",
				DecisionDate:   "01/03/2026",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.ExecuteTest(t, func(config consumer.MockServerConfig) error {
				client := NewClient(http.DefaultClient, fmt.Sprintf("http://127.0.0.1:%d", config.Port))

				feeDecision, err := client.GetFeeDecision(Context{Context: context.Background()}, 1245)

				assert.Equal(t, tc.expectedResponse, feeDecision)
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
