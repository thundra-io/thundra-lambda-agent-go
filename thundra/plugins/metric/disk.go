package metric

import (
	"fmt"

	"github.com/thundra-io/thundra-lambda-agent-go/plugin"
	"github.com/shirou/gopsutil/process"
)

func prepareDiskMetricsData(metric *metric) metricData {
	df := takeDiskFrame(metric)
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
		MetricName:      diskMetric,
		MetricTimestamp: metric.span.metricTimestamp,

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
func takeDiskFrame(metric *metric) *diskFrame {
	if metric.span.endDiskStat == nil || metric.span.startDiskStat == nil {
		return &diskFrame{}
	}
	rb := metric.span.endDiskStat.ReadBytes - metric.span.startDiskStat.ReadBytes
	wb := metric.span.endDiskStat.WriteBytes - metric.span.startDiskStat.WriteBytes

	rc := metric.span.endDiskStat.ReadCount - metric.span.startDiskStat.ReadCount
	wc := metric.span.endDiskStat.WriteCount - metric.span.startDiskStat.WriteCount

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
