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

func TestPersonReferences(t *testing.T) {
	t.Parallel()

	pact, err := newPact2()
	assert.NoError(t, err)

	testCases := []struct {
		name             string
		setup            func()
		expectedResponse []PersonReference
		expectedError    func(int) error
	}{
		{
			name: "OK",
			setup: func() {
				pact.
					AddInteraction().
					Given("A donor with a reference").
					UponReceiving("A request for person references").
					WithCompleteRequest(consumer.Request{
						Method: http.MethodGet,
						Path:   matchers.String("/lpa-api/v1/persons/189/references"),
					}).
					WithCompleteResponse(consumer.Response{
						Status: http.StatusOK,
						Body: matchers.EachLike(map[string]interface{}{
							"referenceId": matchers.Like(768),
							"id":          matchers.Like(189),
							"uid":         matchers.Like(70000000000),
							"displayName": matchers.String("John Doe"),
							"reason":      matchers.String("Friend"),
						}, 1),
						Headers: matchers.MapMatcher{"Content-Type": matchers.String("application/json")},
					})
			},
			expectedResponse: []PersonReference{{
				ReferenceID: 768,
				ID:          189,
				UID:         70000000000,
				DisplayName: "John Doe",
				Reason:      "Friend",
			}},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.ExecuteTest(t, func(config consumer.MockServerConfig) error {
				client := NewClient(http.DefaultClient, fmt.Sprintf("http://127.0.0.1:%d", config.Port))

				caseitem, err := client.PersonReferences(Context{Context: context.Background()}, 189)

				assert.Equal(t, tc.expectedResponse, caseitem)
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
