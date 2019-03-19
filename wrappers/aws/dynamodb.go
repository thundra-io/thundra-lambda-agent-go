package thundraaws

import (
	"encoding/json"
	"strings"

	"github.com/thundra-io/thundra-lambda-agent-go/application"

	"github.com/thundra-io/thundra-lambda-agent-go/constants"

	"github.com/aws/aws-sdk-go/aws/request"
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

func (i *dynamodbIntegration) beforeCall(r *request.Request, span *tracer.RawSpan) {
	span.ClassName = constants.ClassNames["DYNAMODB"]
	span.DomainName = constants.DomainNames["DB"]

	operationName := r.Operation.Name
	operationType := constants.DynamoDBRequestTypes[operationName]
	endpoint := r.ClientInfo.Endpoint
	endpointParts := strings.SplitN(endpoint, "://", 2)
	if len(endpointParts) > 1 {
		endpoint = endpointParts[1]
	}
	tags := map[string]interface{}{
		constants.SpanTags["OPERATION_TYPE"]:          operationType,
		constants.DBTags["DB_INSTANCE"]:               endpoint,
		constants.AwsDynamoDBTags["TABLE_NAME"]:       i.getTableName(r),
		constants.DBTags["DB_STATEMENT_TYPE"]:         operationType,
		constants.AwsSDKTags["REQUEST_NAME"]:          operationName,
		constants.SpanTags["TOPOLOGY_VERTEX"]:         true,
		constants.SpanTags["TRIGGER_OPERATION_NAMES"]: []string{application.FunctionName},
		constants.SpanTags["TRIGGER_DOMAIN_NAME"]:     constants.AwsLambdaApplicationDomain,
		constants.SpanTags["TRIGGER_CLASS_NAME"]:      constants.AwsLambdaApplicationClass,
	}

	span.Tags = tags
	// TODO: Get Key and Item values from request in a safe way to set db statement
	return
}

func (i *dynamodbIntegration) afterCall(r *request.Request, span *tracer.RawSpan) {
	return
}

func init() {
	integrations["DynamoDB"] = &dynamodbIntegration{}
}
