package metric

import (
	"math"

	"github.com/thundra-io/thundra-lambda-agent-go/plugin"
)

type cpuStatsData struct {
	Id                 string `json:"id"`
	TransactionId      string `json:"transactionId"`
	ApplicationName    string `json:"applicationName"`
	ApplicationId      string `json:"applicationId"`
	ApplicationVersion string `json:"applicationVersion"`
	ApplicationProfile string `json:"applicationProfile"`
	ApplicationType    string `json:"applicationType"`
	StatName           string `json:"statName"`
	StatTimestamp      int64  `json:"statTimestamp"`

	// ProcessCPUPercent is the pid usage of the total CPU time
	ProcessCPUPercent float64 `json:"procPercent"`

	// SystemCPUPercent is the system usage of the total CPU time
	SystemCPUPercent float64 `json:"sysPercent"`
}

func prepareCPUStatsData(metric *metric) cpuStatsData {
	return cpuStatsData{
		Id:                 plugin.GenerateNewId(),
		TransactionId:      plugin.TransactionId,
		ApplicationName:    plugin.ApplicationName,
		ApplicationId:      plugin.ApplicationId,
		ApplicationVersion: plugin.ApplicationVersion,
		ApplicationProfile: plugin.ApplicationProfile,
		ApplicationType:    plugin.ApplicationType,
		StatName:           cpuStat,
		StatTimestamp:      metric.span.statTimestamp,
		ProcessCPUPercent:  metric.span.processCpuPercent,
		SystemCPUPercent:   metric.span.systemCpuPercent,
	}
}

func getSystemUsagePercent(metric *metric) float64 {
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

func getProcessUsagePercent(metric *metric) float64 {
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
