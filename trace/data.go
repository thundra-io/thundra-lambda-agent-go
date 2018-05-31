package trace

import (
	"encoding/json"
	"github.com/thundra-io/thundra-lambda-agent-go/plugin"
)

type traceData struct {
	Id                 string                 `json:"id"`
	TransactionId      string                 `json:"transactionId"`
	ApplicationName    string                 `json:"applicationName"`
	ApplicationId      string                 `json:"applicationId"`
	ApplicationVersion string                 `json:"applicationVersion"`
	ApplicationProfile string                 `json:"applicationProfile"`
	ApplicationType    string                 `json:"applicationType"`
	ContextId          string                 `json:"contextId"`
	ContextName        string                 `json:"contextName"`
	ContextType        string                 `json:"contextType"`
	StartTimestamp     int64                  `json:"startTimestamp"`
	EndTimestamp       int64                  `json:"endTimestamp"`
	Duration           int64                  `json:"duration"`
	Errors             []string               `json:"errors"`
	ThrownError        interface{}            `json:"thrownError"`
	ThrownErrorMessage interface{}            `json:"thrownErrorMessage"`
	AuditInfo          map[string]interface{} `json:"auditInfo"`
	Properties         map[string]interface{} `json:"properties"`
}

func prepareTraceData(request json.RawMessage, response interface{}, trace *trace) traceData {
	props := prepareProperties(request, response)
	ai := prepareAuditInfo(trace)

	return traceData{
		Id:                 plugin.ContextId,
		TransactionId:      plugin.TransactionId,
		ApplicationName:    plugin.ApplicationName,
		ApplicationId:      plugin.ApplicationId,
		ApplicationVersion: plugin.ApplicationVersion,
		ApplicationProfile: plugin.ApplicationProfile,
		ApplicationType:    plugin.ApplicationType,
		ContextId:          plugin.ContextId,
		ContextName:        plugin.ApplicationName,
		ContextType:        executionContext,
		StartTimestamp:     trace.startTime,
		EndTimestamp:       trace.endTime,
		Duration:           trace.duration,
		Errors:             trace.errors,
		ThrownError:        trace.thrownError,
		ThrownErrorMessage: trace.thrownErrorMessage,
		AuditInfo:          ai,
		Properties:         props,
	}
}

func prepareProperties(request json.RawMessage, response interface{}) map[string]interface{} {
	coldStart := "true"
	if invocationCount += 1; invocationCount != 1 {
		coldStart = "false"
	}
	if shouldHideRequest() {
		request = nil
	}
	if shouldHideResponse() {
		response = nil
	}
	return map[string]interface{}{
		auditInfoPropertiesRequest:             string(request),
		auditInfoPropertiesResponse:            response,
		auditInfoPropertiesColdStart:           coldStart,
		auditInfoPropertiesFunctionRegion:      plugin.Region,
		auditInfoPropertiesFunctionMemoryLimit: plugin.MemorySize,
	}
}

func prepareAuditInfo(trace *trace) map[string]interface{} {
	var auditErrors []interface{}
	var auditThrownError interface{}

	if trace.panicInfo != nil {
		p := *trace.panicInfo
		auditErrors = append(auditErrors, p)
		auditThrownError = p
	} else if trace.errorInfo != nil {
		e := *trace.errorInfo
		auditErrors = append(auditErrors, e)
		auditThrownError = e
	}

	return map[string]interface{}{
		auditInfoContextName:    plugin.ApplicationName,
		auditInfoId:             plugin.ContextId,
		auditInfoOpenTimestamp:  trace.startTime,
		auditInfoCloseTimestamp: trace.endTime,
		auditInfoErrors:         auditErrors,
		auditInfoThrownError:    auditThrownError,
	}
}
