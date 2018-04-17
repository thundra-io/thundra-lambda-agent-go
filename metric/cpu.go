package metric

import (
	"github.com/thundra-io/thundra-lambda-agent-go/plugin"
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
	CPUPercent         float64 `json:"cpuPercent"`
}

func prepareCPUStatsData(metric *Metric) cpuStatsData {
	return cpuStatsData{
		Id:                 plugin.GenerateNewId(),
		ApplicationName:    metric.applicationName,
		ApplicationId:      metric.applicationId,
		ApplicationVersion: metric.applicationVersion,
		ApplicationProfile: metric.applicationProfile,
		ApplicationType:    plugin.ApplicationType,
		StatName:           cpuStat,
		StatTimestamp:      metric.statTimestamp,
		CPUPercent:         metric.cpuPercent,
	}
}

func getCPUUsagePercentage(p *process.Process) (float64, error) {
	percentage, err := p.CPUPercent()
	if err != nil {
		return 0, err
	}

	return percentage, nil
}
