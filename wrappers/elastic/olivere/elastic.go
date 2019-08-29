package thundraelastic

import (
	"net/http"
	"strings"

	"github.com/thundra-io/thundra-lambda-agent-go/config"

	"github.com/opentracing/opentracing-go"
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

func getNormalizedPath(path string) string {
	depth := config.ESIntegrationUrlPathDepth
	if depth <= 0 {
		return ""
	}

	pathSlice := strings.Split(path, "/")

	//filter empty string
	n := 0
	for _, x := range pathSlice {
		if len(x) > 0 {
			pathSlice[n] = x
			n++
		}
	}
	pathSlice = pathSlice[:n]

	// check out of bounds
	pathLength := len(pathSlice)
	if depth > pathLength {
		depth = pathLength
	}

	//slice till depth
	pathSlice = pathSlice[:depth]
	return "/" + strings.Join(pathSlice, "/")
}

func (t *roundTripperWrapper) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	normalizedPath := getNormalizedPath(req.URL.Path)
	span, _ := opentracing.StartSpanFromContext(
		req.Context(),
		normalizedPath,
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
