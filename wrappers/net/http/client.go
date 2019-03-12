package thundrahttp

import (
	"io"
	"net/http"
	"net/url"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/thundra-io/thundra-lambda-agent-go/constants"
	"github.com/thundra-io/thundra-lambda-agent-go/ext"
)

// ClientWrapper wraps the default http.Client
type ClientWrapper struct {
	http.Client
}

// Wrap wraps the given http.Client with ClientWrapper
func Wrap(c http.Client) ClientWrapper {
	return ClientWrapper{c}
}

// Do wraps the http.Client.Do and starts a new span for the http call
func (c *ClientWrapper) Do(req *http.Request) (resp *http.Response, err error) {
	// TODO: Add recover method
	span := opentracing.StartSpan(
		req.URL.Path,
		ext.ClassName("HTTP"),
		ext.DomainName("API"),
		ext.OperationType(req.Method),
	)
	defer span.Finish()
	addHTTPTags(span, req.URL.String(), req.Method)
	resp, err = c.Client.Do(req)
	return
}

// Get wraps the http.Client.Get and starts a new span for the http call
func (c *ClientWrapper) Get(URL string) (resp *http.Response, err error) {
	// TODO: Add recover method
	// Create and finish(defered) span
	span := opentracing.StartSpan(
		getOperationName(URL),
		ext.ClassName("HTTP"),
		ext.DomainName("API"),
		ext.OperationType("GET"),
	)
	defer span.Finish()
	addHTTPTags(span, URL, "GET")
	resp, err = c.Client.Get(URL)
	return
}

// Post wraps the http.Client.Post and starts a new span for the http call
func (c *ClientWrapper) Post(URL, contentType string, body io.Reader) (resp *http.Response, err error) {
	// TODO: Add recover method
	span := opentracing.StartSpan(
		getOperationName(URL),
		ext.ClassName("HTTP"),
		ext.DomainName("API"),
		ext.OperationType("POST"),
	)
	defer span.Finish()
	addHTTPTags(span, URL, "POST")
	resp, err = c.Client.Post(URL, contentType, body)
	return
}

func addHTTPTags(span opentracing.Span, URL, method string) {
	// Parse URL
	parsedURL, err := url.Parse(URL)

	// Set span tags
	span.SetTag(constants.HttpMethodTag, method)
	if err == nil {
		span.SetTag(constants.HttpURLTag, parsedURL.Host+parsedURL.Path)
		span.SetTag(constants.HttpPathTag, parsedURL.Path)
		span.SetTag(constants.HttpHostTag, parsedURL.Host)
		span.SetTag(constants.HttpQueryParamsTag, parsedURL.Query().Encode())
	}

}

func getOperationName(URL string) string {
	// Parse URL
	parsedURL, err := url.Parse(URL)
	if err != nil {
		return ""
	}
	return parsedURL.Host + parsedURL.Path
}
