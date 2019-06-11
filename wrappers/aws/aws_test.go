package thundraaws

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/aws/aws-sdk-go/aws/request"

	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/firehose"
	"github.com/aws/aws-sdk-go/service/kinesis"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/stretchr/testify/assert"
	"github.com/thundra-io/thundra-lambda-agent-go/application"
	"github.com/thundra-io/thundra-lambda-agent-go/config"
	"github.com/thundra-io/thundra-lambda-agent-go/constants"
	"github.com/thundra-io/thundra-lambda-agent-go/trace"

	"github.com/aws/aws-sdk-go/aws/session"
)

type RoundTripFunc func(req *http.Request) (*http.Response, error)

func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func NewTestClient(f RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: f,
	}
}

var testClient = NewTestClient(func(req *http.Request) (*http.Response, error) {
	return nil, http.ErrServerClosed
})

var sess = session.New(&aws.Config{
	HTTPClient: testClient,
	Region:     aws.String("us-west-2"),
	MaxRetries: aws.Int(0),
})

func getSessionWithDateResponseHeader() *session.Session {
	var sess = session.New(&aws.Config{
		HTTPClient: testClient,
		Region:     aws.String("us-west-2"),
		MaxRetries: aws.Int(0),
	})
	wrappedSess := Wrap(sess)
	wrappedSess.Handlers.Complete.PushFront(func(r *request.Request) {
		r.HTTPResponse = &http.Response{}
		r.HTTPResponse.Header = http.Header{}
		r.HTTPResponse.Header.Set("date", "Thu, 10 Apr 2019 16:00:00 GMT")
	})
	return wrappedSess
}

func getSessionWithSnsResponse() *session.Session {
	var sess = session.New(&aws.Config{
		HTTPClient: testClient,
		Region:     aws.String("us-west-2"),
		MaxRetries: aws.Int(0),
	})
	wrappedSess := Wrap(sess)

	snsData := sns.PublishOutput{MessageId: aws.String("95df01b4-ee98-5cb9-9903-4c221d41eb5e")}
	wrappedSess.Handlers.Complete.PushFront(func(r *request.Request) {
		r.Data = &snsData
	})

	return wrappedSess
}

func getSessionWithKinesisResponse() *session.Session {
	var sess = session.New(&aws.Config{
		HTTPClient: testClient,
		Region:     aws.String("us-west-2"),
		MaxRetries: aws.Int(0),
	})
	kinesisData := &kinesis.PutRecordOutput{
		ShardId:        aws.String("shardId-000000000000"),
		SequenceNumber: aws.String("49568167373333333333333333333333333333333333333333333333"),
	}
	sess = Wrap(sess)
	sess.Handlers.Complete.PushFront(func(r *request.Request) {
		r.Data = kinesisData
	})
	return sess
}

func getSessionWithKinesisPutRecordsResponse() *session.Session {
	var sess = session.New(&aws.Config{
		HTTPClient: testClient,
		Region:     aws.String("us-west-2"),
		MaxRetries: aws.Int(0),
	})
	kinesisData := &kinesis.PutRecordsOutput{
		Records: []*kinesis.PutRecordsResultEntry{
			&kinesis.PutRecordsResultEntry{
				ShardId:        aws.String("shardId-000000000000"),
				SequenceNumber: aws.String("49568167373333333333333333333333333333333333333333333333"),
			},
			&kinesis.PutRecordsResultEntry{
				ShardId:        aws.String("shardId-000000000000"),
				SequenceNumber: aws.String("49568167374444444444444444444444444444444444444444444444"),
			},
		},
	}
	sess = Wrap(sess)
	sess.Handlers.Complete.PushFront(func(r *request.Request) {
		r.Data = kinesisData
	})

	return sess
}

func getSessionWithS3Response() *session.Session {
	var sess = session.New(&aws.Config{
		HTTPClient: testClient,
		Region:     aws.String("us-west-2"),
		MaxRetries: aws.Int(0),
	})
	sess = Wrap(sess)
	sess.Handlers.Complete.PushFront(func(r *request.Request) {
		r.HTTPResponse = &http.Response{}
		r.HTTPResponse.Header = http.Header{}
		r.HTTPResponse.Header.Set("x-amz-request-id", "C3D13FE58DE4C810")
	})

	return sess
}

func getSessionWithLambdaResponse() *session.Session {
	var sess = session.New(&aws.Config{
		HTTPClient: testClient,
		Region:     aws.String("us-west-2"),
		MaxRetries: aws.Int(0),
	})
	sess = Wrap(sess)
	sess.Handlers.Complete.PushFront(func(r *request.Request) {
		r.HTTPResponse = &http.Response{}
		r.HTTPResponse.Header = http.Header{}
		r.HTTPResponse.Header.Set("X-Amzn-Requestid", "C3D13FE58DE4C810")
	})

	return sess
}

func getSessionWithSqsResponse() *session.Session {
	var sess = session.New(&aws.Config{
		HTTPClient: testClient,
		Region:     aws.String("us-west-2"),
		MaxRetries: aws.Int(0),
	})
	sess = Wrap(sess)
	sess.Handlers.Complete.PushFront(func(r *request.Request) {
		r.Data = &sqs.SendMessageOutput{MessageId: aws.String("95df01b4-ee98-5cb9-9903-4c221d41eb5e")}
	})

	return sess
}

func getSessionWithSqsBatchResponse() *session.Session {
	var sess = session.New(&aws.Config{
		HTTPClient: testClient,
		Region:     aws.String("us-west-2"),
		MaxRetries: aws.Int(0),
	})
	mockSqsBatchResult := &sqs.SendMessageBatchResultEntry{MessageId: aws.String("84df12b4-ee98-2cb8-1903-1c234d56eb7e")}
	mockSqsBatchResult2 := &sqs.SendMessageBatchResultEntry{MessageId: aws.String("95df01b4-ee98-5cb9-9903-4c221d41eb5e")}
	sess = Wrap(sess)
	sess.Handlers.Complete.PushFront(func(r *request.Request) {
		r.Data = &sqs.SendMessageBatchOutput{Successful: []*sqs.SendMessageBatchResultEntry{mockSqsBatchResult, mockSqsBatchResult2}}
	})

	return sess
}

func getMockBase64EncodedClientContext() string {
	cc := &lambdacontext.ClientContext{}
	cc.Client.InstallationID = "testId"
	cc.Custom = map[string]string{
		"testKey": "testValue",
	}
	ccByte, _ := json.Marshal(cc)
	return base64.StdEncoding.EncodeToString(ccByte)
}

func getBase64EncodedClientContext() string {
	cc := &lambdacontext.ClientContext{}
	cc.Custom = map[string]string{
		constants.AwsLambdaTriggerOperationName: "test",
		constants.AwsLambdaTriggerDomainName:    "API",
		constants.AwsLambdaTriggerClassName:     "AWS-Lambda",
	}
	ccByte, _ := json.Marshal(cc)
	return base64.StdEncoding.EncodeToString(ccByte)
}

func getBase64EncodedClientContextWithMockParam() string {
	cc := &lambdacontext.ClientContext{}
	cc.Client.InstallationID = "testId"
	cc.Custom = map[string]string{
		constants.AwsLambdaTriggerOperationName: "test",
		constants.AwsLambdaTriggerDomainName:    "API",
		constants.AwsLambdaTriggerClassName:     "AWS-Lambda",
		"testKey":                               "testValue",
	}
	ccByte, _ := json.Marshal(cc)
	return base64.StdEncoding.EncodeToString(ccByte)
}

