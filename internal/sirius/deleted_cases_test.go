package sirius

import (
	"context"
	"fmt"
	"github.com/pact-foundation/pact-go/dsl"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestDeletedCases(t *testing.T) {
	t.Parallel()

	pact := newPact()
	defer pact.Teardown()

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
					WithRequest(dsl.Request{
						Method: http.MethodGet,
						Path:   dsl.String("/lpa-api/v1/deleted-cases"),
						Query: dsl.MapMatcher{
							"uid": dsl.String("700000005555"),
						},
					}).
					WillRespondWith(dsl.Response{
						Status:  http.StatusOK,
						Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/json")},
						Body: dsl.EachLike(map[string]interface{}{
							"uId":            dsl.String("7000-0000-5555"),
							"type":           dsl.String("LPA"),
							"onlineLpaId":    dsl.String("A987654321"),
							"status":         dsl.String("Return - unpaid"),
							"deletedAt":      dsl.String("02/12/2022"),
							"deletionReason": dsl.String("LPA was not paid for after 12 months"),
						}, 1),
					})
			},
			expectedResponse: []DeletedCase{
				{
					UID:         "7000-0000-5555",
					Type:        "LPA",
					OnlineLpaId: "A987654321",
					Status:      "Return - unpaid",
					DeletedAt:   DateString("2022-12-02"),
					Reason:      "LPA was not paid for after 12 months",
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.Verify(func() error {
				client := NewClient(http.DefaultClient, fmt.Sprintf("http://localhost:%d", pact.Server.Port))

				deletedCases, err := client.DeletedCases(Context{Context: context.Background()}, "700000005555")

				assert.Equal(t, tc.expectedResponse, deletedCases)
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
