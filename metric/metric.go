package metric

import (
	"context"
	"encoding/json"
	"runtime"

	"github.com/thundra-io/thundra-lambda-agent-go/v2/config"

	"github.com/shirou/gopsutil/net"
	"github.com/shirou/gopsutil/process"
	"github.com/thundra-io/thundra-lambda-agent-go/v2/plugin"
	"github.com/thundra-io/thundra-lambda-agent-go/v2/utils"
)

var proc *process.Process
var pid string

type metricPlugin struct {
	data                    *metricData
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

// New returns new metric plugin initialized with empty metrics data
func New() *metricPlugin {
	pid = utils.GetPid()
	proc = utils.GetThisProcess()

	return &metricPlugin{
		data: &metricData{},
	}
}

func (mp *metricPlugin) IsEnabled() bool {
	return !config.MetricDisabled
}

func (mp *metricPlugin) Order() uint8 {
	return pluginOrder
}

func (mp *metricPlugin) BeforeExecution(ctx context.Context, request json.RawMessage) context.Context {
	mp.data = &metricData{}
	mp.data.metricTimestamp = utils.GetTimestamp()

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

	return ctx
}

func (mp *metricPlugin) AfterExecution(ctx context.Context, request json.RawMessage, response interface{}, err interface{}) ([]plugin.MonitoringDataWrapper, context.Context) {
	mStats := &runtime.MemStats{}
	runtime.ReadMemStats(mStats)

	var stats []plugin.MonitoringDataWrapper

	base := mp.prepareMetricsData()

	if GetSampler() != nil {
		if !GetSampler().IsSampled(base) {
			return stats, ctx
		}
	}

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

	return stats, ctx
}
