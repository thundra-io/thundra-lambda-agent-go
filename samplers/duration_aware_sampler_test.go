package samplers

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/thundra-io/thundra-lambda-agent-go/tracer"
)

func TestDurationAwareSample(t *testing.T) {

	das := NewDurationAwareSampler(100)
	assert.True(t, das.IsSampled(&tracer.RawSpan{StartTimestamp: 0, EndTimestamp: 10}))
	assert.False(t, das.IsSampled(nil))
}

func TestDurationAwareSampleLongerThan(t *testing.T) {

	das := NewDurationAwareSampler(100, true)
	assert.False(t, das.IsSampled(&tracer.RawSpan{StartTimestamp: 0, EndTimestamp: 10}))
}
