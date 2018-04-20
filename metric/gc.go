package metric

import (
	"runtime"
	"github.com/thundra-io/thundra-lambda-agent-go/plugin"
)

type gcStatsData struct {
	Id                 string `json:"id"`
	ApplicationName    string `json:"applicationName"`
	ApplicationId      string `json:"applicationId"`
	ApplicationVersion string `json:"applicationVersion"`
	ApplicationProfile string `json:"applicationProfile"`
	ApplicationType    string `json:"applicationType"`
	StatName           string `json:"statName"`
	StatTimestamp      int64  `json:"statTimestamp"`

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

func prepareGCStatsData(metric *metric, memStats *runtime.MemStats) gcStatsData {
	return gcStatsData{
		Id:                 plugin.GenerateNewId(),
		ApplicationName:    metric.applicationName,
		ApplicationId:      metric.applicationId,
		ApplicationVersion: metric.applicationVersion,
		ApplicationProfile: metric.applicationProfile,
		ApplicationType:    plugin.ApplicationType,
		StatName:           gcStat,
		StatTimestamp:      metric.statTimestamp,

		PauseTotalNs:      memStats.PauseTotalNs,
		PauseNs:           memStats.PauseNs[(memStats.NumGC+255)%256],
		NumGC:             memStats.NumGC,
		NextGC:            memStats.NextGC,
		GCCPUFraction:     memStats.GCCPUFraction,
		DeltaNumGc:        metric.endGCCount - metric.startGCCount,
		DeltaPauseTotalNs: metric.endPauseTotalNs - metric.startPauseTotalNs,
	}
}
