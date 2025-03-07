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

func TestTaskTypes(t *testing.T) {
	t.Parallel()

	pact, err := newPact2()
	assert.NoError(t, err)

	testCases := []struct {
		name             string
		setup            func()
		expectedResponse []string
		expectedError    func(int) error
	}{
		{
			name: "OK",
			setup: func() {
				pact.
					AddInteraction().
					Given("Some task types exist").
					UponReceiving("A request for task types").
					WithCompleteRequest(consumer.Request{
						Method: http.MethodGet,
						Path:   matchers.String("/lpa-api/v1/tasktypes/lpa"),
					}).
					WithCompleteResponse(consumer.Response{
						Status: http.StatusOK,
						Body: matchers.Like(map[string]interface{}{
							"task_types": matchers.Like(map[string]interface{}{
								"Check Application": matchers.Like(map[string]interface{}{}),
							}),
						}),
						Headers: matchers.MapMatcher{"Content-Type": matchers.String("application/json")},
					})
			},
			expectedResponse: []string{"Check Application"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.ExecuteTest(t, func(config consumer.MockServerConfig) error {
				client := NewClient(http.DefaultClient, fmt.Sprintf("http://127.0.0.1:%d", config.Port))

				types, err := client.TaskTypes(Context{Context: context.Background()})

				assert.Equal(t, tc.expectedResponse, types)
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
