package thundraaws

import (
	"github.com/aws/aws-sdk-go/aws/request"
	opentracing "github.com/opentracing/opentracing-go"
)

type kinesisIntegration struct{}

func (i *kinesisIntegration) getOperationName(r *request.Request) string {
	return "KinesisSpan"
}

func (i *kinesisIntegration) beforeCall(r *request.Request, span opentracing.Span) {
	return
}

func (i *kinesisIntegration) afterCall(r *request.Request, span opentracing.Span) {
	return
}

func init() {
	integrations["Kinesis"] = &kinesisIntegration{}
}
