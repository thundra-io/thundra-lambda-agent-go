package thundraaws

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/firehose"
	"github.com/aws/aws-sdk-go/service/kinesis"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/stretchr/testify/assert"
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

var sess = Wrap(session.New(&aws.Config{
	HTTPClient: testClient,
	Region:     aws.String("us-west-2"),
	MaxRetries: aws.Int(0),
}))

func TestDynamoDBPutItem(t *testing.T) {
	// Initilize trace plugin to set GlobalTracer of opentracing
	tp := trace.New()
	// Create a session and wrap it
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
	// Clear tracer
	tp.Reset()
}

func TestDynamoDBGetItem(t *testing.T) {
	// Initilize trace plugin to set GlobalTracer of opentracing
	tp := trace.New()
	// Create a session and wrap it
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
	// Clear tracer
	tp.Reset()
}

func TestSNSGetTopic(t *testing.T) {
	// Initilize trace plugin to set GlobalTracer of opentracing
	tp := trace.New()
	// Create a session and wrap it
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
	knssc := kinesis.New(sess)
	// Actual call
	knssc.PutRecord(&kinesis.PutRecordInput{
		StreamName: aws.String("Foo Stream"),
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

	// Clear tracer
	tp.Reset()
}

func TestKinesisGetRecord(t *testing.T) {
	// Initilize trace plugin to set GlobalTracer of opentracing
	tp := trace.New()
	// Create a session and wrap it
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
	fhc := firehose.New(sess)
	// Actual call
	fhc.PutRecord(&firehose.PutRecordInput{
		DeliveryStreamName: aws.String("Foo Stream"),
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

	// Clear tracer
	tp.Reset()
}

func TestS3GetObject(t *testing.T) {
	// Initilize trace plugin to set GlobalTracer of opentracing
	tp := trace.New()
	// Create a session and wrap it
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

	// Clear tracer
	tp.Reset()
}
