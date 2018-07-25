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
	statTimestamp     int64
	startGCCount      uint32
	endGCCount        uint32
	startPauseTotalNs uint64
	endPauseTotalNs   uint64
	startCPUTimeStat  *cpuTimesStat
	endCPUTimeStat    *cpuTimesStat
	processCpuPercent float64
	systemCpuPercent  float64
	endDiskStat       *process.IOCountersStat
	startDiskStat     *process.IOCountersStat
	endNetStat        *net.IOCountersStat
	startNetStat      *net.IOCountersStat
}

func (metric *metric) BeforeExecution(ctx context.Context, request json.RawMessage, wg *sync.WaitGroup) {
	metric.span = new(metricSpan)
	metric.span.statTimestamp = plugin.GetTimestamp()

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

func (metric *metric) AfterExecution(ctx context.Context, request json.RawMessage, response interface{}, err interface{}) ([]interface{}, string) {
	mStats := &runtime.MemStats{}
	runtime.ReadMemStats(mStats)

	var stats []interface{}

	if !metric.disableGCStats {
		metric.span.endGCCount = mStats.NumGC
		metric.span.endPauseTotalNs = mStats.PauseTotalNs

		gc := prepareGCStatsData(metric, mStats)
		stats = append(stats, gc)
	}

	if !metric.disableHeapStats {
		h := prepareHeapStatsData(metric, mStats)
		stats = append(stats, h)
	}

	if !metric.disableGoroutineStats {
		g := prepareGoRoutineStatsData(metric)
		stats = append(stats, g)
	}

	if !metric.disableCPUStats {
		metric.span.endCPUTimeStat = sampleCPUtimesStat()

		metric.span.processCpuPercent = getProcessUsagePercent(metric)
		metric.span.systemCpuPercent = getSystemUsagePercent(metric)

		c := prepareCPUStatsData(metric)
		stats = append(stats, c)
	}

	if !metric.disableDiskStats {
		metric.span.endDiskStat = sampleDiskStat()
		d := prepareDiskStatsData(metric)
		stats = append(stats, d)
	}

	if !metric.disableNetStats {
		metric.span.endNetStat = sampleNetStat()
		n := prepareNetStatsData(metric)
		stats = append(stats, n)
	}
	metric.span = nil
	return stats, statDataType
}

//OnPanic just collect the metrics and send them as in the AfterExecution
func (metric *metric) OnPanic(ctx context.Context, request json.RawMessage, err interface{}, stackTrace []byte) ([]interface{}, string) {
	return metric.AfterExecution(ctx, request, nil, err)
}
