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

func TestClearTask(t *testing.T) {
	client := &mockClearTaskClient{}
	client.
		On("Task", mock.Anything, 33).
		Return(sirius.Task{Name: "A task", CaseItems: []sirius.Case{{UID: "M-DIGI-0001-0001", CaseType: "DIGITAL_LPA"}}}, nil)

	template := &mockTemplate{}
	template.On("Func", mock.Anything, clearTaskData{
		Task: sirius.Task{Name: "A task", CaseItems: []sirius.Case{{UID: "M-DIGI-0001-0001", CaseType: "DIGITAL_LPA"}}},
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
	client := &mockClearTaskClient{}
	client.
		On("Task", mock.Anything, 33).
		Return(sirius.Task{ID: 33, Name: "A task", CaseItems: []sirius.Case{{UID: "M-DIGI-0001-0001", CaseType: "DIGITAL_LPA"}}}, nil)
	client.
		On("ClearTask", mock.Anything, 33).
		Return(nil)

	template := &mockTemplate{}

	form := url.Values{
		"id": {"33"},
	}

	r, _ := http.NewRequest(http.MethodPost, "/clear-task?id=33", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := ClearTask(client, template.Func)(w, r)
	resp := w.Result()

	assert.Equal(t, RedirectError("/lpa/M-DIGI-0001-0001"), err)
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
	client := &mockClearTaskClient{}
	client.
		On("Task", mock.Anything, mock.Anything).
		Return(sirius.Task{Name: "A task", CaseItems: []sirius.Case{{UID: "7000-0000-0000", CaseType: "LPA"}}}, expectedError)

	r, _ := http.NewRequest(http.MethodGet, "/?id=123", nil)
	w := httptest.NewRecorder()

	err := ClearTask(client, nil)(w, r)

	assert.Equal(t, expectedError, err)
	mock.AssertExpectationsForObjects(t, client)
}

func TestClearTaskWhenTemplateErrors(t *testing.T) {
	client := &mockClearTaskClient{}
	client.
		On("Task", mock.Anything, mock.Anything).
		Return(sirius.Task{Name: "A task", CaseItems: []sirius.Case{{UID: "7000-0000-0000", CaseType: "LPA"}}}, nil)

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
