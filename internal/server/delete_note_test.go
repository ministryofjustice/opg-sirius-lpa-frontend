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

func TestGetDeleteNote(t *testing.T) {
	client := &mockDeleteNoteClient{}
	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, deleteNoteData{
			DonorId: 123,
		}).
		Return(nil)

	r, _ := http.NewRequest(http.MethodGet, "/?donorId=123&noteId=456", nil)
	w := httptest.NewRecorder()

	err := DeleteNote(client, template.Func)(w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestGetDeleteNoteBadQuery(t *testing.T) {
	testCases := map[string]string{
		"no-donorId":  "/",
		"bad-donorId": "/?donorId=test",
		"no-noteId":   "/?donorId=123",
		"bad-noteId":  "/?donorId=123&noteId=test",
	}

	for name, url := range testCases {
		t.Run(name, func(t *testing.T) {
			r, _ := http.NewRequest(http.MethodGet, url, nil)
			w := httptest.NewRecorder()

			err := DeleteNote(nil, nil)(w, r)

			assert.NotNil(t, err)
		})
	}
}

func TestPostDeleteNote(t *testing.T) {
	client := &mockDeleteNoteClient{}
	client.
		On("DeleteNote", mock.Anything, 456).
		Return(nil)

	template := &mockTemplate{}

	r, _ := http.NewRequest(http.MethodPost, "/?donorId=123&noteId=456", nil)
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := DeleteNote(client, template.Func)(w, r)
	resp := w.Result()

	assert.Equal(t, RedirectError("/donor/123/history"), err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, client, template)
}

func TestPostDeleteNoteError(t *testing.T) {
	client := &mockDeleteNoteClient{}
	client.
		On("DeleteNote", mock.Anything, 456).
		Return(errExample)

	template := &mockTemplate{}

	r, _ := http.NewRequest(http.MethodPost, "/?donorId=123&noteId=456", nil)
	r.Header.Add("Content-Type", formUrlEncoded)
	w := httptest.NewRecorder()

	err := DeleteNote(client, template.Func)(w, r)

	assert.Equal(t, errExample, err)
	mock.AssertExpectationsForObjects(t, client, template)
}
