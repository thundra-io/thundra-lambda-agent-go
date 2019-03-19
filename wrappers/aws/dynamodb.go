package thundraaws

import (
	"github.com/thundra-io/thundra-lambda-agent-go/application"
	"encoding/json"
	"strings"

	"github.com/thundra-io/thundra-lambda-agent-go/constants"

	"github.com/aws/aws-sdk-go/aws/request"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/thundra-io/thundra-lambda-agent-go/tracer"
)

type dynamodbIntegration struct{}

func (i *dynamodbIntegration) getTableName(r *request.Request) string {
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

func (i *dynamodbIntegration) getOperationName(r *request.Request) string {
	return i.getTableName(r)
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
	endpoint := r.ClientInfo.Endpoint
	endpointParts := strings.SplitN(endpoint, "://", 2)
	if len(endpointParts) > 1 {
		endpoint = endpointParts[1]
	}
	span.SetTag(constants.SpanTags["OPERATION_TYPE"], operationType)
	span.SetTag(constants.DBTags["DB_INSTANCE"], endpoint)
	span.SetTag(constants.AwsDynamoDBTags["TABLE_NAME"], i.getTableName(r))
	span.SetTag(constants.DBTags["DB_STATEMENT_TYPE"], operationType)
	span.SetTag(constants.AwsSDKTags["REQUEST_NAME"], operationName)

	// TODO: Get Key and Item values from request in a safe way to set db statement
	
	span.SetTag(constants.SpanTags["TOPOLOGY_VERTEX"], true)
	span.SetTag(constants.SpanTags["TRIGGER_OPERATION_NAMES"], []string{application.FunctionName})
	span.SetTag(constants.SpanTags["TRIGGER_DOMAIN_NAME"], constants.AwsLambdaApplicationDomain)
	span.SetTag(constants.SpanTags["TRIGGER_CLASS_NAME"], constants.AwsLambdaApplicationClass)
	return
}

func (i *dynamodbIntegration) afterCall(r *request.Request, span opentracing.Span) {
	return
}

func init() {
	integrations["DynamoDB"] = &dynamodbIntegration{}
}
