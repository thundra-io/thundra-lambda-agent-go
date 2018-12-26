package metric

import (
	"github.com/thundra-io/thundra-lambda-agent-go/plugin"
	"runtime"
)

func prepareGoRoutineMetricsData(metric *metric) metricData {
	return metricData{
		Id:                        plugin.GenerateNewID(),
		Type:                      metricType,
		AgentVersion:              plugin.AgentVersion,
		DataModelVersion:          plugin.DataModelVersion,
		ApplicationId:             plugin.ApplicationID,
		ApplicationDomainName:     plugin.ApplicationDomainName,
		ApplicationClassName:      plugin.ApplicationClassName,
		ApplicationName:           plugin.FunctionName,
		ApplicationVersion:        plugin.ApplicationVersion,
		ApplicationStage:          plugin.ApplicationStage,
		ApplicationRuntime:        plugin.ApplicationRuntime,
		ApplicationRuntimeVersion: plugin.ApplicationRuntimeVersion,
		ApplicationTags:           map[string]interface{}{},

		TraceId:         plugin.TraceID,
		TransactionId:  plugin.TransactionID,
		SpanId:          plugin.SpanID,
		MetricName:      goroutineMetric,
		MetricTimestamp: metric.span.metricTimestamp,

		Metrics: map[string]interface{}{
			// NumGoroutine is the number of goroutines on execution
			numGoroutine: uint64(runtime.NumGoroutine()),
		},
		Tags: map[string]interface{}{},
	}
}
