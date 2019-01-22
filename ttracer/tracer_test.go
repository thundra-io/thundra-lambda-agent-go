package ttracer

import (
	"testing"
	"time"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/stretchr/testify/assert"
	"github.com/thundra-io/thundra-lambda-agent-go/plugin"
)

const (
	duration      = 500
	operationName = "creating-bubble"
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
