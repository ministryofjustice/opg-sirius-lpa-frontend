package server

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockAssignTaskClient struct {
	mock.Mock
}

func (m *mockAssignTaskClient) AssignTask(ctx sirius.Context, assigneeID, taskID int) error {
	args := m.Called(ctx, assigneeID, taskID)
	return args.Error(0)
}

func (m *mockAssignTaskClient) Teams(ctx sirius.Context) ([]sirius.Team, error) {
	args := m.Called(ctx)
	return args.Get(0).([]sirius.Team), args.Error(1)
}

func (m *mockAssignTaskClient) Task(ctx sirius.Context, id int) (sirius.Task, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(sirius.Task), args.Error(1)
}

func TestGetAssignTask(t *testing.T) {
	client := &mockAssignTaskClient{}
	client.
		On("Teams", mock.Anything).
		Return([]sirius.Team{{ID: 1, DisplayName: "A Team"}}, nil)
	client.
		On("Task", mock.Anything, 123).
		Return(sirius.Task{Name: "A task", CaseItems: []sirius.Case{{UID: "7000-0000-0000", CaseType: "LPA"}}}, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, assignTaskData{
			Teams:  []sirius.Team{{ID: 1, DisplayName: "A Team"}},
			Entity: "LPA 7000-0000-0000: A task",
		}).
		Return(nil)

	r, _ := http.NewRequest(http.MethodGet, "/?id=123", nil)
	w := httptest.NewRecorder()

	err := AssignTask(client, template.Func)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestGetAssignTaskBadQueryString(t *testing.T) {
	testCases := map[string]string{
		"no-id":  "/",
		"bad-id": "/?id=what",
	}

	for name, url := range testCases {
		t.Run(name, func(t *testing.T) {
			r, _ := http.NewRequest(http.MethodGet, url, nil)
			w := httptest.NewRecorder()

			err := AssignTask(nil, nil)(w, r)

			assert.NotNil(t, err)
		})
	}
}

func TestGetAssignTaskWhenTeamsErrors(t *testing.T) {
	expectedError := errors.New("hmm")

	client := &mockAssignTaskClient{}
	client.
		On("Teams", mock.Anything).
		Return([]sirius.Team{}, expectedError)
	client.
		On("Task", mock.Anything, mock.Anything).
		Return(sirius.Task{Name: "A task", CaseItems: []sirius.Case{{UID: "7000-0000-0000", CaseType: "LPA"}}}, nil)

	r, _ := http.NewRequest(http.MethodGet, "/?id=123", nil)
	w := httptest.NewRecorder()

	err := AssignTask(client, nil)(w, r)

	assert.Equal(t, expectedError, err)
	mock.AssertExpectationsForObjects(t, client)
}

func TestGetAssignTaskWhenTaskErrors(t *testing.T) {
	expectedError := errors.New("hmm")

	client := &mockAssignTaskClient{}
	client.
		On("Teams", mock.Anything).
		Return([]sirius.Team{}, nil)
	client.
		On("Task", mock.Anything, mock.Anything).
		Return(sirius.Task{Name: "A task", CaseItems: []sirius.Case{{UID: "7000-0000-0000", CaseType: "LPA"}}}, expectedError)

	r, _ := http.NewRequest(http.MethodGet, "/?id=123", nil)
	w := httptest.NewRecorder()

	err := AssignTask(client, nil)(w, r)

	assert.Equal(t, expectedError, err)
	mock.AssertExpectationsForObjects(t, client)
}

func TestGetAssignTaskWhenTemplateErrors(t *testing.T) {
	expectedError := errors.New("hmm")

	client := &mockAssignTaskClient{}
	client.
		On("Teams", mock.Anything).
		Return([]sirius.Team{}, nil)
	client.
		On("Task", mock.Anything, mock.Anything).
		Return(sirius.Task{Name: "A task", CaseItems: []sirius.Case{{UID: "7000-0000-0000", CaseType: "LPA"}}}, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, mock.Anything).
		Return(expectedError)

	r, _ := http.NewRequest(http.MethodGet, "/?id=123", nil)
	w := httptest.NewRecorder()

	err := AssignTask(client, template.Func)(w, r)

	assert.Equal(t, expectedError, err)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestPostAssignTask(t *testing.T) {
	client := &mockAssignTaskClient{}
	client.
		On("Teams", mock.Anything).
		Return([]sirius.Team{{ID: 1, DisplayName: "A Team"}}, nil)
	client.
		On("Task", mock.Anything, 123).
		Return(sirius.Task{Name: "A task", CaseItems: []sirius.Case{{UID: "7000-0000-0000", CaseType: "LPA"}}}, nil)
	client.
		On("AssignTask", mock.Anything, 66, 123).
		Return(nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, assignTaskData{
			Success:          true,
			Teams:            []sirius.Team{{ID: 1, DisplayName: "A Team"}},
			AssigneeUserName: "System user",
			Entity:           "LPA 7000-0000-0000: A task",
		}).
		Return(nil)

	form := url.Values{
		"assignTo":     {"user"},
		"assigneeUser": {"66:System user"},
	}

	r, _ := http.NewRequest(http.MethodPost, "/?id=123", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := AssignTask(client, template.Func)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestPostAssignTaskWhenAssignTaskFails(t *testing.T) {
	expectedError := errors.New("hmm")

	client := &mockAssignTaskClient{}
	client.
		On("Teams", mock.Anything).
		Return([]sirius.Team{{ID: 1, DisplayName: "A Team"}}, nil)
	client.
		On("Task", mock.Anything, 123).
		Return(sirius.Task{Name: "A task", CaseItems: []sirius.Case{{UID: "7000-0000-0000", CaseType: "LPA"}}}, nil)
	client.
		On("AssignTask", mock.Anything, 66, 123).
		Return(expectedError)

	form := url.Values{
		"assignTo":     {"user"},
		"assigneeUser": {"66:System user"},
	}

	r, _ := http.NewRequest(http.MethodPost, "/?id=123", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := AssignTask(client, nil)(w, r)

	assert.Equal(t, expectedError, err)
	mock.AssertExpectationsForObjects(t, client)
}

func TestPostAssignTaskWhenAssignToNotSet(t *testing.T) {
	client := &mockAssignTaskClient{}
	client.
		On("Teams", mock.Anything).
		Return([]sirius.Team{{ID: 1, DisplayName: "A Team"}}, nil)
	client.
		On("Task", mock.Anything, 123).
		Return(sirius.Task{Name: "A task", CaseItems: []sirius.Case{{UID: "7000-0000-0000", CaseType: "LPA"}}}, nil)
	client.
		On("AssignTask", mock.Anything, mock.Anything, mock.Anything).
		Return(sirius.ValidationError{
			Field: sirius.FieldErrors{
				"assigneeId": {"empty": "Not set"},
			},
		})

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, assignTaskData{
			Teams:  []sirius.Team{{ID: 1, DisplayName: "A Team"}},
			Entity: "LPA 7000-0000-0000: A task",
			Error: sirius.ValidationError{
				Field: sirius.FieldErrors{
					"assignTo": {"": "Assignee not set"},
				},
			},
		}).
		Return(nil)

	form := url.Values{
		"assigneeUser": {"66"},
	}

	r, _ := http.NewRequest(http.MethodPost, "/?id=123", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := AssignTask(client, template.Func)(w, r)
	assert.Nil(t, err)

	resp := w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestPostAssignTaskWhenValidationError(t *testing.T) {
	testCases := map[string]struct {
		field            string
		value            string
		assigneeUserName string
	}{
		"team": {
			field: "assigneeTeam",
			value: "66",
		},
		"user": {
			field:            "assigneeUser",
			value:            "66:Some user",
			assigneeUserName: "Some user",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			client := &mockAssignTaskClient{}
			client.
				On("Teams", mock.Anything).
				Return([]sirius.Team{{ID: 1, DisplayName: "A Team"}}, nil)
			client.
				On("Task", mock.Anything, 123).
				Return(sirius.Task{Name: "A task", CaseItems: []sirius.Case{{UID: "7000-0000-0000", CaseType: "LPA"}}}, nil)
			client.
				On("AssignTask", mock.Anything, mock.Anything, mock.Anything).
				Return(sirius.ValidationError{Field: sirius.FieldErrors{
					"field":      {"reason": "Description"},
					"assigneeId": {"problem": "Because"},
				}})

			expectedErrors := sirius.FieldErrors{
				"field": {"reason": "Description"},
			}
			expectedErrors[tc.field] = map[string]string{"problem": "Because"}

			template := &mockTemplate{}
			template.
				On("Func", mock.Anything, assignTaskData{
					AssignTo:         name,
					Teams:            []sirius.Team{{ID: 1, DisplayName: "A Team"}},
					Entity:           "LPA 7000-0000-0000: A task",
					Error:            sirius.ValidationError{Field: expectedErrors},
					AssigneeUserName: tc.assigneeUserName,
				}).
				Return(nil)

			form := url.Values{
				"assignTo": {name},
			}
			form.Add(tc.field, tc.value)

			r, _ := http.NewRequest(http.MethodPost, "/?id=123", strings.NewReader(form.Encode()))
			r.Header.Add("Content-Type", formUrlEncoded)
			w := httptest.NewRecorder()

			err := AssignTask(client, template.Func)(w, r)
			assert.Nil(t, err)

			resp := w.Result()
			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		})
	}
}
