package metric

import (
	"context"
	"encoding/json"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thundra-io/thundra-lambda-agent-go/v2/test"
)

func TestNewMetricPlugin(t *testing.T) {
	test.PrepareEnvironment()
	mp := New()

	assert.NotNil(t, proc)
	assert.NotNil(t, mp.data)

	assert.False(t, mp.disableDiskMetrics)
	assert.False(t, mp.disableNetMetrics)
	assert.False(t, mp.disableCPUMetrics)
	assert.False(t, mp.disableGoroutineMetrics)
	assert.False(t, mp.disableHeapMetrics)
	assert.False(t, mp.disableGCMetrics)
	test.CleanEnvironment()
}

func TestMetric_BeforeExecution(t *testing.T) {
	const MaxUint32 = ^uint32(0)
	const MaxUint64 = ^uint64(0)

	mp := New()
	mp.data.startGCCount = MaxUint32
	mp.data.startPauseTotalNs = MaxUint64

	mp.BeforeExecution(context.TODO(), json.RawMessage{})
	assert.NotNil(t, mp)

	// In order to ensure startGCCount and startPauseTotalNs are assigned,
	// check it's initial value is changed.
	// Initial values are the maximum numbers to eliminate unlucky conditions from happenning.
	assert.NotEqual(t, MaxUint32, mp.data.startGCCount)
	assert.NotEqual(t, MaxUint64, mp.data.startPauseTotalNs)
}

func TestMetric_AfterExecution(t *testing.T) {

	mp := New()

	stats, _ := mp.AfterExecution(context.TODO(), json.RawMessage{}, nil, nil)

	// Assert all stats are collected, heap, gc, goroutine, cpu, net, disk
	// Note that this fails on MACOSX and returns 6 instead of 7
	if runtime.GOOS != "darwin" {
		assert.Equal(t, 7, len(stats))
	}

	for _, stat := range stats {
		assert.Equal(t, metricType, stat.Type)
	}
}
