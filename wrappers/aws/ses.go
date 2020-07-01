package thundraaws

import (
	// "encoding/json"
	"fmt"
	"github.com/thundra-io/thundra-lambda-agent-go/application"
	"github.com/thundra-io/thundra-lambda-agent-go/constants"
	"github.com/thundra-io/thundra-lambda-agent-go/tracer"
	"github.com/aws/aws-sdk-go/aws/request"
)

type sesIntegration struct{}
type sesData struct {
	Data string
	Charset string
}
type sesParams struct {
	Source string
	Destination []string
	Subject sesData
	Body sesData
	Template string
	TemplateArn string
	TemplateData string
}
type sesSendEmailParams struct {
	Source string
	Destination struct{
		ToAddresses []string
	}
	Message struct{
		Body struct{
			Text sesData
			Html sesData
		}
		Subject sesData
	}
}
type sesSendRawEmailParams struct {
	Source string
	Destinations []string
	RawMessage sesData
}
type sesSendTemplatedEmailParams struct {
	Source string
	Destination struct{
		ToAddresses []string
	}
	Template string
	TemplateArn string
	TemplateData string
}

func (i *sesIntegration) getSesInfo(r *request.Request) *sesParams {
	fields := &sesParams{}
	print(r)
	return fields
}

func (i *sesIntegration) getOperationName(r *request.Request) string {
	return constants.AWSServiceRequest
}

func (i *sesIntegration) beforeCall(r *request.Request, span *tracer.RawSpan) {
	span.ClassName = constants.ClassNames["SES"]
	span.DomainName = constants.DomainNames["STORAGE"]

	operationName := r.Operation.Name
	operationType := getOperationType(operationName, constants.ClassNames["SES"])

	tags := map[string]interface{}{
		constants.SpanTags["OPERATION_TYPE"]:          operationType,
		constants.AwsSDKTags["REQUEST_NAME"]:          operationName,
		constants.SpanTags["TOPOLOGY_VERTEX"]:         true,
		constants.SpanTags["TRIGGER_OPERATION_NAMES"]: []string{application.FunctionName},
		constants.SpanTags["TRIGGER_DOMAIN_NAME"]:     constants.AwsLambdaApplicationDomain,
		constants.SpanTags["TRIGGER_CLASS_NAME"]:      constants.AwsLambdaApplicationClass,
	}


	span.Tags = tags
}

func (i *sesIntegration) afterCall(r *request.Request, span *tracer.RawSpan) {

}

func init() {
	fmt.Print("SES integration init");
	integrations["SES"] = &sesIntegration{}
}
