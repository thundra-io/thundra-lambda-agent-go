package thundragoredis

import (
	"context"
	"testing"

	"github.com/go-redis/redis"
	ot "github.com/opentracing/opentracing-go"
	"github.com/stretchr/testify/assert"
	"github.com/thundra-io/thundra-lambda-agent-go/v2/application"
	"github.com/thundra-io/thundra-lambda-agent-go/v2/constants"
	"github.com/thundra-io/thundra-lambda-agent-go/v2/trace"
	"github.com/thundra-io/thundra-lambda-agent-go/v2/tracer"
)

func TestWithWrongURL(t *testing.T) {
	_ = trace.New()
	// Initilize trace plugin to set GlobalTracer of opentracing
	c := NewClient(&redis.Options{
		Addr: "localhost:1234",
	})
	defer c.Close()

	err := c.Ping().Err()
	assert.Error(t, err)
}

func TestPing(t *testing.T) {
	// Initilize trace plugin to set GlobalTracer of opentracing
	tp := trace.New()

	c := NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	defer c.Close()

	err := c.Ping().Err()
	assert.NoError(t, err)

	span := tp.Recorder.GetSpans()[0]
	assert.Equal(t, constants.ClassNames["REDIS"], span.ClassName)
	assert.Equal(t, constants.DomainNames["CACHE"], span.DomainName)
	assert.Equal(t, "", span.Tags[constants.SpanTags["OPERATION_TYPE"]])
	assert.Equal(t, "localhost", span.Tags[constants.DBTags["DB_INSTANCE"]])
	assert.Equal(t, "PING", span.Tags[constants.DBTags["DB_STATEMENT_TYPE"]])
	assert.Equal(t, "redis", span.Tags[constants.DBTags["DB_TYPE"]])
	assert.Equal(t, "localhost", span.Tags[constants.RedisTags["REDIS_HOST"]])
	assert.Equal(t, "PING", span.Tags[constants.RedisTags["REDIS_COMMAND_TYPE"]])
	assert.Equal(t, constants.AwsLambdaApplicationDomain, span.Tags[constants.SpanTags["TRIGGER_DOMAIN_NAME"]])
	assert.Equal(t, constants.AwsLambdaApplicationClass, span.Tags[constants.SpanTags["TRIGGER_CLASS_NAME"]])
	assert.Equal(t, []string{application.FunctionName}, span.Tags[constants.SpanTags["TRIGGER_OPERATION_NAMES"]])
	assert.Equal(t, true, span.Tags[constants.SpanTags["TOPOLOGY_VERTEX"]])

	assert.Equal(t, "ping", span.Tags[constants.DBTags["DB_STATEMENT"]])
	assert.Equal(t, "ping", span.Tags[constants.RedisTags["REDIS_COMMAND"]])

	tp.Reset()
}

func TestSet(t *testing.T) {
	// Initilize trace plugin to set GlobalTracer of opentracing
	tp := trace.New()

	c := NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	defer c.Close()

	err := c.Set("foo", "bar", 0).Err()
	assert.NoError(t, err)

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

	assert.Equal(t, "set foo bar", span.Tags[constants.DBTags["DB_STATEMENT"]])
	assert.Equal(t, "set foo bar", span.Tags[constants.RedisTags["REDIS_COMMAND"]])

	tp.Reset()
}

func TestSetWithParent(t *testing.T) {
	// Initilize trace plugin to set GlobalTracer of opentracing
	tp := trace.New()

	ps, ctx := ot.StartSpanFromContext(context.Background(), "parent")
	psr, _ := tracer.GetRaw(ps)
	defer ps.Finish()

	c := NewClient(&redis.Options{
		Addr: "localhost:6379",
	}).WithContext(ctx)
	defer c.Close()

	err := c.Set("foo", "bar", 0).Err()
	assert.NoError(t, err)

	span := tp.Recorder.GetSpans()[1]

	assert.Equal(t, psr.Context.SpanID, span.ParentSpanID)

	tp.Reset()
}

