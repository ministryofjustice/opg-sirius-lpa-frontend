package server

import (
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

type mockClearTaskClient struct {
	mock.Mock
}

func (m *mockClearTaskClient) ClearTask(ctx sirius.Context, taskID int) error {
	args := m.Called(ctx, taskID)
	return args.Error(0)
}

func (m *mockClearTaskClient) Task(ctx sirius.Context, id int) (sirius.Task, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(sirius.Task), args.Error(1)
}

func setupMockClearTaskClient(taskID int, taskUID string, CaseType string, expectedError error) *mockClearTaskClient {
	client := &mockClearTaskClient{}
	client.
		On("Task", mock.Anything, taskID).
		Return(sirius.Task{ID: taskID, Name: "A task", CaseItems: []sirius.Case{{UID: taskUID, CaseType: CaseType}}}, expectedError)
	return client
}

func TestClearTask(t *testing.T) {
	client := setupMockClearTaskClient(33, "M-DIGI-0001-0001", "DIGITAL_LPA", nil)

	template := &mockTemplate{}
	template.On("Func", mock.Anything, clearTaskData{
		Task: sirius.Task{ID: 33, Name: "A task", CaseItems: []sirius.Case{{UID: "M-DIGI-0001-0001", CaseType: "DIGITAL_LPA"}}},
		Uid:  "M-DIGI-0001-0001",
	}).Return(nil)

	r, _ := http.NewRequest(http.MethodGet, "/clear-task?id=33", nil)
	w := httptest.NewRecorder()

	err := ClearTask(client, template.Func)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestPostClearTaskRedirects(t *testing.T) {
	client := setupMockClearTaskClient(33, "M-DIGI-0001-0001", "DIGITAL_LPA", nil)
	client.On("ClearTask", mock.Anything, 33).Return(nil)

	template := &mockTemplate{}

	form := url.Values{"id": {"33"}}

	r, _ := http.NewRequest(http.MethodPost, "/clear-task?id=33", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := ClearTask(client, template.Func)(w, r)
	resp := w.Result()

	assert.Equal(t, RedirectError("/lpa/M-DIGI-0001-0001"), err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestPostClearTaskSuccess(t *testing.T) {
	client := setupMockClearTaskClient(66, "7000-0000-0000", "LPA", nil)
	client.On("ClearTask", mock.Anything, 66).Return(nil)

	template := &mockTemplate{}

	template.On("Func", mock.Anything, clearTaskData{
		Success: true,
		Task:    sirius.Task{ID: 66, Name: "A task", CaseItems: []sirius.Case{{UID: "7000-0000-0000", CaseType: "LPA"}}},
		Uid:     "7000-0000-0000",
	}).Return(nil)

	form := url.Values{"id": {"66"}}

	r, _ := http.NewRequest(http.MethodPost, "/clear-task?id=66", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := ClearTask(client, template.Func)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestClearTaskBadQueryString(t *testing.T) {
	testCases := map[string]string{
		"no-id":      "/",
		"bad-id":     "/?id=what",
		"one-bad-id": "/?id=1&id=bad&id=2",
	}

	for name, url := range testCases {
		t.Run(name, func(t *testing.T) {
			r, _ := http.NewRequest(http.MethodGet, url, nil)
			w := httptest.NewRecorder()

			err := ClearTask(nil, nil)(w, r)

			assert.NotNil(t, err)
		})
	}
}

func TestClearTaskWhenTaskErrors(t *testing.T) {
	client := setupMockClearTaskClient(123, "7000-0000-0000", "LPA", expectedError)

	r, _ := http.NewRequest(http.MethodGet, "/?id=123", nil)
	w := httptest.NewRecorder()

	err := ClearTask(client, nil)(w, r)

	assert.Equal(t, expectedError, err)
	mock.AssertExpectationsForObjects(t, client)
}

func TestClearTaskWhenTemplateErrors(t *testing.T) {
	client := setupMockClearTaskClient(123, "7000-0000-0000", "LPA", nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, mock.Anything).
		Return(expectedError)

	r, _ := http.NewRequest(http.MethodGet, "/?id=123", nil)
	w := httptest.NewRecorder()

	err := ClearTask(client, template.Func)(w, r)

	assert.Equal(t, expectedError, err)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestPostClearTaskWhenValidationErrors(t *testing.T) {
	expectedError := sirius.ValidationError{
		Field: sirius.FieldErrors{"field": {"": "problem"}},
	}

	client := setupMockClearTaskClient(35, "M-DIGI-0001-0002", "DIGITAL_LPA", nil)
	client.On("ClearTask", mock.Anything, 35).Return(expectedError)

	template := &mockTemplate{}

	template.On("Func", mock.Anything, clearTaskData{
		Success: false,
		Error:   expectedError,
		Task:    sirius.Task{ID: 35, Name: "A task", CaseItems: []sirius.Case{{UID: "M-DIGI-0001-0002", CaseType: "DIGITAL_LPA"}}},
		Uid:     "M-DIGI-0001-0002",
	}).Return(nil)

	form := url.Values{"id": {"35"}}

	r, _ := http.NewRequest(http.MethodPost, "/clear-task?id=35", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := ClearTask(client, template.Func)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestPostClearTaskWhenOtherError(t *testing.T) {
	client := setupMockClearTaskClient(33, "M-DIGI-0001-0001", "DIGITAL_LPA", nil)
	client.On("ClearTask", mock.Anything, 33).Return(expectedError)

	form := url.Values{"id": {"33"}}

	r, _ := http.NewRequest(http.MethodPost, "/clear-task?id=33", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := ClearTask(client, nil)(w, r)

	assert.Equal(t, expectedError, err)
	mock.AssertExpectationsForObjects(t, client)
}
