package tredigo

import (
	"context"
	"testing"

	ot "github.com/opentracing/opentracing-go"
	"github.com/stretchr/testify/assert"
	"github.com/thundra-io/thundra-lambda-agent-go/application"
	"github.com/thundra-io/thundra-lambda-agent-go/constants"
	"github.com/thundra-io/thundra-lambda-agent-go/trace"
	"github.com/thundra-io/thundra-lambda-agent-go/tracer"
)

func TestRedigoDo(t *testing.T) {
	// Initilize trace plugin to set GlobalTracer of opentracing
	tp := trace.New()

	c, err := Dial("tcp", "localhost:6379")
	defer c.Close()
	assert.Nil(t, err)

	c.Do("SET", "test", "1")

	// Get the span created for Send call
	span := tp.Recorder.GetSpans()[0]
	assert.Equal(t, constants.ClassNames["REDIS"], span.ClassName)
	assert.Equal(t, constants.DomainNames["CACHE"], span.DomainName)
	assert.Equal(t, "WRITE", span.Tags[constants.SpanTags["OPERATION_TYPE"]])
	assert.Equal(t, "localhost", span.Tags[constants.DBTags["DB_INSTANCE"]])
	assert.Equal(t, "SET", span.Tags[constants.DBTags["DB_STATEMENT_TYPE"]])
	assert.Equal(t, "redis", span.Tags[constants.DBTags["DB_TYPE"]])
	assert.Equal(t, "localhost", span.Tags[constants.RedisTags["REDIS_HOST"]])
	assert.Equal(t, "SET", span.Tags[constants.RedisTags["REDIS_COMMAND_TYPE"]])
	assert.Equal(t, constants.AwsLambdaApplicationDomain, span.Tags[constants.SpanTags["TRIGGER_DOMAIN_NAME"]])
	assert.Equal(t, constants.AwsLambdaApplicationClass, span.Tags[constants.SpanTags["TRIGGER_CLASS_NAME"]])
	assert.Equal(t, []string{application.FunctionName}, span.Tags[constants.SpanTags["TRIGGER_OPERATION_NAMES"]])
	assert.Equal(t, true, span.Tags[constants.SpanTags["TOPOLOGY_VERTEX"]])

	assert.Equal(t, "SET test 1", span.Tags[constants.DBTags["DB_STATEMENT"]])
	assert.Equal(t, "SET test 1", span.Tags[constants.RedisTags["REDIS_COMMAND"]])

	tp.Reset()
}

func TestRedigoDialURL(t *testing.T) {
	// Initilize trace plugin to set GlobalTracer of opentracing
	tp := trace.New()

	c, err := DialURL("redis://localhost:6379")
	defer c.Close()
	assert.Nil(t, err)

	c.Do("SET", "test", "1")

	// Get the span created for Send call
	span := tp.Recorder.GetSpans()[0]
	assert.Equal(t, constants.ClassNames["REDIS"], span.ClassName)
	assert.Equal(t, constants.DomainNames["CACHE"], span.DomainName)
	assert.Equal(t, "WRITE", span.Tags[constants.SpanTags["OPERATION_TYPE"]])
	assert.Equal(t, "localhost", span.Tags[constants.DBTags["DB_INSTANCE"]])
	assert.Equal(t, "SET", span.Tags[constants.DBTags["DB_STATEMENT_TYPE"]])
	assert.Equal(t, "redis", span.Tags[constants.DBTags["DB_TYPE"]])
	assert.Equal(t, "localhost", span.Tags[constants.RedisTags["REDIS_HOST"]])
	assert.Equal(t, "SET", span.Tags[constants.RedisTags["REDIS_COMMAND_TYPE"]])
	assert.Equal(t, constants.AwsLambdaApplicationDomain, span.Tags[constants.SpanTags["TRIGGER_DOMAIN_NAME"]])
	assert.Equal(t, constants.AwsLambdaApplicationClass, span.Tags[constants.SpanTags["TRIGGER_CLASS_NAME"]])
	assert.Equal(t, []string{application.FunctionName}, span.Tags[constants.SpanTags["TRIGGER_OPERATION_NAMES"]])
	assert.Equal(t, true, span.Tags[constants.SpanTags["TOPOLOGY_VERTEX"]])

	assert.Equal(t, "SET test 1", span.Tags[constants.DBTags["DB_STATEMENT"]])
	assert.Equal(t, "SET test 1", span.Tags[constants.RedisTags["REDIS_COMMAND"]])

	tp.Reset()
}

