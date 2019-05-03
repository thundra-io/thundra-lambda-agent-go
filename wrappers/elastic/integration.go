package thundraelastic

import (
	"bufio"
	"io"
	"net/http"

	"github.com/thundra-io/thundra-lambda-agent-go/application"
	"github.com/thundra-io/thundra-lambda-agent-go/config"
	"github.com/thundra-io/thundra-lambda-agent-go/constants"
	"github.com/thundra-io/thundra-lambda-agent-go/tracer"
)

type readCloser struct {
	io.Reader
	io.Closer
}

func BeforeCall(span *tracer.RawSpan, req *http.Request) {
	span.ClassName = constants.ClassNames["ELASTICSEARCH"]
	span.DomainName = constants.DomainNames["DB"]

	host := req.URL.Host
	method := req.Method
	esBody := ""

	// Set span tags
	tags := map[string]interface{}{
		constants.SpanTags["OPERATION_TYPE"]:          method,
		constants.EsTags["ES_HOSTS"]:                  []string{host},
		constants.EsTags["ES_URI"]:                    req.URL.Path,
		constants.EsTags["ES_METHOD"]:                 method,
		constants.EsTags["ES_PARAMS"]:                 req.URL.Query().Encode(),
		constants.DBTags["DB_TYPE"]:                   "elasticsearch",
		constants.SpanTags["TRIGGER_DOMAIN_NAME"]:     constants.AwsLambdaApplicationDomain,
		constants.SpanTags["TRIGGER_CLASS_NAME"]:      constants.AwsLambdaApplicationClass,
		constants.SpanTags["TRIGGER_OPERATION_NAMES"]: []string{application.FunctionName},
		constants.SpanTags["TOPOLOGY_VERTEX"]:         true,
	}

	if req.Body != nil {
		esBody, req.Body = readRequestBody(req.Body, int(req.ContentLength))
	}
	if !config.MaskEsBody {
		tags[constants.EsTags["ES_BODY"]] = esBody
	}

	span.Tags = tags
}

func AfterCall(span *tracer.RawSpan, resp *http.Response) {

}

func readRequestBody(body io.ReadCloser, contentLength int) (string, io.ReadCloser) {
	bodySize := constants.MaxTracedHttpBodySize
	if contentLength > 0 && contentLength < bodySize {
		bodySize = contentLength
	}
	rd := bufio.NewReaderSize(body, bodySize)

	rc := readCloser{
		Reader: rd,
		Closer: body,
	}
	bodyLimited, err := rd.Peek(bodySize)
	if err == io.EOF {
		err = nil
	}
	if err != nil {
		return "", rc
	}
	return string(bodyLimited), rc
}
