package metric

import (
	"github.com/aws/aws-lambda-go/lambdacontext"
	"os"
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/thundra-io/thundra-lambda-agent-go/plugin"
	"runtime"
	"time"
	"context"
	"encoding/json"
	"sync"
)

const (
	functionName       = "TestFunctionName"
	logStreamName      = "2018/01/01/[$LATEST]1234567890"
	appId              = "1234567890"
	functionVersion    = "$Version"
	applicationProfile = "TestProfile"
)

func TestMetric_BeforeExecution(t *testing.T) {
	const MaxUint32 = ^uint32(0)
	const MaxUint64 = ^uint64(0)

	m := &Metric{
		EnableGCStats:     true,
		startGCCount:      MaxUint32,
		startPauseTotalNs: MaxUint64,
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	m.BeforeExecution(context.TODO(), json.RawMessage{}, &wg)

	//In order to ensure startGCCount and startPauseTotalNs are assigned,
	//check it's initial value is changed.
	//Initial values are the maximum numbers to eliminate unlucky conditions from happenning.
	assert.NotEqual(t, MaxUint32, m.startGCCount)
	assert.NotEqual(t, MaxUint64, m.startPauseTotalNs)
}

func TestMetric_AfterExecution(t *testing.T) {
	const MaxUint32 = ^uint32(0)
	const MaxUint64 = ^uint64(0)

	m := &Metric{
		EnableHeapStats:      true,
		EnableGCStats:        true,
		EnableGoroutineStats: true,
		endGCCount:           MaxUint32,
		endPauseTotalNs:      MaxUint64,
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	stats, dataType := m.AfterExecution(context.TODO(), json.RawMessage{}, nil, nil)

	//Assert Heap,GC,Goroutine and CPU stats are collected
	assert.Equal(t, 3, len(stats))

	//In order to ensure endGCCount and endPauseTotalNs are assigned,
	//check it's initial value is changed.
	//Initial values are the maximum numbers to eliminate unlucky conditions from happenning.
	assert.NotEqual(t, MaxUint32, m.endGCCount)
	assert.NotEqual(t, MaxUint64, m.endPauseTotalNs)

	now := time.Now().Round(time.Millisecond)
	assert.True(t, m.statTime.Before(now) || m.statTime.Equal(now))
	assert.Equal(t, StatDataType, dataType)
}

func TestPrepareHeapStatsData(t *testing.T) {
	prepareEnvironment()

	metric := NewMetric()
	metric.statTime = time.Now()

	memStats := &runtime.MemStats{}

	heapStatsData := prepareHeapStatsData(metric, memStats)

	assert.Equal(t, functionName, heapStatsData.ApplicationName)
	assert.Equal(t, appId, heapStatsData.ApplicationId)
	assert.Equal(t, functionVersion, heapStatsData.ApplicationVersion)
	assert.Equal(t, applicationProfile, heapStatsData.ApplicationProfile)
	assert.Equal(t, plugin.ApplicationType, heapStatsData.ApplicationType)

	assert.Equal(t, heapStat, heapStatsData.StatName)
	assert.Equal(t, metric.statTime.Format(plugin.TimeFormat), heapStatsData.StatTime)

	assert.Equal(t, memStats.HeapAlloc, heapStatsData.HeapAlloc)
	assert.Equal(t, memStats.HeapSys, heapStatsData.HeapSys)
	assert.Equal(t, memStats.HeapInuse, heapStatsData.HeapInuse)
	assert.Equal(t, memStats.HeapObjects, heapStatsData.HeapObjects)

	cleanEnvironment()
}

func TestPrepareGCStatsData(t *testing.T) {
	prepareEnvironment()

	metric := NewMetric()
	metric.statTime = time.Now()
	metric.startGCCount = 1
	metric.endGCCount = 2

	memStats := &runtime.MemStats{}

	gcStatsData := prepareGCStatsData(metric, memStats)

	assert.Equal(t, functionName, gcStatsData.ApplicationName)
	assert.Equal(t, appId, gcStatsData.ApplicationId)
	assert.Equal(t, functionVersion, gcStatsData.ApplicationVersion)
	assert.Equal(t, applicationProfile, gcStatsData.ApplicationProfile)
	assert.Equal(t, plugin.ApplicationType, gcStatsData.ApplicationType)

	assert.Equal(t, gcStat, gcStatsData.StatName)
	assert.Equal(t, metric.statTime.Format(plugin.TimeFormat), gcStatsData.StatTime)

	assert.Equal(t, memStats.PauseTotalNs, gcStatsData.PauseTotalNs)
	assert.Equal(t, memStats.PauseNs[(memStats.NumGC+255)%256], gcStatsData.PauseNs)
	assert.Equal(t, memStats.NumGC, gcStatsData.NumGC)
	assert.Equal(t, memStats.NextGC, gcStatsData.NextGC)
	assert.Equal(t, memStats.GCCPUFraction, gcStatsData.GCCPUFraction)

	//DeltaGCCount equals to endGCCount - startGCCount
	assert.Equal(t, uint32(1), gcStatsData.DeltaNumGc)

	cleanEnvironment()
}

func TestPrepareGoroutineStatsData(t *testing.T) {
	prepareEnvironment()

	metric := NewMetric()
	metric.statTime = time.Now()
	metric.startGCCount = 1
	metric.endGCCount = 2

	gcStatsData := prepareGoRoutineStatsData(metric)

	assert.Equal(t, functionName, gcStatsData.ApplicationName)
	assert.Equal(t, appId, gcStatsData.ApplicationId)
	assert.Equal(t, functionVersion, gcStatsData.ApplicationVersion)
	assert.Equal(t, applicationProfile, gcStatsData.ApplicationProfile)
	assert.Equal(t, plugin.ApplicationType, gcStatsData.ApplicationType)

	assert.Equal(t, goroutineStat, gcStatsData.StatName)
	assert.Equal(t, metric.statTime.Format(plugin.TimeFormat), gcStatsData.StatTime)

	assert.Equal(t, uint64(runtime.NumGoroutine()), gcStatsData.NumGoroutine)

	cleanEnvironment()
}

func TestPrepareCPUStatsData(t *testing.T) {
	prepareEnvironment()

	metric := NewMetric()
	metric.startGCCount = 1
	metric.endGCCount = 2

	cpuStatsData := prepareCPUStatsData(metric)

	assert.Equal(t, functionName, cpuStatsData.ApplicationName)
	assert.Equal(t, appId, cpuStatsData.ApplicationId)
	assert.Equal(t, functionVersion, cpuStatsData.ApplicationVersion)
	assert.Equal(t, applicationProfile, cpuStatsData.ApplicationProfile)
	assert.Equal(t, plugin.ApplicationType, cpuStatsData.ApplicationType)

	assert.Equal(t, cpuStat, cpuStatsData.StatName)
	assert.Equal(t, metric.statTime.Format(plugin.TimeFormat), cpuStatsData.StatTime)

	cleanEnvironment()
}

func prepareEnvironment() {
	lambdacontext.FunctionName = functionName
	lambdacontext.LogStreamName = logStreamName
	lambdacontext.FunctionVersion = functionVersion
	os.Setenv(plugin.ThundraApplicationProfile, applicationProfile)
}

func cleanEnvironment() {
	lambdacontext.FunctionName = ""
	lambdacontext.MemoryLimitInMB = 0
	lambdacontext.FunctionVersion = ""
	os.Clearenv()
}