func TestDynamoDBPutItem(t *testing.T) {
	// Initilize trace plugin to set GlobalTracer of opentracing
	tp := trace.New()

	// Create a session and wrap it
	sess := getSessionWithDateResponseHeader()
	dynamoc := dynamodb.New(sess)
	// Actual call
	input := &dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{
			"AlbumTitle": {
				S: aws.String("Somewhat Famous"),
			},
			"Artist": {
				S: aws.String("No One You Know"),
			},
			"SongTitle": {
				S: aws.String("Call Me Today"),
			},
		},
		ReturnConsumedCapacity: aws.String("TOTAL"),
		TableName:              aws.String("Music"),
	}
	dynamoc.PutItem(input)
	// Get the span created for dynamo call
	span := tp.Recorder.GetSpans()[0]
	// Test related fields
	assert.Equal(t, constants.ClassNames["DYNAMODB"], span.ClassName)
	assert.Equal(t, constants.DomainNames["DB"], span.DomainName)
	assert.Equal(t, constants.DynamoDBRequestTypes["PutItem"], span.Tags[constants.SpanTags["OPERATION_TYPE"]])
	assert.Equal(t, "dynamodb.us-west-2.amazonaws.com", span.Tags[constants.DBTags["DB_INSTANCE"]])
	assert.Equal(t, "Music", span.Tags[constants.AwsDynamoDBTags["TABLE_NAME"]])
	assert.Equal(t, constants.DynamoDBRequestTypes["PutItem"], span.Tags[constants.DBTags["DB_STATEMENT_TYPE"]])
	assert.Equal(t, "PutItem", span.Tags[constants.AwsSDKTags["REQUEST_NAME"]])
	assert.Equal(t, true, span.Tags[constants.SpanTags["TOPOLOGY_VERTEX"]])
	assert.Equal(t, constants.AwsLambdaApplicationDomain, span.Tags[constants.SpanTags["TRIGGER_DOMAIN_NAME"]])
	assert.Equal(t, constants.AwsLambdaApplicationClass, span.Tags[constants.SpanTags["TRIGGER_CLASS_NAME"]])
	exp, err := json.Marshal(input.Item)
	if err != nil {
		t.Errorf("Couldn't marshal dynamodb input: %v", err)
	}
	got, err := json.Marshal(span.Tags[constants.DBTags["DB_STATEMENT"]])
	if err != nil {
		t.Errorf("Couldn't marshal db_statement tag in the span")
	}
	assert.Equal(t, exp, got)

	expectedTraceLinks := []string{
		"us-west-2:Music:1554912000:SAVE:cd2ecd1787d28c7d589601c6456b2e55",
		"us-west-2:Music:1554912001:SAVE:cd2ecd1787d28c7d589601c6456b2e55",
		"us-west-2:Music:1554912002:SAVE:cd2ecd1787d28c7d589601c6456b2e55",
	}

	assert.Equal(t, expectedTraceLinks, span.Tags[constants.SpanTags["TRACE_LINKS"]])

	// Clear tracer
	tp.Reset()
}

func TestDynamoDBUpdateItem(t *testing.T) {
	// Initilize trace plugin to set GlobalTracer of opentracing
	tp := trace.New()

	// Create a session and wrap it
	sess := getSessionWithDateResponseHeader()
	dynamoc := dynamodb.New(sess)

	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":r": {
				S: aws.String("test"),
			},
		},
		TableName: aws.String("Music"),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String("keyid"),
			},
		},
		ReturnValues:     aws.String("UPDATED_NEW"),
		UpdateExpression: aws.String("SET message = :r"),
	}
	exp, err := json.Marshal(input.Key)
	if err != nil {
		t.Errorf("Couldn't marshal dynamodb input: %v", err)
	}

	// Actual call
	dynamoc.UpdateItem(input)

	// Get the span created for dynamo call
	span := tp.Recorder.GetSpans()[0]
	// Test related fields
	assert.Equal(t, constants.ClassNames["DYNAMODB"], span.ClassName)
	assert.Equal(t, constants.DomainNames["DB"], span.DomainName)
	assert.Equal(t, constants.DynamoDBRequestTypes["UpdateItem"], span.Tags[constants.SpanTags["OPERATION_TYPE"]])
	assert.Equal(t, "dynamodb.us-west-2.amazonaws.com", span.Tags[constants.DBTags["DB_INSTANCE"]])
	assert.Equal(t, "Music", span.Tags[constants.AwsDynamoDBTags["TABLE_NAME"]])
	assert.Equal(t, constants.DynamoDBRequestTypes["UpdateItem"], span.Tags[constants.DBTags["DB_STATEMENT_TYPE"]])
	assert.Equal(t, "UpdateItem", span.Tags[constants.AwsSDKTags["REQUEST_NAME"]])
	assert.Equal(t, true, span.Tags[constants.SpanTags["TOPOLOGY_VERTEX"]])
	assert.Equal(t, constants.AwsLambdaApplicationDomain, span.Tags[constants.SpanTags["TRIGGER_DOMAIN_NAME"]])
	assert.Equal(t, constants.AwsLambdaApplicationClass, span.Tags[constants.SpanTags["TRIGGER_CLASS_NAME"]])

	got, err := json.Marshal(span.Tags[constants.DBTags["DB_STATEMENT"]])
	if err != nil {
		t.Errorf("Couldn't marshal db_statement tag in the span")
	}
	assert.Equal(t, exp, got)

	expectedTraceLinks := []string{
		"us-west-2:Music:1554912000:SAVE:214e7d85ccee118350d24b06f2c33d9c",
		"us-west-2:Music:1554912001:SAVE:214e7d85ccee118350d24b06f2c33d9c",
		"us-west-2:Music:1554912002:SAVE:214e7d85ccee118350d24b06f2c33d9c",
	}

	assert.Equal(t, expectedTraceLinks, span.Tags[constants.SpanTags["TRACE_LINKS"]])

	// Clear tracer
	tp.Reset()
}

func TestDynamoDBUpdateItemAttributeUpdate(t *testing.T) {
	// Initilize trace plugin to set GlobalTracer of opentracing
	tp := trace.New()
	config.DynamoDBTraceInjectionEnabled = true

	// Create a session and wrap it
	sess := getSessionWithDateResponseHeader()
	dynamoc := dynamodb.New(sess)

	input := &dynamodb.UpdateItemInput{
		AttributeUpdates: map[string]*dynamodb.AttributeValueUpdate{
			"Genre": {
				Action: aws.String("PUT"),
				Value: &dynamodb.AttributeValue{
					S: aws.String("Rock"),
				},
			},
		},
		TableName: aws.String("Music"),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String("keyid"),
			},
		},
	}
	exp, err := json.Marshal(input.Key)
	if err != nil {
		t.Errorf("Couldn't marshal dynamodb input: %v", err)
	}

	// Actual call
	dynamoc.UpdateItem(input)

	// Get the span created for dynamo call
	span := tp.Recorder.GetSpans()[0]
	// Test related fields
	assert.Equal(t, constants.ClassNames["DYNAMODB"], span.ClassName)
	assert.Equal(t, constants.DomainNames["DB"], span.DomainName)
	assert.Equal(t, constants.DynamoDBRequestTypes["UpdateItem"], span.Tags[constants.SpanTags["OPERATION_TYPE"]])
	assert.Equal(t, "dynamodb.us-west-2.amazonaws.com", span.Tags[constants.DBTags["DB_INSTANCE"]])
	assert.Equal(t, "Music", span.Tags[constants.AwsDynamoDBTags["TABLE_NAME"]])
	assert.Equal(t, constants.DynamoDBRequestTypes["UpdateItem"], span.Tags[constants.DBTags["DB_STATEMENT_TYPE"]])
	assert.Equal(t, "UpdateItem", span.Tags[constants.AwsSDKTags["REQUEST_NAME"]])
	assert.Equal(t, true, span.Tags[constants.SpanTags["TOPOLOGY_VERTEX"]])
	assert.Equal(t, constants.AwsLambdaApplicationDomain, span.Tags[constants.SpanTags["TRIGGER_DOMAIN_NAME"]])
	assert.Equal(t, constants.AwsLambdaApplicationClass, span.Tags[constants.SpanTags["TRIGGER_CLASS_NAME"]])

	got, err := json.Marshal(span.Tags[constants.DBTags["DB_STATEMENT"]])
	if err != nil {
		t.Errorf("Couldn't marshal db_statement tag in the span")
	}
	assert.Equal(t, exp, got)
	assert.Equal(t, []string{"SAVE:" + span.Context.SpanID}, span.Tags[constants.SpanTags["TRACE_LINKS"]])

	// Clear tracer
	tp.Reset()
}

