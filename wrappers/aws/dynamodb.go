package thundraaws

import (
	"github.com/thundra-io/thundra-lambda-agent-go/config"
	"encoding/json"
	"strings"

	"github.com/thundra-io/thundra-lambda-agent-go/application"

	"github.com/thundra-io/thundra-lambda-agent-go/constants"

	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/thundra-io/thundra-lambda-agent-go/tracer"
)

type dynamodbIntegration struct{}
type dynamodbParams struct {
	TableName string
	Key       map[string]map[string]interface{}
	Item      map[string]map[string]interface{}
}

func (i *dynamodbIntegration) getDynamodbInfo(r *request.Request) *dynamodbParams {
	fields := &dynamodbParams{}
	m, err := json.Marshal(r.Params)
	if err != nil {
		return &dynamodbParams{}
	}
	if err = json.Unmarshal(m, &fields); err != nil {
		return &dynamodbParams{}
	}
	return fields
}

func (i *dynamodbIntegration) getOperationName(r *request.Request) string {
	dynamodbInfo := i.getDynamodbInfo(r)
	if dynamodbInfo.TableName != "" {
		return dynamodbInfo.TableName
	}
	return constants.AWSServiceRequest
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
	dynamodbInfo := i.getDynamodbInfo(r)
	tags := map[string]interface{}{
		constants.SpanTags["OPERATION_TYPE"]:          operationType,
		constants.DBTags["DB_INSTANCE"]:               endpoint,
		constants.AwsDynamoDBTags["TABLE_NAME"]:       dynamodbInfo.TableName,
		constants.DBTags["DB_STATEMENT_TYPE"]:         operationType,
		constants.AwsSDKTags["REQUEST_NAME"]:          operationName,
		constants.SpanTags["TOPOLOGY_VERTEX"]:         true,
		constants.SpanTags["TRIGGER_OPERATION_NAMES"]: []string{application.FunctionName},
		constants.SpanTags["TRIGGER_DOMAIN_NAME"]:     constants.AwsLambdaApplicationDomain,
		constants.SpanTags["TRIGGER_CLASS_NAME"]:      constants.AwsLambdaApplicationClass,
	}

	span.Tags = tags
	
	if !config.MaskDynamoDBStatement {
		if len(dynamodbInfo.Item) > 0 {
			tags[constants.DBTags["DB_STATEMENT"]] = dynamodbInfo.Item
		} else if len(dynamodbInfo.Key) > 0 {
			tags[constants.DBTags["DB_STATEMENT"]] = dynamodbInfo.Key
		}
	}
}

func (i *dynamodbIntegration) afterCall(r *request.Request, span *tracer.RawSpan) {
	return
}

func init() {
	integrations["DynamoDB"] = &dynamodbIntegration{}
}
