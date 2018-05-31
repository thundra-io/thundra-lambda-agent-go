package thundra_log

import (
	"github.com/thundra-io/thundra-lambda-agent-go/plugin"
)

type logData struct {
	Id                          string `json:"id"`
	TransactionId               string `json:"transactionId"`
	ApplicationName             string `json:"applicationName"`
	ApplicationId               string `json:"applicationId"`
	ApplicationVersion          string `json:"applicationVersion"`
	ApplicationProfile          string `json:"applicationProfile"`
	ApplicationType             string `json:"applicationType"`
	LogTimestamp                int64  `json:"logTimestamp"`
	Log                         string `json:"log"`
	LogMessage                  string `json:"logMessage"`
	LoggerName                  string `json:"loggerName"`
	LogLevel                    string `json:"logLevel"`
	LogLevelId                  int    `json:"logLevelId"`
	RootExecutionAuditContextId string `json:"rootExecutionAuditContextId"`
}

type monitoredLog struct {
	log          string
	logMessage   string
	loggerName   string
	logTimestamp int64
	logLevel     string
	logLevelId   int
}

func prepareLogData(log *monitoredLog) logData {
	return logData{
		Id:                          plugin.GenerateNewId(),
		TransactionId:               plugin.TransactionId,
		ApplicationName:             plugin.ApplicationName,
		ApplicationId:               plugin.ApplicationId,
		ApplicationVersion:          plugin.ApplicationVersion,
		ApplicationProfile:          plugin.ApplicationProfile,
		ApplicationType:             plugin.ApplicationType,
		Log:                         log.log,
		LogMessage:                  log.logMessage,
		LoggerName:                  log.loggerName,
		LogTimestamp:                log.logTimestamp,
		LogLevel:                    log.logLevel,
		LogLevelId:                  log.logLevelId,
		RootExecutionAuditContextId: plugin.ContextId,
	}
}
