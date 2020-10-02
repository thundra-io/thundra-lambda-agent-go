package invocation

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/thundra-io/thundra-lambda-agent-go/v2/application"
	"github.com/thundra-io/thundra-lambda-agent-go/v2/constants"
	"github.com/thundra-io/thundra-lambda-agent-go/v2/utils"
)

type requestField struct {
	URI string
}

type cloudFrontField struct {
	Request requestField
}

type recordField struct {
	EventSource string
	CF          cloudFrontField
}

type paramsField struct {
	Header map[string]string
}

type requestContextField struct {
	Stage string
}

type cloudwatchLogsRawDataField struct {
	Data string
}

type triggerEvent struct {
	Records           []recordField
	RequestContext    requestContextField
	Params            paramsField
	AwsLogs           cloudwatchLogsRawDataField
	Context           map[string]string
	Headers           map[string]string
	HTTPMethod        string
	Path              string `json:"path"`
	DeliveryStreamArn string `json:"deliveryStreamArn"`
	DetailType        string `json:"detail-type"`
}

type key struct{}
type eventTypeKey key

var void struct{}

func injectTriggerTagsToInvocation(domainName string, className string, operationNames []string) {
	SetAgentTag(constants.SpanTags["TRIGGER_DOMAIN_NAME"], domainName)
	SetAgentTag(constants.SpanTags["TRIGGER_CLASS_NAME"], className)
	SetAgentTag(constants.SpanTags["TRIGGER_OPERATION_NAMES"], operationNames)
}

func injectTriggerTagsForDynamoDB(payload json.RawMessage) {
	e := events.DynamoDBEvent{}
	err := json.Unmarshal(payload, &e)
	if err != nil {
		return
	}

	domainName := constants.DomainNames["DB"]
	className := constants.ClassNames["DYNAMODB"]
	tableNamesMap := make(map[string]struct{})
	var traceLinks []string
	for _, record := range e.Records {
		tableName := ""
		if len(strings.Split(record.EventSourceArn, "/")) > 1 {
			tableName = strings.Split(record.EventSourceArn, "/")[1]
		}
		tableNamesMap[tableName] = void

		traceLinkFound := false
		if record.EventName == "INSERT" || record.EventName == "MODIFY" {
			if record.Change.NewImage != nil {
				if record.Change.NewImage["x-thundra-span-id"].DataType() == events.DataTypeString {
					traceLinks = append(traceLinks, "SAVE:"+record.Change.NewImage["x-thundra-span-id"].String())
					traceLinkFound = true
				}
			}
		} else if record.EventName == "REMOVE" {
			if record.Change.OldImage != nil {
				if record.Change.OldImage["x-thundra-span-id"].DataType() == events.DataTypeString {
					spanID := record.Change.OldImage["x-thundra-span-id"].String()
					traceLinks = append(traceLinks, "DELETE:"+spanID)
					traceLinkFound = true
				}
			}
		}
		if !traceLinkFound {
			creationTime := record.Change.ApproximateCreationDateTime
			if creationTime != (events.SecondsEpochTime{}) {
				if record.EventName == "INSERT" || record.EventName == "MODIFY" {
					if record.Change.NewImage != nil {
						attributeStr := attributesToStr(record.Change.NewImage)
						addDynamoDBTraceLinks(&traceLinks, &record, "SAVE", tableName, attributeStr)
					}
					if record.Change.Keys != nil {
						attributeStr := attributesToStr(record.Change.Keys)
						addDynamoDBTraceLinks(&traceLinks, &record, "SAVE", tableName, attributeStr)
					}
				} else if record.EventName == "REMOVE" {
					if record.Change.Keys != nil {
						attributeStr := attributesToStr(record.Change.Keys)
						addDynamoDBTraceLinks(&traceLinks, &record, "DELETE", tableName, attributeStr)
					}
				}
			}
		}
	}

	AddIncomingTraceLinks(traceLinks)
	var tableNames []string
	for k := range tableNamesMap {
		tableNames = append(tableNames, k)
	}
	injectTriggerTagsToInvocation(domainName, className, tableNames)
}