func TestDynamoDeleteItem(t *testing.T) {
	// Initilize trace plugin to set GlobalTracer of opentracing
	tp := trace.New()

	// Create a session and wrap it
	sess := getSessionWithDateResponseHeader()
	dynamoc := dynamodb.New(sess)

	input := &dynamodb.DeleteItemInput{
		TableName: aws.String("Music"),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String("keyid"),
			},
		},
	}

	exp, err := json.Marshal(input.Key)
	if err != nil {
		t.Errorf("Couldn't marshal dynamodb input: %v", err)
	}

	dynamoc.DeleteItem(input)

	// Get the span created for dynamo call
	span := tp.Recorder.GetSpans()[0]
	// Test related fields
	assert.Equal(t, constants.ClassNames["DYNAMODB"], span.ClassName)
	assert.Equal(t, constants.DomainNames["DB"], span.DomainName)
	assert.Equal(t, constants.DynamoDBRequestTypes["DeleteItem"], span.Tags[constants.SpanTags["OPERATION_TYPE"]])
	assert.Equal(t, "dynamodb.us-west-2.amazonaws.com", span.Tags[constants.DBTags["DB_INSTANCE"]])
	assert.Equal(t, "Music", span.Tags[constants.AwsDynamoDBTags["TABLE_NAME"]])
	assert.Equal(t, constants.DynamoDBRequestTypes["DeleteItem"], span.Tags[constants.DBTags["DB_STATEMENT_TYPE"]])
	assert.Equal(t, "DeleteItem", span.Tags[constants.AwsSDKTags["REQUEST_NAME"]])
	assert.Equal(t, true, span.Tags[constants.SpanTags["TOPOLOGY_VERTEX"]])
	assert.Equal(t, constants.AwsLambdaApplicationDomain, span.Tags[constants.SpanTags["TRIGGER_DOMAIN_NAME"]])
	assert.Equal(t, constants.AwsLambdaApplicationClass, span.Tags[constants.SpanTags["TRIGGER_CLASS_NAME"]])

	got, err := json.Marshal(span.Tags[constants.DBTags["DB_STATEMENT"]])
	if err != nil {
		t.Errorf("Couldn't marshal db_statement tag in the span")
	}
	assert.Equal(t, exp, got)

	expectedTraceLinks := []string{
		"us-west-2:Music:1554912000:DELETE:214e7d85ccee118350d24b06f2c33d9c",
		"us-west-2:Music:1554912001:DELETE:214e7d85ccee118350d24b06f2c33d9c",
		"us-west-2:Music:1554912002:DELETE:214e7d85ccee118350d24b06f2c33d9c",
	}

	assert.Equal(t, expectedTraceLinks, span.Tags[constants.SpanTags["TRACE_LINKS"]])

	// Clear tracer
	tp.Reset()
}

func TestDynamoDBPutItemTraceEnabled(t *testing.T) {
	// Initilize trace plugin to set GlobalTracer of opentracing
	tp := trace.New()
	config.DynamoDBTraceInjectionEnabled = true
	// Create a session and wrap it
	sess := Wrap(sess)
	dynamoc := dynamodb.New(sess)
	// Actual call
	input := &dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{
			"AlbumTitle": {
				S: aws.String("Somewhat Famous"),
			},
			"Artist": {
				S: aws.String("No One You Know"),
			},
			"SongTitle": {
				S: aws.String("Call Me Today"),
			},
		},
		ReturnConsumedCapacity: aws.String("TOTAL"),
		TableName:              aws.String("Music"),
	}

	exp, err := json.Marshal(input.Item)
	if err != nil {
		t.Errorf("Couldn't marshal dynamodb input: %v", err)
	}

	dynamoc.PutItem(input)
	// Get the span created for dynamo call
	span := tp.Recorder.GetSpans()[0]
	// Test related fields
	assert.Equal(t, constants.ClassNames["DYNAMODB"], span.ClassName)
	assert.Equal(t, constants.DomainNames["DB"], span.DomainName)
	assert.Equal(t, constants.DynamoDBRequestTypes["PutItem"], span.Tags[constants.SpanTags["OPERATION_TYPE"]])
	assert.Equal(t, "dynamodb.us-west-2.amazonaws.com", span.Tags[constants.DBTags["DB_INSTANCE"]])
	assert.Equal(t, "Music", span.Tags[constants.AwsDynamoDBTags["TABLE_NAME"]])
	assert.Equal(t, constants.DynamoDBRequestTypes["PutItem"], span.Tags[constants.DBTags["DB_STATEMENT_TYPE"]])
	assert.Equal(t, "PutItem", span.Tags[constants.AwsSDKTags["REQUEST_NAME"]])
	assert.Equal(t, true, span.Tags[constants.SpanTags["TOPOLOGY_VERTEX"]])
	assert.Equal(t, constants.AwsLambdaApplicationDomain, span.Tags[constants.SpanTags["TRIGGER_DOMAIN_NAME"]])
	assert.Equal(t, constants.AwsLambdaApplicationClass, span.Tags[constants.SpanTags["TRIGGER_CLASS_NAME"]])

	got, err := json.Marshal(span.Tags[constants.DBTags["DB_STATEMENT"]])
	if err != nil {
		t.Errorf("Couldn't marshal db_statement tag in the span")
	}
	assert.Equal(t, exp, got)
	assert.Equal(t, []string{"SAVE:" + span.Context.SpanID}, span.Tags[constants.SpanTags["TRACE_LINKS"]])

	// Clear tracer
	tp.Reset()
}

