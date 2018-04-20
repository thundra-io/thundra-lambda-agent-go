package metric

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/thundra-io/thundra-lambda-agent-go/plugin"
)

func TestNewBuilder(t *testing.T) {
	prepareEnvironment()
	m := NewBuilder().
		EnableDiskStats().
		EnableNetStats().
		EnableCPUStats().
		EnableGoroutineStats().
		EnableHeapStats().
		EnableGCStats().
		Build()

	assert.Equal(t, functionName, m.applicationName)
	assert.Equal(t, appId, m.applicationId)
	assert.Equal(t, functionVersion, m.applicationVersion)
	assert.Equal(t, applicationProfile, m.applicationProfile)
	assert.Equal(t, plugin.ApplicationType, m.applicationType)

	assert.NotNil(t, m.prevDiskStat)
	assert.NotNil(t, m.prevNetStat)
	assert.NotNil(t, m.process)

	assert.True(t, m.enableDiskStats)
	assert.True(t, m.enableNetStats)
	assert.True(t, m.enableCPUStats)
	assert.True(t, m.enableGoroutineStats)
	assert.True(t, m.enableHeapStats)
	assert.True(t, m.enableGCStats)
}
