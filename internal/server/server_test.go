package server

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockTemplate struct {
	mock.Mock
}

func (t *mockTemplate) Func(w io.Writer, data interface{}) error {
	args := t.Called(w, data)
	return args.Error(0)
}

type mockLogger struct {
	mock.Mock
}

func (t *mockLogger) Request(r *http.Request, err error) {
	t.Called(r, err)
}

type anUnauthorizedError struct {
	is bool
}

func (anUnauthorizedError) Error() string {
	return "hey"
}

func (e anUnauthorizedError) IsUnauthorized() bool {
	return e.is
}

func TestNew(t *testing.T) {
	assert.Implements(t, (*http.Handler)(nil), New(nil, nil, nil, "", "", ""))
}

func TestErrorHandlerError(t *testing.T) {
	assert := assert.New(t)

	expectedErr := anUnauthorizedError{is: false}

	logger := &mockLogger{}
	logger.
		On("Request", mock.Anything, expectedErr)

	template := &mockTemplate{}
	template.
		On("Func", mock.Anything, errorVars{
			SiriusURL: "http://sirius",
			Code:      http.StatusInternalServerError,
			Error:     "hey",
		}).
		Return(nil)

	handler := errorHandler(logger, template.Func, "http://prefix", "http://sirius")(func(w http.ResponseWriter, r *http.Request) error {
		return expectedErr
	})

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)

	handler.ServeHTTP(w, r)

	resp := w.Result()

	assert.Equal(http.StatusInternalServerError, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, template, logger)
}

func TestErrorHandlerUnauthorizedError(t *testing.T) {
	assert := assert.New(t)

	expectedErr := anUnauthorizedError{is: true}

	handler := errorHandler(nil, nil, "http://prefix", "http://sirius")(func(w http.ResponseWriter, r *http.Request) error {
		return expectedErr
	})

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)

	handler.ServeHTTP(w, r)

	resp := w.Result()

	assert.Equal(http.StatusFound, resp.StatusCode)
	assert.Equal("http://sirius/auth", resp.Header.Get("Location"))
}

func TestGetContext(t *testing.T) {
	assert := assert.New(t)

	r, _ := http.NewRequest("GET", "/", nil)
	r.AddCookie(&http.Cookie{Name: "XSRF-TOKEN", Value: "z3tVRZ00yx4dHz3KWYv3boLWHZ4/RsCsVAKbvo2SBNc%3D"})
	r.AddCookie(&http.Cookie{Name: "another", Value: "one"})

	ctx := getContext(r)
	assert.Equal(r.Context(), ctx.Context)
	assert.Equal(r.Cookies(), ctx.Cookies)
	assert.Equal("z3tVRZ00yx4dHz3KWYv3boLWHZ4/RsCsVAKbvo2SBNc=", ctx.XSRFToken)
}

func TestGetContextBadXSRFToken(t *testing.T) {
	assert := assert.New(t)

	r, _ := http.NewRequest("GET", "/", nil)
	r.AddCookie(&http.Cookie{Name: "XSRF-TOKEN", Value: "%"})
	r.AddCookie(&http.Cookie{Name: "another", Value: "one"})

	ctx := getContext(r)
	assert.Equal(r.Context(), ctx.Context)
	assert.Equal(r.Cookies(), ctx.Cookies)
	assert.Equal("", ctx.XSRFToken)
}

func TestGetContextMissingXSRFToken(t *testing.T) {
	assert := assert.New(t)

	r, _ := http.NewRequest("GET", "/", nil)
	r.AddCookie(&http.Cookie{Name: "another", Value: "one"})

	ctx := getContext(r)
	assert.Equal(r.Context(), ctx.Context)
	assert.Equal(r.Cookies(), ctx.Cookies)
	assert.Equal("", ctx.XSRFToken)
}

func TestGetContextForPostRequest(t *testing.T) {
	assert := assert.New(t)

	r, _ := http.NewRequest("POST", "/", strings.NewReader("xsrfToken=the-real-one"))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.AddCookie(&http.Cookie{Name: "XSRF-TOKEN", Value: "z3tVRZ00yx4dHz3KWYv3boLWHZ4/RsCsVAKbvo2SBNc%3D"})
	r.AddCookie(&http.Cookie{Name: "another", Value: "one"})

	ctx := getContext(r)
	assert.Equal(r.Context(), ctx.Context)
	assert.Equal(r.Cookies(), ctx.Cookies)
	assert.Equal("the-real-one", ctx.XSRFToken)
}
