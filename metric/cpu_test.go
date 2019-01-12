package metric

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrepareCPUMetricsData(t *testing.T) {
	mp := NewBuilder().Build()
	base := mp.prepareMetricsData()
	cpuStatsData := prepareCPUMetricsData(mp, base)

	assert.True(t, len(cpuStatsData.ID) != 0)
	assert.Equal(t, cpuMetric, cpuStatsData.MetricName)
	assert.Equal(t, mp.data.metricTimestamp, cpuStatsData.MetricTimestamp)
}
