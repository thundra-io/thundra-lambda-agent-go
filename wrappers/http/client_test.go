package thundrahttp

import (
	"bytes"
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/thundra-io/thundra-lambda-agent-go/tracer"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/thundra-io/thundra-lambda-agent-go/constants"

	"github.com/stretchr/testify/assert"
	"github.com/thundra-io/thundra-lambda-agent-go/trace"
)

type RoundTripFunc func(req *http.Request) (*http.Response, error)

func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func NewTestClient(f RoundTripFunc) http.Client {
	return http.Client{
		Transport: f,
	}
}

var client = Wrap(NewTestClient(func(req *http.Request) (*http.Response, error) {
	return nil, http.ErrServerClosed
}))

func TestHTTPGet(t *testing.T) {
	// Initilize trace plugin to set GlobalTracer of opentracing
	tp := trace.New()
	// Actual call
	resp, err := client.Get("https://httpbin.org/get?foo=bar")
	// Get the span created for http call
	span := tp.Recorder.GetSpans()[0]
	// Test HTTP related fields of span
	assert.Nil(t, resp)
	assert.NotNil(t, err)
	assert.Equal(t, constants.HTTPClassName, span.ClassName)
	assert.Equal(t, constants.HTTPDomainName, span.DomainName)
	assert.Equal(t, "httpbin.org", span.Tags[constants.HTTPHostTag].(string))
	assert.Equal(t, http.MethodGet, span.Tags[constants.HTTPMethodTag].(string))
	assert.Equal(t, "/get", span.Tags[constants.HTTPPathTag].(string))
	assert.Equal(t, "foo=bar", span.Tags[constants.HTTPQueryParamsTag].(string))
	assert.Equal(t, "httpbin.org/get", span.Tags[constants.HTTPURLTag].(string))
	assert.True(t, span.Tags[constants.AwsError].(bool))
	assert.Equal(t, "Error", span.Tags[constants.AwsErrorKind].(string))
	assert.Equal(t, "Get https://httpbin.org/get?foo=bar: http: Server closed", 
		span.Tags[constants.AwsErrorMessage].(string))
	// Clear tracer
	tp.Reset()
}

func TestHTTPGetWithContext(t *testing.T) {
	// Initilize trace plugin to set GlobalTracer of opentracing
	tp := trace.New()
	// Create the parent span
	ctx := context.Background()
	parentSpan, ctx := opentracing.StartSpanFromContext(ctx, "parentSpan")
	parentSpanRaw, _ := tracer.GetRaw(parentSpan)
	// Actual call
	resp, err := client.GetWithContext(ctx, "https://httpbin.org/get?foo=bar")
	// Get the span created for http call
	span := tp.Recorder.GetSpans()[1]
	// Check parent span is set 
	assert.Nil(t, resp)
	assert.NotNil(t, err)
	assert.Equal(t, parentSpanRaw.Context.SpanID, span.ParentSpanID)
}

func TestHTTPPost(t *testing.T) {
	// Initilize trace plugin to set GlobalTracer of opentracing
	tp := trace.New()
	// Actual call
	jsonStr := `{"title":"Foobar"}`
	resp, err := client.Post("https://httpbin.org/post?foo=bar", "application/json", bytes.NewBufferString(jsonStr))
	// Get the span created for http call
	span := tp.Recorder.GetSpans()[0]
	// Test HTTP related fields of span
	assert.NotNil(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, constants.HTTPClassName, span.ClassName)
	assert.Equal(t, constants.HTTPDomainName, span.DomainName)
	assert.Equal(t, "httpbin.org", span.Tags[constants.HTTPHostTag].(string))
	assert.Equal(t, http.MethodPost, span.Tags[constants.HTTPMethodTag].(string))
	assert.Equal(t, "/post", span.Tags[constants.HTTPPathTag].(string))
	assert.Equal(t, "foo=bar", span.Tags[constants.HTTPQueryParamsTag].(string))
	assert.Equal(t, "httpbin.org/post", span.Tags[constants.HTTPURLTag].(string))
	assert.True(t, span.Tags[constants.AwsError].(bool))
	assert.Equal(t, "Error", span.Tags[constants.AwsErrorKind].(string))
	assert.Equal(t, "Post https://httpbin.org/post?foo=bar: http: Server closed", 
		span.Tags[constants.AwsErrorMessage].(string))
	// Clear tracer
	tp.Reset()
}

