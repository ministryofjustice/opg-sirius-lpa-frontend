package sirius

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/pact-foundation/pact-go/dsl"
	"github.com/stretchr/testify/assert"
)

func TestSearchUsers(t *testing.T) {
	t.Parallel()

	pact := newPact()
	defer pact.Teardown()

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
					WithRequest(dsl.Request{
						Method: http.MethodGet,
						Path:   dsl.String("/lpa-api/v1/search/users"),
						Query: dsl.MapMatcher{
							"query": dsl.String("admin"),
						},
					}).
					WillRespondWith(dsl.Response{
						Status:  http.StatusOK,
						Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/json")},
						Body: dsl.EachLike(map[string]interface{}{
							"id":          dsl.Like(47),
							"displayName": dsl.String("system admin"),
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

			assert.Nil(t, pact.Verify(func() error {
				client := NewClient(http.DefaultClient, fmt.Sprintf("http://localhost:%d", pact.Server.Port))

				users, err := client.SearchUsers(Context{Context: context.Background()}, "admin")
				assert.Equal(t, tc.expectedResponse, users)
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

func TestSearchUsersTooShort(t *testing.T) {
	client := NewClient(http.DefaultClient, "")

	expectedErr := ValidationError{
		Detail: "Search term must be at least three characters",
		Field: FieldErrors{
			"search": {"reason": "Search term must be at least three characters"},
		},
	}

	users, err := client.SearchUsers(Context{Context: context.Background()}, "ad")
	assert.Nil(t, users)
	assert.Equal(t, expectedErr, err)
}
