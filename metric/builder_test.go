package metric

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thundra-io/thundra-lambda-agent-go/test"
)

func TestNewBuilder(t *testing.T) {
	test.PrepareEnvironment()
	m := NewBuilder().Build()

	assert.NotNil(t, m.prevDiskStat)
	assert.NotNil(t, m.prevNetStat)
	assert.NotNil(t, m.process)

	assert.False(t, m.disableDiskStats)
	assert.False(t, m.disableNetStats)
	assert.False(t, m.disableCPUStats)
	assert.False(t, m.disableGoroutineStats)
	assert.False(t, m.disableHeapStats)
	assert.False(t, m.disableGCStats)
	test.CleanEnvironment()
}
