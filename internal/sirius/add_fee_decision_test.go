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
			name: "OK",
			description: "Valid request Sirius can handle",
			caseId: 800,
			request: map[string]string{
				"decisionType": "DECLINED_REMISSION",
				"decisionReason": "Insufficient evidence",
				"decisionDate": "18/10/2023",
			},
			response: func() dsl.Response {
				return dsl.Response{
					Status:  http.StatusCreated,
					Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/json")},
				}
			},
		},
		{
			name: "ValidationError",
			description: "Request with invalid data",
			caseId: 800,
			request: map[string]string{
				"decisionType": "",
				"decisionReason": "",
				"decisionDate": "18/10/2023",
			},
			response: func() dsl.Response {
				validationError := ValidationError{
					Field: FieldErrors{
						"decisionType": map[string]string{},
					},
					Detail: "",
				}

				bodyBytes, _ := json.Marshal(validationError)

				return dsl.Response{
					Status:  http.StatusBadRequest,
					Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/problem+json")},
					Body:    string(bodyBytes),
				}
			},
			expectedError: func(pactPort int) error {
				return ValidationError{
					Field: FieldErrors{
						"decisionType": map[string]string{},
					},
					Detail: "",
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			request := dsl.Request{
				Method:  http.MethodPost,
				Path:    dsl.String(fmt.Sprintf("/lpa-api/v1/cases/%d/fee-decisions", tc.caseId)),
				Headers: dsl.MapMatcher{
					"Content-Type": dsl.String("application/json"),
				},
				Body: tc.request,
			}

			pact.
				AddInteraction().
				Given("Request to add a fee decision via Sirius API").
				UponReceiving(tc.description).
				WithRequest(request).
				WillRespondWith(tc.response())

			assert.Nil(t, pact.Verify(func() error {
				client := NewClient(http.DefaultClient, fmt.Sprintf("http://localhost:%d", pact.Server.Port))

				// ctx Context, caseID int, decisionType string, decisionReason string, decisionDate DateString
				err := client.AddFeeDecision(
					Context{Context: context.Background()},
					tc.caseId,
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

/*
		{
			name: "InternalServerError",
			description: "Request which Sirius responds to with a 500 error",
			caseId: 111,
			request: map[string]string{
				"decisionType": "DECLINED_REMISSION",
				"decisionReason": "Insufficient evidence",
				"decisionDate": "18/10/2023",
			},
			response: func() dsl.Response {
				return dsl.Response{
					Status:  http.StatusInternalServerError,
					Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/problem+json")},
				}
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
		*/
