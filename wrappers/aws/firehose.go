package thundraaws

import (
	"github.com/aws/aws-sdk-go/aws/request"
	opentracing "github.com/opentracing/opentracing-go"
)

type firehoseIntegration struct{}

func (i *firehoseIntegration) getOperationName(r *request.Request) string {
	return "FirehoseSpan"
}

func (i *firehoseIntegration) beforeCall(r *request.Request, span opentracing.Span) {
	return
}

func (i *firehoseIntegration) afterCall(r *request.Request, span opentracing.Span) {
	return
}

func init() {
	integrations["Firehose"] = &firehoseIntegration{}
}
