package metric

import (
	"github.com/thundra-io/thundra-lambda-agent-go/plugin"
	"runtime"
)

type goRoutineStatsData struct {
	Id                 string `json:"id"`
	ApplicationName    string `json:"applicationName"`
	ApplicationId      string `json:"applicationId"`
	ApplicationVersion string `json:"applicationVersion"`
	ApplicationProfile string `json:"applicationProfile"`
	ApplicationType    string `json:"applicationType"`
	StatName           string `json:"statName"`
	StatTimestamp      int64  `json:"statTimestamp"`

	// NumGoroutine is the number of goroutines on execution
	NumGoroutine       uint64 `json:"numGoroutine"`
}

func prepareGoRoutineStatsData(metric *metric) goRoutineStatsData {
	return goRoutineStatsData{
		Id:                 plugin.GenerateNewId(),
		ApplicationName:    metric.applicationName,
		ApplicationId:      metric.applicationId,
		ApplicationVersion: metric.applicationVersion,
		ApplicationProfile: metric.applicationProfile,
		ApplicationType:    plugin.ApplicationType,
		StatName:           goroutineStat,
		StatTimestamp:      metric.statTimestamp,
		NumGoroutine:       uint64(runtime.NumGoroutine()),
	}
}
