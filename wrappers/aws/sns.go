package thundraaws

import (
	"github.com/aws/aws-sdk-go/aws/request"
	opentracing "github.com/opentracing/opentracing-go"
)

type snsIntegration struct{}

func (i *snsIntegration) getOperationName(r *request.Request) string {
	return "SNSSpan"
}

func (i *snsIntegration) beforeCall(r *request.Request, span opentracing.Span) {
	return
}

func (i *snsIntegration) afterCall(r *request.Request, span opentracing.Span) {
	return
}

func init() {
	integrations["SNS"] = &snsIntegration{}
}
