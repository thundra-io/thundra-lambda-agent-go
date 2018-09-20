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

	disableGCStats        bool
	disableHeapStats      bool
	disableGoroutineStats bool
	disableCPUStats       bool
	disableDiskStats      bool
	disableNetStats       bool
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

	if !metric.disableGCStats {
		m := &runtime.MemStats{}
		runtime.ReadMemStats(m)

		metric.span.startGCCount = m.NumGC
		metric.span.startPauseTotalNs = m.PauseTotalNs
	}

	if !metric.disableCPUStats {
		metric.span.startCPUTimeStat = sampleCPUtimesStat()
	}

	if !metric.disableDiskStats {
		metric.span.startDiskStat = sampleDiskStat()
	}

	if !metric.disableNetStats {
		metric.span.startNetStat = sampleNetStat()
	}

	wg.Done()
}

func (metric *metric) AfterExecution(ctx context.Context, request json.RawMessage, response interface{}, err interface{}) []plugin.MonitoringDataWrapper {
	mStats := &runtime.MemStats{}
	runtime.ReadMemStats(mStats)

	var stats []plugin.MonitoringDataWrapper

	if !metric.disableGCStats {
		metric.span.endGCCount = mStats.NumGC
		metric.span.endPauseTotalNs = mStats.PauseTotalNs

		gc := prepareGCMetricsData(metric, mStats)
		stats = append(stats, plugin.WrapMonitoringData(gc, metricType))
	}

	if !metric.disableHeapStats {
		h := prepareHeapStatsData(metric, mStats)
		stats = append(stats, plugin.WrapMonitoringData(h, metricType))
	}

	if !metric.disableGoroutineStats {
		g := prepareGoRoutineMetricsData(metric)
		stats = append(stats, plugin.WrapMonitoringData(g, metricType))
	}

	if !metric.disableCPUStats {
		metric.span.endCPUTimeStat = sampleCPUtimesStat()

		metric.span.appCpuLoad = getProcessCpuLoad(metric)
		metric.span.systemCpuLoad = getSystemCpuLoad(metric)

		c := prepareCpuMetricsData(metric)
		stats = append(stats, plugin.WrapMonitoringData(c, metricType))
	}

	if !metric.disableDiskStats {
		metric.span.endDiskStat = sampleDiskStat()
		d := prepareDiskMetricData(metric)
		stats = append(stats, plugin.WrapMonitoringData(d, metricType))
	}

	if !metric.disableNetStats {
		metric.span.endNetStat = sampleNetStat()
		n := prepareNetStatsData(metric)
		stats = append(stats, plugin.WrapMonitoringData(n, metricType))
	}
	metric.span = nil
	return stats
}

//OnPanic just collect the metrics and send them as in the AfterExecution
func (metric *metric) OnPanic(ctx context.Context, request json.RawMessage, err interface{}, stackTrace []byte) []plugin.MonitoringDataWrapper {
	return metric.AfterExecution(ctx, request, nil, err)
}
