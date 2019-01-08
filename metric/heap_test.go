package metric

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrepareHeapMetricsData(t *testing.T) {
	mp := NewBuilder().Build()

	memStats := &runtime.MemStats{}
	base := mp.prepareMetricsData()
	heapMetricsData := prepareHeapMetricsData(mp, memStats, base)

	assert.Equal(t, heapMetric, heapMetricsData.MetricName)

	assert.Equal(t, memStats.HeapAlloc, heapMetricsData.Metrics[heapAlloc])
	assert.Equal(t, memStats.HeapSys, heapMetricsData.Metrics[heapSys])
	assert.Equal(t, memStats.HeapInuse, heapMetricsData.Metrics[heapInuse])
	assert.Equal(t, memStats.HeapObjects, heapMetricsData.Metrics[heapObjects])
}
