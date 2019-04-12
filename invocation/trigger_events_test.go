package invocation

import (
	"context"
	"io/ioutil"
	"reflect"
	"testing"

	"github.com/thundra-io/thundra-lambda-agent-go/constants"
	"github.com/thundra-io/thundra-lambda-agent-go/utils"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/stretchr/testify/assert"
)

type expected struct {
	val string
	err error
}

func readJSONFromFile(t *testing.T, inputFile string) []byte {
	inputJSON, err := ioutil.ReadFile(inputFile)
	if err != nil {
		t.Errorf("Could not open test file. Details: %v", err)
	}

	return inputJSON
}

func createMockLambdaTriggerContext() context.Context {
	cc := lambdacontext.ClientContext{}
	lc := &lambdacontext.LambdaContext{}
	cc.Custom = map[string]string{
		constants.AwsLambdaTriggerOperationName: "test-function",
	}
	c := context.Background()
	lc.ClientContext = cc
	lc.AwsRequestID = "aws_request_id"
	c = lambdacontext.NewContext(c, lc)
	return c
}

func TestInvocationTags_SNSTrigger(t *testing.T) {
	ClearTags()
	clearTraceLinks()

	eventMock := readJSONFromFile(t, "./testdata/sns-event.json")
	setInvocationTriggerTags(context.TODO(), eventMock)

	operationNames := []string{"EXAMPLE"}

	assert.Equal(t, invocationTags[constants.SpanTags["TRIGGER_DOMAIN_NAME"]], constants.DomainNames["MESSAGING"])
	assert.Equal(t, invocationTags[constants.SpanTags["TRIGGER_CLASS_NAME"]], constants.ClassNames["SNS"])
	assert.Equal(t, invocationTags[constants.SpanTags["TRIGGER_OPERATION_NAMES"]], operationNames)

	traceLinks := getIncomingTraceLinks()
	assert.ElementsMatch(t, traceLinks, []string{"95df01b4-ee98-5cb9-9903-4c221d41eb5e"})
}

func TestInvocationTags_SNSTriggerFromInputType(t *testing.T) {
	ClearTags()
	clearTraceLinks()

	eventMock := readJSONFromFile(t, "./testdata/sns-event.json")
	ctx := context.Background()
	ctx = utils.SetEventTypeToContext(ctx, reflect.TypeOf(events.SNSEvent{}))
	ok := injectTriggerTagsFromInputType(ctx, eventMock)

	operationNames := []string{"EXAMPLE"}

	assert.True(t, ok)

	assert.Equal(t, invocationTags[constants.SpanTags["TRIGGER_DOMAIN_NAME"]], constants.DomainNames["MESSAGING"])
	assert.Equal(t, invocationTags[constants.SpanTags["TRIGGER_CLASS_NAME"]], constants.ClassNames["SNS"])
	assert.Equal(t, invocationTags[constants.SpanTags["TRIGGER_OPERATION_NAMES"]], operationNames)

	traceLinks := getIncomingTraceLinks()
	assert.ElementsMatch(t, traceLinks, []string{"95df01b4-ee98-5cb9-9903-4c221d41eb5e"})
}

func TestInvocationTags_SQSTrigger(t *testing.T) {
	ClearTags()
	clearTraceLinks()

	eventMock := readJSONFromFile(t, "./testdata/sqs-event.json")
	setInvocationTriggerTags(context.TODO(), eventMock)

	operationNames := []string{"SQSQueue"}

	assert.Equal(t, invocationTags[constants.SpanTags["TRIGGER_DOMAIN_NAME"]], constants.DomainNames["MESSAGING"])
	assert.Equal(t, invocationTags[constants.SpanTags["TRIGGER_CLASS_NAME"]], constants.ClassNames["SQS"])
	assert.Equal(t, invocationTags[constants.SpanTags["TRIGGER_OPERATION_NAMES"]], operationNames)

	traceLinks := getIncomingTraceLinks()
	assert.ElementsMatch(t, traceLinks, []string{"MessageID_1"})
}

