package trace

import (
	"context"
	"sync"
	"time"
	"os"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/satori/go.uuid"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/thundra-io/thundra-lambda-agent-go/plugin"
)

type Trace struct {
	startTime          time.Time
	endTime            time.Time
	duration           time.Duration
	errors             []string
	thrownError        interface{}
	thrownErrorMessage interface{}
	panicInfo          *panicInfo
	errorInfo          *errorInfo
}

var invocationCount uint32 = 0
var uniqueId uuid.UUID

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
	StartTime          string                 `json:"startTime"`
	EndTime            string                 `json:"endTime"`
	Duration           int64                  `json:"duration"`
	Errors             []string               `json:"errors"`
	ThrownError        interface{}            `json:"thrownError"`
	ThrownErrorMessage interface{}            `json:"thrownErrorMessage"`
	AuditInfo          map[string]interface{} `json:"auditInfo"`
	Properties         map[string]interface{} `json:"properties"`
}

func (trace *Trace) BeforeExecution(ctx context.Context, request json.RawMessage, wg *sync.WaitGroup) {
	trace.startTime = time.Now().Round(time.Millisecond)
	cleanBuffer(trace)
	wg.Done()
}

func (trace *Trace) AfterExecution(ctx context.Context, request json.RawMessage, response interface{}, err interface{}) ([]interface{}, string) {
	trace.endTime = time.Now().Round(time.Millisecond)
	trace.duration = trace.endTime.Sub(trace.startTime)

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
	trace.endTime = time.Now()
	trace.duration = trace.endTime.Sub(trace.startTime)

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
	uniqueId = uuid.Must(uuid.NewV4())

	appId := plugin.SplitAppId(lambdacontext.LogStreamName)
	ver := lambdacontext.FunctionVersion

	profile := os.Getenv(plugin.ThundraApplicationProfile)
	if profile == "" {
		profile = plugin.DefaultProfile
	}

	props := prepareProperties(request, response)
	ai := prepareAuditInfo(trace)

	return traceData{
		Id:                 uniqueId.String(),
		ApplicationName:    lambdacontext.FunctionName,
		ApplicationId:      appId,
		ApplicationVersion: ver,
		ApplicationProfile: profile,
		ApplicationType:    plugin.ApplicationType,
		ContextId:          uniqueId.String(),
		ContextName:        lambdacontext.FunctionName,
		ContextType:        executionContext,
		StartTime:          trace.startTime.Format(plugin.TimeFormat),
		EndTime:            trace.endTime.Format(plugin.TimeFormat),
		Duration:           convertToMsec(trace.duration), //Convert it to msec
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
		auditInfoContextName: lambdacontext.FunctionName,
		auditInfoId:          uniqueId,
		auditInfoOpenTime:    trace.startTime.Format(plugin.TimeFormat),
		auditInfoCloseTime:   trace.endTime.Format(plugin.TimeFormat),
		auditInfoErrors:      auditErrors,
		auditInfoThrownError: auditThrownError,
		//"thrownErrorMessage": trace.thrownErrorMessage,
	}
}

func convertToMsec(duration time.Duration) int64 {
	return int64(duration / time.Millisecond)
}
