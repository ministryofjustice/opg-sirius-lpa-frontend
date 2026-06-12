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

func TestGetUserPermissions(t *testing.T) {
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
					UponReceiving("A request for the user's permissions").
					WithCompleteRequest(consumer.Request{
						Method: http.MethodGet,
						Path:   matchers.String("/lpa-api/v1/permissions"),
					}).
					WithCompleteResponse(consumer.Response{
						Status: http.StatusOK,
						Body: matchers.Like(map[string]interface{}{
							"keep-alive": matchers.Like(map[string]interface{}{
								"permissions":  matchers.EachLike("GET", 1),
								"restrictions": matchers.EachLike("GET", 0),
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

				permissions, err := client.GetUserPermissions(Context{Context: context.Background()})

				assert.NotEmpty(t, permissions)

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

func TestIncludesPermission(t *testing.T) {
	t.Run("with permission", func(t *testing.T) {
		assert.True(t, Permissions{"v1-persons": PermissionType{Permissions: []string{"GET"}}}.Includes("v1-persons", "GET"))
	})

	t.Run("without permission type", func(t *testing.T) {
		assert.False(t, Permissions{"v1-persons": PermissionType{Permissions: []string{"GET"}}}.Includes("v1-test", "GET"))
	})

	t.Run("without permission method", func(t *testing.T) {
		assert.False(t, Permissions{"v1-persons": PermissionType{Permissions: []string{"GET"}}}.Includes("v1-persons", "POST"))
	})
}
