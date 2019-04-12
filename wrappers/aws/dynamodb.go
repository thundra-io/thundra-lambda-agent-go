package thundraaws

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/dynamodb"

	"github.com/thundra-io/thundra-lambda-agent-go/application"
	"github.com/thundra-io/thundra-lambda-agent-go/config"
	"github.com/thundra-io/thundra-lambda-agent-go/constants"
	"github.com/thundra-io/thundra-lambda-agent-go/tracer"
	"github.com/thundra-io/thundra-lambda-agent-go/utils"
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

	if config.DynamoDBTraceInjectionEnabled {
		if operationName == "PutItem" {
			i.injectTraceLinkOnPut(r, span)
		} else if operationName == "UpdateItem" {
			i.injectTraceLinkOnUpdate(r, span)
		} else if operationName == "DeleteItem" {
			i.injectTraceLinkOnDelete(r, span)
		}
	}
}

func (i *dynamodbIntegration) afterCall(r *request.Request, span *tracer.RawSpan) {
	if span.GetTag(constants.SpanTags["TRACE_LINKS"]) != nil {
		return
	}
	traceLinks := i.getTraceLinks(r, span)
	if traceLinks != nil {
		span.Tags[constants.SpanTags["TRACE_LINKS"]] = traceLinks
	}
}

func (i *dynamodbIntegration) getTraceLinks(r *request.Request, span *tracer.RawSpan) []string {

	region := ""
	dateStr := ""

	if r.Config.Region != nil {
		region = *r.Config.Region
	}
	if r.HTTPResponse != nil && r.HTTPResponse.Header != nil {
		dateStr = r.HTTPResponse.Header.Get("date")
	}
	dynamodbInfo := i.getDynamodbInfo(r)

	params := reflect.ValueOf(&r.Params).Elem().Interface()
	if r.Operation.Name == "PutItem" {
		switch v := params.(type) {
		case *dynamodb.PutItemInput:
			if v.Item != nil {
				strAttributes := attributesToStr(v.Item)
				return i.generateTraceLinks(region, dateStr, "PUT", dynamodbInfo.TableName, strAttributes)
			}
		}
	} else if r.Operation.Name == "UpdateItem" {
		switch v := params.(type) {
		case *dynamodb.UpdateItemInput:
			if v.Key != nil {
				strAttributes := attributesToStr(v.Key)
				return i.generateTraceLinks(region, dateStr, "UPDATE", dynamodbInfo.TableName, strAttributes)
			}
		}
	} else if r.Operation.Name == "DeleteItem" {

		attributes := reflect.ValueOf(r.Data).Elem().FieldByName("Attributes")
		if attributes != (reflect.Value{}) {
			if attributeValues, ok := attributes.Interface().(map[string]*dynamodb.AttributeValue); ok {
				if attributeValues != nil {
					spanAttr := attributeValues["x-thundra-span-id"]
					if spanAttr.S != nil {
						return []string{"DELETE:" + *spanAttr.S}
					}
				}
			}
		}
		switch v := params.(type) {
		case *dynamodb.DeleteItemInput:
			if v.Key != nil {
				strAttributes := attributesToStr(v.Key)
				return i.generateTraceLinks(region, dateStr, "DELETE", dynamodbInfo.TableName, strAttributes)
			}
		}
	}

	return nil
}

