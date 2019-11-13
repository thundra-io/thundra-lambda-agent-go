package tracer

import (
	"github.com/thundra-io/thundra-lambda-agent-go/constants"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thundra-io/thundra-lambda-agent-go/ext"
)

func TestConfig(t *testing.T) {
	sasl := SecurityAwareSpanListener{
		block: true,
		whitelist: []Operation{
			{
				className: "HTTP",
				tags: map[string]interface{}{
					"http.host":      []string{"www.google.com", "www.google.com"},
					"operation.type": []string{"GET"},
				},
			},
			{
				className: "AWS-DynamoDB",
				tags: map[string]interface{}{
					"aws.dynamodb.table.name": []string{"Users"},
					"operation.type":          []string{"READ"},
				},
			},
		},
		blacklist: []Operation{
			{
				className: "HTTP",
				tags: map[string]interface{}{
					"http.host":      []string{"www.foo.com", "www.foo.com"},
					"operation.type": []string{"POST"},
				},
			},
			{
				className: "AWS-SNS",
				tags: map[string]interface{}{
					"aws.sns.topic.name": []string{"foo-topic"},
					"operation.type":     []string{"WRITE"},
				},
			},
		},
	}

	assert.Equal(t, true, sasl.block)
	assert.Equal(t, 2, len(sasl.whitelist))
	assert.Equal(t, 2, len(sasl.blacklist))

	assert.Equal(t, "HTTP", sasl.whitelist[0].className)
	assert.Equal(t, []string{"www.google.com", "www.google.com"}, sasl.whitelist[0].tags["http.host"])
	assert.Equal(t, []string{"GET"}, sasl.whitelist[0].tags["operation.type"])
	assert.Equal(t, "www.google.com", sasl.whitelist[0].tags["http.host"].([]string)[0])

	assert.Equal(t, "AWS-SNS", sasl.blacklist[1].className)
	assert.Equal(t, []string{"foo-topic"}, sasl.blacklist[1].tags["aws.sns.topic.name"])
	assert.Equal(t, []string{"WRITE"}, sasl.blacklist[1].tags["operation.type"])
	assert.Equal(t, "foo-topic", sasl.blacklist[1].tags["aws.sns.topic.name"].([]string)[0])
}

func TestWhiteListSpan(t *testing.T) {
	sasl := SecurityAwareSpanListener{
		block: true,
		whitelist: []Operation{
			{
				className: "HTTP",
				tags: map[string]interface{}{
					"http.host":      []string{"www.google.com", "www.facebook.com"},
					"operation.type": []string{"GET"},
				},
			},
			{
				className: "AWS-DynamoDB",
				tags: map[string]interface{}{
					"aws.dynamodb.table.name": []string{"Users"},
					"operation.type":          []string{"READ"},
				},
			},
		},
	}

	tracer, _ := newTracerAndRecorder()

	s1 := tracer.StartSpan("foo", ext.ClassName("HTTP"), ext.OperationType("GET"))
	s1.SetTag("http.host", "www.google.com")
	s1.SetTag("topology.vertex", true)

	s2 := tracer.StartSpan("foo", ext.ClassName("HTTP"), ext.OperationType("GET"))
	s2.SetTag("http.host", "www.google.com")
	s2.SetTag("topology.vertex", true)

	sasl.OnSpanStarted(s1.(*spanImpl))
	sasl.OnSpanStarted(s2.(*spanImpl))

	assert.Equal(t, nil, s1.(*spanImpl).raw.GetTag(constants.SecurityTags["BLOCKED"]))
	assert.Equal(t, nil, s1.(*spanImpl).raw.GetTag(constants.SecurityTags["VIOLATED"]))
	assert.Equal(t, nil, s2.(*spanImpl).raw.GetTag(constants.SecurityTags["BLOCKED"]))
	assert.Equal(t, nil, s2.(*spanImpl).raw.GetTag(constants.SecurityTags["VIOLATED"]))

	s3 := tracer.StartSpan("foo", ext.ClassName("HTTP"), ext.OperationType("POST"))
	s3.SetTag("http.host", "www.example.com")
	s3.SetTag("topology.vertex", true)

	var errorPanicked3 error
	func() {
		defer func() {
			errorPanicked3 = recover().(error)
			assert.Equal(t, "Operation was blocked due to security configuration", errorPanicked3.Error())
			assert.Equal(t, true, s3.(*spanImpl).raw.GetTag(constants.SecurityTags["BLOCKED"]))
			assert.Equal(t, true, s3.(*spanImpl).raw.GetTag(constants.SecurityTags["VIOLATED"]))
		}()
		sasl.OnSpanStarted(s3.(*spanImpl))
	}()

	s4 := tracer.StartSpan("foo", ext.ClassName("AWS-DynamoDB"), ext.OperationType("WRITE"))
	s4.SetTag("topology.vertex", true)
	s4.SetTag(constants.SpanTags["OPERATION_TYPE"], "WRITE")

	var errorPanicked4 error
	func() {
		defer func() {
			errorPanicked4 = recover().(error)
			assert.Equal(t, "Operation was blocked due to security configuration", errorPanicked4.Error())
			assert.Equal(t, true, s4.(*spanImpl).raw.GetTag(constants.SecurityTags["BLOCKED"]))
			assert.Equal(t, true, s4.(*spanImpl).raw.GetTag(constants.SecurityTags["VIOLATED"]))
		}()
		sasl.OnSpanStarted(s4.(*spanImpl))
	}()

	s5 := tracer.StartSpan("foo", ext.ClassName("HTTP"), ext.OperationType("POST"))
	s5.SetTag("http.host", "www.google.com")
	s5.SetTag("topology.vertex", true)

	var errorPanicked5 error
	func() {
		defer func() {
			errorPanicked5 = recover().(error)
			assert.Equal(t, "Operation was blocked due to security configuration", errorPanicked5.Error())
			assert.Equal(t, true, s5.(*spanImpl).raw.GetTag(constants.SecurityTags["BLOCKED"]))
			assert.Equal(t, true, s5.(*spanImpl).raw.GetTag(constants.SecurityTags["VIOLATED"]))
		}()
		sasl.OnSpanStarted(s5.(*spanImpl))
	}()
}
