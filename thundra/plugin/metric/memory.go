package metric

import (
	"fmt"

	"github.com/thundra-io/thundra-lambda-agent-go/plugin"
	"github.com/shirou/gopsutil/mem"
)

func prepareMemoryMetricsData(metric *metric) metricData {
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		fmt.Println(err)
	}

	procMemInfo, err := proc.MemoryInfo()
	if err != nil {
		fmt.Println(err)
	}

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
		MetricName:      memoryMetric,
		MetricTimestamp: metric.span.metricTimestamp,

		Metrics: map[string]interface{}{
			appUsedMemory: procMemInfo.RSS,
			appMaxMemory:  plugin.MemoryLimit * 1024 * 1024,
			sysUsedMemory: memInfo.Used,
			sysMaxMemory:  memInfo.Total,
		},
		Tags: map[string]interface{}{},
	}
}
