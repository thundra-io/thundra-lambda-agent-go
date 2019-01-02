package metric

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

const garbageCollectionCount = 5

func TestPrepareGCMetricsData(t *testing.T) {
	mp := NewBuilder().Build()
	mp.startGCCount = 0
	mp.endGCCount = garbageCollectionCount

	makeMultipleGCCalls(garbageCollectionCount)
	memStats := &runtime.MemStats{}
	runtime.ReadMemStats(memStats)

	gcStatsData := prepareGCMetricsData(mp, memStats)

	assert.Equal(t, gcMetric, gcStatsData.MetricName)
	assert.Equal(t, memStats.PauseTotalNs, gcStatsData.Metrics[pauseTotalNs])
	assert.Equal(t, memStats.PauseNs[(memStats.NumGC+255)%256], gcStatsData.Metrics[pauseNs])

	assert.Equal(t, uint32(garbageCollectionCount), gcStatsData.Metrics[numGc])
	assert.Equal(t, memStats.NextGC, gcStatsData.Metrics[nextGc])
	assert.Equal(t, memStats.GCCPUFraction, gcStatsData.Metrics[gcCPUFraction])

	//DeltaGCCount equals to endGCCount - startGCCount
	assert.Equal(t, uint32(garbageCollectionCount), gcStatsData.Metrics[deltaNumGc])
}

func makeMultipleGCCalls(times int) {
	for i := 0; i < times; i++ {
		runtime.GC()
	}
}