func TestInvocationTags_SQSTriggerFromInputType(t *testing.T) {
	ClearTags()
	clearTraceLinks()

	ctx := context.Background()
	ctx = utils.SetEventTypeToContext(ctx, reflect.TypeOf(events.SQSEvent{}))

	eventMock := readJSONFromFile(t, "./testdata/sqs-event.json")

	ok := injectTriggerTagsFromInputType(ctx, eventMock)

	assert.True(t, ok)

	operationNames := []string{"SQSQueue"}

	assert.Equal(t, invocationTags[constants.SpanTags["TRIGGER_DOMAIN_NAME"]], constants.DomainNames["MESSAGING"])
	assert.Equal(t, invocationTags[constants.SpanTags["TRIGGER_CLASS_NAME"]], constants.ClassNames["SQS"])
	assert.Equal(t, invocationTags[constants.SpanTags["TRIGGER_OPERATION_NAMES"]], operationNames)

	traceLinks := getIncomingTraceLinks()
	assert.ElementsMatch(t, traceLinks, []string{"MessageID_1"})
}

func TestInvocationTags_CFTrigger(t *testing.T) {
	ClearTags()

	eventMock := readJSONFromFile(t, "./testdata/cf-event.json")
	setInvocationTriggerTags(context.TODO(), eventMock)

	operationNames := []string{"/test"}

	assert.Equal(t, invocationTags[constants.SpanTags["TRIGGER_DOMAIN_NAME"]], constants.DomainNames["CDN"])
	assert.Equal(t, invocationTags[constants.SpanTags["TRIGGER_CLASS_NAME"]], constants.ClassNames["CLOUDFRONT"])
	assert.Equal(t, invocationTags[constants.SpanTags["TRIGGER_OPERATION_NAMES"]], operationNames)
}

func TestInvocationTags_APIGatewayTrigger(t *testing.T) {
	ClearTags()

	eventMock := readJSONFromFile(t, "./testdata/api-gw.json")
	setInvocationTriggerTags(context.TODO(), eventMock)

	operationNames := []string{"random.execute-api.us-west-2.amazonaws.com/dev/hello"}

	assert.Equal(t, invocationTags[constants.SpanTags["TRIGGER_DOMAIN_NAME"]], constants.DomainNames["API"])
	assert.Equal(t, invocationTags[constants.SpanTags["TRIGGER_CLASS_NAME"]], constants.ClassNames["APIGATEWAY"])
	assert.Equal(t, invocationTags[constants.SpanTags["TRIGGER_OPERATION_NAMES"]], operationNames)
}

func TestInvocationTags_APIGatewayProxyTrigger(t *testing.T) {
	ClearTags()
	clearTraceLinks()

	eventMock := readJSONFromFile(t, "./testdata/apigw-request.json")
	setInvocationTriggerTags(context.TODO(), eventMock)

	operationNames := []string{"gy415nuibc.execute-api.us-east-1.amazonaws.com/testStage/hello/world"}

	assert.Equal(t, invocationTags[constants.SpanTags["TRIGGER_DOMAIN_NAME"]], constants.DomainNames["API"])
	assert.Equal(t, invocationTags[constants.SpanTags["TRIGGER_CLASS_NAME"]], constants.ClassNames["APIGATEWAY"])
	assert.Equal(t, invocationTags[constants.SpanTags["TRIGGER_OPERATION_NAMES"]], operationNames)

	traceLinks := getIncomingTraceLinks()
	assert.Equal(t, traceLinks, []string{"test_span_id"})
}

func TestInvocationTags_APIGatewayProxyTriggerFromInputType(t *testing.T) {
	ClearTags()
	clearTraceLinks()

	ctx := context.Background()
	ctx = utils.SetEventTypeToContext(ctx, reflect.TypeOf(events.APIGatewayProxyRequest{}))
	eventMock := readJSONFromFile(t, "./testdata/apigw-request.json")

	ok := injectTriggerTagsFromInputType(ctx, eventMock)
	assert.True(t, ok)

	operationNames := []string{"gy415nuibc.execute-api.us-east-1.amazonaws.com/testStage/hello/world"}

	assert.Equal(t, invocationTags[constants.SpanTags["TRIGGER_DOMAIN_NAME"]], constants.DomainNames["API"])
	assert.Equal(t, invocationTags[constants.SpanTags["TRIGGER_CLASS_NAME"]], constants.ClassNames["APIGATEWAY"])
	assert.Equal(t, invocationTags[constants.SpanTags["TRIGGER_OPERATION_NAMES"]], operationNames)
}

