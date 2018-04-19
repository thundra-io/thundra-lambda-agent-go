package metric

import (
	"github.com/aws/aws-lambda-go/lambdacontext"
	"os"
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/thundra-io/thundra-lambda-agent-go/plugin"
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

func TestNewMetric(t *testing.T) {
	prepareEnvironment()
	m := NewMetric()

	assert.Equal(t, functionName, m.applicationName)
	assert.Equal(t, appId, m.applicationId)
	assert.Equal(t, functionVersion, m.applicationVersion)
	assert.Equal(t, applicationProfile, m.applicationProfile)
	assert.Equal(t, plugin.ApplicationType, m.applicationType)

	assert.NotNil(t, m.prevDiskStat)
	assert.NotNil(t, m.prevNetStat)
}

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

	assert.True(t, m.statTimestamp <= plugin.GetTimestamp())
	assert.Equal(t, StatDataType, dataType)
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
