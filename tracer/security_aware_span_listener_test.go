package tracer

import (
	"testing"

	"github.com/thundra-io/thundra-lambda-agent-go/v2/constants"

	"github.com/stretchr/testify/assert"
	"github.com/thundra-io/thundra-lambda-agent-go/v2/ext"
)

func TestConfig(t *testing.T) {
	sasl := SecurityAwareSpanListener{
		block: true,
		whitelist: &[]Operation{
			{
				ClassName: "HTTP",
				Tags: map[string][]string{
					"http.host":      {"www.google.com", "www.google.com"},
					"operation.type": {"GET"},
				},
			},
			{
				ClassName: "AWS-DynamoDB",
				Tags: map[string][]string{
					"aws.dynamodb.table.name": {"Users"},
					"operation.type":          {"READ"},
				},
			},
		},
		blacklist: &[]Operation{
			{
				ClassName: "HTTP",
				Tags: map[string][]string{
					"http.host":      {"www.foo.com", "www.foo.com"},
					"operation.type": {"POST"},
				},
			},
			{
				ClassName: "AWS-SNS",
				Tags: map[string][]string{
					"aws.sns.topic.name": {"foo-topic"},
					"operation.type":     {"WRITE"},
				},
			},
		},
	}

	assert.Equal(t, true, sasl.block)
	assert.Equal(t, 2, len(*sasl.whitelist))
	assert.Equal(t, 2, len(*sasl.blacklist))

	assert.Equal(t, "HTTP", (*sasl.whitelist)[0].ClassName)
	assert.Equal(t, []string{"www.google.com", "www.google.com"}, (*sasl.whitelist)[0].Tags["http.host"])
	assert.Equal(t, []string{"GET"}, (*sasl.whitelist)[0].Tags["operation.type"])

	assert.Equal(t, "AWS-SNS", (*sasl.blacklist)[1].ClassName)
	assert.Equal(t, []string{"foo-topic"}, (*sasl.blacklist)[1].Tags["aws.sns.topic.name"])
	assert.Equal(t, []string{"WRITE"}, (*sasl.blacklist)[1].Tags["operation.type"])
}

