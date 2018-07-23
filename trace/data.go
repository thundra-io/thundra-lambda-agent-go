package trace

import (
	"encoding/json"
	"github.com/thundra-io/thundra-lambda-agent-go/plugin"
	"context"
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

func (tr *trace) prepareTraceData(ctx context.Context, request json.RawMessage, response interface{}) traceData {
	props := tr.prepareProperties(ctx, request, response)
	ai := tr.prepareAuditInfo()

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
		StartTimestamp:     tr.startTime,
		EndTimestamp:       tr.endTime,
		Duration:           tr.duration,
		Errors:             tr.errors,
		ThrownError:        tr.thrownError,
		ThrownErrorMessage: tr.thrownErrorMessage,
		AuditInfo:          ai,
		Properties:         props,
	}
}

func (tr *trace) prepareProperties(ctx context.Context, request json.RawMessage, response interface{}) map[string]interface{} {
	coldStart := "true"
	if invocationCount += 1; invocationCount != 1 {
		coldStart = "false"
	}
	// If the agent's user doesn't want to send their request and response data, hide them.
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
		auditInfoPropertiesLogGroupName:        plugin.LogGroupName,
		auditInfoPropertiesLogStreamName:       plugin.LogStreamName,
		auditInfoPropertiesFunctionRegion:      plugin.Region,
		auditInfoPropertiesFunctionMemoryLimit: plugin.MemorySize,
		auditInfoPropertiesFunctionARN:         plugin.GetInvokedFunctionArn(ctx),
		auditInfoPropertiesRequestId:           plugin.GetAwsRequestID(ctx),
		auditInfoPropertiesTimeout:             tr.timeout,
	}
}

func (tr *trace) prepareAuditInfo() map[string]interface{} {
	var auditErrors []interface{}
	var auditThrownError interface{}

	if tr.panicInfo != nil {
		p := *tr.panicInfo
		auditErrors = append(auditErrors, p)
		auditThrownError = p
	} else if tr.errorInfo != nil {
		e := *tr.errorInfo
		auditErrors = append(auditErrors, e)
		auditThrownError = e
	}

	return map[string]interface{}{
		auditInfoContextName:    plugin.ApplicationName,
		auditInfoId:             plugin.ContextId,
		auditInfoOpenTimestamp:  tr.startTime,
		auditInfoCloseTimestamp: tr.endTime,
		auditInfoErrors:         auditErrors,
		auditInfoThrownError:    auditThrownError,
	}
}
