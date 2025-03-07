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

func TestDeletedCases(t *testing.T) {
	t.Parallel()

	pact, err := newPact2()
	assert.NoError(t, err)

	testCases := []struct {
		name             string
		setup            func()
		expectedResponse []DeletedCase
		expectedError    func(int) error
	}{
		{
			name: "OK",
			setup: func() {
				pact.
					AddInteraction().
					Given("I have deleted a case").
					UponReceiving("A search for the deleted case").
					WithCompleteRequest(consumer.Request{
						Method: http.MethodGet,
						Path:   matchers.String("/lpa-api/v1/deleted-cases"),
						Query: matchers.MapMatcher{
							"uid": matchers.String("700000005555"),
						},
					}).
					WithCompleteResponse(consumer.Response{
						Status:  http.StatusOK,
						Headers: matchers.MapMatcher{"Content-Type": matchers.String("application/json")},
						Body: matchers.EachLike(map[string]interface{}{
							"uId":            matchers.String("7000-0000-5555"),
							"type":           matchers.String("LPA"),
							"status":         matchers.String("Return - unpaid"),
							"deletedAt":      matchers.String("02/12/2022"),
							"deletionReason": matchers.String("LPA was not paid for after 12 months"),
						}, 1),
					})
			},
			expectedResponse: []DeletedCase{
				{
					UID:       "7000-0000-5555",
					Type:      "LPA",
					Status:    "Return - unpaid",
					DeletedAt: DateString("2022-12-02"),
					Reason:    "LPA was not paid for after 12 months",
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.ExecuteTest(t, func(config consumer.MockServerConfig) error {
				client := NewClient(http.DefaultClient, fmt.Sprintf("http://127.0.0.1:%d", config.Port))

				deletedCases, err := client.DeletedCases(Context{Context: context.Background()}, "700000005555")

				assert.Equal(t, tc.expectedResponse, deletedCases)
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