func TestWhiteListSpan(t *testing.T) {
	sasl := SecurityAwareSpanListener{
		block: true,
		whitelist: &[]Operation{
			{
				ClassName: "HTTP",
				Tags: map[string][]string{
					"http.host":      []string{"www.google.com", "www.facebook.com"},
					"operation.type": []string{"GET"},
				},
			},
			{
				ClassName: "AWS-DynamoDB",
				Tags: map[string][]string{
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
	s5.SetTag("topology.vertex", true)
	s5.SetTag("http.host", "www.google.com")
	s4.SetTag(constants.SpanTags["OPERATION_TYPE"], "POST")

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

func TestBlackListSpan(t *testing.T) {
	sasl := SecurityAwareSpanListener{
		block: true,
		blacklist: &[]Operation{
			{
				ClassName: "HTTP",
				Tags: map[string][]string{
					"http.host":      []string{"host1.com", "host2.com"},
					"operation.type": []string{"GET"},
				},
			},
			{
				ClassName: "AWS-DynamoDB",
				Tags: map[string][]string{
					"aws.dynamodb.table.name": []string{"users"},
					"operation.type":          []string{"READ"},
				},
			},
			{
				ClassName: "ELASTICSEARCH",
				Tags: map[string][]string{
					"elasticsearch.normalized_uri": {
						"/twitter",
					},
					"operation.type": {
						"POST",
					},
				},
			},
		},
	}

	tracer, _ := newTracerAndRecorder()

	s1 := tracer.StartSpan("foo", ext.ClassName("HTTP"), ext.OperationType("POST"))
	s1.SetTag("http.host", "host1.com")
	s1.SetTag("topology.vertex", true)
	s1.SetTag(constants.SpanTags["OPERATION_TYPE"], "POST")

	s2 := tracer.StartSpan("bar", ext.ClassName("HTTP"), ext.OperationType("POST"))
	s2.SetTag("http.host", "host2.com")
	s2.SetTag("topology.vertex", true)
	s2.SetTag(constants.SpanTags["OPERATION_TYPE"], "POST")

	sasl.OnSpanStarted(s1.(*spanImpl))
	sasl.OnSpanStarted(s2.(*spanImpl))

	assert.Equal(t, nil, s1.(*spanImpl).raw.GetTag(constants.SecurityTags["BLOCKED"]))
	assert.Equal(t, nil, s1.(*spanImpl).raw.GetTag(constants.SecurityTags["VIOLATED"]))
	assert.Equal(t, nil, s2.(*spanImpl).raw.GetTag(constants.SecurityTags["BLOCKED"]))
	assert.Equal(t, nil, s2.(*spanImpl).raw.GetTag(constants.SecurityTags["VIOLATED"]))

	s3 := tracer.StartSpan("foo", ext.ClassName("HTTP"), ext.OperationType("GET"))
	s3.SetTag("http.host", "host1.com")
	s3.SetTag("topology.vertex", true)
	s3.SetTag(constants.SpanTags["OPERATION_TYPE"], "GET")

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

	s4 := tracer.StartSpan("foo", ext.ClassName("AWS-DynamoDB"), ext.OperationType("READ"))
	s4.SetTag("topology.vertex", true)
	s4.SetTag("aws.dynamodb.table.name", "users")
	s4.SetTag(constants.SpanTags["OPERATION_TYPE"], "READ")

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

	s5 := tracer.StartSpan("foo", ext.ClassName("AWS-DynamoDB"), ext.OperationType("WRITE"))
	s5.SetTag("topology.vertex", true)
	s5.SetTag("aws.dynamodb.table.name", "users")
	s5.SetTag(constants.SpanTags["OPERATION_TYPE"], "WRITE")

	assert.Equal(t, nil, s1.(*spanImpl).raw.GetTag(constants.SecurityTags["BLOCKED"]))
	assert.Equal(t, nil, s1.(*spanImpl).raw.GetTag(constants.SecurityTags["VIOLATED"]))
}

func TestViolateBlacklistSpan(t *testing.T) {
	sasl := SecurityAwareSpanListener{
		block: false,
		blacklist: &[]Operation{
			{
				ClassName: "HTTP",
				Tags: map[string][]string{
					"http.host":      {"host1.com", "host2.com"},
					"operation.type": {"GET"},
				},
			},
			{
				ClassName: "AWS-DynamoDB",
				Tags: map[string][]string{
					"aws.dynamodb.table.name": {"users"},
					"operation.type":          {"READ"},
				},
			},
		},
	}

	tracer, _ := newTracerAndRecorder()

	s1 := tracer.StartSpan("foo", ext.ClassName("HTTP"), ext.OperationType("GET"))
	s1.SetTag("http.host", "host1.com")
	s1.SetTag("topology.vertex", true)
	s1.SetTag(constants.SpanTags["OPERATION_TYPE"], "GET")

	sasl.OnSpanStarted(s1.(*spanImpl))

	assert.Equal(t, nil, s1.(*spanImpl).raw.GetTag(constants.SecurityTags["BLOCKED"]))
	assert.Equal(t, true, s1.(*spanImpl).raw.GetTag(constants.SecurityTags["VIOLATED"]))

	s2 := tracer.StartSpan("foo", ext.ClassName("AWS-DynamoDB"), ext.OperationType("READ"))
	s2.SetTag("topology.vertex", true)
	s2.SetTag("aws.dynamodb.table.name", "users")
	s2.SetTag(constants.SpanTags["OPERATION_TYPE"], "READ")

	sasl.OnSpanStarted(s2.(*spanImpl))

	assert.Equal(t, nil, s2.(*spanImpl).raw.GetTag(constants.SecurityTags["BLOCKED"]))
	assert.Equal(t, true, s2.(*spanImpl).raw.GetTag(constants.SecurityTags["VIOLATED"]))

}

func TestEsBlackList(t *testing.T) {
	sasl := SecurityAwareSpanListener{
		block: true,
		blacklist: &[]Operation{
			{
				ClassName: "ELASTICSEARCH",
				Tags: map[string][]string{
					"elasticsearch.normalized_uri": {
						"/twitter",
					},
					"operation.type": {
						"POST",
					},
				},
			},
		},
	}

	tracer, _ := newTracerAndRecorder()

	s1 := tracer.StartSpan("foo", ext.ClassName("ELASTICSEARCH"), ext.OperationType("POST"))
	s1.SetTag("elasticsearch.normalized_uri", "/twitter")
	s1.SetTag("topology.vertex", true)
	s1.SetTag(constants.SpanTags["OPERATION_TYPE"], "POST")

	assert.Equal(t, nil, s1.(*spanImpl).raw.GetTag(constants.SecurityTags["BLOCKED"]))
	assert.Equal(t, nil, s1.(*spanImpl).raw.GetTag(constants.SecurityTags["VIOLATED"]))

	var errorPanicked error
	func() {
		defer func() {
			errorPanicked = recover().(error)
			assert.Equal(t, "Operation was blocked due to security configuration", errorPanicked.Error())
			assert.Equal(t, true, s1.(*spanImpl).raw.GetTag(constants.SecurityTags["BLOCKED"]))
			assert.Equal(t, true, s1.(*spanImpl).raw.GetTag(constants.SecurityTags["VIOLATED"]))
		}()
		sasl.OnSpanStarted(s1.(*spanImpl))
	}()

}

func TestAPIGWBlackList(t *testing.T) {
	sasl := SecurityAwareSpanListener{
		block: true,
		blacklist: &[]Operation{
			{
				ClassName: "HTTP",
				Tags: map[string][]string{
					"http.host": {
						"34zsqapxkj.execute-api.eu-west-1.amazonaws.com",
					},
					"operation.type": {
						"GET",
					},
				},
			},
		},
	}

	tracer, _ := newTracerAndRecorder()

	s1 := tracer.StartSpan("foo", ext.ClassName("HTTP"), ext.OperationType("GET"))
	s1.SetTag("http.host", "34zsqapxkj.execute-api.eu-west-1.amazonaws.com")
	s1.SetTag("topology.vertex", true)
	s1.SetTag(constants.SpanTags["OPERATION_TYPE"], "GET")

	assert.Equal(t, nil, s1.(*spanImpl).raw.GetTag(constants.SecurityTags["BLOCKED"]))
	assert.Equal(t, nil, s1.(*spanImpl).raw.GetTag(constants.SecurityTags["VIOLATED"]))

	var errorPanicked error
	func() {
		defer func() {
			errorPanicked = recover().(error)
			assert.Equal(t, "Operation was blocked due to security configuration", errorPanicked.Error())
			assert.Equal(t, true, s1.(*spanImpl).raw.GetTag(constants.SecurityTags["BLOCKED"]))
			assert.Equal(t, true, s1.(*spanImpl).raw.GetTag(constants.SecurityTags["VIOLATED"]))
		}()
		sasl.OnSpanStarted(s1.(*spanImpl))
	}()

}

func TestEmptyWhileList(t *testing.T) {
	sasl := SecurityAwareSpanListener{
		block:     true,
		whitelist: &[]Operation{},
	}

	tracer, _ := newTracerAndRecorder()

	s1 := tracer.StartSpan("foo", ext.ClassName("HTTP"), ext.OperationType("GET"))
	s1.SetTag("http.host", "34zsqapxkj.execute-api.eu-west-1.amazonaws.com")
	s1.SetTag("topology.vertex", true)
	s1.SetTag(constants.SpanTags["OPERATION_TYPE"], "GET")

	var errorPanicked error
	func() {
		defer func() {
			errorPanicked = recover().(error)
			assert.Equal(t, "Operation was blocked due to security configuration", errorPanicked.Error())
			assert.Equal(t, true, s1.(*spanImpl).raw.GetTag(constants.SecurityTags["BLOCKED"]))
			assert.Equal(t, true, s1.(*spanImpl).raw.GetTag(constants.SecurityTags["VIOLATED"]))
		}()
		sasl.OnSpanStarted(s1.(*spanImpl))
	}()

}
