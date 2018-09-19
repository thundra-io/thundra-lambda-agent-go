package metric

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrepareHeapStatsData(t *testing.T) {
	metric := NewBuilder().Build()

	memStats := &runtime.MemStats{}

	heapStatsData := prepareHeapStatsData(metric, memStats)

	assert.Equal(t, heapMetric, heapStatsData.MetricName)

	assert.Equal(t, memStats.HeapAlloc, heapStatsData.Metrics[heapAlloc])
	assert.Equal(t, memStats.HeapSys, heapStatsData.Metrics[heapSys])
	assert.Equal(t, memStats.HeapInuse, heapStatsData.Metrics[heapInuse])
	assert.Equal(t, memStats.HeapObjects, heapStatsData.Metrics[heapObjects])
}
