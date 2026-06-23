package sirius

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/pact-foundation/pact-go/v2/consumer"
	"github.com/pact-foundation/pact-go/v2/matchers"

	"github.com/stretchr/testify/assert"
)

func TestMiConfig(t *testing.T) {
	t.Parallel()

	pact, err := newPact()
	assert.NoError(t, err)

	testCases := []struct {
		name           string
		setup          func()
		expectedResult MiConfig
		expectedError  func(int) error
	}{
		{
			name: "OK",
			setup: func() {
				pact.
					AddInteraction().
					Given("I am a System Admin").
					UponReceiving("A request for the MI config").
					WithCompleteRequest(consumer.Request{
						Method: http.MethodGet,
						Path:   matchers.String("/api/reporting/config"),
					}).
					WithCompleteResponse(consumer.Response{
						Status:  http.StatusOK,
						Headers: matchers.MapMatcher{"Content-Type": matchers.String("application/json")},
						Body: matchers.Like(map[string]interface{}{
							"reports": matchers.Like(map[string]interface{}{
								"epasReceived": matchers.Like(map[string]interface{}{
									"class":       matchers.String("EPAReceived"),
									"description": matchers.String("Number of EPAs received"),
									"fields": matchers.EachLike(map[string]interface{}{
										"name": matchers.String("startDate"),
									}, 1),
								}),
							}),
							"fields": matchers.Like(map[string]interface{}{
								"startDate": matchers.Like(map[string]interface{}{
									"type":    matchers.String("date"),
									"maxDate": matchers.String("today"),
								}),
								"applicationType": matchers.Like(map[string]interface{}{
									"type": matchers.String("checkbox"),
									"options": matchers.EachLike(map[string]interface{}{
										"value": matchers.String("HW"),
										"label": matchers.String("HW"),
									}, 1),
								}),
							}),
						}),
					})
			},
			expectedResult: MiConfig{
				Reports: map[string]MiReportConfig{
					"epasReceived": {
						Class:       "EPAReceived",
						Description: "Number of EPAs received",
						Fields: []struct {
							Name     string `json:"name"`
							Optional bool   `json:"optional"`
						}{
							{Name: "startDate", Optional: false},
						},
					},
				},
				Fields: map[string]MiConfigField{
					"startDate": {Type: "date", MaxDate: "today"},
					"applicationType": {Type: "checkbox", Options: []MiConfigEnum{
						{Value: "HW", Label: "HW"},
					}},
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
