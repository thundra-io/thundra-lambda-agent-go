package thundraaws

import (
	"github.com/aws/aws-sdk-go/aws/request"
	opentracing "github.com/opentracing/opentracing-go"
)

type sqsIntegration struct{}

func (i *sqsIntegration) getOperationName(r *request.Request) string {
	return "SQSSpan"
}

func (i *sqsIntegration) beforeCall(r *request.Request, span opentracing.Span) {
	return
}

func (i *sqsIntegration) afterCall(r *request.Request, span opentracing.Span) {
	return
}

func init() {
	integrations["SQS"] = &sqsIntegration{}
}
