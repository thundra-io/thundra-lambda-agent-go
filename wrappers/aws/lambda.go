package thundraaws

import (
	"encoding/base64"
	"encoding/json"

	"github.com/thundra-io/thundra-lambda-agent-go/application"
	"github.com/thundra-io/thundra-lambda-agent-go/constants"

	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/thundra-io/thundra-lambda-agent-go/tracer"
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
		return lambdaInfo.FunctionName
	}
	return constants.AWSServiceRequest
}

func (i *lambdaIntegration) beforeCall(r *request.Request, span *tracer.RawSpan) {
	span.ClassName = constants.ClassNames["LAMBDA"]
	span.DomainName = constants.DomainNames["API"]

	operationName := r.Operation.Name
	operationType := constants.LambdaRequestTypes[operationName]

	lambdaInfo := i.getLambdaInfo(r)

	tags := map[string]interface{}{
		constants.AwsLambdaTags["FUNCTION_NAME"]:      lambdaInfo.FunctionName,
		constants.SpanTags["OPERATION_TYPE"]:          operationType,
		constants.AwsSDKTags["REQUEST_NAME"]:          operationName,
		constants.SpanTags["TOPOLOGY_VERTEX"]:         true,
		constants.SpanTags["TRIGGER_OPERATION_NAMES"]: []string{application.FunctionName},
		constants.SpanTags["TRIGGER_DOMAIN_NAME"]:     constants.AwsLambdaApplicationDomain,
		constants.SpanTags["TRIGGER_CLASS_NAME"]:      constants.AwsLambdaApplicationClass,
	}

	if lambdaInfo.Payload != "" {
		tags[constants.AwsLambdaTags["INVOCATION_PAYLOAD"]] = lambdaInfo.Payload
	}
	if lambdaInfo.Qualifier != "" {
		tags[constants.AwsLambdaTags["FUNCTION_QUALIFIER"]] = lambdaInfo.Qualifier
	}
	if lambdaInfo.InvocationType != "" {
		tags[constants.AwsLambdaTags["INVOCATION_TYPE"]] = lambdaInfo.InvocationType
	}

	span.Tags = tags
	i.injectSpanIntoClientContext(r)
}

func (i *lambdaIntegration) afterCall(r *request.Request, span *tracer.RawSpan) {
	xAmzRequestID := r.HTTPResponse.Header.Get("X-Amzn-Requestid")
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
