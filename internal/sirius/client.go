package sirius

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Context struct {
	Context   context.Context
	Cookies   []*http.Cookie
	XSRFToken string
}

func (ctx Context) With(c context.Context) Context {
	return Context{
		Context:   c,
		Cookies:   ctx.Cookies,
		XSRFToken: ctx.XSRFToken,
	}
}

func NewClient(httpClient *http.Client, baseURL string) *Client {
	return &Client{
		http:    httpClient,
		baseURL: baseURL,
	}
}

type Client struct {
	http    *http.Client
	baseURL string
}

func (c *Client) newRequest(ctx Context, method, path string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx.Context, method, c.baseURL+path, body)
	if err != nil {
		return nil, err
	}

	for _, c := range ctx.Cookies {
		req.AddCookie(c)
	}

	req.Header.Add("OPG-Bypass-Membrane", "1")
	req.Header.Add("X-XSRF-TOKEN", ctx.XSRFToken)
	if body != nil {
		req.Header.Add("Content-Type", "application/json")
	}

	return req, err
}

func (c *Client) get(ctx Context, path string, v interface{}) error {
	req, err := c.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return newStatusError(resp)
	}

	return json.NewDecoder(resp.Body).Decode(&v)
}

type StatusError struct {
	Code          int    `json:"code"`
	URL           string `json:"url"`
	Method        string `json:"method"`
	CorrelationId string
}

func newStatusError(resp *http.Response) StatusError {
	return StatusError{
		Code:          resp.StatusCode,
		URL:           resp.Request.URL.String(),
		Method:        resp.Request.Method,
		CorrelationId: resp.Header.Get("Correlation-Id"),
	}
}

func (e StatusError) IsUnauthorized() bool {
	return e.Code == http.StatusUnauthorized
}

func (e StatusError) Error() string {
	return fmt.Sprintf("%s %s returned %d", e.Method, e.URL, e.Code)
}

func (StatusError) Title() string {
	return "unexpected response from Sirius"
}

func (e StatusError) Data() interface{} {
	return e
}

type FieldErrors map[string]map[string]string
type flexibleFieldErrors map[string]json.RawMessage

type ValidationError struct {
	Detail string      `json:"detail"`
	Field  FieldErrors `json:"validation_errors"`
}

func (f flexibleFieldErrors) toFieldErrors() (FieldErrors, error) {
	s := FieldErrors{}

	for k, v := range f {
		var asSlice []string
		if err := json.Unmarshal(v, &asSlice); err == nil {
			s[k] = map[string]string{"": strings.Join(asSlice, "")}
			continue
		}

		var asMap map[string]string
		if err := json.Unmarshal(v, &asMap); err == nil {
			s[k] = asMap
			continue
		}

		return nil, errors.New("could not parse field validation_errors")
	}
	return s, nil
}

func (e ValidationError) Any() bool {
	return len(e.Detail) > 0 || len(e.Field) > 0
}

func (ValidationError) Error() string {
	return "validation error"
}
