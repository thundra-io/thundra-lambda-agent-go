package thundraaws

import (
	"encoding/json"
	"reflect"
	"strings"

	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/thundra-io/thundra-lambda-agent-go/application"
	"github.com/thundra-io/thundra-lambda-agent-go/constants"
	"github.com/thundra-io/thundra-lambda-agent-go/tracer"
	"github.com/thundra-io/thundra-lambda-agent-go/utils"
)

type snsIntegration struct{}

func (i *snsIntegration) getTopicName(r *request.Request) string {
	fields := struct {
		Name      string
		TopicArn  string
		TargetArn string
	}{}
	m, err := json.Marshal(r.Params)
	if err != nil {
		return ""
	}
	if err = json.Unmarshal(m, &fields); err != nil {
		return ""
	}
	if len(fields.Name) > 0 {
		return fields.Name
	} else if len(fields.TopicArn) > 0 {
		arnParts := strings.Split(fields.TopicArn, ":")
		return arnParts[len(arnParts)-1]
	} else if len(fields.TargetArn) > 0 {
		arnParts := strings.Split(fields.TargetArn, ":")
		return arnParts[len(arnParts)-1]
	}
	return ""
}

func (i *snsIntegration) getOperationName(r *request.Request) string {
	topicName := i.getTopicName(r)
	if len(topicName) > 0 {
		return topicName
	}
	return constants.AWSServiceRequest
}

func (i *snsIntegration) beforeCall(r *request.Request, span *tracer.RawSpan) {
	span.ClassName = constants.ClassNames["SNS"]
	span.DomainName = constants.DomainNames["MESSAGING"]

	operationName := r.Operation.Name
	operationType := constants.SNSRequestTypes[operationName]

	tags := map[string]interface{}{
		constants.AwsSNSTags["TOPIC_NAME"]:            i.getTopicName(r),
		constants.SpanTags["OPERATION_TYPE"]:          operationType,
		constants.AwsSDKTags["REQUEST_NAME"]:          operationName,
		constants.SpanTags["TOPOLOGY_VERTEX"]:         true,
		constants.SpanTags["TRIGGER_OPERATION_NAMES"]: []string{application.FunctionName},
		constants.SpanTags["TRIGGER_DOMAIN_NAME"]:     constants.AwsLambdaApplicationDomain,
		constants.SpanTags["TRIGGER_CLASS_NAME"]:      constants.AwsLambdaApplicationClass,
	}

	span.Tags = tags
}

func (i *snsIntegration) afterCall(r *request.Request, span *tracer.RawSpan) {
	responseValue := reflect.ValueOf(r.Data).Elem()
	messageID, _ := utils.GetStringFieldFromValue(responseValue, "MessageId")

	span.Tags[constants.SpanTags["TRACE_LINKS"]] = []string{messageID}
}

func init() {
	integrations["SNS"] = &snsIntegration{}
}
