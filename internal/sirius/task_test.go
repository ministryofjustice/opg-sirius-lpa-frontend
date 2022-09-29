package sirius

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/pact-foundation/pact-go/dsl"
	"github.com/stretchr/testify/assert"
)

func TestTask(t *testing.T) {
	t.Parallel()

	pact := newPact()
	defer pact.Teardown()

	testCases := []struct {
		name             string
		setup            func()
		expectedResponse Task
		expectedError    func(int) error
	}{
		{
			name: "OK",
			setup: func() {
				pact.
					AddInteraction().
					Given("I have a case with an open task assigned").
					UponReceiving("A request for a task").
					WithRequest(dsl.Request{
						Method: http.MethodGet,
						Path:   dsl.String("/lpa-api/v1/tasks/990"),
					}).
					WillRespondWith(dsl.Response{
						Status: http.StatusOK,
						Body: dsl.Like(map[string]interface{}{
							"id":      dsl.Like(990),
							"status":  dsl.String("Not Started"),
							"dueDate": dsl.String("10/01/2022"),
							"name":    dsl.String("Create physical case file"),
							"caseItems": dsl.EachLike(map[string]interface{}{
								"uId":      dsl.String("7000-0000-0001"),
								"caseType": dsl.String("LPA"),
							}, 1),
						}),
						Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/json")},
					})
			},
			expectedResponse: Task{
				ID:      990,
				Status:  "Not Started",
				DueDate: DateString("2022-01-10"),
				Name:    "Create physical case file",
				CaseItems: []Case{{
					UID:      "7000-0000-0001",
					CaseType: "LPA",
				}},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.Verify(func() error {
				client := NewClient(http.DefaultClient, fmt.Sprintf("http://localhost:%d", pact.Server.Port))

				task, err := client.Task(Context{Context: context.Background()}, 990)

				assert.Equal(t, tc.expectedResponse, task)
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

func TestTaskSummary(t *testing.T) {
	task := Task{
		Name: "Review case details",
		CaseItems: []Case{
			{
				CaseType: "LPA",
				UID:      "7000-0420-0130",
			},
		},
	}

	assert.Equal(t, task.Summary(), "LPA 7000-0420-0130: Review case details")
}

func TestTaskSummaryMultipleCases(t *testing.T) {
	task := Task{
		Name: "Review case details",
		CaseItems: []Case{
			{
				CaseType: "LPA",
				UID:      "7000-0420-0130",
			},
			{
				CaseType: "LPA",
				UID:      "7000-2839-1010",
			},
		},
	}

	assert.Equal(t, task.Summary(), "LPA 7000-0420-0130: Review case details")
}

func TestTaskSummaryNoCase(t *testing.T) {
	task := Task{
		Name: "Review case details",
	}

	assert.Equal(t, task.Summary(), "Review case details")
}
