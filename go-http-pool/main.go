package main

import (
    "context"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "net/url"
    "os"

    "golang.org/x/sync/errgroup"
)

type HttpClient interface {
    Do(req *http.Request) (*http.Response, error)
}

type Context struct {
    Context context.Context
    Cookies []*http.Cookie
    XSRFToken string
}

type Client struct {
    http HttpClient
    baseURL string
    maxPoolSize int
    requests map[string]*http.Request
    receivers map[string]interface{}
}

// key is a unique identifier for the request and its receiver; if you enqueue
// another request with the same key, it is ignored;
// if you try to enqueue a request already in the queue (i.e. same path and method),
// that will also be ignored
func (c *Client) Enqueue(key string, r *http.Request, rec interface{}) bool {
    // If key is already in the queue, ignore request
    _, exists := c.requests[key]
    if exists {
        return false
    }

    // If method is GET and path matches existing request, ignore request
    requestUri := r.URL.RequestURI()

    if r.Method == http.MethodGet {
        for _, req := range c.requests {
            if req.URL.RequestURI() == requestUri && req.Method == http.MethodGet {
                return false
            }
        }
    }

    c.requests[key] = r
    c.receivers[key] = rec

    return true
}

func (c *Client) NewRequest(ctx Context, method, path string, query url.Values, body io.Reader) (*http.Request, error) {
    req, err := http.NewRequestWithContext(ctx.Context, method, c.baseURL + path, body)
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

    querystring := req.URL.Query()
    for k, values := range query {
        for _, v := range values {
            querystring.Add(k, v)
        }
    }
    req.URL.RawQuery = querystring.Encode()

    return req, err
}

func (c *Client) send(req *http.Request, receiver interface{}) error {
    resp, err := c.http.Do(req)
    if err != nil {
        return err
    }

    defer resp.Body.Close() //#nosec G307 false positive

    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("%s", resp.Body)
    }

    return json.NewDecoder(resp.Body).Decode(&receiver)
}

func (c *Client) Dispatch(ctx context.Context) error {
    g, ctx := errgroup.WithContext(ctx)
    g.SetLimit(c.maxPoolSize)

    for key, request := range c.requests {
        key, request := key, request // see https://go.dev/doc/faq#closures_and_goroutines
        g.Go(func() error {
            receiver := c.receivers[key]
            err := c.send(request, &receiver)

            // request is done, remove it;
            // should we retain if it threw an error?
            delete(c.requests, key)

            // alternatively, we could be building a map from key to response|error
            // and return that from Dispatch(); note that the func passed to Go()
            // has to have this signature and can only return an error
            return err
        })
    }

    return g.Wait()
}

func NewClient() *Client {
    return &Client{
        http: http.DefaultClient,
        baseURL: "http://localhost:5000",
        maxPoolSize: 2,
        requests: map[string]*http.Request{},
        receivers: map[string]interface{}{},
    }
}

// runtime stuff; these structs marry up with test-server.py's endpoints
type HelloWorldResult struct {
    A string `json:"a"`
    B int `json:"b"`
}

type GoodbyeCruelWorldResult struct {
    C int `json:"c"`
    D int `json:"d"`
}

type SayItResult struct {
    Say string `json:"say"`
}

type HelloSunshineF struct {
    G string `json:"g"`
    H string `json:"h"`
}

type HelloSunshineResult struct {
    E int `json:"e"`
    F HelloSunshineF
}

type MultiResult struct {
    Foo HelloWorldResult
    Bar HelloSunshineResult
    Baz GoodbyeCruelWorldResult
    Boo GoodbyeCruelWorldResult
    Faz GoodbyeCruelWorldResult
}

// this currently writes into a struct; but if a request is rejected because
// it is already in the queue, the target part of the struct will remain unchanged
func main() {
    c := NewClient()
    ctx := Context{
        Context: context.Background(),
    }

    result := MultiResult{
        Foo: HelloWorldResult{},
        Bar: HelloSunshineResult{},
        Baz: GoodbyeCruelWorldResult{},
        Boo: GoodbyeCruelWorldResult{},
        Faz: GoodbyeCruelWorldResult{},
    }

    // try to enqueue 7 requests; two are rejected: one is a duplicate, the other an existing key
    r1, _ := c.NewRequest(ctx, http.MethodGet, "/", nil, nil)
    c.Enqueue("foo", r1, &result.Foo)

    r2, _ := c.NewRequest(ctx, http.MethodGet, "/hello_sunshine", nil, nil)
    c.Enqueue("bar", r2, &result.Bar)

    r3, _ := c.NewRequest(ctx, http.MethodGet, "/goodbye_cruel_world?page=1", nil, nil)
    c.Enqueue("baz", r3, &result.Baz)

    // this request is honoured as the querystring is different
    q1 := url.Values{}
    q1.Add("page", "2")
    r4, _ := c.NewRequest(ctx, http.MethodGet, "/goodbye_cruel_world", q1, nil)
    c.Enqueue("boo", r4, &result.Boo)

    // this request is not honoured as the path+querystring matches one already in the queue
    q2 := url.Values{}
    q2.Add("page", "1")
    r5, _ := c.NewRequest(ctx, http.MethodGet, "/goodbye_cruel_world", q2, nil)
    c.Enqueue("faz", r5, &result.Faz)

    // this request is not honoured as the key matches one already in the queue
    r6, _ := c.NewRequest(ctx, http.MethodGet, "/goodbye_cruel_world", nil, nil)
    c.Enqueue("bar", r6, &result.Faz)

    // send the requests in parallel
    err := c.Dispatch(context.Background())
    if err != nil {
        fmt.Fprintln(os.Stderr, err)
    }
    fmt.Println(result)

    // use something from a previous request in next request
    result2 := SayItResult{}

    q3 := url.Values{}
    q3.Add("say_this", result.Foo.A)

    r7, _ := c.NewRequest(ctx, http.MethodGet, "/say_it", q3, nil)

    c.Enqueue("say", r7, &result2)
    err = c.Dispatch(context.Background())
    if err != nil {
        fmt.Fprintln(os.Stderr, err)
    }
    fmt.Println(result2)
}
