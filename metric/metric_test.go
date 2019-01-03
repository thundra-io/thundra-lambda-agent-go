package metric

import (
	"context"
	"encoding/json"
	"sync"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMetric_BeforeExecution(t *testing.T) {
	const MaxUint32 = ^uint32(0)
	const MaxUint64 = ^uint64(0)

	mp := NewBuilder().Build()
	mp.data.startGCCount = MaxUint32
	mp.data.startPauseTotalNs = MaxUint64

	wg := sync.WaitGroup{}
	wg.Add(1)
	mp.BeforeExecution(context.TODO(), json.RawMessage{}, &wg)
	assert.NotNil(t, mp)

	// In order to ensure startGCCount and startPauseTotalNs are assigned,
	// check it's initial value is changed.
	// Initial values are the maximum numbers to eliminate unlucky conditions from happenning.
	assert.NotEqual(t, MaxUint32, mp.data.startGCCount)
	assert.NotEqual(t, MaxUint64, mp.data.startPauseTotalNs)
}

func TestMetric_AfterExecution(t *testing.T) {

	mp := NewBuilder().Build()

	wg := sync.WaitGroup{}
	wg.Add(1)
	stats := mp.AfterExecution(context.TODO(), json.RawMessage{}, nil, nil)

	// Assert all stats are collected, heap, gc, goroutine, cpu, net, disk
	// Note that this fails on MACOSX and returns 6 instead of 7
	if runtime.GOOS != "darwin" {
		assert.Equal(t, 7, len(stats))
	}

	for _, stat := range stats {
		assert.Equal(t, metricType, stat.Type)
	}
}