func addDynamoDBTraceLinks(traceLinks *[]string, record *events.DynamoDBEventRecord, operationType string, tableName string, attributesStr string) {
	creationTime := record.Change.ApproximateCreationDateTime.Unix() - 1
	region := record.AWSRegion

	b := md5.Sum([]byte(attributesStr))
	dataMD5 := hex.EncodeToString(b[:])

	for j := 0; j < 3; j++ {
		*traceLinks = append(*traceLinks, region+":"+tableName+":"+strconv.FormatInt(creationTime+int64(j), 10)+":"+operationType+":"+dataMD5)
	}
}

func attributesToStr(attr map[string]events.DynamoDBAttributeValue) string {
	var keys []string
	for k := range attr {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	attributesStr := ""
	first := true
	for _, key := range keys {
		dynamoAttrValue := attr[key]

		valueStr, err := utils.AttributeValuetoStr(dynamoAttrValue)
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

func injectTriggerTagsForKinesis(payload json.RawMessage) {
	e := events.KinesisEvent{}
	err := json.Unmarshal(payload, &e)
	if err != nil {
		return
	}

	domainName := constants.DomainNames["STREAM"]
	className := constants.ClassNames["KINESIS"]
	streamNamesMap := make(map[string]struct{})

	var links []string
	for _, record := range e.Records {
		streamName := ""
		i := strings.Index(record.EventSourceArn, "/")
		if i != -1 && (i+1) < len(record.EventSourceArn) {
			streamName = record.EventSourceArn[i+1:]
		}
		streamNamesMap[streamName] = void
		link := record.AwsRegion + ":" + streamName + ":" + record.EventID
		links = append(links, link)
	}

	var streamNames []string
	for k := range streamNamesMap {
		streamNames = append(streamNames, k)
	}

	AddIncomingTraceLinks(links)
	injectTriggerTagsToInvocation(domainName, className, streamNames)
}

func injectTriggerTagsForKinesisFirehose(payload json.RawMessage) {
	e := events.KinesisFirehoseEvent{}
	err := json.Unmarshal(payload, &e)
	if err != nil {
		return
	}

	domainName := constants.DomainNames["STREAM"]
	className := constants.ClassNames["FIREHOSE"]
	streamName := ""
	i := strings.Index(e.DeliveryStreamArn, "/")

	var traceLinks []string
	if i != -1 && (i+1) < len(e.DeliveryStreamArn) {
		streamName = e.DeliveryStreamArn[i+1:]

		for _, record := range e.Records {
			if record.ApproximateArrivalTimestamp != (events.MilliSecondsEpochTime{}) {
				timestamp := record.ApproximateArrivalTimestamp.Unix()
				b := md5.Sum(record.Data)
				dataMD5 := hex.EncodeToString(b[:])
				addFirehoseLink(&traceLinks, e.Region, dataMD5, timestamp, streamName)
			}
		}
	}
	AddIncomingTraceLinks(traceLinks)
	var streamNames = []string{streamName}
	injectTriggerTagsToInvocation(domainName, className, streamNames)
}

func addFirehoseLink(traceLinks *[]string, region string, dataMD5 string, timestamp int64, streamName string) {
	for j := 0; j < 3; j++ {
		*traceLinks = append(*traceLinks, region+":"+streamName+":"+strconv.FormatInt(timestamp+int64(j), 10)+":"+dataMD5)
	}
}

func injectTriggerTagsForSNS(payload json.RawMessage) {
	e := events.SNSEvent{}
	err := json.Unmarshal(payload, &e)
	if err != nil {
		return
	}

	domainName := constants.DomainNames["MESSAGING"]
	className := constants.ClassNames["SNS"]
	streamNamesMap := make(map[string]struct{})
	var traceLinks []string

	for _, record := range e.Records {
		topicSlice := strings.Split(record.SNS.TopicArn, ":")
		topicName := topicSlice[len(topicSlice)-1]
		streamNamesMap[topicName] = void
		if record.SNS.MessageID != "" {
			traceLinks = append(traceLinks, record.SNS.MessageID)
		}
	}

	var topicNames []string
	for k := range streamNamesMap {
		topicNames = append(topicNames, k)
	}

	AddIncomingTraceLinks(traceLinks)
	injectTriggerTagsToInvocation(domainName, className, topicNames)
}

func injectTriggerTagsForSQS(payload json.RawMessage) {
	e := events.SQSEvent{}
	err := json.Unmarshal(payload, &e)
	if err != nil {
		return
	}

	domainName := constants.DomainNames["MESSAGING"]
	className := constants.ClassNames["SQS"]
	queueNamesMap := make(map[string]struct{})
	var traceLinks []string

	for _, record := range e.Records {
		queueSlice := strings.Split(record.EventSourceARN, ":")
		topicName := queueSlice[len(queueSlice)-1]
		queueNamesMap[topicName] = void
		if record.MessageId != "" {
			traceLinks = append(traceLinks, record.MessageId)
		}
	}

	var queueNames []string
	for k := range queueNamesMap {
		queueNames = append(queueNames, k)
	}

	AddIncomingTraceLinks(traceLinks)
	injectTriggerTagsToInvocation(domainName, className, queueNames)
}

func injectTriggerTagsForS3(payload json.RawMessage) {
	e := events.S3Event{}
	err := json.Unmarshal(payload, &e)
	if err != nil {
		return
	}

	domainName := constants.DomainNames["STORAGE"]
	className := constants.ClassNames["S3"]
	bucketNamesMap := make(map[string]struct{})
	var traceLinks []string

	for _, record := range e.Records {
		bucketName := record.S3.Bucket.Name
		bucketNamesMap[bucketName] = void
		link := record.ResponseElements["x-amz-request-id"]
		if link != "" {
			traceLinks = append(traceLinks, link)
		}
	}

	var bucketNames []string
	for k := range bucketNamesMap {
		bucketNames = append(bucketNames, k)
	}

	AddIncomingTraceLinks(traceLinks)
	injectTriggerTagsToInvocation(domainName, className, bucketNames)
}

func injectTriggerTagsForCloudFront(payload json.RawMessage) {
	e := triggerEvent{}
	err := json.Unmarshal(payload, &e)
	if err != nil {
		return
	}

	domainName := constants.DomainNames["CDN"]
	className := constants.ClassNames["CLOUDFRONT"]
	urisMap := make(map[string]struct{})
	for _, record := range e.Records {
		uri := record.CF.Request.URI
		urisMap[uri] = void
	}

	var uris []string
	for k := range urisMap {
		uris = append(uris, k)
	}

	injectTriggerTagsToInvocation(domainName, className, uris)
}

func injectTriggerTagsForAPIGateway(payload json.RawMessage) {
	e := triggerEvent{}
	err := json.Unmarshal(payload, &e)
	if err != nil {
		return
	}

	domainName := constants.DomainNames["API"]
	className := constants.ClassNames["APIGATEWAY"]

	path := ""
	if e.Params.Header["Host"] != "" && e.Context["stage"] != "" {
		path = e.Params.Header["Host"] + "/" + e.Context["stage"] + e.Context["resource-path"]
	}

	var operationNames = []string{path}

	injectTriggerTagsToInvocation(domainName, className, operationNames)
}

func injectTriggerTagsForAPIGatewayProxy(payload json.RawMessage) {
	e := events.APIGatewayProxyRequest{}
	err := json.Unmarshal(payload, &e)
	if err != nil {
		return
	}

	domainName := constants.DomainNames["API"]
	className := constants.ClassNames["APIGATEWAY"]

	var operationNames = []string{e.Resource}

	if e.Headers != nil {
		spanID := e.Headers["x-thundra-span-id"]
		if spanID != "" {
			AddIncomingTraceLinks([]string{spanID})
		}
	}

	injectTriggerTagsToInvocation(domainName, className, operationNames)
}

func injectTriggerTagsForCloudWatchLogs(payload json.RawMessage) {
	e := events.CloudwatchLogsEvent{}
	err := json.Unmarshal(payload, &e)
	if err != nil {
		return
	}

	domainName := constants.DomainNames["LOG"]
	className := constants.ClassNames["CLOUDWATCHLOG"]

	data, err := e.AWSLogs.Parse()
	var operationNames []string
	if err == nil {
		operationNames = append(operationNames, data.LogGroup)
	}
	injectTriggerTagsToInvocation(domainName, className, operationNames)

}

func injectTriggerTagsForSchedule(payload json.RawMessage) {
	e := events.CloudWatchEvent{}
	err := json.Unmarshal(payload, &e)
	if err != nil {
		return
	}

	domainName := constants.DomainNames["SCHEDULE"]
	className := constants.ClassNames["SCHEDULE"]

	scheduleNamesMap := make(map[string]struct{})
	for _, resource := range e.Resources {
		scheduleSlice := strings.Split(resource, "/")
		scheduleName := scheduleSlice[len(scheduleSlice)-1]
		scheduleNamesMap[scheduleName] = void
	}
	var operationNames []string
	for k := range scheduleNamesMap {
		operationNames = append(operationNames, k)
	}

	injectTriggerTagsToInvocation(domainName, className, operationNames)
}

func injectTriggerTagsForLambda(ctx context.Context) {
	domainName := constants.DomainNames["API"]
	className := constants.ClassNames["LAMBDA"]

	clientContext, ok := application.GetClientContext(ctx)
	if ok {
		operationName := clientContext.Custom[constants.AwsLambdaTriggerOperationName]
		operationNames := []string{operationName}
		injectTriggerTagsToInvocation(domainName, className, operationNames)
	}
	awsRequestID := application.GetAwsRequestID(ctx)
	if awsRequestID != "" {
		AddIncomingTraceLinks([]string{awsRequestID})
	}
}

func setInvocationTriggerTags(ctx context.Context, payload json.RawMessage) {
	ok := injectTriggerTagsFromInputType(ctx, payload)
	if !ok {
		injectTriggerTagsFromPayload(ctx, payload)
	}
}

func injectTriggerTagsFromInputType(ctx context.Context, payload json.RawMessage) bool {
	eventType, ok := utils.GetEventTypeFromContext(ctx).(reflect.Type)
	if !ok {
		return false
	}
	switch eventType {
	case reflect.TypeOf(events.DynamoDBEvent{}):
		injectTriggerTagsForDynamoDB(payload)
	case reflect.TypeOf(events.KinesisEvent{}):
		injectTriggerTagsForKinesis(payload)
	case reflect.TypeOf(events.KinesisFirehoseEvent{}):
		injectTriggerTagsForKinesisFirehose(payload)
	case reflect.TypeOf(events.SNSEvent{}):
		injectTriggerTagsForSNS(payload)
	case reflect.TypeOf(events.SQSEvent{}):
		injectTriggerTagsForSQS(payload)
	case reflect.TypeOf(events.S3Event{}):
		injectTriggerTagsForS3(payload)
	case reflect.TypeOf(events.CloudwatchLogsEvent{}):
		injectTriggerTagsForCloudWatchLogs(payload)
	case reflect.TypeOf(events.CloudWatchEvent{}):
		injectTriggerTagsForSchedule(payload)
	case reflect.TypeOf(events.APIGatewayProxyRequest{}):
		injectTriggerTagsForAPIGatewayProxy(payload)
	default:
		ok = false
	}
	return ok
}

func injectTriggerTagsFromPayload(ctx context.Context, payload json.RawMessage) {
	clientContext, ok := application.GetClientContext(ctx)
	if ok {
		if clientContext.Custom[constants.AwsLambdaTriggerOperationName] != "" {
			injectTriggerTagsForLambda(ctx)
			return
		}
	}

	var rawEvent triggerEvent
	err := json.Unmarshal(payload, &rawEvent)
	if err != nil {
		return
	}
	if len(rawEvent.Records) > 0 {
		switch rawEvent.Records[0].EventSource {
		case "aws:dynamodb":
			injectTriggerTagsForDynamoDB(payload)
		case "aws:kinesis":
			injectTriggerTagsForKinesis(payload)
		case "aws:sns":
			injectTriggerTagsForSNS(payload)
		case "aws:sqs":
			injectTriggerTagsForSQS(payload)
		case "aws:s3":
			injectTriggerTagsForS3(payload)
		}
		if rawEvent.DeliveryStreamArn != "" {
			injectTriggerTagsForKinesisFirehose(payload)
		} else if rawEvent.Records[0].CF != (cloudFrontField{}) {
			injectTriggerTagsForCloudFront(payload)
		}
	} else if rawEvent.RequestContext != (requestContextField{}) && rawEvent.Headers != nil && len(rawEvent.HTTPMethod) > 0 {
		injectTriggerTagsForAPIGatewayProxy(payload)
	} else if _, ok := rawEvent.Context["http-method"]; ok {
		injectTriggerTagsForAPIGateway(payload)
	} else if rawEvent.AwsLogs.Data != "" {
		injectTriggerTagsForCloudWatchLogs(payload)
	} else if rawEvent.DetailType == "Scheduled Event" {
		injectTriggerTagsForSchedule(payload)
	}
}
