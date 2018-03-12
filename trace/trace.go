package trace

import (
	"context"
	"sync"
	"time"
	"github.com/satori/go.uuid"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"os"
	"encoding/json"
	"strings"
	"fmt"
	"thundra-agent-go/plugin"
	"reflect"
)

type trace struct {
	startTime          time.Time
	endTime            time.Time
	duration           time.Duration
	errors             []string
	thrownError        interface{}
	thrownErrorMessage interface{}
	panicInfo          *PanicInfo
	errorInfo          *ErrorInfo
}

type TraceFactory struct{}

type PanicInfo struct {
	ErrMessage string `json:"errorMessage"`
	StackTrace string `json:"error"`
	ErrType    string `json:"errorType"`
}

type ErrorInfo struct {
	ErrMessage string `json:"errorMessage"`
	ErrType    string `json:"errorType"`
}

func (t *TraceFactory) Create() plugin.Plugin {
	return &trace{}
}

var invocationCount uint32 = 0
var uniqueId uuid.UUID

type TraceData struct {
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

func (trace *trace) BeforeExecution(ctx context.Context, request interface{}, wg *sync.WaitGroup) {
	trace.startTime = time.Now().Round(time.Millisecond)
	cleanBuffer(trace)
	wg.Done()
}
func cleanBuffer(trace *trace) {
	trace.errors = nil
}

func (trace *trace) AfterExecution(ctx context.Context, request interface{}, response interface{}, err interface{}, wg *sync.WaitGroup) (interface{}, string) {
	defer wg.Done()
	trace.endTime = time.Now().Round(time.Millisecond)
	trace.duration = trace.endTime.Sub(trace.startTime)

	if err != nil {
		errMessage := getErrorMessage(err)
		errType := getErrorType(err)

		ei := &ErrorInfo{
			errMessage,
			errType,
		}

		trace.errorInfo = ei
		trace.thrownError = errType
		trace.thrownErrorMessage = errMessage
		trace.errors = append(trace.errors, errType)
	}

	td := prepareReport(request, response, err, trace)
	return td, TraceDataType
}

func (trace *trace) OnPanic(ctx context.Context, request json.RawMessage, err interface{}, stackTrace []byte, wg *sync.WaitGroup) (interface{}, string) {
	defer wg.Done()
	trace.endTime = time.Now()
	trace.duration = trace.endTime.Sub(trace.startTime)

	errMessage := err.(error).Error()
	errType := getErrorType(err)
	pi := &PanicInfo{
		errMessage,
		string(stackTrace),
		errType,
	}

	trace.panicInfo = pi
	trace.thrownError = errType
	trace.thrownErrorMessage = getErrorMessage(err)
	trace.errors = append(trace.errors, errType)

	td := prepareReport(request, nil, nil, trace)
	return td, TraceDataType
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

func prepareReport(request interface{}, response interface{}, err interface{}, trace *trace) interface{} {
	uniqueId = uuid.Must(uuid.NewV4())

	props := prepareProperties(request, response)
	ai := prepareAuditInfo(trace)
	td := prepareTraceData(trace, err, props, ai)
	return td
}

func prepareProperties(request interface{}, response interface{}) map[string]interface{} {
	coldStart := "true"
	if invocationCount += 1; invocationCount != 1 {
		coldStart = "false"
	}
	return map[string]interface{}{
		AuditInfoPropertiesRequest:             request,
		AuditInfoPropertiesResponse:            response,
		AuditInfoPropertiesColdStart:           coldStart,
		AuditInfoPropertiesFunctionRegion:      os.Getenv(AwsDefaultRegion),
		AuditInfoPropertiesFunctionMemoryLimit: lambdacontext.MemoryLimitInMB,
	}
}

func prepareAuditInfo(trace *trace) map[string]interface{} {
	var auditErrors []interface{}
	var auditThrownError interface{}

	if trace.panicInfo != nil {
		fmt.Println("PANIC NOT NULL")
		p := *trace.panicInfo
		auditErrors = append(auditErrors, p)
		auditThrownError = p
	} else if trace.errorInfo != nil {
		fmt.Println("ERROR NOT NULL")
		e := *trace.errorInfo
		auditErrors = append(auditErrors, e)
		auditThrownError = e
	}

	return map[string]interface{}{
		AuditInfoContextName: lambdacontext.FunctionName,
		AuditInfoId:          uniqueId,
		AuditInfoOpenTime:    trace.startTime.Format(TimeFormat),
		AuditInfoCloseTime:   trace.endTime.Format(TimeFormat),
		AuditInfoErrors:      auditErrors,
		AuditInfoThrownError: auditThrownError,
		//"thrownErrorMessage": trace.thrownErrorMessage,
	}
}

func prepareTraceData(trace *trace, err interface{}, props map[string]interface{}, auditInfo map[string]interface{}) TraceData {
	appId := splitAppId(lambdacontext.LogStreamName)
	ver := lambdacontext.FunctionVersion

	profile := os.Getenv(ThundraApplicationProfile)
	if profile == "" {
		profile = DefaultProfile
	}

	return TraceData{
		uniqueId.String(),
		lambdacontext.FunctionName,
		appId,
		ver,
		profile,
		ApplicationType,
		uniqueId.String(),
		lambdacontext.FunctionName,
		ExecutionContext,
		trace.startTime.Format(TimeFormat),
		trace.endTime.Format(TimeFormat),
		convertToMsec(trace.duration), //Convert it to msec
		trace.errors,
		trace.thrownError,
		trace.thrownErrorMessage,
		auditInfo,
		props,
	}
}

func splitAppId(logStreamName string) string {
	s := strings.Split(logStreamName, "]")
	if len(s) > 1 {
		return s[1]
	} else {
		return ""
	}
}

func convertToMsec(duration time.Duration) int64 {
	return int64(duration / time.Millisecond)
}
