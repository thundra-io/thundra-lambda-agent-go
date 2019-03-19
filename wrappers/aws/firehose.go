package thundraaws

import (
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/thundra-io/thundra-lambda-agent-go/tracer"
)

type firehoseIntegration struct{}

func (i *firehoseIntegration) getOperationName(r *request.Request) string {
	return "FirehoseSpan"
}

func (i *firehoseIntegration) beforeCall(r *request.Request, span *tracer.RawSpan) {
	return
}

func (i *firehoseIntegration) afterCall(r *request.Request, span *tracer.RawSpan) {
	return
}

func init() {
	integrations["Firehose"] = &firehoseIntegration{}
}
