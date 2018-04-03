package metric

import (
	"context"
	"encoding/json"
	"sync"
	"time"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"os"
	"github.com/thundra-io/thundra-lambda-agent-go/plugin"
	"fmt"
	"runtime"
)

const StatDataType = "StatData"

type Metric struct {
	statData
	statTime          time.Time
	startGCCount      uint32
	endGCCount        uint32
	startPauseTotalNs uint64
	endPauseTotalNs   uint64

	EnableGCStats        bool
	EnableHeapStats      bool
	EnableGoroutineStats bool
	EnableCPUStats       bool
}

type statData struct {
	ApplicationName    string `json:"applicationName"`
	ApplicationId      string `json:"applicationId"`
	ApplicationVersion string `json:"applicationVersion"`
	ApplicationProfile string `json:"applicationProfile"`
	ApplicationType    string `json:"applicationType"`
}

type heapStatsData struct {
	Id                 string `json:"id"`
	ApplicationName    string `json:"applicationName"`
	ApplicationId      string `json:"applicationId"`
	ApplicationVersion string `json:"applicationVersion"`
	ApplicationProfile string `json:"applicationProfile"`
	ApplicationType    string `json:"applicationType"`
	StatName           string `json:"statName"`
	StatTime           string `json:"statTime"`

	// HeapAlloc is bytes of allocated heap objects.
	//
	// "Allocated" heap objects include all reachable objects, as
	// well as unreachable objects that the garbage collector has
	// not yet freed.
	HeapAlloc uint64 `json:"heapAlloc"`

	// HeapSys estimates the largest size the heap has had.
	HeapSys uint64 `json:"heapSys"`

	// HeapInuse is bytes in in-use spans.

	// In-use spans have at least one object in them. These spans
	// can only be used for other objects of roughly the same
	// size.
	HeapInuse uint64 `json:"heapInuse"`

	// HeapObjects is the number of allocated heap objects.
	HeapObjects uint64 `json:"heapObjects"`
}

type gcStatsData struct {
	Id                 string `json:"id"`
	ApplicationName    string `json:"applicationName"`
	ApplicationId      string `json:"applicationId"`
	ApplicationVersion string `json:"applicationVersion"`
	ApplicationProfile string `json:"applicationProfile"`
	ApplicationType    string `json:"applicationType"`
	StatName           string `json:"statName"`
	StatTime           string `json:"statTime"`

	// PauseTotalNs is the cumulative nanoseconds in GC
	// stop-the-world pauses since the program started.
	PauseTotalNs uint64 `json:"pauseTotalNs"`

	// PauseNs is recent GC stop-the-world pause time in nanoseconds.
	PauseNs uint64 `json:"pauseNs"`

	// NumGC is the number of completed GC cycles.
	NumGC uint32 `json:"numGC"`

	// NextGC is the target heap size of the next GC cycle.
	NextGC uint64 `json:"nextGC"`

	// GCCPUFraction is the fraction of this program's available
	// CPU time used by the GC since the program started.
	GCCPUFraction float64 `json:"gcCPUFraction"`

	//DeltaNumGc is the change in NUMGC from before execution to after execution
	DeltaNumGc uint32 `json:"deltaNumGC"`

	//DeltaPauseTotalNs is pause total change from before execution to after execution
	DeltaPauseTotalNs uint64 `json:"deltaPauseTotalNs"`
}

type goRoutineStatsData struct {
	Id                 string `json:"id"`
	ApplicationName    string `json:"applicationName"`
	ApplicationId      string `json:"applicationId"`
	ApplicationVersion string `json:"applicationVersion"`
	ApplicationProfile string `json:"applicationProfile"`
	ApplicationType    string `json:"applicationType"`
	StatName           string `json:"statName"`
	StatTime           string `json:"statTime"`
	NumGoroutine       uint64 `json:"numGoroutine"`
}

type cpuStatsData struct {
	Id                 string `json:"id"`
	ApplicationName    string `json:"applicationName"`
	ApplicationId      string `json:"applicationId"`
	ApplicationVersion string `json:"applicationVersion"`
	ApplicationProfile string `json:"applicationProfile"`
	ApplicationType    string `json:"applicationType"`
	StatName           string `json:"statName"`
	StatTime           string `json:"statTime"`
	NumCPU             uint64 `json:"numCPU"`
}

func (metric *Metric) BeforeExecution(ctx context.Context, request json.RawMessage, wg *sync.WaitGroup) {
	m := &runtime.MemStats{}
	runtime.ReadMemStats(m)

	if metric.EnableGCStats {
		metric.startGCCount = m.NumGC
		metric.startPauseTotalNs = m.PauseTotalNs
	}

	initStatData(metric)
	wg.Done()
}

