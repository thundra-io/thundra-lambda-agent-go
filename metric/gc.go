package metric

import (
	"runtime"

	"github.com/thundra-io/thundra-lambda-agent-go/plugin"
)

func prepareGCMetricsData(metric *metric, memStats *runtime.MemStats) metricData {
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
		MetricName:      gcMetric,
		MetricTimestamp: metric.span.metricTimestamp,

		Metrics: map[string]interface{}{
			// PauseTotalNs is the cumulative nanoseconds in GC
			// stop-the-world pauses since the program started.
			pauseTotalNs: memStats.PauseTotalNs,
			// PauseNs is recent GC stop-the-world pause time in nanoseconds.
			pauseNs: memStats.PauseNs[(memStats.NumGC+255)%256],
			// NumGC is the number of completed GC cycles.
			numGc: memStats.NumGC,
			// NextGC is the target heap size of the next GC cycle.
			nextGc: memStats.NextGC,
			// GCCPUFraction is the fraction of this program's available
			// CPU time used by the GC since the program started.
			gcCpuFraction: memStats.GCCPUFraction,
			//DeltaNumGc is the change in NUMGC from before execution to after execution.
			deltaNumGc: metric.span.endGCCount - metric.span.startGCCount,
			//DeltaPauseTotalNs is pause total change from before execution to after execution.
			deltaPauseTotalNs: metric.span.endPauseTotalNs - metric.span.startPauseTotalNs,
		},
		Tags: map[string]interface{}{},
	}
}
