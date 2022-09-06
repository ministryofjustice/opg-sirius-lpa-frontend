package sirius

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/pact-foundation/pact-go/dsl"
	"github.com/stretchr/testify/assert"
)

func TestMiConfig(t *testing.T) {
	t.Parallel()

	pact := newPact()
	defer pact.Teardown()

	testCases := []struct {
		name           string
		setup          func()
		cookies        []*http.Cookie
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
					WithRequest(dsl.Request{
						Method: http.MethodGet,
						Path:   dsl.String("/api/reporting/config"),
					}).
					WillRespondWith(dsl.Response{
						Status:  http.StatusOK,
						Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/json")},
						Body: dsl.Like(map[string]interface{}{
							"data": dsl.Like(map[string]interface{}{
								"items": dsl.EachLike(map[string]interface{}{
									"properties": dsl.Like(map[string]interface{}{
										"reportType": dsl.Like(map[string]interface{}{
											"description": dsl.String("radio"),
											"type":        dsl.String("reportType"),
											"enum": dsl.EachLike(map[string]interface{}{
												"name":        dsl.String("epasReceived"),
												"description": dsl.String("Number of EPAs received"),
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

			assert.Nil(t, pact.Verify(func() error {
				client := NewClient(http.DefaultClient, fmt.Sprintf("http://localhost:%d", pact.Server.Port))

				config, err := client.MiConfig(Context{Context: context.Background()})

				assert.Equal(t, tc.expectedResult, config)
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
