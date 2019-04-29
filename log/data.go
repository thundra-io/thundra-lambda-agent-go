package log

import (
	"github.com/thundra-io/thundra-lambda-agent-go/plugin"
	"github.com/thundra-io/thundra-lambda-agent-go/utils"
)

type logData struct {
	//Base fields
	plugin.BaseDataModel
	ID             string                 `json:"id"`
	Type           string                 `json:"type"`
	TraceID        string                 `json:"traceId"`
	TransactionID  string                 `json:"transactionId"`
	SpanID         string                 `json:"spanId"`
	LogMessage     string                 `json:"logMessage"`
	LogContextName string                 `json:"logContextName"`
	LogTimestamp   int64                  `json:"logTimestamp"`
	LogLevel       string                 `json:"logLevel"`
	LogLevelCode   int                    `json:"logLevelCode"`
	Tags           map[string]interface{} `json:"tags"`
}

type monitoringLog struct {
	logMessage     string
	logContextName string
	logTimestamp   int64
	logLevel       string
	logLevelCode   int
	spanID         string
}

func prepareLogData(log *monitoringLog) logData {
	return logData{
		BaseDataModel:  plugin.GetBaseData(),
		ID:             utils.GenerateNewID(),
		Type:           logType,
		TraceID:        plugin.TraceID,
		TransactionID:  plugin.TransactionID,
		SpanID:         log.spanID,
		LogMessage:     log.logMessage,
		LogContextName: log.logContextName,
		LogTimestamp:   log.logTimestamp,
		LogLevel:       log.logLevel,
		LogLevelCode:   log.logLevelCode,
		Tags:           map[string]interface{}{},
	}
}
