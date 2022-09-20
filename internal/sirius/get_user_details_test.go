package sirius

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/pact-foundation/pact-go/dsl"
	"github.com/stretchr/testify/assert"
)

func TestGetUserDetails(t *testing.T) {
	t.Parallel()

	pact := newPact()
	defer pact.Teardown()

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
					Given("A donor exists").
					UponReceiving("A request for the current user").
					WithRequest(dsl.Request{
						Method: http.MethodGet,
						Path:   dsl.String("/lpa-api/v1/users/current"),
					}).
					WillRespondWith(dsl.Response{
						Status: http.StatusOK,
						Body: dsl.Like(map[string]interface{}{
							"id":          dsl.Like(103),
							"displayName": dsl.String("Test User"),
							"roles":       dsl.Like([]string{"OPG User", "Case Manager"}),
						}),
						Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/json")},
					})
			},
			expectedResponse: User{
				ID:          103,
				DisplayName: "Test User",
				Roles:       []string{"OPG User", "Case Manager"},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.Verify(func() error {
				client := NewClient(http.DefaultClient, fmt.Sprintf("http://localhost:%d", pact.Server.Port))

				currentUser, err := client.GetUserDetails(Context{Context: context.Background()})

				assert.Equal(t, tc.expectedResponse, currentUser)
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

func TestHasRole(t *testing.T) {
	t.Run("with reduced fees user role", func(t *testing.T) {
		assert.True(t, User{Roles: []string{"OPG User", "Reduced Fees User"}}.HasRole("Reduced Fees User"))
	})

	t.Run("without reduced fees user role", func(t *testing.T) {
		assert.False(t, User{Roles: []string{"OPG User", "Case Manager"}}.HasRole("Reduced Fees User"))
	})
}
