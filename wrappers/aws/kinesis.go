package thundraaws

import (
	"encoding/json"

	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/thundra-io/thundra-lambda-agent-go/application"
	"github.com/thundra-io/thundra-lambda-agent-go/constants"
	"github.com/thundra-io/thundra-lambda-agent-go/tracer"
)

type kinesisIntegration struct{}

func (i *kinesisIntegration) getStreamName(r *request.Request) string {
	fields := struct {
		StreamName string `json:"StreamName"`
	}{}
	m, err := json.Marshal(r.Params)
	if err != nil {
		return ""
	}
	if err = json.Unmarshal(m, &fields); err != nil {
		return ""
	}
	if len(fields.StreamName) > 0 {
		return fields.StreamName
	}
	return ""
}

func (i *kinesisIntegration) getOperationName(r *request.Request) string {
	streamName := i.getStreamName(r)
	if len(streamName) > 0 {
		return streamName
	}
	return constants.AWSServiceRequest
}

func (i *kinesisIntegration) beforeCall(r *request.Request, span *tracer.RawSpan) {
	span.ClassName = constants.ClassNames["KINESIS"]
	span.DomainName = constants.DomainNames["STREAM"]

	operationName := r.Operation.Name
	operationType := constants.KinesisRequestTypes[operationName]

	tags := map[string]interface{}{
		constants.AwsKinesisTags["STREAM_NAME"]:       i.getStreamName(r),
		constants.SpanTags["OPERATION_TYPE"]:          operationType,
		constants.AwsSDKTags["REQUEST_NAME"]:          operationName,
		constants.SpanTags["TOPOLOGY_VERTEX"]:         true,
		constants.SpanTags["TRIGGER_OPERATION_NAMES"]: []string{application.FunctionName},
		constants.SpanTags["TRIGGER_DOMAIN_NAME"]:     constants.AwsLambdaApplicationDomain,
		constants.SpanTags["TRIGGER_CLASS_NAME"]:      constants.AwsLambdaApplicationClass,
	}

	span.Tags = tags
}

func (i *kinesisIntegration) afterCall(r *request.Request, span *tracer.RawSpan) {
	return
}

func init() {
	integrations["Kinesis"] = &kinesisIntegration{}
}
