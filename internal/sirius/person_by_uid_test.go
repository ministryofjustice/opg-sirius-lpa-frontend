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

func TestPersonByUid(t *testing.T) {
	t.Parallel()

	pact, err := newPact()
	assert.NoError(t, err)

	testCases := []struct {
		name             string
		setup            func()
		expectedResponse Person
		expectedError    func(int) error
	}{
		{
			name: "OK",
			setup: func() {
				pact.
					AddInteraction().
					Given("A donor exists").
					UponReceiving("A request for the person by UID").
					WithCompleteRequest(consumer.Request{
						Method: http.MethodGet,
						Path:   matchers.String("/lpa-api/v1/persons/by-uid/7000-0000-0001"),
					}).
					WithCompleteResponse(consumer.Response{
						Status: http.StatusOK,
						Body: matchers.Like(map[string]interface{}{
							"id":        matchers.Like(103),
							"uId":       matchers.Term("7000-0000-0001", `7\d{3}-\d{4}-\d{4}`),
							"firstname": matchers.String("John"),
							"surname":   matchers.String("Doe"),
							"dob":       matchers.Term("01/01/1970", `^\d{1,2}/\d{1,2}/\d{4}$`),
						}),
						Headers: matchers.MapMatcher{"Content-Type": matchers.String("application/json")},
					})
			},
			expectedResponse: Person{
				ID:          103,
				UID:         "7000-0000-0001",
				Firstname:   "John",
				Surname:     "Doe",
				DateOfBirth: DateString("1970-01-01"),
			},
		},
		{
			name: "OK with children",
			setup: func() {
				pact.
					AddInteraction().
					Given("A donor exists with children").
					UponReceiving("A request for the person by UID").
					WithCompleteRequest(consumer.Request{
						Method: http.MethodGet,
						Path:   matchers.String("/lpa-api/v1/persons/by-uid/7000-0000-0001"),
					}).
					WithCompleteResponse(consumer.Response{
						Status: http.StatusOK,
						Body: matchers.Like(map[string]interface{}{
							"id":        matchers.Like(103),
							"uId":       matchers.Term("7000-0000-0001", `7\d{3}-\d{4}-\d{4}`),
							"firstname": matchers.String("John"),
							"surname":   matchers.String("Doe"),
							"dob":       matchers.Term("01/01/1970", `^\d{1,2}/\d{1,2}/\d{4}$`),
							"children": matchers.EachLike(map[string]interface{}{
								"id": matchers.Like(104),
							}, 1),
						}),
						Headers: matchers.MapMatcher{"Content-Type": matchers.String("application/json")},
					})
			},
			expectedResponse: Person{
				ID:          103,
				UID:         "7000-0000-0001",
				Firstname:   "John",
				Surname:     "Doe",
				DateOfBirth: DateString("1970-01-01"),
				Children: []Person{
					{
						ID: 104,
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.ExecuteTest(t, func(config consumer.MockServerConfig) error {
				client := NewClient(http.DefaultClient, fmt.Sprintf("http://127.0.0.1:%d", config.Port))

				caseitem, err := client.PersonByUid(Context{Context: context.Background()}, "7000-0000-0001")

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