func TestHTTPPostWithContext(t *testing.T) {
	// Initilize trace plugin to set GlobalTracer of opentracing
	tp := trace.New()
	// Create the parent span
	ctx := context.Background()
	parentSpan, ctx := opentracing.StartSpanFromContext(ctx, "parentSpan")
	parentSpanRaw, _ := tracer.GetRaw(parentSpan)
	// Actual call
	jsonStr := `{"title":"Foobar"}`
	resp, err := client.PostWithContext(ctx, "https://httpbin.org/post?foo=bar", "application/json", bytes.NewBufferString(jsonStr))
	// Get the span created for http call
	span := tp.Recorder.GetSpans()[1]
	// Check parent span is set
	assert.NotNil(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, parentSpanRaw.Context.SpanID, span.ParentSpanID)
}

func TestHTTPPostForm(t *testing.T) {
	// Initilize trace plugin to set GlobalTracer of opentracing
	tp := trace.New()
	// Actual call
	v := url.Values{}
	v.Set("name", "Ava")
	v.Add("friend", "Jess")
	v.Add("friend", "Sarah")
	resp, err := client.PostForm("https://httpbin.org/post?foo=bar", v)
	// Get the span created for http call
	span := tp.Recorder.GetSpans()[0]
	// Test HTTP related fields of span
	assert.NotNil(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, constants.HTTPClassName, span.ClassName)
	assert.Equal(t, constants.HTTPDomainName, span.DomainName)
	assert.Equal(t, "httpbin.org", span.Tags[constants.HTTPHostTag].(string))
	assert.Equal(t, http.MethodPost, span.Tags[constants.HTTPMethodTag].(string))
	assert.Equal(t, "/post", span.Tags[constants.HTTPPathTag].(string))
	assert.Equal(t, "foo=bar", span.Tags[constants.HTTPQueryParamsTag].(string))
	assert.Equal(t, "httpbin.org/post", span.Tags[constants.HTTPURLTag].(string))
	assert.True(t, span.Tags[constants.AwsError].(bool))
	assert.Equal(t, "Error", span.Tags[constants.AwsErrorKind].(string))
	assert.Equal(t, "Post https://httpbin.org/post?foo=bar: http: Server closed", 
		span.Tags[constants.AwsErrorMessage].(string))
	// Clear tracer
	tp.Reset()
}

func TestHTTPPostFormWithContext(t *testing.T) {
	// Initilize trace plugin to set GlobalTracer of opentracing
	tp := trace.New()
	// Create the parent span
	ctx := context.Background()
	parentSpan, ctx := opentracing.StartSpanFromContext(ctx, "parentSpan")
	parentSpanRaw, _ := tracer.GetRaw(parentSpan)
	// Actual call
	v := url.Values{}
	v.Set("name", "Ava")
	v.Add("friend", "Jess")
	v.Add("friend", "Sarah")
	resp, err := client.PostFormWithContext(ctx, "https://httpbin.org/post?foo=bar", v)
	// Get the span created for http call
	span := tp.Recorder.GetSpans()[1]
	// Check parent span is set
	assert.NotNil(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, parentSpanRaw.Context.SpanID, span.ParentSpanID)
}

