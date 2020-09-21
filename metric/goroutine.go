package metric

import (
	"runtime"

	uuid "github.com/google/uuid"
)

func prepareGoRoutineMetricsData(mp *metricPlugin, base metricDataModel) metricDataModel {
	base.ID = uuid.New().String()
	base.MetricName = goroutineMetric
	base.Metrics = map[string]interface{}{
		// NumGoroutine is the number of goroutines on execution
		numGoroutine: uint64(runtime.NumGoroutine()),
	}

	return base
}
