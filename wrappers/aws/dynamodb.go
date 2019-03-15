package thundraaws

import (
	"github.com/thundra-io/thundra-lambda-agent-go/tracer"
	"github.com/aws/aws-sdk-go/aws/request"
	opentracing "github.com/opentracing/opentracing-go"
)

type dynamodbIntegration struct{}

func (i *dynamodbIntegration) getOperationName(r *request.Request) string {
	return "DynamoSpan"
}

func (i *dynamodbIntegration) beforeCall(r *request.Request, span opentracing.Span) {
	rawSpan, _ := tracer.GetRaw(span)
	// TODO: Check if ok is false
	rawSpan.ClassName = "AWS-DynamoDB"
	rawSpan.DomainName = "DB"
	return
}

func (i *dynamodbIntegration) afterCall(r *request.Request, span opentracing.Span) {
	return
}

func init() {
	integrations["DynamoDB"] = &dynamodbIntegration{}
}
