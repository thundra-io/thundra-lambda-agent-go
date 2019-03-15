package thundraaws

import (
	"github.com/aws/aws-sdk-go/aws/request"
	opentracing "github.com/opentracing/opentracing-go"
)

type integration interface {
	beforeCall(r *request.Request, span opentracing.Span)
	afterCall(r *request.Request, span opentracing.Span)
	getOperationName(r *request.Request) string
}

var integrations = make(map[string]integration, 8)
