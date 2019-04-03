package thundraaws

import (
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/thundra-io/thundra-lambda-agent-go/tracer"
)

type integration interface {
	beforeCall(r *request.Request, span *tracer.RawSpan)
	afterCall(r *request.Request, span *tracer.RawSpan)
	getOperationName(r *request.Request) string
}

var integrations = make(map[string]integration, 8)
