package log

import (
	"context"
	"encoding/json"
	"sync"
	"github.com/thundra-io/thundra-lambda-agent-go/plugin"
)

type logPlugin struct{}

func New() *logPlugin {
	return &logPlugin{}
}

func (p *logPlugin) BeforeExecution(ctx context.Context, request json.RawMessage, wg *sync.WaitGroup) {
	logManager.clearLogs()
	wg.Done()
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
