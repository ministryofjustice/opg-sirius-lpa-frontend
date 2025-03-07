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

func TestCreateDocument(t *testing.T) {
	t.Parallel()

	pact, err := newPact2()
	assert.NoError(t, err)

	testCases := []struct {
		name           string
		setup          func()
		expectedResult Document
		expectedError  func(int) error
	}{
		{
			name: "OK",
			setup: func() {
				pact.
					AddInteraction().
					Given("I have a pending case assigned").
					UponReceiving("A request to create a draft document on the case").
					WithCompleteRequest(consumer.Request{
						Method: http.MethodPost,
						Path:   matchers.String("/lpa-api/v1/lpas/800/documents/draft"),
						Headers: matchers.MapMatcher{
							"Content-Type": matchers.String("application/json"),
						},
						Body: matchers.Like(map[string]interface{}{
							"templateId":      matchers.String("DD"),
							"inserts":         []string{"DD1"},
							"correspondentId": 189,
						}),
					}).
					WithCompleteResponse(consumer.Response{
						Status:  http.StatusCreated,
						Headers: matchers.MapMatcher{"Content-Type": matchers.String("application/json")},
						Body: matchers.Like(map[string]interface{}{
							"id": matchers.Like(1),
						}),
					})
			},
			expectedResult: Document{
				ID: 1,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.ExecuteTest(t, func(config consumer.MockServerConfig) error {
				client := NewClient(http.DefaultClient, fmt.Sprintf("http://127.0.0.1:%d", config.Port))

				result, err := client.CreateDocument(Context{Context: context.Background()}, 800, 189, "DD", []string{"DD1"})

				assert.Equal(t, tc.expectedResult, result)
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