func TestDynamoDBUpdateItemTraceEnabled(t *testing.T) {
	// Initilize trace plugin to set GlobalTracer of opentracing
	tp := trace.New()
	config.DynamoDBTraceInjectionEnabled = true
	// Create a session and wrap it
	sess := Wrap(sess)
	dynamoc := dynamodb.New(sess)

	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":r": {
				S: aws.String("test"),
			},
		},
		TableName: aws.String("Music"),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String("keyid"),
			},
		},
		ReturnValues:     aws.String("UPDATED_NEW"),
		UpdateExpression: aws.String("SET message = :r"),
	}
	exp, err := json.Marshal(input.Key)
	if err != nil {
		t.Errorf("Couldn't marshal dynamodb input: %v", err)
	}

	// Actual call
	dynamoc.UpdateItem(input)

	// Get the span created for dynamo call
	span := tp.Recorder.GetSpans()[0]
	// Test related fields
	assert.Equal(t, constants.ClassNames["DYNAMODB"], span.ClassName)
	assert.Equal(t, constants.DomainNames["DB"], span.DomainName)
	assert.Equal(t, constants.DynamoDBRequestTypes["UpdateItem"], span.Tags[constants.SpanTags["OPERATION_TYPE"]])
	assert.Equal(t, "dynamodb.us-west-2.amazonaws.com", span.Tags[constants.DBTags["DB_INSTANCE"]])
	assert.Equal(t, "Music", span.Tags[constants.AwsDynamoDBTags["TABLE_NAME"]])
	assert.Equal(t, constants.DynamoDBRequestTypes["UpdateItem"], span.Tags[constants.DBTags["DB_STATEMENT_TYPE"]])
	assert.Equal(t, "UpdateItem", span.Tags[constants.AwsSDKTags["REQUEST_NAME"]])
	assert.Equal(t, true, span.Tags[constants.SpanTags["TOPOLOGY_VERTEX"]])
	assert.Equal(t, constants.AwsLambdaApplicationDomain, span.Tags[constants.SpanTags["TRIGGER_DOMAIN_NAME"]])
	assert.Equal(t, constants.AwsLambdaApplicationClass, span.Tags[constants.SpanTags["TRIGGER_CLASS_NAME"]])

	got, err := json.Marshal(span.Tags[constants.DBTags["DB_STATEMENT"]])
	if err != nil {
		t.Errorf("Couldn't marshal db_statement tag in the span")
	}
	assert.Equal(t, exp, got)
	assert.Equal(t, []string{"SAVE:" + span.Context.SpanID}, span.Tags[constants.SpanTags["TRACE_LINKS"]])

	// Clear tracer
	tp.Reset()
}

func TestDynamoDeleteItemTraceEnabled(t *testing.T) {
	// Initilize trace plugin to set GlobalTracer of opentracing
	tp := trace.New()
	config.DynamoDBTraceInjectionEnabled = true
	// Create a session and wrap it
	sess := Wrap(sess)
	dynamoc := dynamodb.New(sess)

	input := &dynamodb.DeleteItemInput{
		TableName: aws.String("Music"),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String("keyid"),
			},
		},
	}

	exp, err := json.Marshal(input.Key)
	if err != nil {
		t.Errorf("Couldn't marshal dynamodb input: %v", err)
	}

	dynamoc.DeleteItem(input)

	// Get the span created for dynamo call
	span := tp.Recorder.GetSpans()[0]
	// Test related fields
	assert.Equal(t, constants.ClassNames["DYNAMODB"], span.ClassName)
	assert.Equal(t, constants.DomainNames["DB"], span.DomainName)
	assert.Equal(t, constants.DynamoDBRequestTypes["DeleteItem"], span.Tags[constants.SpanTags["OPERATION_TYPE"]])
	assert.Equal(t, "dynamodb.us-west-2.amazonaws.com", span.Tags[constants.DBTags["DB_INSTANCE"]])
	assert.Equal(t, "Music", span.Tags[constants.AwsDynamoDBTags["TABLE_NAME"]])
	assert.Equal(t, constants.DynamoDBRequestTypes["DeleteItem"], span.Tags[constants.DBTags["DB_STATEMENT_TYPE"]])
	assert.Equal(t, "DeleteItem", span.Tags[constants.AwsSDKTags["REQUEST_NAME"]])
	assert.Equal(t, true, span.Tags[constants.SpanTags["TOPOLOGY_VERTEX"]])
	assert.Equal(t, constants.AwsLambdaApplicationDomain, span.Tags[constants.SpanTags["TRIGGER_DOMAIN_NAME"]])
	assert.Equal(t, constants.AwsLambdaApplicationClass, span.Tags[constants.SpanTags["TRIGGER_CLASS_NAME"]])

	got, err := json.Marshal(span.Tags[constants.DBTags["DB_STATEMENT"]])
	if err != nil {
		t.Errorf("Couldn't marshal db_statement tag in the span")
	}
	assert.Equal(t, exp, got)

	assert.Equal(t, input.ReturnValues, aws.String("ALL_OLD"))

	// Clear tracer
	tp.Reset()
}

func TestDynamoDBGetItem(t *testing.T) {
	// Initilize trace plugin to set GlobalTracer of opentracing
	tp := trace.New()
	// Create a session and wrap it
	sess := Wrap(sess)
	dynamoc := dynamodb.New(sess)
	// Actual call
	input := &dynamodb.GetItemInput{
		TableName: aws.String("users-staging"),
		Key: map[string]*dynamodb.AttributeValue{
			"id": &dynamodb.AttributeValue{
				S: aws.String("1001"),
			},
			"name": &dynamodb.AttributeValue{
				B: []byte("{1:10, 2:20}"),
			},
		},
	}
	dynamoc.GetItem(input)
	// Get the span created for dynamo call
	span := tp.Recorder.GetSpans()[0]
	// Test related fields
	assert.Equal(t, constants.ClassNames["DYNAMODB"], span.ClassName)
	assert.Equal(t, constants.DomainNames["DB"], span.DomainName)
	assert.Equal(t, constants.DynamoDBRequestTypes["GetItem"], span.Tags[constants.SpanTags["OPERATION_TYPE"]])
	assert.Equal(t, "dynamodb.us-west-2.amazonaws.com", span.Tags[constants.DBTags["DB_INSTANCE"]])
	assert.Equal(t, "users-staging", span.Tags[constants.AwsDynamoDBTags["TABLE_NAME"]])
	assert.Equal(t, constants.DynamoDBRequestTypes["GetItem"], span.Tags[constants.DBTags["DB_STATEMENT_TYPE"]])
	assert.Equal(t, "GetItem", span.Tags[constants.AwsSDKTags["REQUEST_NAME"]])
	assert.Equal(t, true, span.Tags[constants.SpanTags["TOPOLOGY_VERTEX"]])
	assert.Equal(t, constants.AwsLambdaApplicationDomain, span.Tags[constants.SpanTags["TRIGGER_DOMAIN_NAME"]])
	assert.Equal(t, constants.AwsLambdaApplicationClass, span.Tags[constants.SpanTags["TRIGGER_CLASS_NAME"]])

	exp, err := json.Marshal(input.Key)
	if err != nil {
		t.Errorf("Couldn't marshal dynamodb input: %v", err)
	}
	got, err := json.Marshal(span.Tags[constants.DBTags["DB_STATEMENT"]])
	if err != nil {
		t.Errorf("Couldn't marshal db_statement tag in the span")
	}
	assert.Equal(t, exp, got)

	assert.Equal(t, nil, span.Tags[constants.SpanTags["TRACE_LINKS"]])

	// Clear tracer
	tp.Reset()
}

func TestSNSPublish(t *testing.T) {
	config.MaskSNSMessage = false
	// Initilize trace plugin to set GlobalTracer of opentracing
	tp := trace.New()

	// Create a session and wrap it
	sess := getSessionWithSnsResponse()
	snsc := sns.New(sess)

	// Params will be sent to the publish call included here is the bare minimum params to send a message
	message := "foobar"
	params := &sns.PublishInput{
		Message:  aws.String(message),
		TopicArn: aws.String("arn:aws:sns:us-west-2:123456789012:gsg-signup-notifications"),
	}

	// Call to publish message
	snsc.Publish(params)

	// Get the span created for dynamo call
	span := tp.Recorder.GetSpans()[0]
	// Test related fields
	assert.Equal(t, constants.ClassNames["SNS"], span.ClassName)
	assert.Equal(t, constants.DomainNames["MESSAGING"], span.DomainName)
	assert.Equal(t, "gsg-signup-notifications", span.Tags[constants.AwsSNSTags["TOPIC_NAME"]])
	assert.Equal(t, "WRITE", span.Tags[constants.SpanTags["OPERATION_TYPE"]])
	assert.Equal(t, "Publish", span.Tags[constants.AwsSDKTags["REQUEST_NAME"]])
	assert.Equal(t, true, span.Tags[constants.SpanTags["TOPOLOGY_VERTEX"]])
	assert.Equal(t, constants.AwsLambdaApplicationDomain, span.Tags[constants.SpanTags["TRIGGER_DOMAIN_NAME"]])
	assert.Equal(t, constants.AwsLambdaApplicationClass, span.Tags[constants.SpanTags["TRIGGER_CLASS_NAME"]])
	assert.Equal(t, message, span.Tags[constants.AwsSNSTags["MESSAGE"]])

	assert.Equal(t, []string{"95df01b4-ee98-5cb9-9903-4c221d41eb5e"}, span.Tags[constants.SpanTags["TRACE_LINKS"]])

	// Clear tracer
	tp.Reset()
}

