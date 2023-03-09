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

func TestPersonFiltersInactiveActors(t *testing.T) {
	actor1 := Person{ID: 1, SystemStatus: true}
	actor2 := Person{ID: 2, SystemStatus: true}
	inactiveActor1 := Person{ID: 3, SystemStatus: false}
	inactiveActor2 := Person{ID: 4, SystemStatus: false}
	persons := []Person{actor1, actor2, inactiveActor1, inactiveActor2}
	activeActors := FilterInactiveAttorneys(persons)

	assert.Equal(t, 2, len(activeActors))
	assert.NotContains(t, activeActors, inactiveActor1)
	assert.NotContains(t, activeActors, inactiveActor2)
}
