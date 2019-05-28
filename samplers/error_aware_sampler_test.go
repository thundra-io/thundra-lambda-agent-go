package samplers

import (
	"testing"

	ot "github.com/opentracing/opentracing-go"
	"github.com/stretchr/testify/assert"
	"github.com/thundra-io/thundra-lambda-agent-go/tracer"
)

func TestErrorAwareSample(t *testing.T) {

	eas := NewErrorAwareSampler()
	assert.True(t, eas.IsSampled(&tracer.RawSpan{Tags: ot.Tags{"error": true}}))
	assert.False(t, eas.IsSampled(&tracer.RawSpan{}))
}
