package metric

import (
	"context"
	"encoding/json"
	"sync"
	"time"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/thundra-io/thundra-lambda-agent-go/plugin"
	"runtime"

	"github.com/shirou/gopsutil/process"
	"fmt"
	"github.com/shirou/gopsutil/net"
)

const StatDataType = "StatData"

//const selfStatFile = "/proc/1/stat"

type Metric struct {
	statData
	statTime          time.Time
	startGCCount      uint32
	endGCCount        uint32
	startPauseTotalNs uint64
	endPauseTotalNs   uint64
	process           *process.Process
	cpuPercent        float64
	prevIOStat        *process.IOCountersStat
	currIOStat        *process.IOCountersStat
	prevNetIOStat     *net.IOCountersStat
	currNetIOStat     *net.IOCountersStat

	EnableGCStats        bool
	EnableHeapStats      bool
	EnableGoroutineStats bool
	EnableCPUStats       bool
	EnableIOStats        bool
	EnableNetworkIOStats bool
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
			applicationId:      plugin.GetAppId(lambdacontext.LogStreamName),
			applicationVersion: plugin.GetApplicationVersion(),
			applicationProfile: plugin.GetApplicationProfile(),
			applicationType:    plugin.GetApplicationType(),
		},

		//Initialize with empty objects
		prevIOStat:    &process.IOCountersStat{},
		prevNetIOStat: &net.IOCountersStat{},
	}
}

func (metric *Metric) BeforeExecution(ctx context.Context, request json.RawMessage, wg *sync.WaitGroup) {
	m := &runtime.MemStats{}
	runtime.ReadMemStats(m)

	if metric.EnableGCStats {
		metric.startGCCount = m.NumGC
		metric.startPauseTotalNs = m.PauseTotalNs
	}

	if metric.EnableCPUStats || metric.EnableIOStats {
		metric.process = plugin.GetThisProcess()
	}

	wg.Done()
}

func (metric *Metric) AfterExecution(ctx context.Context, request json.RawMessage, response interface{}, err interface{}) ([]interface{}, string) {
	mStats := &runtime.MemStats{}
	runtime.ReadMemStats(mStats)

	metric.statTime = time.Now().Round(time.Millisecond)

	var stats []interface{}
	metric.statTime = time.Now().Round(time.Millisecond)

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
			fmt.Println("CpuPercent: ", p)
			metric.cpuPercent = p
			c := prepareCPUStatsData(metric)
			stats = append(stats, c)
		}
	}

	if metric.EnableIOStats {
		ioStat, err := metric.process.IOCounters()
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("ioStat: ", ioStat)
			metric.currIOStat = ioStat
			io := prepareIOStatsData(metric)
			stats = append(stats, io)
		}
	}

	if metric.EnableNetworkIOStats {
		netIOStat, err := net.IOCounters(false)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("netIOSTat: ", netIOStat)
			metric.currNetIOStat = &netIOStat[ALL]
			n := prepareNetIOStatsData(metric)
			stats = append(stats, n)
		}
	}

	return stats, StatDataType
}

//OnPanic just collect the metrics and send them as in the AfterExecution
func (metric *Metric) OnPanic(ctx context.Context, request json.RawMessage, err interface{}, stackTrace []byte) ([]interface{}, string) {
	return metric.AfterExecution(ctx, request, nil, err)
}
