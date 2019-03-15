package thundraaws

import (
	"github.com/aws/aws-sdk-go/aws/request"
	opentracing "github.com/opentracing/opentracing-go"
)

type lambdaIntegration struct{}

func (i *lambdaIntegration) getOperationName(r *request.Request) string {
	return "LambdaSpan"
}

func (i *lambdaIntegration) beforeCall(r *request.Request, span opentracing.Span) {
	return
}

func (i *lambdaIntegration) afterCall(r *request.Request, span opentracing.Span) {
	return
}

func init() {
	integrations["Lambda"] = &lambdaIntegration{}
}
