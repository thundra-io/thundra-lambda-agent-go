package metric

import (
	"github.com/thundra-io/thundra-lambda-agent-go/plugin"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/shirou/gopsutil/process"
	"github.com/shirou/gopsutil/net"
)

type Builder interface {
	EnableGCStats() Builder
	EnableHeapStats() Builder
	EnableGoroutineStats() Builder
	EnableCPUStats() Builder
	EnableDiskStats() Builder
	EnableNetStats() Builder
	Build() *metric
}

type builder struct {
	enableGCStats        bool
	enableHeapStats      bool
	enableGoroutineStats bool
	enableCPUStats       bool
	enableDiskStats      bool
	enableNetStats       bool
	prevDiskStat         *process.IOCountersStat
	prevNetStat          *net.IOCountersStat
	process              *process.Process
}

func (b *builder) EnableGCStats() Builder {
	b.enableGCStats = true
	return b
}

func (b *builder) EnableHeapStats() Builder {
	b.enableHeapStats = true
	return b
}

func (b *builder) EnableGoroutineStats() Builder {
	b.enableGoroutineStats = true
	return b
}

func (b *builder) EnableCPUStats() Builder {
	b.enableCPUStats = true
	return b
}

func (b *builder) EnableDiskStats() Builder {
	b.enableDiskStats = true
	return b
}

func (b *builder) EnableNetStats() Builder {
	b.enableNetStats = true
	return b
}

func (b *builder) Build() *metric {
	//Initialize with empty objects
	if b.enableDiskStats {
		b.prevDiskStat = &process.IOCountersStat{}
	}

	if b.enableNetStats {
		b.prevNetStat = &net.IOCountersStat{}
	}

	if b.enableCPUStats || b.enableDiskStats {
		b.process = plugin.GetThisProcess()
	}

	return &metric{
		statData: statData{
			applicationName:    plugin.GetApplicationName(),
			applicationId:      plugin.GetAppIdFromStreamName(lambdacontext.LogStreamName),
			applicationVersion: plugin.GetApplicationVersion(),
			applicationProfile: plugin.GetApplicationProfile(),
			applicationType:    plugin.GetApplicationType(),
		},

		prevDiskStat: b.prevDiskStat,
		prevNetStat:  b.prevNetStat,
		process:      b.process,

		enableGCStats:        b.enableGCStats,
		enableHeapStats:      b.enableHeapStats,
		enableGoroutineStats: b.enableGoroutineStats,
		enableCPUStats:       b.enableCPUStats,
		enableDiskStats:      b.enableDiskStats,
		enableNetStats:       b.enableNetStats,
	}
}

func NewBuilder() Builder {
	return &builder{}
}