func TestSNSPublishWithMaskedMessage(t *testing.T) {
	config.MaskSNSMessage = true
	// Initilize trace plugin to set GlobalTracer of opentracing
	tp := trace.New()

	// Create a session and wrap it
	sess := getSessionWithSnsResponse()
	snsc := sns.New(sess)

	// Params will be sent to the publish call included here is the bare minimum params to send a message
	message := "foobar"
	params := &sns.PublishInput{
		Message:  aws.String(message),
		TopicArn: aws.String("arn:aws:sns:us-west-2:123456789012:gsg-signup-notifications"),
	}

	// Call to publish message
	snsc.Publish(params)

	// Get the span created for dynamo call
	span := tp.Recorder.GetSpans()[0]
	// Test related fields
	assert.Equal(t, constants.ClassNames["SNS"], span.ClassName)
	assert.Equal(t, constants.DomainNames["MESSAGING"], span.DomainName)
	assert.Equal(t, "gsg-signup-notifications", span.Tags[constants.AwsSNSTags["TOPIC_NAME"]])
	assert.Equal(t, "WRITE", span.Tags[constants.SpanTags["OPERATION_TYPE"]])
	assert.Equal(t, "Publish", span.Tags[constants.AwsSDKTags["REQUEST_NAME"]])
	assert.Equal(t, true, span.Tags[constants.SpanTags["TOPOLOGY_VERTEX"]])
	assert.Equal(t, constants.AwsLambdaApplicationDomain, span.Tags[constants.SpanTags["TRIGGER_DOMAIN_NAME"]])
	assert.Equal(t, constants.AwsLambdaApplicationClass, span.Tags[constants.SpanTags["TRIGGER_CLASS_NAME"]])
	assert.Nil(t, span.Tags[constants.AwsSNSTags["MESSAGE"]])

	assert.Equal(t, []string{"95df01b4-ee98-5cb9-9903-4c221d41eb5e"}, span.Tags[constants.SpanTags["TRACE_LINKS"]])

	// Clear tracer
	tp.Reset()
}

func TestSNSGetTopic(t *testing.T) {
	// Initilize trace plugin to set GlobalTracer of opentracing
	tp := trace.New()
	// Create a session and wrap it
	sess := Wrap(sess)
	snsc := sns.New(sess)
	// Actual call
	snsc.GetTopicAttributes(&sns.GetTopicAttributesInput{
		TopicArn: aws.String("arn:aws:sns:us-west-2:123456789012:gsg-signup-notifications"),
	})
	// Get the span created for dynamo call
	span := tp.Recorder.GetSpans()[0]
	// Test related fields
	assert.Equal(t, constants.ClassNames["SNS"], span.ClassName)
	assert.Equal(t, constants.DomainNames["MESSAGING"], span.DomainName)
	assert.Equal(t, "gsg-signup-notifications", span.Tags[constants.AwsSNSTags["TOPIC_NAME"]])
	assert.Equal(t, "", span.Tags[constants.SpanTags["OPERATION_TYPE"]])
	assert.Equal(t, "GetTopicAttributes", span.Tags[constants.AwsSDKTags["REQUEST_NAME"]])
	assert.Equal(t, true, span.Tags[constants.SpanTags["TOPOLOGY_VERTEX"]])
	assert.Equal(t, constants.AwsLambdaApplicationDomain, span.Tags[constants.SpanTags["TRIGGER_DOMAIN_NAME"]])
	assert.Equal(t, constants.AwsLambdaApplicationClass, span.Tags[constants.SpanTags["TRIGGER_CLASS_NAME"]])

	// Clear tracer
	tp.Reset()
}

func TestSNSCreateTopic(t *testing.T) {
	// Initilize trace plugin to set GlobalTracer of opentracing
	tp := trace.New()
	// Create a session and wrap it
	sess := Wrap(sess)
	snsc := sns.New(sess)
	// Actual call
	snsc.CreateTopic(&sns.CreateTopicInput{
		Name: aws.String("foobar"),
	})
	// Get the span created for dynamo call
	span := tp.Recorder.GetSpans()[0]
	// Test related fields
	assert.Equal(t, constants.ClassNames["SNS"], span.ClassName)
	assert.Equal(t, constants.DomainNames["MESSAGING"], span.DomainName)
	assert.Equal(t, "foobar", span.Tags[constants.AwsSNSTags["TOPIC_NAME"]])
	assert.Equal(t, "", span.Tags[constants.SpanTags["OPERATION_TYPE"]])
	assert.Equal(t, "CreateTopic", span.Tags[constants.AwsSDKTags["REQUEST_NAME"]])
	assert.Equal(t, true, span.Tags[constants.SpanTags["TOPOLOGY_VERTEX"]])
	assert.Equal(t, constants.AwsLambdaApplicationDomain, span.Tags[constants.SpanTags["TRIGGER_DOMAIN_NAME"]])
	assert.Equal(t, constants.AwsLambdaApplicationClass, span.Tags[constants.SpanTags["TRIGGER_CLASS_NAME"]])

	// Clear tracer
	tp.Reset()
}

func TestKinesisPutRecord(t *testing.T) {
	// Initilize trace plugin to set GlobalTracer of opentracing
	tp := trace.New()

	// Create a session and wrap it
	sess := getSessionWithKinesisResponse()
	knssc := kinesis.New(sess)
	// Actual call
	knssc.PutRecord(&kinesis.PutRecordInput{
		Data:         []byte("message"),
		PartitionKey: aws.String("key1"),
		StreamName:   aws.String("Foo Stream"),
	})
	// Get the span created for dynamo call
	span := tp.Recorder.GetSpans()[0]
	// Test related fields
	assert.Equal(t, constants.ClassNames["KINESIS"], span.ClassName)
	assert.Equal(t, constants.DomainNames["STREAM"], span.DomainName)
	assert.Equal(t, "Foo Stream", span.Tags[constants.AwsKinesisTags["STREAM_NAME"]])
	assert.Equal(t, "WRITE", span.Tags[constants.SpanTags["OPERATION_TYPE"]])
	assert.Equal(t, "PutRecord", span.Tags[constants.AwsSDKTags["REQUEST_NAME"]])
	assert.Equal(t, true, span.Tags[constants.SpanTags["TOPOLOGY_VERTEX"]])
	assert.Equal(t, constants.AwsLambdaApplicationDomain, span.Tags[constants.SpanTags["TRIGGER_DOMAIN_NAME"]])
	assert.Equal(t, constants.AwsLambdaApplicationClass, span.Tags[constants.SpanTags["TRIGGER_CLASS_NAME"]])

	expectedTraceLinks := []string{"us-west-2:Foo Stream:shardId-000000000000:49568167373333333333333333333333333333333333333333333333"}

	assert.Equal(t, expectedTraceLinks, span.Tags[constants.SpanTags["TRACE_LINKS"]])

	// Clear tracer
	tp.Reset()
}