func TestUnknownCommand(t *testing.T) {
	// Initilize trace plugin to set GlobalTracer of opentracing
	tp := trace.New()

	c := NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	defer c.Close()

	err := c.Do("UNKNOWN", "test", "1").Err()
	assert.Error(t, err)

	// Get the span created for Send call
	span := tp.Recorder.GetSpans()[0]
	assert.Equal(t, constants.ClassNames["REDIS"], span.ClassName)
	assert.Equal(t, constants.DomainNames["CACHE"], span.DomainName)
	assert.Equal(t, "", span.Tags[constants.SpanTags["OPERATION_TYPE"]])
	assert.Equal(t, "localhost", span.Tags[constants.DBTags["DB_INSTANCE"]])
	assert.Equal(t, "UNKNOWN", span.Tags[constants.DBTags["DB_STATEMENT_TYPE"]])
	assert.Equal(t, "redis", span.Tags[constants.DBTags["DB_TYPE"]])
	assert.Equal(t, "localhost", span.Tags[constants.RedisTags["REDIS_HOST"]])
	assert.Equal(t, "UNKNOWN", span.Tags[constants.RedisTags["REDIS_COMMAND_TYPE"]])
	assert.Equal(t, constants.AwsLambdaApplicationDomain, span.Tags[constants.SpanTags["TRIGGER_DOMAIN_NAME"]])
	assert.Equal(t, constants.AwsLambdaApplicationClass, span.Tags[constants.SpanTags["TRIGGER_CLASS_NAME"]])
	assert.Equal(t, []string{application.FunctionName}, span.Tags[constants.SpanTags["TRIGGER_OPERATION_NAMES"]])
	assert.Equal(t, true, span.Tags[constants.SpanTags["TOPOLOGY_VERTEX"]])

	assert.Equal(t, "UNKNOWN", span.Tags[constants.DBTags["DB_STATEMENT"]])
	assert.Equal(t, "UNKNOWN", span.Tags[constants.RedisTags["REDIS_COMMAND"]])

	tp.Reset()
}

func TestPipeline(t *testing.T) {
	// Initilize trace plugin to set GlobalTracer of opentracing
	tp := trace.New()

	c := NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	pipe := c.Pipeline()
	defer c.Close()

	pipe.Set("foo", "bar", 0)
	pipe.Get("foo")
	pipe.Incr("pipeline_counter")

	_, err := pipe.Exec()

	assert.NoError(t, err)

	span := tp.Recorder.GetSpans()[0]
	assert.Equal(t, constants.ClassNames["REDIS"], span.ClassName)
	assert.Equal(t, constants.DomainNames["CACHE"], span.DomainName)
	assert.Equal(t, "", span.Tags[constants.SpanTags["OPERATION_TYPE"]])
	assert.Equal(t, "localhost", span.Tags[constants.DBTags["DB_INSTANCE"]])
	assert.Equal(t, "PIPELINE", span.Tags[constants.DBTags["DB_STATEMENT_TYPE"]])
	assert.Equal(t, "redis", span.Tags[constants.DBTags["DB_TYPE"]])
	assert.Equal(t, "localhost", span.Tags[constants.RedisTags["REDIS_HOST"]])
	assert.Equal(t, "PIPELINE", span.Tags[constants.RedisTags["REDIS_COMMAND_TYPE"]])
	assert.Equal(t, constants.AwsLambdaApplicationDomain, span.Tags[constants.SpanTags["TRIGGER_DOMAIN_NAME"]])
	assert.Equal(t, constants.AwsLambdaApplicationClass, span.Tags[constants.SpanTags["TRIGGER_CLASS_NAME"]])
	assert.Equal(t, []string{application.FunctionName}, span.Tags[constants.SpanTags["TRIGGER_OPERATION_NAMES"]])
	assert.Equal(t, true, span.Tags[constants.SpanTags["TOPOLOGY_VERTEX"]])

	assert.Equal(t, "set foo bar\nget foo\nincr pipeline_counter", span.Tags[constants.DBTags["DB_STATEMENT"]])
	assert.Equal(t, "set foo bar\nget foo\nincr pipeline_counter", span.Tags[constants.RedisTags["REDIS_COMMAND"]])

	tp.Reset()
}
