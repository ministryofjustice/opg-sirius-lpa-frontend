package server

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const formUrlEncoded = "application/x-www-form-urlencoded"

var expectedError = errors.New("error")

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

	handler := errorHandler(nil, nil, "/prefix", "http://sirius")(func(w http.ResponseWriter, r *http.Request) error {
		return expectedErr
	})

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)

	handler.ServeHTTP(w, r)

	resp := w.Result()

	assert.Equal(http.StatusFound, resp.StatusCode)
	assert.Equal("http://sirius/auth?redirect=%2Fprefix%2Fpath", resp.Header.Get("Location"))
}

func TestErrorHandlerJsonError(t *testing.T) {
	assert := assert.New(t)

	expectedError := sirius.ValidationError{
		Detail: "Not valid complaint",
		Field: sirius.FieldErrors{
			"title": {
				"tooShort": "The title must be at least 5 characters",
			},
		},
	}

	logger := &mockLogger{}
	logger.
		On("Request", mock.Anything, expectedError)

	handler := errorHandler(logger, nil, "http://prefix", "http://sirius")(func(w http.ResponseWriter, r *http.Request) error {
		return expectedError
	})

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)
	r.Header.Add("Accept", "application/json")

	handler.ServeHTTP(w, r)

	resp := w.Result()

	assert.Equal(http.StatusBadRequest, resp.StatusCode)
	assert.Equal("application/problem+json", resp.Header.Get("Content-Type"))

	body := new(bytes.Buffer)
	_, _ = body.ReadFrom(resp.Body)
	assert.Equal(`{"title":"validation error","detail":"Not valid complaint","validationErrors":{"title":{"tooShort":"The title must be at least 5 characters"}}}`, strings.Trim(body.String(), "\n"))
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
	r.Header.Add("Content-Type", formUrlEncoded)
	r.AddCookie(&http.Cookie{Name: "XSRF-TOKEN", Value: "z3tVRZ00yx4dHz3KWYv3boLWHZ4/RsCsVAKbvo2SBNc%3D"})
	r.AddCookie(&http.Cookie{Name: "another", Value: "one"})

	ctx := getContext(r)
	assert.Equal(r.Context(), ctx.Context)
	assert.Equal(r.Cookies(), ctx.Cookies)
	assert.Equal("the-real-one", ctx.XSRFToken)
}

func TestPostFormString(t *testing.T) {
	r, _ := http.NewRequest(http.MethodPost, "/?name=ignored", strings.NewReader("name=%20%20%09%0Ahello%0A%0A%20%20%09%09"))
	r.Header.Add("Content-Type", formUrlEncoded)

	assert.Equal(t, "hello", postFormString(r, "name"))
}

func TestPostFormInt(t *testing.T) {
	r, _ := http.NewRequest(http.MethodPost, "/?name=ignored", strings.NewReader("name=%20%20%09%0A123%0A%0A%20%20%09%09"))
	r.Header.Add("Content-Type", formUrlEncoded)

	n, err := postFormInt(r, "name")
	assert.Equal(t, 123, n)
	assert.Nil(t, err)
}

func TestPostFormDateString(t *testing.T) {
	r, _ := http.NewRequest(http.MethodPost, "/?name=ignored", strings.NewReader("name=%20%20%09%0A2022-01-02%0A%0A%20%20%09%09"))
	r.Header.Add("Content-Type", formUrlEncoded)

	assert.Equal(t, sirius.DateString("2022-01-02"), postFormDateString(r, "name"))
}

func TestCancelledContext(t *testing.T) {
	assert := assert.New(t)

	expectedErr := context.Canceled

	logger := &mockLogger{}
	template := &mockTemplate{}

	handler := errorHandler(logger, template.Func, "http://prefix", "http://sirius")(func(w http.ResponseWriter, r *http.Request) error {
		return expectedErr
	})

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/", strings.NewReader("xsrfToken=the-real-one"))

	handler.ServeHTTP(w, r)

	resp := w.Result()

	assert.Equal(499, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, template, logger)
}
