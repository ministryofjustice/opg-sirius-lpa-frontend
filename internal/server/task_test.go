package server

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockTaskClient struct {
	mock.Mock
}

func (m *mockTaskClient) CreateTask(ctx sirius.Context, caseID int, task sirius.TaskRequest) error {
	args := m.Called(ctx, caseID, task)
	return args.Error(0)
}

func (m *mockTaskClient) Teams(ctx sirius.Context) ([]sirius.Team, error) {
	args := m.Called(ctx)
	return args.Get(0).([]sirius.Team), args.Error(1)
}

func (m *mockTaskClient) TaskTypes(ctx sirius.Context) ([]string, error) {
	args := m.Called(ctx)
	return args.Get(0).([]string), args.Error(1)
}

func (m *mockTaskClient) Case(ctx sirius.Context, id int) (sirius.Case, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(sirius.Case), args.Error(1)
}

func (m *mockTaskClient) GetUserDetails(ctx sirius.Context) (sirius.User, error) {
	args := m.Called(ctx)
	return args.Get(0).(sirius.User), args.Error(1)
}

func TestGetTask(t *testing.T) {
	client := &mockTaskClient{}
	client.
		On("TaskTypes", mock.Anything).
		Return([]string{"a", "b"}, nil)
	client.
		On("Teams", mock.Anything).
		Return([]sirius.Team{{ID: 1, DisplayName: "A Team"}}, nil)
	client.
		On("Case", mock.Anything, 123).
		Return(sirius.Case{UID: "7000-0000-0000", CaseType: "LPA"}, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, taskData{
			TaskTypes: []string{"a", "b"},
			Teams:     []sirius.Team{{ID: 1, DisplayName: "A Team"}},
			Entity:    "LPA 7000-0000-0000",
			CaseUID:   "7000-0000-0000",
		}).
		Return(nil)

	r, _ := http.NewRequest(http.MethodGet, "/?id=123", nil)
	w := httptest.NewRecorder()

	err := Task(client, template.Func)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestGetTaskBadQueryString(t *testing.T) {
	testCases := map[string]string{
		"no-id":  "/",
		"bad-id": "/?id=what",
	}

	for name, url := range testCases {
		t.Run(name, func(t *testing.T) {
			r, _ := http.NewRequest(http.MethodGet, url, nil)
			w := httptest.NewRecorder()

			err := Task(nil, nil)(w, r)

			assert.NotNil(t, err)
		})
	}
}

func TestGetTaskWhenTaskTypeErrors(t *testing.T) {
	client := &mockTaskClient{}
	client.
		On("TaskTypes", mock.Anything).
		Return([]string{}, errExample)
	client.
		On("Teams", mock.Anything).
		Return([]sirius.Team{}, nil)
	client.
		On("Case", mock.Anything, mock.Anything).
		Return(sirius.Case{}, nil)

	r, _ := http.NewRequest(http.MethodGet, "/?id=123", nil)
	w := httptest.NewRecorder()

	err := Task(client, nil)(w, r)

	assert.Equal(t, errExample, err)
}

func TestGetTaskWhenTeamsErrors(t *testing.T) {
	client := &mockTaskClient{}
	client.
		On("TaskTypes", mock.Anything).
		Return([]string{}, nil)
	client.
		On("Teams", mock.Anything).
		Return([]sirius.Team{}, errExample)
	client.
		On("Case", mock.Anything, mock.Anything).
		Return(sirius.Case{}, nil)

	r, _ := http.NewRequest(http.MethodGet, "/?id=123", nil)
	w := httptest.NewRecorder()

	err := Task(client, nil)(w, r)

	assert.Equal(t, errExample, err)
}

func TestGetTaskWhenCaseErrors(t *testing.T) {
	client := &mockTaskClient{}
	client.
		On("TaskTypes", mock.Anything).
		Return([]string{}, nil)
	client.
		On("Teams", mock.Anything).
		Return([]sirius.Team{}, nil)
	client.
		On("Case", mock.Anything, mock.Anything).
		Return(sirius.Case{}, errExample)

	r, _ := http.NewRequest(http.MethodGet, "/?id=123", nil)
	w := httptest.NewRecorder()

	err := Task(client, nil)(w, r)

	assert.Equal(t, errExample, err)
}

func TestGetTaskWhenTemplateErrors(t *testing.T) {
	client := &mockTaskClient{}
	client.
		On("TaskTypes", mock.Anything).
		Return([]string{}, nil)
	client.
		On("Teams", mock.Anything).
		Return([]sirius.Team{}, nil)
	client.
		On("Case", mock.Anything, mock.Anything).
		Return(sirius.Case{}, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, mock.Anything).
		Return(errExample)

	r, _ := http.NewRequest(http.MethodGet, "/?id=123", nil)
	w := httptest.NewRecorder()

	err := Task(client, template.Func)(w, r)

	assert.Equal(t, errExample, err)
}

func TestPostTask(t *testing.T) {
	client := &mockTaskClient{}
	client.
		On("TaskTypes", mock.Anything).
		Return([]string{"a", "b"}, nil)
	client.
		On("Teams", mock.Anything).
		Return([]sirius.Team{{ID: 1, DisplayName: "A Team"}}, nil)
	client.
		On("Case", mock.Anything, 123).
		Return(sirius.Case{UID: "7000-0000-0000", CaseType: "LPA"}, nil)
	client.
		On("CreateTask", mock.Anything, 123, sirius.TaskRequest{
			Type:        "Some task type",
			DueDate:     "2022-03-04",
			Name:        "Do this",
			Description: "Please",
			AssigneeID:  66,
		}).
		Return(nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, taskData{
			Success:          true,
			TaskTypes:        []string{"a", "b"},
			Teams:            []sirius.Team{{ID: 1, DisplayName: "A Team"}},
			AssigneeUserName: "System user",
			Entity:           "LPA 7000-0000-0000",
			CaseUID:          "7000-0000-0000",
		}).
		Return(nil)

	form := url.Values{
		"assignTo":     {"user"},
		"assigneeUser": {"66:System user"},
		"taskType":     {"Some task type"},
		"dueDate":      {"2022-03-04"},
		"name":         {"Do this"},
		"description":  {"Please"},
	}

	r, _ := http.NewRequest(http.MethodPost, "/?id=123", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := Task(client, template.Func)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestPostTaskWhenAssignToMe(t *testing.T) {
	user := sirius.User{ID: 66, DisplayName: "Me", Roles: []string{"OPG User", "Reduced Fees User"}}

	client := &mockTaskClient{}
	client.
		On("TaskTypes", mock.Anything).
		Return([]string{"a", "b"}, nil)
	client.
		On("Teams", mock.Anything).
		Return([]sirius.Team{{ID: 1, DisplayName: "A Team"}}, nil)
	client.
		On("Case", mock.Anything, 123).
		Return(sirius.Case{UID: "7000-0000-0000", CaseType: "LPA"}, nil)
	client.
		On("CreateTask", mock.Anything, 123, sirius.TaskRequest{
			Type:        "Some task type",
			DueDate:     "2022-03-04",
			Name:        "Do this",
			Description: "Please",
			AssigneeID:  66,
		}).
		Return(nil)
	client.
		On("GetUserDetails", mock.Anything).
		Return(user, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, taskData{
			Success:          true,
			TaskTypes:        []string{"a", "b"},
			Teams:            []sirius.Team{{ID: 1, DisplayName: "A Team"}},
			AssigneeUserName: "Me",
			Entity:           "LPA 7000-0000-0000",
			CaseUID:          "7000-0000-0000",
		}).
		Return(nil)

	form := url.Values{
		"assignTo":    {"me"},
		"taskType":    {"Some task type"},
		"dueDate":     {"2022-03-04"},
		"name":        {"Do this"},
		"description": {"Please"},
	}

	r, _ := http.NewRequest(http.MethodPost, "/?id=123", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := Task(client, template.Func)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestGetTaskWhenGetUserDetailsErrors(t *testing.T) {
	client := &mockTaskClient{}
	client.
		On("TaskTypes", mock.Anything).
		Return([]string{}, nil)
	client.
		On("Teams", mock.Anything).
		Return([]sirius.Team{}, nil)
	client.
		On("Case", mock.Anything, mock.Anything).
		Return(sirius.Case{}, nil)
	client.
		On("CreateTask", mock.Anything, 123, sirius.TaskRequest{
			Type:        "Some task type",
			DueDate:     "2022-03-04",
			Name:        "Do this",
			Description: "Please",
			AssigneeID:  66,
		}).
		Return(nil)
	client.
		On("GetUserDetails", mock.Anything).
		Return(sirius.User{}, errExample)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, taskData{
			Success:          true,
			TaskTypes:        []string{"a", "b"},
			Teams:            []sirius.Team{{ID: 1, DisplayName: "A Team"}},
			AssigneeUserName: "Me",
			Entity:           "LPA 7000-0000-0000",
			CaseUID:          "7000-0000-0000",
		}).
		Return(nil)

	form := url.Values{
		"assignTo":    {"me"},
		"taskType":    {"Some task type"},
		"dueDate":     {"2022-03-04"},
		"name":        {"Do this"},
		"description": {"Please"},
	}

	r, _ := http.NewRequest(http.MethodPost, "/?id=123", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := Task(client, template.Func)(w, r)

	assert.Equal(t, errExample, err)
}

func TestPostTaskWhenCreateTaskFails(t *testing.T) {
	client := &mockTaskClient{}
	client.
		On("TaskTypes", mock.Anything).
		Return([]string{"a", "b"}, nil)
	client.
		On("Teams", mock.Anything).
		Return([]sirius.Team{{ID: 1, DisplayName: "A Team"}}, nil)
	client.
		On("Case", mock.Anything, 123).
		Return(sirius.Case{UID: "7000-0000-0000", CaseType: "LPA"}, nil)
	client.
		On("CreateTask", mock.Anything, mock.Anything, mock.Anything).
		Return(errExample)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, taskData{
			TaskTypes: []string{"a", "b"},
			Teams:     []sirius.Team{{ID: 1, DisplayName: "A Team"}},
			Entity:    "LPA 7000-0000-0000",
		}).
		Return(nil)

	form := url.Values{
		"assignTo":     {"user"},
		"assigneeUser": {"66"},
		"taskType":     {"Some task type"},
		"dueDate":      {"2022-03-04"},
		"name":         {"Do this"},
		"description":  {"Please"},
	}

	r, _ := http.NewRequest(http.MethodPost, "/?id=123", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := Task(client, template.Func)(w, r)
	assert.Equal(t, errExample, err)
}

func TestPostTaskWhenAssignToNotSet(t *testing.T) {
	client := &mockTaskClient{}
	client.
		On("TaskTypes", mock.Anything).
		Return([]string{"a", "b"}, nil)
	client.
		On("Teams", mock.Anything).
		Return([]sirius.Team{{ID: 1, DisplayName: "A Team"}}, nil)
	client.
		On("Case", mock.Anything, 123).
		Return(sirius.Case{UID: "7000-0000-0000", CaseType: "LPA"}, nil)
	client.
		On("CreateTask", mock.Anything, mock.Anything, mock.Anything).
		Return(sirius.ValidationError{
			Field: sirius.FieldErrors{
				"assigneeId": {"empty": "Not set"},
			},
		})

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, taskData{
			TaskTypes: []string{"a", "b"},
			Teams:     []sirius.Team{{ID: 1, DisplayName: "A Team"}},
			Entity:    "LPA 7000-0000-0000",
			CaseUID:   "7000-0000-0000",
			Error: sirius.ValidationError{
				Field: sirius.FieldErrors{
					"assignTo": {"": "Assignee not set"},
				},
			},
			Task: sirius.TaskRequest{
				Type:        "Some task type",
				DueDate:     "2022-03-04",
				Name:        "Do this",
				Description: "Please",
			},
		}).
		Return(nil)

	form := url.Values{
		"assigneeUser": {"66"},
		"taskType":     {"Some task type"},
		"dueDate":      {"2022-03-04"},
		"name":         {"Do this"},
		"description":  {"Please"},
	}

	r, _ := http.NewRequest(http.MethodPost, "/?id=123", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := Task(client, template.Func)(w, r)
	assert.Nil(t, err)

	resp := w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestPostTaskWhenValidationError(t *testing.T) {
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
			client := &mockTaskClient{}
			client.
				On("TaskTypes", mock.Anything).
				Return([]string{"a", "b"}, nil)
			client.
				On("Teams", mock.Anything).
				Return([]sirius.Team{{ID: 1, DisplayName: "A Team"}}, nil)
			client.
				On("Case", mock.Anything, 123).
				Return(sirius.Case{UID: "7000-0000-0000", CaseType: "LPA"}, nil)
			client.
				On("CreateTask", mock.Anything, mock.Anything, mock.Anything).
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
				On("Func", mock.Anything, taskData{
					AssignTo:  name,
					TaskTypes: []string{"a", "b"},
					Teams:     []sirius.Team{{ID: 1, DisplayName: "A Team"}},
					Entity:    "LPA 7000-0000-0000",
					CaseUID:   "7000-0000-0000",
					Error:     sirius.ValidationError{Field: expectedErrors},
					Task: sirius.TaskRequest{
						Type:        "Some task type",
						DueDate:     "2022-03-04",
						Name:        "Do this",
						Description: "Please",
						AssigneeID:  66,
					},
					AssigneeUserName: tc.assigneeUserName,
				}).
				Return(nil)

			form := url.Values{
				"assignTo":    {name},
				"taskType":    {"Some task type"},
				"dueDate":     {"2022-03-04"},
				"name":        {"Do this"},
				"description": {"Please"},
			}
			form.Add(tc.field, tc.value)

			r, _ := http.NewRequest(http.MethodPost, "/?id=123", strings.NewReader(form.Encode()))
			r.Header.Add("Content-Type", formUrlEncoded)
			w := httptest.NewRecorder()

			err := Task(client, template.Func)(w, r)
			assert.Nil(t, err)

			resp := w.Result()
			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		})
	}
}
