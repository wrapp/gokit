package trace

import (
	"net/http"
	"net/url"

	"io"

	"strings"

	"github.com/sethgrid/pester"
	"github.com/wrapp/gokit/env"
	"github.com/wrapp/gokit/middleware/requestidmw"
)

// TraceClient struct provides the data which is required to make http requests. It contains a
// func which can generate request-ids. It contains the UserAgent which should be set to all
// outgoing request headers. It also contains the base http client for making the requests.
// TraceClient uses pester client underlying to provide retries for http requests. By default
// it retries 3 times if a request fails but it can be customized by passing a custom pester
// client if needed.
type TraceClient struct {
	RequestIDFunc RequestIDFunc
	UserAgent     string
	client        *pester.Client
}

// A function type that generates a request-id as a string
type RequestIDFunc func() string

// Do performs the passed http.Request. This method can be used to perform any custom
// http requests. It returns the http.Response object or an error if there was a problem
// performing this request.
func (t *TraceClient) Do(req *http.Request) (*http.Response, error) {
	req.Header.Set("User-Agent", t.UserAgent)
	requestidmw.SetIDInHeader(&req.Header, t.RequestIDFunc())
	return t.client.Do(req)
}

// Get sends a GET request to passed url. It returns the http.Response object or an error
// if there was a problem performing this request.
func (t *TraceClient) Get(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	return t.Do(req)
}

// Head sends a HEAD request to passed url. It returns the http.Response object or an error
// if there was a problem performing this request.
func (t *TraceClient) Head(url string) (*http.Response, error) {
	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		return nil, err
	}
	return t.Do(req)
}

// Post sends a POST request to passed url. Post also accepts content-type and the request body
// which should be sent in the request. It returns the http.Response object or an error if there
// was a problem performing this request.
func (t *TraceClient) Post(url string, contentType string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)
	return t.Do(req)
}

// PostForm sends a POST request to passed url. Post accepts form values to be sent as a form request.
// It returns the http.Response object or an error if there was a problem performing this request.
func (t *TraceClient) PostForm(url string, data url.Values) (*http.Response, error) {
	return t.Post(url, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
}

// SetUserAgent sets the user-agent header for each request sent from the client.
func (t *TraceClient) SetUserAgent(agent string) {
	t.UserAgent = agent
}

// New creates a new TraceClient. It accepts a function which can generate request-ids to be set
// for the outgoing request. The returned client will retry the request 3 times with a linear
// backoff if the request failed to execute.
func New(rIdFunc RequestIDFunc) *TraceClient {
	client := pester.New()
	client.Backoff = pester.LinearBackoff
	client.MaxRetries = 3
	return NewExtendedClient(rIdFunc, client)
}

// NewExtendedClient generates an extended client which uses the passed pester client to perform
// http requests. RequestIDFunc is used to get the request id which is set in the header of the
// request.
func NewExtendedClient(rIdFunc RequestIDFunc, client *pester.Client) *TraceClient {
	return &TraceClient{
		RequestIDFunc: rIdFunc,
		UserAgent:     env.ServiceName(),
		client:        client,
	}
}
