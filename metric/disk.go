package metric

import (
	"github.com/thundra-io/thundra-lambda-agent-go/plugin"
)

type diskStatsData struct {
	Id                 string `json:"id"`
	TransactionId      string `json:"transactionId"`
	ApplicationName    string `json:"applicationName"`
	ApplicationId      string `json:"applicationId"`
	ApplicationVersion string `json:"applicationVersion"`
	ApplicationProfile string `json:"applicationProfile"`
	ApplicationType    string `json:"applicationType"`
	StatName           string `json:"statName"`
	StatTimestamp      int64  `json:"statTimestamp"`

	// ReadBytes is the number of bytes read from disk
	ReadBytes uint64 `json:"readBytes"`

	// WriteBytes is the number of bytes write to disk
	WriteBytes uint64 `json:"writeBytes"`

	// ReadCount is the number read operations from disk
	ReadCount uint64 `json:"readCount"`

	// WriteCount is the number write operations to disk
	WriteCount uint64 `json:"writeCount"`
}

func prepareDiskStatsData(metric *metric) diskStatsData {
	df := takeDiskFrame(metric)

	return diskStatsData{
		Id:                 plugin.GenerateNewId(),
		TransactionId:      plugin.TransactionId,
		ApplicationName:    plugin.ApplicationName,
		ApplicationId:      plugin.ApplicationId,
		ApplicationVersion: plugin.ApplicationVersion,
		ApplicationProfile: plugin.ApplicationProfile,
		ApplicationType:    plugin.ApplicationType,
		StatName:           diskStat,
		StatTimestamp:      metric.span.statTimestamp,

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
func takeDiskFrame(metric *metric) *diskFrame {
	rb := metric.span.currDiskStat.ReadBytes - metric.span.prevDiskStat.ReadBytes
	wb := metric.span.currDiskStat.WriteBytes - metric.span.prevDiskStat.WriteBytes

	rc := metric.span.currDiskStat.ReadCount - metric.span.prevDiskStat.ReadCount
	wc := metric.span.currDiskStat.WriteCount - metric.span.prevDiskStat.WriteCount

	metric.span.prevDiskStat = metric.span.currDiskStat
	return &diskFrame{
		readBytes:  rb,
		writeBytes: wb,
		readCount:  rc,
		writeCount: wc,
	}
}
