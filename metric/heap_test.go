package metric

import (
	"testing"
	"runtime"
	"github.com/stretchr/testify/assert"
)

func TestPrepareHeapStatsData(t *testing.T) {
	metric := NewMetric()

	memStats := &runtime.MemStats{}

	heapStatsData := prepareHeapStatsData(metric, memStats)

	assert.Equal(t, heapStat, heapStatsData.StatName)

	assert.Equal(t, memStats.HeapAlloc, heapStatsData.HeapAlloc)
	assert.Equal(t, memStats.HeapSys, heapStatsData.HeapSys)
	assert.Equal(t, memStats.HeapInuse, heapStatsData.HeapInuse)
	assert.Equal(t, memStats.HeapObjects, heapStatsData.HeapObjects)
}
