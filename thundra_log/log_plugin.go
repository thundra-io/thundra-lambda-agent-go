package thundra_log

import (
	"context"
	"encoding/json"
	"sync"
)

type logPlugin struct{}

func New() *logPlugin {
	return &logPlugin{}
}

func (p *logPlugin) BeforeExecution(ctx context.Context, request json.RawMessage, wg *sync.WaitGroup) {
	logManager.clearLogs()
	wg.Done()
}

func (p *logPlugin) AfterExecution(ctx context.Context, request json.RawMessage, response interface{}, err interface{}) ([]interface{}, string) {
	var collectedData []interface{}
	for _, l := range logManager.logs {
		data := prepareLogData(l)
		collectedData = append(collectedData, data)
	}
	return collectedData, logDataType
}

func (p *logPlugin) OnPanic(ctx context.Context, request json.RawMessage, err interface{}, stackTrace []byte) ([]interface{}, string) {
	var collectedData []interface{}
	for _, l := range logManager.logs {
		data := prepareLogData(l)
		collectedData = append(collectedData, data)
	}
	return collectedData, logDataType
}
