package thundrahttp

import (
	"github.com/thundra-io/thundra-lambda-agent-go/utils"
	"context"
	"io"
	"net/http"
	gourl "net/url"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/thundra-io/thundra-lambda-agent-go/constants"
	"github.com/thundra-io/thundra-lambda-agent-go/ext"
)

var emptyCtx = context.Background()

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
	return c.DoWithContext(emptyCtx, req)
}

// DoWithContext wraps the http.Client.Do and starts a new span for
// the http call. The newly created span will be a child of the span
// whose context is is passed using the ctx parameter
func (c *ClientWrapper) DoWithContext(ctx context.Context, req *http.Request) (resp *http.Response, err error) {
	// TODO: Add recover method
	span, _ := opentracing.StartSpanFromContext(
		ctx,
		req.URL.Path,
		ext.ClassName(constants.HTTPClassName),
		ext.DomainName(constants.HTTPDomainName),
		ext.OperationType(req.Method),
	)
	defer span.Finish()
	addHTTPTags(span, req.URL.String(), req.Method)
	resp, err = c.Client.Do(req)
	if err != nil {
		utils.SetSpanError(span, err)
	}
	return
}

// Get wraps the http.Client.Get and starts a new span for the http call
func (c *ClientWrapper) Get(url string) (resp *http.Response, err error) {
	return c.GetWithContext(emptyCtx, url)
}

// GetWithContext wraps the http.Client.Get and starts a new span for
// the http call. The newly created span will be a child of the span
// whose context is is passed using the ctx parameter
func (c *ClientWrapper) GetWithContext(ctx context.Context, url string) (resp *http.Response, err error) {
	// TODO: Add recover method
	span, _ := opentracing.StartSpanFromContext(
		ctx,
		getOperationName(url),
		ext.ClassName(constants.HTTPClassName),
		ext.DomainName(constants.HTTPDomainName),
		ext.OperationType(http.MethodGet),
	)
	defer span.Finish()
	addHTTPTags(span, url, http.MethodGet)
	resp, err = c.Client.Get(url)
	if err != nil {
		utils.SetSpanError(span, err)
	}
	return
}

// Post wraps the http.Client.Post and starts a new span for the http call
func (c *ClientWrapper) Post(url, contentType string, body io.Reader) (resp *http.Response, err error) {
	return c.PostWithContext(emptyCtx, url, contentType, body)
}

// PostWithContext wraps the http.Client.Post and starts a new span for
// the http call. The newly created span will be a child of the span
// whose context is is passed using the ctx parameter
func (c *ClientWrapper) PostWithContext(ctx context.Context, url, contentType string, body io.Reader) (resp *http.Response, err error) {
	// TODO: Add recover method
	span, _ := opentracing.StartSpanFromContext(
		ctx,
		getOperationName(url),
		ext.ClassName(constants.HTTPClassName),
		ext.DomainName(constants.HTTPDomainName),
		ext.OperationType(http.MethodPost),
	)
	defer span.Finish()
	addHTTPTags(span, url, http.MethodPost)
	resp, err = c.Client.Post(url, contentType, body)
	if err != nil {
		utils.SetSpanError(span, err)
	}
	return
}

// PostForm wraps the http.Client.PostForm and starts a new span for the http call
func (c *ClientWrapper) PostForm(url string, data gourl.Values) (resp *http.Response, err error) {
	return c.PostFormWithContext(emptyCtx, url, data)
}

// PostFormWithContext wraps the http.Client.PostForm and starts a new span
// for the http call. The newly created span will be a child of the span
// whose context is is passed using the ctx parameter
func (c *ClientWrapper) PostFormWithContext(ctx context.Context, url string, data gourl.Values) (resp *http.Response, err error) {
	// Parse URL
	span, _ := opentracing.StartSpanFromContext(
		ctx,
		getOperationName(url),
		ext.ClassName(constants.HTTPClassName),
		ext.DomainName(constants.HTTPDomainName),
		ext.OperationType(http.MethodPost),
	)
	defer span.Finish()
	addHTTPTags(span, url, http.MethodPost)
	resp, err = c.Client.PostForm(url, data)
	if err != nil {
		utils.SetSpanError(span, err)
	}
	return
}

// Head wraps the http.Client.Head and starts a new span for the http call
func (c *ClientWrapper) Head(url string) (resp *http.Response, err error) {
	return c.HeadWithContext(emptyCtx, url)
}

// HeadWithContext wraps the http.Client.Head and starts a new span
// for the http call. The newly created span will be a child of the span
// whose context is is passed using the ctx parameter
func (c *ClientWrapper) HeadWithContext(ctx context.Context, url string) (resp *http.Response, err error) {
	// Parse URL
	span, _ := opentracing.StartSpanFromContext(
		ctx,
		getOperationName(url),
		ext.ClassName(constants.HTTPClassName),
		ext.DomainName(constants.HTTPDomainName),
		ext.OperationType(http.MethodHead),
	)
	defer span.Finish()
	addHTTPTags(span, url, http.MethodHead)
	resp, err = c.Client.Head(url)
	if err != nil {
		utils.SetSpanError(span, err)
	}
	return
}

func addHTTPTags(span opentracing.Span, url, method string) {
	// Parse URL
	parsedURL, err := gourl.Parse(url)

	// Set span tags
	span.SetTag(constants.HTTPMethodTag, method)
	if err == nil {
		span.SetTag(constants.HTTPURLTag, parsedURL.Host+parsedURL.Path)
		span.SetTag(constants.HTTPPathTag, parsedURL.Path)
		span.SetTag(constants.HTTPHostTag, parsedURL.Host)
		span.SetTag(constants.HTTPQueryParamsTag, parsedURL.Query().Encode())
	}

}

func getOperationName(url string) string {
	// Parse URLs
	parsedURL, err := gourl.Parse(url)
	if err != nil {
		return ""
	}
	return parsedURL.Host + parsedURL.Path
}
