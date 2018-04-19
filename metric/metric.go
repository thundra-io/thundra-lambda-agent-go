package metric

import (
	"context"
	"encoding/json"
	"sync"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/thundra-io/thundra-lambda-agent-go/plugin"
	"runtime"

	"github.com/shirou/gopsutil/process"
	"fmt"
	"github.com/shirou/gopsutil/net"
	"github.com/shirou/gopsutil/cpu"
)

const StatDataType = "StatData"

type Metric struct {
	statData
	statTimestamp int64
	//TODO separate as GC objects, process objects etc.
	startGCCount      uint32
	endGCCount        uint32
	startPauseTotalNs uint64
	endPauseTotalNs   uint64
	process           *process.Process
	cpuPercent        float64
	currDiskStat      *process.IOCountersStat
	prevDiskStat      *process.IOCountersStat
	currNetStat       *net.IOCountersStat
	prevNetStat       *net.IOCountersStat
	currTimeStat      *cpu.TimesStat
	prevTimeStat      *cpu.TimesStat

	EnableGCStats        bool
	EnableHeapStats      bool
	EnableGoroutineStats bool
	EnableCPUStats       bool
	EnableDiskStats      bool
	EnableNetStats       bool
}

type statData struct {
	applicationName    string
	applicationId      string
	applicationVersion string
	applicationProfile string
	applicationType    string
}

func NewMetric() *Metric {

	return &Metric{
		statData: statData{
			applicationName:    plugin.GetApplicationName(),
			applicationId:      plugin.GetAppIdFromStreamName(lambdacontext.LogStreamName),
			applicationVersion: plugin.GetApplicationVersion(),
			applicationProfile: plugin.GetApplicationProfile(),
			applicationType:    plugin.GetApplicationType(),
		},

		//Initialize with empty objects
		prevDiskStat: &process.IOCountersStat{},
		prevNetStat:  &net.IOCountersStat{},
		prevTimeStat: &cpu.TimesStat{},
	}
}

func (metric *Metric) BeforeExecution(ctx context.Context, request json.RawMessage, wg *sync.WaitGroup) {
	metric.statTimestamp = plugin.GetTimestamp()

	if metric.EnableGCStats {
		m := &runtime.MemStats{}
		runtime.ReadMemStats(m)

		metric.startGCCount = m.NumGC
		metric.startPauseTotalNs = m.PauseTotalNs
	}

	if metric.EnableCPUStats || metric.EnableDiskStats {
		metric.process = plugin.GetThisProcess()
	}

	wg.Done()
}

func (metric *Metric) AfterExecution(ctx context.Context, request json.RawMessage, response interface{}, err interface{}) ([]interface{}, string) {
	mStats := &runtime.MemStats{}
	runtime.ReadMemStats(mStats)

	var stats []interface{}

	if metric.EnableHeapStats {
		h := prepareHeapStatsData(metric, mStats)
		stats = append(stats, h)
	}

	if metric.EnableGCStats {
		metric.endGCCount = mStats.NumGC
		metric.endPauseTotalNs = mStats.PauseTotalNs

		gc := prepareGCStatsData(metric, mStats)
		stats = append(stats, gc)
	}

	if metric.EnableGoroutineStats {
		g := prepareGoRoutineStatsData(metric)
		stats = append(stats, g)
	}

	if metric.EnableCPUStats {
		p, err := getCPUUsagePercentage(metric.process)
		if err != nil {
			fmt.Println(err)
		} else {
			metric.cpuPercent = p
			c := prepareCPUStatsData(metric)
			stats = append(stats, c)
		}
	}

	if metric.EnableDiskStats {
		diskStat, err := metric.process.IOCounters()
		if err != nil {
			fmt.Println(err)
		} else {
			metric.currDiskStat = diskStat
			d := prepareDiskStatsData(metric)
			stats = append(stats, d)
		}
	}

	if metric.EnableNetStats {
		netIOStat, err := net.IOCounters(false)
		if err != nil {
			fmt.Println(err)
		} else {
			metric.currNetStat = &netIOStat[ALL]
			n := prepareNetStatsData(metric)
			stats = append(stats, n)
		}
	}

	return stats, StatDataType
}

//OnPanic just collect the metrics and send them as in the AfterExecution
func (metric *Metric) OnPanic(ctx context.Context, request json.RawMessage, err interface{}, stackTrace []byte) ([]interface{}, string) {
	return metric.AfterExecution(ctx, request, nil, err)
}
