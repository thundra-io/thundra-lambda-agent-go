package elastic

import (
	"net/http"
	"strings"

	"github.com/thundra-io/thundra-lambda-agent-go/utils"

	"github.com/thundra-io/thundra-lambda-agent-go/application"
	"github.com/thundra-io/thundra-lambda-agent-go/config"
	"github.com/thundra-io/thundra-lambda-agent-go/constants"
	"github.com/thundra-io/thundra-lambda-agent-go/tracer"
)

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
		constants.EsTags["ES_NORMALIZED_URI"]:         GetNormalizedPath(req.URL.Path),
		constants.EsTags["ES_METHOD"]:                 method,
		constants.EsTags["ES_PARAMS"]:                 req.URL.Query().Encode(),
		constants.DBTags["DB_TYPE"]:                   "elasticsearch",
		constants.SpanTags["TRIGGER_DOMAIN_NAME"]:     constants.AwsLambdaApplicationDomain,
		constants.SpanTags["TRIGGER_CLASS_NAME"]:      constants.AwsLambdaApplicationClass,
		constants.SpanTags["TRIGGER_OPERATION_NAMES"]: []string{application.FunctionName},
		constants.SpanTags["TOPOLOGY_VERTEX"]:         true,
	}

	if req != nil && req.Body != nil && !config.MaskEsBody {
		esBody, req.Body = utils.ReadRequestBody(req.Body, int(req.ContentLength))
		tags[constants.EsTags["ES_BODY"]] = esBody
	}

	span.Tags = tags
}

func AfterCall(span *tracer.RawSpan, resp *http.Response) {

}

func GetNormalizedPath(path string) string {
	depth := config.EsIntegrationUrlPathDepth
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
