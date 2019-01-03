package metric

import (
	"fmt"

	"github.com/thundra-io/thundra-lambda-agent-go/plugin"
	"github.com/shirou/gopsutil/process"
)

func prepareDiskMetricsData(mp *metricPlugin) metricDataModel {
	df := takeDiskFrame(mp)
	return metricDataModel{
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
		MetricName:      diskMetric,
		MetricTimestamp: mp.data.metricTimestamp,

		Metrics: map[string]interface{}{
			// ReadBytes is the number of bytes read from disk
			readBytes: df.readBytes,
			// WriteBytes is the number of bytes write to disk
			writeBytes: df.writeBytes,
			// ReadCount is the number read operations from disk
			readCount: df.readCount,
			// WriteCount is the number write operations to disk
			writeCount: df.writeCount,
		},
		Tags: map[string]interface{}{},
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
func takeDiskFrame(mp *metricPlugin) *diskFrame {
	if mp.data.endDiskStat == nil || mp.data.startDiskStat == nil {
		return &diskFrame{}
	}
	rb := mp.data.endDiskStat.ReadBytes - mp.data.startDiskStat.ReadBytes
	wb := mp.data.endDiskStat.WriteBytes - mp.data.startDiskStat.WriteBytes

	rc := mp.data.endDiskStat.ReadCount - mp.data.startDiskStat.ReadCount
	wc := mp.data.endDiskStat.WriteCount - mp.data.startDiskStat.WriteCount

	return &diskFrame{
		readBytes:  rb,
		writeBytes: wb,
		readCount:  rc,
		writeCount: wc,
	}
}

func sampleDiskStat() *process.IOCountersStat {
	diskStat, err := proc.IOCounters()
	if err != nil {
		fmt.Println("Error sampling disk stat", err)
	}
	return diskStat
}
