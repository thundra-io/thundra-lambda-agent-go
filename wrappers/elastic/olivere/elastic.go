package thundraelastic

import (
	"net/http"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/thundra-io/thundra-lambda-agent-go/tracer"
	"github.com/thundra-io/thundra-lambda-agent-go/utils"
	"github.com/thundra-io/thundra-lambda-agent-go/wrappers/elastic"
)

type roundTripperWrapper struct {
	http.RoundTripper
}

// Wrap wraps the Transport of given http.Client to trace http requests
func Wrap(c *http.Client) *http.Client {
	if c.Transport == nil {
		c.Transport = http.DefaultTransport
	}
	c.Transport = &roundTripperWrapper{c.Transport}
	return c
}

func (t *roundTripperWrapper) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	span, _ := opentracing.StartSpanFromContext(
		req.Context(),
		req.URL.Path,
	)
	defer span.Finish()
	rawSpan, ok := tracer.GetRaw(span)
	if ok {
		elastic.BeforeCall(rawSpan, req)
	}
	tracer.OnSpanStarted(span)
	resp, err = t.RoundTripper.RoundTrip(req)
	if err != nil {
		utils.SetSpanError(span, err)
	} else if ok {
		elastic.AfterCall(rawSpan, resp)
	}
	return
}
