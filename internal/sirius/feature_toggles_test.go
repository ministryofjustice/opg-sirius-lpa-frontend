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

func TestFeatureToggles(t *testing.T) {
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
					Given("I am a System Admin").
					UponReceiving("A request for the feature toggles").
					WithCompleteRequest(consumer.Request{
						Method: http.MethodGet,
						Path:   matchers.String("/config"),
					}).
					WithCompleteResponse(consumer.Response{
						Status: http.StatusOK,
						Body: matchers.Like(map[string]interface{}{
							"featureToggles": matchers.Like(map[string]bool{
								"poasNewLpaForm": true,
							}),
						}),
						Headers: matchers.MapMatcher{"Content-Type": matchers.String("application/json")},
					})
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.ExecuteTest(t, func(config consumer.MockServerConfig) error {
				client := NewClient(http.DefaultClient, fmt.Sprintf("http://127.0.0.1:%d", config.Port))

				featureToggles, err := client.FeatureToggles(Context{Context: context.Background()})

				assert.NotEmpty(t, featureToggles)

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

func TestEnabled(t *testing.T) {
	t.Run("feature toggle set", func(t *testing.T) {
		assert.True(t, FeatureToggles{"toggle": true}.Enabled("toggle"))
	})

	t.Run("feature toggle not set", func(t *testing.T) {
		assert.False(t, FeatureToggles{"toggle": false}.Enabled("toggle"))
	})

	t.Run("feature toggle does not exist", func(t *testing.T) {
		assert.False(t, FeatureToggles{"notToggle": true}.Enabled("toggle"))
	})
}