func TestKinesisPutRecords(t *testing.T) {
	// Initilize trace plugin to set GlobalTracer of opentracing
	tp := trace.New()

	// Create a session and wrap it
	sess := getSessionWithKinesisPutRecordsResponse()
	knssc := kinesis.New(sess)

	entries := []*kinesis.PutRecordsRequestEntry{
		&kinesis.PutRecordsRequestEntry{
			Data:         []byte("1"),
			PartitionKey: aws.String("key1"),
		},
		&kinesis.PutRecordsRequestEntry{
			Data:         []byte("2"),
			PartitionKey: aws.String("key2"),
		},
	}
	// Actual call
	knssc.PutRecords(&kinesis.PutRecordsInput{
		Records:    entries,
		StreamName: aws.String("Foo Stream"),
	})
	// Get the span created for dynamo call
	span := tp.Recorder.GetSpans()[0]
	// Test related fields
	assert.Equal(t, constants.ClassNames["KINESIS"], span.ClassName)
	assert.Equal(t, constants.DomainNames["STREAM"], span.DomainName)
	assert.Equal(t, "Foo Stream", span.Tags[constants.AwsKinesisTags["STREAM_NAME"]])
	assert.Equal(t, "WRITE", span.Tags[constants.SpanTags["OPERATION_TYPE"]])
	assert.Equal(t, "PutRecords", span.Tags[constants.AwsSDKTags["REQUEST_NAME"]])
	assert.Equal(t, true, span.Tags[constants.SpanTags["TOPOLOGY_VERTEX"]])
	assert.Equal(t, constants.AwsLambdaApplicationDomain, span.Tags[constants.SpanTags["TRIGGER_DOMAIN_NAME"]])
	assert.Equal(t, constants.AwsLambdaApplicationClass, span.Tags[constants.SpanTags["TRIGGER_CLASS_NAME"]])

	expectedTraceLinks := []string{
		"us-west-2:Foo Stream:shardId-000000000000:49568167373333333333333333333333333333333333333333333333",
		"us-west-2:Foo Stream:shardId-000000000000:49568167374444444444444444444444444444444444444444444444",
	}

	assert.ElementsMatch(t, expectedTraceLinks, span.Tags[constants.SpanTags["TRACE_LINKS"]])

	// Clear tracer
	tp.Reset()
}

func TestKinesisGetRecord(t *testing.T) {
	// Initilize trace plugin to set GlobalTracer of opentracing
	tp := trace.New()
	// Create a session and wrap it
	sess := Wrap(sess)
	knssc := kinesis.New(sess)
	// Actual call
	knssc.GetRecords(&kinesis.GetRecordsInput{
		ShardIterator: aws.String("foo"),
	})
	// Get the span created for dynamo call
	span := tp.Recorder.GetSpans()[0]
	// Test related fields
	assert.Equal(t, constants.ClassNames["KINESIS"], span.ClassName)
	assert.Equal(t, constants.DomainNames["STREAM"], span.DomainName)
	assert.Equal(t, "", span.Tags[constants.AwsKinesisTags["STREAM_NAME"]])
	assert.Equal(t, "READ", span.Tags[constants.SpanTags["OPERATION_TYPE"]])
	assert.Equal(t, "GetRecords", span.Tags[constants.AwsSDKTags["REQUEST_NAME"]])
	assert.Equal(t, true, span.Tags[constants.SpanTags["TOPOLOGY_VERTEX"]])
	assert.Equal(t, constants.AwsLambdaApplicationDomain, span.Tags[constants.SpanTags["TRIGGER_DOMAIN_NAME"]])
	assert.Equal(t, constants.AwsLambdaApplicationClass, span.Tags[constants.SpanTags["TRIGGER_CLASS_NAME"]])

	// Clear tracer
	tp.Reset()
}

func TestFirehosePutRecord(t *testing.T) {
	// Initilize trace plugin to set GlobalTracer of opentracing
	tp := trace.New()

	// Create a session and wrap it
	sess := getSessionWithDateResponseHeader()
	fhc := firehose.New(sess)
	// Actual call
	fhc.PutRecord(&firehose.PutRecordInput{
		DeliveryStreamName: aws.String("Foo Stream"),
		Record:             &firehose.Record{Data: []byte("test")},
	})
	// Get the span created for dynamo call
	span := tp.Recorder.GetSpans()[0]
	// Test related fields
	assert.Equal(t, constants.ClassNames["FIREHOSE"], span.ClassName)
	assert.Equal(t, constants.DomainNames["STREAM"], span.DomainName)
	assert.Equal(t, "Foo Stream", span.Tags[constants.AwsFirehoseTags["STREAM_NAME"]])
	assert.Equal(t, "WRITE", span.Tags[constants.SpanTags["OPERATION_TYPE"]])
	assert.Equal(t, "PutRecord", span.Tags[constants.AwsSDKTags["REQUEST_NAME"]])
	assert.Equal(t, true, span.Tags[constants.SpanTags["TOPOLOGY_VERTEX"]])
	assert.Equal(t, constants.AwsLambdaApplicationDomain, span.Tags[constants.SpanTags["TRIGGER_DOMAIN_NAME"]])
	assert.Equal(t, constants.AwsLambdaApplicationClass, span.Tags[constants.SpanTags["TRIGGER_CLASS_NAME"]])

	expectedTraceLinks := []string{
		"us-west-2:Foo Stream:1554912000:098f6bcd4621d373cade4e832627b4f6",
		"us-west-2:Foo Stream:1554912001:098f6bcd4621d373cade4e832627b4f6",
		"us-west-2:Foo Stream:1554912002:098f6bcd4621d373cade4e832627b4f6",
	}

	assert.ElementsMatch(t, expectedTraceLinks, span.Tags[constants.SpanTags["TRACE_LINKS"]])

	// Clear tracer
	tp.Reset()
}

func TestFirehosePutRecordBatch(t *testing.T) {
	// Initilize trace plugin to set GlobalTracer of opentracing
	tp := trace.New()

	// Create a session and wrap it
	sess := getSessionWithDateResponseHeader()
	fhc := firehose.New(sess)

	recordsBatchInput := &firehose.PutRecordBatchInput{}
	recordsBatchInput = recordsBatchInput.SetDeliveryStreamName(*aws.String("Foo Stream"))

	records := []*firehose.Record{
		&firehose.Record{Data: []byte("test")},
	}

	recordsBatchInput = recordsBatchInput.SetRecords(records)

	// Actual call
	fhc.PutRecordBatch(recordsBatchInput)

	// Get the span created for dynamo call
	span := tp.Recorder.GetSpans()[0]
	// Test related fields
	assert.Equal(t, constants.ClassNames["FIREHOSE"], span.ClassName)
	assert.Equal(t, constants.DomainNames["STREAM"], span.DomainName)
	assert.Equal(t, "Foo Stream", span.Tags[constants.AwsFirehoseTags["STREAM_NAME"]])
	assert.Equal(t, "WRITE", span.Tags[constants.SpanTags["OPERATION_TYPE"]])
	assert.Equal(t, "PutRecordBatch", span.Tags[constants.AwsSDKTags["REQUEST_NAME"]])
	assert.Equal(t, true, span.Tags[constants.SpanTags["TOPOLOGY_VERTEX"]])
	assert.Equal(t, constants.AwsLambdaApplicationDomain, span.Tags[constants.SpanTags["TRIGGER_DOMAIN_NAME"]])
	assert.Equal(t, constants.AwsLambdaApplicationClass, span.Tags[constants.SpanTags["TRIGGER_CLASS_NAME"]])

	expectedTraceLinks := []string{
		"us-west-2:Foo Stream:1554912000:098f6bcd4621d373cade4e832627b4f6",
		"us-west-2:Foo Stream:1554912001:098f6bcd4621d373cade4e832627b4f6",
		"us-west-2:Foo Stream:1554912002:098f6bcd4621d373cade4e832627b4f6",
	}

	assert.ElementsMatch(t, expectedTraceLinks, span.Tags[constants.SpanTags["TRACE_LINKS"]])

	// Clear tracer
	tp.Reset()
}

