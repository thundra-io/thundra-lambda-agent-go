package metric

import (
	"runtime"

	"github.com/thundra-io/thundra-lambda-agent-go/plugin"
)

func prepareGoRoutineMetricsData(mp *metricPlugin) metricDataModel {
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
		TransactionID:   plugin.TransactionID,
		SpanID:          plugin.SpanID,
		MetricName:      goroutineMetric,
		MetricTimestamp: mp.data.metricTimestamp,

		Metrics: map[string]interface{}{
			// NumGoroutine is the number of goroutines on execution
			numGoroutine: uint64(runtime.NumGoroutine()),
		},
		Tags: map[string]interface{}{},
	}
}