func TestRedigoDoCommandError(t *testing.T) {
	// Initilize trace plugin to set GlobalTracer of opentracing
	tp := trace.New()

	c, err := Dial("tcp", "localhost:6379")
	defer c.Close()
	assert.Nil(t, err)

	c.Do("UNKNOWN_COMMAND", "test", "1")

	// Get the span created for Send call
	span := tp.Recorder.GetSpans()[0]
	assert.Equal(t, constants.ClassNames["REDIS"], span.ClassName)
	assert.Equal(t, constants.DomainNames["CACHE"], span.DomainName)
	assert.Equal(t, "", span.Tags[constants.SpanTags["OPERATION_TYPE"]])
	assert.Equal(t, "localhost", span.Tags[constants.DBTags["DB_INSTANCE"]])
	assert.Equal(t, "UNKNOWN_COMMAND", span.Tags[constants.DBTags["DB_STATEMENT_TYPE"]])
	assert.Equal(t, "redis", span.Tags[constants.DBTags["DB_TYPE"]])
	assert.Equal(t, "localhost", span.Tags[constants.RedisTags["REDIS_HOST"]])
	assert.Equal(t, "UNKNOWN_COMMAND", span.Tags[constants.RedisTags["REDIS_COMMAND_TYPE"]])
	assert.Equal(t, constants.AwsLambdaApplicationDomain, span.Tags[constants.SpanTags["TRIGGER_DOMAIN_NAME"]])
	assert.Equal(t, constants.AwsLambdaApplicationClass, span.Tags[constants.SpanTags["TRIGGER_CLASS_NAME"]])
	assert.Equal(t, []string{application.FunctionName}, span.Tags[constants.SpanTags["TRIGGER_OPERATION_NAMES"]])
	assert.Equal(t, true, span.Tags[constants.SpanTags["TOPOLOGY_VERTEX"]])

	assert.Equal(t, "UNKNOWN_COMMAND test 1", span.Tags[constants.DBTags["DB_STATEMENT"]])
	assert.Equal(t, "UNKNOWN_COMMAND test 1", span.Tags[constants.RedisTags["REDIS_COMMAND"]])

	assert.Equal(t, true, span.Tags["error"])

	tp.Reset()
}

func TestRedigoDoParentSpan(t *testing.T) {
	// Initilize trace plugin to set GlobalTracer of opentracing
	tp := trace.New()

	parentSpan, ctxWithSpan := ot.StartSpanFromContext(context.Background(), "set")
	parentSpanRaw, _ := tracer.GetRaw(parentSpan)
	defer parentSpan.Finish()

	c, err := Dial("tcp", "localhost:6379")
	defer c.Close()
	assert.Nil(t, err)

	c.Do("SET", "test", "1", ctxWithSpan)

	// Get the span created for Send call
	span := tp.Recorder.GetSpans()[1]

	assert.Equal(t, parentSpanRaw.Context.SpanID, span.ParentSpanID)
	tp.Reset()
}

func TestRedigoDialConnectionError(t *testing.T) {
	_, err := Dial("tcp", "localhost:6380")
	assert.NotNil(t, err)
}

func TestRedigoDialURLConnectionError(t *testing.T) {
	_, err := DialURL("redis://localhost:6380")
	assert.NotNil(t, err)
}

func TestRedigoDialURLWrongAddress(t *testing.T) {
	_, err := DialURL("test")
	assert.NotNil(t, err)
}