func TestS3GetObject(t *testing.T) {
	// Initilize trace plugin to set GlobalTracer of opentracing
	tp := trace.New()

	// Create a session and wrap it
	sess := getSessionWithS3Response()
	s3c := s3.New(sess)
	// Actual call
	s3c.GetObject(&s3.GetObjectInput{
		Bucket: aws.String("some-bucket-name"),
		Key:    aws.String("some-object-key"),
	})
	// Get the span created for dynamo call
	span := tp.Recorder.GetSpans()[0]
	// Test related fields
	assert.Equal(t, constants.ClassNames["S3"], span.ClassName)
	assert.Equal(t, constants.DomainNames["STORAGE"], span.DomainName)
	assert.Equal(t, "some-bucket-name", span.Tags[constants.AwsS3Tags["BUCKET_NAME"]])
	assert.Equal(t, "some-object-key", span.Tags[constants.AwsS3Tags["OBJECT_NAME"]])
	assert.Equal(t, "READ", span.Tags[constants.SpanTags["OPERATION_TYPE"]])
	assert.Equal(t, "GetObject", span.Tags[constants.AwsSDKTags["REQUEST_NAME"]])
	assert.Equal(t, true, span.Tags[constants.SpanTags["TOPOLOGY_VERTEX"]])
	assert.Equal(t, constants.AwsLambdaApplicationDomain, span.Tags[constants.SpanTags["TRIGGER_DOMAIN_NAME"]])
	assert.Equal(t, constants.AwsLambdaApplicationClass, span.Tags[constants.SpanTags["TRIGGER_CLASS_NAME"]])

	assert.Equal(t, []string{"C3D13FE58DE4C810"}, span.Tags[constants.SpanTags["TRACE_LINKS"]])

	// Clear tracer
	tp.Reset()
}

func TestLambdaInvoke(t *testing.T) {
	// Set application name
	application.ApplicationName = "test"
	// Initilize trace plugin to set GlobalTracer of opentracing
	tp := trace.New()

	// Create a session and wrap it
	sess := getSessionWithLambdaResponse()
	lambdac := lambda.New(sess)
	// Actual call
	input := &lambda.InvokeInput{
		FunctionName:   aws.String("a-lambda-function:42"),
		Payload:        []byte("\"foobar\""),
		InvocationType: aws.String("RequestResponse"),
		Qualifier:      aws.String("function-qualifier"),
	}
	lambdac.Invoke(input)
	// Get the span created for dynamo call
	span := tp.Recorder.GetSpans()[0]
	// Test related fields
	assert.Equal(t, constants.ClassNames["LAMBDA"], span.ClassName)
	assert.Equal(t, constants.DomainNames["API"], span.DomainName)
	assert.Equal(t, "a-lambda-function", span.Tags[constants.AwsLambdaTags["FUNCTION_NAME"]])
	assert.Equal(t, "RequestResponse", span.Tags[constants.AwsLambdaTags["INVOCATION_TYPE"]])
	assert.Equal(t, "function-qualifier", span.Tags[constants.AwsLambdaTags["FUNCTION_QUALIFIER"]])
	assert.Equal(t, "CALL", span.Tags[constants.SpanTags["OPERATION_TYPE"]])
	assert.Equal(t, "Invoke", span.Tags[constants.AwsSDKTags["REQUEST_NAME"]])
	assert.Equal(t, true, span.Tags[constants.SpanTags["TOPOLOGY_VERTEX"]])
	assert.Equal(t, constants.AwsLambdaApplicationDomain, span.Tags[constants.SpanTags["TRIGGER_DOMAIN_NAME"]])
	assert.Equal(t, constants.AwsLambdaApplicationClass, span.Tags[constants.SpanTags["TRIGGER_CLASS_NAME"]])

	assert.Equal(t, []string{"C3D13FE58DE4C810"}, span.Tags[constants.SpanTags["TRACE_LINKS"]])

	exp, err := json.Marshal(input.Payload)
	if err != nil {
		t.Errorf("Couldn't marshal lambda payload: %v", err)
	}
	got, err := json.Marshal(span.Tags[constants.AwsLambdaTags["INVOCATION_PAYLOAD"]])
	if err != nil {
		t.Errorf("Couldn't marshal lambda payload from span tags: %v", err)
	}
	assert.Equal(t, exp, got)

	clientContextExp := getBase64EncodedClientContext()
	clientContextGot := *input.ClientContext

	assert.Equal(t, clientContextExp, string(clientContextGot))
	// Clear tracer
	tp.Reset()
}

func TestLambdaInvokeFunctionArn(t *testing.T) {
	// Initilize trace plugin to set GlobalTracer of opentracing
	tp := trace.New()

	// Create a session and wrap it
	sess := getSessionWithLambdaResponse()
	lambdac := lambda.New(sess)
	// Actual call
	input := &lambda.InvokeInput{
		FunctionName:   aws.String("arn:aws:lambda:us-west-2:123456789012:function:a-lambda-function"),
		Payload:        []byte("\"foobar\""),
		InvocationType: aws.String("RequestResponse"),
		Qualifier:      aws.String("function-qualifier"),
	}
	lambdac.Invoke(input)
	// Get the span created for dynamo call
	span := tp.Recorder.GetSpans()[0]
	// Test related fields
	assert.Equal(t, "a-lambda-function", span.Tags[constants.AwsLambdaTags["FUNCTION_NAME"]])

	// Clear tracer
	tp.Reset()
}

func TestLambdaInvokeWithClientContext(t *testing.T) {
	// Set application name
	application.ApplicationName = "test"

	// Initilize trace plugin to set GlobalTracer of opentracing
	tp := trace.New()

	// Create a session and wrap it
	sess := getSessionWithLambdaResponse()
	lambdac := lambda.New(sess)
	// Actual call
	input := &lambda.InvokeInput{
		FunctionName:   aws.String("a-lambda-function"),
		Payload:        []byte("\"foobar\""),
		InvocationType: aws.String("RequestResponse"),
		Qualifier:      aws.String("function-qualifier"),
		ClientContext:  aws.String(getMockBase64EncodedClientContext()),
	}
	lambdac.Invoke(input)
	// Get the span created for dynamo call
	span := tp.Recorder.GetSpans()[0]
	// Test related fields
	assert.Equal(t, constants.ClassNames["LAMBDA"], span.ClassName)
	assert.Equal(t, constants.DomainNames["API"], span.DomainName)
	assert.Equal(t, "a-lambda-function", span.Tags[constants.AwsLambdaTags["FUNCTION_NAME"]])
	assert.Equal(t, "RequestResponse", span.Tags[constants.AwsLambdaTags["INVOCATION_TYPE"]])
	assert.Equal(t, "function-qualifier", span.Tags[constants.AwsLambdaTags["FUNCTION_QUALIFIER"]])
	assert.Equal(t, "CALL", span.Tags[constants.SpanTags["OPERATION_TYPE"]])
	assert.Equal(t, "Invoke", span.Tags[constants.AwsSDKTags["REQUEST_NAME"]])
	assert.Equal(t, true, span.Tags[constants.SpanTags["TOPOLOGY_VERTEX"]])
	assert.Equal(t, constants.AwsLambdaApplicationDomain, span.Tags[constants.SpanTags["TRIGGER_DOMAIN_NAME"]])
	assert.Equal(t, constants.AwsLambdaApplicationClass, span.Tags[constants.SpanTags["TRIGGER_CLASS_NAME"]])

	assert.Equal(t, []string{"C3D13FE58DE4C810"}, span.Tags[constants.SpanTags["TRACE_LINKS"]])

	exp, err := json.Marshal(input.Payload)
	if err != nil {
		t.Errorf("Couldn't marshal lambda payload: %v", err)
	}
	got, err := json.Marshal(span.Tags[constants.AwsLambdaTags["INVOCATION_PAYLOAD"]])
	if err != nil {
		t.Errorf("Couldn't marshal lambda payload from span tags: %v", err)
	}
	assert.Equal(t, exp, got)

	clientContextExp := getBase64EncodedClientContextWithMockParam()
	clientContextGot := *input.ClientContext

	assert.Equal(t, clientContextExp, string(clientContextGot))
	// Clear tracer
	tp.Reset()
}

