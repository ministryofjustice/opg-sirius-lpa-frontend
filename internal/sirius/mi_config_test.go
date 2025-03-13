package sirius

import (
	"context"
	"fmt"
	"github.com/pact-foundation/pact-go/v2/consumer"
	"github.com/pact-foundation/pact-go/v2/matchers"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMiConfig(t *testing.T) {
	t.Parallel()

	pact, err := newPact()
	assert.NoError(t, err)

	testCases := []struct {
		name           string
		setup          func()
		expectedResult map[string]MiConfigProperty
		expectedError  func(int) error
	}{
		{
			name: "OK",
			setup: func() {
				pact.
					AddInteraction().
					Given("I have a pending case assigned").
					UponReceiving("A request for the MI config").
					WithCompleteRequest(consumer.Request{
						Method: http.MethodGet,
						Path:   matchers.String("/api/reporting/config"),
					}).
					WithCompleteResponse(consumer.Response{
						Status:  http.StatusOK,
						Headers: matchers.MapMatcher{"Content-Type": matchers.String("application/json")},
						Body: matchers.Like(map[string]interface{}{
							"data": matchers.Like(map[string]interface{}{
								"items": matchers.EachLike(map[string]interface{}{
									"properties": matchers.Like(map[string]interface{}{
										"reportType": matchers.Like(map[string]interface{}{
											"description": matchers.String("radio"),
											"type":        matchers.String("reportType"),
											"enum": matchers.EachLike(map[string]interface{}{
												"name":        matchers.String("epasReceived"),
												"description": matchers.String("Number of EPAs received"),
											}, 1),
										}),
									}),
								}, 1),
							}),
						}),
					})
			},
			expectedResult: map[string]MiConfigProperty{
				"reportType": {
					Description: "radio",
					Type:        "reportType",
					Enum: []MiConfigEnum{
						{Name: "epasReceived", Description: "Number of EPAs received"},
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.ExecuteTest(t, func(config consumer.MockServerConfig) error {
				client := NewClient(http.DefaultClient, fmt.Sprintf("http://127.0.0.1:%d", config.Port))

				miConfig, err := client.MiConfig(Context{Context: context.Background()})

				assert.Equal(t, tc.expectedResult, miConfig)
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
