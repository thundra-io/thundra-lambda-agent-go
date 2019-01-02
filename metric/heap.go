package metric

import (
	"fmt"
	"runtime"

	"github.com/thundra-io/thundra-lambda-agent-go/plugin"
)

func prepareHeapMetricsData(metric *metric, memStats *runtime.MemStats) metricData {
	mp, err := proc.MemoryPercent()
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
		MetricName:      heapMetric,
		MetricTimestamp: metric.span.metricTimestamp,

		Metrics: map[string]interface{}{
			// heapAlloc is bytes of allocated heap objects.
			//
			// "Allocated" heap objects include all reachable objects, as
			// well as unreachable objects that the garbage collector has
			// not yet freed.
			heapAlloc: memStats.HeapAlloc,
			// heapSys estimates the largest size the heap has had.
			heapSys: memStats.HeapSys,
			// heapInuse is bytes in in-use spans.
			// In-use spans have at least one object in them. These spans
			// can only be used for other objects of roughly the same
			// size.
			heapInuse: memStats.HeapInuse,
			// heapObjects is the number of allocated heap objects.
			heapObjects: memStats.HeapObjects,
			// memoryPercent returns how many percent of the total RAM this process uses
			memoryPercent: mp,
		},
		Tags: map[string]interface{}{},
	}
}
