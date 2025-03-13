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

func TestGetUserDetails(t *testing.T) {
	t.Parallel()

	pact, err := newPact()
	assert.NoError(t, err)

	testCases := []struct {
		name             string
		setup            func()
		expectedResponse User
		expectedError    func(int) error
	}{
		{
			name: "OK",
			setup: func() {
				pact.
					AddInteraction().
					Given("I am a reduced fees user").
					UponReceiving("A request for the current user").
					WithCompleteRequest(consumer.Request{
						Method: http.MethodGet,
						Path:   matchers.String("/lpa-api/v1/users/current"),
					}).
					WithCompleteResponse(consumer.Response{
						Status: http.StatusOK,
						Body: matchers.Like(map[string]interface{}{
							"id":          matchers.Like(104),
							"displayName": matchers.String("Test User"),
							"roles":       matchers.Like([]string{"OPG User", "Reduced Fees User"}),
						}),
						Headers: matchers.MapMatcher{"Content-Type": matchers.String("application/json")},
					})
			},
			expectedResponse: User{
				ID:          104,
				DisplayName: "Test User",
				Roles:       []string{"OPG User", "Reduced Fees User"},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.ExecuteTest(t, func(config consumer.MockServerConfig) error {
				client := NewClient(http.DefaultClient, fmt.Sprintf("http://127.0.0.1:%d", config.Port))

				currentUser, err := client.GetUserDetails(Context{Context: context.Background()})

				assert.Equal(t, tc.expectedResponse, currentUser)
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

func TestHasRole(t *testing.T) {
	t.Run("with reduced fees user role", func(t *testing.T) {
		assert.True(t, User{Roles: []string{"OPG User", "Reduced Fees User"}}.HasRole("Reduced Fees User"))
	})

	t.Run("without reduced fees user role", func(t *testing.T) {
		assert.False(t, User{Roles: []string{"OPG User", "Case Manager"}}.HasRole("Reduced Fees User"))
	})
}
