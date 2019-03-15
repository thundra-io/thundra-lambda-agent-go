package thundraaws

import (
	"github.com/aws/aws-sdk-go/aws/request"
	opentracing "github.com/opentracing/opentracing-go"
)

type s3Integration struct{}

func (i *s3Integration) getOperationName(r *request.Request) string {
	return "S3Span"
}

func (i *s3Integration) beforeCall(r *request.Request, span opentracing.Span) {
	return
}

func (i *s3Integration) afterCall(r *request.Request, span opentracing.Span) {
	return
}

func init() {
	integrations["S3"] = &s3Integration{}
}
