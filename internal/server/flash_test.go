package server

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockResponseWriter struct {
	mock.Mock
}

func (m *mockResponseWriter) Header() http.Header {
	args := m.Called()
	return args.Get(0).(http.Header)
}

func (m *mockResponseWriter) Write([]byte) (int, error) {
	args := m.Called()
	return args.Get(0).(int), args.Error(1)
}
func (m *mockResponseWriter) WriteHeader(statusCode int) {
	m.Called(statusCode)
}

func TestSetFlashAddsHeader(t *testing.T) {
	header := http.Header{}
	w := &mockResponseWriter{}
	w.
		On("Header").
		Return(header)

	SetFlash(w, FlashNotification{
		Title:       "title",
		Description: "description",
	})

	assert.Equal(t, "flash-lpa-frontend=eyJuYW1lIjoidGl0bGUiLCJkZXNjcmlwdGlvbiI6ImRlc2NyaXB0aW9uIn0=; HttpOnly", header.Get("Set-Cookie"))
}

func TestGetFlashGetsHeader(t *testing.T) {
	header := http.Header{}
	w := &mockResponseWriter{}
	w.
		On("Header").
		Return(header)

	r, _ := http.NewRequest("GET", "/some-url", nil)
	r.AddCookie(&http.Cookie{
		Name:  "flash-lpa-frontend",
		Value: "eyJuYW1lIjoidGl0bGUiLCJkZXNjcmlwdGlvbiI6ImRlc2NyaXB0aW9uIn0=",
	})

	notification, err := GetFlash(w, r)

	assert.Nil(t, err)
	assert.Equal(t, "title", notification.Title)
	assert.Equal(t, "description", notification.Description)
	assert.Contains(t, header.Get("Set-Cookie"), "flash-lpa-frontend=;")
}

func TestGetFlashReturnsEmptyIfNoCookie(t *testing.T) {
	header := http.Header{}
	w := &mockResponseWriter{}
	w.
		On("Header").
		Return(header)

	r, _ := http.NewRequest("GET", "/some-url", nil)

	notification, err := GetFlash(w, r)

	assert.Nil(t, err)
	assert.Equal(t, FlashNotification{}, notification)
}

func TestGetFlashReturnsErrorIfCannotDecodeBase64(t *testing.T) {
	header := http.Header{}
	w := &mockResponseWriter{}
	w.
		On("Header").
		Return(header)

	r, _ := http.NewRequest("GET", "/some-url", nil)
	r.AddCookie(&http.Cookie{
		Name:  "flash-lpa-frontend",
		Value: "badstring",
	})

	notification, err := GetFlash(w, r)

	assert.NotNil(t, err)
	assert.Equal(t, FlashNotification{}, notification)
}

func TestGetFlashReturnsErrorIfCannotDecodeJSON(t *testing.T) {
	header := http.Header{}
	w := &mockResponseWriter{}
	w.
		On("Header").
		Return(header)

	r, _ := http.NewRequest("GET", "/some-url", nil)
	r.AddCookie(&http.Cookie{
		Name:  "flash-lpa-frontend",
		Value: "dGhpcyBpcyBub3QgSlNPTg==",
	})

	notification, err := GetFlash(w, r)

	assert.NotNil(t, err)
	assert.Equal(t, FlashNotification{}, notification)
}