func TestSQSSendMessage(t *testing.T) {
	config.MaskSQSMessage = false
	// Initilize trace plugin to set GlobalTracer of opentracing
	tp := trace.New()

	// Create a session and wrap it
	sess := getSessionWithSqsResponse()
	sqsc := sqs.New(sess)

	// Params will be sent to the publish call included here is the bare minimum params to send a message
	message := "foobar"
	params := &sqs.SendMessageInput{
		MessageBody: aws.String(message),
		QueueUrl:    aws.String("https://sqs.us-west-2.amazonaws.com/123456789012/test-queue"),
	}

	sqsc.SendMessage(params)

	// Get the span created for dynamo call
	span := tp.Recorder.GetSpans()[0]
	// Test related fields
	assert.Equal(t, constants.ClassNames["SQS"], span.ClassName)
	assert.Equal(t, constants.DomainNames["MESSAGING"], span.DomainName)
	assert.Equal(t, "test-queue", span.Tags[constants.AwsSQSTags["QUEUE_NAME"]])
	assert.Equal(t, "WRITE", span.Tags[constants.SpanTags["OPERATION_TYPE"]])
	assert.Equal(t, "SendMessage", span.Tags[constants.AwsSDKTags["REQUEST_NAME"]])
	assert.Equal(t, true, span.Tags[constants.SpanTags["TOPOLOGY_VERTEX"]])
	assert.Equal(t, constants.AwsLambdaApplicationDomain, span.Tags[constants.SpanTags["TRIGGER_DOMAIN_NAME"]])
	assert.Equal(t, constants.AwsLambdaApplicationClass, span.Tags[constants.SpanTags["TRIGGER_CLASS_NAME"]])
	assert.Equal(t, message, span.Tags[constants.AwsSQSTags["MESSAGE"]])

	assert.Equal(t, []string{"95df01b4-ee98-5cb9-9903-4c221d41eb5e"}, span.Tags[constants.SpanTags["TRACE_LINKS"]])

	// Clear tracer
	tp.Reset()
}

func TestSQSSendMessageWithMaskedMessage(t *testing.T) {
	config.MaskSQSMessage = true
	// Initilize trace plugin to set GlobalTracer of opentracing
	tp := trace.New()

	// Create a session and wrap it
	sess := getSessionWithSqsResponse()
	sqsc := sqs.New(sess)

	// Params will be sent to the publish call included here is the bare minimum params to send a message
	message := "foobar"
	params := &sqs.SendMessageInput{
		MessageBody: aws.String(message),
		QueueUrl:    aws.String("https://sqs.us-west-2.amazonaws.com/123456789012/test-queue"),
	}

	sqsc.SendMessage(params)

	// Get the span created for dynamo call
	span := tp.Recorder.GetSpans()[0]
	// Test related fields
	assert.Equal(t, constants.ClassNames["SQS"], span.ClassName)
	assert.Equal(t, constants.DomainNames["MESSAGING"], span.DomainName)
	assert.Equal(t, "test-queue", span.Tags[constants.AwsSQSTags["QUEUE_NAME"]])
	assert.Equal(t, "WRITE", span.Tags[constants.SpanTags["OPERATION_TYPE"]])
	assert.Equal(t, "SendMessage", span.Tags[constants.AwsSDKTags["REQUEST_NAME"]])
	assert.Equal(t, true, span.Tags[constants.SpanTags["TOPOLOGY_VERTEX"]])
	assert.Equal(t, constants.AwsLambdaApplicationDomain, span.Tags[constants.SpanTags["TRIGGER_DOMAIN_NAME"]])
	assert.Equal(t, constants.AwsLambdaApplicationClass, span.Tags[constants.SpanTags["TRIGGER_CLASS_NAME"]])
	assert.Nil(t, span.Tags[constants.AwsSQSTags["MESSAGE"]])

	assert.Equal(t, []string{"95df01b4-ee98-5cb9-9903-4c221d41eb5e"}, span.Tags[constants.SpanTags["TRACE_LINKS"]])

	// Clear tracer
	tp.Reset()
}

func TestSQSSendMessageBatch(t *testing.T) {
	// Initilize trace plugin to set GlobalTracer of opentracing
	tp := trace.New()

	// Create a session and wrap it
	sess := getSessionWithSqsBatchResponse()
	sqsc := sqs.New(sess)

	// Params will be sent to the publish call included here is the bare minimum params to send a message

	entries := []*sqs.SendMessageBatchRequestEntry{
		&sqs.SendMessageBatchRequestEntry{
			Id:          aws.String("1"),
			MessageBody: aws.String("test"),
		},
		&sqs.SendMessageBatchRequestEntry{
			Id:          aws.String("2"),
			MessageBody: aws.String("test"),
		},
	}
	params := &sqs.SendMessageBatchInput{
		Entries:  entries,
		QueueUrl: aws.String("https://sqs.us-west-2.amazonaws.com/123456789012/test-queue"),
	}

	sqsc.SendMessageBatch(params)

	// Get the span created for dynamo call
	span := tp.Recorder.GetSpans()[0]
	// Test related fields
	assert.Equal(t, constants.ClassNames["SQS"], span.ClassName)
	assert.Equal(t, constants.DomainNames["MESSAGING"], span.DomainName)
	assert.Equal(t, "test-queue", span.Tags[constants.AwsSQSTags["QUEUE_NAME"]])
	assert.Equal(t, "WRITE", span.Tags[constants.SpanTags["OPERATION_TYPE"]])
	assert.Equal(t, "SendMessageBatch", span.Tags[constants.AwsSDKTags["REQUEST_NAME"]])
	assert.Equal(t, true, span.Tags[constants.SpanTags["TOPOLOGY_VERTEX"]])
	assert.Equal(t, constants.AwsLambdaApplicationDomain, span.Tags[constants.SpanTags["TRIGGER_DOMAIN_NAME"]])
	assert.Equal(t, constants.AwsLambdaApplicationClass, span.Tags[constants.SpanTags["TRIGGER_CLASS_NAME"]])

	expectedTraceLinks := []string{"84df12b4-ee98-2cb8-1903-1c234d56eb7e", "95df01b4-ee98-5cb9-9903-4c221d41eb5e"}
	assert.ElementsMatch(t, expectedTraceLinks, span.Tags[constants.SpanTags["TRACE_LINKS"]])

	// Clear tracer
	tp.Reset()
}

func TestNonTracedService(t *testing.T) {
	// Initilize trace plugin to set GlobalTracer of opentracing
	tp := trace.New()
	// Create a session and wrap it
	sess := Wrap(sess)
	cwc := cloudwatch.New(sess)

	cwc.GetDashboard(&cloudwatch.GetDashboardInput{
		DashboardName: aws.String("foo"),
	})
	assert.Equal(t, 0, len(tp.Recorder.GetSpans()))
	tp.Reset()
}
