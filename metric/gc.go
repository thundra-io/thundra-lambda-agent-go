package metric

import (
	"runtime"

	uuid "github.com/google/uuid"
)

func prepareGCMetricsData(mp *metricPlugin, memStats *runtime.MemStats, base metricDataModel) metricDataModel {
	base.ID = uuid.New().String()
	base.MetricName = gcMetric
	base.Metrics = map[string]interface{}{
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
		deltaNumGc: mp.data.endGCCount - mp.data.startGCCount,
		//DeltaPauseTotalNs is pause total change from before execution to after execution.
		deltaPauseTotalNs: mp.data.endPauseTotalNs - mp.data.startPauseTotalNs,
	}

	return base
}
