package metric

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/thundra-io/thundra-lambda-agent-go/plugin"
)

func TestNewBuilder(t *testing.T) {
	prepareEnvironment()
	m := NewBuilder().Build()

	assert.Equal(t, functionName, m.applicationName)
	assert.Equal(t, appId, m.applicationId)
	assert.Equal(t, functionVersion, m.applicationVersion)
	assert.Equal(t, applicationProfile, m.applicationProfile)
	assert.Equal(t, plugin.ApplicationType, m.applicationType)

	assert.NotNil(t, m.prevDiskStat)
	assert.NotNil(t, m.prevNetStat)
	assert.NotNil(t, m.process)

	assert.False(t, m.disableDiskStats)
	assert.False(t, m.disableNetStats)
	assert.False(t, m.disableCPUStats)
	assert.False(t, m.disableGoroutineStats)
	assert.False(t, m.disableHeapStats)
	assert.False(t, m.disableGCStats)
	cleanEnvironment()
}
