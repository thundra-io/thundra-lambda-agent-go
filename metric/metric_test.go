package metric

import (
	"context"
	"encoding/json"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"runtime"
)

func TestMetric_BeforeExecution(t *testing.T) {
	const MaxUint32 = ^uint32(0)
	const MaxUint64 = ^uint64(0)

	m := NewBuilder().Build()
	m.span.startGCCount = MaxUint32
	m.span.startPauseTotalNs = MaxUint64

	wg := sync.WaitGroup{}
	wg.Add(1)
	m.BeforeExecution(context.TODO(), json.RawMessage{}, &wg)
	assert.NotNil(t, m.span)

	// In order to ensure startGCCount and startPauseTotalNs are assigned,
	// check it's initial value is changed.
	// Initial values are the maximum numbers to eliminate unlucky conditions from happenning.
	assert.NotEqual(t, MaxUint32, m.span.startGCCount)
	assert.NotEqual(t, MaxUint64, m.span.startPauseTotalNs)
}

func TestMetric_AfterExecution(t *testing.T) {

	m := NewBuilder().Build()

	wg := sync.WaitGroup{}
	wg.Add(1)
	stats, dataType := m.AfterExecution(context.TODO(), json.RawMessage{}, nil, nil)

	// Assert all stats are collected, heap, gc, goroutine, cpu, net, disk
	// Note that this fails on MACOSX and returns 5 instead of 6
	if runtime.GOOS != "darwin" {
		assert.Equal(t, 6, len(stats))
	}

	assert.Nil(t, m.span)
	assert.Equal(t, statDataType, dataType)
}
