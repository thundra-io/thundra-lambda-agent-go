package metric

import (
	"context"
	"encoding/json"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thundra-io/thundra-lambda-agent-go/plugin"
	"runtime"
)

func TestMetric_BeforeExecution(t *testing.T) {
	const MaxUint32 = ^uint32(0)
	const MaxUint64 = ^uint64(0)

	m := NewBuilder().Build()
	m.startGCCount = MaxUint32
	m.startPauseTotalNs = MaxUint64

	wg := sync.WaitGroup{}
	wg.Add(1)
	m.BeforeExecution(context.TODO(), json.RawMessage{}, &wg)

	// In order to ensure startGCCount and startPauseTotalNs are assigned,
	// check it's initial value is changed.
	// Initial values are the maximum numbers to eliminate unlucky conditions from happenning.
	assert.NotEqual(t, MaxUint32, m.startGCCount)
	assert.NotEqual(t, MaxUint64, m.startPauseTotalNs)
}

func TestMetric_AfterExecution(t *testing.T) {
	const MaxUint32 = ^uint32(0)
	const MaxUint64 = ^uint64(0)

	m := NewBuilder().Build()
	m.endGCCount = MaxUint32
	m.endPauseTotalNs = MaxUint64

	wg := sync.WaitGroup{}
	wg.Add(1)
	stats, dataType := m.AfterExecution(context.TODO(), json.RawMessage{}, nil, nil)

	// Assert all stats are collected, heap, gc, goroutine, cpu, net, disk
	// Note that this fails on MACOSX and returns 5 instead of 6
	if runtime.GOOS != "darwin" {
		assert.Equal(t, 6, len(stats))
	}

	// In order to ensure endGCCount and endPauseTotalNs are assigned,
	// check it's initial value is changed.
	// Initial values are the maximum numbers to eliminate unlucky conditions from happenning.
	assert.NotEqual(t, MaxUint32, m.endGCCount)
	assert.NotEqual(t, MaxUint64, m.endPauseTotalNs)

	assert.True(t, m.statTimestamp <= plugin.GetTimestamp())
	assert.Equal(t, statDataType, dataType)
}
