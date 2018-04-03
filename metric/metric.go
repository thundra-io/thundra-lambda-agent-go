package metric

import (
	"context"
	"encoding/json"
	"sync"
	"time"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"strings"
	"os"
	"github.com/thundra-io/thundra-lambda-agent-go/plugin"
	"github.com/satori/go.uuid"
	"fmt"
	"runtime"
)

const StatDataType = "StatData"

var uniqueId uuid.UUID

type Metric struct {
	statData
	statTime     time.Time
	startGCCount uint32
	endGCCount   uint32
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
	HeapObjects uint64 `json:"heapObject"`
}

type stackStatsData struct {
	Id                 string `json:"id"`
	ApplicationName    string `json:"applicationName"`
	ApplicationId      string `json:"applicationId"`
	ApplicationVersion string `json:"applicationVersion"`
	ApplicationProfile string `json:"applicationProfile"`
	ApplicationType    string `json:"applicationType"`
	StatName           string `json:"statName"`
	StatTime           string `json:"statTime"`

	// StackInuse is bytes in stack spans.
	StackInuse uint64 `json:"stackInuse"`

	// StackSys is bytes of stack memory obtained from the OS.
	//
	// StackSys is StackInuse, plus any memory obtained directly
	// from the OS for OS thread stacks (which should be minimal).
	StackSys uint64 `json:"stackSys"`
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
	RecentPauseNs uint64 `json:"recentPauseNs"`

	// NumGC is the number of completed GC cycles.
	NumGC uint32 `json:"numGC"`

	// GCCPUFraction is the fraction of this program's available
	// CPU time used by the GC since the program started.
	GCCPUFraction float64 `json:"gcCPUFraction"`

	DeltaGcCount uint32 `json:"deltaGcCount"`
}

func (metric *Metric) BeforeExecution(ctx context.Context, request json.RawMessage, wg *sync.WaitGroup) {
	fmt.Println("Metric: BeforeExecution")

	m := &runtime.MemStats{}
	runtime.ReadMemStats(m)

	metric.startGCCount = m.NumGC
	initStatData(metric)
	wg.Done()
}

func initStatData(metric *Metric) {
	appId := splitAppId(lambdacontext.LogStreamName)
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

func (metric *Metric) AfterExecution(ctx context.Context, request json.RawMessage, response interface{}, err interface{}, wg *sync.WaitGroup) (interface{}, string) {
	defer wg.Done()
	fmt.Println("Metric: AfterExecution")

	m := &runtime.MemStats{}
	runtime.ReadMemStats(m)

	metric.endGCCount = m.NumGC
	metric.statTime = time.Now().Round(time.Millisecond)

	heap := prepareHeapStatsData(metric, m)
	stack := prepareStackStatsData(metric, m)
	gc := prepareGCStatsData(metric, m)
	fmt.Print("Heap Metrics: ", heap)
	fmt.Print("Stack Metrics: ", stack)
	fmt.Print("GC Metrics: ", gc)
	//TODO return all types fo data in an array
	return gc, StatDataType
}

func (metric *Metric) OnPanic(ctx context.Context, request json.RawMessage, err interface{}, stackTrace []byte, wg *sync.WaitGroup) (interface{}, string) {
	defer wg.Done()

	m := &runtime.MemStats{}
	runtime.ReadMemStats(m)

	gm := prepareGCStatsData(metric, m)

	//TODO return all types fo data in an array
	return gm, StatDataType
}

func prepareHeapStatsData(metric *Metric, memStats *runtime.MemStats) heapStatsData {
	return heapStatsData{
		Id:                 generateNewId(),
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

func prepareStackStatsData(metric *Metric, memStats *runtime.MemStats) stackStatsData {
	return stackStatsData{
		Id:                 generateNewId(),
		ApplicationName:    metric.ApplicationName,
		ApplicationId:      metric.ApplicationId,
		ApplicationVersion: metric.ApplicationVersion,
		ApplicationProfile: metric.ApplicationProfile,
		ApplicationType:    plugin.ApplicationType,
		StatName:           stackStat,
		StatTime:           metric.statTime.Format(plugin.TimeFormat),

		StackInuse: memStats.StackInuse,
		StackSys:   memStats.StackSys,
	}
}

func prepareGCStatsData(metric *Metric, memStats *runtime.MemStats) gcStatsData {

	return gcStatsData{
		Id:                 generateNewId(),
		ApplicationName:    metric.ApplicationName,
		ApplicationId:      metric.ApplicationId,
		ApplicationVersion: metric.ApplicationVersion,
		ApplicationProfile: metric.ApplicationProfile,
		ApplicationType:    plugin.ApplicationType,
		StatName:           gcStat,
		StatTime:           metric.statTime.Format(plugin.TimeFormat),

		PauseTotalNs:  memStats.PauseTotalNs,
		RecentPauseNs: memStats.PauseNs[(memStats.NumGC+255)%256],
		NumGC:         memStats.NumGC,
		GCCPUFraction: memStats.GCCPUFraction,
		DeltaGcCount:  metric.endGCCount - metric.startGCCount,
	}
}

func generateNewId() string {
	return uuid.Must(uuid.NewV4()).String()
}

func splitAppId(logStreamName string) string {
	s := strings.Split(logStreamName, "]")
	if len(s) > 1 {
		return s[1]
	} else {
		return ""
	}
}
