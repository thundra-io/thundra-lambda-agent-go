package thundrahttp

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

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
	return &http.Response{
		StatusCode: 200,
		Body:       ioutil.NopCloser(bytes.NewBufferString(`OK`)),
		Header:     make(http.Header),
	}, nil
}))

func TestHTTPGet(t *testing.T) {
	// Initilize trace plugin to set GlobalTracer of opentracing
	tp := trace.New()
	// Actual call
	resp, err := client.Get("https://httpbin.org/get?foo=bar")
	// Get the span created for http call
	span := tp.Recorder.GetSpans()[0]
	// Test HTTP related fields of span
	assert.Nil(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, span.ClassName, constants.HTTPClassName)
	assert.Equal(t, span.DomainName, constants.HTTPDomainName)
	assert.Equal(t, span.Tags[constants.HTTPHostTag].(string), "httpbin.org")
	assert.Equal(t, span.Tags[constants.HTTPMethodTag].(string), http.MethodGet)
	assert.Equal(t, span.Tags[constants.HTTPPathTag].(string), "/get")
	assert.Equal(t, span.Tags[constants.HTTPQueryParamsTag].(string), "foo=bar")
	assert.Equal(t, span.Tags[constants.HTTPURLTag].(string), "httpbin.org/get")
}