func TestInvocationTags_S3Trigger(t *testing.T) {
	ClearTags()
	clearTraceLinks()

	eventMock := readJSONFromFile(t, "./testdata/s3-event.json")
	setInvocationTriggerTags(context.TODO(), eventMock)

	operationNames := []string{"sourcebucket"}

	assert.Equal(t, invocationTags[constants.SpanTags["TRIGGER_DOMAIN_NAME"]], constants.DomainNames["STORAGE"])
	assert.Equal(t, invocationTags[constants.SpanTags["TRIGGER_CLASS_NAME"]], constants.ClassNames["S3"])
	assert.Equal(t, invocationTags[constants.SpanTags["TRIGGER_OPERATION_NAMES"]], operationNames)
}

func TestInvocationTags_S3TriggerFromInputType(t *testing.T) {
	ClearTags()

	ctx := context.Background()
	ctx = utils.SetEventTypeToContext(ctx, reflect.TypeOf(events.S3Event{}))
	eventMock := readJSONFromFile(t, "./testdata/s3-event.json")

	ok := injectTriggerTagsFromInputType(ctx, eventMock)
	assert.True(t, ok)

	operationNames := []string{"sourcebucket"}

	assert.Equal(t, invocationTags[constants.SpanTags["TRIGGER_DOMAIN_NAME"]], constants.DomainNames["STORAGE"])
	assert.Equal(t, invocationTags[constants.SpanTags["TRIGGER_CLASS_NAME"]], constants.ClassNames["S3"])
	assert.Equal(t, invocationTags[constants.SpanTags["TRIGGER_OPERATION_NAMES"]], operationNames)
}

func TestInvocationTags_ScheduleTrigger(t *testing.T) {
	ClearTags()

	eventMock := readJSONFromFile(t, "./testdata/schedule-event.json")
	setInvocationTriggerTags(context.TODO(), eventMock)

	operationNames := []string{"ExampleRule"}

	assert.Equal(t, invocationTags[constants.SpanTags["TRIGGER_DOMAIN_NAME"]], constants.DomainNames["SCHEDULE"])
	assert.Equal(t, invocationTags[constants.SpanTags["TRIGGER_CLASS_NAME"]], constants.ClassNames["SCHEDULE"])
	assert.Equal(t, invocationTags[constants.SpanTags["TRIGGER_OPERATION_NAMES"]], operationNames)
}

func TestInvocationTags_ScheduleTriggerFromInputType(t *testing.T) {
	ClearTags()

	ctx := context.Background()
	ctx = utils.SetEventTypeToContext(ctx, reflect.TypeOf(events.CloudWatchEvent{}))
	eventMock := readJSONFromFile(t, "./testdata/schedule-event.json")

	ok := injectTriggerTagsFromInputType(ctx, eventMock)
	assert.True(t, ok)

	operationNames := []string{"ExampleRule"}

	assert.Equal(t, invocationTags[constants.SpanTags["TRIGGER_DOMAIN_NAME"]], constants.DomainNames["SCHEDULE"])
	assert.Equal(t, invocationTags[constants.SpanTags["TRIGGER_CLASS_NAME"]], constants.ClassNames["SCHEDULE"])
	assert.Equal(t, invocationTags[constants.SpanTags["TRIGGER_OPERATION_NAMES"]], operationNames)
}

func TestInvocationTags_KinesisTrigger(t *testing.T) {
	ClearTags()
	clearTraceLinks()

	eventMock := readJSONFromFile(t, "./testdata/kinesis-event.json")
	setInvocationTriggerTags(context.TODO(), eventMock)

	operationNames := []string{"simple-stream"}

	assert.Equal(t, invocationTags[constants.SpanTags["TRIGGER_DOMAIN_NAME"]], constants.DomainNames["STREAM"])
	assert.Equal(t, invocationTags[constants.SpanTags["TRIGGER_CLASS_NAME"]], constants.ClassNames["KINESIS"])
	assert.Equal(t, invocationTags[constants.SpanTags["TRIGGER_OPERATION_NAMES"]], operationNames)

	traceLinks := getIncomingTraceLinks()
	expectedLinks := []string{
		"us-east-1:simple-stream:shardId-000000000000:49568167373333333333333333333333333333333333333333333333",
		"us-east-1:simple-stream:shardId-000000000000:49568167373333333334444444444444444444444444444444444444",
	}
	assert.ElementsMatch(t, traceLinks, expectedLinks)
}

