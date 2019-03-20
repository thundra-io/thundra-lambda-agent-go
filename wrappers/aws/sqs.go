package thundraaws

import (
	"encoding/json"
	"strings"

	"github.com/thundra-io/thundra-lambda-agent-go/application"
	"github.com/thundra-io/thundra-lambda-agent-go/constants"

	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/thundra-io/thundra-lambda-agent-go/tracer"
)

type sqsIntegration struct{}

func (i *sqsIntegration) getQueueName(r *request.Request) string {
	fields := struct {
		QueueName string `json:"QueueName"`
		QueueURL  string `json:"QueueUrl"`
	}{}
	m, err := json.Marshal(r.Params)
	if err != nil {
		return ""
	}
	if err = json.Unmarshal(m, &fields); err != nil {
		return ""
	}
	if len(fields.QueueName) > 0 {
		return fields.QueueName
	} else if len(fields.QueueURL) > 0 {
		urlParts := strings.Split(fields.QueueURL, "/")
		return urlParts[len(urlParts)-1]
	}
	return ""
}

func (i *sqsIntegration) getOperationName(r *request.Request) string {
	queueName := i.getQueueName(r)
	if len(queueName) > 0 {
		return queueName
	}
	return constants.AWSServiceRequest
}

func (i *sqsIntegration) beforeCall(r *request.Request, span *tracer.RawSpan) {
	span.ClassName = constants.ClassNames["SQS"]
	span.DomainName = constants.DomainNames["MESSAGING"]

	operationName := r.Operation.Name
	operationType := constants.SQSRequestTypes[operationName]

	tags := map[string]interface{}{
		constants.AwsSQSTags["QUEUE_NAME"]:            i.getQueueName(r),
		constants.SpanTags["OPERATION_TYPE"]:          operationType,
		constants.AwsSDKTags["REQUEST_NAME"]:          operationName,
		constants.SpanTags["TOPOLOGY_VERTEX"]:         true,
		constants.SpanTags["TRIGGER_OPERATION_NAMES"]: []string{application.FunctionName},
		constants.SpanTags["TRIGGER_DOMAIN_NAME"]:     constants.AwsLambdaApplicationDomain,
		constants.SpanTags["TRIGGER_CLASS_NAME"]:      constants.AwsLambdaApplicationClass,
	}

	span.Tags = tags
}

func (i *sqsIntegration) afterCall(r *request.Request, span *tracer.RawSpan) {
	return
}

func init() {
	integrations["SQS"] = &sqsIntegration{}
}
