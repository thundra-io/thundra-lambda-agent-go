package thundra_log

import (
	"context"
	"encoding/json"
	"os"

	"github.com/thundra-io/thundra-lambda-agent-go/plugin"
)

type logPlugin struct{}

func New() *logPlugin {
	return &logPlugin{}
}

func (p *logPlugin) IsEnabled() bool {
	if os.Getenv(plugin.ThundraDisableLog) == "true" {
		return false
	}

	return true
}

func (p *logPlugin) BeforeExecution(ctx context.Context, request json.RawMessage) context.Context {
	logManager.clearLogs()
	return ctx
}

func (p *logPlugin) AfterExecution(ctx context.Context, request json.RawMessage, response interface{}, err interface{}) []plugin.MonitoringDataWrapper {
	var collectedData []plugin.MonitoringDataWrapper
	for _, l := range logManager.logs {
		data := prepareLogData(l)
		collectedData = append(collectedData, plugin.WrapMonitoringData(data, logType))
	}
	return collectedData
}

func (p *logPlugin) OnPanic(ctx context.Context, request json.RawMessage, err interface{}, stackTrace []byte) []plugin.MonitoringDataWrapper {
	var collectedData []plugin.MonitoringDataWrapper
	for _, l := range logManager.logs {
		data := prepareLogData(l)
		collectedData = append(collectedData, plugin.WrapMonitoringData(data, logType))
	}
	return collectedData
}
