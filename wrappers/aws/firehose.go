package thundraaws

import (
	"encoding/json"

	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/thundra-io/thundra-lambda-agent-go/application"
	"github.com/thundra-io/thundra-lambda-agent-go/constants"
	"github.com/thundra-io/thundra-lambda-agent-go/tracer"
)

type firehoseIntegration struct{}

func (i *firehoseIntegration) getDeliveryStreamName(r *request.Request) string {
	fields := struct {
		DeliveryStreamName string `json:"DeliveryStreamName"`
	}{}
	m, err := json.Marshal(r.Params)
	if err != nil {
		return ""
	}
	if err = json.Unmarshal(m, &fields); err != nil {
		return ""
	}
	if len(fields.DeliveryStreamName) > 0 {
		return fields.DeliveryStreamName
	}
	return ""
}

func (i *firehoseIntegration) getOperationName(r *request.Request) string {
	return i.getDeliveryStreamName(r)
}

func (i *firehoseIntegration) beforeCall(r *request.Request, span *tracer.RawSpan) {
	span.ClassName = constants.ClassNames["FIREHOSE"]
	span.DomainName = constants.DomainNames["STREAM"]

	operationName := r.Operation.Name
	operationType := constants.FirehoseRequestTypes[operationName]

	tags := map[string]interface{}{
		constants.AwsFirehoseTags["STREAM_NAME"]:      i.getDeliveryStreamName(r),
		constants.SpanTags["OPERATION_TYPE"]:          operationType,
		constants.AwsSDKTags["REQUEST_NAME"]:          operationName,
		constants.SpanTags["TOPOLOGY_VERTEX"]:         true,
		constants.SpanTags["TRIGGER_OPERATION_NAMES"]: []string{application.FunctionName},
		constants.SpanTags["TRIGGER_DOMAIN_NAME"]:     constants.AwsLambdaApplicationDomain,
		constants.SpanTags["TRIGGER_CLASS_NAME"]:      constants.AwsLambdaApplicationClass,
	}

	span.Tags = tags
}

func (i *firehoseIntegration) afterCall(r *request.Request, span *tracer.RawSpan) {
	return
}

func init() {
	integrations["Firehose"] = &firehoseIntegration{}
}
