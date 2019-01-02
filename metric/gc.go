package metric

import (
	"runtime"

	"github.com/thundra-io/thundra-lambda-agent-go/plugin"
)

func prepareGCMetricsData(mp *metricPlugin, memStats *runtime.MemStats) metricData {
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
		SpanID:          plugin.SpanID,
		MetricName:      gcMetric,
		MetricTimestamp: mp.metricTimestamp,

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
			gcCPUFraction: memStats.GCCPUFraction,
			//DeltaNumGc is the change in NUMGC from before execution to after execution.
			deltaNumGc: mp.endGCCount - mp.startGCCount,
			//DeltaPauseTotalNs is pause total change from before execution to after execution.
			deltaPauseTotalNs: mp.endPauseTotalNs - mp.startPauseTotalNs,
		},
		Tags: map[string]interface{}{},
	}
}
