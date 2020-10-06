package thundrahttp

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"net/http"
	gourl "net/url"

	"github.com/thundra-io/thundra-lambda-agent-go/v2/config"

	"github.com/thundra-io/thundra-lambda-agent-go/v2/application"
	"github.com/thundra-io/thundra-lambda-agent-go/v2/tracer"
	"github.com/thundra-io/thundra-lambda-agent-go/v2/utils"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/thundra-io/thundra-lambda-agent-go/v2/constants"
)

var emptyCtx = context.Background()

// ClientWrapper wraps the default http.Client
type ClientWrapper struct {
	http.Client
}

// Wrap wraps the given http.Client with ClientWrapper
func Wrap(c http.Client) *ClientWrapper {
	return &ClientWrapper{c}
}

// Do wraps the http.Client.Do and starts a new span for the http call
func (c *ClientWrapper) Do(req *http.Request) (resp *http.Response, err error) {
	return c.DoWithContext(emptyCtx, req)
}

// DoWithContext wraps the http.Client.Do and starts a new span for
// the http call. The newly created span will be a child of the span
// whose context is is passed using the ctx parameter
func (c *ClientWrapper) DoWithContext(ctx context.Context, req *http.Request) (resp *http.Response, err error) {
	span, _ := opentracing.StartSpanFromContext(
		ctx,
		req.URL.Host+getNormalizedPath(req.URL.Path),
	)
	defer span.Finish()
	rawSpan, ok := tracer.GetRaw(span)
	if ok {
		beforeCall(rawSpan, req.URL.String(), req.Method, req, req.Body)
	}
	tracer.OnSpanStarted(span)
	resp, err = c.Client.Do(req)
	if err != nil {
		utils.SetSpanError(span, err)
	} else if ok {
		afterCall(rawSpan, resp)
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
	span, _ := opentracing.StartSpanFromContext(
		ctx,
		getOperationName(url),
	)
	defer span.Finish()
	rawSpan, ok := tracer.GetRaw(span)
	if ok {
		beforeCall(rawSpan, url, http.MethodGet, nil, nil)
	}
	tracer.OnSpanStarted(span)
	resp, err = c.Client.Get(url)
	if err != nil {
		utils.SetSpanError(span, err)
	} else if ok {
		afterCall(rawSpan, resp)
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
	span, _ := opentracing.StartSpanFromContext(
		ctx,
		getOperationName(url),
	)
	defer span.Finish()
	rawSpan, ok := tracer.GetRaw(span)
	if ok {
		beforeCall(rawSpan, url, http.MethodPost, nil, ioutil.NopCloser(body))
	}
	tracer.OnSpanStarted(span)
	resp, err = c.Client.Post(url, contentType, body)
	if err != nil {
		utils.SetSpanError(span, err)
	} else if ok {
		afterCall(rawSpan, resp)
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
	span, _ := opentracing.StartSpanFromContext(
		ctx,
		getOperationName(url),
	)
	defer span.Finish()
	rawSpan, ok := tracer.GetRaw(span)
	if ok {
		beforeCall(rawSpan, url, http.MethodPost, nil, ioutil.NopCloser(bytes.NewBufferString(data.Encode())))
	}
	tracer.OnSpanStarted(span)
	resp, err = c.Client.PostForm(url, data)
	if err != nil {
		utils.SetSpanError(span, err)
	} else if ok {
		afterCall(rawSpan, resp)
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
	span, _ := opentracing.StartSpanFromContext(
		ctx,
		getOperationName(url),
	)
	defer span.Finish()
	rawSpan, ok := tracer.GetRaw(span)
	if ok {
		beforeCall(rawSpan, url, http.MethodHead, nil, nil)
	}
	tracer.OnSpanStarted(span)
	resp, err = c.Client.Head(url)
	if err != nil {
		utils.SetSpanError(span, err)
	} else if ok {
		afterCall(rawSpan, resp)
	}
	return
}

func getOperationName(url string) string {
	// Parse URLs
	parsedURL, err := gourl.Parse(url)
	if err != nil {
		return ""
	}
	return parsedURL.Host + getNormalizedPath(parsedURL.Path)
}

func getNormalizedPath(urlPath string) string {
	depth := config.HTTPIntegrationUrlPathDepth

	pathSeperatorCount := 0
	buff := bytes.NewBufferString("")
	var prevChar byte = '-'

	for i := 0; i < len(urlPath); i++ {
		if string(urlPath[i]) == "/" && prevChar != '/' {
			pathSeperatorCount++
		}
		if pathSeperatorCount > depth {
			break
		}
		buff.WriteByte(urlPath[i])

		prevChar = urlPath[i]
	}

	return buff.String()
}

func beforeCall(span *tracer.RawSpan, url, method string, req *http.Request, body io.ReadCloser) {
	span.ClassName = constants.ClassNames["HTTP"]
	span.DomainName = constants.DomainNames["API"]

	// Parse URL
	parsedURL, err := gourl.Parse(url)

	// Set span tags
	tags := map[string]interface{}{
		constants.SpanTags["OPERATION_TYPE"]: method,
		constants.HTTPTags["METHOD"]:         method,
	}
	if err == nil {
		tags[constants.HTTPTags["URL"]] = parsedURL.Host + parsedURL.Path
		tags[constants.HTTPTags["PATH"]] = parsedURL.Path
		tags[constants.HTTPTags["HOST"]] = parsedURL.Host
		tags[constants.HTTPTags["QUERY_PARAMS"]] = parsedURL.Query().Encode()
		tags[constants.SpanTags["TRIGGER_DOMAIN_NAME"]] = constants.AwsLambdaApplicationDomain
		tags[constants.SpanTags["TRIGGER_CLASS_NAME"]] = constants.AwsLambdaApplicationClass
		tags[constants.SpanTags["TRIGGER_OPERATION_NAMES"]] = []string{application.FunctionName}
		tags[constants.SpanTags["TOPOLOGY_VERTEX"]] = true
	}

	var bodyLen int64 = constants.MaxTracedHttpBodySize

	if req != nil {
		req.Header.Add("x-thundra-span-id", span.Context.SpanID)
		tags[constants.SpanTags["TRACE_LINKS"]] = []string{span.Context.SpanID}
		bodyLen = req.ContentLength
	}

	if !config.MaskHTTPBody && body != nil {
		bodyRead, newReadCloser := utils.ReadRequestBody(body, int(bodyLen))
		if req != nil {
			req.Body = newReadCloser
		}
		tags[constants.HTTPTags["BODY"]] = bodyRead
	}

	span.Tags = tags
}

func afterCall(span *tracer.RawSpan, resp *http.Response) {
	if resp != nil {
		span.Tags[constants.HTTPTags["STATUS"]] = resp.StatusCode
		if _, ok := resp.Header["X-Amz-Apigw-Id"]; ok {
			span.ClassName = constants.ClassNames["APIGATEWAY"]
		}

		if _, ok := resp.Header["X-Thundra-Resource-Name"]; ok {
			resourceNameHeader := resp.Header["X-Thundra-Resource-Name"]
			if len(resourceNameHeader) > 0 {
				span.OperationName = resourceNameHeader[0]
			}
		}
		if !config.Http4xxErrorDisabled && resp.StatusCode >= 400 && resp.StatusCode <= 499 {
			span.Tags[constants.AwsError] = true
			span.Tags[constants.AwsErrorKind] = "HttpError"
			span.Tags[constants.AwsErrorMessage] = resp.Status
		}
		if !config.Http5xxErrorDisabled && resp.StatusCode >= 500 && resp.StatusCode <= 599 {
			span.Tags[constants.AwsError] = true
			span.Tags[constants.AwsErrorKind] = "HttpError"
			span.Tags[constants.AwsErrorMessage] = resp.Status
		}
	}
}
