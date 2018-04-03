package metric

import (
	"github.com/aws/aws-lambda-go/lambdacontext"
	"os"
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/thundra-io/thundra-lambda-agent-go/plugin"
	"runtime"
	"time"
)

const (
	functionName       = "TestFunctionName"
	logStreamName      = "2018/01/01/[$LATEST]1234567890"
	appId              = "1234567890"
	functionVersion    = "$Version"
	applicationProfile = "TestProfile"
)

func TestInitStatData(t *testing.T) {
	prepareEnvironment()

	m := &Metric{}
	initStatData(m)
	assert.Equal(t, functionName, m.ApplicationName)
	assert.Equal(t, appId, m.ApplicationId)
	assert.Equal(t, functionVersion, m.ApplicationVersion)
	assert.Equal(t, applicationProfile, m.ApplicationProfile)
	assert.Equal(t, plugin.ApplicationType, m.ApplicationType)

	cleanEnvironment()
}

func TestPrepareHeapStatsData(t *testing.T) {
	prepareEnvironment()

	metric := &Metric{
		statTime: time.Now(),
	}
	initStatData(metric)

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

func TestPrepareStackStatsData(t *testing.T) {
	prepareEnvironment()

	metric := &Metric{
		statTime: time.Now(),
	}
	initStatData(metric)
	memStats := &runtime.MemStats{}

	stackStatsData := prepareStackStatsData(metric, memStats)

	assert.Equal(t, functionName, stackStatsData.ApplicationName)
	assert.Equal(t, appId, stackStatsData.ApplicationId)
	assert.Equal(t, functionVersion, stackStatsData.ApplicationVersion)
	assert.Equal(t, applicationProfile, stackStatsData.ApplicationProfile)
	assert.Equal(t, plugin.ApplicationType, stackStatsData.ApplicationType)

	assert.Equal(t, stackStat, stackStatsData.StatName)
	assert.Equal(t, metric.statTime.Format(plugin.TimeFormat), stackStatsData.StatTime)

	assert.Equal(t, memStats.StackInuse, stackStatsData.StackInuse)
	assert.Equal(t, memStats.StackSys, stackStatsData.StackSys)

	cleanEnvironment()
}

func TestPrepareGCStatsData(t *testing.T) {
	prepareEnvironment()

	metric := &Metric{
		statTime:     time.Now(),
		startGCCount: 1,
		endGCCount:   2,
	}
	initStatData(metric)
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
	assert.Equal(t, memStats.PauseNs[(memStats.NumGC+255)%256], gcStatsData.RecentPauseNs)
	assert.Equal(t, memStats.NumGC, gcStatsData.NumGC)
	assert.Equal(t, memStats.GCCPUFraction, gcStatsData.GCCPUFraction)

	//DeltaGCCount equals to endGCCount - startGCCount
	assert.Equal(t, uint32(1), gcStatsData.DeltaGcCount)

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
