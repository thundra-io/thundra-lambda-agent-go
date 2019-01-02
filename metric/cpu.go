package metric

import (
	"math"

	"github.com/thundra-io/thundra-lambda-agent-go/plugin"
)

func prepareCPUMetricsData(mp *metricPlugin) metricData {
	return metricData{
		ID:                        plugin.GenerateNewID(),
		Type:                      metricType,
		AgentVersion:              plugin.AgentVersion,
		DataModelVersion:          plugin.DataModelVersion,
		ApplicationID:             plugin.ApplicationID,
		ApplicationDomainName:     plugin.ApplicationDomainName,
		ApplicationClassName:      plugin.ApplicationClassName,
		ApplicationName:           plugin.FunctionName,
		ApplicationVersion:        plugin.ApplicationVersion,
		ApplicationStage:          plugin.ApplicationStage,
		ApplicationRuntime:        plugin.ApplicationRuntime,
		ApplicationRuntimeVersion: plugin.ApplicationRuntimeVersion,
		ApplicationTags:           map[string]interface{}{},

		TraceID:         plugin.TraceID,
		TransactionID:  plugin.TransactionID,
		// SpanId:          "", // Optional
		MetricName:      cpuMetric,
		MetricTimestamp: mp.metricTimestamp,

		Metrics: map[string]interface{}{
			appCPULoad: mp.appCPULoad,
			sysCPULoad: mp.systemCPULoad,
		},
		Tags: map[string]interface{}{},
	}
}

func getSystemCPULoad(mp *metricPlugin) float64 {
	// Skip test
	if mp.startCPUTimeStat == nil {
		return 0
	}
	dSysUsed := mp.endCPUTimeStat.sysUsed() - mp.startCPUTimeStat.sysUsed()
	dTotal := mp.endCPUTimeStat.total() - mp.startCPUTimeStat.total()
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

func getProcessCPULoad(mp *metricPlugin) float64 {
	// Skip test
	if mp.startCPUTimeStat == nil {
		return 0
	}
	dProcUsed := mp.endCPUTimeStat.procUsed() - mp.startCPUTimeStat.procUsed()
	dTotal := mp.endCPUTimeStat.total() - mp.startCPUTimeStat.total()
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