func TestInvocationTags_KinesisTriggerFromInputType(t *testing.T) {
	ClearTags()
	clearTraceLinks()

	ctx := context.Background()
	ctx = utils.SetEventTypeToContext(ctx, reflect.TypeOf(events.KinesisEvent{}))
	eventMock := readJSONFromFile(t, "./testdata/kinesis-event.json")

	ok := injectTriggerTagsFromInputType(ctx, eventMock)
	assert.True(t, ok)

	operationNames := []string{"simple-stream"}

	assert.Equal(t, invocationTags[constants.SpanTags["TRIGGER_DOMAIN_NAME"]], constants.DomainNames["STREAM"])
	assert.Equal(t, invocationTags[constants.SpanTags["TRIGGER_CLASS_NAME"]], constants.ClassNames["KINESIS"])
	assert.Equal(t, invocationTags[constants.SpanTags["TRIGGER_OPERATION_NAMES"]], operationNames)

	traceLinks := getIncomingTraceLinks()
	expectedLinks := []string{
		"us-east-1:simple-stream:shardId-000000000000:49568167373333333333333333333333333333333333333333333333",
		"us-east-1:simple-stream:shardId-000000000000:49568167373333333334444444444444444444444444444444444444",
	}
	assert.ElementsMatch(t, traceLinks, expectedLinks)
}

func TestInvocationTags_KinesisFirehoseTrigger(t *testing.T) {
	ClearTags()
	clearTraceLinks()

	eventMock := readJSONFromFile(t, "./testdata/kinesis-firehose-event.json")
	setInvocationTriggerTags(context.TODO(), eventMock)

	operationNames := []string{"exampleStream"}

	assert.Equal(t, invocationTags[constants.SpanTags["TRIGGER_DOMAIN_NAME"]], constants.DomainNames["STREAM"])
	assert.Equal(t, invocationTags[constants.SpanTags["TRIGGER_CLASS_NAME"]], constants.ClassNames["FIREHOSE"])
	assert.Equal(t, invocationTags[constants.SpanTags["TRIGGER_OPERATION_NAMES"]], operationNames)
}

func TestInvocationTags_KinesisFirehoseTriggerFromInputType(t *testing.T) {
	ClearTags()
	clearTraceLinks()

	ctx := context.Background()
	ctx = utils.SetEventTypeToContext(ctx, reflect.TypeOf(events.KinesisFirehoseEvent{}))
	eventMock := readJSONFromFile(t, "./testdata/kinesis-firehose-event.json")

	ok := injectTriggerTagsFromInputType(ctx, eventMock)
	assert.True(t, ok)

	operationNames := []string{"exampleStream"}

	assert.Equal(t, invocationTags[constants.SpanTags["TRIGGER_DOMAIN_NAME"]], constants.DomainNames["STREAM"])
	assert.Equal(t, invocationTags[constants.SpanTags["TRIGGER_CLASS_NAME"]], constants.ClassNames["FIREHOSE"])
	assert.Equal(t, invocationTags[constants.SpanTags["TRIGGER_OPERATION_NAMES"]], operationNames)
}

func TestInvocationTags_DynamoDBTrigger(t *testing.T) {
	ClearTags()
	clearTraceLinks()

	eventMock := readJSONFromFile(t, "./testdata/dynamodb-event.json")
	setInvocationTriggerTags(context.TODO(), eventMock)

	operationNames := []string{"Example-Table", "Example-Table2"}

	assert.Equal(t, invocationTags[constants.SpanTags["TRIGGER_DOMAIN_NAME"]], constants.DomainNames["DB"])
	assert.Equal(t, invocationTags[constants.SpanTags["TRIGGER_CLASS_NAME"]], constants.ClassNames["DYNAMODB"])
	assert.ElementsMatch(t, invocationTags[constants.SpanTags["TRIGGER_OPERATION_NAMES"]], operationNames)

}

func TestInvocationTags_DynamoDBTriggerFromInputType(t *testing.T) {
	ClearTags()
	clearTraceLinks()

	ctx := context.Background()
	ctx = utils.SetEventTypeToContext(ctx, reflect.TypeOf(events.DynamoDBEvent{}))
	eventMock := readJSONFromFile(t, "./testdata/dynamodb-event.json")

	ok := injectTriggerTagsFromInputType(ctx, eventMock)
	assert.True(t, ok)

	operationNames := []string{"Example-Table", "Example-Table2"}

	assert.Equal(t, invocationTags[constants.SpanTags["TRIGGER_DOMAIN_NAME"]], constants.DomainNames["DB"])
	assert.Equal(t, invocationTags[constants.SpanTags["TRIGGER_CLASS_NAME"]], constants.ClassNames["DYNAMODB"])
	assert.ElementsMatch(t, invocationTags[constants.SpanTags["TRIGGER_OPERATION_NAMES"]], operationNames)
}

