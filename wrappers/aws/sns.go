package thundraaws

import (
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/thundra-io/thundra-lambda-agent-go/tracer"
)

type snsIntegration struct{}

func (i *snsIntegration) getOperationName(r *request.Request) string {
	return "SNSSpan"
}

func (i *snsIntegration) beforeCall(r *request.Request, span *tracer.RawSpan) {
	return
}

func (i *snsIntegration) afterCall(r *request.Request, span *tracer.RawSpan) {
	return
}

func init() {
	integrations["SNS"] = &snsIntegration{}
}
