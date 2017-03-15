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

type TraceClient struct {
	RequestIDFunc RequestIDFunc
	UserAgent     string
	client        *pester.Client
}

type RequestIDFunc func() string

func (t *TraceClient) Do(req *http.Request) (*http.Response, error) {
	req.Header.Set("User-Agent", t.UserAgent)
	requestidmw.SetID(&req.Header, t.RequestIDFunc())
	return t.client.Do(req)
}

func (t *TraceClient) Get(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	return t.Do(req)
}

func (t *TraceClient) Head(url string) (*http.Response, error) {
	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		return nil, err
	}
	return t.Do(req)
}

func (t *TraceClient) Post(url string, contentType string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)
	return t.Do(req)
}

func (t *TraceClient) PostForm(url string, data url.Values) (*http.Response, error) {
	return t.Post(url, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
}

func (t *TraceClient) SetUserAgent(agent string) {
	t.UserAgent = agent
}

func NewClient(rIdFunc RequestIDFunc) *TraceClient {
	client := pester.New()
	client.Backoff = pester.LinearBackoff
	client.MaxRetries = 3
	return NewExtendedClient(rIdFunc, client)
}

func NewExtendedClient(rIdFunc RequestIDFunc, client *pester.Client) *TraceClient {
	return &TraceClient{
		RequestIDFunc: rIdFunc,
		UserAgent:     env.ServiceName(),
		client:        client,
	}
}
