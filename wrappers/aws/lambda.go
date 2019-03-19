package thundraaws

import (
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/thundra-io/thundra-lambda-agent-go/tracer"
)

type lambdaIntegration struct{}

func (i *lambdaIntegration) getOperationName(r *request.Request) string {
	return "LambdaSpan"
}

func (i *lambdaIntegration) beforeCall(r *request.Request, span *tracer.RawSpan) {
	return
}

func (i *lambdaIntegration) afterCall(r *request.Request, span *tracer.RawSpan) {
	return
}

func init() {
	integrations["Lambda"] = &lambdaIntegration{}
}
