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

type metricPlugin struct {
	metricTimestamp   int64
	startGCCount      uint32
	endGCCount        uint32
	startPauseTotalNs uint64
	endPauseTotalNs   uint64
	startCPUTimeStat  *cpuTimesStat
	endCPUTimeStat    *cpuTimesStat
	appCPULoad        float64
	systemCPULoad     float64
	endDiskStat       *process.IOCountersStat
	startDiskStat     *process.IOCountersStat
	endNetStat        *net.IOCountersStat
	startNetStat      *net.IOCountersStat

	disableGCMetrics        bool
	disableHeapMetrics      bool
	disableGoroutineMetrics bool
	disableCPUMetrics       bool
	disableDiskMetrics      bool
	disableNetMetrics       bool
	disableMemoryMetrics    bool
}

func (mp *metricPlugin) BeforeExecution(ctx context.Context, request json.RawMessage, wg *sync.WaitGroup) {
	mp = new(metricPlugin)
	mp.metricTimestamp = plugin.GetTimestamp()

	if !mp.disableGCMetrics {
		m := &runtime.MemStats{}
		runtime.ReadMemStats(m)

		mp.startGCCount = m.NumGC
		mp.startPauseTotalNs = m.PauseTotalNs
	}

	if !mp.disableCPUMetrics {
		mp.startCPUTimeStat = sampleCPUtimesStat()
	}

	if !mp.disableDiskMetrics {
		mp.startDiskStat = sampleDiskStat()
	}

	if !mp.disableNetMetrics {
		mp.startNetStat = sampleNetStat()
	}

	wg.Done()
}

func (mp *metricPlugin) AfterExecution(ctx context.Context, request json.RawMessage, response interface{}, err interface{}) []plugin.MonitoringDataWrapper {
	mStats := &runtime.MemStats{}
	runtime.ReadMemStats(mStats)

	var stats []plugin.MonitoringDataWrapper

	if !mp.disableGCMetrics {
		mp.endGCCount = mStats.NumGC
		mp.endPauseTotalNs = mStats.PauseTotalNs

		gc := prepareGCMetricsData(mp, mStats)
		stats = append(stats, plugin.WrapMonitoringData(gc, metricType))
	}

	if !mp.disableHeapMetrics {
		h := prepareHeapMetricsData(mp, mStats)
		stats = append(stats, plugin.WrapMonitoringData(h, metricType))
	}

	if !mp.disableGoroutineMetrics {
		g := prepareGoRoutineMetricsData(mp)
		stats = append(stats, plugin.WrapMonitoringData(g, metricType))
	}

	if !mp.disableCPUMetrics {
		mp.endCPUTimeStat = sampleCPUtimesStat()

		mp.appCPULoad = getProcessCPULoad(mp)
		mp.systemCPULoad = getSystemCPULoad(mp)

		c := prepareCPUMetricsData(mp)
		stats = append(stats, plugin.WrapMonitoringData(c, metricType))
	}

	if !mp.disableDiskMetrics {
		mp.endDiskStat = sampleDiskStat()
		d := prepareDiskMetricsData(mp)
		stats = append(stats, plugin.WrapMonitoringData(d, metricType))
	}

	if !mp.disableNetMetrics {
		mp.endNetStat = sampleNetStat()
		n := prepareNetMetricsData(mp)
		stats = append(stats, plugin.WrapMonitoringData(n, metricType))
	}

	if !mp.disableMemoryMetrics {
		mm := prepareMemoryMetricsData(mp)
		stats = append(stats, plugin.WrapMonitoringData(mm, metricType))
	}
	mp = nil
	return stats
}

//OnPanic just collect the metrics and send them as in the AfterExecution
func (mp *metricPlugin) OnPanic(ctx context.Context, request json.RawMessage, err interface{}, stackTrace []byte) []plugin.MonitoringDataWrapper {
	return mp.AfterExecution(ctx, request, nil, err)
}
