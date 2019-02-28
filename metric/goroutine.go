package metric

import (
	"runtime"

	uuid "github.com/satori/go.uuid"
)

func prepareGoRoutineMetricsData(mp *metricPlugin, base metricDataModel) metricDataModel {
	base.ID = uuid.NewV4().String()
	base.MetricName = goroutineMetric
	base.Metrics = map[string]interface{}{
		// NumGoroutine is the number of goroutines on execution
		numGoroutine: uint64(runtime.NumGoroutine()),
	}

	return base
}
