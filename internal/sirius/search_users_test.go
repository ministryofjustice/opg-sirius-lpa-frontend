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

func TestSearchUsers(t *testing.T) {
	t.Parallel()

	pact, err := newPact2()
	assert.NoError(t, err)

	testCases := []struct {
		name             string
		setup            func()
		expectedResponse []User
		expectedError    func(int) error
	}{
		{
			name: "OK",
			setup: func() {
				pact.
					AddInteraction().
					Given("User exists").
					UponReceiving("A search for admin users").
					WithCompleteRequest(consumer.Request{
						Method: http.MethodGet,
						Path:   matchers.String("/lpa-api/v1/search/users"),
						Query: matchers.MapMatcher{
							"query": matchers.String("admin"),
						},
					}).
					WithCompleteResponse(consumer.Response{
						Status:  http.StatusOK,
						Headers: matchers.MapMatcher{"Content-Type": matchers.String("application/json")},
						Body: matchers.EachLike(map[string]interface{}{
							"id":          matchers.Like(47),
							"displayName": matchers.String("system admin"),
						}, 1),
					})
			},
			expectedResponse: []User{
				{
					ID:          47,
					DisplayName: "system admin",
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.ExecuteTest(t, func(config consumer.MockServerConfig) error {
				client := NewClient(http.DefaultClient, fmt.Sprintf("http://127.0.0.1:%d", config.Port))

				users, err := client.SearchUsers(Context{Context: context.Background()}, "admin")
				assert.Equal(t, tc.expectedResponse, users)
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

func TestSearchUsersTooShort(t *testing.T) {
	client := NewClient(http.DefaultClient, "")

	expectedErr := ValidationError{
		Detail: "Search term must be at least three characters",
		Field: FieldErrors{
			"term": {"reason": "Search term must be at least three characters"},
		},
	}

	users, err := client.SearchUsers(Context{Context: context.Background()}, "ad")
	assert.Nil(t, users)
	assert.Equal(t, expectedErr, err)
}
