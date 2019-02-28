package metric

import (
	"fmt"
	"runtime"

	uuid "github.com/satori/go.uuid"
)

func prepareHeapMetricsData(metric *metricPlugin, memStats *runtime.MemStats, base metricDataModel) metricDataModel {
	base.ID = uuid.NewV4().String()
	base.MetricName = heapMetric

	memPercent, err := proc.MemoryPercent()
	if err != nil {
		fmt.Println(err)
	}

	base.Metrics = map[string]interface{}{
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
		memoryPercent: memPercent,
	}

	return base
}
