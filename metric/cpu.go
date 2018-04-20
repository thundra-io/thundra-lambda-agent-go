package metric

import (
	"github.com/thundra-io/thundra-lambda-agent-go/plugin"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/process"
)

type cpuStatsData struct {
	Id                 string  `json:"id"`
	ApplicationName    string  `json:"applicationName"`
	ApplicationId      string  `json:"applicationId"`
	ApplicationVersion string  `json:"applicationVersion"`
	ApplicationProfile string  `json:"applicationProfile"`
	ApplicationType    string  `json:"applicationType"`
	StatName           string  `json:"statName"`
	StatTimestamp      int64   `json:"statTimestamp"`
	ProcessCPUPercent  float64 `json:"procPercent"`
	SystemCPUPercent   float64 `json:"sysPercent"`
}

func prepareCPUStatsData(metric *metric) cpuStatsData {
	return cpuStatsData{
		Id:                 plugin.GenerateNewId(),
		ApplicationName:    metric.applicationName,
		ApplicationId:      metric.applicationId,
		ApplicationVersion: metric.applicationVersion,
		ApplicationProfile: metric.applicationProfile,
		ApplicationType:    plugin.ApplicationType,
		StatName:           cpuStat,
		StatTimestamp:      metric.statTimestamp,
		ProcessCPUPercent:  metric.processCpuPercent,
		SystemCPUPercent:   metric.systemCpuPercent,
	}
}

func getCPUUsagePercentage(p *process.Process) (float64, float64, error) {
	sysUsage, err := cpu.Percent(0, false)
	if err != nil {
		return 0, 0, err
	}

	processUsage, err := p.Percent(0)
	if err != nil {
		return 0, 0, err
	}

	return processUsage, sysUsage[0], nil
}
