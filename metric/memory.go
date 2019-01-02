package metric

import (
	"fmt"

	"github.com/thundra-io/thundra-lambda-agent-go/plugin"
	"github.com/shirou/gopsutil/mem"
)

func prepareMemoryMetricsData(mp *metricPlugin) metricData {
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		fmt.Println(err)
	}

	procMemInfo, err := proc.MemoryInfo()
	if err != nil {
		fmt.Println(err)
	}

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
		MetricName:      memoryMetric,
		MetricTimestamp: mp.metricTimestamp,

		Metrics: map[string]interface{}{
			appUsedMemory: procMemInfo.RSS,
			appMaxMemory:  plugin.MemoryLimit * 1024 * 1024,
			sysUsedMemory: memInfo.Used,
			sysMaxMemory:  memInfo.Total,
		},
		Tags: map[string]interface{}{},
	}
}
