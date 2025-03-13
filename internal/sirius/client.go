package sirius

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type Context struct {
	Context   context.Context
	Cookies   []*http.Cookie
	XSRFToken string
}

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

func (ctx Context) With(c context.Context) Context {
	return Context{
		Context:   c,
		Cookies:   ctx.Cookies,
		XSRFToken: ctx.XSRFToken,
	}
}

func NewClient(httpClient HttpClient, baseURL string) *Client {
	return &Client{
		http:    httpClient,
		baseURL: baseURL,
	}
}

type Client struct {
	http    HttpClient
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

func (c *Client) newRequestWithQuery(ctx Context, method, path string, query url.Values, body io.Reader) (*http.Request, error) {
	req, err := c.newRequest(ctx, method, path, body)
	if err != nil {
		return nil, err
	}

	querystring := req.URL.Query()
	for k, values := range query {
		for _, v := range values {
			querystring.Add(k, v)
		}
	}
	req.URL.RawQuery = querystring.Encode()

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

	defer resp.Body.Close() //#nosec G307 false positive

	if resp.StatusCode != http.StatusOK {
		return newStatusError(resp)
	}

	return json.NewDecoder(resp.Body).Decode(&v)
}

func (c *Client) post(ctx Context, path string, body interface{}, response interface{}) error {
	data := []byte{}

	if body != nil {
		var err error
		if data, err = json.Marshal(body); err != nil {
			return err
		}
	}

	req, err := c.newRequest(ctx, http.MethodPost, path, bytes.NewReader(data))
	if err != nil {
		return err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close() //#nosec G307 false positive

	if resp.StatusCode == http.StatusBadRequest {
		var v ValidationError
		if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
			return err
		}
		return v
	}

	if resp.StatusCode != http.StatusOK &&
		resp.StatusCode != http.StatusCreated &&
		resp.StatusCode != http.StatusNoContent {
		return newStatusError(resp)
	}

	var buf []byte
	if buf, err = io.ReadAll(resp.Body); err != nil {
		return err
	}

	if len(buf) == 0 {
		return nil
	}

	if err := json.Unmarshal(buf, &response); err != nil {
		return err
	}

	return nil
}

func (c *Client) put(ctx Context, path string, body interface{}, response interface{}) error {
	data := []byte{}
	var err error

	if body != nil {
		if data, err = json.Marshal(body); err != nil {
			return err
		}
	}

	req, err := c.newRequest(ctx, http.MethodPut, path, bytes.NewReader(data))
	if err != nil {
		return err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close() //#nosec G307 false positive

	if resp.StatusCode == http.StatusBadRequest {
		var v ValidationError
		if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
			return err
		}
		return v
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return newStatusError(resp)
	}

	var buf []byte
	if buf, err = io.ReadAll(resp.Body); err != nil {
		return err
	}

	if len(buf) == 0 {
		return nil
	}

	if err := json.Unmarshal(buf, &response); err != nil {
		return err
	}

	return nil
}

func (c *Client) delete(ctx Context, path string) error {
	req, err := c.newRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close() //#nosec G307 false positive

	if resp.StatusCode != http.StatusNoContent {
		return newStatusError(resp)
	}

	return nil
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

func (e ValidationError) Error() string {
	if len(e.Detail) > 0 {
		return e.Detail
	}

	return "validation error"
}
