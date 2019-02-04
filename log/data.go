package log

import (
	"github.com/thundra-io/thundra-lambda-agent-go/plugin"
)

type logData struct {
	//Base fields
	ID                        string                 `json:"id"`
	Type                      string                 `json:"type"`
	AgentVersion              string                 `json:"agentVersion"`
	DataModelVersion          string                 `json:"dataModelVersion"`
	ApplicationID             string                 `json:"applicationId"`
	ApplicationDomainName     string                 `json:"applicationDomainName"`
	ApplicationClassName      string                 `json:"applicationClassName"`
	ApplicationName           string                 `json:"applicationName"`
	ApplicationVersion        string                 `json:"applicationVersion"`
	ApplicationStage          string                 `json:"applicationStage"`
	ApplicationRuntime        string                 `json:"applicationRuntime"`
	ApplicationRuntimeVersion string                 `json:"applicationRuntimeVersion"`
	ApplicationTags           map[string]interface{} `json:"applicationTags"`

	TraceID        string                 `json:"traceId"`
	TransactionID string                 `json:"transactionId"`
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
}

func prepareLogData(log *monitoringLog) logData {
	return logData{
		ID:                        plugin.GenerateNewID(),
		Type:                      logType,
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

		TraceID:        plugin.TraceID,
		TransactionID: plugin.TransactionID,
		SpanID:         plugin.SpanID,
		LogMessage:     log.logMessage,
		LogContextName: log.logContextName,
		LogTimestamp:   log.logTimestamp,
		LogLevel:       log.logLevel,
		LogLevelCode:   log.logLevelCode,
		Tags:           map[string]interface{}{},
	}
}
