package thundraaws

import (
	"encoding/json"
	"strings"

	"github.com/thundra-io/thundra-lambda-agent-go/constants"

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
	rawSpan, ok := tracer.GetRaw(span)
	if !ok {
		return
	}

	rawSpan.ClassName = "AWS-DynamoDB"
	rawSpan.DomainName = "DB"

	operationName := r.Operation.Name
	operationType := constants.DynamoDBRequestTypes[operationName]
	endpoint := strings.SplitN(r.ClientInfo.Endpoint, "://", 2)[1]

	span.SetTag(constants.SpanTags["OPERATION_TYPE"], operationType)
	span.SetTag("db.host", endpoint)
	span.SetTag(constants.SpanTags["TOPOLOGY_VERTEX"], true)
	return
}

func (i *dynamodbIntegration) afterCall(r *request.Request, span opentracing.Span) {
	return
}

func init() {
	integrations["DynamoDB"] = &dynamodbIntegration{}
}
