package metric

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrepareHeapMetricsData(t *testing.T) {
	metric := NewBuilder().Build()

	memStats := &runtime.MemStats{}

	heapMetricsData := prepareHeapMetricsData(metric, memStats)

	assert.Equal(t, heapMetric, heapMetricsData.MetricName)

	assert.Equal(t, memStats.HeapAlloc, heapMetricsData.Metrics[heapAlloc])
	assert.Equal(t, memStats.HeapSys, heapMetricsData.Metrics[heapSys])
	assert.Equal(t, memStats.HeapInuse, heapMetricsData.Metrics[heapInuse])
	assert.Equal(t, memStats.HeapObjects, heapMetricsData.Metrics[heapObjects])
}
