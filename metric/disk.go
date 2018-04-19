package metric

import (
	"github.com/thundra-io/thundra-lambda-agent-go/plugin"
)

type diskStatsData struct {
	Id                 string `json:"id"`
	ApplicationName    string `json:"applicationName"`
	ApplicationId      string `json:"applicationId"`
	ApplicationVersion string `json:"applicationVersion"`
	ApplicationProfile string `json:"applicationProfile"`
	ApplicationType    string `json:"applicationType"`
	StatName           string `json:"statName"`
	StatTimestamp      int64  `json:"statTimestamp"`

	ReadBytes  uint64 `json:"readBytes"`
	WriteBytes uint64 `json:"writeBytes"`
	ReadCount  uint64 `json:"readCount"`
	WriteCount uint64 `json:"writeCount"`
}

func prepareDiskStatsData(metric *Metric) diskStatsData {
	df := takeDiskFrame(metric)

	return diskStatsData{
		Id:                 plugin.GenerateNewId(),
		ApplicationName:    metric.applicationName,
		ApplicationId:      metric.applicationId,
		ApplicationVersion: metric.applicationVersion,
		ApplicationProfile: metric.applicationProfile,
		ApplicationType:    plugin.ApplicationType,
		StatName:           diskStat,
		StatTimestamp:      metric.statTimestamp,

		ReadBytes:  df.readBytes,
		WriteBytes: df.writeBytes,
		ReadCount:  df.readCount,
		WriteCount: df.writeCount,
	}
}

type diskFrame struct {
	readBytes  uint64
	writeBytes uint64
	readCount  uint64
	writeCount uint64
}

//Since lambda works continuously we should subtract io values in order to get correct results per invocation
//takeDiskFrame returns IO operations count for a specific time range
func takeDiskFrame(metric *Metric) *diskFrame {
	rb := metric.currDiskStat.ReadBytes - metric.prevDiskStat.ReadBytes
	wb := metric.currDiskStat.WriteBytes - metric.prevDiskStat.WriteBytes

	rc := metric.currDiskStat.ReadCount - metric.prevDiskStat.ReadCount
	wc := metric.currDiskStat.WriteCount - metric.prevDiskStat.WriteCount

	metric.prevDiskStat = metric.currDiskStat
	return &diskFrame{
		readBytes:  rb,
		writeBytes: wb,
		readCount:  rc,
		writeCount: wc,
	}
}
