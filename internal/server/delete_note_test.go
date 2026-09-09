package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockDeleteNoteClient struct {
	mock.Mock
}

func (m *mockDeleteNoteClient) DeleteNote(ctx sirius.Context, id int) error {
	return m.Called(ctx, id).Error(0)
}

func (m *mockDeleteNoteClient) GetEvents(ctx sirius.Context, donorId string, caseIds []string, sourceTypes []string, eventIds []string, sortBy string) (sirius.LpaEventsResponse, error) {
	args := m.Called(ctx, donorId, caseIds, sourceTypes, eventIds, sortBy)
	return args.Get(0).(sirius.LpaEventsResponse), args.Error(1)
}

func TestGetDeleteNote(t *testing.T) {
	event := sirius.LpaEvent{ID: 456, SourceNote: sirius.SourceNote{ID: 789}}

	client := &mockDeleteNoteClient{}
	client.
		On("GetEvents", mock.Anything, "123", []string{}, []string{}, []string{"456"}, "desc").
		Return(sirius.LpaEventsResponse{
			Events: []sirius.LpaEvent{event},
		}, nil)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, deleteNoteData{
			DonorId: "123",
			Event:   event,
		}).
		Return(nil)

	r, _ := http.NewRequest(http.MethodGet, "/?donorId=123&eventId=456", nil)
	w := httptest.NewRecorder()

	err := DeleteNote(client, template.Func)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestGetDeleteNoteGetEventsFails(t *testing.T) {
	client := &mockDeleteNoteClient{}
	client.
		On("GetEvents", mock.Anything, "123", []string{}, []string{}, []string{"456"}, "desc").
		Return(sirius.LpaEventsResponse{}, errExample)

	template := &mockTemplate{}

	r, _ := http.NewRequest(http.MethodGet, "/?donorId=123&eventId=456", nil)
	w := httptest.NewRecorder()

	err := DeleteNote(client, template.Func)(w, r)

	assert.Equal(t, errExample, err)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestPostDeleteNote(t *testing.T) {
	event := sirius.LpaEvent{ID: 456, SourceNote: sirius.SourceNote{ID: 789}}

	client := &mockDeleteNoteClient{}
	client.
		On("GetEvents", mock.Anything, "123", []string{}, []string{}, []string{"456"}, "desc").
		Return(sirius.LpaEventsResponse{
			Events: []sirius.LpaEvent{event},
		}, nil)
	client.
		On("DeleteNote", mock.Anything, 789).
		Return(nil)

	template := &mockTemplate{}

	r, _ := http.NewRequest(http.MethodPost, "/?donorId=123&eventId=456", nil)
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := DeleteNote(client, template.Func)(w, r)
	resp := w.Result()

	assert.Equal(t, RedirectError("/donor/123/history"), err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestPostDeleteNoteError(t *testing.T) {
	event := sirius.LpaEvent{ID: 456, SourceNote: sirius.SourceNote{ID: 789}}

	client := &mockDeleteNoteClient{}
	client.
		On("GetEvents", mock.Anything, "123", []string{}, []string{}, []string{"456"}, "desc").
		Return(sirius.LpaEventsResponse{
			Events: []sirius.LpaEvent{event},
		}, nil)
	client.
		On("DeleteNote", mock.Anything, 789).
		Return(errExample)

	template := &mockTemplate{}

	r, _ := http.NewRequest(http.MethodPost, "/?donorId=123&eventId=456", nil)
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := DeleteNote(client, template.Func)(w, r)

	assert.Equal(t, errExample, err)
	mock.AssertExpectationsForObjects(t, client, template)
}
