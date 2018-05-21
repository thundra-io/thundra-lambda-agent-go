package metric

import (
	"github.com/shirou/gopsutil/net"
	"github.com/shirou/gopsutil/process"
	"github.com/thundra-io/thundra-lambda-agent-go/plugin"
)

type mBuilder interface {
	DisableGCStats() mBuilder
	DisableHeapStats() mBuilder
	DisableGoroutineStats() mBuilder
	DisableCPUStats() mBuilder
	DisableDiskStats() mBuilder
	DisableNetStats() mBuilder
	Build() *metric
}

type builder struct {
	disableGCStats        bool
	disableHeapStats      bool
	disableGoroutineStats bool
	disableCPUStats       bool
	disableDiskStats      bool
	disableNetStats       bool
	prevDiskStat          *process.IOCountersStat
	prevNetStat           *net.IOCountersStat
	process               *process.Process
}

// DisableGCStats disables gc metrics collection. Check gcStatsData to see which metrics are collected.
func (b *builder) DisableGCStats() mBuilder {
	b.disableGCStats = true
	return b
}

// DisableHeapStats disables heap stats collection. Check heapStatsData to see which metrics are collected.
func (b *builder) DisableHeapStats() mBuilder {
	b.disableHeapStats = true
	return b
}

// DisableGoroutineStats disables heap stats collection. Check goRoutineStatsData to see which metrics are collected.
func (b *builder) DisableGoroutineStats() mBuilder {
	b.disableGoroutineStats = true
	return b
}

// DisableCPUStats disables cpu stats collection. Check cpuStatsData to see which metrics are collected.
func (b *builder) DisableCPUStats() mBuilder {
	b.disableCPUStats = true
	return b
}

// DisableDiskStats disables disk stats collection. Check diskStatsData to see which metrics are collected.
func (b *builder) DisableDiskStats() mBuilder {
	b.disableDiskStats = true
	return b
}

// DisableNetStats disables net stats collection. Check netStatsData to see which metrics are collected.
func (b *builder) DisableNetStats() mBuilder {
	b.disableNetStats = true
	return b
}

// Builds and returns the metric plugin that you can pass to a thundra object while building it using AddPlugin().
func (b *builder) Build() *metric {
	//Initialize with empty objects
	if !b.disableDiskStats {
		b.prevDiskStat = &process.IOCountersStat{}
	}

	if !b.disableNetStats {
		b.prevNetStat = &net.IOCountersStat{}
	}

	if !b.disableCPUStats || !b.disableDiskStats || !b.disableHeapStats {
		b.process = plugin.GetThisProcess()
	}

	return &metric{
		statData: statData{
			applicationName:    plugin.GetApplicationName(),
			applicationId:      plugin.GetAppId(),
			applicationVersion: plugin.GetApplicationVersion(),
			applicationProfile: plugin.GetApplicationProfile(),
			applicationType:    plugin.GetApplicationType(),
		},

		prevDiskStat: b.prevDiskStat,
		prevNetStat:  b.prevNetStat,
		process:      b.process,

		disableGCStats:        b.disableGCStats,
		disableHeapStats:      b.disableHeapStats,
		disableGoroutineStats: b.disableGoroutineStats,
		disableCPUStats:       b.disableCPUStats,
		disableDiskStats:      b.disableDiskStats,
		disableNetStats:       b.disableNetStats,
	}
}

// NewBuilder returns a new metric builder.
func NewBuilder() mBuilder {
	return &builder{}
}