func initStatData(metric *Metric) {
	appId := plugin.SplitAppId(lambdacontext.LogStreamName)
	ver := lambdacontext.FunctionVersion
	profile := os.Getenv(plugin.ThundraApplicationProfile)
	if profile == "" {
		profile = plugin.DefaultProfile
	}

	metric.ApplicationName = lambdacontext.FunctionName
	metric.ApplicationId = appId
	metric.ApplicationVersion = ver
	metric.ApplicationProfile = profile
	metric.ApplicationType = plugin.ApplicationType
}

func (metric *Metric) AfterExecution(ctx context.Context, request json.RawMessage, response interface{}, err interface{}) ([]interface{}, string) {
	mStats := &runtime.MemStats{}
	runtime.ReadMemStats(mStats)

	metric.statTime = time.Now().Round(time.Millisecond)

	var stats []interface{}
	if metric.EnableHeapStats {
		heap := prepareHeapStatsData(metric, mStats)
		stats = append(stats, heap)
	}

	if metric.EnableGCStats {
		metric.endGCCount = mStats.NumGC
		metric.endPauseTotalNs = mStats.PauseTotalNs

		gc := prepareGCStatsData(metric, mStats)
		stats = append(stats, gc)
	}

	if metric.EnableGoroutineStats {
		gs := prepareGoRoutineStatsData(metric)
		stats = append(stats, gs)
	}

	if metric.EnableCPUStats {
		cs := prepareCPUStatsData(metric)
		stats = append(stats, cs)
	}

	fmt.Println("Metrics: ", stats)
	return stats, StatDataType
}

func (metric *Metric) OnPanic(ctx context.Context, request json.RawMessage, err interface{}, stackTrace []byte) ([]interface{}, string) {
	//TODO return all types fo data in an array
	return nil, StatDataType
}

func prepareHeapStatsData(metric *Metric, memStats *runtime.MemStats) heapStatsData {
	return heapStatsData{
		Id:                 plugin.GenerateNewId(),
		ApplicationName:    metric.ApplicationName,
		ApplicationId:      metric.ApplicationId,
		ApplicationVersion: metric.ApplicationVersion,
		ApplicationProfile: metric.ApplicationProfile,
		ApplicationType:    plugin.ApplicationType,
		StatName:           heapStat,
		StatTime:           metric.statTime.Format(plugin.TimeFormat),

		HeapAlloc:   memStats.HeapAlloc,
		HeapSys:     memStats.HeapSys,
		HeapInuse:   memStats.HeapInuse,
		HeapObjects: memStats.HeapObjects,
	}
}

func prepareGCStatsData(metric *Metric, memStats *runtime.MemStats) gcStatsData {
	return gcStatsData{
		Id:                 plugin.GenerateNewId(),
		ApplicationName:    metric.ApplicationName,
		ApplicationId:      metric.ApplicationId,
		ApplicationVersion: metric.ApplicationVersion,
		ApplicationProfile: metric.ApplicationProfile,
		ApplicationType:    plugin.ApplicationType,
		StatName:           gcStat,
		StatTime:           metric.statTime.Format(plugin.TimeFormat),

		PauseTotalNs:      memStats.PauseTotalNs,
		PauseNs:           memStats.PauseNs[(memStats.NumGC+255)%256],
		NumGC:             memStats.NumGC,
		NextGC:            memStats.NextGC,
		GCCPUFraction:     memStats.GCCPUFraction,
		DeltaNumGc:        metric.endGCCount - metric.startGCCount,
		DeltaPauseTotalNs: metric.endPauseTotalNs - metric.startPauseTotalNs,
	}
}

func prepareGoRoutineStatsData(metric *Metric) goRoutineStatsData {
	return goRoutineStatsData{
		Id:                 plugin.GenerateNewId(),
		ApplicationName:    metric.ApplicationName,
		ApplicationId:      metric.ApplicationId,
		ApplicationVersion: metric.ApplicationVersion,
		ApplicationProfile: metric.ApplicationProfile,
		ApplicationType:    plugin.ApplicationType,
		StatName:           goroutineStat,
		StatTime:           metric.statTime.Format(plugin.TimeFormat),
		NumGoroutine:       uint64(runtime.NumGoroutine()),
	}
}

func prepareCPUStatsData(metric *Metric) cpuStatsData {
	return cpuStatsData{
		Id:                 plugin.GenerateNewId(),
		ApplicationName:    metric.ApplicationName,
		ApplicationId:      metric.ApplicationId,
		ApplicationVersion: metric.ApplicationVersion,
		ApplicationProfile: metric.ApplicationProfile,
		ApplicationType:    plugin.ApplicationType,
		StatName:           cpuStat,
		StatTime:           metric.statTime.Format(plugin.TimeFormat),
		NumCPU:             uint64(runtime.NumCPU()),
	}
}
