package thundraaws

import (
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/thundra-io/thundra-lambda-agent-go/tracer"
)

type s3Integration struct{}

func (i *s3Integration) getOperationName(r *request.Request) string {
	return "S3Span"
}

func (i *s3Integration) beforeCall(r *request.Request, span *tracer.RawSpan) {
	return
}

func (i *s3Integration) afterCall(r *request.Request, span *tracer.RawSpan) {
	return
}

func init() {
	integrations["S3"] = &s3Integration{}
}
