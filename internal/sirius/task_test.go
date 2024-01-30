package sirius

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/pact-foundation/pact-go/dsl"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockTaskHttpClient struct {
	mock.Mock
}

func (m *mockTaskHttpClient) Do(req *http.Request) (*http.Response, error) {
	args := m.Called(req)
	return args.Get(0).(*http.Response), args.Error(1)
}

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

func TestTasksForCaseBadContext(t *testing.T) {
	client := NewClient(http.DefaultClient, "http://localhost")
	_, err := client.TasksForCase(Context{Context: nil}, 21)
	assert.Equal(t, "net/http: nil Context", err.Error())
}

func TestTasksForCaseNetworkError(t *testing.T) {
	mockTaskHttpClient := &mockTaskHttpClient{}

	// NB the structure of the error returned here does not match real http errors
	// returned by Do(), but we only care whether the error is handled
	mockTaskHttpClient.On("Do", mock.Anything).Return(&http.Response{}, errors.New("Networking issue"))

	client := NewClient(mockTaskHttpClient, "http://localhost")
	_, err := client.TasksForCase(Context{Context: context.Background()}, 777)
	assert.Equal(t, "Networking issue", err.Error())
}

func TestTasksForCaseBadJson(t *testing.T) {
	mockTaskHttpClient := &mockTaskHttpClient{}

	badJsonResponse := http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewBufferString("[bad json time")),
	}

	mockTaskHttpClient.On("Do", mock.Anything).Return(&badJsonResponse, nil)

	client := NewClient(mockTaskHttpClient, "http://localhost")
	_, err := client.TasksForCase(Context{Context: context.Background()}, 8888)

	assert.Equal(t, "invalid character 'b' looking for beginning of value", err.Error())
}

func TestTasksForCase(t *testing.T) {
	t.Parallel()

	pact := newPact()
	defer pact.Teardown()

	testCases := []struct {
		name             string
		id               int
		setup            func()
		expectedResponse []Task
		expectedError    func(int) error
	}{
		{
			name: "OK",
			id:   10,
			setup: func() {
				pact.
					AddInteraction().
					Given("I have a case with an open task assigned").
					UponReceiving("A request for the tasks for a case").
					WithRequest(dsl.Request{
						Method: http.MethodGet,
						Path:   dsl.Like("/lpa-api/v1/cases/10/tasks"),
					}).
					WillRespondWith(dsl.Response{
						Status: http.StatusOK,
						Body: dsl.Like(map[string]interface{}{
							"tasks": dsl.EachLike(map[string]interface{}{
								"id":      dsl.Like(12),
								"status":  dsl.String("Not started"),
								"dueDate": dsl.String("05/09/2023"),
								"name":    dsl.String("Review reduced fees request"),
								"assignee": dsl.Like(map[string]interface{}{
									"displayName": dsl.String("Consuela"),
								}),
							}, 1),
						}),
						Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/json")},
					})
			},
			expectedResponse: []Task{
				{
					ID:      12,
					Status:  "Not started",
					DueDate: DateString("2023-09-05"),
					Name:    "Review reduced fees request",
					Assignee: User{
						DisplayName: "Consuela",
					},
				},
			},
		},
		{
			name: "404",
			id:   9012929,
			setup: func() {
				pact.
					AddInteraction().
					Given("There is no case with tasks with the specified ID").
					UponReceiving("A request for the tasks for a non-existent case").
					WithRequest(dsl.Request{
						Method: http.MethodGet,
						Path:   dsl.Like("/lpa-api/v1/cases/9012929/tasks"),
					}).
					WillRespondWith(dsl.Response{
						Status: http.StatusNotFound,
					})
			},
			expectedError: func(port int) error {
				return StatusError{
					Code:          404,
					URL:           fmt.Sprintf("http://localhost:%d/lpa-api/v1/cases/9012929/tasks?filter=status%%3ANot+started%%2Cactive%%3Atrue&limit=99&sort=duedate%%3AASC", port),
					Method:        "GET",
					CorrelationId: "",
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.Verify(func() error {
				client := NewClient(http.DefaultClient, fmt.Sprintf("http://localhost:%d", pact.Server.Port))

				tasks, err := client.TasksForCase(Context{Context: context.Background()}, tc.id)

				assert.Equal(t, tc.expectedResponse, tasks)
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
