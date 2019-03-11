package thundrahttp

import (
	"github.com/thundra-io/thundra-lambda-agent-go/ext"
	"github.com/opentracing/opentracing-go"
	"net/http"
)

// ClientWrapper wraps the default http.Client
type ClientWrapper struct {
	http.Client
}

// Wrap wraps the given http.Client with ClientWrapper
func Wrap(c http.Client) ClientWrapper {
	return ClientWrapper{c}
}

// Do wraps the http.Client.Do, starts a new span for the http call
func (c *ClientWrapper) Do(req *http.Request) (resp *http.Response, err error) {
	// TODO: Add recover method
	span := opentracing.StartSpan(
		req.URL.Path,
		ext.ClassName("HTTP"),
		ext.DomainName("API"),
	)
	defer span.Finish()
	resp, err = c.Client.Do(req)
	return
}
