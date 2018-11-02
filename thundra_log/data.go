package thundra_log

import (
	"github.com/thundra-io/thundra-lambda-agent-go/plugin"
)

type logData struct {
	//Base fields
	Id                        string                 `json:"id"`
	Type                      string                 `json:"type"`
	AgentVersion              string                 `json:"agentVersion"`
	DataModelVersion          string                 `json:"dataModelVersion"`
	ApplicationId             string                 `json:"applicationId"`
	ApplicationDomainName     string                 `json:"applicationDomainName"`
	ApplicationClassName      string                 `json:"applicationClassName"`
	ApplicationName           string                 `json:"applicationName"`
	ApplicationVersion        string                 `json:"applicationVersion"`
	ApplicationStage          string                 `json:"applicationStage"`
	ApplicationRuntime        string                 `json:"applicationRuntime"`
	ApplicationRuntimeVersion string                 `json:"applicationRuntimeVersion"`
	ApplicationTags           map[string]interface{} `json:"applicationTags"`

	TraceId        string                 `json:"traceId"`
	TransactionId string                 `json:"transactionId"`
	SpanId         string                 `json:"spanId"`
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
		Id:                        plugin.GenerateNewId(),
		Type:                      logType,
		AgentVersion:              plugin.AgentVersion,
		DataModelVersion:          plugin.DataModelVersion,
		ApplicationId:             plugin.ApplicationId,
		ApplicationDomainName:     plugin.ApplicationDomainName,
		ApplicationClassName:      plugin.ApplicationClassName,
		ApplicationName:           plugin.FunctionName,
		ApplicationVersion:        plugin.ApplicationVersion,
		ApplicationStage:          plugin.ApplicationStage,
		ApplicationRuntime:        plugin.ApplicationRuntime,
		ApplicationRuntimeVersion: plugin.ApplicationRuntimeVersion,
		ApplicationTags:           map[string]interface{}{},

		TraceId:        plugin.TraceId,
		TransactionId: plugin.TransactionId,
		SpanId:         plugin.SpanId,
		LogMessage:     log.logMessage,
		LogContextName: log.logContextName,
		LogTimestamp:   log.logTimestamp,
		LogLevel:       log.logLevel,
		LogLevelCode:   log.logLevelCode,
		Tags:           map[string]interface{}{},
	}
}
