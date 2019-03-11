package log

import (
	"context"
	"encoding/json"

	"github.com/thundra-io/thundra-lambda-agent-go/config"

	"github.com/thundra-io/thundra-lambda-agent-go/plugin"
)

type logPlugin struct{}

// New creates and returns new logPlugin
func New() *logPlugin {
	return &logPlugin{}
}

func (p *logPlugin) IsEnabled() bool {
	return !config.LogDisabled
}

func (p *logPlugin) Order() uint8 {
	return pluginOrder
}

func (p *logPlugin) BeforeExecution(ctx context.Context, request json.RawMessage) context.Context {
	logManager.clearLogs()
	return ctx
}

func (p *logPlugin) AfterExecution(ctx context.Context, request json.RawMessage, response interface{}, err interface{}) ([]plugin.MonitoringDataWrapper, context.Context) {
	var collectedData []plugin.MonitoringDataWrapper
	for _, l := range logManager.logs {
		data := prepareLogData(l)
		collectedData = append(collectedData, plugin.WrapMonitoringData(data, logType))
	}
	return collectedData, ctx
}

func (p *logPlugin) OnPanic(ctx context.Context, request json.RawMessage, err interface{}, stackTrace []byte) []plugin.MonitoringDataWrapper {
	var collectedData []plugin.MonitoringDataWrapper
	for _, l := range logManager.logs {
		data := prepareLogData(l)
		collectedData = append(collectedData, plugin.WrapMonitoringData(data, logType))
	}
	return collectedData
}
