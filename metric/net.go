package metric

import (
	"github.com/thundra-io/thundra-lambda-agent-go/plugin"
)

const ALL = 0

type netIOStatsData struct {
	Id                 string `json:"id"`
	ApplicationName    string `json:"applicationName"`
	ApplicationId      string `json:"applicationId"`
	ApplicationVersion string `json:"applicationVersion"`
	ApplicationProfile string `json:"applicationProfile"`
	ApplicationType    string `json:"applicationType"`
	StatName           string `json:"statName"`
	StatTime           string `json:"statTime"`

	BytesRecv uint64 `json:"bytesRecv"`
	BytesSent uint64 `json:"bytesSent"`
	ErrIn     uint64 `json:"errIn"`
	ErrOut    uint64 `json:"errOut"`
}

func prepareNetIOStatsData(metric *Metric) netIOStatsData {
	br, bs, ei, eo := takeNetIOFrame(metric)

	return netIOStatsData{
		Id:                 plugin.GenerateNewId(),
		ApplicationName:    metric.applicationName,
		ApplicationId:      metric.applicationId,
		ApplicationVersion: metric.applicationVersion,
		ApplicationProfile: metric.applicationProfile,
		ApplicationType:    plugin.ApplicationType,
		StatName:           netIOStat,
		StatTime:           metric.statTime.Format(plugin.TimeFormat),

		BytesRecv: br,
		BytesSent: bs,
		ErrIn:     ei,
		ErrOut:    eo,
	}
}

//Since lambda works continuously we should substract io values in order to get correct results per invocation
func takeNetIOFrame(metric *Metric) (uint64, uint64, uint64, uint64) {
	br := metric.currNetIOStat.BytesRecv - metric.prevNetIOStat.BytesRecv
	bs := metric.currNetIOStat.BytesSent - metric.prevNetIOStat.BytesSent
	ei := metric.currNetIOStat.Errin - metric.prevNetIOStat.Errin
	eo := metric.currNetIOStat.Errout - metric.prevNetIOStat.Errout

	metric.prevNetIOStat = metric.currNetIOStat
	return br, bs, ei, eo
}
