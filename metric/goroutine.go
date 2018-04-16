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
	StatTime           string `json:"statTime"`
	NumGoroutine       uint64 `json:"numGoroutine"`
}

func prepareGoRoutineStatsData(metric *Metric) goRoutineStatsData {
	return goRoutineStatsData{
		Id:                 plugin.GenerateNewId(),
		ApplicationName:    metric.applicationName,
		ApplicationId:      metric.applicationId,
		ApplicationVersion: metric.applicationVersion,
		ApplicationProfile: metric.applicationProfile,
		ApplicationType:    plugin.ApplicationType,
		StatName:           goroutineStat,
		StatTime:           metric.statTime.Format(plugin.TimeFormat),
		NumGoroutine:       uint64(runtime.NumGoroutine()),
	}
}