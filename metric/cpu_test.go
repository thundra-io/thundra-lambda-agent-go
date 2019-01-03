package metric

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrepareCPUMetricsData(t *testing.T) {
	mp := NewBuilder().Build()

	cpuStatsData := prepareCPUMetricsData(mp)

	assert.Equal(t, cpuMetric, cpuStatsData.MetricName)
	assert.Equal(t, mp.data.metricTimestamp, cpuStatsData.MetricTimestamp)
}
