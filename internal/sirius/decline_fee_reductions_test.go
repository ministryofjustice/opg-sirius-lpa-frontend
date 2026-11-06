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
			id:   8521,
			setup: func() {
				pact.
					AddInteraction().
					Given("I have a case with declined fee reductions").
					UponReceiving("A request for declined fee reductions").
					WithCompleteRequest(consumer.Request{
						Method: http.MethodGet,
						Path:   matchers.String("/lpa-api/v1/cases/8521/fee-decisions"),
					}).
					WithCompleteResponse(consumer.Response{
						Status: http.StatusOK,
						Body: matchers.EachLike(map[string]interface{}{
							"id":             matchers.Like(1245),
							"decisionDate":   matchers.String("23/08/2022"),
							"decisionType":   matchers.String("DECLINED_EXEMPTION"),
							"decisionReason": matchers.String("Cant have it sorry"),
						}, 1),
						Headers: matchers.MapMatcher{"Content-Type": matchers.String("application/json")},
					})
			},
			expectedResponse: []DeclineFeeReductions{
				{
					ID:             1245,
					DecisionDate:   "2022-08-23",
					DecisionType:   shared.DecisionTypeExemption,
					DecisionReason: "Cant have it sorry",
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
