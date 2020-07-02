package thundraaws

import (
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/thundra-io/thundra-lambda-agent-go/application"
	"github.com/thundra-io/thundra-lambda-agent-go/config"
	"github.com/thundra-io/thundra-lambda-agent-go/constants"
	"github.com/thundra-io/thundra-lambda-agent-go/tracer"
)

type sesIntegration struct{}
type sesData struct {
	Data string
	Charset string
}
type sesBody struct {
	Text sesData
	Html sesData
}
type sesParams struct {
	Source string
	Destination []string
	Subject sesData
	Body sesBody
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
		Body sesBody
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
	switch operationName := i.getOperationName(r); operationName {
	case "SendEmail":
		params := &sesSendEmailParams{}
		m, err := json.Marshal(r.Params)
		if err != nil { return &sesParams{} }
		if err = json.Unmarshal(m, &params); err != nil { return &sesParams{} }
		fields.Source = params.Source
		fields.Destination = params.Destination.ToAddresses
		fields.Subject = params.Message.Subject
		fields.Body = params.Message.Body
	case "SendRawEmail":
		params := &sesSendRawEmailParams{}
		m, err := json.Marshal(r.Params)
		if err != nil { return &sesParams{} }
		if err = json.Unmarshal(m, &params); err != nil { return &sesParams{} }
		fields.Source = params.Source
		fields.Destination = params.Destinations
	case "SendTemplatedEmail":
		params := &sesSendTemplatedEmailParams{}
		m, err := json.Marshal(r.Params)
		if err != nil { return &sesParams{} }
		if err = json.Unmarshal(m, &params); err != nil { return &sesParams{} }
		fields.Source = params.Source
		fields.Destination = params.Destination.ToAddresses
		fields.Template = params.Template
		fields.TemplateArn = params.TemplateArn
		fields.TemplateData = params.TemplateData
	}
	return fields
}

func (i *sesIntegration) getOperationName(r *request.Request) string {
	if r.Operation.Name != "" {
		return r.Operation.Name
	} else {
		return constants.AWSServiceRequest
	}
}

func (i *sesIntegration) beforeCall(r *request.Request, span *tracer.RawSpan) {
	span.ClassName = constants.ClassNames["SES"]
	span.DomainName = constants.DomainNames["MESSAGING"]

	operationName := r.Operation.Name
	operationType := getOperationType(operationName, constants.ClassNames["SES"])

	sesInfo := i.getSesInfo(r)

	tags := map[string]interface{}{
		constants.SpanTags["OPERATION_TYPE"]:          operationType,
		constants.AwsSDKTags["REQUEST_NAME"]:          operationName,
		constants.SpanTags["TOPOLOGY_VERTEX"]:         true,
		constants.SpanTags["TRIGGER_OPERATION_NAMES"]: []string{application.FunctionName},
		constants.SpanTags["TRIGGER_DOMAIN_NAME"]:     constants.AwsLambdaApplicationDomain,
		constants.SpanTags["TRIGGER_CLASS_NAME"]:      constants.AwsLambdaApplicationClass,
	}

	if operationName == "SendEmail" && !config.MaskSESMail {
		if sesInfo.Subject.Data != "" { tags[constants.AwsSESTags["SUBJECT"]] = sesInfo.Subject }
		if sesInfo.Body.Text.Data != "" || sesInfo.Body.Html.Data != "" {
			tags[constants.AwsSESTags["BODY"]] = sesInfo.Body
		}
	}

	if operationName == "SendTemplatedEmail" {
		if sesInfo.Template != "" { tags[constants.AwsSESTags["TEMPLATE_NAME"]] = sesInfo.Template }
		if sesInfo.TemplateArn != "" { tags[constants.AwsSESTags["TEMPLATE_ARN"]] = sesInfo.TemplateArn }
		if sesInfo.TemplateData != "" && !config.MaskSESMail {
			tags[constants.AwsSESTags["TEMPLATE_DATA"]] = sesInfo.TemplateData
		}
	}

	if sesInfo.Source != "" { tags[constants.AwsSESTags["SOURCE"]] = sesInfo.Source }
	if len(sesInfo.Destination) > 0 && !config.MaskSESDestination {
		tags[constants.AwsSESTags["DESTINATION"]] = sesInfo.Destination
	}

	span.Tags = tags
}

func (i *sesIntegration) afterCall(r *request.Request, span *tracer.RawSpan) {

}

func init() {
	integrations["SES"] = &sesIntegration{}
}
