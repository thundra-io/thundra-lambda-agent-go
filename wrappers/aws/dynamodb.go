package thundraaws

import (
	"encoding/json"

	"github.com/aws/aws-sdk-go/aws/request"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/thundra-io/thundra-lambda-agent-go/tracer"
)

type dynamodbIntegration struct{}

func (i *dynamodbIntegration) getOperationName(r *request.Request) string {
	fields := struct {
		TableName string `json:"TableName"`
	}{}
	m, err := json.Marshal(r.Params)
	if err != nil {
		return ""
	}
	err = json.Unmarshal(m, &fields)
	if err != nil {
		return ""
	}
	return fields.TableName
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
