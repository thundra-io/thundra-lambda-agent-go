package metric

import (
	"log"

	uuid "github.com/google/uuid"
	"github.com/shirou/gopsutil/net"
)

const all = 0

func prepareNetMetricsData(mp *metricPlugin, base metricDataModel) metricDataModel {
	base.ID = uuid.New().String()
	base.MetricName = netMetric
	nf := takeNetFrame(mp)
	base.Metrics = map[string]interface{}{
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
	}

	return base
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
	if mp.data.endNetStat == nil || mp.data.startNetStat == nil {
		return &netFrame{}
	}

	br := mp.data.endNetStat.BytesRecv - mp.data.startNetStat.BytesRecv
	bs := mp.data.endNetStat.BytesSent - mp.data.startNetStat.BytesSent
	ps := mp.data.endNetStat.PacketsSent - mp.data.startNetStat.PacketsSent
	pr := mp.data.endNetStat.PacketsRecv - mp.data.startNetStat.PacketsRecv
	ei := mp.data.endNetStat.Errin - mp.data.startNetStat.Errin
	eo := mp.data.endNetStat.Errout - mp.data.startNetStat.Errout

	return &netFrame{
		bytesRecv:   br,
		bytesSent:   bs,
		packetsRecv: pr,
		packetsSent: ps,
		errin:       ei,
		errout:      eo,
	}
}

func sampleNetStat() *net.IOCountersStat {
	netIOStat, err := net.IOCounters(false)
	if err != nil {
		log.Println("Error sampling net stat", err)
		return nil
	}
	return &netIOStat[all]
}
