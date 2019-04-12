package thundraaws

import (
	"encoding/json"
	"reflect"
	"strings"

	"github.com/thundra-io/thundra-lambda-agent-go/application"
	"github.com/thundra-io/thundra-lambda-agent-go/constants"
	"github.com/thundra-io/thundra-lambda-agent-go/utils"

	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/thundra-io/thundra-lambda-agent-go/tracer"
)

type sqsIntegration struct{}

func (i *sqsIntegration) getQueueName(r *request.Request) string {
	fields := struct {
		QueueName string
		QueueURL  string
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
	responseValue := reflect.ValueOf(r.Data).Elem()
	links := i.getTraceLinks(r.Operation.Name, &responseValue)
	if links != nil {
		span.Tags[constants.SpanTags["TRACE_LINKS"]] = links
	}
}

func (i *sqsIntegration) getTraceLinks(operationName string, responseValue *reflect.Value) []string {
	if operationName == "SendMessage" {
		messageID, _ := utils.GetStringFieldFromValue(*responseValue, "MessageId")
		return []string{messageID}

	} else if operationName == "SendMessageBatch" {
		successful := responseValue.FieldByName("Successful")
		if successful != (reflect.Value{}) && successful.Len() > 0 {
			var links []string
			for i := 0; i < successful.Len(); i++ {
				messageID, _ := utils.GetStringFieldFromValue(successful.Index(i).Elem(), "MessageId")
				links = append(links, messageID)
			}
			return links
		}
	}
	return nil
}

func init() {
	integrations["SQS"] = &sqsIntegration{}
}
