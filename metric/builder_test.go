package metric

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thundra-io/thundra-lambda-agent-go/test"
)

func TestNewBuilder(t *testing.T) {
	test.PrepareEnvironment()
	mp := NewBuilder().Build()

	assert.NotNil(t, proc)
	assert.NotNil(t, mp.data)

	assert.False(t, mp.disableDiskMetrics)
	assert.False(t, mp.disableNetMetrics)
	assert.False(t, mp.disableCPUMetrics)
	assert.False(t, mp.disableGoroutineMetrics)
	assert.False(t, mp.disableHeapMetrics)
	assert.False(t, mp.disableGCMetrics)
	test.CleanEnvironment()
}
