package metric

import (
	"github.com/thundra-io/thundra-lambda-agent-go/plugin"
)

const all = 0

type netStatsData struct {
	Id                 string `json:"id"`
	TransactionId      string `json:"transactionId"`
	ApplicationName    string `json:"applicationName"`
	ApplicationId      string `json:"applicationId"`
	ApplicationVersion string `json:"applicationVersion"`
	ApplicationProfile string `json:"applicationProfile"`
	ApplicationType    string `json:"applicationType"`
	StatName           string `json:"statName"`
	StatTimestamp      int64  `json:"statTimestamp"`

	// BytesRecv is how many bytes received from network
	BytesRecv uint64 `json:"bytesRecv"`

	// BytesSent is how many bytes sent to network
	BytesSent uint64 `json:"bytesSent"`

	// PacketsRecv is how many packets received from network
	PacketsRecv uint64 `json:"packetsRecv"`

	// PacketsSent is how many packets sent to network
	PacketsSent uint64 `json:"packetsSent"`

	// ErrIn is the number of errors while sending packet
	ErrIn uint64 `json:"errIn"`

	// ErrOut is the number of errors while receiving packet
	ErrOut uint64 `json:"errOut"`
}

func prepareNetStatsData(metric *metric) netStatsData {
	nf := takeNetFrame(metric)

	return netStatsData{
		Id:                 plugin.GenerateNewId(),
		TransactionId:      plugin.TransactionId,
		ApplicationName:    plugin.ApplicationName,
		ApplicationId:      plugin.ApplicationId,
		ApplicationVersion: plugin.ApplicationVersion,
		ApplicationProfile: plugin.ApplicationProfile,
		ApplicationType:    plugin.ApplicationType,
		StatName:           netStat,
		StatTimestamp:      metric.span.statTimestamp,

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
func takeNetFrame(metric *metric) *netFrame {
	br := metric.span.currNetStat.BytesRecv - metric.span.prevNetStat.BytesRecv
	bs := metric.span.currNetStat.BytesSent - metric.span.prevNetStat.BytesSent
	ps := metric.span.currNetStat.PacketsSent - metric.span.prevNetStat.PacketsSent
	pr := metric.span.currNetStat.PacketsRecv - metric.span.prevNetStat.PacketsRecv
	ei := metric.span.currNetStat.Errin - metric.span.prevNetStat.Errin
	eo := metric.span.currNetStat.Errout - metric.span.prevNetStat.Errout

	metric.span.prevNetStat = metric.span.currNetStat
	return &netFrame{
		bytesRecv:   br,
		bytesSent:   bs,
		packetsRecv: pr,
		packetsSent: ps,
		errin:       ei,
		errout:      eo,
	}
}
