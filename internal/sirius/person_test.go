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

func TestPerson(t *testing.T) {
	t.Parallel()

	pact, err := newPact2()
	assert.NoError(t, err)

	testCases := []struct {
		name             string
		setup            func()
		expectedResponse Person
		expectedError    func(int) error
	}{
		{
			name: "OK without children",
			setup: func() {
				pact.
					AddInteraction().
					Given("A donor exists").
					UponReceiving("A request for the person").
					WithCompleteRequest(consumer.Request{
						Method: http.MethodGet,
						Path:   matchers.String("/lpa-api/v1/persons/188"),
					}).
					WithCompleteResponse(consumer.Response{
						Status: http.StatusOK,
						Body: matchers.Like(map[string]interface{}{
							"id":        matchers.Like(188),
							"uId":       matchers.Term("7000-0000-0007", `7\d{3}-\d{4}-\d{4}`),
							"firstname": matchers.String("John"),
							"surname":   matchers.String("Doe"),
							"dob":       matchers.Term("05/05/1970", `^\d{1,2}/\d{1,2}/\d{4}$`),
						}),
						Headers: matchers.MapMatcher{"Content-Type": matchers.String("application/json")},
					})
			},
			expectedResponse: Person{
				ID:          188,
				UID:         "7000-0000-0007",
				Firstname:   "John",
				Surname:     "Doe",
				DateOfBirth: DateString("1970-05-05"),
			},
		},
		{
			name: "OK with children",
			setup: func() {
				pact.
					AddInteraction().
					Given("A donor exists with children").
					UponReceiving("A request for the person with children").
					WithCompleteRequest(consumer.Request{
						Method: http.MethodGet,
						Path:   matchers.String("/lpa-api/v1/persons/189"),
					}).
					WithCompleteResponse(consumer.Response{
						Status: http.StatusOK,
						Body: matchers.Like(map[string]interface{}{
							"id":        matchers.Like(189),
							"uId":       matchers.Term("7000-0000-0001", `7\d{3}-\d{4}-\d{4}`),
							"firstname": matchers.String("John"),
							"surname":   matchers.String("Doe"),
							"dob":       matchers.Term("01/01/1970", `^\d{1,2}/\d{1,2}/\d{4}$`),
							"children": matchers.Like([]map[string]interface{}{
								{
									"id":        matchers.Like(105),
									"uId":       matchers.Term("7000-0000-0002", `7\d{3}-\d{4}-\d{4}`),
									"firstname": matchers.String("Child"),
									"surname":   matchers.String("One"),
								},
							}),
						}),
						Headers: matchers.MapMatcher{"Content-Type": matchers.String("application/json")},
					})
			},
			expectedResponse: Person{
				ID:          189,
				UID:         "7000-0000-0001",
				Firstname:   "John",
				Surname:     "Doe",
				DateOfBirth: DateString("1970-01-01"),
				Children: []Person{
					{
						ID:        105,
						UID:       "7000-0000-0002",
						Firstname: "Child",
						Surname:   "One",
					},
				},
			},
		},
		{
			name: "OK with multiple cases",
			setup: func() {
				pact.
					AddInteraction().
					Given("A donor exists with more than 1 case").
					UponReceiving("A request for the person with multiple cases").
					WithCompleteRequest(consumer.Request{
						Method: http.MethodGet,
						Path:   matchers.String("/lpa-api/v1/persons/400"),
					}).
					WithCompleteResponse(consumer.Response{
						Status: http.StatusOK,
						Body: matchers.Like(map[string]interface{}{
							"id":        matchers.Like(400),
							"uId":       matchers.Term("7000-0000-0001", `7\d{3}-\d{4}-\d{4}`),
							"firstname": matchers.String("John"),
							"surname":   matchers.String("Doe"),
							"dob":       matchers.Term("01/01/1970", `^\d{1,2}/\d{1,2}/\d{4}$`),
							"cases": []map[string]interface{}{
								{
									"id":          matchers.Like(405),
									"uId":         matchers.Term("7000-5382-4438", `\d{4}-\d{4}-\d{4}`),
									"caseSubtype": matchers.Term("pfa", "hw|pfa"),
									"caseType":    matchers.Like("LPA"),
								},
								{
									"id":          matchers.Like(406),
									"uId":         matchers.Term("7000-5382-8764", `\d{4}-\d{4}-\d{4}`),
									"caseSubtype": matchers.Term("hw", "hw|pfa"),
									"caseType":    matchers.Like("LPA"),
								},
							},
						}),
						Headers: matchers.MapMatcher{"Content-Type": matchers.String("application/json")},
					})
			},
			expectedResponse: Person{
				ID:          400,
				UID:         "7000-0000-0001",
				Firstname:   "John",
				Surname:     "Doe",
				DateOfBirth: DateString("1970-01-01"),
				Cases: []*Case{
					{
						ID:       405,
						UID:      "7000-5382-4438",
						CaseType: "LPA",
						SubType:  "pfa",
					},
					{
						ID:       406,
						UID:      "7000-5382-8764",
						CaseType: "LPA",
						SubType:  "hw",
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

				caseitem, err := client.Person(Context{Context: context.Background()}, tc.expectedResponse.ID)

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
