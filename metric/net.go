package metric

import (
	"github.com/thundra-io/thundra-lambda-agent-go/plugin"
)

const ALL = 0

type netStatsData struct {
	Id                 string `json:"id"`
	ApplicationName    string `json:"applicationName"`
	ApplicationId      string `json:"applicationId"`
	ApplicationVersion string `json:"applicationVersion"`
	ApplicationProfile string `json:"applicationProfile"`
	ApplicationType    string `json:"applicationType"`
	StatName           string `json:"statName"`
	StatTimestamp      int64  `json:"statTimestamp"`

	BytesRecv   uint64 `json:"bytesRecv"`
	BytesSent   uint64 `json:"bytesSent"`
	PacketsRecv uint64 `json:"packetsRecv"`
	PacketsSent uint64 `json:"packetsSent"`
	ErrIn       uint64 `json:"errIn"`
	ErrOut      uint64 `json:"errOut"`
}

func prepareNetStatsData(metric *Metric) netStatsData {
	nf := takeNetFrame(metric)

	return netStatsData{
		Id:                 plugin.GenerateNewId(),
		ApplicationName:    metric.applicationName,
		ApplicationId:      metric.applicationId,
		ApplicationVersion: metric.applicationVersion,
		ApplicationProfile: metric.applicationProfile,
		ApplicationType:    plugin.ApplicationType,
		StatName:           netStat,
		StatTimestamp:      metric.statTimestamp,

		BytesRecv:   nf.bytesRecv,
		BytesSent:   nf.bytesSent,
		PacketsRecv: nf.packetsRecv,
		PacketsSent: nf.packetsSent,
		ErrIn:       nf.errin,
		ErrOut:      nf.errout,
	}
}

type netFrame struct {
	bytesSent   uint64
	bytesRecv   uint64
	packetsRecv uint64
	packetsSent uint64
	errin       uint64
	errout      uint64
}

//Since lambda works continuously we should subtract io values in order to get correct results per invocation
func takeNetFrame(metric *Metric) *netFrame {
	br := metric.currNetIOStat.BytesRecv - metric.prevNetIOStat.BytesRecv
	bs := metric.currNetIOStat.BytesSent - metric.prevNetIOStat.BytesSent
	ps := metric.currNetIOStat.PacketsSent - metric.prevNetIOStat.PacketsSent
	pr := metric.currNetIOStat.PacketsRecv - metric.prevNetIOStat.PacketsRecv
	ei := metric.currNetIOStat.Errin - metric.prevNetIOStat.Errin
	eo := metric.currNetIOStat.Errout - metric.prevNetIOStat.Errout

	metric.prevNetIOStat = metric.currNetIOStat
	return &netFrame{
		bytesRecv:   br,
		bytesSent:   bs,
		packetsRecv: pr,
		packetsSent: ps,
		errin:       ei,
		errout:      eo,
	}
}