func TestInvocationTags_CloudWatchLogsTrigger(t *testing.T) {
	ClearTags()

	eventMock := readJSONFromFile(t, "./testdata/cloudwatch-logs-event.json")
	setInvocationTriggerTags(context.TODO(), eventMock)

	operationNames := []string{"testLogGroup"}

	assert.Equal(t, invocationTags[constants.SpanTags["TRIGGER_DOMAIN_NAME"]], constants.DomainNames["LOG"])
	assert.Equal(t, invocationTags[constants.SpanTags["TRIGGER_CLASS_NAME"]], constants.ClassNames["CLOUDWATCHLOG"])
	assert.Equal(t, invocationTags[constants.SpanTags["TRIGGER_OPERATION_NAMES"]], operationNames)
}

func TestInvocationTags_CloudWatchLogsTriggerFromInputType(t *testing.T) {
	ClearTags()

	ctx := context.Background()
	ctx = utils.SetEventTypeToContext(ctx, reflect.TypeOf(events.CloudwatchLogsEvent{}))

	eventMock := readJSONFromFile(t, "./testdata/cloudwatch-logs-event.json")

	ok := injectTriggerTagsFromInputType(ctx, eventMock)
	assert.True(t, ok)

	operationNames := []string{"testLogGroup"}

	assert.Equal(t, invocationTags[constants.SpanTags["TRIGGER_DOMAIN_NAME"]], constants.DomainNames["LOG"])
	assert.Equal(t, invocationTags[constants.SpanTags["TRIGGER_CLASS_NAME"]], constants.ClassNames["CLOUDWATCHLOG"])
	assert.Equal(t, invocationTags[constants.SpanTags["TRIGGER_OPERATION_NAMES"]], operationNames)
}

func TestInvocationTags_LambdaTrigger(t *testing.T) {
	ClearTags()
	clearTraceLinks()

	c := createMockLambdaTriggerContext()
	setInvocationTriggerTags(c, nil)
	operationNames := []string{"test-function"}

	assert.Equal(t, invocationTags[constants.SpanTags["TRIGGER_DOMAIN_NAME"]], constants.DomainNames["API"])
	assert.Equal(t, invocationTags[constants.SpanTags["TRIGGER_CLASS_NAME"]], constants.ClassNames["LAMBDA"])
	assert.Equal(t, invocationTags[constants.SpanTags["TRIGGER_OPERATION_NAMES"]], operationNames)

	traceLinks := getIncomingTraceLinks()
	assert.ElementsMatch(t, traceLinks, []string{"aws_request_id"})
}

func TestInvocationTags_NilEvent(t *testing.T) {
	ClearTags()

	setInvocationTriggerTags(context.TODO(), nil)

	assert.Equal(t, invocationTags[constants.SpanTags["TRIGGER_DOMAIN_NAME"]], nil)
	assert.Equal(t, invocationTags[constants.SpanTags["TRIGGER_CLASS_NAME"]], nil)
	assert.Equal(t, invocationTags[constants.SpanTags["TRIGGER_OPERATION_NAMES"]], nil)
}

func TestInvocationTags_UnknownInputEventFromInput(t *testing.T) {
	ClearTags()

	ctx := context.Background()
	ctx = utils.SetEventTypeToContext(ctx, reflect.TypeOf(events.CognitoEvent{}))
	eventMock := readJSONFromFile(t, "./testdata/alb-lambda-target-request-headers-only.json")

	ok := injectTriggerTagsFromInputType(ctx, eventMock)
	assert.False(t, ok)

	assert.Equal(t, invocationTags[constants.SpanTags["TRIGGER_DOMAIN_NAME"]], nil)
	assert.Equal(t, invocationTags[constants.SpanTags["TRIGGER_CLASS_NAME"]], nil)
	assert.Equal(t, invocationTags[constants.SpanTags["TRIGGER_OPERATION_NAMES"]], nil)
}

