package invocation

import (
	"testing"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/stretchr/testify/assert"
	"github.com/thundra-io/thundra-lambda-agent-go/constants"
	"github.com/thundra-io/thundra-lambda-agent-go/ext"
	"github.com/thundra-io/thundra-lambda-agent-go/trace"
)

func createMockSpans() {
	spans := []opentracing.Span{
		opentracing.StartSpan(
			"www.test.com",
			ext.ClassName("HTTP"),
			opentracing.Tag{Key: constants.SpanTags["OPERATION_TYPE"], Value: "GET"},
			opentracing.Tag{Key: constants.SpanTags["TOPOLOGY_VERTEX"], Value: true},
		), opentracing.StartSpan(
			"www.test.com",
			ext.ClassName("HTTP"),
			opentracing.Tag{Key: constants.SpanTags["OPERATION_TYPE"], Value: "GET"},
			opentracing.Tag{Key: constants.SpanTags["TOPOLOGY_VERTEX"], Value: true},
		),
		opentracing.StartSpan(
			"localhost",
			ext.ClassName("Redis"),
			opentracing.Tag{Key: constants.SpanTags["OPERATION_TYPE"], Value: "READ"},
			opentracing.Tag{Key: constants.SpanTags["TOPOLOGY_VERTEX"], Value: true},
			opentracing.Tag{Key: constants.AwsError, Value: true},
			opentracing.Tag{Key: constants.AwsErrorKind, Value: "testErr"},
		),
	}
	defer func() {
		for _, s := range spans {
			s.Finish()
		}
	}()
}

func TestGetResources(t *testing.T) {
	tp := trace.GetInstance()
	tp.Reset()

	createMockSpans()
	resources := getResources("")

	var resource1, resource2 Resource
	if resources[0].ResourceType == "HTTP" {
		resource1 = resources[0]
		resource2 = resources[1]
	} else {
		resource2 = resources[0]
		resource1 = resources[1]
	}

	assert.Equal(t, "HTTP", resource1.ResourceType)
	assert.Equal(t, "www.test.com", resource1.ResourceName)
	assert.Equal(t, "GET", resource1.ResourceOperation)
	assert.Equal(t, 2, resource1.ResourceCount)
	assert.Equal(t, 0, resource1.ResourceErrorCount)

	assert.Equal(t, "Redis", resource2.ResourceType)
	assert.Equal(t, "localhost", resource2.ResourceName)
	assert.Equal(t, "READ", resource2.ResourceOperation)
	assert.Equal(t, 1, resource2.ResourceCount)
	assert.Equal(t, 1, resource2.ResourceErrorCount)
	assert.ElementsMatch(t, []string{"testErr"}, resource2.ResourceErrors)

	tp.Reset()
}
