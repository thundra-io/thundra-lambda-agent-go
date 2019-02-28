package log

import (
	"github.com/thundra-io/thundra-lambda-agent-go/application"
	"github.com/thundra-io/thundra-lambda-agent-go/constants"
	"github.com/thundra-io/thundra-lambda-agent-go/plugin"
	"github.com/thundra-io/thundra-lambda-agent-go/utils"
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
	TraceID                   string                 `json:"traceId"`
	TransactionID             string                 `json:"transactionId"`
	SpanID                    string                 `json:"spanId"`
	LogMessage                string                 `json:"logMessage"`
	LogContextName            string                 `json:"logContextName"`
	LogTimestamp              int64                  `json:"logTimestamp"`
	LogLevel                  string                 `json:"logLevel"`
	LogLevelCode              int                    `json:"logLevelCode"`
	Tags                      map[string]interface{} `json:"tags"`
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
		ID:                        utils.GenerateNewID(),
		Type:                      logType,
		AgentVersion:              constants.AgentVersion,
		DataModelVersion:          constants.DataModelVersion,
		ApplicationID:             application.ApplicationID,
		ApplicationDomainName:     application.ApplicationDomainName,
		ApplicationClassName:      application.ApplicationClassName,
		ApplicationName:           application.ApplicationName,
		ApplicationVersion:        application.ApplicationVersion,
		ApplicationStage:          application.ApplicationStage,
		ApplicationRuntime:        application.ApplicationRuntime,
		ApplicationRuntimeVersion: application.ApplicationRuntimeVersion,
		ApplicationTags:           application.ApplicationTags,
		TraceID:                   plugin.TraceID,
		TransactionID:             plugin.TransactionID,
		SpanID:                    log.spanID,
		LogMessage:                log.logMessage,
		LogContextName:            log.logContextName,
		LogTimestamp:              log.logTimestamp,
		LogLevel:                  log.logLevel,
		LogLevelCode:              log.logLevelCode,
		Tags:                      map[string]interface{}{},
	}
}
