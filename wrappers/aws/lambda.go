package thundraaws

import (
	"encoding/base64"
	"encoding/json"
	"strings"

	"github.com/thundra-io/thundra-lambda-agent-go/v2/application"
	"github.com/thundra-io/thundra-lambda-agent-go/v2/config"
	"github.com/thundra-io/thundra-lambda-agent-go/v2/constants"

	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/thundra-io/thundra-lambda-agent-go/v2/tracer"
)

type lambdaIntegration struct{}
type lambdaParams struct {
	FunctionName   string
	Qualifier      string
	InvocationType string
	Payload        string
	ClientContext  string
}

func (i *lambdaIntegration) getLambdaInfo(r *request.Request) *lambdaParams {
	fields := &lambdaParams{}
	m, err := json.Marshal(r.Params)
	if err != nil {
		return &lambdaParams{}
	}
	if err = json.Unmarshal(m, fields); err != nil {
		return &lambdaParams{}
	}
	return fields
}

func (i *lambdaIntegration) getOperationName(r *request.Request) string {
	lambdaInfo := i.getLambdaInfo(r)
	if len(lambdaInfo.FunctionName) > 0 {
		return i.getFunctionName(lambdaInfo.FunctionName)
	}
	return constants.AWSServiceRequest
}

func (i *lambdaIntegration) getFunctionName(name string) string {
	functionName := name
	pos := strings.LastIndex(name, ":function:")
	if pos != -1 {
		posAfter := pos + len(":function:")
		if posAfter >= len(name) {
			functionName = ""
		}
		functionName = name[posAfter:len(name)]
	}

	// Strip version number if exists
	pos = strings.IndexByte(functionName, ':')
	if pos != -1 && pos < len(functionName) {
		functionName = functionName[:pos]
	}
	return functionName
}

func (i *lambdaIntegration) beforeCall(r *request.Request, span *tracer.RawSpan) {
	span.ClassName = constants.ClassNames["LAMBDA"]
	span.DomainName = constants.DomainNames["API"]

	operationName := r.Operation.Name
	operationType := getOperationType(operationName, constants.ClassNames["LAMBDA"])

	lambdaInfo := i.getLambdaInfo(r)

	tags := map[string]interface{}{
		constants.AwsLambdaTags["FUNCTION_NAME"]:      i.getFunctionName(lambdaInfo.FunctionName),
		constants.SpanTags["OPERATION_TYPE"]:          operationType,
		constants.AwsSDKTags["REQUEST_NAME"]:          operationName,
		constants.SpanTags["TOPOLOGY_VERTEX"]:         true,
		constants.SpanTags["TRIGGER_OPERATION_NAMES"]: []string{application.FunctionName},
		constants.SpanTags["TRIGGER_DOMAIN_NAME"]:     constants.AwsLambdaApplicationDomain,
		constants.SpanTags["TRIGGER_CLASS_NAME"]:      constants.AwsLambdaApplicationClass,
	}

	if !config.MaskLambdaPayload && lambdaInfo.Payload != "" {
		tags[constants.AwsLambdaTags["INVOCATION_PAYLOAD"]] = lambdaInfo.Payload
	}
	if lambdaInfo.Qualifier != "" {
		tags[constants.AwsLambdaTags["FUNCTION_QUALIFIER"]] = lambdaInfo.Qualifier
	}
	if lambdaInfo.InvocationType != "" {
		tags[constants.AwsLambdaTags["INVOCATION_TYPE"]] = lambdaInfo.InvocationType
	}

	span.Tags = tags

	if !config.LambdaTraceInjectionDisabled {
		i.injectSpanIntoClientContext(r)
	}
}

func (i *lambdaIntegration) afterCall(r *request.Request, span *tracer.RawSpan) {
	xAmzRequestID := ""
	if r.HTTPResponse != nil && r.HTTPResponse.Header != nil {
		xAmzRequestID = r.HTTPResponse.Header.Get("X-Amzn-Requestid")
	}
	if xAmzRequestID != "" {
		span.Tags[constants.SpanTags["TRACE_LINKS"]] = []string{xAmzRequestID}
	}
}

func (i *lambdaIntegration) injectSpanIntoClientContext(r *request.Request) {
	input, ok := r.Params.(*lambda.InvokeInput)

	if !ok {
		return
	}
	clientContext := &lambdacontext.ClientContext{}
	if input.ClientContext != nil {
		data, err := base64.StdEncoding.DecodeString(*input.ClientContext)
		if err != nil {
			return
		}
		if err = json.Unmarshal(data, clientContext); err != nil {
			return
		}
	}
	if clientContext.Custom == nil {
		clientContext.Custom = make(map[string]string, 3)
	}
	clientContext.Custom[constants.AwsLambdaTriggerOperationName] = application.ApplicationName
	clientContext.Custom[constants.AwsLambdaTriggerDomainName] = application.ApplicationDomainName
	clientContext.Custom[constants.AwsLambdaTriggerClassName] = application.ApplicationClassName

	clientContextJSON, err := json.Marshal(clientContext)
	if err != nil {
		return
	}

	encodedClientContextJSON := base64.StdEncoding.EncodeToString(clientContextJSON)
	input.ClientContext = &encodedClientContextJSON
}

func init() {
	integrations["Lambda"] = &lambdaIntegration{}
}
