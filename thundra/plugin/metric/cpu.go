package metric

import (
	"math"

	"github.com/thundra-io/thundra-lambda-agent-go/plugin"
)

func prepareCpuMetricsData(metric *metric) metricData {
	return metricData{
		Id:                        plugin.GenerateNewId(),
		Type:                      metricType,
		AgentVersion:              plugin.AgentVersion,
		DataModelVersion:          plugin.DataModelVersion,
		ApplicationId:             plugin.ApplicationId,
		ApplicationDomainName:     plugin.ApplicationDomainName,
		ApplicationClassName:      plugin.ApplicationClassName,
		ApplicationName:           plugin.FunctionName,
		ApplicationVersion:        plugin.ApplicationVersion,
		ApplicationStage:          plugin.ApplicationStage,
		ApplicationRuntime:        plugin.ApplicationRuntime,
		ApplicationRuntimeVersion: plugin.ApplicationRuntimeVersion,
		ApplicationTags:           map[string]interface{}{},

		TraceId:         plugin.TraceId,
		TracnsactionId:  plugin.TransactionId,
		SpanId:          plugin.SpanId,
		MetricName:      cpuMetric,
		MetricTimestamp: metric.span.metricTimestamp,

		Metrics: map[string]interface{}{
			appCpuLoad: metric.span.appCpuLoad,
			sysCpuLoad: metric.span.systemCpuLoad,
		},
		Tags: map[string]interface{}{},
	}
}

func getSystemCpuLoad(metric *metric) float64 {
	// Skip test
	if metric.span.startCPUTimeStat == nil {
		return 0
	}
	dSysUsed := metric.span.endCPUTimeStat.sys_used() - metric.span.startCPUTimeStat.sys_used()
	dTotal := metric.span.endCPUTimeStat.total() - metric.span.startCPUTimeStat.total()
	s := float64(dSysUsed) / float64(dTotal)
	if s <= 0 {
		s = 0
	} else if s >= 1 {
		s = 1
	} else if math.IsNaN(s) {
		s = 0;
	}
	return s
}

func getProcessCpuLoad(metric *metric) float64 {
	// Skip test
	if metric.span.startCPUTimeStat == nil {
		return 0
	}
	dProcUsed := metric.span.endCPUTimeStat.proc_used() - metric.span.startCPUTimeStat.proc_used()
	dTotal := metric.span.endCPUTimeStat.total() - metric.span.startCPUTimeStat.total()
	p := float64(dProcUsed) / float64(dTotal)
	if p <= 0 {
		p = 0
	} else if p >= 1 {
		p = 1
	} else if math.IsNaN(p) {
		p = 0;
	}
	return p
}
