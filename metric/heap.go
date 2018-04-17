package metric

import (
	"runtime"
	"github.com/thundra-io/thundra-lambda-agent-go/plugin"
)

type heapStatsData struct {
	Id                 string `json:"id"`
	ApplicationName    string `json:"applicationName"`
	ApplicationId      string `json:"applicationId"`
	ApplicationVersion string `json:"applicationVersion"`
	ApplicationProfile string `json:"applicationProfile"`
	ApplicationType    string `json:"applicationType"`
	StatName           string `json:"statName"`
	StatTimestamp      int64  `json:"statTimestamp"`

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

func prepareHeapStatsData(metric *Metric, memStats *runtime.MemStats) heapStatsData {
	return heapStatsData{
		Id:                 plugin.GenerateNewId(),
		ApplicationName:    metric.applicationName,
		ApplicationId:      metric.applicationId,
		ApplicationVersion: metric.applicationVersion,
		ApplicationProfile: metric.applicationProfile,
		ApplicationType:    plugin.ApplicationType,
		StatName:           heapStat,
		StatTimestamp:      metric.statTimestamp,

		HeapAlloc:   memStats.HeapAlloc,
		HeapSys:     memStats.HeapSys,
		HeapInuse:   memStats.HeapInuse,
		HeapObjects: memStats.HeapObjects,
	}
}
