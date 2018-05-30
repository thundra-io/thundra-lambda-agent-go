package metric

import (
	"context"
	"encoding/json"
	"fmt"
	"runtime"
	"sync"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/net"
	"github.com/shirou/gopsutil/process"
	"github.com/thundra-io/thundra-lambda-agent-go/plugin"
)

type metric struct {
	statTimestamp     int64
	startGCCount      uint32
	endGCCount        uint32
	startPauseTotalNs uint64
	endPauseTotalNs   uint64
	process           *process.Process
	processCpuPercent float64
	systemCpuPercent  float64
	currDiskStat      *process.IOCountersStat
	prevDiskStat      *process.IOCountersStat
	currNetStat       *net.IOCountersStat
	prevNetStat       *net.IOCountersStat

	disableGCStats        bool
	disableHeapStats      bool
	disableGoroutineStats bool
	disableCPUStats       bool
	disableDiskStats      bool
	disableNetStats       bool
}

func (metric *metric) BeforeExecution(ctx context.Context, request json.RawMessage, wg *sync.WaitGroup) {
	metric.statTimestamp = plugin.GetTimestamp()

	if !metric.disableGCStats {
		m := &runtime.MemStats{}
		runtime.ReadMemStats(m)

		metric.startGCCount = m.NumGC
		metric.startPauseTotalNs = m.PauseTotalNs
	}

	if !metric.disableCPUStats {
		// We need to calculate process and system percentages here to register cpu times
		// Later we'll use them to calculate cpu usage percentage
		metric.process.Percent(0)
		cpu.Percent(0, false)
	}

	wg.Done()
}

func (metric *metric) AfterExecution(ctx context.Context, request json.RawMessage, response interface{}, err interface{}) ([]interface{}, string) {
	mStats := &runtime.MemStats{}
	runtime.ReadMemStats(mStats)

	var stats []interface{}

	if !metric.disableHeapStats {
		h := prepareHeapStatsData(metric, mStats)
		stats = append(stats, h)
	}

	if !metric.disableGCStats {
		metric.endGCCount = mStats.NumGC
		metric.endPauseTotalNs = mStats.PauseTotalNs

		gc := prepareGCStatsData(metric, mStats)
		stats = append(stats, gc)
	}

	if !metric.disableGoroutineStats {
		g := prepareGoRoutineStatsData(metric)
		stats = append(stats, g)
	}

	if !metric.disableCPUStats {
		p, s, err := getCPUUsagePercentage(metric.process)
		if err != nil {
			fmt.Println(err)
		} else {
			metric.processCpuPercent = p
			metric.systemCpuPercent = s
			c := prepareCPUStatsData(metric)
			stats = append(stats, c)
		}
	}

	if !metric.disableDiskStats {
		diskStat, err := metric.process.IOCounters()
		if err != nil {
			fmt.Println(err)
		} else {
			metric.currDiskStat = diskStat
			d := prepareDiskStatsData(metric)
			stats = append(stats, d)
		}
	}

	if !metric.disableNetStats {
		netIOStat, err := net.IOCounters(false)
		if err != nil {
			fmt.Println(err)
		} else {
			metric.currNetStat = &netIOStat[all]
			n := prepareNetStatsData(metric)
			stats = append(stats, n)
		}
	}

	return stats, statDataType
}

//OnPanic just collect the metrics and send them as in the AfterExecution
func (metric *metric) OnPanic(ctx context.Context, request json.RawMessage, err interface{}, stackTrace []byte) ([]interface{}, string) {
	return metric.AfterExecution(ctx, request, nil, err)
}
