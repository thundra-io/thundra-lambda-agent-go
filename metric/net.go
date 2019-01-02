package metric

import (
	"fmt"

	"github.com/thundra-io/thundra-lambda-agent-go/plugin"
	"github.com/shirou/gopsutil/net"
)

const all = 0

func prepareNetMetricsData(mp *metricPlugin) metricData {
	nf := takeNetFrame(mp)
	return metricData{
		ID:                        plugin.GenerateNewID(),
		Type:                      metricType,
		AgentVersion:              plugin.AgentVersion,
		DataModelVersion:          plugin.DataModelVersion,
		ApplicationID:             plugin.ApplicationID,
		ApplicationDomainName:     plugin.ApplicationDomainName,
		ApplicationClassName:      plugin.ApplicationClassName,
		ApplicationName:           plugin.FunctionName,
		ApplicationVersion:        plugin.ApplicationVersion,
		ApplicationStage:          plugin.ApplicationStage,
		ApplicationRuntime:        plugin.ApplicationRuntime,
		ApplicationRuntimeVersion: plugin.ApplicationRuntimeVersion,
		ApplicationTags:           map[string]interface{}{},

		TraceID:         plugin.TraceID,
		TransactionID:  plugin.TransactionID,
		// SpanID:          plugin.SpanID, // Optional
		MetricName:      netMetric,
		MetricTimestamp: mp.metricTimestamp,

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
func takeNetFrame(mp *metricPlugin) *netFrame {
	// If nil, return an empty netFrame
	if mp.endNetStat == nil || mp.startNetStat == nil {
		return &netFrame{}
	}

	br := mp.endNetStat.BytesRecv - mp.startNetStat.BytesRecv
	bs := mp.endNetStat.BytesSent - mp.startNetStat.BytesSent
	ps := mp.endNetStat.PacketsSent - mp.startNetStat.PacketsSent
	pr := mp.endNetStat.PacketsRecv - mp.startNetStat.PacketsRecv
	ei := mp.endNetStat.Errin - mp.startNetStat.Errin
	eo := mp.endNetStat.Errout - mp.startNetStat.Errout

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
