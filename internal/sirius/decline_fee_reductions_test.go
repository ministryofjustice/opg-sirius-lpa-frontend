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

func TestDeclineFeeReductions(t *testing.T) {
	t.Parallel()

	pact, err := newPact()
	assert.NoError(t, err)

	testCases := []struct {
		name             string
		id               int
		setup            func()
		expectedResponse []DeclineFeeReductions
		expectedError    func(int) error
	}{
		{
			name: "OK",
			id:   742,
			setup: func() {
				pact.
					AddInteraction().
					Given("I have a case with declined fee reductions").
					UponReceiving("A request for declined fee reductions").
					WithCompleteRequest(consumer.Request{
						Method: http.MethodGet,
						Path:   matchers.String("/lpa-api/v1/cases/742/fee-decisions"),
					}).
					WithCompleteResponse(consumer.Response{
						Status: http.StatusOK,
						Body: matchers.EachLike(map[string]interface{}{
							"id":             matchers.Like(1),
							"decisionDate":   matchers.String("2022-08-23"),
							"decisionType":   matchers.String("DECLINED_EXEMPTION"),
							"decisionReason": matchers.String("Cant have it sorry"),
						}, 1),
						Headers: matchers.MapMatcher{"Content-Type": matchers.String("application/json")},
					})
			},
			expectedResponse: []DeclineFeeReductions{
				{
					ID:             1,
					DecisionDate:   DateString("2022-08-23"),
					DecisionType:   shared.DecisionTypeExemption,
					DecisionReason: "Cant have it sorry",
				},
			},
		},
		{
			name: "Empty",
			id:   901,
			setup: func() {
				pact.
					AddInteraction().
					Given("I have a case with no declined fee reductions").
					UponReceiving("A request for declined fee reductions on a case with none").
					WithCompleteRequest(consumer.Request{
						Method: http.MethodGet,
						Path:   matchers.String("/lpa-api/v1/cases/901/fee-decisions"),
					}).
					WithCompleteResponse(consumer.Response{
						Status:  http.StatusOK,
						Body:    matchers.Like([]interface{}{}),
						Headers: matchers.MapMatcher{"Content-Type": matchers.String("application/json")},
					})
			},
			expectedResponse: []DeclineFeeReductions{},
		},
		{
			name: "404",
			id:   9999,
			setup: func() {
				pact.
					AddInteraction().
					Given("There is no case with the specified ID").
					UponReceiving("A request for declined fee reductions for a non-existent case").
					WithCompleteRequest(consumer.Request{
						Method: http.MethodGet,
						Path:   matchers.String("/lpa-api/v1/cases/9999/fee-decisions"),
					}).
					WithCompleteResponse(consumer.Response{
						Status: http.StatusNotFound,
					})
			},
			expectedError: func(port int) error {
				return StatusError{
					Code:          404,
					URL:           fmt.Sprintf("http://127.0.0.1:%d/lpa-api/v1/cases/9999/fee-decisions", port),
					Method:        "GET",
					CorrelationId: "",
				}
			},
		},
		{
			name: "Multiple",
			id:   500,
			setup: func() {
				pact.
					AddInteraction().
					Given("I have a case with multiple declined fee reductions").
					UponReceiving("A request for multiple declined fee reductions").
					WithCompleteRequest(consumer.Request{
						Method: http.MethodGet,
						Path:   matchers.String("/lpa-api/v1/cases/500/fee-decisions"),
					}).
					WithCompleteResponse(consumer.Response{
						Status: http.StatusOK,
						Body: matchers.EachLike(map[string]interface{}{
							"id":             matchers.Like(1),
							"decisionDate":   matchers.String("2022-08-23"),
							"decisionType":   matchers.String("DECLINED_EXEMPTION"),
							"decisionReason": matchers.String("Reason 1"),
						}, 2),
						Headers: matchers.MapMatcher{"Content-Type": matchers.String("application/json")},
					})
			},
			expectedResponse: []DeclineFeeReductions{
				{
					ID:             1,
					DecisionDate:   DateString("2022-08-23"),
					DecisionType:   shared.DecisionTypeExemption,
					DecisionReason: "Reason 1",
				},
				{
					ID:             1,
					DecisionDate:   DateString("2022-08-23"),
					DecisionType:   shared.DecisionTypeExemption,
					DecisionReason: "Reason 1",
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.ExecuteTest(t, func(config consumer.MockServerConfig) error {
				client := NewClient(http.DefaultClient, fmt.Sprintf("http://127.0.0.1:%d", config.Port))

				declineFeeReductions, err := client.DeclineFeeReductions(Context{Context: context.Background()}, tc.id)

				assert.Equal(t, tc.expectedResponse, declineFeeReductions)
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
