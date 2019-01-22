package ttracer

import (
	"testing"
	"time"

	"github.com/thundra-io/thundra-lambda-agent-go/ext"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/stretchr/testify/assert"
	"github.com/thundra-io/thundra-lambda-agent-go/plugin"
)

const (
	duration      = 500
	operationName = "creating-bubble"
	className     = "Test Class"
	domainName    = "Test Domain"
)

func TestStartSpan(t *testing.T) {
	r := NewInMemoryRecorder()
	tracer := New(r)

	f := func() opentracing.Span {
		span := tracer.StartSpan(operationName)
		defer span.Finish()
		time.Sleep(time.Millisecond * duration)
		return span
	}

	span, ok := f().(*spanImpl)
	assert.True(t, ok)
	assert.True(t, span.raw.Duration() >= int64(duration))
	assert.True(t, span.OperationName() == operationName)
	assert.True(t, span.raw.ClassName == plugin.DefaultClassName)
	assert.True(t, span.raw.DomainName == plugin.DefaultDomainName)
}

func TestStartSpanWithOptions(t *testing.T) {
	r := NewInMemoryRecorder()
	tracer := New(r)

	f := func() opentracing.Span {
		span := tracer.StartSpan(
			operationName,
			ext.ClassName(className),
			ext.DomainName(domainName),
			opentracing.Tag{Key: "stage", Value: "testing"},
		)
		defer span.Finish()

		time.Sleep(time.Millisecond * duration)
		return span
	}

	span, ok := f().(*spanImpl)
	assert.True(t, ok)
	assert.True(t, span.raw.Duration() >= int64(duration))
	assert.True(t, span.OperationName() == operationName)
	assert.True(t, span.raw.ClassName == className)
	assert.True(t, span.raw.DomainName == domainName)
	assert.True(t, len(span.raw.GetTags()) == 1)
	assert.True(t, span.raw.GetTags()["stage"].(string) == "testing")
}
