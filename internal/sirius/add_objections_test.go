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

func TestAddObjection(t *testing.T) {
	t.Parallel()

	pact, err := newPact()
	assert.NoError(t, err)

	testCases := []struct {
		name           string
		objectionsData AddObjections
		setup          func()
		expectedError  func(int) error
	}{
		{
			name: "OK",
			objectionsData: AddObjections{
				LpaUids:       []string{"M-1234-9876-4567"},
				ReceivedDate:  "2025-01-02",
				ObjectionType: "factual",
				Notes:         "test",
			},
			setup: func() {
				pact.
					AddInteraction().
					Given("A digital LPA exists").
					UponReceiving("A request to add an objection").
					WithCompleteRequest(consumer.Request{
						Method: http.MethodPost,
						Path:   matchers.String("/lpa-api/v1/objections"),
						Headers: matchers.MapMatcher{
							"Content-Type": matchers.String("application/json"),
						},
						Body: matchers.Like(map[string]interface{}{
							"lpaUids":       []string{"M-1234-9876-4567"},
							"receivedDate":  matchers.Like("02/01/2025"),
							"objectionType": matchers.Like("factual"),
							"notes":         matchers.Like("test"),
						}),
					}).
					WithCompleteResponse(consumer.Response{
						Status:  http.StatusCreated,
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

				err := client.AddObjections(Context{Context: context.Background()}, tc.objectionsData)
				if (tc.expectedError) == nil {
					assert.Nil(t, err)
				} else {
					assert.Equal(t, tc.expectedError(config.Port), err)
				}
				return nil
			}))
		})
	}
}
