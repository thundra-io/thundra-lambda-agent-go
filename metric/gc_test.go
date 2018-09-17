package metric

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

const numGC = 5

func TestPrepareGCStatsData(t *testing.T) {
	metric := NewBuilder().Build()
	metric.span.startGCCount = 0
	metric.span.endGCCount = numGC

	makeMultipleGCCalls(numGC)
	memStats := &runtime.MemStats{}
	runtime.ReadMemStats(memStats)

	gcStatsData := prepareGCStatsData(metric, memStats)

	assert.Equal(t, gcStat, gcStatsData.StatName)
	assert.Equal(t, memStats.PauseTotalNs, gcStatsData.PauseTotalNs)
	assert.Equal(t, memStats.PauseNs[(memStats.NumGC+255)%256], gcStatsData.PauseNs)

	assert.Equal(t, uint32(numGC), gcStatsData.NumGC)
	assert.Equal(t, memStats.NextGC, gcStatsData.NextGC)
	assert.Equal(t, memStats.GCCPUFraction, gcStatsData.GCCPUFraction)

	//DeltaGCCount equals to endGCCount - startGCCount
	assert.Equal(t, uint32(numGC), gcStatsData.DeltaNumGc)
}

func makeMultipleGCCalls(times int) {
	for i := 0; i < times; i++ {
		runtime.GC()
	}
}
