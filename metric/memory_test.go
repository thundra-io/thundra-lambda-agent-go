package metric

import (
	"github.com/thundra-io/thundra-lambda-agent-go/application"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrepareMemoryMetricsData(t *testing.T) {
	mp := New()

	base := mp.prepareMetricsData()
	memoryMetricsData := prepareMemoryMetricsData(mp, base)

	assert.True(t, len(memoryMetricsData.ID) != 0)
	assert.Equal(t, memoryMetric, memoryMetricsData.MetricName)
	assert.Equal(t, application.MemoryUsed, int(memoryMetricsData.Metrics[appUsedMemory].(uint64) / miBToB))
}
