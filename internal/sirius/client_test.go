package sirius

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func getContext(cookies []*http.Cookie) Context {
	return Context{
		Context:   context.Background(),
		Cookies:   cookies,
		XSRFToken: "abcde",
	}
}

func TestStatusError(t *testing.T) {
	req, _ := http.NewRequest(http.MethodPost, "/some/url", nil)

	resp := &http.Response{
		StatusCode: http.StatusTeapot,
		Request:    req,
	}

	err := newStatusError(resp)

	assert.Equal(t, "POST /some/url returned 418", err.Error())
	assert.Equal(t, "unexpected response from Sirius", err.Title())
	assert.Equal(t, err, err.Data())
	assert.False(t, err.IsUnauthorized())
}

func TestStatusErrorIsUnauthorized(t *testing.T) {
	req, _ := http.NewRequest(http.MethodPost, "/some/url", nil)

	resp := &http.Response{
		StatusCode: http.StatusUnauthorized,
		Request:    req,
	}

	err := newStatusError(resp)

	assert.Equal(t, "POST /some/url returned 401", err.Error())
	assert.Equal(t, "unexpected response from Sirius", err.Title())
	assert.Equal(t, err, err.Data())
	assert.True(t, err.IsUnauthorized())
}