func attributesToStr(attr map[string]*dynamodb.AttributeValue) string {
	var keys []string
	for k := range attr {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	attributesStr := ""
	first := true
	for _, key := range keys {
		attrValue := attr[key]
		valueStr, err := utils.AttributeValuetoStr(*attrValue)
		if err == nil {
			if !first {
				attributesStr += ", "

			} else {
				first = false
			}
			attributesStr += key + "="
			attributesStr += valueStr
		}
	}
	return attributesStr
}

func (i *dynamodbIntegration) generateTraceLinks(region string, dateStr string, operationType string, tableName string, attributesStr string) []string {
	var traceLinks []string
	timestamp := getTimeStamp(dateStr)

	b := md5.Sum([]byte(attributesStr))
	dataMD5 := hex.EncodeToString(b[:])

	for j := 0; j < 3; j++ {
		traceLinks = append(traceLinks, region+":"+tableName+":"+strconv.FormatInt(timestamp+int64(j), 10)+":"+operationType+":"+dataMD5)
	}

	return traceLinks
}

func (i *dynamodbIntegration) injectTraceLinkOnPut(r *request.Request, span *tracer.RawSpan) {
	paramsValue := reflect.ValueOf(&r.Params)
	if paramsValue == (reflect.Value{}) {
		return
	}
	params := paramsValue.Elem().Interface()
	switch v := params.(type) {
	case *dynamodb.PutItemInput:
		thundraSpanAttr := &dynamodb.AttributeValue{
			S: aws.String(span.Context.SpanID),
		}
		v.Item["x-thundra-span-id"] = thundraSpanAttr
		span.Tags[constants.SpanTags["TRACE_LINKS"]] = []string{"SAVE:" + span.Context.SpanID}
	}
}

func (i *dynamodbIntegration) injectTraceLinkOnDelete(r *request.Request, span *tracer.RawSpan) {
	paramsValue := reflect.ValueOf(&r.Params)
	if paramsValue == (reflect.Value{}) {
		return
	}
	params := paramsValue.Elem().Interface()
	switch v := params.(type) {
	case *dynamodb.DeleteItemInput:
		v.ReturnValues = aws.String("ALL_OLD")
	}
}

func (i *dynamodbIntegration) injectTraceLinkOnUpdate(r *request.Request, span *tracer.RawSpan) {
	paramsValue := reflect.ValueOf(&r.Params)
	if paramsValue == (reflect.Value{}) {
		return
	}
	params := paramsValue.Elem().Interface()
	thundraSpanAttr := &dynamodb.AttributeValue{
		S: aws.String(span.Context.SpanID),
	}
	switch v := params.(type) {
	case *dynamodb.UpdateItemInput:
		if v.AttributeUpdates != nil {
			action := "PUT"
			attributeUpdate := dynamodb.AttributeValueUpdate{Action: &action, Value: thundraSpanAttr}
			v.AttributeUpdates["x-thundra-span-id"] = &attributeUpdate
			span.Tags[constants.SpanTags["TRACE_LINKS"]] = []string{"SAVE:" + span.Context.SpanID}

		} else if v.UpdateExpression != nil {
			if v.ExpressionAttributeNames != nil {
				v.ExpressionAttributeNames["#xThundraSpanId"] = aws.String("x-thundra-span-id")
			} else {
				v.ExpressionAttributeNames = map[string]*string{
					"#xThundraSpanId": aws.String("x-thundra-span-id"),
				}
			}

			if v.ExpressionAttributeValues != nil {
				v.ExpressionAttributeValues[":xThundraSpanId"] = thundraSpanAttr
			} else {
				v.ExpressionAttributeValues = map[string]*dynamodb.AttributeValue{
					":xThundraSpanId": thundraSpanAttr,
				}
			}
			if strings.Index((*v.UpdateExpression), "SET") == -1 {
				v.SetUpdateExpression("SET #xThundraSpanId = :xThundraSpanId " + *v.UpdateExpression)

			} else {
				pattern := regexp.MustCompile("SET (.)")
				repl := "SET #xThundraSpanId = :xThundraSpanId, ${1}$2"
				output := pattern.ReplaceAllString(*v.UpdateExpression, repl)
				v.SetUpdateExpression(output)
			}
			span.Tags[constants.SpanTags["TRACE_LINKS"]] = []string{"SAVE:" + span.Context.SpanID}
		}
	}
}

func init() {
	integrations["DynamoDB"] = &dynamodbIntegration{}
}
