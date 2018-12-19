package metric

import (
	"fmt"

	"github.com/thundra-io/thundra-lambda-agent-go/plugin"
	"github.com/shirou/gopsutil/net"
)

const all = 0

func prepareNetMetricsData(metric *metric) metricData {
	nf := takeNetFrame(metric)
	return metricData{
		Id:                        plugin.GenerateNewId(),
		Type:                      metricType,
		AgentVersion:              plugin.AgentVersion,
		DataModelVersion:          plugin.DataModelVersion,
		ApplicationId:             plugin.ApplicationId,
		ApplicationDomainName:     plugin.ApplicationDomainName,
		ApplicationClassName:      plugin.ApplicationClassName,
		ApplicationName:           plugin.FunctionName,
		ApplicationVersion:        plugin.ApplicationVersion,
		ApplicationStage:          plugin.ApplicationStage,
		ApplicationRuntime:        plugin.ApplicationRuntime,
		ApplicationRuntimeVersion: plugin.ApplicationRuntimeVersion,
		ApplicationTags:           map[string]interface{}{},

		TraceId:         plugin.TraceId,
		TracnsactionId:  plugin.TransactionId,
		SpanId:          plugin.SpanId,
		MetricName:      netMetric,
		MetricTimestamp: metric.span.metricTimestamp,

		Metrics: map[string]interface{}{
			// BytesRecv is how many bytes received from network
			bytesRecv: nf.bytesRecv,
			// BytesSent is how many bytes sent to network
			bytesSent: nf.bytesSent,
			// PacketsRecv is how many packets received from network
			packetsRecv: nf.packetsRecv,
			// PacketsSent is how many packets sent to network
			packetsSent: nf.packetsSent,
			// ErrIn is the number of errors while sending packet
			errIn: nf.errin,
			// ErrOut is the number of errors while receiving packet
			errOut: nf.errout,
		},
		Tags: map[string]interface{}{},
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
	// If nil, return an empty netFrame
	if metric.span.endNetStat == nil || metric.span.startNetStat == nil {
		return &netFrame{}
	}

	br := metric.span.endNetStat.BytesRecv - metric.span.startNetStat.BytesRecv
	bs := metric.span.endNetStat.BytesSent - metric.span.startNetStat.BytesSent
	ps := metric.span.endNetStat.PacketsSent - metric.span.startNetStat.PacketsSent
	pr := metric.span.endNetStat.PacketsRecv - metric.span.startNetStat.PacketsRecv
	ei := metric.span.endNetStat.Errin - metric.span.startNetStat.Errin
	eo := metric.span.endNetStat.Errout - metric.span.startNetStat.Errout

	return &netFrame{
		bytesRecv:   br,
		bytesSent:   bs,
		packetsRecv: pr,
		packetsSent: ps,
		errin:       ei,
		errout:      eo,
	}
}

func sampleNetStat() (*net.IOCountersStat) {
	netIOStat, err := net.IOCounters(false)
	if err != nil {
		fmt.Println("Error sampling net stat", err)
		return nil
	}
	return &netIOStat[all]
}
