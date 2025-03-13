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

func TestCreateInvestigation(t *testing.T) {
	t.Parallel()

	pact, err := newPact2()
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
					Given("I am an investigations user with a pending case assigned").
					UponReceiving("A request to create an investigation on the case").
					WithCompleteRequest(consumer.Request{
						Method: http.MethodPost,
						Path:   matchers.String("/lpa-api/v1/lpas/800/investigations"),
						Headers: matchers.MapMatcher{
							"Content-Type": matchers.String("application/json"),
						},
						Body: matchers.Like(map[string]interface{}{
							"investigationTitle":        matchers.String("Test Investigation"),
							"additionalInformation":     matchers.String("This is an investigation"),
							"type":                      matchers.String("Priority"),
							"investigationReceivedDate": matchers.String("05/04/2022"),
						}),
					}).
					WithCompleteResponse(consumer.Response{
						Status: http.StatusCreated,
					})
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.ExecuteTest(t, func(config consumer.MockServerConfig) error {
				client := NewClient(http.DefaultClient, fmt.Sprintf("http://127.0.0.1:%d", config.Port))

				err := client.CreateInvestigation(Context{Context: context.Background()}, 800, CaseTypeLpa, Investigation{
					Title:        "Test Investigation",
					Information:  "This is an investigation",
					Type:         "Priority",
					DateReceived: DateString("2022-04-05"),
				})

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
