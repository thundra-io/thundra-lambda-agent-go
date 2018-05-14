package trace

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"sync"

	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/thundra-io/thundra-lambda-agent-go/plugin"
)

type Trace struct {
	startTime          int64
	endTime            int64
	duration           int64
	errors             []string
	thrownError        interface{}
	thrownErrorMessage interface{}
	panicInfo          *panicInfo
	errorInfo          *errorInfo
}

var invocationCount uint32 = 0
var uniqueId string

type traceData struct {
	Id                 string                 `json:"id"`
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

// NewTrace returns a new trace object.
func NewTrace() *Trace {
	return &Trace{}
}

func (trace *Trace) BeforeExecution(ctx context.Context, request json.RawMessage, wg *sync.WaitGroup) {
	trace.startTime = plugin.GetTimestamp()
	cleanBuffer(trace)
	wg.Done()
}

func (trace *Trace) AfterExecution(ctx context.Context, request json.RawMessage, response interface{}, err interface{}) ([]interface{}, string) {
	trace.endTime = plugin.GetTimestamp()
	trace.duration = trace.endTime - trace.startTime

	if err != nil {
		errMessage := getErrorMessage(err)
		errType := getErrorType(err)

		ei := &errorInfo{
			errMessage,
			errType,
		}

		trace.errorInfo = ei
		trace.thrownError = errType
		trace.thrownErrorMessage = errMessage
		trace.errors = append(trace.errors, errType)
	}

	td := prepareTraceData(request, response, err, trace)
	var traceArr []interface{}
	traceArr = append(traceArr, td)
	return traceArr, TraceDataType
}

func (trace *Trace) OnPanic(ctx context.Context, request json.RawMessage, err interface{}, stackTrace []byte) ([]interface{}, string) {
	trace.endTime = plugin.GetTimestamp()
	trace.duration = trace.endTime - trace.startTime

	errMessage := err.(error).Error()
	errType := getErrorType(err)
	pi := &panicInfo{
		errMessage,
		string(stackTrace),
		errType,
	}

	trace.panicInfo = pi
	trace.thrownError = errType
	trace.thrownErrorMessage = getErrorMessage(err)
	trace.errors = append(trace.errors, errType)

	td := prepareTraceData(request, nil, nil, trace)
	var traceArr []interface{}
	traceArr = append(traceArr, td)
	return traceArr, TraceDataType
}

func cleanBuffer(trace *Trace) {
	trace.errors = nil
}

func getErrorType(err interface{}) string {
	errorType := reflect.TypeOf(err)
	if errorType.Kind() == reflect.Ptr {
		return errorType.Elem().Name()
	}
	return errorType.Name()
}

func getErrorMessage(err interface{}) string {
	return err.(error).Error()
}

func prepareTraceData(request json.RawMessage, response interface{}, err interface{}, trace *Trace) traceData {
	uniqueId = plugin.GenerateNewId()
	props := prepareProperties(request, response)
	ai := prepareAuditInfo(trace)

	return traceData{
		Id:                 uniqueId,
		ApplicationName:    plugin.GetApplicationName(),
		ApplicationId:      plugin.GetAppIdFromStreamName(lambdacontext.LogStreamName),
		ApplicationVersion: plugin.GetApplicationVersion(),
		ApplicationProfile: plugin.GetApplicationProfile(),
		ApplicationType:    plugin.GetApplicationType(),
		ContextId:          uniqueId,
		ContextName:        plugin.GetApplicationName(),
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
		auditInfoPropertiesFunctionRegion:      os.Getenv(awsDefaultRegion),
		auditInfoPropertiesFunctionMemoryLimit: lambdacontext.MemoryLimitInMB,
	}
}

func prepareAuditInfo(trace *Trace) map[string]interface{} {
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
		auditInfoContextName:    lambdacontext.FunctionName,
		auditInfoId:             uniqueId,
		auditInfoOpenTimestamp:  trace.startTime,
		auditInfoCloseTimestamp: trace.endTime,
		auditInfoErrors:         auditErrors,
		auditInfoThrownError:    auditThrownError,
		//"thrownErrorMessage": trace.thrownErrorMessage,
	}
}
