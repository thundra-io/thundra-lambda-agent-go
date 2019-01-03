package metric

import (
	"github.com/thundra-io/thundra-lambda-agent-go/plugin"
)

type mBuilder interface {
	DisableGCMetrics() mBuilder
	DisableHeapMetrics() mBuilder
	DisableGoroutineMetrics() mBuilder
	DisableCPUMetrics() mBuilder
	DisableDiskMetrics() mBuilder
	DisableNetMetrics() mBuilder
	Build() *metricPlugin
}

var pid string

type builder struct {
	disableGCMetrics        bool
	disableHeapMetrics      bool
	disableGoroutineMetrics bool
	disableCPUMetrics       bool
	disableDiskMetrics      bool
	disableNetMetrics       bool
	disableMemoryMetrics    bool
}

// New initializes a new metric object which collects all types of metrics. If you want to disable a metric that
// you don't want to collect use NewBuilder() instead.
func New() *metricPlugin {
	pid = plugin.GetPid()
	return NewBuilder().Build()
}

// NewBuilder returns a builder that you can use to disable the metrics that you don't want to collect.
func NewBuilder() mBuilder {
	return &builder{}
}

// DisableGCMetrics disables gc metrics collection. Check gcMetricsData to see which metrics are collected.
func (b *builder) DisableGCMetrics() mBuilder {
	b.disableGCMetrics = true
	return b
}

// DisableHeapMetrics disables heap metrics collection. Check heapMetricsData to see which metrics are collected.
func (b *builder) DisableHeapMetrics() mBuilder {
	b.disableHeapMetrics = true
	return b
}

// DisableGoroutineMetrics disables goroutines metrics collection. Check goRoutineMetricsData to see which metrics are collected.
func (b *builder) DisableGoroutineMetrics() mBuilder {
	b.disableGoroutineMetrics = true
	return b
}

// DisableCPUMetrics disables cpu metrics collection. Check cpuMetricsData to see which metrics are collected.
func (b *builder) DisableCPUMetrics() mBuilder {
	b.disableCPUMetrics = true
	return b
}

// DisableDiskMetrics disables disk metrics collection. Check diskMetricsData to see which metrics are collected.
func (b *builder) DisableDiskMetrics() mBuilder {
	b.disableDiskMetrics = true
	return b
}

// DisableNetMetrics disables net metrics collection. Check netMetricsData to see which metrics are collected.
func (b *builder) DisableNetMetrics() mBuilder {
	b.disableNetMetrics = true
	return b
}

// DisableMemoryMetrics disables memory metrics collection. Check memoryMetricsData to see which metrics are collected.
func (b *builder) DisableMemoryMetrics() mBuilder {
	b.disableMemoryMetrics = true
	return b
}

// Builds and returns the metric plugin that you can pass to a thundra object while building it using AddPlugin().
func (b *builder) Build() *metricPlugin {
	proc = plugin.GetThisProcess()

	return &metricPlugin{
		data:                    new(metricData),
		disableGCMetrics:        b.disableGCMetrics,
		disableHeapMetrics:      b.disableHeapMetrics,
		disableGoroutineMetrics: b.disableGoroutineMetrics,
		disableCPUMetrics:       b.disableCPUMetrics,
		disableDiskMetrics:      b.disableDiskMetrics,
		disableNetMetrics:       b.disableNetMetrics,
		disableMemoryMetrics:    b.disableMemoryMetrics,
	}
}
