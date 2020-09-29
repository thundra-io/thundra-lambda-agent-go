package metric

import (
	"math"

	uuid "github.com/google/uuid"
)

func prepareCPUMetricsData(mp *metricPlugin, base metricDataModel) metricDataModel {
	base.ID = uuid.New().String()
	base.MetricName = cpuMetric
	base.Metrics = map[string]interface{}{
		appCPULoad: mp.data.appCPULoad,
		sysCPULoad: mp.data.systemCPULoad,
	}

	return base
}

func getSystemCPULoad(mp *metricPlugin) float64 {
	// Skip test
	if mp.data.startCPUTimeStat == nil {
		return 0
	}
	dSysUsed := mp.data.endCPUTimeStat.sysUsed() - mp.data.startCPUTimeStat.sysUsed()
	dTotal := mp.data.endCPUTimeStat.total() - mp.data.startCPUTimeStat.total()
	s := float64(dSysUsed) / float64(dTotal)
	if s <= 0 {
		s = 0
	} else if s >= 1 {
		s = 1
	} else if math.IsNaN(s) {
		s = 0
	}
	return s
}

func getProcessCPULoad(mp *metricPlugin) float64 {
	// Skip test
	if mp.data.startCPUTimeStat == nil {
		return 0
	}
	dProcUsed := mp.data.endCPUTimeStat.procUsed() - mp.data.startCPUTimeStat.procUsed()
	dTotal := mp.data.endCPUTimeStat.total() - mp.data.startCPUTimeStat.total()
	p := float64(dProcUsed) / float64(dTotal)
	if p <= 0 {
		p = 0
	} else if p >= 1 {
		p = 1
	} else if math.IsNaN(p) {
		p = 0
	}
	return p
}
