package tracer

import (
	"testing"
	"time"

	"github.com/thundra-io/thundra-lambda-agent-go/v2/constants"
	"github.com/thundra-io/thundra-lambda-agent-go/v2/ext"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/stretchr/testify/assert"
)

const (
	duration      = 10
	operationName = "creating-bubble"
	className     = "Test Class"
	domainName    = "Test Domain"
)

func TestStartSpan(t *testing.T) {
	tracer, r := newTracerAndRecorder()

	f := func() {
		span := tracer.StartSpan(operationName)
		defer span.Finish()
		time.Sleep(time.Millisecond * duration)
	}

	f()

	spans := r.GetSpans()
	span := spans[0]

	assert.True(t, len(spans) == 1)
	assert.True(t, span.Duration() >= int64(duration))
	assert.True(t, span.OperationName == operationName)
	assert.True(t, span.ClassName == constants.DefaultClassName)
	assert.True(t, span.DomainName == constants.DefaultDomainName)
}

func TestStartSpanWithOptions(t *testing.T) {
	tracer, r := newTracerAndRecorder()

	f := func() {
		span := tracer.StartSpan(
			operationName,
			ext.ClassName(className),
			ext.DomainName(domainName),
			opentracing.Tag{Key: "stage", Value: "testing"},
		)
		defer span.Finish()

		time.Sleep(time.Millisecond * duration)
	}

	f()
	spans := r.GetSpans()
	span := spans[0]

	assert.True(t, len(spans) == 1)
	assert.True(t, span.Duration() >= int64(duration))
	assert.True(t, span.OperationName == operationName)
	assert.True(t, span.ClassName == className)
	assert.True(t, span.DomainName == domainName)
	assert.True(t, len(span.GetTags()) == 1)
	assert.True(t, span.GetTags()["stage"] == "testing")
}

func TestParentChildRelation(t *testing.T) {
	tracer, r := newTracerAndRecorder()

	f := func() {
		parentSpan := tracer.StartSpan("parentSpan")
		time.Sleep(time.Millisecond * duration)
		childSpan := tracer.StartSpan("childSpan", opentracing.ChildOf(parentSpan.Context()))
		time.Sleep(time.Millisecond * duration)
		childSpan.Finish()
		time.Sleep(time.Millisecond * duration)
		parentSpan.Finish()
	}

	f()

	spans := r.GetSpans()
	parentSpan, childSpan := spans[0], spans[1]

	assert.True(t, len(spans) == 2)
	assert.True(t, parentSpan.OperationName == "parentSpan")
	assert.True(t, childSpan.OperationName == "childSpan")
	assert.True(t, parentSpan.ParentSpanID == "")
	assert.True(t, childSpan.ParentSpanID == parentSpan.Context.SpanID)
	assert.True(t, childSpan.Context.TraceID == parentSpan.Context.TraceID)
	assert.True(t, childSpan.Duration() >= duration)
	assert.True(t, parentSpan.Duration() >= 3*duration)
}

func newTracerAndRecorder() (opentracing.Tracer, *InMemorySpanRecorder) {
	r := NewInMemoryRecorder()
	tracer := New(r)

	return tracer, r
}
