package metric

import (
	"context"
	"encoding/json"
	"runtime"
	"sync"

	"github.com/shirou/gopsutil/net"
	"github.com/shirou/gopsutil/process"
	"github.com/thundra-io/thundra-lambda-agent-go/plugin"
)

var proc *process.Process

type metric struct {
	span *metricSpan

	disableGCMetrics        bool
	disableHeapMetrics      bool
	disableGoroutineMetrics bool
	disableCPUMetrics       bool
	disableDiskMetrics      bool
	disableNetMetrics       bool
	disableMemoryMetrics    bool
}

// metricSpan collects information related to metric plugin per invocation.
type metricSpan struct {
	metricTimestamp   int64
	startGCCount      uint32
	endGCCount        uint32
	startPauseTotalNs uint64
	endPauseTotalNs   uint64
	startCPUTimeStat  *cpuTimesStat
	endCPUTimeStat    *cpuTimesStat
	appCpuLoad        float64
	systemCpuLoad     float64
	endDiskStat       *process.IOCountersStat
	startDiskStat     *process.IOCountersStat
	endNetStat        *net.IOCountersStat
	startNetStat      *net.IOCountersStat
}

func (metric *metric) BeforeExecution(ctx context.Context, request json.RawMessage, wg *sync.WaitGroup) {
	metric.span = new(metricSpan)
	metric.span.metricTimestamp = plugin.GetTimestamp()

	if !metric.disableGCMetrics {
		m := &runtime.MemStats{}
		runtime.ReadMemStats(m)

		metric.span.startGCCount = m.NumGC
		metric.span.startPauseTotalNs = m.PauseTotalNs
	}

	if !metric.disableCPUMetrics {
		metric.span.startCPUTimeStat = sampleCPUtimesStat()
	}

	if !metric.disableDiskMetrics {
		metric.span.startDiskStat = sampleDiskStat()
	}

	if !metric.disableNetMetrics {
		metric.span.startNetStat = sampleNetStat()
	}

	wg.Done()
}

func (metric *metric) AfterExecution(ctx context.Context, request json.RawMessage, response interface{}, err interface{}) []plugin.MonitoringDataWrapper {
	mStats := &runtime.MemStats{}
	runtime.ReadMemStats(mStats)

	var stats []plugin.MonitoringDataWrapper

	if !metric.disableGCMetrics {
		metric.span.endGCCount = mStats.NumGC
		metric.span.endPauseTotalNs = mStats.PauseTotalNs

		gc := prepareGCMetricsData(metric, mStats)
		stats = append(stats, plugin.WrapMonitoringData(gc, metricType))
	}

	if !metric.disableHeapMetrics {
		h := prepareHeapMetricsData(metric, mStats)
		stats = append(stats, plugin.WrapMonitoringData(h, metricType))
	}

	if !metric.disableGoroutineMetrics {
		g := prepareGoRoutineMetricsData(metric)
		stats = append(stats, plugin.WrapMonitoringData(g, metricType))
	}

	if !metric.disableCPUMetrics {
		metric.span.endCPUTimeStat = sampleCPUtimesStat()

		metric.span.appCpuLoad = getProcessCPULoad(metric)
		metric.span.systemCpuLoad = getSystemCPULoad(metric)

		c := prepareCPUMetricsData(metric)
		stats = append(stats, plugin.WrapMonitoringData(c, metricType))
	}

	if !metric.disableDiskMetrics {
		metric.span.endDiskStat = sampleDiskStat()
		d := prepareDiskMetricsData(metric)
		stats = append(stats, plugin.WrapMonitoringData(d, metricType))
	}

	if !metric.disableNetMetrics {
		metric.span.endNetStat = sampleNetStat()
		n := prepareNetMetricsData(metric)
		stats = append(stats, plugin.WrapMonitoringData(n, metricType))
	}

	if !metric.disableMemoryMetrics {
		mm := prepareMemoryMetricsData(metric)
		stats = append(stats, plugin.WrapMonitoringData(mm, metricType))
	}
	metric.span = nil
	return stats
}

//OnPanic just collect the metrics and send them as in the AfterExecution
func (metric *metric) OnPanic(ctx context.Context, request json.RawMessage, err interface{}, stackTrace []byte) []plugin.MonitoringDataWrapper {
	return metric.AfterExecution(ctx, request, nil, err)
}
