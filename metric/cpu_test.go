package metric

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrepareCPUStatsData(t *testing.T) {
	metric := NewBuilder().Build()

	cpuStatsData := prepareCPUStatsData(metric)

	assert.Equal(t, cpuStat, cpuStatsData.StatName)
	assert.Equal(t, metric.statTimestamp, cpuStatsData.StatTimestamp)
}
