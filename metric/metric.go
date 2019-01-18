package metric

import (
	"context"
	"encoding/json"
	"os"
	"runtime"
	"sync"

	"github.com/shirou/gopsutil/net"
	"github.com/shirou/gopsutil/process"
	"github.com/thundra-io/thundra-lambda-agent-go/plugin"
)

var proc *process.Process

type metricPlugin struct {
	data *metricData

	disableGCMetrics        bool
	disableHeapMetrics      bool
	disableGoroutineMetrics bool
	disableCPUMetrics       bool
	disableDiskMetrics      bool
	disableNetMetrics       bool
	disableMemoryMetrics    bool
}

type metricData struct {
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
}

func (mp *metricPlugin) IsEnabled() bool {
	if os.Getenv(plugin.ThundraDisableMetric) == "true" {
		return false
	}

	return true
}

func (mp *metricPlugin) BeforeExecution(ctx context.Context, request json.RawMessage, wg *sync.WaitGroup) {
	mp.data = new(metricData)
	mp.data.metricTimestamp = plugin.GetTimestamp()

	if !mp.disableGCMetrics {
		m := &runtime.MemStats{}
		runtime.ReadMemStats(m)

		mp.data.startGCCount = m.NumGC
		mp.data.startPauseTotalNs = m.PauseTotalNs
	}

	if !mp.disableCPUMetrics {
		mp.data.startCPUTimeStat = sampleCPUtimesStat()
	}

	if !mp.disableDiskMetrics {
		mp.data.startDiskStat = sampleDiskStat()
	}

	if !mp.disableNetMetrics {
		mp.data.startNetStat = sampleNetStat()
	}

	wg.Done()
}

func (mp *metricPlugin) AfterExecution(ctx context.Context, request json.RawMessage, response interface{}, err interface{}) []plugin.MonitoringDataWrapper {
	mStats := &runtime.MemStats{}
	runtime.ReadMemStats(mStats)

	var stats []plugin.MonitoringDataWrapper

	base := mp.prepareMetricsData()

	if !mp.disableGCMetrics {
		mp.data.endGCCount = mStats.NumGC
		mp.data.endPauseTotalNs = mStats.PauseTotalNs

		gc := prepareGCMetricsData(mp, mStats, base)
		stats = append(stats, plugin.WrapMonitoringData(gc, metricType))
	}

	if !mp.disableHeapMetrics {
		h := prepareHeapMetricsData(mp, mStats, base)
		stats = append(stats, plugin.WrapMonitoringData(h, metricType))
	}

	if !mp.disableGoroutineMetrics {
		g := prepareGoRoutineMetricsData(mp, base)
		stats = append(stats, plugin.WrapMonitoringData(g, metricType))
	}

	if !mp.disableCPUMetrics {
		mp.data.endCPUTimeStat = sampleCPUtimesStat()

		mp.data.appCPULoad = getProcessCPULoad(mp)
		mp.data.systemCPULoad = getSystemCPULoad(mp)

		c := prepareCPUMetricsData(mp, base)
		stats = append(stats, plugin.WrapMonitoringData(c, metricType))
	}

	if !mp.disableDiskMetrics {
		mp.data.endDiskStat = sampleDiskStat()
		d := prepareDiskMetricsData(mp, base)
		stats = append(stats, plugin.WrapMonitoringData(d, metricType))
	}

	if !mp.disableNetMetrics {
		mp.data.endNetStat = sampleNetStat()
		n := prepareNetMetricsData(mp, base)
		stats = append(stats, plugin.WrapMonitoringData(n, metricType))
	}

	if !mp.disableMemoryMetrics {
		mm := prepareMemoryMetricsData(mp, base)
		stats = append(stats, plugin.WrapMonitoringData(mm, metricType))
	}
	mp.data = nil
	return stats
}

//OnPanic just collect the metrics and send them as in the AfterExecution
func (mp *metricPlugin) OnPanic(ctx context.Context, request json.RawMessage, err interface{}, stackTrace []byte) []plugin.MonitoringDataWrapper {
	return mp.AfterExecution(ctx, request, nil, err)
}