func TestInvocationTags_MalformedPayloadJSON(t *testing.T) {
	ClearTags()

	eventMock := readJSONFromFile(t, "./testdata/dynamodb-event-malformed.json")
	setInvocationTriggerTags(context.Background(), eventMock)

	assert.Equal(t, invocationTags[constants.SpanTags["TRIGGER_DOMAIN_NAME"]], nil)
	assert.Equal(t, invocationTags[constants.SpanTags["TRIGGER_CLASS_NAME"]], nil)
	assert.Equal(t, invocationTags[constants.SpanTags["TRIGGER_OPERATION_NAMES"]], nil)
}

func TestInvocationTags_MalformedEventWithInputType(t *testing.T) {
	ClearTags()

	ctx := context.Background()
	ctx = utils.SetEventTypeToContext(ctx, reflect.TypeOf(events.DynamoDBEvent{}))
	eventMock := readJSONFromFile(t, "./testdata/dynamodb-event-malformed.json")
	setInvocationTriggerTags(ctx, eventMock)

	assert.Equal(t, invocationTags[constants.SpanTags["TRIGGER_DOMAIN_NAME"]], nil)
	assert.Equal(t, invocationTags[constants.SpanTags["TRIGGER_CLASS_NAME"]], nil)
	assert.Equal(t, invocationTags[constants.SpanTags["TRIGGER_OPERATION_NAMES"]], nil)
}

func TestInvocationTags_IncorrectInputEventTypeForPayload(t *testing.T) {
	ClearTags()
	ctx := context.Background()
	ctx = utils.SetEventTypeToContext(ctx, reflect.TypeOf(events.SQSEvent{}))

	eventMock := readJSONFromFile(t, "./testdata/sns-event.json")

	ok := injectTriggerTagsFromInputType(ctx, eventMock)
	assert.True(t, ok)

	assert.Equal(t, invocationTags[constants.SpanTags["TRIGGER_DOMAIN_NAME"]], constants.DomainNames["MESSAGING"])
	assert.Equal(t, invocationTags[constants.SpanTags["TRIGGER_CLASS_NAME"]], constants.ClassNames["SQS"])
	assert.Equal(t, invocationTags[constants.SpanTags["TRIGGER_OPERATION_NAMES"]], []string{""})
}

func TestInvocationTags_DynamoDBPayloadWithInvalidField(t *testing.T) {
	ClearTags()

	eventMock := readJSONFromFile(t, "./testdata/dynamodb-event-wrong-arn.json")
	setInvocationTriggerTags(context.Background(), eventMock)

	assert.Equal(t, invocationTags[constants.SpanTags["TRIGGER_DOMAIN_NAME"]], constants.DomainNames["DB"])
	assert.Equal(t, invocationTags[constants.SpanTags["TRIGGER_CLASS_NAME"]], constants.ClassNames["DYNAMODB"])
	assert.Equal(t, invocationTags[constants.SpanTags["TRIGGER_OPERATION_NAMES"]], []string{""})
}

func TestInvocationTags_APIGatewayPayloadWithMissingFields(t *testing.T) {
	ClearTags()

	eventMock := readJSONFromFile(t, "./testdata/api-gw-missing-fields.json")
	setInvocationTriggerTags(context.Background(), eventMock)

	assert.Equal(t, invocationTags[constants.SpanTags["TRIGGER_DOMAIN_NAME"]], constants.DomainNames["API"])
	assert.Equal(t, invocationTags[constants.SpanTags["TRIGGER_CLASS_NAME"]], constants.ClassNames["APIGATEWAY"])
	assert.Equal(t, invocationTags[constants.SpanTags["TRIGGER_OPERATION_NAMES"]], []string{""})
}

func TestInvocationTags_APIGatewayProxyPayloadWithMissingFields(t *testing.T) {
	ClearTags()

	eventMock := readJSONFromFile(t, "./testdata/apigw-request-missing-fields.json")
	setInvocationTriggerTags(context.Background(), eventMock)

	assert.Equal(t, invocationTags[constants.SpanTags["TRIGGER_DOMAIN_NAME"]], constants.DomainNames["API"])
	assert.Equal(t, invocationTags[constants.SpanTags["TRIGGER_CLASS_NAME"]], constants.ClassNames["APIGATEWAY"])
	assert.Equal(t, invocationTags[constants.SpanTags["TRIGGER_OPERATION_NAMES"]], []string{""})
}
