package metric

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thundra-io/thundra-lambda-agent-go/test"
)

func TestNewBuilder(t *testing.T) {
	test.PrepareEnvironment()
	m := NewBuilder().Build()

	assert.NotNil(t, m.span)
	assert.NotNil(t, proc)

	assert.False(t, m.disableDiskMetrics)
	assert.False(t, m.disableNetMetrics)
	assert.False(t, m.disableCPUMetrics)
	assert.False(t, m.disableGoroutineMetrics)
	assert.False(t, m.disableHeapMetrics)
	assert.False(t, m.disableGCMetrics)
	test.CleanEnvironment()
}