func TestHTTPDo(t *testing.T) {
	// Initilize trace plugin to set GlobalTracer of opentracing
	tp := trace.New()
	// Actual call
	req, _ := http.NewRequest(http.MethodGet, "https://httpbin.org/get?foo=bar", nil)
	resp, err := client.Do(req)
	// Get the span created for http call
	span := tp.Recorder.GetSpans()[0]
	// Test HTTP related fields of span
	assert.NotNil(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, constants.HTTPClassName, span.ClassName)
	assert.Equal(t, constants.HTTPDomainName, span.DomainName)
	assert.Equal(t, "httpbin.org", span.Tags[constants.HTTPHostTag].(string))
	assert.Equal(t, http.MethodGet, span.Tags[constants.HTTPMethodTag].(string))
	assert.Equal(t, "/get", span.Tags[constants.HTTPPathTag].(string))
	assert.Equal(t, "foo=bar", span.Tags[constants.HTTPQueryParamsTag].(string))
	assert.Equal(t, "httpbin.org/get", span.Tags[constants.HTTPURLTag].(string))
	assert.True(t, span.Tags[constants.AwsError].(bool))
	assert.Equal(t, "Error", span.Tags[constants.AwsErrorKind].(string))
	assert.Equal(t, "Get https://httpbin.org/get?foo=bar: http: Server closed", 
		span.Tags[constants.AwsErrorMessage].(string))
	// Clear tracer
	tp.Reset()
}

func TestHTTPDoWithContext(t *testing.T) {
	// Initilize trace plugin to set GlobalTracer of opentracing
	tp := trace.New()
	// Create the parent span
	ctx := context.Background()
	parentSpan, ctx := opentracing.StartSpanFromContext(ctx, "parentSpan")
	parentSpanRaw, _ := tracer.GetRaw(parentSpan)
	// Actual call
	req, _ := http.NewRequest(http.MethodGet, "https://httpbin.org/get?foo=bar", nil)
	resp, err := client.DoWithContext(ctx, req)
	// Get the span created for http call
	span := tp.Recorder.GetSpans()[1]
	// Check parent span is set
	assert.NotNil(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, parentSpanRaw.Context.SpanID, span.ParentSpanID)
}

func TestHTTPHead(t *testing.T) {
	// Initilize trace plugin to set GlobalTracer of opentracing
	tp := trace.New()
	// Actual call
	resp, err := client.Head("https://httpbin.org/head?foo=bar")
	// Get the span created for http call
	span := tp.Recorder.GetSpans()[0]
	// Test HTTP related fields of span
	assert.NotNil(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, constants.HTTPClassName, span.ClassName)
	assert.Equal(t, constants.HTTPDomainName, span.DomainName)
	assert.Equal(t, "httpbin.org", span.Tags[constants.HTTPHostTag].(string))
	assert.Equal(t, http.MethodHead, span.Tags[constants.HTTPMethodTag].(string))
	assert.Equal(t, "/head", span.Tags[constants.HTTPPathTag].(string))
	assert.Equal(t, "foo=bar", span.Tags[constants.HTTPQueryParamsTag].(string))
	assert.Equal(t, "httpbin.org/head", span.Tags[constants.HTTPURLTag].(string))
	assert.True(t, span.Tags[constants.AwsError].(bool))
	assert.Equal(t, "Error", span.Tags[constants.AwsErrorKind].(string))
	assert.Equal(t, "Head https://httpbin.org/head?foo=bar: http: Server closed", 
		span.Tags[constants.AwsErrorMessage].(string))
	// Clear tracer
	tp.Reset()
}

func TestHTTPHeadWithContext(t *testing.T) {
	// Initilize trace plugin to set GlobalTracer of opentracing
	tp := trace.New()
	// Create the parent span
	ctx := context.Background()
	parentSpan, ctx := opentracing.StartSpanFromContext(ctx, "parentSpan")
	parentSpanRaw, _ := tracer.GetRaw(parentSpan)
	// Actual call
	resp, err := client.HeadWithContext(ctx, "https://httpbin.org/head?foo=bar")
	// Get the span created for http call
	span := tp.Recorder.GetSpans()[1]
	// Test HTTP related fields of span
	assert.NotNil(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, parentSpanRaw.Context.SpanID, span.ParentSpanID)
}
