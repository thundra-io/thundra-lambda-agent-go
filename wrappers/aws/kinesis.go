package thundraaws

import (
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/thundra-io/thundra-lambda-agent-go/tracer"
)

type kinesisIntegration struct{}

func (i *kinesisIntegration) getOperationName(r *request.Request) string {
	return "KinesisSpan"
}

func (i *kinesisIntegration) beforeCall(r *request.Request, span *tracer.RawSpan) {
	return
}

func (i *kinesisIntegration) afterCall(r *request.Request, span *tracer.RawSpan) {
	return
}

func init() {
	integrations["Kinesis"] = &kinesisIntegration{}
}
