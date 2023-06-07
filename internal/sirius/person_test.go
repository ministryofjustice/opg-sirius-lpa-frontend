package sirius

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/pact-foundation/pact-go/dsl"
	"github.com/stretchr/testify/assert"
)

func TestPerson(t *testing.T) {
	t.Parallel()

	pact := newPact()
	defer pact.Teardown()

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
					WithRequest(dsl.Request{
						Method: http.MethodGet,
						Path:   dsl.String("/lpa-api/v1/persons/188"),
					}).
					WillRespondWith(dsl.Response{
						Status: http.StatusOK,
						Body: dsl.Like(map[string]interface{}{
							"id":        dsl.Like(188),
							"uId":       dsl.Term("7000-0000-0007", `7\d{3}-\d{4}-\d{4}`),
							"firstname": dsl.String("John"),
							"surname":   dsl.String("Doe"),
							"dob":       dsl.Term("05/05/1970", `^\d{1,2}/\d{1,2}/\d{4}$`),
						}),
						Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/json")},
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
					WithRequest(dsl.Request{
						Method: http.MethodGet,
						Path:   dsl.String("/lpa-api/v1/persons/189"),
					}).
					WillRespondWith(dsl.Response{
						Status: http.StatusOK,
						Body: dsl.Like(map[string]interface{}{
							"id":        dsl.Like(189),
							"uId":       dsl.Term("7000-0000-0001", `7\d{3}-\d{4}-\d{4}`),
							"firstname": dsl.String("John"),
							"surname":   dsl.String("Doe"),
							"dob":       dsl.Term("01/01/1970", `^\d{1,2}/\d{1,2}/\d{4}$`),
							"children": dsl.Like([]map[string]interface{}{
								{
									"id":        dsl.Like(105),
									"uId":       dsl.Term("7000-0000-0002", `7\d{3}-\d{4}-\d{4}`),
									"firstname": dsl.String("Child"),
									"surname":   dsl.String("One"),
								},
							}),
						}),
						Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/json")},
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
					WithRequest(dsl.Request{
						Method: http.MethodGet,
						Path:   dsl.String("/lpa-api/v1/persons/400"),
					}).
					WillRespondWith(dsl.Response{
						Status: http.StatusOK,
						Body: dsl.Like(map[string]interface{}{
							"id":        dsl.Like(400),
							"uId":       dsl.Term("7000-0000-0001", `7\d{3}-\d{4}-\d{4}`),
							"firstname": dsl.String("John"),
							"surname":   dsl.String("Doe"),
							"dob":       dsl.Term("01/01/1970", `^\d{1,2}/\d{1,2}/\d{4}$`),
							"cases": []map[string]interface{}{
								{
									"id":          dsl.Like(405),
									"uId":         dsl.Term("7000-5382-4438", `\d{4}-\d{4}-\d{4}`),
									"caseSubtype": dsl.Term("pfa", "hw|pfa"),
									"status":      dsl.Like("Perfect"),
									"caseType":    dsl.Like("LPA"),
								},
								{
									"id":          dsl.Like(406),
									"uId":         dsl.Term("7000-5382-8764", `\d{4}-\d{4}-\d{4}`),
									"caseSubtype": dsl.Term("hw", "hw|pfa"),
									"status":      dsl.Like("Pending"),
									"caseType":    dsl.Like("LPA"),
								},
							},
						}),
						Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/json")},
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
						Status:   "Perfect",
					},
					{
						ID:       406,
						UID:      "7000-5382-8764",
						CaseType: "LPA",
						SubType:  "hw",
						Status:   "Pending",
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.Verify(func() error {
				client := NewClient(http.DefaultClient, fmt.Sprintf("http://localhost:%d", pact.Server.Port))

				caseitem, err := client.Person(Context{Context: context.Background()}, tc.expectedResponse.ID)

				assert.Equal(t, tc.expectedResponse, caseitem)
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
