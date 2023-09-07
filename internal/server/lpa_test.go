package server

import (
	"errors"
	"net/http"
	"testing"

	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockLpaClient struct {
	mock.Mock
}

func (m *mockLpaClient) DigitalLpa(ctx sirius.Context, uid string) (sirius.DigitalLpa, error) {
	args := m.Called(ctx, uid)
	return args.Get(0).(sirius.DigitalLpa), args.Error(1)
}

func (m *mockLpaClient) TasksForCase(ctx sirius.Context, id int) ([]sirius.Task, error) {
	args := m.Called(ctx, id)
	arg := args.Get(0)

	if arg == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]sirius.Task), args.Error(1)
}

func TestGetLpaErrorRetrievingTasksForCase(t *testing.T) {
	digitalLpa := sirius.DigitalLpa{
		ID:      88,
		UID:     "M-AAAA-9876-9876",
		Subtype: "hw",
	}

	client := &mockLpaClient{}
	client.
		On("DigitalLpa", mock.Anything, "M-AAAA-9876-9876").
		Return(digitalLpa, nil)
	client.
		On("TasksForCase", mock.Anything, 88).
		Return(nil, errors.New("foo"))

	template := &mockTemplate{}
	/*
	template.
		On("Func", mock.Anything, lpaData{
			Lpa: digitalLpa,
			TaskList: taskList,
		}).
		Return(nil)
	*/

	server := newMockServer("/lpa/{uid}", Lpa(client, template.Func))

	req, _ := http.NewRequest(http.MethodGet, "/lpa/M-AAAA-9876-9876", nil)
	_, err := server.serve(req)

	assert.Equal(t, "foo", err.Error())
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestGetLpa(t *testing.T) {
	digitalLpa := sirius.DigitalLpa{
		ID:      22,
		UID:     "M-9876-9876-9876",
		Subtype: "hw",
	}

	taskList := []sirius.Task{}

	client := &mockLpaClient{}
	client.
		On("DigitalLpa", mock.Anything, "M-9876-9876-9876").
		Return(digitalLpa, nil)
	client.
		On("TasksForCase", mock.Anything, 22).
		Return(taskList, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, lpaData{
			Lpa: digitalLpa,
			TaskList: taskList,
		}).
		Return(nil)

	server := newMockServer("/lpa/{uid}", Lpa(client, template.Func))

	req, _ := http.NewRequest(http.MethodGet, "/lpa/M-9876-9876-9876", nil)
	_, err := server.serve(req)

	assert.Nil(t, err)
	mock.AssertExpectationsForObjects(t, client, template)
}
