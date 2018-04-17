package metric

import (
	"github.com/thundra-io/thundra-lambda-agent-go/plugin"
)

type ioStatsData struct {
	Id                 string `json:"id"`
	ApplicationName    string `json:"applicationName"`
	ApplicationId      string `json:"applicationId"`
	ApplicationVersion string `json:"applicationVersion"`
	ApplicationProfile string `json:"applicationProfile"`
	ApplicationType    string `json:"applicationType"`
	StatName           string `json:"statName"`
	StatTime           string `json:"statTime"`

	ReadBytes  uint64 `json:"readBytes"`
	WriteBytes uint64 `json:"writeBytes"`
}

func prepareIOStatsData(metric *Metric) ioStatsData {
	rb, wb := takeIOFrame(metric)

	return ioStatsData{
		Id:                 plugin.GenerateNewId(),
		ApplicationName:    metric.applicationName,
		ApplicationId:      metric.applicationId,
		ApplicationVersion: metric.applicationVersion,
		ApplicationProfile: metric.applicationProfile,
		ApplicationType:    plugin.ApplicationType,
		StatName:           ioStat,
		StatTime:           metric.statTime.Format(plugin.TimeFormat),

		ReadBytes:  rb,
		WriteBytes: wb,
	}
}

//Since lambda works continuously we should substract io values in order to get correct results per invocation
func takeIOFrame(metric *Metric) (uint64, uint64) {

	rb := metric.currIOStat.ReadBytes - metric.prevIOStat.ReadBytes
	wb := metric.currIOStat.WriteBytes - metric.prevIOStat.WriteBytes

	metric.prevIOStat = metric.currIOStat
	return rb, wb
}
