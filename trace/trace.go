package trace

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/thundra-io/thundra-lambda-agent-go/plugin"
)

type trace struct {
	startTime          int64
	endTime            int64
	duration           int64
	errors             []string
	thrownError        interface{}
	thrownErrorMessage interface{}
	panicInfo          *panicInfo
	errorInfo          *errorInfo
}

var invocationCount uint32
var UniqueId string

// NewTrace returns a new trace object.
func NewTrace() *trace {
	return &trace{}
}

func (trace *trace) BeforeExecution(ctx context.Context, request json.RawMessage, wg *sync.WaitGroup) {
	cleanBuffer(trace)
	trace.startTime = plugin.GetTimestamp()
	UniqueId = plugin.GenerateNewId()
	wg.Done()
}

func (trace *trace) AfterExecution(ctx context.Context, request json.RawMessage, response interface{}, err interface{}) ([]interface{}, string) {
	trace.endTime = plugin.GetTimestamp()
	trace.duration = trace.endTime - trace.startTime

	if err != nil {
		errMessage := plugin.GetErrorMessage(err)
		errType := plugin.GetErrorType(err)

		ei := &errorInfo{
			errMessage,
			errType,
		}

		trace.errorInfo = ei
		trace.thrownError = errType
		trace.thrownErrorMessage = errMessage
		trace.errors = append(trace.errors, errType)
	}

	td := prepareTraceData(request, response, trace)
	var traceArr []interface{}
	traceArr = append(traceArr, td)
	return traceArr, traceDataType
}

func (trace *trace) OnPanic(ctx context.Context, request json.RawMessage, err interface{}, stackTrace []byte) ([]interface{}, string) {
	trace.endTime = plugin.GetTimestamp()
	trace.duration = trace.endTime - trace.startTime

	errMessage := plugin.GetErrorMessage(err)
	errType := plugin.GetErrorType(err)
	pi := &panicInfo{
		errMessage,
		string(stackTrace),
		errType,
	}

	trace.panicInfo = pi
	trace.thrownError = errType
	trace.thrownErrorMessage = plugin.GetErrorMessage(err)
	trace.errors = append(trace.errors, errType)

	td := prepareTraceData(request, nil, trace)
	var traceArr []interface{}
	traceArr = append(traceArr, td)
	return traceArr, traceDataType
}

func cleanBuffer(trace *trace) {
	trace.errors = nil
}

func prepareTraceData(request json.RawMessage, response interface{}, trace *trace) traceData {
	props := prepareProperties(request, response)
	ai := prepareAuditInfo(trace)

	return traceData{
		Id:                 UniqueId,
		TransactionId:      plugin.TransactionId,
		ApplicationName:    plugin.ApplicationName,
		ApplicationId:      plugin.ApplicationId,
		ApplicationVersion: plugin.ApplicationVersion,
		ApplicationProfile: plugin.ApplicationProfile,
		ApplicationType:    plugin.ApplicationType,
		ContextId:          UniqueId,
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
		fmt.Println("Panic is not null")
		p := *trace.panicInfo
		auditErrors = append(auditErrors, p)
		auditThrownError = p
	} else if trace.errorInfo != nil {
		fmt.Println("Error is not null")
		e := *trace.errorInfo
		auditErrors = append(auditErrors, e)
		auditThrownError = e
	}

	return map[string]interface{}{
		auditInfoContextName:    plugin.ApplicationName,
		auditInfoId:             UniqueId,
		auditInfoOpenTimestamp:  trace.startTime,
		auditInfoCloseTimestamp: trace.endTime,
		auditInfoErrors:         auditErrors,
		auditInfoThrownError:    auditThrownError,
		//"thrownErrorMessage": trace.thrownErrorMessage,
	}
}
