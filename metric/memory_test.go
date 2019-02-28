package metric

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrepareMemoryMetricsData(t *testing.T) {
	mp := New()

	base := mp.prepareMetricsData()
	memoryMetricsData := prepareMemoryMetricsData(mp, base)

	assert.True(t, len(memoryMetricsData.ID) != 0)
	assert.Equal(t, memoryMetric, memoryMetricsData.MetricName)
}
