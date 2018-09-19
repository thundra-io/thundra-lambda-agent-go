package metric

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrepareCPUStatsData(t *testing.T) {
	metric := NewBuilder().Build()

	cpuStatsData := prepareCpuMetricsData(metric)

	assert.Equal(t, cpuMetric, cpuStatsData.MetricName)
	assert.Equal(t, metric.span.metricTimestamp, cpuStatsData.MetricTimestamp)
}
